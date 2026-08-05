package transport

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher владеет одним confirm-каналом и допускает одну публикацию в полёте.
type publisher struct {
	mu      sync.Mutex
	channel *amqp.Channel
	returns chan amqp.Return
}

func (p *publisher) publish(ctx context.Context, connection *amqp.Connection, exchange, routingKey string, body []byte) (err error) {
	p.mu.Lock()
	defer func() {
		// После транспортной ошибки судьба publish неизвестна, поэтому канал нельзя безопасно переиспользовать; basic.return канал не повреждает.
		if err != nil && !errors.Is(err, ErrUnroutable) {
			p.closeChannel()
		}
		p.mu.Unlock()
	}()

	if p.channel == nil || p.channel.IsClosed() {
		p.channel, err = connection.Channel()
		if err != nil {
			return fmt.Errorf("%w: open publisher channel: %w", ErrUnavailable, err)
		}
		if err = p.channel.Confirm(false); err != nil {
			return fmt.Errorf("%w: enable publisher confirms: %w", ErrUnavailable, err)
		}
		p.returns = p.channel.NotifyReturn(make(chan amqp.Return, 1))
	}

	confirmation, err := p.channel.PublishWithDeferredConfirmWithContext(ctx, exchange, routingKey, true, false, amqp.Publishing{
		Body:         body,
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Timestamp:    time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("%w: publish message: %w", ErrUnavailable, err)
	}

	// RabbitMQ отправляет basic.return до publisher confirm, а одна публикация в полёте однозначно связывает return с текущим сообщением.
	acknowledged, err := confirmation.WaitContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: wait for publisher confirm: %w", ErrUnavailable, err)
	}
	if !acknowledged {
		return fmt.Errorf("%w: broker nacked published message", ErrUnavailable)
	}

	select {
	case returned := <-p.returns:
		return fmt.Errorf("%w: %s", ErrUnroutable, returned.ReplyText)
	default:
		return nil
	}
}

func (p *publisher) close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closeChannel()
}

func (p *publisher) closeChannel() {
	if p.channel != nil && !p.channel.IsClosed() {
		_ = p.channel.Close()
	}
	p.channel = nil
	p.returns = nil
}

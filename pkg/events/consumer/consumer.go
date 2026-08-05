package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/scrynaxx/tr-core/pkg/events/transport"
)

// Config описывает topology и конкурентность одной подписки.
type Config struct {
	Exchange       string
	Subscriber     string
	EventName      string
	Concurrency    uint
	RetryDelay     time.Duration
	ReconnectDelay time.Duration
	HandleDelivery func(queue, dlx string, amqp amqp.Delivery)
}

// Consumer управляет topology и AMQP-consumer-ами одной подписки.
type Consumer struct {
	transport      *transport.Transport
	queue          string
	consumer       string
	concurrency    uint
	reconnectDelay time.Duration
	handle         func(queue, dlx string, amqp amqp.Delivery)
	dlx            string
}

// New объявляет topology подписки и возвращает готовый consumer.
func New(ctx context.Context, amqpTransport *transport.Transport, config Config) (*Consumer, error) {
	queueName := fmt.Sprintf("%s-%s-%s", config.Exchange, config.Subscriber, config.EventName)
	dlxName := fmt.Sprintf("%s-dlx", queueName)
	retryQueueName := fmt.Sprintf("%s-dlq", queueName)
	deadQueueName := fmt.Sprintf("%s-dead", queueName)
	retryRoutingKey := fmt.Sprintf("%s-retry", queueName)

	// Ошибка handler-а ведёт сообщение из основной очереди через DLX в retry-очередь, а её TTL возвращает сообщение по отдельному routing key.

	channel, err := amqpTransport.Channel(ctx)
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}
	defer channel.Close()

	if err = channel.ExchangeDeclare(config.Exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	if err = channel.ExchangeDeclare(dlxName, "direct", true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("declare DLX: %w", err)
	}

	if _, err = channel.QueueDeclare(retryQueueName, true, false, false, false, amqp.Table{
		"x-message-ttl":             config.RetryDelay.Milliseconds(),
		"x-dead-letter-exchange":    config.Exchange,
		"x-dead-letter-routing-key": retryRoutingKey,
	}); err != nil {
		return nil, fmt.Errorf("declare retry queue: %w", err)
	}

	if err = channel.QueueBind(retryQueueName, "failed", dlxName, false, nil); err != nil {
		return nil, fmt.Errorf("bind retry queue: %w", err)
	}

	if _, err = channel.QueueDeclare(deadQueueName, true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("declare dead queue: %w", err)
	}

	if err = channel.QueueBind(deadQueueName, "dead", dlxName, false, nil); err != nil {
		return nil, fmt.Errorf("bind dead queue: %w", err)
	}

	if _, err = channel.QueueDeclare(queueName, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    dlxName,
		"x-dead-letter-routing-key": "failed",
	}); err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	if err = channel.QueueBind(queueName, config.EventName, config.Exchange, false, nil); err != nil {
		return nil, fmt.Errorf("bind queue: %w", err)
	}

	if err = channel.QueueBind(queueName, retryRoutingKey, config.Exchange, false, nil); err != nil {
		return nil, fmt.Errorf("bind retry route: %w", err)
	}

	return &Consumer{
		transport:      amqpTransport,
		queue:          queueName,
		consumer:       fmt.Sprintf("%s-consumer", queueName),
		concurrency:    config.Concurrency,
		reconnectDelay: config.ReconnectDelay,
		handle:         config.HandleDelivery,
		dlx:            dlxName,
	}, nil
}

// Run запускает настроенное число независимых AMQP-consumer-ов и ждёт их завершения.
func (c *Consumer) Run(ctx context.Context) {
	var workers sync.WaitGroup

	for worker := range c.concurrency {
		name := fmt.Sprintf("%s-%d", c.consumer, worker+1)
		workers.Go(func() {
			c.runWorker(ctx, name)
		})
	}

	workers.Wait()
}

func (c *Consumer) runWorker(ctx context.Context, name string) {
	for ctx.Err() == nil {
		channel, err := c.transport.Channel(ctx)
		if err != nil {
			return
		}

		// Каждый worker владеет своим каналом; после разрыва consumer пересоздаётся целиком вместе с QoS.
		err = c.consume(ctx, channel, name)
		_ = channel.Close()
		if err != nil && ctx.Err() == nil {
			slog.Warn("[event bus] consumer stopped, reconnecting",
				slog.String("queue", c.queue),
				slog.Any("error", err),
			)

			select {
			case <-ctx.Done():
				return
			case <-time.After(c.reconnectDelay):
			}
		}
	}
}

func (c *Consumer) consume(ctx context.Context, channel *amqp.Channel, name string) error {
	if err := channel.Qos(1, 0, false); err != nil {
		return fmt.Errorf("set qos: %w", err)
	}

	deliveries, err := channel.Consume(c.queue, name, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return channel.Cancel(name, false)
		case delivery, ok := <-deliveries:
			if !ok {
				return transport.ErrUnavailable
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}

			c.handle(c.queue, c.dlx, delivery)
		}
	}
}

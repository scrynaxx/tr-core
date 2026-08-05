package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime/debug"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/scrynaxx/tr-core/pkg/events/consumer"
)

// HandleFunc обрабатывает типизированное событие.
type HandleFunc[TData any] = func(ctx context.Context, event Event[TData]) error

// Интерфейс унифицирует запуск обработчиков событий с разными типами данных.
type subscriptionRunner interface {
	// Метод запускает подготовленный обработчик.
	run(bus *Bus)
}

// Подписка хранит параметры регистрации типизированного обработчика.
type subscription[TData any] struct {
	// Идентифицирует микросервис-подписчик и участвует в имени очереди.
	subscriber string

	// Описывает имя и тип данных события.
	descriptor Descriptor[TData]

	// Обрабатывает доставленное событие.
	handle HandleFunc[TData]

	replicas uint
	consumer *consumer.Consumer
}

func (s *subscription[TData]) register(b *Bus) error {
	replicas := s.replicas
	if replicas == 0 {
		replicas = 1
	}
	s.replicas = replicas

	registeredConsumer, err := consumer.New(b.runCtx, b.transport, consumer.Config{
		Exchange:       exchange,
		Subscriber:     s.subscriber,
		EventName:      s.descriptor.Name,
		Concurrency:    replicas,
		RetryDelay:     b.config.Subscriber.RetryDelay,
		ReconnectDelay: b.config.ConnectionRetryDelay,
		HandleDelivery: func(queue, dlx string, delivery amqp.Delivery) {
			handleDelivery(b, queue, dlx, s.descriptor, s.handle, delivery)
		},
	})
	if err != nil {
		return err
	}

	s.consumer = registeredConsumer

	return nil
}

func (s *subscription[TData]) run(b *Bus) {
	b.runWG.Go(func() {
		s.consumer.Run(b.runCtx)
	})

	slog.Info("[event bus] consumer registered",
		slog.String("event_name", s.descriptor.Name),
		slog.Int("replicas", int(s.replicas)),
	)
}

func handleDelivery[TData any](b *Bus, queueName, dlxName string, desc Descriptor[TData], handle HandleFunc[TData], delivery amqp.Delivery) {
	var event Event[TData]
	if err := json.Unmarshal(delivery.Body, &event); err != nil {
		slog.Error("[event bus] malformed event received",
			slog.String("queue", queueName),
			slog.Any("error", err),
		)

		// Исходное сообщение подтверждается только после confirmed-публикации в dead queue, иначе RabbitMQ продолжит retry-маршрут.
		if deadErr := b.moveToDead(dlxName, delivery.Body); deadErr != nil {
			slog.Error("[event bus] store malformed event in dead queue failed",
				slog.String("queue", queueName),
				slog.Any("error", deadErr),
			)
			if nackErr := delivery.Nack(false, false); nackErr != nil {
				slog.Error("[event bus] reject malformed event failed",
					slog.String("queue", queueName),
					slog.Any("error", nackErr),
				)
			}
			return
		}

		if ackErr := delivery.Ack(false); ackErr != nil {
			slog.Error("[event bus] acknowledge malformed dead event failed",
				slog.String("queue", queueName),
				slog.Any("error", ackErr),
			)
		}
		return
	}

	// Panic преобразуется в обычную ошибку обработки, чтобы доставка обязательно завершилась Ack или Nack и не заняла QoS-слот навсегда.
	handleErr := func() (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf("panic: %v", recovered)
				slog.Error("[event bus] consumer panic recovered",
					slog.String("queue", queueName),
					slog.String("event_id", event.EventID.String()),
					slog.Any("panic", recovered),
					slog.String("stack", string(debug.Stack())),
				)
			}
		}()
		return handle(b.runCtx, event)
	}()
	if handleErr == nil {
		if err := delivery.Ack(false); err != nil {
			slog.Error("[event bus] acknowledge event failed",
				slog.String("queue", queueName),
				slog.String("event_id", event.EventID.String()),
				slog.Any("error", err),
			)
		}
		return
	}

	attempt := deliveryAttempt(delivery.Headers, queueName)
	if attempt >= int64(b.config.Subscriber.MaxAttempts) {
		if err := b.moveToDead(dlxName, delivery.Body); err != nil {
			slog.Error("[event bus] store event in dead queue failed",
				slog.String("event_name", desc.Name),
				slog.String("queue", queueName),
				slog.String("event_id", event.EventID.String()),
				slog.Int64("attempt", attempt),
				slog.Any("error", err),
			)

			if nackErr := delivery.Nack(false, false); nackErr != nil {
				slog.Error("[event bus] reject event after dead queue failure failed",
					slog.String("queue", queueName),
					slog.String("event_id", event.EventID.String()),
					slog.Any("error", nackErr),
				)
			}
			return
		}

		slog.Error("[event bus] event moved to dead queue",
			slog.String("event_name", desc.Name),
			slog.String("queue", queueName),
			slog.String("event_id", event.EventID.String()),
			slog.Int64("attempt", attempt),
			slog.Any("error", handleErr),
		)
		if err := delivery.Ack(false); err != nil {
			slog.Error("[event bus] acknowledge dead event failed",
				slog.String("queue", queueName),
				slog.String("event_id", event.EventID.String()),
				slog.Any("error", err),
			)
		}
		return
	}

	slog.Error("[event bus] event handler failed",
		slog.String("event_name", desc.Name),
		slog.String("queue", queueName),
		slog.String("event_id", event.EventID.String()),
		slog.Int64("attempt", attempt),
		slog.Any("error", handleErr),
	)
	if err := delivery.Nack(false, false); err != nil {
		slog.Error("[event bus] reject event for retry failed",
			slog.String("queue", queueName),
			slog.String("event_id", event.EventID.String()),
			slog.Any("error", err),
		)
	}
}

func (b *Bus) moveToDead(dlxName string, body []byte) error {
	// Уже принятое сообщение успевает попасть в dead queue во время shutdown, но отдельный таймаут не позволяет задержать остановку бесконечно.
	ctx, cancel := context.WithTimeout(context.WithoutCancel(b.runCtx), b.config.Subscriber.DeadLetterTimeout)
	defer cancel()
	return b.transport.Publish(ctx, dlxName, "dead", body)
}

func deliveryAttempt(headers amqp.Table, queueName string) int64 {
	if headers == nil {
		return 1
	}

	deaths, ok := headers["x-death"]
	if !ok {
		return 1
	}

	items, ok := deaths.([]any)
	if !ok {
		return 1
	}

	// RabbitMQ хранит отдельную x-death запись для каждой очереди маршрута; число попыток берётся только из основной очереди subscriber-а.
	for _, item := range items {
		death, ok := item.(amqp.Table)
		if !ok {
			continue
		}

		queue, ok := death["queue"].(string)
		if !ok || queue != queueName {
			continue
		}

		switch count := death["count"].(type) {
		case int64:
			return count + 1
		case int32:
			return int64(count) + 1
		case int:
			return int64(count) + 1
		}
	}

	return 1
}

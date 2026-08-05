package events

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/pkg/events/transport"
)

const exchange = "events"

// Params содержит параметры подключения шины к RabbitMQ.
type Params struct {
	User     string
	Password string
	Address  string
	Vhost    string
}

// Bus объединяет RabbitMQ-транспорт и доставку сохранённых outbox-событий.
type Bus struct {
	// Контекст запрещает запуск новых consumer- и outbox-операций при остановке.
	runCtx    context.Context
	runCancel context.CancelFunc

	transport     *transport.Transport
	repository    outboxRepository
	subscriptions []subscriptionRunner
	config        Config
	started       bool

	// Идентификатор отличает экземпляр шины при распределённой аренде событий.
	workerID string

	runWG sync.WaitGroup
}

// New создаёт outbox-хранилище, открывает соединение с RabbitMQ и возвращает шину событий.
func New(ctx context.Context, params Params, pool *pgxpool.Pool, schemaName string, config *Config) (*Bus, error) {
	repo, err := newRepository(ctx, pool, schemaName)
	if err != nil {
		return nil, fmt.Errorf("create event repository: %w", err)
	}

	// Жизненным циклом после создания управляют Start и Stop, поэтому отмена контекста инициализации не останавливает шину.
	runCtx, runCancel := context.WithCancel(context.WithoutCancel(ctx))

	normalizedConfig := normalizeConfig(config)
	amqpTransport, err := transport.New(ctx, transport.Config{
		User:                 params.User,
		Password:             params.Password,
		Address:              params.Address,
		Vhost:                params.Vhost,
		PublisherPoolSize:    normalizedConfig.PublisherPoolSize,
		ConnectionRetryDelay: normalizedConfig.ConnectionRetryDelay,
		ConnectionTimeout:    normalizedConfig.ConnectionTimeout,
	})
	if err != nil {
		runCancel()
		return nil, fmt.Errorf("create connection: %w", err)
	}

	bus := &Bus{
		runCtx:     runCtx,
		runCancel:  runCancel,
		transport:  amqpTransport,
		repository: repo,
		workerID:   uuid.NewString(),
		config:     normalizedConfig,
	}

	return bus, nil
}

// StoreEvent сохраняет событие в outbox, включая текущую PostgreSQL-транзакцию.
func (b *Bus) StoreEvent(ctx context.Context, event Message) error {
	return b.repository.StoreEvent(ctx, event)
}

// AddSubscriber добавляет типизированный обработчик подписчика; нулевое число реплик означает одну реплику.
// AddSubscriber возвращает ошибку, если параметры подписчика неполны.
func AddSubscriber[TData any](b *Bus, subscriber string, desc Descriptor[TData], handle HandleFunc[TData], replicas uint) error {
	if b == nil {
		return fmt.Errorf("event bus is required")
	}

	if subscriber == "" {
		return fmt.Errorf("subscriber is required")
	}

	if desc.Name == "" {
		return fmt.Errorf("event name is required")
	}

	if handle == nil {
		return fmt.Errorf("event handler is required")
	}

	if b.started {
		return fmt.Errorf("event bus is already started")
	}

	subscription := &subscription[TData]{
		subscriber: subscriber,
		descriptor: desc,
		handle:     handle,
		replicas:   replicas,
	}
	if err := subscription.register(b); err != nil {
		return fmt.Errorf("register consumer: %w", err)
	}

	b.subscriptions = append(b.subscriptions, subscription)

	return nil
}

// Start запускает ранее зарегистрированные consumer-ы и публикацию outbox-событий.
func (b *Bus) Start() {
	b.started = true
	for i := range b.subscriptions {
		b.subscriptions[i].run(b)
	}

	b.runWG.Go(func() {
		b.runOutbox(b.runCtx)
	})
}

// Stop завершает consumer-ы и outbox-публикацию перед закрытием RabbitMQ.
func (b *Bus) Stop() {
	b.runCancel()
	b.runWG.Wait()
	b.transport.Close()
}

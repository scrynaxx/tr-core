package messaging

import (
	"context"
	"fmt"
	"sync"

	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/pkg/messaging/transport"
)

const exchange = "events"

// Options содержит параметры подключения шины к RabbitMQ.
type Options struct {
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
	outboxes      []outboxRepository
	outboxSchemas map[string]struct{}
	defaultOutbox OutboxRepository
	subscriptions []subscriptionRunner
	config        Config
	started       bool

	// Идентификатор отличает экземпляр шины при распределённой аренде событий.
	workerID string

	runWG sync.WaitGroup
}

// New создаёт шину с единственным outbox и сохраняет совместимость с отдельными сервисами.
func New(ctx context.Context, params Options, pool *pgxpool.Pool, schema string, config *Config) (*Bus, error) {
	bus, err := NewBus(ctx, params, config)
	if err != nil {
		return nil, err
	}

	outbox, err := bus.AddOutbox(ctx, pool, schema)
	if err != nil {
		return nil, bus.Shutdown()
	}

	bus.defaultOutbox = outbox
	return bus, nil
}

// NewBus открывает одно соединение с RabbitMQ; модульные outbox регистрируются отдельно.
func NewBus(ctx context.Context, params Options, config *Config) (*Bus, error) {
	// Жизненным циклом после создания управляют Start и Shutdown, поэтому отмена контекста инициализации не останавливает шину.
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
		runCtx:        runCtx,
		runCancel:     runCancel,
		transport:     amqpTransport,
		outboxSchemas: make(map[string]struct{}),
		workerID:      uuid.New().String(),
		config:        normalizedConfig,
	}

	return bus, nil
}

// StoreEvent сохраняет событие в default outbox отдельного сервиса.
func (b *Bus) StoreEvent(ctx context.Context, event Message) error {
	if b == nil || b.defaultOutbox == nil {
		return fmt.Errorf("default outbox is not configured")
	}

	return b.defaultOutbox.StoreEvent(ctx, event)
}

// AddOutbox создаёт и регистрирует отдельный transactional outbox для схемы модуля.
func (b *Bus) AddOutbox(ctx context.Context, pool *pgxpool.Pool, schema string) (OutboxRepository, error) {
	if b == nil {
		return nil, fmt.Errorf("event bus is required")
	}
	if b.started {
		return nil, fmt.Errorf("event bus is already started")
	}
	if _, exists := b.outboxSchemas[schema]; exists {
		return nil, fmt.Errorf("%s outbox is already registered", schema)
	}

	repo, err := newRepository(ctx, pool, schema)
	if err != nil {
		return nil, fmt.Errorf("create %s outbox: %w", schema, err)
	}

	b.outboxes = append(b.outboxes, repo)
	b.outboxSchemas[schema] = struct{}{}
	return repo, nil
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

	sub := &subscription[TData]{
		subscriber: subscriber,
		descriptor: desc,
		handle:     handle,
		replicas:   replicas,
	}
	if err := sub.register(b); err != nil {
		return fmt.Errorf("register consumer: %w", err)
	}

	b.subscriptions = append(b.subscriptions, sub)

	return nil
}

// Start запускает ранее зарегистрированные consumer-ы и публикацию outbox-событий.
func (b *Bus) Start() {
	b.started = true
	for i := range b.subscriptions {
		b.subscriptions[i].run(b)
	}

	for _, outbox := range b.outboxes {
		b.runWG.Go(func() {
			b.runOutbox(b.runCtx, outbox)
		})
	}
}

// Shutdown завершает consumer-ы и outbox-публикацию перед закрытием RabbitMQ.
func (b *Bus) Shutdown() error {
	b.runCancel()
	b.runWG.Wait()
	return b.transport.Close()
}

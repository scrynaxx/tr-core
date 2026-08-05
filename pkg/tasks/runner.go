package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Runner получает задачи из PostgreSQL и выполняет зарегистрированные handler-ы.
type Runner struct {
	repo     *repository
	config   Config
	handlers map[string]handler

	pollCtx  context.Context
	stopPoll context.CancelFunc
	workCtx  context.Context
	stopWork context.CancelFunc
	pollWG   sync.WaitGroup
	workWG   sync.WaitGroup

	mu     sync.Mutex
	active int
}

// NewRunner создаёт runner и инициализирует его таблицу в schemaName.
func NewRunner(ctx context.Context, pool *pgxpool.Pool, schemaName string, config *Config) (*Runner, error) {
	repo, err := newRepository(ctx, pool, schemaName)
	if err != nil {
		return nil, err
	}

	pollCtx, stopPoll := context.WithCancel(context.Background())
	workCtx, stopWork := context.WithCancel(context.Background())
	return &Runner{
		repo:     repo,
		config:   normalizeConfig(config),
		handlers: make(map[string]handler),
		pollCtx:  pollCtx,
		stopPoll: stopPoll,
		workCtx:  workCtx,
		stopWork: stopWork,
	}, nil
}

func (r *Runner) register(source string, handler handler) error {
	if source == "" {
		return fmt.Errorf("task source is empty")
	}
	if handler == nil {
		return fmt.Errorf("task handler is nil: source=%s", source)
	}
	if _, exists := r.handlers[source]; exists {
		return fmt.Errorf("task handler is already registered: source=%s", source)
	}

	r.handlers[source] = handler
	return nil
}

// AddHandler регистрирует обработчик и перед его вызовом декодирует JSON задачи в тип дескриптора.
// Ошибка декодирования возвращается как ошибка выполнения задачи и проходит через общий механизм retry.
func AddHandler[TPayload any](runner *Runner, descriptor Descriptor[TPayload], handler HandleFunc[TPayload]) error {
	if runner == nil {
		return fmt.Errorf("task runner is required")
	}
	if handler == nil {
		return fmt.Errorf("task handler is nil: source=%s", descriptor.Source)
	}

	return runner.register(descriptor.Source, func(ctx context.Context, data json.RawMessage) error {
		var payload TPayload
		if data != nil {
			if err := json.Unmarshal(data, &payload); err != nil {
				return fmt.Errorf("decode task payload: source=%s: %w", descriptor.Source, err)
			}
		}

		return handler(ctx, payload)
	})
}

// Enqueue идемпотентно ставит задачу по паре Source и SourceID.
func (r *Runner) Enqueue(ctx context.Context, task Task) error {
	if task.Source == "" {
		return fmt.Errorf("task source is empty")
	}
	if task.SourceID == "" {
		return fmt.Errorf("task source id is empty")
	}
	if _, ok := r.handlers[task.Source]; !ok {
		return fmt.Errorf("task handler is not registered: source=%s", task.Source)
	}
	if task.RunAt.IsZero() {
		task.RunAt = time.Now()
	}

	return r.repo.enqueue(ctx, task)
}

// Ensure добавляет задачу либо возобновляет failed-задачу с теми же source и source ID.
// При возобновлении заменяет payload, serial key и время запуска, сбрасывает число попыток
// и последнюю ошибку. Существующую pending, running или completed задачу не изменяет.
func (r *Runner) Ensure(ctx context.Context, task Task) error {
	if task.Source == "" {
		return fmt.Errorf("task source is empty")
	}
	if task.SourceID == "" {
		return fmt.Errorf("task source id is empty")
	}
	if _, ok := r.handlers[task.Source]; !ok {
		return fmt.Errorf("task handler is not registered: source=%s", task.Source)
	}
	if task.RunAt.IsZero() {
		task.RunAt = time.Now()
	}

	return r.repo.ensure(ctx, task)
}

// Start запускает получение и очистку задач; повторный вызов ничего не делает.
func (r *Runner) Start() {
	r.pollWG.Go(r.poll)
	r.pollWG.Go(r.cleanup)
}

// Stop прекращает claim новых задач и ожидает завершения уже принятых задач.
func (r *Runner) Stop() {
	r.stopPoll()
	r.pollWG.Wait()

	r.workWG.Wait()
	r.stopWork()
}

func (r *Runner) poll() {
	defer r.pollWG.Done()

	ticker := time.NewTicker(r.config.PollInterval)
	defer ticker.Stop()

	for {
		if err := r.claim(); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("[task] claim tasks failed", slog.Any("error", err))
		}

		select {
		case <-r.pollCtx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r *Runner) claim() error {
	limit := r.claimLimit()
	if limit == 0 {
		return nil
	}

	tasks, err := r.repo.claim(r.pollCtx, limit, r.config.LeaseDuration)
	if err != nil {
		return err
	}

	for i := range tasks {
		r.activate(tasks[i])
		r.workWG.Add(1)
		go r.execute(tasks[i])
	}

	return nil
}

func (r *Runner) execute(claimed ClaimedTask) {
	defer r.workWG.Done()
	defer r.deactivate(claimed)

	executor := taskExecutor{
		ctx:           r.workCtx,
		repo:          r.repo,
		task:          claimed,
		handler:       r.handlers[claimed.Source],
		leaseDuration: r.config.LeaseDuration,
	}
	err := executor.execute()
	if errors.Is(err, ErrLeaseLost) {
		return
	}

	finishCtx, finishCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer finishCancel()

	if err == nil {
		if completeErr := r.repo.complete(finishCtx, claimed); completeErr != nil && !errors.Is(completeErr, ErrLeaseLost) {
			slog.Error("[task] complete task failed", slog.String("source", claimed.Source),
				slog.String("source_id", claimed.SourceID),
				slog.Any("error", completeErr),
			)
		}
		return
	}

	if claimed.Attempt >= r.config.MaxAttempts {
		if failErr := r.repo.fail(finishCtx, claimed, err.Error()); failErr != nil && !errors.Is(failErr, ErrLeaseLost) {
			slog.Error("[task] fail task failed",
				slog.String("source", claimed.Source),
				slog.String("source_id", claimed.SourceID),
				slog.Any("error", failErr),
			)
		} else if failErr == nil {
			slog.Warn("[task] attempts exhausted",
				slog.String("source", claimed.Source),
				slog.String("source_id", claimed.SourceID),
				slog.Any("error", err),
			)
		}
		return
	}

	retryAt := time.Now().Add(r.config.RetryDelay * time.Duration(1<<(claimed.Attempt-1)))
	var retryErr *retryAtError
	if errors.As(err, &retryErr) {
		retryAt = retryErr.retryAt
	}

	if retryErr := r.repo.retry(finishCtx, claimed, retryAt, err.Error()); retryErr != nil && !errors.Is(retryErr, ErrLeaseLost) {
		slog.Error("[task] retry task failed",
			slog.String("source", claimed.Source),
			slog.String("source_id", claimed.SourceID),
			slog.Any("error", retryErr),
		)
	}
}

func callHandler(ctx context.Context, handler handler, payload []byte) (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = fmt.Errorf("panic: %v\n%s", value, debug.Stack())
		}
	}()

	return handler(ctx, payload)
}

func (r *Runner) cleanup() {
	defer r.pollWG.Done()

	ticker := time.NewTicker(r.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-r.pollCtx.Done():
			return
		case <-ticker.C:
			if err := r.repo.cleanup(r.pollCtx, time.Now().Add(-r.config.Retention)); err != nil && !errors.Is(err, context.Canceled) {
				slog.Error("[task] cleanup tasks failed", slog.Any("error", err))
			}
		}
	}
}

func (r *Runner) claimLimit() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	limit := min(r.config.BatchSize, r.config.Concurrency-r.active)
	return max(limit, 0)
}

func (r *Runner) activate(claimed ClaimedTask) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.active++
}

func (r *Runner) deactivate(claimed ClaimedTask) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.active--
}

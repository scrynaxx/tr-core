package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron     *cron.Cron
	shutdown func()
}

// NewScheduler создаёт UTC-планировщик, который пропускает повторный запуск задачи,
// пока её предыдущий вызов ещё выполняется.
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(cron.SkipIfStillRunning(new(logger))),
			cron.WithLocation(time.UTC),
		),
	}
}

// AddJob проверяет cron-выражение и добавляет задачу в расписание до его запуска.
func (s *Scheduler) AddJob(expression, description string, fn func(ctx context.Context) error) error {
	cmd := func() {
		if err := fn(context.Background()); err != nil {
			slog.Error("[scheduler] job failed",
				slog.String("expression", expression),
				slog.String("description", description),
				slog.Any("error", err),
			)
		}
	}

	if _, err := s.cron.AddFunc(expression, cmd); err != nil {
		return fmt.Errorf("add job %q with expression %q: %w", description, expression, err)
	}

	return nil
}

// Start запускает ранее добавленные задачи.
func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	shutdown := func() {
		stop := s.cron.Stop()
		<-stop.Done()
	}

	shutdown()
}

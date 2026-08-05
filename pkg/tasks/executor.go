package tasks

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// taskExecutor выполняет одну попытку задачи и поддерживает её lease.
type taskExecutor struct {
	ctx           context.Context
	repo          *repository
	task          ClaimedTask
	handler       handler
	leaseDuration time.Duration
}

func (e *taskExecutor) execute() error {
	ctx, cancel := context.WithCancel(e.ctx)
	defer cancel()

	leaseResult := make(chan error, 1)
	go func() {
		leaseResult <- e.renewLease(ctx, cancel)
	}()

	handlerErr := callHandler(ctx, e.handler, e.task.Payload)
	cancel()
	leaseErr := <-leaseResult
	if leaseErr != nil {
		return leaseErr
	}

	return handlerErr
}

func (e *taskExecutor) renewLease(ctx context.Context, cancel context.CancelFunc) error {
	ticker := time.NewTicker(max(e.leaseDuration/3, time.Nanosecond))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := e.repo.renew(ctx, e.task, e.leaseDuration); err != nil {
				if ctx.Err() != nil {
					return nil
				}
				cancel()
				if errors.Is(err, ErrLeaseLost) {
					return ErrLeaseLost
				}

				return fmt.Errorf("renew task lease: %w", err)
			}
		}
	}
}

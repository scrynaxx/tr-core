package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

func (b *Bus) runOutbox(ctx context.Context) {
	pollTimer := time.NewTimer(b.config.Outbox.PollInterval)
	cleanupTicker := time.NewTicker(b.config.Outbox.CleanupInterval)
	defer pollTimer.Stop()
	defer cleanupTicker.Stop()

	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return
		case <-pollTimer.C:
			fullBatch, err := b.dispatchOutbox(ctx)
			if err != nil && ctx.Err() == nil {
				slog.Error("[event bus] outbox dispatch failed", slog.Any("error", err))
			}

			// Полная пачка означает возможный backlog: следующая аренда начинается сразу, пустая или неполная возвращает обычный polling.
			delay := b.config.Outbox.PollInterval
			if err == nil && fullBatch {
				delay = 0
			}
			pollTimer.Reset(delay)
		case <-cleanupTicker.C:
			if err := b.cleanupOutbox(ctx); err != nil {
				slog.Error("[event bus] outbox cleanup failed", slog.Any("error", err))
			}
		}
	}
}

func (b *Bus) dispatchOutbox(ctx context.Context) (bool, error) {
	// Реальный параллелизм ограничен и числом задач, и доступными publisher-каналами; lease должен покрывать все последовательные волны пачки.
	concurrency := min(b.config.Outbox.Concurrency, b.config.PublisherPoolSize)
	waves := (b.config.Outbox.BatchSize + concurrency - 1) / concurrency
	lease := time.Duration(waves)*b.config.Outbox.PublishTimeout + b.config.Outbox.LeaseMargin

	events, err := b.repository.Claim(ctx, b.workerID, b.config.Outbox.BatchSize, lease)
	if err != nil {
		return false, fmt.Errorf("claim events: %w", err)
	}

	var group errgroup.Group
	group.SetLimit(concurrency)

	// Уже арендованная пачка завершает публикации при shutdown, после чего runOutbox прекращает брать новые события.
	workCtx := context.WithoutCancel(ctx)

	for i := range events {
		event := events[i]
		group.Go(func() error {
			return b.publishEvent(workCtx, event)
		})
	}

	if err = group.Wait(); err != nil {
		return false, err
	}

	return len(events) == b.config.Outbox.BatchSize, nil
}

func (b *Bus) publishEvent(ctx context.Context, event ClaimedEvent) error {
	publishCtx, cancel := context.WithTimeout(ctx, b.config.Outbox.PublishTimeout)
	defer cancel()

	body, err := json.Marshal(event.Message)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	err = b.transport.Publish(publishCtx, exchange, event.Name, body)
	return b.handlePublishResult(ctx, event, err)
}

func (b *Bus) handlePublishResult(ctx context.Context, event ClaimedEvent, publishErr error) error {
	if publishErr == nil {
		if err := b.repository.MarkPublished(ctx, b.workerID, event.EventID); err != nil {
			return fmt.Errorf("mark event %s published: %w", event.EventID, err)
		}
		return nil
	}

	if isTemporary(publishErr) {
		retryAt := time.Now().UTC().Add(b.config.Outbox.RetryDelay)
		if err := b.repository.Reschedule(ctx, b.workerID, event.EventID, retryAt, publishErr.Error()); err != nil {
			return fmt.Errorf("reschedule unavailable event %s after publish error %v: %w", event.EventID, publishErr, err)
		}
		return nil
	}

	if event.RetryCount >= b.config.Outbox.MaxRetries {
		if err := b.repository.MarkDead(ctx, b.workerID, event.EventID, publishErr.Error()); err != nil {
			return fmt.Errorf("mark event %s dead after publish error %v: %w", event.EventID, publishErr, err)
		}
		return nil
	}

	retryAt := time.Now().UTC().Add(b.config.Outbox.RetryDelay)
	if err := b.repository.MarkFailed(ctx, b.workerID, event.EventID, retryAt, publishErr.Error()); err != nil {
		return fmt.Errorf("mark event %s failed after publish error %v: %w", event.EventID, publishErr, err)
	}

	return nil
}

func (b *Bus) cleanupOutbox(ctx context.Context) error {
	if err := b.repository.Cleanup(ctx, time.Now().UTC().Add(-b.config.Outbox.Retention)); err != nil {
		return fmt.Errorf("cleanup events: %w", err)
	}

	return nil
}

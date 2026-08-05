package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// OutboxRepository описывает сохранение событий transactional outbox.
type OutboxRepository interface {
	// StoreEvent сохраняет событие в текущем контексте, включая открытую транзакцию.
	StoreEvent(ctx context.Context, event Message) error
}

// repository описывает управление событиями transactional outbox.
type outboxRepository interface {
	OutboxRepository

	// Claim атомарно арендует доступную пачку событий для указанного worker-а.
	Claim(ctx context.Context, workerID string, limit int, lease time.Duration) ([]ClaimedEvent, error)

	// MarkPublished завершает аренду и отмечает событие опубликованным.
	MarkPublished(ctx context.Context, workerID string, eventID uuid.UUID) error

	// Reschedule переносит временно недоступное событие без увеличения счётчика ошибок.
	Reschedule(ctx context.Context, workerID string, eventID uuid.UUID, retryAt time.Time, cause string) error

	// MarkFailed планирует повтор после учитываемой ошибки публикации.
	MarkFailed(ctx context.Context, workerID string, eventID uuid.UUID, retryAt time.Time, cause string) error

	// MarkDead окончательно прекращает публикацию события после исчерпания повторов.
	MarkDead(ctx context.Context, workerID string, eventID uuid.UUID, cause string) error

	// Cleanup удаляет завершённые события старше указанного времени.
	Cleanup(ctx context.Context, before time.Time) error
}

// ClaimedEvent содержит арендованное сообщение и состояние его публикации.
type ClaimedEvent struct {
	Message

	// RetryCount хранит число начатых повторных публикаций после первой неинфраструктурной ошибки.
	RetryCount int `json:"-"`
}

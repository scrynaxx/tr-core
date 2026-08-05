package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
)

// Repository хранит outbox-события в схеме конкретного микросервиса.
type repository struct {
	pool *pgxpool.Pool

	// Имя содержит безопасный schema-qualified путь к таблице outbox.
	table string
}

type record struct {
	ID         uuid.UUID       `db:"id"`
	Name       string          `db:"name"`
	Payload    json.RawMessage `db:"payload"`
	OccurredAt time.Time       `db:"occurred_at"`
	RetryCount int             `db:"retry_count"`
}

func newRepository(ctx context.Context, pool *pgxpool.Pool, schemaName string) (*repository, error) {
	if pool == nil {
		return nil, fmt.Errorf("postgres pool is required")
	}

	if schemaName == "" {
		return nil, fmt.Errorf("postgres schema is required")
	}

	schema := pgx.Identifier{schemaName}.Sanitize()
	table := pgx.Identifier{schemaName, "schema_outbox_events"}.Sanitize()
	pendingIndex := pgx.Identifier{"events_outbox_pending_idx"}.Sanitize()
	completedIndex := pgx.Identifier{"events_outbox_completed_idx"}.Sanitize()

	sql := fmt.Sprintf(`
		CREATE SCHEMA IF NOT EXISTS %[1]s;

		CREATE TABLE IF NOT EXISTS %[2]s (
			id uuid PRIMARY KEY,
			name text NOT NULL,
			payload jsonb NOT NULL,
			retry_count integer NOT NULL DEFAULT 0,
			locked_by text,
			locked_until timestamptz,
			last_error text,
			published_at timestamptz,
			failed_at timestamptz,
			occurred_at timestamptz NOT NULL,
			available_at timestamptz NOT NULL DEFAULT now()
		);

		COMMENT ON TABLE %[2]s IS 'События, ожидающие публикации в event bus';
		COMMENT ON COLUMN %[2]s.id IS 'Уникальный идентификатор события';
		COMMENT ON COLUMN %[2]s.name IS 'Имя события и routing key для RabbitMQ';
		COMMENT ON COLUMN %[2]s.payload IS 'Сериализованная JSON-полезная нагрузка';
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1
				FROM pg_attribute
				WHERE attrelid = '%[2]s'::regclass
				  AND attname = 'attempts'
				  AND NOT attisdropped
			) THEN
				ALTER TABLE %[2]s RENAME COLUMN attempts TO retry_count;
			END IF;
		END $$;

		COMMENT ON COLUMN %[2]s.retry_count IS 'Количество начатых повторных публикаций после первой неинфраструктурной ошибки';
		COMMENT ON COLUMN %[2]s.locked_by IS 'Идентификатор экземпляра, арендовавшего событие';
		COMMENT ON COLUMN %[2]s.locked_until IS 'Момент окончания аренды события';
		COMMENT ON COLUMN %[2]s.last_error IS 'Текст последней ошибки публикации';
		COMMENT ON COLUMN %[2]s.published_at IS 'Момент успешной публикации события';
		COMMENT ON COLUMN %[2]s.failed_at IS 'Момент окончательного прекращения публикации';
		COMMENT ON COLUMN %[2]s.occurred_at IS 'Момент возникновения события';
		COMMENT ON COLUMN %[2]s.available_at IS 'Момент доступности следующей попытки публикации';

		CREATE INDEX IF NOT EXISTS %[3]s
			ON %[2]s (available_at, occurred_at)
			WHERE published_at IS NULL AND failed_at IS NULL;

		CREATE INDEX IF NOT EXISTS %[4]s
			ON %[2]s (COALESCE(published_at, failed_at))
			WHERE published_at IS NOT NULL OR failed_at IS NOT NULL`, schema, table, pendingIndex, completedIndex)

	if _, err := pool.Exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("initialize outbox: %w", err)
	}

	return &repository{
		pool:  pool,
		table: table,
	}, nil
}

// StoreEvent сохраняет событие в outbox, используя транзакцию из контекста при её наличии.
func (p *repository) StoreEvent(ctx context.Context, event Message) error {
	sql := fmt.Sprintf(`
		INSERT INTO %s (id, name, payload, occurred_at)
		VALUES (@id, @name, @payload, @occurred_at)`, p.table)

	return postgres.Exec(ctx, p.pool, sql, pgx.NamedArgs{
		"id":          event.EventID,
		"name":        event.Name,
		"payload":     event.Data,
		"occurred_at": event.Timestamp,
	})
}

// Claim атомарно арендует доступные события через FOR UPDATE SKIP LOCKED.
func (p *repository) Claim(ctx context.Context, workerID string, limit int, lease time.Duration) ([]ClaimedEvent, error) {
	// CTE пропускает строки, уже заблокированные конкурентами, а UPDATE в том же statement атомарно фиксирует владельца и срок аренды.
	sql := fmt.Sprintf(`
		WITH candidates AS (
			SELECT id
			FROM %[1]s
			WHERE published_at IS NULL
			  AND failed_at IS NULL
			  AND available_at <= now()
			  AND (locked_until IS NULL OR locked_until <= now())
			ORDER BY occurred_at, id
			LIMIT @limit
			FOR UPDATE SKIP LOCKED
		)
		UPDATE %[1]s AS events
		SET locked_by = @worker_id,
			locked_until = now() + @lease_seconds * interval '1 second'
		FROM candidates
		WHERE events.id = candidates.id
		RETURNING events.id, events.name, events.payload, events.occurred_at, events.retry_count`, p.table)

	recs, err := postgres.Select[record](ctx, p.pool, sql, pgx.NamedArgs{
		"limit":         limit,
		"worker_id":     workerID,
		"lease_seconds": lease.Seconds(),
	})
	if err != nil {
		return nil, err
	}

	events := make([]ClaimedEvent, len(recs))
	for i := range recs {
		events[i] = ClaimedEvent{
			Message: Message{
				EventID:   recs[i].ID,
				Timestamp: recs[i].OccurredAt,
				Name:      recs[i].Name,
				Data:      recs[i].Payload,
			},
			RetryCount: recs[i].RetryCount,
		}
	}

	return events, nil
}

// MarkPublished отмечает успешно опубликованное событие, пока worker владеет его арендой.
func (p *repository) MarkPublished(ctx context.Context, workerID string, eventID uuid.UUID) error {
	sql := fmt.Sprintf(`
		UPDATE %s
		SET published_at = now(),
			locked_by = NULL,
			locked_until = NULL,
			last_error = NULL
		WHERE id = @id
		  AND locked_by = @worker_id`, p.table)

	affected, err := postgres.ExecWithAffected(ctx, p.pool, sql, pgx.NamedArgs{
		"id":        eventID,
		"worker_id": workerID,
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrLeaseLost
	}

	return nil
}

// Reschedule освобождает аренду после временной ошибки, не расходуя лимит повторов.
func (p *repository) Reschedule(ctx context.Context, workerID string, eventID uuid.UUID, retryAt time.Time, cause string) error {
	sql := fmt.Sprintf(`
		UPDATE %s
		SET available_at = @retry_at,
			locked_by = NULL,
			locked_until = NULL,
			last_error = @last_error
		WHERE id = @id
		  AND locked_by = @worker_id`, p.table)

	affected, err := postgres.ExecWithAffected(ctx, p.pool, sql, pgx.NamedArgs{
		"id":         eventID,
		"worker_id":  workerID,
		"retry_at":   retryAt,
		"last_error": cause,
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrLeaseLost
	}

	return nil
}

// MarkFailed освобождает аренду и увеличивает число учитываемых ошибок публикации.
func (p *repository) MarkFailed(ctx context.Context, workerID string, eventID uuid.UUID, retryAt time.Time, cause string) error {
	sql := fmt.Sprintf(`
		UPDATE %s
		SET available_at = @retry_at,
			locked_by = NULL,
			locked_until = NULL,
			last_error = @last_error,
			retry_count = retry_count + 1
		WHERE id = @id
		  AND locked_by = @worker_id`, p.table)

	affected, err := postgres.ExecWithAffected(ctx, p.pool, sql, pgx.NamedArgs{
		"id":         eventID,
		"worker_id":  workerID,
		"retry_at":   retryAt,
		"last_error": cause,
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrLeaseLost
	}

	return nil
}

// MarkDead окончательно прекращает публикацию события и освобождает его аренду.
func (p *repository) MarkDead(ctx context.Context, workerID string, eventID uuid.UUID, cause string) error {
	sql := fmt.Sprintf(`
		UPDATE %s
		SET failed_at = now(),
			locked_by = NULL,
			locked_until = NULL,
			last_error = @last_error
		WHERE id = @id
		  AND locked_by = @worker_id`, p.table)

	affected, err := postgres.ExecWithAffected(ctx, p.pool, sql, pgx.NamedArgs{
		"id":         eventID,
		"worker_id":  workerID,
		"last_error": cause,
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrLeaseLost
	}

	return nil
}

// Cleanup удаляет завершённые события старше указанного времени.
func (p *repository) Cleanup(ctx context.Context, before time.Time) error {
	sql := fmt.Sprintf(`
		DELETE FROM %s
		WHERE COALESCE(published_at, failed_at) < @before`, p.table)

	return postgres.Exec(ctx, p.pool, sql, pgx.NamedArgs{"before": before})
}

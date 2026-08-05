package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
)

type repository struct {
	pool  *pgxpool.Pool
	table string
}

type record struct {
	Source    string          `db:"source"`
	SourceID  string          `db:"source_id"`
	SerialKey *string         `db:"serial_key"`
	Payload   json.RawMessage `db:"payload"`
	RunAt     time.Time       `db:"run_at"`
	Attempt   int             `db:"attempt"`
}

func newRepository(ctx context.Context, pool *pgxpool.Pool, schemaName string) (*repository, error) {
	schema := pgx.Identifier{schemaName}.Sanitize()
	table := pgx.Identifier{schemaName, "schema_outbox_tasks"}.Sanitize()
	pendingIndex := pgx.Identifier{"tasks_pending_idx"}.Sanitize()
	runningIndex := pgx.Identifier{"tasks_running_idx"}.Sanitize()
	serialIndex := pgx.Identifier{"tasks_serial_idx"}.Sanitize()
	completedIndex := pgx.Identifier{"tasks_completed_idx"}.Sanitize()

	sql := fmt.Sprintf(`
		CREATE SCHEMA IF NOT EXISTS %[1]s;

		CREATE TABLE IF NOT EXISTS %[2]s (
			source text NOT NULL,
			source_id text NOT NULL,
			serial_key text,
			payload jsonb,
			status text NOT NULL,
			attempt integer NOT NULL DEFAULT 0,
			run_at timestamptz NOT NULL,
			lease_until timestamptz,
			last_error text,
			created_at timestamptz NOT NULL DEFAULT now(),
			updated_at timestamptz NOT NULL DEFAULT now(),
			PRIMARY KEY (source, source_id),
			CONSTRAINT tasks_status_check CHECK (status IN ('pending', 'running', 'completed', 'failed'))
		);

		ALTER TABLE %[2]s ADD COLUMN IF NOT EXISTS serial_key text;
		ALTER TABLE %[2]s ALTER COLUMN payload DROP NOT NULL;

		COMMENT ON TABLE %[2]s IS 'Фоновые задачи микросервиса с повторными попытками и распределённой арендой';
		COMMENT ON COLUMN %[2]s.source IS 'Тип задачи и имя зарегистрированного обработчика';
		COMMENT ON COLUMN %[2]s.source_id IS 'Уникальный идентификатор задачи внутри типа для идемпотентной постановки';
		COMMENT ON COLUMN %[2]s.serial_key IS 'Ключ FIFO-группы задач с последовательным выполнением между репликами';
		COMMENT ON COLUMN %[2]s.payload IS 'Входные данные обработчика в формате JSON; NULL для задач без входных данных';
		COMMENT ON COLUMN %[2]s.status IS 'Состояние задачи: pending, running, completed или failed';
		COMMENT ON COLUMN %[2]s.attempt IS 'Номер текущей попытки, используемый как fencing token';
		COMMENT ON COLUMN %[2]s.run_at IS 'Самый ранний момент следующего запуска задачи';
		COMMENT ON COLUMN %[2]s.lease_until IS 'Момент окончания аренды выполняющей задачу репликой';
		COMMENT ON COLUMN %[2]s.last_error IS 'Текст последней ошибки выполнения без чувствительных данных';
		COMMENT ON COLUMN %[2]s.created_at IS 'Момент постановки задачи и основа порядка внутри FIFO-группы';
		COMMENT ON COLUMN %[2]s.updated_at IS 'Момент последнего изменения состояния задачи';

		CREATE INDEX IF NOT EXISTS %[3]s ON %[2]s (run_at, source, source_id) WHERE status = 'pending';
		CREATE INDEX IF NOT EXISTS %[4]s ON %[2]s (lease_until) WHERE status = 'running';
		CREATE INDEX IF NOT EXISTS %[5]s ON %[2]s (serial_key, created_at, source, source_id)
			WHERE serial_key IS NOT NULL AND status IN ('pending', 'running');
		CREATE INDEX IF NOT EXISTS %[6]s ON %[2]s (updated_at) WHERE status IN ('completed', 'failed');`,
		schema, table, pendingIndex, runningIndex, serialIndex, completedIndex,
	)

	if _, err := pool.Exec(ctx, sql); err != nil {
		return nil, fmt.Errorf("initialize tasks: %w", err)
	}

	return &repository{pool: pool, table: table}, nil
}

func (r *repository) enqueue(ctx context.Context, task Task) error {
	sql := fmt.Sprintf(`
		INSERT INTO %s (source, source_id, serial_key, payload, status, run_at)
		VALUES (@source, @source_id, @serial_key, @payload, 'pending', @run_at)
		ON CONFLICT (source, source_id) DO NOTHING`, r.table)

	args := pgx.NamedArgs{
		"source":     task.Source,
		"source_id":  task.SourceID,
		"serial_key": task.SerialKey,
		"payload":    task.Payload,
		"run_at":     task.RunAt,
	}

	return postgres.Exec(ctx, r.pool, sql, args)
}

func (r *repository) ensure(ctx context.Context, task Task) error {
	sql := fmt.Sprintf(`
		INSERT INTO %s AS tasks (source, source_id, serial_key, payload, status, run_at)
		VALUES (@source, @source_id, @serial_key, @payload, 'pending', @run_at)
		ON CONFLICT (source, source_id) DO UPDATE SET
			serial_key = EXCLUDED.serial_key,
			payload = EXCLUDED.payload,
			status = 'pending',
			attempt = 0,
			run_at = EXCLUDED.run_at,
			lease_until = NULL,
			last_error = NULL,
			updated_at = now()
		WHERE tasks.status = 'failed'`, r.table)

	args := pgx.NamedArgs{
		"source":     task.Source,
		"source_id":  task.SourceID,
		"serial_key": task.SerialKey,
		"payload":    task.Payload,
		"run_at":     task.RunAt,
	}

	return postgres.Exec(ctx, r.pool, sql, args)
}

func (r *repository) claim(ctx context.Context, limit int, lease time.Duration) ([]ClaimedTask, error) {
	sql := fmt.Sprintf(`
		WITH candidates AS (
			SELECT candidate.source, candidate.source_id
			FROM %[1]s AS candidate
			WHERE ((candidate.status = 'pending' AND candidate.run_at <= now()) OR
			       (candidate.status = 'running' AND candidate.lease_until <= now()))
			  AND (candidate.serial_key IS NULL OR NOT EXISTS (
				SELECT 1
				FROM %[1]s AS predecessor
				WHERE predecessor.serial_key = candidate.serial_key
				  AND predecessor.status IN ('pending', 'running')
				  AND (predecessor.created_at, predecessor.source, predecessor.source_id) <
				      (candidate.created_at, candidate.source, candidate.source_id)
			  ))
			ORDER BY candidate.run_at, candidate.source, candidate.source_id
			LIMIT @limit
			FOR UPDATE OF candidate SKIP LOCKED
		)
		UPDATE %[1]s AS tasks
		SET status = 'running',
			attempt = tasks.attempt + 1,
			lease_until = now() + @lease_seconds * interval '1 second',
			updated_at = now()
		FROM candidates
		WHERE tasks.source = candidates.source
		  AND tasks.source_id = candidates.source_id
		RETURNING tasks.source, tasks.source_id, tasks.serial_key,
			tasks.payload, tasks.run_at, tasks.attempt`, r.table)

	args := pgx.NamedArgs{
		"lease_seconds": lease.Seconds(),
		"limit":         limit,
	}

	recs, err := postgres.Select[record](ctx, r.pool, sql, args)
	if err != nil {
		return nil, err
	}

	tasks := make([]ClaimedTask, len(recs))
	for i := range recs {
		tasks[i] = ClaimedTask{
			Task: Task{
				Source:    recs[i].Source,
				SourceID:  recs[i].SourceID,
				SerialKey: recs[i].SerialKey,
				Payload:   recs[i].Payload,
				RunAt:     recs[i].RunAt,
			},
			Attempt: recs[i].Attempt,
		}
	}

	return tasks, nil
}

func (r *repository) complete(ctx context.Context, task ClaimedTask) error {
	return r.finish(ctx, task, "completed", time.Time{}, "")
}

func (r *repository) renew(ctx context.Context, task ClaimedTask, lease time.Duration) error {
	sql := fmt.Sprintf(`
		UPDATE %s
		SET lease_until = now() + @lease_seconds * interval '1 second',
			updated_at = now()
		WHERE source = @source
		  AND source_id = @source_id
		  AND status = 'running'
		  AND attempt = @attempt
		  AND lease_until > now()`, r.table)

	args := pgx.NamedArgs{
		"source":        task.Source,
		"source_id":     task.SourceID,
		"attempt":       task.Attempt,
		"lease_seconds": lease.Seconds(),
	}

	affected, err := postgres.ExecWithAffected(ctx, r.pool, sql, args)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrLeaseLost
	}

	return nil
}

func (r *repository) retry(ctx context.Context, task ClaimedTask, runAt time.Time, cause string) error {
	return r.finish(ctx, task, "pending", runAt, cause)
}

func (r *repository) fail(ctx context.Context, task ClaimedTask, cause string) error {
	return r.finish(ctx, task, "failed", time.Time{}, cause)
}

func (r *repository) finish(ctx context.Context, task ClaimedTask, status string, runAt time.Time, cause string) error {
	sql := fmt.Sprintf(`
		UPDATE %s
		SET status = @status,
			run_at = CASE WHEN @run_at::timestamptz IS NULL THEN run_at ELSE @run_at END,
			lease_until = NULL,
			last_error = NULLIF(@last_error, ''),
			updated_at = now()
		WHERE source = @source
		  AND source_id = @source_id
		  AND status = 'running'
		  AND attempt = @attempt`, r.table)

	var nextRunAt *time.Time
	if !runAt.IsZero() {
		nextRunAt = &runAt
	}
	args := pgx.NamedArgs{
		"source":     task.Source,
		"source_id":  task.SourceID,
		"attempt":    task.Attempt,
		"status":     status,
		"run_at":     nextRunAt,
		"last_error": cause,
	}

	affected, err := postgres.ExecWithAffected(ctx, r.pool, sql, args)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrLeaseLost
	}

	return nil
}

func (r *repository) release(ctx context.Context, task ClaimedTask) error {
	return r.retry(ctx, task, time.Now(), "")
}

func (r *repository) unclaim(ctx context.Context, task ClaimedTask) error {
	sql := fmt.Sprintf(`
		UPDATE %s
		SET status = 'pending',
			attempt = attempt - 1,
			lease_until = NULL,
			updated_at = now()
		WHERE source = @source
		  AND source_id = @source_id
		  AND status = 'running'
		  AND attempt = @attempt`, r.table)

	args := pgx.NamedArgs{
		"source":    task.Source,
		"source_id": task.SourceID,
		"attempt":   task.Attempt,
	}

	affected, err := postgres.ExecWithAffected(ctx, r.pool, sql, args)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrLeaseLost
	}

	return nil
}

func (r *repository) cleanup(ctx context.Context, before time.Time) error {
	sql := fmt.Sprintf(`DELETE FROM %s WHERE status IN ('completed', 'failed') AND updated_at < @before`, r.table)
	args := pgx.NamedArgs{"before": before}
	return postgres.Exec(ctx, r.pool, sql, args)
}

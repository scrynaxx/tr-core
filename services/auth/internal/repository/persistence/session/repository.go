package session

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
	"github.com/scrynaxx/tr-core/services/auth/internal/model"
	"github.com/scrynaxx/tr-core/services/auth/internal/repository/persistence/session/record"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Get(ctx context.Context, refreshHash, userAgent string) (model.Session, error) {
	sql := `
	SELECT *
	FROM auth.sessions
	WHERE refresh_hash = @refresh_hash
		AND user_agent = @user_agent`

	rec, err := postgres.Get[record.Session](ctx, r.pool, sql, pgx.NamedArgs{
		"refresh_hash": refreshHash,
		"user_agent":   userAgent,
	}, postgres.WithNotFound(model.ErrSessionNotFound))
	if err != nil {
		return model.Session{}, err
	}

	return record.ToSession(rec), nil
}

func (r *Repository) Create(ctx context.Context, session model.Session) (model.Session, error) {
	sql := `
	INSERT INTO auth.sessions (id, employee_id, user_agent, refresh_hash, expires_at)
	VALUES (@id, @employee_id, @user_agent, @refresh_hash, @expires_at)
	RETURNING *`

	rec, err := postgres.Get[record.Session](ctx, r.pool, sql, pgx.NamedArgs{
		"id":           session.ID,
		"employee_id":  session.EmployeeID,
		"user_agent":   session.UserAgent,
		"refresh_hash": session.RefreshHash,
		"expires_at":   session.ExpiresAt,
	})
	if err != nil {
		return model.Session{}, err
	}

	return record.ToSession(rec), nil
}

func (r *Repository) Update(ctx context.Context, session model.Session) (model.Session, error) {
	sql := `
	UPDATE auth.sessions
	SET refresh_hash = @refresh_hash, 
		expires_at = @expires_at, 
		updated_at = now()
	WHERE id = @id
	RETURNING *`

	rec, err := postgres.Get[record.Session](ctx, r.pool, sql, pgx.NamedArgs{
		"id":           session.ID,
		"refresh_hash": session.RefreshHash,
		"expires_at":   session.ExpiresAt,
	}, postgres.WithNotFound(model.ErrSessionNotFound))
	if err != nil {
		return model.Session{}, err
	}

	return record.ToSession(rec), nil
}

func (r *Repository) Delete(ctx context.Context, employeeID, sessionID uuid.UUID) (int64, error) {
	sql := `
	DELETE FROM auth.sessions
	WHERE employee_id = @employee_id 
	  AND id = @session_id`

	return postgres.ExecWithAffected(ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id": employeeID,
		"session_id":  sessionID,
	})
}

func (r *Repository) DeleteByEmployee(ctx context.Context, employeeID uuid.UUID) ([]uuid.UUID, error) {
	sql := `
	DELETE FROM auth.sessions
	WHERE employee_id = @employee_id
	RETURNING id`

	return postgres.Select[uuid.UUID](ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id": employeeID,
	})
}

func (r *Repository) DeleteByUserAgent(ctx context.Context, employeeID uuid.UUID, userAgent string) ([]uuid.UUID, error) {
	sql := `
	DELETE FROM auth.sessions
	WHERE employee_id = @employee_id 
	  AND user_agent = @user_agent
	RETURNING id`

	return postgres.Select[uuid.UUID](ctx, r.pool, sql, pgx.NamedArgs{
		"employee_id": employeeID,
		"user_agent":  userAgent,
	})
}

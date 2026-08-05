package session

import (
	"context"

	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
	"github.com/scrynaxx/tr-core/internal/modules/auth/repository"
	"github.com/scrynaxx/tr-core/internal/modules/auth/repository/persistence/session/record"
	"github.com/scrynaxx/tr-core/internal/security"
	"github.com/scrynaxx/tr-core/pkg/database/postgres"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) repository.Session {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Get(ctx context.Context, actor security.Actor, refreshHash, userAgent string) (model.Session, error) {
	const query = `
		SELECT s.*
		FROM auth.sessions s
		LEFT JOIN auth.accounts a on a.id = s.account_id
		WHERE refresh_hash = @refresh_hash 
	  		AND user_agent = @user_agent
	  		AND a.actor = @actor`

	rec, err := postgres.Get[record.Session](ctx, r.pool, query, pgx.NamedArgs{
		"actor":        actor,
		"refresh_hash": refreshHash,
		"user_agent":   userAgent,
	}, postgres.WithNotFound(model.ErrSessionNotFound))
	if err != nil {
		return model.Session{}, err
	}

	return record.ToSession(rec), nil
}

func (r *Repository) Create(ctx context.Context, session model.Session) error {
	const query = `
		INSERT INTO auth.sessions (id, account_id, user_agent, refresh_hash, expires_at)
		VALUES (@id, @account_id, @user_agent, @refresh_hash, @expires_at)`

	return postgres.Exec(ctx, r.pool, query, pgx.NamedArgs{
		"id":           session.ID,
		"account_id":   session.AccountID,
		"user_agent":   session.UserAgent,
		"refresh_hash": session.RefreshHash,
		"expires_at":   session.ExpiresAt,
	})
}

func (r *Repository) Update(ctx context.Context, session model.Session) error {
	const query = `
		UPDATE auth.sessions
		SET refresh_hash = @refresh_hash, 
		    expires_at = @expires_at, 
		    updated_at = now()
		WHERE id = @id`

	return postgres.Exec(ctx, r.pool, query, pgx.NamedArgs{
		"id":           session.ID,
		"account_id":   session.AccountID,
		"user_agent":   session.UserAgent,
		"refresh_hash": session.RefreshHash,
		"expires_at":   session.ExpiresAt,
	})
}

func (r *Repository) Delete(ctx context.Context, accountID, sessionID uuid.UUID) (bool, error) {
	const query = `
		DELETE FROM auth.sessions 
	   	WHERE account_id = @account_id 
			AND id = @session_id`

	affected, err := postgres.ExecWithAffected(ctx, r.pool, query, pgx.NamedArgs{
		"account_id": accountID,
		"session_id": sessionID,
	})

	return affected > 0, err
}

func (r *Repository) DeleteByAccount(ctx context.Context, accountID uuid.UUID) ([]uuid.UUID, error) {
	const query = `
		DELETE FROM auth.sessions 
	   	WHERE account_id = @account_id 
	   	RETURNING id`

	return postgres.Select[uuid.UUID](ctx, r.pool, query, pgx.NamedArgs{"account_id": accountID})
}

func (r *Repository) DeleteByUserAgent(ctx context.Context, accountID uuid.UUID, userAgent string) ([]uuid.UUID, error) {
	const query = `
		DELETE FROM auth.sessions 
	   	WHERE account_id = @account_id 
	  		AND user_agent = @user_agent 
	   	RETURNING id`

	return postgres.Select[uuid.UUID](ctx, r.pool, query, pgx.NamedArgs{
		"account_id": accountID,
		"user_agent": userAgent,
	})
}

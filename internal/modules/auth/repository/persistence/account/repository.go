package account

import (
	"context"

	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scrynaxx/tr-core/internal/modules/auth/model"
	"github.com/scrynaxx/tr-core/internal/modules/auth/repository"
	"github.com/scrynaxx/tr-core/internal/modules/auth/repository/persistence/account/record"
	"github.com/scrynaxx/tr-core/internal/security"

	"github.com/scrynaxx/tr-core/pkg/database/postgres"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) repository.Account {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Get(ctx context.Context, accountID uuid.UUID) (model.Account, error) {
	const query = `
		SELECT *
		FROM auth.accounts
		WHERE id = @account_id`

	rec, err := postgres.Get[record.Account](ctx, r.pool, query, pgx.NamedArgs{"account_id": accountID}, postgres.WithNotFound(model.ErrAccountNotFound))
	if err != nil {
		return model.Account{}, err
	}

	return record.ToAccount(rec), nil
}

func (r *Repository) GetByEmail(ctx context.Context, actor security.Actor, email string) (model.Account, error) {
	const query = `
		SELECT *
		FROM auth.accounts
		WHERE actor = @actor
			AND email = @email`

	rec, err := postgres.Get[record.Account](ctx, r.pool, query, pgx.NamedArgs{
		"actor": actor,
		"email": email,
	}, postgres.WithNotFound(model.ErrAccountNotFound))
	if err != nil {
		return model.Account{}, err
	}

	return record.ToAccount(rec), nil
}

func (r *Repository) Create(ctx context.Context, account model.Account) error {
	const query = `
		INSERT INTO auth.accounts (id, actor, email, password_hash)
		VALUES (@id, @actor, @email, @password_hash)`

	return postgres.Exec(ctx, r.pool, query, pgx.NamedArgs{
		"id":            account.ID,
		"actor":         account.Actor,
		"email":         account.Email,
		"password_hash": account.PasswordHash,
	}, postgres.WithExists(model.ErrEmailAlreadyTaken))
}

func (r *Repository) Update(ctx context.Context, account model.Account) error {
	const query = `
		UPDATE auth.accounts
		SET email = @email, 
		    password_hash = @password_hash
		WHERE id = @account_id`

	return postgres.Exec(ctx, r.pool, query, pgx.NamedArgs{
		"id":            account.ID,
		"email":         account.Email,
		"password_hash": account.PasswordHash,
	})
}

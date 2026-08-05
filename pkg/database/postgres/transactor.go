package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type Transactor struct {
	pool *pgxpool.Pool
}

func NewTransactor(pool *pgxpool.Pool) *Transactor {
	return &Transactor{pool: pool}
}

func (t *Transactor) Call(ctx context.Context, fn func(ctx context.Context) error) error {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok && tx != nil {
		return fn(ctx)
	}

	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	ctx = context.WithValue(ctx, txKey{}, tx)
	if err = fn(ctx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

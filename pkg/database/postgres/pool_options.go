package postgres

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PoolOption выполняет дополнительную настройку после создания и проверки pool.
type PoolOption func(context.Context, *pgxpool.Pool) error

// WithMigrations применяет миграции указанной схемы при создании pool.
func WithMigrations(schema string, migrations fs.FS) PoolOption {
	return func(ctx context.Context, pool *pgxpool.Pool) error {
		if err := Migrate(ctx, pool, schema, migrations); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}

		return nil
	}
}

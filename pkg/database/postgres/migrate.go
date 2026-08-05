package postgres

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

const (
	// Table Таблица миграций
	tableName = "schema_migrations"

	// Driver Драйвер для работы с БД
	driverName = "pgx5"
)

// Migrate создаёт схему и применяет все ещё не выполненные SQL-миграции через существующий pool.
func Migrate(ctx context.Context, pool *pgxpool.Pool, schema string, migrations fs.FS) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	sql := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, pgx.Identifier{schema}.Sanitize())
	if _, err := pool.Exec(ctx, sql); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		SchemaName:      schema,
		MigrationsTable: tableName,
	})
	if err != nil {
		return err
	}
	defer driver.Close()

	source, err := iofs.New(migrations, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}
	defer source.Close()

	m, err := migrate.NewWithInstance("iofs", source, driverName, driver)
	if err != nil {
		return fmt.Errorf("create migration: %w", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

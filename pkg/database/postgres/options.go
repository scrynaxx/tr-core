package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorHandler преобразует ошибку PostgreSQL в предметную ошибку.
type ErrorHandler func(error) error

// GetErrorHandler преобразует ошибку чтения одной записи.
type GetErrorHandler = ErrorHandler

// ExecErrorHandler преобразует ошибку выполнения команды.
type ExecErrorHandler = ErrorHandler

// WithNotFound заменяет pgx.ErrNoRows указанной ошибкой.
func WithNotFound(notFoundErr error) GetErrorHandler {
	return func(err error) error {
		if errors.Is(err, pgx.ErrNoRows) {
			return notFoundErr
		}

		return err
	}
}

// WithExists заменяет ошибку нарушения уникальности указанной предметной ошибкой.
func WithExists(existsErr error) GetErrorHandler {
	return func(err error) error {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return existsErr
		}

		return err
	}
}

package postgres

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"
)

func Get[TResult any](ctx context.Context, q Querier, sql string, args pgx.NamedArgs, errorHandlers ...GetErrorHandler) (TResult, error) {
	var zero TResult

	mapper, err := getMapper[TResult]()
	if err != nil {
		return zero, fmt.Errorf("get mapper: %w", err)
	}

	rows, err := getQuerier(ctx, q).Query(ctx, sql, args)
	if err != nil {
		for _, handler := range errorHandlers {
			err = handler(err)
		}

		return zero, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	rec, err := pgx.CollectOneRow(rows, mapper)
	if err != nil {
		for _, handler := range errorHandlers {
			err = handler(err)
		}

		return zero, err
	}

	return rec, nil
}

func GetPtr[TResult any](ctx context.Context, q Querier, sql string, args pgx.NamedArgs, errorHandlers ...GetErrorHandler) (*TResult, error) {
	mapper, err := getPtrMapper[TResult]()
	if err != nil {
		return nil, fmt.Errorf("get mapper: %w", err)
	}

	rows, err := getQuerier(ctx, q).Query(ctx, sql, args)
	if err != nil {
		for _, handler := range errorHandlers {
			err = handler(err)
		}

		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	rec, err := pgx.CollectOneRow(rows, mapper)
	if err != nil {
		for _, handler := range errorHandlers {
			err = handler(err)
		}

		return nil, err
	}

	return rec, nil
}

func Find[TResult any](ctx context.Context, q Querier, sql string, args pgx.NamedArgs) (*TResult, error) {
	mapper, err := getPtrMapper[TResult]()
	if err != nil {
		return nil, fmt.Errorf("get mapper: %w", err)
	}

	rows, err := getQuerier(ctx, q).Query(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	rec, err := pgx.CollectOneRow(rows, mapper)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("collect one row: %w", err)
	}

	return rec, nil
}

func Select[TResult any](ctx context.Context, q Querier, sql string, args pgx.NamedArgs) ([]TResult, error) {
	mapper, err := getMapper[TResult]()
	if err != nil {
		return nil, fmt.Errorf("get mapper: %w", err)
	}

	rows, err := getQuerier(ctx, q).Query(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, mapper)
}

func Exec(ctx context.Context, q Querier, sql string, args pgx.NamedArgs, errorHandlers ...ExecErrorHandler) error {
	q = getQuerier(ctx, q)

	if _, err := q.Exec(ctx, sql, args); err != nil {
		for _, handler := range errorHandlers {
			err = handler(err)
		}

		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func ExecWithAffected(ctx context.Context, q Querier, sql string, args pgx.NamedArgs, errorHandlers ...ExecErrorHandler) (int64, error) {
	q = getQuerier(ctx, q)

	ct, err := q.Exec(ctx, sql, args)
	if err != nil {
		for _, handler := range errorHandlers {
			err = handler(err)
		}

		return 0, fmt.Errorf("exec: %w", err)
	}

	return ct.RowsAffected(), nil
}

func Exists(ctx context.Context, q Querier, sql string, args pgx.NamedArgs) (bool, error) {
	rows, err := getQuerier(ctx, q).Query(ctx, sql, args)
	if err != nil {
		return false, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("collect row: %w", err)
	}

	return result, nil
}

func SendBatch(ctx context.Context, q Querier, batch *pgx.Batch) error {
	return getQuerier(ctx, q).SendBatch(ctx, batch).Close()
}

func getQuerier(ctx context.Context, q Querier) Querier {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok && tx != nil {
		return tx
	}

	return q
}

func getMapper[TResult any]() (pgx.RowToFunc[TResult], error) {
	kind := reflect.TypeFor[TResult]().Kind()

	switch kind {
	case reflect.Struct:
		return pgx.RowToStructByName[TResult], nil
	default:
		return pgx.RowTo[TResult], nil
	}
}

func getPtrMapper[TResult any]() (pgx.RowToFunc[*TResult], error) {
	t := reflect.TypeFor[TResult]()

	if t == reflect.TypeFor[time.Time]() {
		return pgx.RowToAddrOf[TResult], nil
	}

	switch t.Kind() {
	case reflect.Struct:
		return pgx.RowToAddrOfStructByName[TResult], nil
	default:
		return pgx.RowToAddrOf[TResult], nil
	}
}

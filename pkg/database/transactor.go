package database

import "context"

// Transactor описывает выполнение нескольких операций в одной транзакции.
type Transactor interface {
	// Call открывает транзакцию и передаёт её через контекст callback-функции.
	Call(ctx context.Context, fn func(ctx context.Context) error) error
}

// Tx открывает транзакцию и передаёт её через контекст callback-функции.
func Tx(ctx context.Context, transactor Transactor, fn func(context.Context) error) error {
	return transactor.Call(ctx, fn)
}

// TxResult открывает транзакцию и передаёт её через контекст callback-функции, возвращая результат выполнения функции.
func TxResult[T any](ctx context.Context, transactor Transactor, fn func(context.Context) (T, error)) (T, error) {
	var result T

	err := transactor.Call(ctx, func(ctx context.Context) error {
		var err error

		result, err = fn(ctx)
		return err
	})
	if err != nil {
		var zero T
		return zero, err
	}

	return result, nil
}

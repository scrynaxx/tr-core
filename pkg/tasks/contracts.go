package tasks

import "context"

// Enqueuer ставит задачи в очередь на выполнение.
type Enqueuer interface {
	// Enqueue добавляет задачу, если сочетание source и source ID ещё не существует.
	// Существующую pending, running, completed или failed задачу не изменяет.
	Enqueue(ctx context.Context, task Task) error

	// Ensure добавляет задачу либо возобновляет failed-задачу с теми же source и source ID.
	// При возобновлении заменяет payload, serial key и время запуска, сбрасывает число попыток
	// и последнюю ошибку. Существующую pending, running или completed задачу не изменяет.
	Ensure(ctx context.Context, task Task) error
}

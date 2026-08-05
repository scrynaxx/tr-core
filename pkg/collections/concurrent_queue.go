package collections

import "sync"

type ConcurrentQueue[T any] struct {
	mu    sync.Mutex
	queue Queue[T]
}

func (q *ConcurrentQueue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.queue.Enqueue(item)
}

func (q *ConcurrentQueue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.queue.Dequeue()
}

func (q *ConcurrentQueue[T]) Peek() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.queue.Peek()
}

func (q *ConcurrentQueue[T]) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.queue.IsEmpty()
}

func (q *ConcurrentQueue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.queue.Size()
}

func (q *ConcurrentQueue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.queue.Clear()
}

func (q *ConcurrentQueue[T]) Values() []T {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.queue.Values()
}

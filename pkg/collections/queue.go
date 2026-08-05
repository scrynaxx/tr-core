package collections

type Queue[T any] struct {
	buf  []T
	head int
	tail int
	size int
}

func (q *Queue[T]) Enqueue(item T) {
	if q.size == len(q.buf) {
		q.grow()
	}

	q.buf[q.tail] = item
	q.tail = (q.tail + 1) % len(q.buf)
	q.size++
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	item := q.buf[q.head]

	var zero T
	q.buf[q.head] = zero

	q.head = (q.head + 1) % len(q.buf)
	q.size--

	if q.size == 0 {
		q.head = 0
		q.tail = 0
	}

	return item, true
}

func (q *Queue[T]) Peek() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	return q.buf[q.head], true
}

func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

func (q *Queue[T]) Size() int {
	return q.size
}

func (q *Queue[T]) Clear() {
	var zero T

	if q.size == 0 {
		return
	}

	if q.head < q.tail {
		for i := q.head; i < q.tail; i++ {
			q.buf[i] = zero
		}
	} else {
		for i := q.head; i < len(q.buf); i++ {
			q.buf[i] = zero
		}
		for i := 0; i < q.tail; i++ {
			q.buf[i] = zero
		}
	}

	q.head = 0
	q.tail = 0
	q.size = 0
}

func (q *Queue[T]) Values() []T {
	if q.size == 0 {
		return nil
	}

	out := make([]T, q.size)

	if q.head < q.tail {
		copy(out, q.buf[q.head:q.tail])
		return out
	}

	n := copy(out, q.buf[q.head:])
	copy(out[n:], q.buf[:q.tail])

	return out
}

func (q *Queue[T]) grow() {
	newCap := 4
	if len(q.buf) > 0 {
		newCap = len(q.buf) * 2
	}

	newBuf := make([]T, newCap)

	if q.size > 0 {
		if q.head < q.tail {
			copy(newBuf, q.buf[q.head:q.tail])
		} else {
			n := copy(newBuf, q.buf[q.head:])
			copy(newBuf[n:], q.buf[:q.tail])
		}
	}

	q.buf = newBuf
	q.head = 0
	q.tail = q.size
}

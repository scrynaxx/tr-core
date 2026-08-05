package collections

import "sync"

type ConcurrentLock[TKey comparable] struct {
	ma sync.Map
}

func (l *ConcurrentLock[TKey]) Lock(key TKey) *sync.Mutex {
	if v, ok := l.ma.Load(key); ok {
		mu := v.(*sync.Mutex)
		mu.Lock()

		return mu
	}

	v, _ := l.ma.LoadOrStore(key, l.setter())
	mu := v.(*sync.Mutex)
	mu.Lock()

	return mu
}

func (l *ConcurrentLock[TKey]) setter() *sync.Mutex {
	return new(sync.Mutex)
}

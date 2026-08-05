package collections

import "sync"

type ConcurrentMap[TKey comparable, TValue any] struct {
	ma sync.Map
}

func (m *ConcurrentMap[TKey, TValue]) Get(key TKey) (TValue, bool) {
	if value, ok := m.ma.Load(key); ok {
		return value.(TValue), true
	}

	var zero TValue

	return zero, false
}

func (m *ConcurrentMap[TKey, TValue]) Store(key TKey, value TValue) {
	m.ma.Store(key, value)
}

func (m *ConcurrentMap[TKey, TValue]) Delete(key TKey) {
	m.ma.Delete(key)
}

func (m *ConcurrentMap[TKey, TValue]) Clear() {
	m.ma.Clear()
}

func (m *ConcurrentMap[TKey, TValue]) LoadOrStore(key TKey, setter func() TValue) TValue {
	if value, ok := m.ma.Load(key); ok {
		return value.(TValue)
	}

	value, _ := m.ma.LoadOrStore(key, setter())
	return value.(TValue)
}

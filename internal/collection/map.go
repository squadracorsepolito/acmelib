package collection

import (
	"iter"
	"maps"
)

type Map[K comparable, V any] struct {
	m map[K]V
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

func (m *Map[K, V]) Set(key K, value V) {
	m.m[key] = value
}

func (m *Map[K, V]) Delete(key K) {
	delete(m.m, key)
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	val, ok := m.m[key]
	return val, ok
}

func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.m[key]
	return ok
}

func (m *Map[K, V]) Size() int {
	return len(m.m)
}

func (m *Map[K, V]) Clear() {
	clear(m.m)
}

func (m *Map[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.m)
}

func (m *Map[K, V]) Values() iter.Seq[V] {
	return maps.Values(m.m)
}

func (m *Map[K, V]) Entries() iter.Seq2[K, V] {
	return maps.All(m.m)
}

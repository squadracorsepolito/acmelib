package collection

import (
	"iter"
	"maps"
)

// Map is a generic map.
type Map[K comparable, V any] struct {
	m map[K]V
}

// NewMap returns a new [Map].
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

// Set sets the value for the given key.
func (m *Map[K, V]) Set(key K, value V) {
	m.m[key] = value
}

// Delete deletes the value for the given key.
func (m *Map[K, V]) Delete(key K) {
	delete(m.m, key)
}

// Get returns the value for the given key.
// It returns false if the key does not exist.
func (m *Map[K, V]) Get(key K) (V, bool) {
	val, ok := m.m[key]
	return val, ok
}

// Has returns whether the given key exists.
func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.m[key]
	return ok
}

// Size returns the number of elements in the map.
func (m *Map[K, V]) Size() int {
	return len(m.m)
}

// Clear clears the map.
func (m *Map[K, V]) Clear() {
	clear(m.m)
}

// Keys returns an iterator over all keys.
func (m *Map[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.m)
}

// Values returns an iterator over all values.
func (m *Map[K, V]) Values() iter.Seq[V] {
	return maps.Values(m.m)
}

// Entries returns an iterator over all entries (key, value).
func (m *Map[K, V]) Entries() iter.Seq2[K, V] {
	return maps.All(m.m)
}

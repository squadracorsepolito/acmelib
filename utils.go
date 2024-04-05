package acmelib

import (
	"fmt"
)

const maxSize = 64

func calcSizeFromValue(val int) int {
	for i := 0; i < maxSize; i++ {
		if val < 1<<i {
			return i
		}
	}
	return maxSize
}

type set[K comparable, V any] struct {
	m         map[K]V
	errPrefix string
}

func newSet[K comparable, V any](errPrefix string) *set[K, V] {
	return &set[K, V]{
		m:         make(map[K]V),
		errPrefix: errPrefix,
	}
}

func (s *set[K, V]) verifyKey(key K) error {
	if _, ok := s.m[key]; ok {
		return fmt.Errorf(`%s "%v" is duplicated`, s.errPrefix, key)
	}
	return nil
}

func (s *set[K, V]) add(key K, val V) {
	s.m[key] = val
}

func (s *set[K, V]) remove(key K) {
	delete(s.m, key)
}

func (s *set[K, V]) hasKey(key K) bool {
	_, ok := s.m[key]
	return ok
}

func (s *set[K, V]) modifyKey(oldKey, newKey K, val V) {
	s.remove(oldKey)
	s.add(newKey, val)
}

func (s *set[K, V]) getValue(key K) (V, error) {
	val, ok := s.m[key]
	if ok {
		return val, nil
	}
	return val, fmt.Errorf(`%s "%v" not found`, s.errPrefix, key)
}

func (s *set[K, V]) getKeys() []K {
	count := len(s.m)
	keys := make([]K, count)
	i := 0
	for k := range s.m {
		keys[i] = k
		i++
	}
	return keys
}

func (s *set[K, V]) size() int {
	return len(s.m)
}

func (s *set[K, V]) clear() {
	for key := range s.m {
		delete(s.m, key)
	}
}

func (s *set[K, V]) entries() map[K]V {
	return s.m
}

func (s *set[K, V]) getValues() []V {
	values := make([]V, s.size())
	i := 0
	for _, val := range s.m {
		values[i] = val
		i++
	}
	return values
}

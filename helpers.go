package acmelib

import (
	"math"
	"strings"
)

const maxSize = 64

func calcSizeFromValue(val int) int {
	if val == 0 {
		return 1
	}

	for i := 0; i < maxSize; i++ {
		if val < 1<<i {
			return i
		}
	}

	return maxSize
}

func getSizeFromCount(count int) int {
	if count == 0 {
		return 0
	}

	for i := range maxSize {
		if count <= 1<<i {
			return i
		}
	}

	return maxSize
}

func calcValueFromSize(size int) int {
	if size <= 0 {
		return 1
	}
	return 1 << size
}

func isDecimal(val float64) bool {
	return math.Mod(val, 1.0) != 0
}

func clearSpaces(str string) string {
	return strings.ReplaceAll(strings.TrimSpace(str), " ", "_")
}

type set[K comparable, V any] struct {
	m map[K]V
}

func newSet[K comparable, V any]() *set[K, V] {
	return &set[K, V]{
		m: make(map[K]V),
	}
}

func (s *set[K, V]) verifyKeyUnique(key K) error {
	if _, ok := s.m[key]; ok {
		return ErrIsDuplicated
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
	return val, ErrNotFound
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

func getTabString(tabs int) string {
	tabStr := ""
	for i := 0; i < tabs; i++ {
		tabStr += "\t"
	}
	return tabStr
}

type stack[T any] struct {
	items []T
}

func newStack[T any]() *stack[T] {
	return &stack[T]{
		items: []T{},
	}
}

func (s *stack[T]) push(item T) {
	s.items = append(s.items, item)
}

func (s *stack[T]) pop() T {
	lastIdx := s.size() - 1
	item := s.items[lastIdx]
	s.items = s.items[:lastIdx]
	return item
}

func (s *stack[T]) size() int {
	return len(s.items)
}

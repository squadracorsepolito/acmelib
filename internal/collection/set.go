package collection

// Set is a generic set.
type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet returns a new [Set].
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		m: make(map[T]struct{}),
	}
}

// Add adds a value to the set.
func (s *Set[T]) Add(value T) {
	s.m[value] = struct{}{}
}

// Delete deletes the given value from the set.
func (s *Set[T]) Delete(value T) {
	delete(s.m, value)
}

// Has returns whether the given value exists in the set.
func (s *Set[T]) Has(value T) bool {
	_, ok := s.m[value]
	return ok
}

// Size returns the number of elements in the set.
func (s *Set[T]) Size() int {
	return len(s.m)
}

// Clear clears the set.
func (s *Set[T]) Clear() {
	clear(s.m)
}

package collection

// Stack is a generic stack.
type Stack[T any] struct {
	items []T
}

// NewStack creates a new [Stack].
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: []T{},
	}
}

// Push adds an item on top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop extracts the first item from the stack.
func (s *Stack[T]) Pop() T {
	if len(s.items) == 0 {
		return *new(T)
	}

	lastIdx := len(s.items) - 1
	item := s.items[lastIdx]
	s.items = s.items[:lastIdx]
	return item
}

// Size returns the number of items in the stack.
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// IsEmpty returns true if the stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

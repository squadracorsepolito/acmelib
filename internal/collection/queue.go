package collection

// Queue is a generic queue.
type Queue[T any] struct {
	items []T
}

// NewQueue creates a new [Queue].
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: []T{},
	}
}

// Push adds an item to the queue.
func (q *Queue[T]) Push(item T) {
	q.items = append(q.items, item)
}

// Pop removes an item from the queue.
func (q *Queue[T]) Pop() T {
	if len(q.items) == 0 {
		return *new(T)
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}

// Size returns the number of items in the queue.
func (q *Queue[T]) Size() int {
	return len(q.items)
}

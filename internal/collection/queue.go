package collection

type Queue[T any] struct {
	items []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: []T{},
	}
}

func (q *Queue[T]) Push(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Pop() T {
	if len(q.items) == 0 {
		return *new(T)
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}

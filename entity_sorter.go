package acmelib

import (
	"slices"
	"strings"
	"sync"
	"time"
)

type sortableEntity interface {
	getName() string
	getCreateTime() time.Time
	getUpdateTime() time.Time
}

type entitySorterMethod[N ~string, E sortableEntity] struct {
	name N
	fn   func([]E) []E
}

func newEntitySorterMethod[N ~string, E sortableEntity](name N, fn func([]E) []E) *entitySorterMethod[N, E] {
	return &entitySorterMethod[N, E]{name, fn}
}

type entitySorter[N ~string, E sortableEntity] struct {
	selectedMethod *entitySorterMethod[N, E]
	mux            sync.RWMutex
	methods        map[N]*entitySorterMethod[N, E]
}

func newEntitySorter[N ~string, E sortableEntity](sortingMethods ...*entitySorterMethod[N, E]) *entitySorter[N, E] {
	sorter := &entitySorter[N, E]{
		methods: make(map[N]*entitySorterMethod[N, E]),
	}

	for _, method := range sortingMethods {
		sorter.methods[method.name] = method
	}

	if len(sortingMethods) > 0 {
		sorter.selectedMethod = sortingMethods[0]
	}

	return sorter
}

func (es *entitySorter[N, E]) selectSortingMethod(methodName N) {
	es.mux.Lock()
	defer es.mux.Unlock()

	method, ok := es.methods[methodName]
	if ok {
		es.selectedMethod = method
	}
}

func (es *entitySorter[N, E]) sortEntities(entities []E) []E {
	es.mux.RLock()
	defer es.mux.RUnlock()

	if es.selectedMethod == nil {
		return entities
	}

	return es.selectedMethod.fn(entities)
}

func sortByName[E sortableEntity](entities []E) []E {
	slices.SortFunc(entities, func(a, b E) int { return strings.Compare(a.getName(), b.getName()) })
	return entities
}

func sortByCreateTime[E sortableEntity](entities []E) []E {
	slices.SortFunc(entities, func(a, b E) int { return b.getCreateTime().Compare(a.getCreateTime()) })
	return entities
}

func sortByUpdateTime[E sortableEntity](entities []E) []E {
	slices.SortFunc(entities, func(a, b E) int { return b.getUpdateTime().Compare(a.getUpdateTime()) })
	return entities
}

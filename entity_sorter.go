package acmelib

import (
	"slices"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

type entitySorterMethod[N ~string, E any] struct {
	name N
	fn   func([]E) []E
}

func newEntitySorterMethod[N ~string, E any](name N, fn func([]E) []E) *entitySorterMethod[N, E] {
	return &entitySorterMethod[N, E]{name, fn}
}

type entitySorter[N ~string, E any] struct {
	methods map[N]*entitySorterMethod[N, E]
}

func newEntitySorter[N ~string, E any](sortingMethods ...*entitySorterMethod[N, E]) *entitySorter[N, E] {
	sorter := &entitySorter[N, E]{
		methods: make(map[N]*entitySorterMethod[N, E]),
	}

	for _, method := range sortingMethods {
		sorter.methods[method.name] = method
	}

	return sorter
}

func (es *entitySorter[N, E]) sortEntities(methodName N, entities []E) []E {
	if method, ok := es.methods[methodName]; ok {
		return method.fn(entities)
	}

	return entities
}

func (es *entitySorter[N, E]) listSortingMethodNames() []N {
	return maps.Keys(es.methods)
}

type sortableByName interface {
	getName() string
}

func sortByName[E sortableByName](entities []E) []E {
	slices.SortFunc(entities, func(a, b E) int { return strings.Compare(a.getName(), b.getName()) })
	return entities
}

type sortableByCreateTime interface {
	getCreateTime() time.Time
}

func sortByCreateTime[E sortableByCreateTime](entities []E) []E {
	slices.SortFunc(entities, func(a, b E) int { return b.getCreateTime().Compare(a.getCreateTime()) })
	return entities
}

type sortableByUpdateTime interface {
	getUpdateTime() time.Time
}

func sortByUpdateTime[E sortableByUpdateTime](entities []E) []E {
	slices.SortFunc(entities, func(a, b E) int { return b.getUpdateTime().Compare(a.getUpdateTime()) })
	return entities
}

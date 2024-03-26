package acmelib

import (
	"slices"
	"strings"
	"time"
)

const maxSize = 64

func calcSizeFromValue(val int) int {
	for i := 0; i < maxSize; i++ {
		if val <= 1<<i {
			return i + 1
		}
	}
	return maxSize
}

type sortableByName interface {
	GetName() string
}

func sortByName[E sortableByName](entities []E) []E {
	slices.SortFunc(entities, func(a, b E) int { return strings.Compare(a.GetName(), b.GetName()) })
	return entities
}

type sortableByCreateTime interface {
	GetCreateTime() time.Time
}

func sortByCreateTime[E sortableByCreateTime](entities []E) []E {
	slices.SortFunc(entities, func(a, b E) int { return b.GetCreateTime().Compare(a.GetCreateTime()) })
	return entities
}

type sortableByUpdateTime interface {
	GetUpdateTime() time.Time
}

func sortByUpdateTime[E sortableByUpdateTime](entities []E) []E {
	slices.SortFunc(entities, func(a, b E) int { return b.GetUpdateTime().Compare(a.GetUpdateTime()) })
	return entities
}

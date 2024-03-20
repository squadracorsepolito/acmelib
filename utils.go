package acmelib

import (
	"slices"
	"strings"
	"time"
)

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

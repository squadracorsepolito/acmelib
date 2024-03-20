package acmelib

import (
	"fmt"
	"sync"

	"golang.org/x/exp/maps"
)

type collectableEntity interface {
	getID() EntityID
	getName() string
}

type entityCollection[E collectableEntity, SM ~string] struct {
	entityMap map[EntityID]E
	nameSet   map[string]EntityID

	sorter                  *entitySorter[SM, E]
	availableSortingMethods []SM
	sortedMux               sync.RWMutex
	sortedEntities          map[SM][]E
}

func newEntityCollection[E collectableEntity, SM ~string](sorter *entitySorter[SM, E]) *entityCollection[E, SM] {
	return &entityCollection[E, SM]{
		entityMap: make(map[EntityID]E),
		nameSet:   make(map[string]EntityID),

		sorter:                  sorter,
		availableSortingMethods: sorter.listSortingMethodNames(),
		sortedEntities:          make(map[SM][]E),
	}
}

func (ec *entityCollection[E, SM]) size() int {
	return len(ec.entityMap)
}

func (ec *entityCollection[E, SM]) updateSortedEntities() {
	ec.sortedMux.Lock()
	defer ec.sortedMux.Unlock()

	entities := maps.Values(ec.entityMap)
	for _, sm := range ec.availableSortingMethods {
		ec.sortedEntities[sm] = ec.sorter.sortEntities(sm, entities)
	}
}

func (ec *entityCollection[E, SM]) listEntities(sortingMethod SM) []E {
	ec.sortedMux.RLock()
	defer ec.sortedMux.RUnlock()

	if sorted, ok := ec.sortedEntities[sortingMethod]; ok {
		return sorted
	}

	return maps.Values(ec.entityMap)
}

func (ec *entityCollection[E, SM]) verifyName(name string) error {
	if _, ok := ec.nameSet[name]; ok {
		return fmt.Errorf(`duplicated name "%s"`, name)
	}
	return nil
}

func (ec *entityCollection[E, SM]) addEntity(entity E) error {
	name := entity.getName()
	if err := ec.verifyName(name); err != nil {
		return err
	}

	id := entity.getID()
	if _, ok := ec.entityMap[id]; ok {
		return fmt.Errorf(`duplicated id "%s"`, id)
	}

	ec.entityMap[id] = entity
	ec.nameSet[name] = id

	ec.updateSortedEntities()

	return nil
}

func (ec *entityCollection[E, SM]) removeEntity(id EntityID) error {
	e, err := ec.getEntityByID(id)
	if err != nil {
		return err
	}

	delete(ec.nameSet, e.getName())
	delete(ec.entityMap, id)

	ec.updateSortedEntities()

	return nil
}

func (ec *entityCollection[E, SM]) updateName(id EntityID, oldName, newName string) error {
	if oldName == newName {
		return fmt.Errorf(`"%s" is not a new name`, newName)
	}

	if err := ec.verifyName(newName); err != nil {
		return err
	}

	ec.nameSet[newName] = id
	delete(ec.nameSet, oldName)

	ec.updateSortedEntities()

	return nil
}

func (ec *entityCollection[E, SM]) getEntityByID(id EntityID) (E, error) {
	e, ok := ec.entityMap[id]
	if ok {
		return e, nil
	}
	return e, fmt.Errorf(`id "%s" not found`, id)
}

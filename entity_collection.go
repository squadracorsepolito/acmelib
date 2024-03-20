package acmelib

import (
	"fmt"

	"golang.org/x/exp/maps"
)

type collectableEntity interface {
	getID() EntityID
	getName() string
}

type entityCollection[E collectableEntity] struct {
	entityMap map[EntityID]E
	nameSet   map[string]EntityID
}

func newEntityCollection[E collectableEntity]() *entityCollection[E] {
	return &entityCollection[E]{
		entityMap: make(map[EntityID]E),
		nameSet:   make(map[string]EntityID),
	}
}

func (ec *entityCollection[E]) size() int {
	return len(ec.entityMap)
}

func (ec *entityCollection[E]) listEntities() []E {
	return maps.Values(ec.entityMap)
}

func (ec *entityCollection[E]) verifyName(name string) error {
	if _, ok := ec.nameSet[name]; ok {
		return fmt.Errorf(`duplicated name "%s"`, name)
	}
	return nil
}

func (ec *entityCollection[E]) addEntity(entity E) error {
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

	return nil
}

func (ec *entityCollection[E]) updateName(id EntityID, oldName, newName string) error {
	if oldName == newName {
		return fmt.Errorf(`"%s" is not a new name`, newName)
	}

	if err := ec.verifyName(newName); err != nil {
		return err
	}

	ec.nameSet[newName] = id
	delete(ec.nameSet, oldName)

	return nil
}

func (ec *entityCollection[E]) getEntityByID(id EntityID) (E, error) {
	e, ok := ec.entityMap[id]
	if ok {
		return e, nil
	}
	return e, fmt.Errorf(`id "%s" not found`, id)
}

func (ec *entityCollection[E]) removeEntity(id EntityID) error {
	e, err := ec.getEntityByID(id)
	if err != nil {
		return err
	}

	delete(ec.nameSet, e.getName())
	delete(ec.entityMap, id)

	return nil
}

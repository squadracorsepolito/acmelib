package acmelib

import (
	"fmt"

	"golang.org/x/exp/maps"
)

type collectableEntity interface {
	getEntityID() EntityID
	getName() string
}

type entityCollection[E collectableEntity] struct {
	entities    map[EntityID]E
	entityNames map[string]EntityID
}

func newEntityCollection[E collectableEntity]() *entityCollection[E] {
	return &entityCollection[E]{
		entities:    make(map[EntityID]E),
		entityNames: make(map[string]EntityID),
	}
}

func (ec *entityCollection[E]) getSize() int {
	return len(ec.entities)
}

func (ec *entityCollection[E]) listEntities() []E {
	return maps.Values(ec.entities)
}

func (ec *entityCollection[E]) verifyEntityName(name string) error {
	if _, ok := ec.entityNames[name]; ok {
		return fmt.Errorf(`duplicated name "%s"`, name)
	}
	return nil
}

func (ec *entityCollection[E]) addEntity(entity E) error {
	name := entity.getName()
	if err := ec.verifyEntityName(name); err != nil {
		return err
	}

	id := entity.getEntityID()
	if _, ok := ec.entities[id]; ok {
		return fmt.Errorf(`duplicated id "%s"`, id)
	}

	ec.entities[id] = entity
	ec.entityNames[name] = id

	return nil
}

func (ec *entityCollection[E]) removeEntity(id EntityID) error {
	e, err := ec.getEntityByID(id)
	if err != nil {
		return err
	}

	delete(ec.entityNames, e.getName())
	delete(ec.entities, id)

	return nil
}

func (ec *entityCollection[E]) updateEntityName(id EntityID, oldName, newName string) error {
	if oldName == newName {
		return fmt.Errorf(`"%s" is not a new name`, newName)
	}

	if err := ec.verifyEntityName(newName); err != nil {
		return err
	}

	ec.entityNames[newName] = id
	delete(ec.entityNames, oldName)

	return nil
}

func (ec *entityCollection[E]) getEntityByID(id EntityID) (e E, err error) {
	e, ok := ec.entities[id]
	if ok {
		return e, nil
	}
	return e, fmt.Errorf(`id "%s" not found`, id)
}

func (ec *entityCollection[E]) getEntityByName(name string) (e E, err error) {
	id, ok := ec.entityNames[name]
	if !ok {
		return e, fmt.Errorf(`name "%s" not found`, name)
	}

	return ec.getEntityByID(id)
}

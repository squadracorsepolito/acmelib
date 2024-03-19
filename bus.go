package acmelib

import "fmt"

type BusSortingMethod string

const (
	BusesByName       BusSortingMethod = "buses_by_name"
	BusesByCreatedAt  BusSortingMethod = "buses_by_created_at"
	BusesByModifiedAt BusSortingMethod = "buses_by_modified_at"
)

var busesSorter = newEntitySorter(
	newEntitySorterMethod(BusesByName, func(buses []*Bus) []*Bus { return sortByName(buses) }),
	newEntitySorterMethod(BusesByCreatedAt, func(buses []*Bus) []*Bus { return sortByCreateTime(buses) }),
	newEntitySorterMethod(BusesByModifiedAt, func(buses []*Bus) []*Bus { return sortByUpdateTime(buses) }),
)

func SelectBusesSortingMethod(method BusSortingMethod) {
	busesSorter.selectSortingMethod(method)
}

type Bus struct {
	*entity
	ParentProject *Project

	nodes *entityCollection[*Node]
}

func NewBus(name, desc string) *Bus {
	return &Bus{
		entity: newEntity(name, desc),

		nodes: newEntityCollection[*Node](),
	}
}

func (b *Bus) errorf(err error) error {
	return b.ParentProject.errorf(fmt.Errorf("bus %s: %v", b.Name, err))
}

func (b *Bus) UpdateName(name string) error {
	if err := b.ParentProject.buses.updateName(b.ID, b.Name, name); err != nil {
		return b.errorf(err)
	}

	return b.entity.UpdateName(name)
}

func (b *Bus) AddNode(node *Node) error {
	if err := b.nodes.addEntity(node); err != nil {
		return b.errorf(err)
	}

	node.ParentNode = b
	b.setUpdateTimeNow()

	return nil
}

func (b *Bus) ListNodes() []*Node {
	return nodesSorter.sortEntities(b.nodes.listEntities())
}

func (b *Bus) RemoveNode(nodeID EntityID) error {
	if err := b.nodes.removeEntity(nodeID); err != nil {
		return b.errorf(err)
	}

	b.setUpdateTimeNow()

	return nil
}

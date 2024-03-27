package acmelib

import "fmt"

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
	busErr := fmt.Errorf(`bus "%s": %v`, b.name, err)
	if b.ParentProject != nil {
		return b.ParentProject.errorf(busErr)
	}
	return busErr
}

func (b *Bus) UpdateName(name string) error {
	if b.ParentProject != nil {
		if err := b.ParentProject.buses.updateEntityName(b.entityID, b.name, name); err != nil {
			return b.errorf(err)
		}
	}

	return b.entity.UpdateName(name)
}

func (b *Bus) AddNode(node *Node) error {
	if err := b.nodes.addEntity(node); err != nil {
		return b.errorf(err)
	}

	node.ParentBus = b

	return nil
}

func (b *Bus) ListNodes() []*Node {
	return b.nodes.listEntities()
}

func (b *Bus) RemoveNode(nodeID EntityID) error {
	if err := b.nodes.removeEntity(nodeID); err != nil {
		return b.errorf(err)
	}

	return nil
}

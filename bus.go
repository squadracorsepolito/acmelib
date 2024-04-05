package acmelib

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type Bus struct {
	*attributeEntity

	parentNetwork *Network

	nodes     *set[EntityID, *Node]
	nodeNames *set[string, EntityID]
	nodeIDs   *set[NodeID, EntityID]

	baudrate uint
}

func NewBus(name, desc string) *Bus {
	return &Bus{
		attributeEntity: newAttributeEntity(name, desc, AttributeRefKindBus),

		parentNetwork: nil,

		nodes:     newSet[EntityID, *Node]("node"),
		nodeNames: newSet[string, EntityID]("node name"),
		nodeIDs:   newSet[NodeID, EntityID]("node id"),

		baudrate: 0,
	}
}

func (b *Bus) hasParentNetwork() bool {
	return b.parentNetwork != nil
}

func (b *Bus) setParentNetwork(net *Network) {
	b.parentNetwork = net
}

func (b *Bus) errorf(err error) error {
	busErr := fmt.Errorf(`bus "%s": %w`, b.name, err)
	if b.hasParentNetwork() {
		return b.parentNetwork.errorf(busErr)
	}
	return busErr
}

func (b *Bus) modifyNodeName(nodeEntID EntityID, newName string) {
	node, err := b.nodes.getValue(nodeEntID)
	if err != nil {
		panic(err)
	}

	oldName := node.name
	b.nodeNames.modifyKey(oldName, newName, nodeEntID)
}

func (b *Bus) UpdateName(newName string) error {
	if b.name == newName {
		return nil
	}

	if b.hasParentNetwork() {
		if err := b.parentNetwork.busNames.verifyKey(newName); err != nil {
			return b.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}

		b.parentNetwork.modifyBusName(b.entityID, newName)
	}

	b.name = newName

	return nil
}

func (b *Bus) AddNode(node *Node) error {
	if err := b.nodeNames.verifyKey(node.name); err != nil {
		return b.errorf(fmt.Errorf(`cannot add node "%s" : %w`, node.name, err))
	}

	if err := b.nodeIDs.verifyKey(node.id); err != nil {
		return b.errorf(fmt.Errorf(`cannot add node "%s" : %w`, node.name, err))
	}

	node.parentBuses.add(b.entityID, b)

	b.nodes.add(node.entityID, node)
	b.nodeNames.add(node.name, node.entityID)
	b.nodeIDs.add(node.id, node.entityID)

	return nil
}

func (b *Bus) RemoveNode(nodeEntityID EntityID) error {
	node, err := b.nodes.getValue(nodeEntityID)
	if err != nil {
		return b.errorf(fmt.Errorf(`cannot remove node with entity id "%s" : %w`, nodeEntityID, err))
	}

	node.parentBuses.remove(b.entityID)

	b.nodes.remove(nodeEntityID)
	b.nodeNames.remove(node.name)
	b.nodeIDs.remove(node.id)

	return nil
}

func (b *Bus) RemoveAllNodes() {
	for _, tmpNode := range b.nodes.entries() {
		tmpNode.parentBuses.remove(b.entityID)
	}

	b.nodes.clear()
	b.nodeNames.clear()
	b.nodeIDs.clear()
}

func (b *Bus) Nodes() []*Node {
	nodeSlice := b.nodes.getValues()
	slices.SortFunc(nodeSlice, func(a, b *Node) int { return int(a.id) - int(b.id) })
	return nodeSlice
}

func (b *Bus) SetBaudrate(baudrate uint) {
	b.baudrate = baudrate
}

func (b *Bus) Baudrate() uint {
	return b.baudrate
}

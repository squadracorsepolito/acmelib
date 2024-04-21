package acmelib

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// Bus is the virtual representation of physical CAN bus cable.
// It holds a list of nodes that are connected to it.
type Bus struct {
	*attributeEntity

	parentNetwork *Network

	nodes     *set[EntityID, *Node]
	nodeNames *set[string, EntityID]
	nodeIDs   *set[NodeID, EntityID]

	baudrate uint
}

// NewBus creates a new [Bus] with the given name and description.
func NewBus(name string) *Bus {
	return &Bus{
		attributeEntity: newAttributeEntity(name, AttributeRefKindBus),

		parentNetwork: nil,

		nodes:     newSet[EntityID, *Node](),
		nodeNames: newSet[string, EntityID](),
		nodeIDs:   newSet[NodeID, EntityID](),

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

func (b *Bus) stringify(builder *strings.Builder, tabs int) {
	b.entity.stringify(builder, tabs)

	tabStr := getTabString(tabs)

	builder.WriteString(fmt.Sprintf("%sbaudrate: %d\n", tabStr, b.baudrate))

	if b.nodes.size() == 0 {
		return
	}

	builder.WriteString(fmt.Sprintf("%snodes:\n", tabStr))
	for _, node := range b.Nodes() {
		node.stringify(builder, tabs+1)
		builder.WriteRune('\n')
	}
}

func (b *Bus) String() string {
	builder := new(strings.Builder)
	b.stringify(builder, 0)
	return builder.String()
}

// UpdateName updates the name of the [Bus].
// It may return an error if the new name is already in use within a network.
func (b *Bus) UpdateName(newName string) error {
	if b.name == newName {
		return nil
	}

	if b.hasParentNetwork() {
		if err := b.parentNetwork.busNames.verifyKeyUnique(newName); err != nil {
			return b.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}

		b.parentNetwork.modifyBusName(b.entityID, newName)
	}

	b.name = newName

	return nil
}

// ParentNetwork returns the [Network] that the [Bus] is part of.
// If the [Bus] is not part of a [Network], it returns nil.
func (b *Bus) ParentNetwork() *Network {
	return b.parentNetwork
}

// AddNode adds the given [Node] to the [Bus].
// It may return an error if the node name or the node id is already used by the bus.
func (b *Bus) AddNode(node *Node) error {
	if err := b.nodeNames.verifyKeyUnique(node.name); err != nil {
		return b.errorf(fmt.Errorf(`cannot add node "%s" : %w`, node.name, err))
	}

	if err := b.nodeIDs.verifyKeyUnique(node.id); err != nil {
		return b.errorf(fmt.Errorf(`cannot add node "%s" : %w`, node.name, err))
	}

	node.parentBuses.add(b.entityID, b)

	b.nodes.add(node.entityID, node)
	b.nodeNames.add(node.name, node.entityID)
	b.nodeIDs.add(node.id, node.entityID)

	return nil
}

// RemoveNode removes a [Node] that matches the given entity id from the [Bus].
// It may return an error if the node with the given entity id is not part of the bus.
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

// RemoveAllNodes removes all nodes from the [Bus].
func (b *Bus) RemoveAllNodes() {
	for _, tmpNode := range b.nodes.entries() {
		tmpNode.parentBuses.remove(b.entityID)
	}

	b.nodes.clear()
	b.nodeNames.clear()
	b.nodeIDs.clear()
}

// Nodes returns a slice of all nodes in the [Bus] sorted by node id.
func (b *Bus) Nodes() []*Node {
	nodeSlice := b.nodes.getValues()
	slices.SortFunc(nodeSlice, func(a, b *Node) int { return int(a.id) - int(b.id) })
	return nodeSlice
}

// GetNodeByName returns the [Node] with the given name from the [Bus].
// It may return an error if the node with the given name is not part of the bus.
func (b *Bus) GetNodeByName(nodeName string) (*Node, error) {
	id, err := b.nodeNames.getValue(nodeName)
	if err != nil {
		return nil, b.errorf(fmt.Errorf(`cannot get node with name "%s" : %w`, nodeName, err))
	}

	node, err := b.nodes.getValue(id)
	if err != nil {
		return nil, b.errorf(fmt.Errorf(`cannot get node with name "%s" : %w`, nodeName, err))
	}

	return node, nil
}

// SetBaudrate sets the baudrate of the [Bus].
func (b *Bus) SetBaudrate(baudrate uint) {
	b.baudrate = baudrate
}

// Baudrate returns the baudrate of the [Bus].
func (b *Bus) Baudrate() uint {
	return b.baudrate
}

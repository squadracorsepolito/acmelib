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

	baudrate int
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
	busErr := &EntityError{
		Kind:     "bus",
		EntityID: b.entityID,
		Name:     b.name,
		Err:      err,
	}

	if b.hasParentNetwork() {
		return b.parentNetwork.errorf(busErr)
	}

	return busErr
}

func (b *Bus) verifyNodeName(name string) error {
	err := b.nodeNames.verifyKeyUnique(name)
	if err != nil {
		return &NameError{
			Name: name,
			Err:  err,
		}
	}
	return nil
}

func (b *Bus) verifyNodeID(nodeID NodeID) error {
	err := b.nodeIDs.verifyKeyUnique(nodeID)
	if err != nil {
		return &NodeIDError{
			NodeID: nodeID,
			Err:    err,
		}
	}
	return nil
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
			return b.errorf(&UpdateNameError{Err: err})
		}

		b.parentNetwork.busNames.modifyKey(b.name, newName, b.entityID)
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
	addNodeErr := &AddEntityError{
		EntityID: node.entityID,
		Name:     node.name,
	}

	if err := b.verifyNodeName(node.name); err != nil {
		addNodeErr.Err = err
		return b.errorf(addNodeErr)
	}

	if err := b.verifyNodeID(node.id); err != nil {
		addNodeErr.Err = err
		return b.errorf(addNodeErr)
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
		return b.errorf(&RemoveEntityError{
			EntityID: nodeEntityID,
			Err:      err,
		})
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
		return nil, b.errorf(&GetEntityError{
			Err: &NameError{Err: err},
		})
	}

	node, err := b.nodes.getValue(id)
	if err != nil {
		return nil, b.errorf(&GetEntityError{
			EntityID: id,
			Err:      err,
		})
	}

	return node, nil
}

// SetBaudrate sets the baudrate of the [Bus].
func (b *Bus) SetBaudrate(baudrate int) {
	b.baudrate = baudrate
}

// Baudrate returns the baudrate of the [Bus].
func (b *Bus) Baudrate() int {
	return b.baudrate
}

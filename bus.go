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

	canIDBuilder *CANIDBuilder

	nodeInts  *set[EntityID, *NodeInterface]
	nodeNames *set[string, EntityID]
	nodeIDs   *set[NodeID, EntityID]

	baudrate int
}

// NewBus creates a new [Bus] with the given name and description.
func NewBus(name string) *Bus {
	return &Bus{
		attributeEntity: newAttributeEntity(name, AttributeRefKindBus),

		parentNetwork: nil,

		canIDBuilder: defaulCANIDBuilder,

		nodeInts:  newSet[EntityID, *NodeInterface](),
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

	b.canIDBuilder.stringify(builder, tabs)

	if b.nodeInts.size() == 0 {
		return
	}

	builder.WriteString(fmt.Sprintf("%sattached_node_interfaces:\n", tabStr))
	for _, nodeInt := range b.Nodes() {
		nodeInt.stringify(builder, tabs+1)
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

// AddNodeInterface adds a [NodeInterface] to the [Bus].
//
// It returns an [ArgumentError] if the given node interface is nil
// or a [NameError]/[NodeIDError] if the node name/id is already used.
func (b *Bus) AddNodeInterface(nodeInterface *NodeInterface) error {
	if nodeInterface == nil {
		return &ArgumentError{
			Name: "nodeInterface",
			Err:  ErrIsNil,
		}
	}

	addNodeIntErr := &AddEntityError{
		EntityID: nodeInterface.entityID,
		Name:     nodeInterface.name,
	}

	node := nodeInterface.node

	if err := b.verifyNodeName(node.name); err != nil {
		addNodeIntErr.Err = err
		return b.errorf(addNodeIntErr)
	}

	if err := b.verifyNodeID(node.id); err != nil {
		addNodeIntErr.Err = err
		return b.errorf(addNodeIntErr)
	}

	nodeInterface.parentBus = b

	b.nodeInts.add(nodeInterface.entityID, nodeInterface)
	b.nodeNames.add(node.name, nodeInterface.entityID)
	b.nodeIDs.add(node.id, nodeInterface.entityID)

	return nil
}

// RemoveNodeInterface removes a [NodeInterface] from the [Bus].
//
// It returns an [ErrNotFound] if the given entity id does not match
// any node interface.
func (b *Bus) RemoveNodeInterface(nodeInterfaceEntityID EntityID) error {
	nodeInt, err := b.nodeInts.getValue(nodeInterfaceEntityID)
	if err != nil {
		return b.errorf(&RemoveEntityError{
			EntityID: nodeInterfaceEntityID,
			Err:      err,
		})
	}

	nodeInt.parentBus = nil

	b.nodeInts.remove(nodeInterfaceEntityID)
	b.nodeNames.remove(nodeInt.node.name)
	b.nodeIDs.remove(nodeInt.node.id)

	return nil
}

// RemoveAllNodeInterfaces removes all node interfaces from the [Bus].
func (b *Bus) RemoveAllNodeInterfaces() {
	for _, tmpNodeInt := range b.nodeInts.entries() {
		tmpNodeInt.parentBus = nil
	}

	b.nodeInts.clear()
	b.nodeNames.clear()
	b.nodeIDs.clear()
}

// Nodes returns a slice of all nodes in the [Bus] sorted by node id.
func (b *Bus) Nodes() []*NodeInterface {
	nodeSlice := b.nodeInts.getValues()
	slices.SortFunc(nodeSlice, func(a, b *NodeInterface) int { return int(a.node.id) - int(b.node.id) })
	return nodeSlice
}

// GetNodeInterfaceByNodeName returns the [NodeInterface] with the given node name.
//
// It returns an [ErrNotFound] wrapped by a [NameError]
// if the node name does not match any node interface.
func (b *Bus) GetNodeInterfaceByNodeName(nodeName string) (*NodeInterface, error) {
	id, err := b.nodeNames.getValue(nodeName)
	if err != nil {
		return nil, b.errorf(&GetEntityError{
			Err: &NameError{Err: err},
		})
	}

	nodeInt, err := b.nodeInts.getValue(id)
	if err != nil {
		panic(err)
	}

	return nodeInt, nil
}

// SetBaudrate sets the baudrate of the [Bus].
func (b *Bus) SetBaudrate(baudrate int) {
	b.baudrate = baudrate
}

// Baudrate returns the baudrate of the [Bus].
func (b *Bus) Baudrate() int {
	return b.baudrate
}

// SetCANIDBuilder sets the [CANIDBuilder] of the [Bus].
func (b *Bus) SetCANIDBuilder(canIDBuilder *CANIDBuilder) {
	b.canIDBuilder = canIDBuilder
}

// CANIDBuilder returns the [CANIDBuilder] of the [Bus].
// If it is not set, it returns the default CAN-ID builder.
func (b *Bus) CANIDBuilder() *CANIDBuilder {
	return b.canIDBuilder
}

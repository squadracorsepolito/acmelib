package acmelib

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// BusType is the type of a [Bus].
type BusType int

const (
	// BusTypeCAN2A represents a CAN 2.0A bus.
	BusTypeCAN2A BusType = iota
)

func (bt BusType) String() string {
	switch bt {
	case BusTypeCAN2A:
		return "CAN_2.0A"
	default:
		return "unknown"
	}
}

// Bus is the virtual representation of physical CAN bus cable.
// It holds a list of nodes that are connected to it.
type Bus struct {
	*entity
	*withAttributes

	parentNetwork *Network

	canIDBuilder      *CANIDBuilder
	isDefCANIDBuilder bool

	nodeInts  *set[EntityID, *NodeInterface]
	nodeNames *set[string, EntityID]
	nodeIDs   *set[NodeID, EntityID]

	messageStaticCANIDs *set[CANID, EntityID]

	baudrate int
	typ      BusType
}

func newBusFromEntity(ent *entity) *Bus {
	return &Bus{
		entity:         ent,
		withAttributes: newWithAttributes(),

		parentNetwork: nil,

		canIDBuilder:      defaulCANIDBuilder,
		isDefCANIDBuilder: true,

		nodeInts:  newSet[EntityID, *NodeInterface](),
		nodeNames: newSet[string, EntityID](),
		nodeIDs:   newSet[NodeID, EntityID](),

		messageStaticCANIDs: newSet[CANID, EntityID](),

		baudrate: 0,
		typ:      BusTypeCAN2A,
	}
}

// NewBus creates a new [Bus] with the given name and description.
// By default, the bus is set to be of type CAN 2.0A.
func NewBus(name string) *Bus {
	return newBusFromEntity(newEntity(name, EntityKindBus))
}

func (b *Bus) hasParentNetwork() bool {
	return b.parentNetwork != nil
}

func (b *Bus) setParentNetwork(net *Network) {
	b.parentNetwork = net
}

func (b *Bus) errorf(err error) error {
	busErr := &EntityError{
		Kind:     EntityKindBus,
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

func (b *Bus) verifyStaticCANID(staticCANID CANID) error {
	if err := b.messageStaticCANIDs.verifyKeyUnique(staticCANID); err != nil {
		return &CANIDError{
			CANID: staticCANID,
			Err:   err,
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
	for _, nodeInt := range b.NodeInterfaces() {
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

	node := nodeInterface.node

	if err := b.verifyNodeName(node.name); err != nil {
		return b.errorf(err)
	}

	if err := b.verifyNodeID(node.id); err != nil {
		return b.errorf(err)
	}

	messages := nodeInterface.sentMessages.getValues()
	msgStaticCANIDs := make(map[CANID]EntityID)
	for _, tmpMsg := range messages {
		if tmpMsg.hasStaticCANID {
			err := b.verifyStaticCANID(tmpMsg.staticCANID)
			if err != nil {
				return b.errorf(err)
			}

			msgStaticCANIDs[tmpMsg.staticCANID] = tmpMsg.entityID
		}
	}
	for canID, entID := range msgStaticCANIDs {
		b.messageStaticCANIDs.add(canID, entID)
	}

	nodeInterface.parentBus = b
	nodeEntID := node.entityID

	b.nodeInts.add(nodeEntID, nodeInterface)
	b.nodeNames.add(node.name, nodeEntID)
	b.nodeIDs.add(node.id, nodeEntID)

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

	for _, tmpMsg := range nodeInt.sentMessages.getValues() {
		if tmpMsg.hasStaticCANID {
			b.messageStaticCANIDs.remove(tmpMsg.staticCANID)
		}
	}

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
	b.messageStaticCANIDs.clear()
}

// NodeInterfaces returns a slice of all node interfaces connected to the [Bus] sorted by node id.
func (b *Bus) NodeInterfaces() []*NodeInterface {
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
	if b.canIDBuilder != nil {
		b.canIDBuilder.removeRef(b.entityID)
	}
	b.canIDBuilder = canIDBuilder
	b.isDefCANIDBuilder = false
}

// CANIDBuilder returns the [CANIDBuilder] of the [Bus].
// If it is not set, it returns the default CAN-ID builder.
func (b *Bus) CANIDBuilder() *CANIDBuilder {
	return b.canIDBuilder
}

// SetType sets the type of the [Bus].
func (b *Bus) SetType(typ BusType) {
	b.typ = typ
}

// Type returns the type of the [Bus].
func (b *Bus) Type() BusType {
	return b.typ
}

// AssignAttribute assigns the given attribute/value pair to the [Bus].
//
// It returns an [ArgumentError] if the attribute is nil,
// or an [AttributeValueError] if the value does not conform to the attribute.
func (b *Bus) AssignAttribute(attribute Attribute, value any) error {
	if err := b.addAttributeAssignment(attribute, b, value); err != nil {
		return b.errorf(err)
	}
	return nil
}

// RemoveAttributeAssignment removes the [AttributeAssignment]
// with the given attribute entity id from the [Bus].
//
// It returns an [ErrNotFound] if the provided attribute entity id is not found.
func (b *Bus) RemoveAttributeAssignment(attributeEntityID EntityID) error {
	if err := b.removeAttributeAssignment(attributeEntityID); err != nil {
		return b.errorf(err)
	}
	return nil
}

// GetAttributeAssignment returns the [AttributeAssignment]
// with the given attribute entity id from the [Bus].
//
// It returns an [ErrNotFound] if the provided attribute entity id is not found.
func (b *Bus) GetAttributeAssignment(attributeEntityID EntityID) (*AttributeAssignment, error) {
	attAss, err := b.getAttributeAssignment(attributeEntityID)
	if err != nil {
		return nil, b.errorf(err)
	}
	return attAss, nil
}

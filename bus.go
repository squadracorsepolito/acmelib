package acmelib

import (
	"cmp"
	"slices"

	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
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

	nodeInts  *collection.Map[EntityID, *NodeInterface]
	nodeNames *collection.Map[string, EntityID]
	nodeIDs   *collection.Map[NodeID, EntityID]

	messageStaticCANIDs *collection.Map[CANID, EntityID]

	baudrate int
	typ      BusType
}

func newBusFromEntity(ent *entity) *Bus {
	bus := &Bus{
		entity:         ent,
		withAttributes: newWithAttributes(),

		parentNetwork: nil,

		isDefCANIDBuilder: true,

		nodeInts:  collection.NewMap[EntityID, *NodeInterface](),
		nodeNames: collection.NewMap[string, EntityID](),
		nodeIDs:   collection.NewMap[NodeID, EntityID](),

		messageStaticCANIDs: collection.NewMap[CANID, EntityID](),

		baudrate: 0,
		typ:      BusTypeCAN2A,
	}

	builder := newDefaultCANIDBuilder()
	bus.canIDBuilder = builder
	builder.addRef(bus)

	return bus
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
	if b.nodeNames.Has(name) {
		return newNameError(name, ErrIsDuplicated)
	}
	return nil
}

func (b *Bus) verifyNodeID(nodeID NodeID) error {
	if b.nodeIDs.Has(nodeID) {
		return newNodeIDError(nodeID, ErrIsDuplicated)
	}
	return nil
}

func (b *Bus) verifyStaticCANID(staticCANID CANID) error {
	if b.messageStaticCANIDs.Has(staticCANID) {
		return newCANIDError(staticCANID, ErrIsDuplicated)
	}
	return nil
}

func (b *Bus) verifyMessageSize(sizeByte int) error {
	switch b.typ {
	case BusTypeCAN2A:
		if sizeByte <= 8 {
			return nil
		}
	}

	return newSizeError(sizeByte, ErrTooBig)
}

func (b *Bus) stringify(s *stringer.Stringer) {
	b.entity.stringify(s)

	s.Write("baudrate: %d\n", b.baudrate)

	b.canIDBuilder.stringify(s)

	if b.nodeInts.Size() == 0 {
		return
	}

	b.withAttributes.stringify(s)
}

func (b *Bus) String() string {
	s := stringer.New()
	s.Write("bus:\n")
	b.stringify(s)
	return s.String()
}

// UpdateName updates the name of the [Bus].
// It may return an error if the new name is already in use within a network.
func (b *Bus) UpdateName(newName string) error {
	if b.name == newName {
		return nil
	}

	if b.hasParentNetwork() {
		if err := b.parentNetwork.verifyBusName(newName); err != nil {
			return b.errorf(err)
		}

		b.parentNetwork.busNames.Delete(b.name)
		b.parentNetwork.busNames.Set(newName, b.entityID)
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
// It returns:
//   - [ArgError] if the given node interface is nil.
//   - [NameError] if the node name is invalid.
//   - [NodeIDError] if the node id is invalid.
//   - [MessageSizeError] if one of the size of a message sent by the node is invalid.
//   - [CANIDError] if one of the static CAN-ID of a message sent by the node is invalid.
func (b *Bus) AddNodeInterface(nodeInterface *NodeInterface) error {
	if nodeInterface == nil {
		return &ArgError{
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

	msgStaticCANIDs := make(map[CANID]EntityID)
	for tmpMsg := range nodeInterface.sentMessages.Values() {
		if err := b.verifyMessageSize(tmpMsg.sizeByte); err != nil {
			return b.errorf(err)
		}

		if tmpMsg.hasStaticCANID {
			err := b.verifyStaticCANID(tmpMsg.staticCANID)
			if err != nil {
				return b.errorf(err)
			}

			msgStaticCANIDs[tmpMsg.staticCANID] = tmpMsg.entityID
		}
	}
	for canID, entID := range msgStaticCANIDs {
		b.messageStaticCANIDs.Set(canID, entID)
	}

	nodeInterface.parentBus = b
	nodeEntID := node.entityID

	b.nodeInts.Set(nodeEntID, nodeInterface)
	b.nodeNames.Set(node.name, nodeEntID)
	b.nodeIDs.Set(node.id, nodeEntID)

	return nil
}

// RemoveNodeInterface removes a [NodeInterface] from the [Bus].
//
// It returns an [ErrNotFound] if the given entity id does not match
// any node interface.
func (b *Bus) RemoveNodeInterface(nodeInterfaceEntityID EntityID) error {
	nodeInt, ok := b.nodeInts.Get(nodeInterfaceEntityID)
	if !ok {
		return b.errorf(ErrNotFound)
	}

	nodeInt.parentBus = nil

	b.nodeInts.Delete(nodeInterfaceEntityID)
	b.nodeNames.Delete(nodeInt.node.name)
	b.nodeIDs.Delete(nodeInt.node.id)

	for tmpMsg := range nodeInt.sentMessages.Values() {
		if tmpMsg.hasStaticCANID {
			b.messageStaticCANIDs.Delete(tmpMsg.staticCANID)
		}
	}

	return nil
}

// RemoveAllNodeInterfaces removes all node interfaces from the [Bus].
func (b *Bus) RemoveAllNodeInterfaces() {
	for tmpNodeInt := range b.nodeInts.Values() {
		tmpNodeInt.parentBus = nil
	}

	b.nodeInts.Clear()
	b.nodeNames.Clear()
	b.nodeIDs.Clear()
	b.messageStaticCANIDs.Clear()
}

// NodeInterfaces returns a slice of all node interfaces connected to the [Bus] sorted by node id.
func (b *Bus) NodeInterfaces() []*NodeInterface {
	nodeSlice := slices.Collect(b.nodeInts.Values())
	slices.SortFunc(nodeSlice, func(a, b *NodeInterface) int {
		return cmp.Compare(a.node.id, b.node.id)
	})
	return nodeSlice
}

// GetNodeInterfaceByNodeName returns the [NodeInterface] with the given node name.
//
// It returns an [ErrNotFound] wrapped by a [NameError]
// if the node name does not match any node interface.
func (b *Bus) GetNodeInterfaceByNodeName(nodeName string) (*NodeInterface, error) {
	id, ok := b.nodeNames.Get(nodeName)
	if !ok {
		return nil, b.errorf(ErrNotFound)
	}

	nodeInt, ok := b.nodeInts.Get(id)
	if !ok {
		return nil, b.errorf(ErrNotFound)
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
// It returns an [ArgError] if the attribute is nil,
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

// ToBus returns the bus itself.
func (b *Bus) ToBus() (*Bus, error) {
	return b, nil
}

package acmelib

import (
	"fmt"

	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// NodeID is a unique identifier for a [Node].
type NodeID uint32

func (nid NodeID) String() string {
	return fmt.Sprintf("%d", nid)
}

// Node is the representation of an ECU or an electronic component capable
// to send messages over a [Bus] through one or more [NodeInterface].
// It holds a list of interfaces that can send messages on the bus.
type Node struct {
	*entity
	*withAttributes

	interfaces []*NodeInterface
	intErrNum  int

	id             NodeID
	interfaceCount int
}

func newNodeFromEntity(ent *entity, id NodeID, intCount int) *Node {
	node := &Node{
		entity:         ent,
		withAttributes: newWithAttributes(),

		interfaces: []*NodeInterface{},
		intErrNum:  -1,

		id:             id,
		interfaceCount: intCount,
	}

	node.interfaces = make([]*NodeInterface, intCount)
	for i := 0; i < intCount; i++ {
		node.interfaces[i] = newNodeInterface(i, node)
	}

	return node
}

// NewNode creates a new [Node] with the given name, id and count of interfaces.
// The id must be unique among all nodes within a bus.
func NewNode(name string, id NodeID, interfaceCount int) *Node {
	return newNodeFromEntity(newEntity(name, EntityKindNode), id, interfaceCount)
}

func (n *Node) errorf(err error) error {
	nodeErr := &EntityError{
		Kind:     EntityKindNode,
		EntityID: n.entityID,
		Name:     n.name,
		Err:      err,
	}

	if len(n.interfaces) > 0 {
		if n.intErrNum >= 0 {
			nodeInt := n.interfaces[n.intErrNum]
			n.intErrNum = -1
			return nodeInt.errorf(nodeErr)
		}
	}

	return nodeErr
}

func (n *Node) stringify(s *stringer.Stringer) {
	n.entity.stringify(s)

	s.Write("node_id: %s\n", n.id.String())

	n.withAttributes.stringify(s)
}

func (n *Node) String() string {
	s := stringer.New()
	s.Write("node:\n")
	n.stringify(s)
	return s.String()
}

// UpdateName updates the name of the [Node].
// By updating the name, also the name of the interfaces are updated.
//
// It may return a [NameError] that wraps the cause of the error.
func (n *Node) UpdateName(newName string) error {
	if n.name == newName {
		return nil
	}

	buses := []*Bus{}
	for _, tmpInt := range n.interfaces {
		if !tmpInt.hasParentBus() {
			continue
		}

		tmpBus := tmpInt.parentBus
		if err := tmpBus.verifyNodeName(newName); err != nil {
			n.intErrNum = tmpInt.number
			return n.errorf(err)
		}

		buses = append(buses, tmpBus)
	}

	for _, tmpBus := range buses {
		tmpBus.nodeNames.Delete(n.name)
		tmpBus.nodeNames.Set(newName, n.entityID)
	}

	n.name = newName

	return nil
}

// UpdateID updates the id of the [Node].
//
// It returns a [NodeError] if the new id is invalid.
func (n *Node) UpdateID(newID NodeID) error {
	if newID == n.id {
		return nil
	}

	buses := []*Bus{}
	for _, tmpNodeInt := range n.interfaces {
		if !tmpNodeInt.hasParentBus() {
			continue
		}

		if err := tmpNodeInt.parentBus.verifyNodeID(newID); err != nil {
			return n.errorf(err)
		}

		buses = append(buses, tmpNodeInt.parentBus)
	}

	for _, tmpBus := range buses {
		tmpBus.nodeIDs.Delete(n.id)
		tmpBus.nodeIDs.Set(newID, n.entityID)
	}

	n.id = newID

	return nil
}

// AddInterface adds a new interface to the [Node].
// The interface will be assigned with the next available interface number.
func (n *Node) AddInterface() {
	n.interfaces = append(n.interfaces, newNodeInterface(n.interfaceCount, n))
	n.interfaceCount++
}

// RemoveInterface removes the interface with the given interface number from the [Node].
// It will update the interface numbers of the remaining interfaces
// in order to keep the interface numbers contiguous.
//
// It returns an [ArgError] if the interface number is negative or out of bounds.
func (n *Node) RemoveInterface(interfaceNumber int) error {
	if interfaceNumber < 0 {
		return &ArgError{
			Name: "interfaceNumber",
			Err:  ErrIsNegative,
		}
	}

	if interfaceNumber >= n.interfaceCount {
		return &ArgError{
			Name: "interfaceNumber",
			Err:  ErrOutOfBounds,
		}
	}

	found := false
	newInterfaces := make([]*NodeInterface, 0, n.interfaceCount-1)
	for _, tmpInt := range n.interfaces {
		if tmpInt.number == interfaceNumber {
			if tmpInt.hasParentBus() {
				if err := tmpInt.parentBus.RemoveNodeInterface(n.entityID); err != nil {
					return err
				}
			}

			found = true
			continue
		}

		if found {
			tmpInt.number--
		}

		newInterfaces = append(newInterfaces, tmpInt)
	}

	n.interfaceCount--

	return nil
}

// Interfaces returns a slice with all the interfaces of the [Node].
func (n *Node) Interfaces() []*NodeInterface {
	return n.interfaces
}

// GetInterface returns the [NodeInterface] with the given interface number.
// It retruns nil if the interface is not found.
func (n *Node) GetInterface(interfaceNumber int) *NodeInterface {
	if interfaceNumber < 0 {
		return nil
	}

	if interfaceNumber >= n.interfaceCount {
		return nil
	}

	return n.interfaces[interfaceNumber]
}

// ID returns the id of the [Node].
func (n *Node) ID() NodeID {
	return n.id
}

// AssignAttribute assigns the given attribute/value pair to the [Node].
//
// It returns an [ArgError] if the attribute is nil,
// or an [AttributeValueError] if the value does not conform to the attribute.
func (n *Node) AssignAttribute(attribute Attribute, value any) error {
	if err := n.addAttributeAssignment(attribute, n, value); err != nil {
		return n.errorf(err)
	}
	return nil
}

// RemoveAttributeAssignment removes the [AttributeAssignment]
// with the given attribute entity id from the [Node].
//
// It returns an [ErrNotFound] if the provided attribute entity id is not found.
func (n *Node) RemoveAttributeAssignment(attributeEntityID EntityID) error {
	if err := n.removeAttributeAssignment(attributeEntityID); err != nil {
		return n.errorf(err)
	}
	return nil
}

// GetAttributeAssignment returns the [AttributeAssignment]
// with the given attribute entity id from the [Node].
//
// It returns an [ErrNotFound] if the provided attribute entity id is not found.
func (n *Node) GetAttributeAssignment(attributeEntityID EntityID) (*AttributeAssignment, error) {
	attAss, err := n.getAttributeAssignment(attributeEntityID)
	if err != nil {
		return nil, n.errorf(err)
	}
	return attAss, nil
}

// ToNode returns the node itself.
func (n *Node) ToNode() (*Node, error) {
	return n, nil
}

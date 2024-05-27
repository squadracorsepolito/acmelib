package acmelib

import (
	"fmt"
	"strings"
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
	// *attributeEntity
	*entity
	*withAttributes

	interfaces []*NodeInterface
	intErrNum  int

	id             NodeID
	interfaceCount int
}

// NewNode creates a new [Node] with the given name, id and count of interfaces.
// The id must be unique among all nodes within a bus.
func NewNode(name string, id NodeID, interfaceCount int) *Node {
	node := &Node{
		// attributeEntity: newAttributeEntity(name, AttributeRefKindNode),
		entity:         newEntity(name),
		withAttributes: newWithAttributes(),

		interfaces: []*NodeInterface{},
		intErrNum:  -1,

		id:             id,
		interfaceCount: interfaceCount,
	}

	node.interfaces = make([]*NodeInterface, interfaceCount)
	for i := 0; i < interfaceCount; i++ {
		node.interfaces[i] = newNodeInterface(i, node)
	}

	return node
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

func (n *Node) stringify(b *strings.Builder, tabs int) {
	n.entity.stringify(b, tabs)
	tabStr := getTabString(tabs)
	b.WriteString(fmt.Sprintf("%snode_id: %s\n", tabStr, n.id.String()))
}

func (n *Node) String() string {
	builder := new(strings.Builder)
	n.stringify(builder, 0)
	return builder.String()
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
	intEntIDs := []EntityID{}
	for _, tmpInt := range n.interfaces {
		if !tmpInt.hasParentBus() {
			continue
		}

		tmpBus := tmpInt.parentBus
		if err := tmpBus.verifyNodeName(newName); err != nil {
			n.intErrNum = tmpInt.number
			return n.errorf(&UpdateNameError{Err: err})
		}

		buses = append(buses, tmpBus)
		intEntIDs = append(intEntIDs, tmpInt.entityID)
	}

	for idx, tmpBus := range buses {
		tmpBus.nodeNames.modifyKey(n.name, newName, intEntIDs[idx])
	}

	for _, tmpInt := range n.interfaces {
		tmpInt.name = tmpInt.setName(newName)
	}

	n.name = newName

	return nil
}

// Interfaces returns a slice with all the interfaces of the [Node].
func (n *Node) Interfaces() []*NodeInterface {
	return n.interfaces
}

// ID returns the id of the [Node].
func (n *Node) ID() NodeID {
	return n.id
}

package acmelib

import (
	"fmt"
	"strings"
)

// NodeID is a unique identifier for a node.
// It must be manually assigned by the user.
type NodeID uint32

func (nid NodeID) String() string {
	return fmt.Sprintf("%d", nid)
}

// Node is the representation of an ECU.
// It holds a list of messages which are sent to other nodes thought the bus.
// A node can be assigned to more then 1 bus.
type Node struct {
	*attributeEntity

	// parentBuses *set[EntityID, *Bus]
	// parErrID    EntityID

	// messages     *set[EntityID, *Message]
	// messageNames *set[string, EntityID]
	// messageIDs   *set[MessageID, EntityID]

	interfaces []*NodeInterface

	id NodeID
}

// NewNode creates a new [Node] with the given name and id.
// The id must be unique among all nodes within a bus.
func NewNode(name string, id NodeID, interfaceCount int) *Node {
	node := &Node{
		attributeEntity: newAttributeEntity(name, AttributeRefKindNode),

		id: id,
	}

	node.interfaces = make([]*NodeInterface, interfaceCount)
	for i := 0; i < interfaceCount; i++ {
		node.interfaces[i] = newNodeInterface(i, node)
	}

	return node

	// return &Node{
	// 	attributeEntity: newAttributeEntity(name, AttributeRefKindNode),

	// 	interfaces:   newSet[EntityID, *NodeInterface](),
	// 	interfaceIDs: newSet[int, EntityID](),

	// 	// parentBuses: newSet[EntityID, *Bus](),
	// 	// parErrID:    "",

	// 	// messages:     newSet[EntityID, *Message](),
	// 	// messageNames: newSet[string, EntityID](),
	// 	// messageIDs:   newSet[MessageID, EntityID](),

	// 	id: id,
	// }
}

func (n *Node) errorf(err error) error {
	nodeIntErr := &EntityError{
		Kind:     EntityKindNode,
		EntityID: n.entityID,
		Name:     n.name,
		Err:      err,
	}
	return nodeIntErr
}

func (n *Node) Interfaces() []*NodeInterface {
	return n.interfaces
}

// func (n *Node) errorf(err error) error {
// 	nodeErr := &EntityError{
// 		Kind:     EntityKindNode,
// 		EntityID: n.entityID,
// 		Name:     n.name,
// 		Err:      err,
// 	}

// 	if n.parentBuses.size() > 0 {
// 		if n.parErrID != "" {
// 			parBus, err := n.parentBuses.getValue(n.parErrID)
// 			if err != nil {
// 				panic(err)
// 			}

// 			n.parErrID = ""
// 			return parBus.errorf(nodeErr)
// 		}

// 		return n.parentBuses.getValues()[0].errorf(nodeErr)
// 	}

// 	return nodeErr
// }

// func (n *Node) verifyMessageName(name string) error {
// 	err := n.messageNames.verifyKeyUnique(name)
// 	if err != nil {
// 		return &NameError{
// 			Name: name,
// 			Err:  err,
// 		}
// 	}
// 	return nil
// }

// func (n *Node) verifyMessageID(msgID MessageID) error {
// 	err := n.messageIDs.verifyKeyUnique(msgID)
// 	if err != nil {
// 		return &MessageIDError{
// 			MessageID: msgID,
// 			Err:       err,
// 		}
// 	}
// 	return nil
// }

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

func (n *Node) UpdateName(newName string) error {
	if n.name == newName {
		return nil
	}

	buses := []*Bus{}
	for _, nodeInt := range n.interfaces {
		if !nodeInt.hasParentBus() {
			continue
		}

		tmpBus := nodeInt.parentBus
		if err := tmpBus.verifyNodeName(newName); err != nil {
			return n.errorf(&UpdateNameError{Err: err})
		}

		buses = append(buses, tmpBus)
	}

	// buses := []*Bus{}
	// for _, tmpBus := range n.parentBuses.entries() {
	// 	if err := tmpBus.verifyNodeName(newName); err != nil {
	// 		return n.errorf(&UpdateNameError{Err: err})
	// 	}

	// 	buses = append(buses, tmpBus)
	// }

	for _, tmpBus := range buses {
		nodeIntEntID, err := tmpBus.nodeNames.getValue(n.name)
		if err != nil {
			panic(err)
		}
		tmpBus.nodeNames.modifyKey(n.name, newName, nodeIntEntID)
	}

	n.name = newName

	return nil
}

// // UpdateName updates the name of the [Node].
// // It may return an error if the new name is already in use within a bus.
// func (n *Node) UpdateName(newName string) error {
// 	if n.name == newName {
// 		return nil
// 	}

// 	buses := []*Bus{}
// 	for _, tmpBus := range n.parentBuses.entries() {
// 		if err := tmpBus.verifyNodeName(newName); err != nil {
// 			return n.errorf(&UpdateNameError{Err: err})
// 		}

// 		buses = append(buses, tmpBus)
// 	}

// 	for _, tmpBus := range buses {
// 		tmpBus.nodeNames.modifyKey(n.name, newName, n.entityID)
// 	}

// 	n.name = newName

// 	return nil
// }

// // ParentBuses returns a slice of [Bus]es that the [Node] is part of.
// func (n *Node) ParentBuses() []*Bus {
// 	return n.parentBuses.getValues()
// }

// // AddMessage adds a [Message] to the [Node].
// // This means that the given message will be sent by the node.
// // It may return an error if the message name or the message id is already used by the node.
// func (n *Node) AddMessage(message *Message) error {
// 	if message == nil {
// 		return &ArgumentError{
// 			Name: "message",
// 			Err:  ErrIsNil,
// 		}
// 	}

// 	addMsgErr := &AddEntityError{
// 		EntityID: message.entityID,
// 		Name:     message.name,
// 	}

// 	if err := n.verifyMessageName(message.name); err != nil {
// 		addMsgErr.Err = err
// 		return n.errorf(addMsgErr)
// 	}

// 	// message.generateID(n.messages.size()+1, n.id)

// 	if !message.hasStaticCANID {
// 		if err := n.verifyMessageID(message.id); err != nil {
// 			addMsgErr.Err = err
// 			return n.errorf(addMsgErr)
// 		}

// 		n.messageIDs.add(message.id, message.entityID)
// 	}

// 	n.messages.add(message.entityID, message)
// 	n.messageNames.add(message.name, message.entityID)

// 	message.senderNode = n

// 	return nil
// }

// // RemoveMessage removes a [Message] that matches the given entity id from the [Node].
// // It may return an error if the message with the given entity id is not sent by the node.
// func (n *Node) RemoveMessage(messageEntityID EntityID) error {
// 	msg, err := n.messages.getValue(messageEntityID)
// 	if err != nil {
// 		return n.errorf(&RemoveEntityError{
// 			EntityID: messageEntityID,
// 			Err:      err,
// 		})
// 	}

// 	msg.senderNode = nil
// 	// msg.resetID()

// 	n.messages.remove(messageEntityID)
// 	n.messageNames.remove(msg.name)
// 	n.messageIDs.remove(msg.id)
// 	// n.messageIDs.clear()
// 	// for idx, tmpMsg := range n.Messages() {
// 	// 	tmpMsg.generateID(idx+1, n.id)
// 	// 	n.messageIDs.add(tmpMsg.id, tmpMsg.entityID)
// 	// }

// 	return nil
// }

// // RemoveAllMessages removes all messages from the [Node].
// func (n *Node) RemoveAllMessages() {
// 	for _, tmpMsg := range n.messages.entries() {
// 		// tmpMsg.resetID()
// 		tmpMsg.senderNode = nil
// 	}

// 	n.messages.clear()
// 	n.messageNames.clear()
// 	n.messageIDs.clear()
// }

// // Messages returns a slice of messages that the [Node] sends sorted by message id.
// func (n *Node) Messages() []*Message {
// 	msgSlice := n.messages.getValues()
// 	slices.SortFunc(msgSlice, func(a, b *Message) int {
// 		// TODO! revisit this
// 		if n.parentBuses.size() > 0 {
// 			busID := n.parentBuses.getValues()[0].entityID
// 			return int(a.GetCANID(busID)) - int(b.GetCANID(busID))
// 		}
// 		return int(a.id) - int(b.id)
// 	})
// 	return msgSlice
// }

// ID returns the node id.
func (n *Node) ID() NodeID {
	return n.id
}

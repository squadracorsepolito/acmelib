package acmelib

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
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

	parentBuses *set[EntityID, *Bus]
	parErrID    EntityID

	messages     *set[EntityID, *Message]
	messageNames *set[string, EntityID]
	messageIDs   *set[MessageID, EntityID]

	id NodeID
}

// NewNode creates a new [Node] with the given name and id.
// The id must be unique among all nodes within a bus.
func NewNode(name string, id NodeID) *Node {
	return &Node{
		attributeEntity: newAttributeEntity(name, AttributeRefKindNode),

		parentBuses: newSet[EntityID, *Bus]("parent bus"),
		parErrID:    "",

		messages:     newSet[EntityID, *Message]("message"),
		messageNames: newSet[string, EntityID]("message name"),
		messageIDs:   newSet[MessageID, EntityID]("message id"),

		id: id,
	}
}

func (n *Node) modifyMessageName(msgEntID EntityID, newName string) error {
	msg, err := n.messages.getValue(msgEntID)
	if err != nil {
		return err
	}

	oldName := msg.name
	n.messageNames.modifyKey(oldName, newName, msgEntID)

	return nil
}

func (n *Node) errorf(err error) error {
	nodeErr := fmt.Errorf(`node "%s" : %w`, n.name, err)

	if n.parentBuses.size() > 0 {
		if n.parErrID != "" {
			parBus, err := n.parentBuses.getValue(n.parErrID)
			if err != nil {
				panic(err)
			}

			n.parErrID = ""
			return parBus.errorf(nodeErr)
		}

		return n.parentBuses.getValues()[0].errorf(nodeErr)
	}

	return nodeErr
}

func (n *Node) stringify(b *strings.Builder, tabs int) {
	n.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%snode_id: %s\n", tabStr, n.id.String()))

	if n.messages.size() == 0 {
		return
	}

	b.WriteString(fmt.Sprintf("%smessages:\n", tabStr))
	for _, msg := range n.Messages() {
		msg.stringify(b, tabs+1)
		b.WriteRune('\n')
	}
}

func (n *Node) String() string {
	builder := new(strings.Builder)
	n.stringify(builder, 0)
	return builder.String()
}

// UpdateName updates the name of the [Node].
// It may return an error if the new name is already in use within a bus.
func (n *Node) UpdateName(newName string) error {
	if n.name == newName {
		return nil
	}

	for _, tmpBus := range n.parentBuses.entries() {
		if err := tmpBus.nodeNames.verifyKey(newName); err != nil {
			n.parErrID = tmpBus.entityID
			return n.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}

		tmpBus.modifyNodeName(n.entityID, newName)
	}

	n.name = newName

	return nil
}

// ParentBuses returns a slice of [Bus]es that the [Node] is part of.
func (n *Node) ParentBuses() []*Bus {
	return n.parentBuses.getValues()
}

// AddMessage adds a [Message] to the [Node].
// This means that the given message will be sent by the node.
// It may return an error if the message name or the message id is already used by the node.
func (n *Node) AddMessage(message *Message) error {
	if err := n.messageNames.verifyKey(message.name); err != nil {
		return n.errorf(fmt.Errorf(`cannot add message "%s" : %w`, message.name, err))
	}

	message.generateID(n.messages.size()+1, n.id)

	if err := n.messageIDs.verifyKey(message.id); err != nil {
		return n.errorf(fmt.Errorf(`cannot add message "%s" : %w`, message.name, err))
	}

	n.messages.add(message.entityID, message)
	n.messageNames.add(message.name, message.entityID)
	n.messageIDs.add(message.id, message.entityID)

	message.parentNodes.add(n.entityID, n)

	return nil
}

// RemoveMessage removes a [Message] that matches the given entity id from the [Node].
// It may return an error if the message with the given entity id is not sent by the node.
func (n *Node) RemoveMessage(messageEntityID EntityID) error {
	msg, err := n.messages.getValue(messageEntityID)
	if err != nil {
		return n.errorf(fmt.Errorf(`cannot remove message with entity id "%s" : %w`, messageEntityID, err))
	}

	msg.parentNodes.remove(n.entityID)
	msg.resetID()

	n.messages.remove(messageEntityID)
	n.messageNames.remove(msg.name)
	n.messageIDs.clear()
	for idx, tmpMsg := range n.Messages() {
		tmpMsg.generateID(idx+1, n.id)
		n.messageIDs.add(tmpMsg.id, tmpMsg.entityID)
	}

	return nil
}

// RemoveAllMessages removes all messages from the [Node].
func (n *Node) RemoveAllMessages() {
	for _, tmpMsg := range n.messages.entries() {
		tmpMsg.resetID()
		tmpMsg.parentNodes.remove(n.entityID)
	}

	n.messages.clear()
	n.messageNames.clear()
	n.messageIDs.clear()
}

// Messages returns a slice of messages that the [Node] sends sorted by message id.
func (n *Node) Messages() []*Message {
	msgSlice := n.messages.getValues()
	slices.SortFunc(msgSlice, func(a, b *Message) int { return int(a.id) - int(b.id) })
	return msgSlice
}

// ID returns the node id.
func (n *Node) ID() NodeID {
	return n.id
}

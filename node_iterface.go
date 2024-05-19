package acmelib

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// NodeInterface represents an interface between a [Bus] and a [Node].
type NodeInterface struct {
	*entity

	parentBus *Bus

	messages     *set[EntityID, *Message]
	messageNames *set[string, EntityID]
	messageIDs   *set[MessageID, EntityID]

	number int
	node   *Node
}

func newNodeInterface(number int, node *Node) *NodeInterface {
	ni := &NodeInterface{
		parentBus: nil,

		messages:     newSet[EntityID, *Message](),
		messageNames: newSet[string, EntityID](),
		messageIDs:   newSet[MessageID, EntityID](),

		number: number,
		node:   node,
	}

	ni.entity = newEntity(ni.setName(node.name))

	return ni
}

func (ni *NodeInterface) setName(nodeName string) string {
	return fmt.Sprintf("%s/int%d", nodeName, ni.number)
}

func (ni *NodeInterface) hasParentBus() bool {
	return ni.parentBus != nil
}

func (ni *NodeInterface) errorf(err error) error {
	nodeIntErr := &EntityError{
		Kind:     EntityKindNode,
		EntityID: ni.entityID,
		Name:     ni.name,
		Err:      err,
	}

	if ni.hasParentBus() {
		return ni.parentBus.errorf(nodeIntErr)
	}

	return nodeIntErr
}

func (ni *NodeInterface) stringify(b *strings.Builder, tabs int) {
	ni.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%snumber: %d\n", tabStr, ni.number))

	b.WriteString(fmt.Sprintf("%snode:\n", tabStr))
	ni.node.stringify(b, tabs+1)

	if ni.messages.size() == 0 {
		return
	}

	b.WriteString(fmt.Sprintf("%ssended_messages:\n", tabStr))
	for _, msg := range ni.Messages() {
		msg.stringify(b, tabs+1)
		b.WriteRune('\n')
	}
}

func (ni *NodeInterface) String() string {
	builder := new(strings.Builder)
	ni.stringify(builder, 0)
	return builder.String()
}

func (ni *NodeInterface) verifyMessageName(name string) error {
	err := ni.messageNames.verifyKeyUnique(name)
	if err != nil {
		return &NameError{
			Name: name,
			Err:  err,
		}
	}
	return nil
}

func (ni *NodeInterface) verifyMessageID(msgID MessageID) error {
	err := ni.messageIDs.verifyKeyUnique(msgID)
	if err != nil {
		return &MessageIDError{
			MessageID: msgID,
			Err:       err,
		}
	}
	return nil
}

// AddMessage adds a [Message] that the [NodeInterface] can send.
//
// It returns an [ArgumentError] if the given message is nil or
// a [NameError]/[MessageIDError] if the message name/id is already used.
func (ni *NodeInterface) AddMessage(message *Message) error {
	if message == nil {
		return &ArgumentError{
			Name: "message",
			Err:  ErrIsNil,
		}
	}

	addMsgErr := &AddEntityError{
		EntityID: message.entityID,
		Name:     message.name,
	}

	if err := ni.verifyMessageName(message.name); err != nil {
		addMsgErr.Err = err
		return ni.errorf(addMsgErr)
	}

	if !message.hasStaticCANID {
		if err := ni.verifyMessageID(message.id); err != nil {
			addMsgErr.Err = err
			return ni.errorf(addMsgErr)
		}

		ni.messageIDs.add(message.id, message.entityID)
	}

	ni.messages.add(message.entityID, message)
	ni.messageNames.add(message.name, message.entityID)

	message.senderNodeInt = ni

	return nil
}

// RemoveMessage removes a [Message] sent by the [NodeInterface].
//
// It returns an [ErrNotFound] if the given entity id does not match
// any message.
func (ni *NodeInterface) RemoveMessage(messageEntityID EntityID) error {
	msg, err := ni.messages.getValue(messageEntityID)
	if err != nil {
		return ni.errorf(&RemoveEntityError{
			EntityID: messageEntityID,
			Err:      err,
		})
	}

	msg.senderNodeInt = nil

	ni.messages.remove(messageEntityID)
	ni.messageNames.remove(msg.name)
	ni.messageIDs.remove(msg.id)

	return nil
}

// RemoveAllMessages removes all the messages sent by the [NodeInterface].
func (ni *NodeInterface) RemoveAllMessages() {
	for _, tmpMsg := range ni.messages.entries() {
		tmpMsg.senderNodeInt = nil
	}

	ni.messages.clear()
	ni.messageNames.clear()
	ni.messageIDs.clear()
}

// Messages returns a slice of messages sended by the [NodeInterface].
func (ni *NodeInterface) Messages() []*Message {
	msgSlice := ni.messages.getValues()
	slices.SortFunc(msgSlice, func(a, b *Message) int {
		return int(a.id) - int(b.id)
	})
	return msgSlice
}

// Node returns the [Node] that owns the [NodeInterface].
func (ni *NodeInterface) Node() *Node {
	return ni.node
}

// ParentBus returns the [Bus] attached to the [NodeInterface].
func (ni *NodeInterface) ParentBus() *Bus {
	return ni.parentBus
}

// Number returns the number of the [NodeInterface].
// The number is unique among all the interfaces within a [Node]
// and it cannot be manually assigned.
func (ni *NodeInterface) Number() int {
	return ni.number
}

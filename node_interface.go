package acmelib

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// NodeInterface represents an interface between a [Bus] and a [Node].
type NodeInterface struct {
	parentBus *Bus

	sentMessages            *set[EntityID, *Message]
	sentMessageNames        *set[string, EntityID]
	sentMessageIDs          *set[MessageID, EntityID]
	sentMessageStaticCANIDs *set[CANID, EntityID]

	receivedMessages *set[EntityID, *Message]

	number int
	node   *Node
}

func newNodeInterface(number int, node *Node) *NodeInterface {
	return &NodeInterface{
		parentBus: nil,

		sentMessages:            newSet[EntityID, *Message](),
		sentMessageNames:        newSet[string, EntityID](),
		sentMessageIDs:          newSet[MessageID, EntityID](),
		sentMessageStaticCANIDs: newSet[CANID, EntityID](),

		receivedMessages: newSet[EntityID, *Message](),

		number: number,
		node:   node,
	}
}

func (ni *NodeInterface) hasParentBus() bool {
	return ni.parentBus != nil
}

func (ni *NodeInterface) errorf(err error) error {
	if ni.hasParentBus() {
		return ni.parentBus.errorf(err)
	}
	return err
}

func (ni *NodeInterface) stringify(b *strings.Builder, tabs int) {
	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%snumber: %d\n", tabStr, ni.number))

	b.WriteString(fmt.Sprintf("%snode:\n", tabStr))
	ni.node.stringify(b, tabs+1)

	if ni.sentMessages.size() == 0 {
		return
	}

	b.WriteString(fmt.Sprintf("%ssent_messages:\n", tabStr))
	for _, msg := range ni.SentMessages() {
		msg.stringify0(b, tabs+1)
		b.WriteRune('\n')
	}

	b.WriteString(fmt.Sprintf("%sreceived_messages:\n", tabStr))
	for _, msg := range ni.ReceivedMessages() {
		msg.stringify0(b, tabs+1)
		b.WriteRune('\n')
	}
}

func (ni *NodeInterface) String() string {
	builder := new(strings.Builder)
	ni.stringify(builder, 0)
	return builder.String()
}

func (ni *NodeInterface) verifyMessageName(name string) error {
	err := ni.sentMessageNames.verifyKeyUnique(name)
	if err != nil {
		return &NameError{
			Name: name,
			Err:  err,
		}
	}
	return nil
}

func (ni *NodeInterface) verifyMessageID(msgID MessageID) error {
	err := ni.sentMessageIDs.verifyKeyUnique(msgID)
	if err != nil {
		return &MessageIDError{
			MessageID: msgID,
			Err:       err,
		}
	}
	return nil
}

func (ni *NodeInterface) verifyStaticCANID(staticCANID CANID) error {
	if err := ni.sentMessageStaticCANIDs.verifyKeyUnique(staticCANID); err != nil {
		return &CANIDError{
			CANID: staticCANID,
			Err:   err,
		}
	}

	if ni.hasParentBus() {
		return ni.parentBus.verifyStaticCANID(staticCANID)
	}

	return nil
}

func (ni *NodeInterface) verifyMessageSize(sizeByte int) error {
	if ni.hasParentBus() {
		return ni.parentBus.verifyMessageSize(sizeByte)
	}

	return nil
}

func (ni *NodeInterface) addReceivedMessage(msg *Message) error {
	if ni.sentMessages.hasKey(msg.entityID) {
		return ErrReceiverIsSender
	}

	ni.receivedMessages.add(msg.entityID, msg)
	msg.receivers.add(ni.node.entityID, ni)

	return nil
}

func (ni *NodeInterface) removeReceivedMessage(msg *Message) {
	ni.receivedMessages.remove(msg.entityID)
	msg.receivers.remove(ni.node.entityID)
}

// AddSentMessage adds a [Message] that the [NodeInterface] can send.
//
// It returns:
//   - [ArgError] if the given message is nil.
//   - [AddEntityError] that wraps a [NameError] if the message name is invalid.
//   - [AddEntityError] that wraps a [CANIDError] if the message static can id is invalid.
//   - [AddEntityError] that wraps a [MessageIDError] if the message id is invalid.
//   - [AddEntityError] that wraps a [MessageSizeError] if the message size is invalid.
func (ni *NodeInterface) AddSentMessage(message *Message) error {
	if message == nil {
		return &ArgError{
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

	if err := ni.verifyMessageSize(message.sizeByte); err != nil {
		addMsgErr.Err = err
		return ni.errorf(addMsgErr)
	}

	if message.hasStaticCANID {
		if err := ni.verifyStaticCANID(message.staticCANID); err != nil {
			addMsgErr.Err = err
			return ni.errorf(addMsgErr)
		}

		if ni.hasParentBus() {
			ni.parentBus.messageStaticCANIDs.add(message.staticCANID, message.entityID)
		}

		ni.sentMessageStaticCANIDs.add(message.staticCANID, message.entityID)

	} else {
		if err := ni.verifyMessageID(message.id); err != nil {
			addMsgErr.Err = err
			return ni.errorf(addMsgErr)
		}
		ni.sentMessageIDs.add(message.id, message.entityID)
	}

	ni.sentMessages.add(message.entityID, message)
	ni.sentMessageNames.add(message.name, message.entityID)

	message.senderNodeInt = ni

	return nil
}

// RemoveSentMessage removes a [Message] sent by the [NodeInterface].
//
// It returns an [ErrNotFound] if the given entity id does not match
// any message.
func (ni *NodeInterface) RemoveSentMessage(messageEntityID EntityID) error {
	msg, err := ni.sentMessages.getValue(messageEntityID)
	if err != nil {
		return ni.errorf(&RemoveEntityError{
			EntityID: messageEntityID,
			Err:      err,
		})
	}

	msg.senderNodeInt = nil

	ni.sentMessages.remove(messageEntityID)
	ni.sentMessageNames.remove(msg.name)

	if msg.hasStaticCANID {
		ni.sentMessageStaticCANIDs.remove(msg.staticCANID)

		if ni.hasParentBus() {
			ni.parentBus.messageStaticCANIDs.remove(msg.staticCANID)
		}

	} else {
		ni.sentMessageIDs.remove(msg.id)
	}

	return nil
}

// RemoveAllSentMessages removes all the messages sent by the [NodeInterface].
func (ni *NodeInterface) RemoveAllSentMessages() {
	for _, tmpMsg := range ni.sentMessages.entries() {
		tmpMsg.senderNodeInt = nil

		if ni.hasParentBus() && tmpMsg.hasStaticCANID {
			ni.parentBus.messageStaticCANIDs.remove(tmpMsg.staticCANID)
		}
	}

	ni.sentMessages.clear()
	ni.sentMessageNames.clear()
	ni.sentMessageIDs.clear()
	ni.sentMessageStaticCANIDs.clear()
}

// GetSentMessageByName returns the sent [Message] with the given name.
//
// It returns an [ErrNotFound] wrapped by a [NameError]
// if the name does not match any message.
func (ni *NodeInterface) GetSentMessageByName(name string) (*Message, error) {
	id, err := ni.sentMessageNames.getValue(name)
	if err != nil {
		return nil, ni.errorf(&NameError{Name: name, Err: err})
	}

	msg, err := ni.sentMessages.getValue(id)
	if err != nil {
		panic(err)
	}

	return msg, nil
}

// SentMessages returns a slice of messages sent by the [NodeInterface].
func (ni *NodeInterface) SentMessages() []*Message {
	msgSlice := ni.sentMessages.getValues()
	slices.SortFunc(msgSlice, func(a, b *Message) int {
		return int(a.id) - int(b.id)
	})
	return msgSlice
}

// AddReceivedMessage adds a [Message] that the [NodeInterface] can receive.
//
// It returns an [ArgError] if the given message is nil or
// a [ErrReceiverIsSender] wrapped by a [AddEntityError]
// if the message is already sent by the [NodeInterface].
func (ni *NodeInterface) AddReceivedMessage(message *Message) error {
	if message == nil {
		return ni.errorf(&ArgError{
			Name: "message",
			Err:  ErrIsNil,
		})
	}

	if err := ni.addReceivedMessage(message); err != nil {
		return ni.errorf(&AddEntityError{
			EntityID: message.entityID,
			Name:     message.name,
			Err:      err,
		})
	}

	return nil
}

// RemoveReceivedMessage removes a [Message] received by the [NodeInterface].
//
// It returns an [ErrNotFound] wrapped by a [RemoveEntityError]
// if the given entity id does not match any received message.
func (ni *NodeInterface) RemoveReceivedMessage(messageEntityID EntityID) error {
	msg, err := ni.receivedMessages.getValue(messageEntityID)
	if err != nil {
		return ni.errorf(&RemoveEntityError{
			EntityID: messageEntityID,
			Err:      err,
		})
	}

	ni.removeReceivedMessage(msg)

	return nil
}

// RemoveAllReceivedMessages removes all the messages received by the [NodeInterface].
func (ni *NodeInterface) RemoveAllReceivedMessages() {
	for _, tmpMsg := range ni.receivedMessages.entries() {
		tmpMsg.receivers.remove(ni.node.entityID)
	}

	ni.receivedMessages.clear()
}

// ReceivedMessages returns a slice of messages received by the [NodeInterface].
func (ni *NodeInterface) ReceivedMessages() []*Message {
	msgSlice := ni.receivedMessages.getValues()
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

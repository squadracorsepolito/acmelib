package acmelib

import (
	"cmp"
	"slices"

	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// NodeInterface represents an interface between a [Bus] and a [Node].
type NodeInterface struct {
	parentBus *Bus

	sentMessages            *collection.Map[EntityID, *Message]
	sentMessageNames        *collection.Map[string, EntityID]
	sentMessageIDs          *collection.Map[MessageID, EntityID]
	sentMessageStaticCANIDs *collection.Map[CANID, EntityID]

	receivedMessages *collection.Map[EntityID, *Message]

	number int
	node   *Node
}

func newNodeInterface(number int, node *Node) *NodeInterface {
	return &NodeInterface{
		parentBus: nil,

		sentMessages:            collection.NewMap[EntityID, *Message](),
		sentMessageNames:        collection.NewMap[string, EntityID](),
		sentMessageIDs:          collection.NewMap[MessageID, EntityID](),
		sentMessageStaticCANIDs: collection.NewMap[CANID, EntityID](),

		receivedMessages: collection.NewMap[EntityID, *Message](),

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

func (ni *NodeInterface) stringify(s *stringer.Stringer) {
	s.Write("number: %d\n", ni.number)

	s.Write("node:\n")
	s.Indent()
	ni.node.stringify(s)
	s.Unindent()

	if ni.sentMessages.Size() > 0 {
		s.Write("sent_messages:\n")
		s.Indent()
		for _, msg := range ni.SentMessages() {
			msg.stringify(s)
		}
		s.Unindent()
	}

	if ni.receivedMessages.Size() > 0 {
		s.Write("received_messages:\n")
		s.Indent()
		for _, msg := range ni.ReceivedMessages() {
			s.Write("entity_id: %s; name: %s\n", msg.EntityID(), msg.Name())
		}
		s.Unindent()
	}
}

func (ni *NodeInterface) String() string {
	s := stringer.New()
	s.Write("node_interface:\n")
	ni.stringify(s)
	return s.String()
}

func (ni *NodeInterface) verifyMessageName(name string) error {
	if ni.sentMessageNames.Has(name) {
		return newNameError(name, ErrIsDuplicated)
	}
	return nil
}

func (ni *NodeInterface) verifyMessageID(msgID MessageID) error {
	if ni.sentMessageIDs.Has(msgID) {
		return newMessageIDError(msgID, ErrIsDuplicated)
	}
	return nil
}

func (ni *NodeInterface) verifyStaticCANID(staticCANID CANID) error {
	if ni.sentMessageStaticCANIDs.Has(staticCANID) {
		return newCANIDError(staticCANID, ErrIsDuplicated)
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
	if ni.sentMessages.Has(msg.entityID) {
		return ErrReceiverIsSender
	}

	ni.receivedMessages.Set(msg.entityID, msg)
	msg.receivers.Set(ni.node.entityID, ni)

	return nil
}

func (ni *NodeInterface) removeReceivedMessage(msg *Message) {
	ni.receivedMessages.Delete(msg.entityID)
	msg.receivers.Delete(ni.node.entityID)
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
		return newArgError("message", ErrIsNil)
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
			ni.parentBus.messageStaticCANIDs.Set(message.staticCANID, message.entityID)
		}

		ni.sentMessageStaticCANIDs.Set(message.staticCANID, message.entityID)

	} else {
		if err := ni.verifyMessageID(message.id); err != nil {
			addMsgErr.Err = err
			return ni.errorf(addMsgErr)
		}
		ni.sentMessageIDs.Set(message.id, message.entityID)
	}

	ni.sentMessages.Set(message.entityID, message)
	ni.sentMessageNames.Set(message.name, message.entityID)

	message.senderNodeInt = ni

	return nil
}

// RemoveSentMessage removes a [Message] sent by the [NodeInterface].
//
// It returns an [ErrNotFound] if the given entity id does not match
// any message.
func (ni *NodeInterface) RemoveSentMessage(messageEntityID EntityID) error {
	msg, ok := ni.sentMessages.Get(messageEntityID)
	if !ok {
		return ni.errorf(ErrNotFound)
	}

	msg.senderNodeInt = nil

	ni.sentMessages.Delete(messageEntityID)
	ni.sentMessageNames.Delete(msg.name)

	if msg.hasStaticCANID {
		ni.sentMessageStaticCANIDs.Delete(msg.staticCANID)

		if ni.hasParentBus() {
			ni.parentBus.messageStaticCANIDs.Delete(msg.staticCANID)
		}

	} else {
		ni.sentMessageIDs.Delete(msg.id)
	}

	return nil
}

// RemoveAllSentMessages removes all the messages sent by the [NodeInterface].
func (ni *NodeInterface) RemoveAllSentMessages() {
	for tmpMsg := range ni.sentMessages.Values() {
		tmpMsg.senderNodeInt = nil

		if ni.hasParentBus() && tmpMsg.hasStaticCANID {
			ni.parentBus.messageStaticCANIDs.Delete(tmpMsg.staticCANID)
		}
	}

	ni.sentMessages.Clear()
	ni.sentMessageNames.Clear()
	ni.sentMessageIDs.Clear()
	ni.sentMessageStaticCANIDs.Clear()
}

// GetSentMessageByName returns the sent [Message] with the given name.
//
// It returns [ErrNotFound] if the name does not match any message.
func (ni *NodeInterface) GetSentMessageByName(name string) (*Message, error) {
	id, ok := ni.sentMessageNames.Get(name)
	if !ok {
		return nil, ni.errorf(ErrNotFound)
	}

	msg, ok := ni.sentMessages.Get(id)
	if !ok {
		return nil, ni.errorf(ErrNotFound)
	}

	return msg, nil
}

// SentMessages returns a slice of messages sent by the [NodeInterface].
func (ni *NodeInterface) SentMessages() []*Message {
	msgSlice := slices.Collect(ni.sentMessages.Values())
	slices.SortFunc(msgSlice, func(a, b *Message) int {
		return cmp.Compare(a.id, b.id)
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
	msg, ok := ni.receivedMessages.Get(messageEntityID)
	if !ok {
		return ni.errorf(ErrNotFound)
	}

	ni.removeReceivedMessage(msg)

	return nil
}

// RemoveAllReceivedMessages removes all the messages received by the [NodeInterface].
func (ni *NodeInterface) RemoveAllReceivedMessages() {
	for tmpMsg := range ni.receivedMessages.Values() {
		tmpMsg.receivers.Delete(ni.node.entityID)
	}

	ni.receivedMessages.Clear()
}

// ReceivedMessages returns a slice of messages received by the [NodeInterface].
func (ni *NodeInterface) ReceivedMessages() []*Message {
	msgSlice := slices.Collect(ni.receivedMessages.Values())
	slices.SortFunc(msgSlice, func(a, b *Message) int {
		return cmp.Compare(a.id, b.id)
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

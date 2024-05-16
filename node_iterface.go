package acmelib

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

type NodeInterface struct {
	entityID   EntityID
	createTime time.Time
	number     int
	node       *Node

	parentBus *Bus

	messages     *set[EntityID, *Message]
	messageNames *set[string, EntityID]
	messageIDs   *set[MessageID, EntityID]
}

func newNodeInterface(number int, node *Node) *NodeInterface {
	return &NodeInterface{
		entityID:   newEntityID(),
		createTime: time.Now(),
		number:     number,

		node:      node,
		parentBus: nil,

		messages:     newSet[EntityID, *Message](),
		messageNames: newSet[string, EntityID](),
		messageIDs:   newSet[MessageID, EntityID](),
	}
}

func (ni *NodeInterface) EntityID() EntityID {
	return ni.entityID
}

func (ni *NodeInterface) CreateTime() time.Time {
	return ni.createTime
}

func (ni *NodeInterface) Number() int {
	return ni.number
}

func (ni *NodeInterface) hasParentBus() bool {
	return ni.parentBus != nil
}

func (ni *NodeInterface) errorf(err error) error {
	nodeIntErr := &EntityError{
		Kind:     EntityKindNode,
		EntityID: ni.entityID,
		Name:     ni.GetName(),
		Err:      err,
	}

	if ni.hasParentBus() {
		return ni.parentBus.errorf(nodeIntErr)
	}

	return nodeIntErr
}

func (ni *NodeInterface) stringify(b *strings.Builder, tabs int) {
	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%sentity_id: %s\n", tabStr, ni.entityID))
	b.WriteString(fmt.Sprintf("%sname: %s\n", tabStr, ni.GetName()))
	b.WriteString(fmt.Sprintf("%screate_time: %s\n", tabStr, ni.createTime.Format(time.RFC3339)))
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

func (ni *NodeInterface) RemoveAllMessages() {
	for _, tmpMsg := range ni.messages.entries() {
		tmpMsg.senderNodeInt = nil
	}

	ni.messages.clear()
	ni.messageNames.clear()
	ni.messageIDs.clear()
}

func (ni *NodeInterface) Messages() []*Message {
	msgSlice := ni.messages.getValues()
	slices.SortFunc(msgSlice, func(a, b *Message) int {
		return int(a.id) - int(b.id)
	})
	return msgSlice
}

func (ni *NodeInterface) GetName() string {
	return fmt.Sprintf("%s/int%d", ni.node.name, ni.number)
}

func (ni *NodeInterface) Node() *Node {
	return ni.node
}

func (ni *NodeInterface) ParentBus() *Bus {
	return ni.parentBus
}

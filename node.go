package acmelib

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type NodeID uint32

type Node struct {
	*entityWithAttributes

	parentBuses *set[EntityID, *Bus]
	parErrID    EntityID

	messages     *set[EntityID, *Message]
	messageNames *set[string, EntityID]
	messageIDs   *set[MessageID, EntityID]

	id NodeID
}

func NewNode(name, desc string, id NodeID) *Node {
	return &Node{
		entityWithAttributes: newEntityWithAttributes(name, desc, AttributeReferenceKindNode),

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

func (n *Node) ParentBuses() []*Bus {
	return n.parentBuses.getValues()
}

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

func (n *Node) RemoveAllMessages() {
	for _, tmpMsg := range n.messages.entries() {
		tmpMsg.resetID()
		tmpMsg.parentNodes.remove(n.entityID)
	}

	n.messages.clear()
	n.messageNames.clear()
	n.messageIDs.clear()
}

func (n *Node) Messages() []*Message {
	msgSlice := n.messages.getValues()
	slices.SortFunc(msgSlice, func(a, b *Message) int { return int(a.id) - int(b.id) })
	return msgSlice
}

func (n *Node) ID() NodeID {
	return n.id
}

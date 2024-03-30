package acmelib

import "fmt"

type NodeID int

type Node struct {
	*entity

	parentBus *Bus

	messages     map[EntityID]*Message
	messageNames map[string]EntityID
	messageIDs   map[MessageID]EntityID

	id NodeID
}

func NewNode(name, desc string) *Node {
	return &Node{
		entity: newEntity(name, desc),

		parentBus: nil,

		messages:     make(map[EntityID]*Message),
		messageNames: make(map[string]EntityID),
		messageIDs:   make(map[MessageID]EntityID),

		id: 0,
	}
}

func (n *Node) hasParent() bool {
	return n.parentBus != nil
}

func (n *Node) getMessageByEntID(msgEntID EntityID) (*Message, error) {
	if msg, ok := n.messages[msgEntID]; ok {
		return msg, nil
	}
	return nil, fmt.Errorf("message not found")
}

func (n *Node) addMessageName(msgEntID EntityID, name string) {
	n.messageNames[name] = msgEntID
}

func (n *Node) removeMessageName(name string) {
	delete(n.messageNames, name)
}

func (n *Node) verifyMessageName(name string) error {
	if _, ok := n.messageNames[name]; ok {
		return fmt.Errorf(`message name "%s" is duplicated`, name)
	}
	return nil
}

func (n *Node) modifyMessageName(msgEntID EntityID, newName string) error {
	msg, err := n.getMessageByEntID(msgEntID)
	if err != nil {
		return err
	}

	oldName := msg.Name()

	n.removeMessageName(oldName)
	n.addMessageName(msgEntID, newName)

	return nil
}

func (n *Node) verifyMessageID(msgID MessageID) error {
	if _, ok := n.messageIDs[msgID]; ok {
		return fmt.Errorf(`message id "%d" is duplicated`, msgID)
	}
	return nil
}

func (n *Node) errorf(err error) error {
	nodeErr := fmt.Errorf(`node "%s": %w`, n.name, err)
	if n.hasParent() {
		return n.parentBus.errorf(nodeErr)
	}
	return nodeErr
}

func (n *Node) setParent(bus *Bus) {
	n.parentBus = bus
}

func (n *Node) AddMessage(message *Message) error {
	if err := n.verifyMessageName(message.Name()); err != nil {
		return n.errorf(fmt.Errorf(`cannot add message "%s" : %w`, message.Name(), err))
	}

	message.generateID(len(n.messages)+1, n.id)
	if err := n.verifyMessageID(message.ID()); err != nil {
		return n.errorf(fmt.Errorf(`cannot add message "%s" : %w`, message.Name(), err))
	}

	entID := message.EntityID()
	n.messages[entID] = message
	n.addMessageName(entID, message.Name())
	n.messageIDs[message.ID()] = entID

	message.setParent(n)

	return nil
}

func (n *Node) UpdateName(newName string) error {
	if n.name == newName {
		return nil
	}

	if n.hasParent() {
		if err := n.parentBus.verifyNodeName(newName); err != nil {
			return n.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}

		if err := n.parentBus.modifyNodeName(n.entityID, newName); err != nil {
			return n.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}
	}

	n.name = newName

	return nil
}

func (n *Node) SetID(nodeID NodeID) {
	n.id = nodeID
}

func (n *Node) ID() NodeID {
	return n.id
}

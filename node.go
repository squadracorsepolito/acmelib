package acmelib

import "fmt"

type Node struct {
	*entity
	ParentBus *Bus

	messages *entityCollection[*Message]
}

func NewNode(name, desc string) *Node {
	return &Node{
		entity: newEntity(name, desc),

		messages: newEntityCollection[*Message](),
	}
}

func (n *Node) errorf(err error) error {
	nodeErr := fmt.Errorf(`node "%s": %v`, n.name, err)
	if n.ParentBus != nil {
		return n.ParentBus.errorf(nodeErr)
	}
	return nodeErr
}

func (n *Node) UpdateName(name string) error {
	if n.ParentBus != nil {
		if err := n.ParentBus.nodes.updateEntityName(n.entityID, n.name, name); err != nil {
			return n.errorf(err)
		}
	}

	return n.entity.UpdateName(name)
}

// func (n *Node) AddMessage(message *Message) error {
// 	if err := n.messages.addEntity(message); err != nil {
// 		return n.errorf(err)
// 	}

// 	message.ParentNode = n
// 	n.setUpdateTimeNow()

// 	return nil
// }

func (n *Node) ListMessages() []*Message {
	return n.messages.listEntities()
}

func (n *Node) RemoveMessage(messageID EntityID) error {
	if err := n.messages.removeEntity(messageID); err != nil {
		return n.errorf(err)
	}

	return nil
}

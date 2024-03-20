package acmelib

import "fmt"

type NodeSortingMethod string

const (
	NodesByName       NodeSortingMethod = "nodes_by_name"
	NodesByCreateTime NodeSortingMethod = "nodes_by_create_time"
	NodesByUpdateTime NodeSortingMethod = "nodes_by_update_time"
)

func newNodeSorter() *entitySorter[NodeSortingMethod, *Node] {
	return newEntitySorter(
		newEntitySorterMethod(NodesByName, func(nodes []*Node) []*Node { return sortByName(nodes) }),
		newEntitySorterMethod(NodesByCreateTime, func(nodes []*Node) []*Node { return sortByCreateTime(nodes) }),
		newEntitySorterMethod(NodesByUpdateTime, func(nodes []*Node) []*Node { return sortByUpdateTime(nodes) }),
	)
}

type Node struct {
	*entity
	ParentNode *Bus

	messages *entityCollection[*Message, MessageSortingMethod]
}

func NewNode(name, desc string) *Node {
	return &Node{
		entity: newEntity(name, desc),

		messages: newEntityCollection(newMessageSorter()),
	}
}

func (n *Node) errorf(err error) error {
	return n.ParentNode.errorf(fmt.Errorf("node %s: %v", n.Name, err))
}

func (n *Node) UpdateName(name string) error {
	if err := n.ParentNode.nodes.updateName(n.ID, n.Name, name); err != nil {
		return n.errorf(err)
	}

	return n.entity.UpdateName(name)
}

func (n *Node) AddMessage(message *Message) error {
	if err := n.messages.addEntity(message); err != nil {
		return n.errorf(err)
	}

	message.ParentMessage = n
	n.setUpdateTimeNow()

	return nil
}

func (n *Node) ListMessages(sortingMathod MessageSortingMethod) []*Message {
	return n.messages.listEntities(sortingMathod)
}

func (n *Node) RemoveMessage(messageID EntityID) error {
	if err := n.messages.removeEntity(messageID); err != nil {
		return n.errorf(err)
	}

	n.setUpdateTimeNow()

	return nil
}

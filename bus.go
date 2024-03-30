package acmelib

import "fmt"

type Bus struct {
	*entity

	parentProject *Project

	nodes     map[EntityID]*Node
	nodeNames map[string]EntityID
	nodeIDs   map[NodeID]EntityID
}

func NewBus(name, desc string) *Bus {
	return &Bus{
		entity: newEntity(name, desc),

		parentProject: nil,

		nodes:     make(map[EntityID]*Node),
		nodeNames: make(map[string]EntityID),
		nodeIDs:   make(map[NodeID]EntityID),
	}
}

func (b *Bus) hasParent() bool {
	return b.parentProject != nil
}

func (b *Bus) errorf(err error) error {
	busErr := fmt.Errorf(`bus "%s": %w`, b.name, err)
	if b.hasParent() {
		return b.parentProject.errorf(busErr)
	}
	return busErr
}

func (b *Bus) getNodeByEntID(nodeEndID EntityID) (*Node, error) {
	if node, ok := b.nodes[nodeEndID]; ok {
		return node, nil
	}
	return nil, fmt.Errorf("node not found")
}

func (b *Bus) addNodeName(nodeEndID EntityID, name string) {
	b.nodeNames[name] = nodeEndID
}

func (b *Bus) removeNodeName(name string) {
	delete(b.nodeNames, name)
}

func (b *Bus) verifyNodeName(name string) error {
	if _, ok := b.nodeNames[name]; ok {
		return fmt.Errorf(`node name "%s" is duplicated`, name)
	}
	return nil
}

func (b *Bus) modifyNodeName(nodeEndID EntityID, newName string) error {
	node, err := b.getNodeByEntID(nodeEndID)
	if err != nil {
		return err
	}

	oldName := node.Name()

	b.removeNodeName(oldName)
	b.addNodeName(nodeEndID, newName)

	return nil
}

func (b *Bus) verifyNodeID(nodeID NodeID) error {
	if _, ok := b.nodeIDs[nodeID]; ok {
		return fmt.Errorf(`node id "%d" is duplicated`, nodeID)
	}
	return nil
}

func (b *Bus) UpdateName(newName string) error {
	if b.name == newName {
		return nil
	}

	b.name = newName

	return nil
}

func (b *Bus) AddNode(node *Node) error {
	if err := b.verifyNodeName(node.Name()); err != nil {
		return b.errorf(fmt.Errorf(`cannot add node "%s" : %w`, node.Name(), err))
	}

	node.SetID(NodeID(len(b.nodes) + 1))
	if err := b.verifyNodeID(node.ID()); err != nil {
		return b.errorf(fmt.Errorf(`cannot add node "%s" : %w`, node.Name(), err))
	}

	entID := node.EntityID()
	b.nodes[entID] = node
	b.addNodeName(entID, node.Name())
	b.nodeIDs[node.ID()] = entID

	node.setParent(b)

	return nil
}

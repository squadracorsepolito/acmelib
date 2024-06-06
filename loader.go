package acmelib

import acmelibv1 "github.com/squadracorsepolito/acmelib/proto/gen/go/acmelib/v1"

type loader struct {
	refNodes map[string]*Node
}

func newLoader() *loader {
	return &loader{
		refNodes: make(map[string]*Node),
	}
}

func (l *loader) loadNetwork(pNet *acmelibv1.Network) (*Network, error) {
	net := NewNetwork(pNet.Entity.Name)
	if len(pNet.Entity.Desc) > 0 {
		net.SetDesc(pNet.Entity.Desc)
	}
	net.entityID = EntityID(pNet.Entity.EntityId)
	net.createTime = pNet.Entity.CreateTime.AsTime()

	for _, pNode := range pNet.Nodes {
		node, err := l.loadNode(pNode)
		if err != nil {
			return nil, err
		}
		l.refNodes[node.entityID.String()] = node
	}

	for _, pBus := range pNet.Buses {
		bus, err := l.loadBus(pBus)
		if err != nil {
			return nil, err
		}
		net.AddBus(bus)
	}

	return net, nil
}

func (l *loader) loadNode(pNode *acmelibv1.Node) (*Node, error) {
	node := NewNode(pNode.Entity.Name, NodeID(pNode.NodeId), int(pNode.InterfaceCount))
	if len(pNode.Entity.Desc) > 0 {
		node.SetDesc(pNode.Entity.Desc)
	}
	return node, nil
}

func (l *loader) loadBus(pBus *acmelibv1.Bus) (*Bus, error) {
	bus := NewBus(pBus.Entity.Name)
	if len(pBus.Entity.Desc) > 0 {
		bus.SetDesc(pBus.Entity.Desc)
	}
	bus.entityID = EntityID(pBus.Entity.EntityId)
	bus.createTime = pBus.Entity.CreateTime.AsTime()

	return bus, nil
}

func (l *loader) loadNodeInterface(pNodeInt *acmelibv1.NodeInterface) (*NodeInterface, error) {
	return nil, nil
}

func (l *loader) loadMessage(pMsg *acmelibv1.Message) (*Message, error) {
	msg := NewMessage(pMsg.Entity.Name, MessageID(pMsg.MessageId), int(pMsg.SizeByte))
	if len(pMsg.Entity.Desc) > 0 {
		msg.SetDesc(pMsg.Entity.Desc)
	}
	msg.entityID = EntityID(pMsg.Entity.EntityId)
	msg.createTime = pMsg.Entity.CreateTime.AsTime()

	return msg, nil
}

func (l *loader) loadSignal(pSig *acmelibv1.Signal) (Signal, error) {
	var sig Signal

	switch pSig.Kind {
	case acmelibv1.SignalKind_SIGNAL_KIND_STANDARD:

	}

	return sig, nil
}

func (l *loader) loadStandardSignal(pSig *acmelibv1.Signal) (*StandardSignal, error) {
	pStdSig := pSig.GetStandard()

	return nil, nil
}

func (l *loader) loadSignalType(pSigType *acmelibv1.SignalType) (*SignalType, error) {
	switch pSigType.Kind {
	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_CUSTOM:
		sigType, err := NewCustomSignalType(pSigType.Entity.Name, int(pSigType.Size), pSigType.Signed, pSigType.Min, pSigType.Max, pSigType.Scale, pSigType.Offset)
		if err != nil {
			return nil, err
		}
		if len(pSigType.Entity.Desc) > 0 {
			sigType.SetDesc(pSigType.Entity.Desc)
		}

	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_FLAG:

	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_INTEGER:

	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_DECIMAL:
	}

	return nil, nil
}

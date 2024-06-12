package acmelib

import (
	"time"

	acmelibv1 "github.com/squadracorsepolito/acmelib/proto/gen/go/acmelib/v1"
)

type loader struct {
	refNodes map[string]*Node
}

func newLoader() *loader {
	return &loader{
		refNodes: make(map[string]*Node),
	}
}

func (l *loader) getEntity(pEnt *acmelibv1.Entity, entKind EntityKind) *entity {
	var cTime time.Time
	if pEnt.CreateTime.IsValid() {
		cTime = pEnt.CreateTime.AsTime()
	} else {
		cTime = time.Now()
	}

	return &entity{
		entityID:   EntityID(pEnt.EntityId),
		name:       pEnt.Name,
		desc:       pEnt.Desc,
		entityKind: entKind,
		createTime: cTime,
	}
}

func (l *loader) loadNetwork(pNet *acmelibv1.Network) (*Network, error) {
	net := newNetworkFromEntity(l.getEntity(pNet.Entity, EntityKindNetwork))

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
	node := newNodeFromEntity(l.getEntity(pNode.Entity, EntityKindNode), NodeID(pNode.NodeId), int(pNode.InterfaceCount))
	return node, nil
}

func (l *loader) loadBus(pBus *acmelibv1.Bus) (*Bus, error) {
	bus := newBusFromEntity(l.getEntity(pBus.Entity, EntityKindBus))

	for _, pNodeInt := range pBus.NodeInterfaces {
		nodeInt, err := l.loadNodeInterface(pNodeInt)
		if err != nil {
			return nil, err
		}

		if err := bus.AddNodeInterface(nodeInt); err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func (l *loader) loadNodeInterface(pNodeInt *acmelibv1.NodeInterface) (*NodeInterface, error) {
	node, ok := l.refNodes[pNodeInt.NodeEntityId]
	if !ok {
		return nil, ErrNotFound
	}

	nodeInt, err := node.GetInterface(int(pNodeInt.Number))
	if err != nil {
		return nil, err
	}

	for _, pMsg := range pNodeInt.Messages {
		msg, err := l.loadMessage(pMsg)
		if err != nil {
			return nil, err
		}

		if err := nodeInt.AddMessage(msg); err != nil {
			return nil, err
		}
	}

	return nodeInt, nil
}

func (l *loader) loadMessage(pMsg *acmelibv1.Message) (*Message, error) {
	msg := newMessageFromEntity(l.getEntity(pMsg.Entity, EntityKindMessage), MessageID(pMsg.MessageId), int(pMsg.SizeByte))

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
	// pStdSig := pSig.GetStandard()

	return nil, nil
}

func (l *loader) loadSignalType(pSigType *acmelibv1.SignalType) (*SignalType, error) {
	var kind SignalTypeKind
	switch pSigType.Kind {
	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_CUSTOM:
		kind = SignalTypeKindCustom

	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_FLAG:
		kind = SignalTypeKindFlag

	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_INTEGER:
		kind = SignalTypeKindInteger

	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_DECIMAL:
		kind = SignalTypeKindDecimal
	}

	sigType, err := newSignalType("", kind, int(pSigType.Size), pSigType.Signed, pSigType.Min, pSigType.Max, pSigType.Scale, pSigType.Offset)
	if err != nil {
		return nil, err
	}
	sigType.entity = l.getEntity(pSigType.Entity, EntityKindSignalType)

	return sigType, nil
}

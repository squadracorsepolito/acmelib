package acmelib

import (
	"time"

	acmelibv1 "github.com/squadracorsepolito/acmelib/proto/gen/go/acmelib/v1"
)

type loader struct {
	refCANIDBuilders map[string]*CANIDBuilder
	refNodes         map[string]*Node
	refSigTypes      map[string]*SignalType
	refSigUnits      map[string]*SignalUnit
	refSigEnums      map[string]*SignalEnum
	refAttributes    map[string]Attribute
}

func newLoader() *loader {
	return &loader{
		refCANIDBuilders: make(map[string]*CANIDBuilder),
		refNodes:         make(map[string]*Node),
		refSigTypes:      make(map[string]*SignalType),
		refSigUnits:      make(map[string]*SignalUnit),
		refSigEnums:      make(map[string]*SignalEnum),
		refAttributes:    make(map[string]Attribute),
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

	for _, pBuilder := range pNet.CanidBuilders {
		l.refCANIDBuilders[pBuilder.Entity.EntityId] = l.loadCANIDBuilder(pBuilder)
	}

	for _, pNode := range pNet.Nodes {
		node, err := l.loadNode(pNode)
		if err != nil {
			return nil, err
		}
		l.refNodes[pNode.Entity.EntityId] = node
	}

	for _, pSigType := range pNet.SignalTypes {
		sigType, err := l.loadSignalType(pSigType)
		if err != nil {
			return nil, err
		}
		l.refSigTypes[pSigType.Entity.EntityId] = sigType
	}

	for _, pSigUnit := range pNet.SignalUnits {
		l.refSigUnits[pSigUnit.Entity.EntityId] = l.loadSignalUnit(pSigUnit)
	}

	for _, pSigEnum := range pNet.SignalEnums {
		sigEnum, err := l.loadSignalEnum(pSigEnum)
		if err != nil {
			return nil, err
		}
		l.refSigEnums[pSigEnum.Entity.EntityId] = sigEnum
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

func (l *loader) loadCANIDBuilder(pBuilder *acmelibv1.CANIDBuilder) *CANIDBuilder {
	builder := newCANIDBuilderFromEntity(l.getEntity(pBuilder.Entity, EntityKindCANIDBuilder))

	for _, pBuilderOp := range pBuilder.Operations {
		builder.operations = append(builder.operations, l.loadCANIDBuilderOp(pBuilderOp))
	}

	return builder
}

func (l *loader) loadCANIDBuilderOp(pBuilderOp *acmelibv1.CANIDBuilderOp) *CANIDBuilderOp {
	var kind CANIDBuilderOpKind
	switch pBuilderOp.Kind {
	case acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_PRIORITY:
		kind = CANIDBuilderOpKindMessagePriority
	case acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_ID:
		kind = CANIDBuilderOpKindMessageID
	case acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_NODE_ID:
		kind = CANIDBuilderOpKindNodeID
	case acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_BIT_MASK:
		kind = CANIDBuilderOpKindBitMask
	}
	return newCANIDBuilderOp(kind, int(pBuilderOp.From), int(pBuilderOp.Len))
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

	ent := l.getEntity(pSigType.Entity, EntityKindSignalType)
	return newSignalTypeFromEntity(ent, kind, int(pSigType.Size), pSigType.Signed, pSigType.Min, pSigType.Max, pSigType.Scale, pSigType.Offset)
}

func (l *loader) loadSignalUnit(pSigUnit *acmelibv1.SignalUnit) *SignalUnit {
	var kind SignalUnitKind
	switch pSigUnit.Kind {
	case acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_CUSTOM:
		kind = SignalUnitKindCustom
	case acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_TEMPERATURE:
		kind = SignalUnitKindTemperature
	case acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_ELECTRICAL:
		kind = SignalUnitKindElectrical
	case acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_POWER:
		kind = SignalUnitKindPower
	}
	return newSignalUnitFromEntity(l.getEntity(pSigUnit.Entity, EntityKindSignalUnit), kind, pSigUnit.Symbol)
}

func (l *loader) loadSignalEnum(pSigEnum *acmelibv1.SignalEnum) (*SignalEnum, error) {
	enum := newSignalEnumFromEntity(l.getEntity(pSigEnum.Entity, EntityKindSignalEnum))

	for _, pVal := range pSigEnum.Values {
		val := l.loadSignalEnumValue(pVal)
		if err := enum.AddValue(val); err != nil {
			return nil, err
		}
	}

	if pSigEnum.MinSize != 0 {
		enum.minSize = int(pSigEnum.MinSize)
	}

	return enum, nil
}

func (l *loader) loadSignalEnumValue(pVal *acmelibv1.SignalEnumValue) *SignalEnumValue {
	return newSignalEnumValueFromEntity(l.getEntity(pVal.Entity, EntityKindSignalEnumValue), int(pVal.Index))
}

func (l *loader) loadAttribute(pAtt *acmelibv1.Attribute) (Attribute, error) {
	var att Attribute

	return att, nil
}

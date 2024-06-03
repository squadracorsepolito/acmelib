package acmelib

import (
	acmelibv1 "github.com/squadracorsepolito/acmelib/proto/gen/go/acmelib/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type saver struct {
	refCANIDBuilders map[EntityID]*CANIDBuilder
	refNodes         map[EntityID]*Node
	refSigTypes      map[EntityID]*SignalType
	refSigUnits      map[EntityID]*SignalUnit
	refSigEnums      map[EntityID]*SignalEnum
	refAttributes    map[EntityID]Attribute
}

func newSaver() *saver {
	return &saver{
		refCANIDBuilders: make(map[EntityID]*CANIDBuilder),
		refNodes:         make(map[EntityID]*Node),
		refSigTypes:      make(map[EntityID]*SignalType),
		refSigUnits:      make(map[EntityID]*SignalUnit),
		refSigEnums:      make(map[EntityID]*SignalEnum),
		refAttributes:    make(map[EntityID]Attribute),
	}
}

func (s *saver) getEntityKind(ek EntityKind) acmelibv1.EntityKind {
	switch ek {
	case EntityKindNetwork:
		return acmelibv1.EntityKind_ENTITY_KIND_NETWORK
	case EntityKindBus:
		return acmelibv1.EntityKind_ENTITY_KIND_BUS
	case EntityKindNode:
		return acmelibv1.EntityKind_ENTITY_KIND_NODE
	case EntityKindNodeInterface:
		return acmelibv1.EntityKind_ENTITY_KIND_NODE_INTERFACE
	case EntityKindMessage:
		return acmelibv1.EntityKind_ENTITY_KIND_MESSAGE
	case EntityKindSignal:
		return acmelibv1.EntityKind_ENTITY_KIND_SIGNAL
	case EntityKindSignalType:
		return acmelibv1.EntityKind_ENTITY_KIND_SIGNAL_TYPE
	case EntityKindSignalUnit:
		return acmelibv1.EntityKind_ENTITY_KIND_SIGNAL_UNIT
	case EntityKindSignalEnum:
		return acmelibv1.EntityKind_ENTITY_KIND_SIGNAL_ENUM
	case EntityKindSignalEnumValue:
		return acmelibv1.EntityKind_ENTITY_KIND_SIGNAL_ENUM_VALUE
	case EntityKindAttribute:
		return acmelibv1.EntityKind_ENTITY_KIND_ATTRIBUTE
	case EntityKindCANIDBuilder:
		return acmelibv1.EntityKind_ENTITY_KIND_CANID_BUILDER
	default:
		return acmelibv1.EntityKind_ENTITY_KIND_UNSPECIFIED
	}
}

func (s *saver) saveEntity(e *entity) *acmelibv1.Entity {
	pEnt := new(acmelibv1.Entity)

	pEnt.EntityId = e.entityID.String()
	pEnt.Desc = e.desc
	pEnt.Name = e.name
	pEnt.CreateTime = timestamppb.New(e.createTime)
	pEnt.EntityKind = s.getEntityKind(e.entityKind)

	return pEnt
}

func (s *saver) saveNetwork(net *Network) *acmelibv1.Network {
	pNet := new(acmelibv1.Network)

	pNet.Entity = s.saveEntity(net.entity)

	for _, bus := range net.Buses() {
		pNet.Buses = append(pNet.Buses, s.saveBus(bus))
	}

	for _, canIDBuilder := range s.refCANIDBuilders {
		pNet.CanidBuilders = append(pNet.CanidBuilders, s.saveCANIDBuilder(canIDBuilder))
	}

	for _, node := range s.refNodes {
		pNet.Nodes = append(pNet.Nodes, s.saveNode(node))
	}

	for _, sigType := range s.refSigTypes {
		pNet.SignalTypes = append(pNet.SignalTypes, s.saveSignalType(sigType))
	}

	for _, sigUnit := range s.refSigUnits {
		pNet.SignalUnits = append(pNet.SignalUnits, s.saveSignalUnit(sigUnit))
	}

	return pNet
}

func (s *saver) saveAttributeAssignments(attAss []*AttributeAssignment) []*acmelibv1.AttributeAssignment {
	pAttAss := []*acmelibv1.AttributeAssignment{}

	for _, tmpAttAss := range attAss {
		tmpAtt := tmpAttAss.attribute

		pTmpAttAss := new(acmelibv1.AttributeAssignment)

		pTmpAttAss.EntityId = tmpAttAss.EntityID().String()
		pTmpAttAss.AttributeEntityId = tmpAtt.EntityID().String()
		pTmpAttAss.EntityKind = s.getEntityKind(tmpAttAss.entity.EntityKind())

		switch tmpAtt.Type() {
		case AttributeTypeString, AttributeTypeEnum:
			pTmpAttAss.Value = &acmelibv1.AttributeAssignment_ValueString{
				ValueString: tmpAttAss.value.(string),
			}
		case AttributeTypeInteger:
			pTmpAttAss.Value = &acmelibv1.AttributeAssignment_ValueInt{
				ValueInt: tmpAttAss.value.(int32),
			}
		case AttributeTypeFloat:
			pTmpAttAss.Value = &acmelibv1.AttributeAssignment_ValueDouble{
				ValueDouble: tmpAttAss.value.(float64),
			}
		}

		pAttAss = append(pAttAss, pTmpAttAss)

		s.refAttributes[tmpAtt.EntityID()] = tmpAtt
	}

	return pAttAss
}

func (s *saver) saveBus(bus *Bus) *acmelibv1.Bus {
	pBus := new(acmelibv1.Bus)

	pBus.Entity = s.saveEntity(bus.entity)
	pBus.AttributeAssignments = s.saveAttributeAssignments(bus.AttributeAssignments())

	pBus.Baudrate = uint32(bus.baudrate)

	pBusType := acmelibv1.BusType_BUS_TYPE_UNSPECIFIED
	switch bus.typ {
	case BusTypeCAN2A:
		pBusType = acmelibv1.BusType_BUS_TYPE_CAN_2A
	}
	pBus.Type = pBusType

	for _, nodeInt := range bus.NodeInterfaces() {
		pBus.NodeInterfaces = append(pBus.NodeInterfaces, s.saveNodeInterface(nodeInt))
	}

	if bus.isDefCANIDBuilder {
		return pBus
	}

	if bus.canIDBuilder.ReferenceCount() > 0 {
		entID := bus.canIDBuilder.entityID
		pBus.CanidBuilder = &acmelibv1.Bus_CanidBuilderEntityId{
			CanidBuilderEntityId: entID.String(),
		}
		s.refCANIDBuilders[entID] = bus.canIDBuilder

		return pBus
	}

	pBus.CanidBuilder = &acmelibv1.Bus_EmbeddedCanidBuilder{
		EmbeddedCanidBuilder: s.saveCANIDBuilder(bus.canIDBuilder),
	}

	return pBus
}

func (s *saver) saveCANIDBuilder(builder *CANIDBuilder) *acmelibv1.CANIDBuilder {
	pBuilder := new(acmelibv1.CANIDBuilder)

	pBuilder.Entity = s.saveEntity(builder.entity)

	for _, tmpOp := range builder.operations {
		pBuilder.Operations = append(pBuilder.Operations, s.saveCANIDBuilderOp(tmpOp))
	}

	return pBuilder
}

func (s *saver) saveCANIDBuilderOp(builderOp *CANIDBuilderOp) *acmelibv1.CANIDBuilderOp {
	pBuilderOp := new(acmelibv1.CANIDBuilderOp)

	pOpKind := acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_UNSPECIFIED
	switch builderOp.kind {
	case CANIDBuilderOpKindMessagePriority:
		pOpKind = acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_PRIORITY
	case CANIDBuilderOpKindMessageID:
		pOpKind = acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_ID
	case CANIDBuilderOpKindNodeID:
		pOpKind = acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_NODE_ID
	case CANIDBuilderOpKindBitMask:
		pOpKind = acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_BIT_MASK
	}
	pBuilderOp.Kind = pOpKind

	pBuilderOp.From = uint32(builderOp.from)
	pBuilderOp.Len = uint32(builderOp.len)

	return pBuilderOp
}

func (s *saver) saveNode(node *Node) *acmelibv1.Node {
	pNode := new(acmelibv1.Node)

	pNode.Entity = s.saveEntity(node.entity)
	pNode.AttributeAssignments = s.saveAttributeAssignments(node.AttributeAssignments())

	pNode.NodeId = uint32(node.id)
	pNode.InterfaceCount = uint32(node.interfaceCount)

	return pNode
}

func (s *saver) saveNodeInterface(nodeInt *NodeInterface) *acmelibv1.NodeInterface {
	pNodeint := new(acmelibv1.NodeInterface)

	pNodeint.Entity = s.saveEntity(nodeInt.entity)
	pNodeint.Number = int32(nodeInt.number)

	for _, msg := range nodeInt.Messages() {
		pNodeint.Messages = append(pNodeint.Messages, s.saveMessage(msg))
	}

	if nodeInt.node.interfaceCount > 1 {
		entID := nodeInt.node.entityID
		pNodeint.Node = &acmelibv1.NodeInterface_NodeEntityId{
			NodeEntityId: entID.String(),
		}
		s.refNodes[entID] = nodeInt.node

		return pNodeint
	}

	pNodeint.Node = &acmelibv1.NodeInterface_EmbeddedNode{
		EmbeddedNode: s.saveNode(nodeInt.node),
	}

	return pNodeint
}

func (s *saver) saveMessage(msg *Message) *acmelibv1.Message {
	pMsg := new(acmelibv1.Message)

	pMsg.Entity = s.saveEntity(msg.entity)
	pMsg.AttributeAssignments = s.saveAttributeAssignments(msg.AttributeAssignments())

	pMsg.SizeByte = uint32(msg.sizeByte)
	pMsg.MessageId = uint32(msg.id)

	for _, sig := range msg.Signals() {
		pMsg.Signals = append(pMsg.Signals, s.saveSignal(sig))
	}

	pMsg.Payload = s.saveSignalPayload(msg.signalPayload)

	pMsg.StaticCanId = uint32(msg.staticCANID)
	pMsg.HasStaticCanId = msg.hasStaticCANID

	pPriority := acmelibv1.MessagePriority_MESSAGE_PRIORITY_UNSPECIFIED
	switch msg.priority {
	case MessagePriorityVeryHigh:
		pPriority = acmelibv1.MessagePriority_MESSAGE_PRIORITY_VERY_HIGH
	case MessagePriorityHigh:
		pPriority = acmelibv1.MessagePriority_MESSAGE_PRIORITY_HIGH
	case MessagePriorityMedium:
		pPriority = acmelibv1.MessagePriority_MESSAGE_PRIORITY_MEDIUM
	case MessagePriorityLow:
		pPriority = acmelibv1.MessagePriority_MESSAGE_PRIORITY_LOW
	}
	pMsg.Priority = pPriority

	pByteOrder := acmelibv1.MessageByteOrder_MESSAGE_BYTE_ORDER_UNSPECIFIED
	switch msg.byteOrder {
	case MessageByteOrderLittleEndian:
		pByteOrder = acmelibv1.MessageByteOrder_MESSAGE_BYTE_ORDER_LITTLE_ENDIAN
	case MessageByteOrderBigEndian:
		pByteOrder = acmelibv1.MessageByteOrder_MESSAGE_BYTE_ORDER_BIG_ENDIAN
	}
	pMsg.ByteOrder = pByteOrder

	pMsg.CycleTime = uint32(msg.cycleTime)

	pSendType := acmelibv1.MessageSendType_MESSAGE_SEND_TYPE_UNSPECIFIED
	switch msg.sendType {
	case MessageSendTypeCyclic:
		pSendType = acmelibv1.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC
	case MessageSendTypeCyclicIfActive:
		pSendType = acmelibv1.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE
	case MessageSendTypeCyclicAndTriggered:
		pSendType = acmelibv1.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_AND_TRIGGERED
	case MessageSendTypeCyclicIfActiveAndTriggered:
		pSendType = acmelibv1.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE_AND_TRIGGERED
	}
	pMsg.SendType = pSendType

	pMsg.DelayTime = uint32(msg.delayTime)
	pMsg.StartDelayTime = uint32(msg.startDelayTime)

	for _, rec := range msg.Receivers() {
		pMsg.ReceiverIds = append(pMsg.ReceiverIds, rec.entityID.String())
	}

	return pMsg
}

func (s *saver) saveSignalPayload(payload *signalPayload) *acmelibv1.SignalPayload {
	pPayload := new(acmelibv1.SignalPayload)

	for _, sig := range payload.signals {
		pPayload.Refs = append(pPayload.Refs, &acmelibv1.SignalPayloadRef{
			SignalEntityId: sig.EntityID().String(),
			RelStartBit:    uint32(sig.getRelStartBit()),
		})
	}

	return pPayload
}

func (s *saver) saveSignal(sig Signal) *acmelibv1.Signal {
	pSig := new(acmelibv1.Signal)

	pSig.AttributeAssignments = s.saveAttributeAssignments(sig.AttributeAssignments())

	pSendType := acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_UNSPECIFIED
	switch sig.SendType() {
	case SignalSendTypeCyclic:
		pSendType = acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_CYCLIC
	case SignalSendTypeOnWrite:
		pSendType = acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_ON_WRITE
	case SignalSendTypeOnWriteWithRepetition:
		pSendType = acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_ON_WRITE_WITH_REPETITION
	case SignalSendTypeOnChange:
		pSendType = acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_ON_CHANGE
	case SignalSendTypeOnChangeWithRepetition:
		pSendType = acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_ON_CHANGE_WITH_REPETITION
	case SignalSendTypeIfActive:
		pSendType = acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_IF_ACTIVE
	case SignalSendTypeIfActiveWithRepetition:
		pSendType = acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_IF_ACTIVE_WITH_REPETITION
	}
	pSig.SendType = pSendType

	pSig.StartValue = int64(sig.StartValue())

	pKind := acmelibv1.SignalKind_SIGNAL_KIND_UNSPECIFIED
	switch sig.Kind() {
	case SignalKindStandard:
		stdSig, err := sig.ToStandard()
		if err != nil {
			panic(err)
		}

		pKind = acmelibv1.SignalKind_SIGNAL_KIND_STANDARD
		pSig.Entity = s.saveEntity(stdSig.entity)
		pSig.Signal = &acmelibv1.Signal_Standard{
			Standard: s.saveStandardSignal(stdSig),
		}

	case SignalKindEnum:
		enumSig, err := sig.ToEnum()
		if err != nil {
			panic(err)
		}
		s.saveEnumSignal(enumSig)
		pKind = acmelibv1.SignalKind_SIGNAL_KIND_ENUM

	case SignalKindMultiplexer:
		muxSig, err := sig.ToMultiplexer()
		if err != nil {
			panic(err)
		}
		s.saveMultiplexerSignal(muxSig)
		pKind = acmelibv1.SignalKind_SIGNAL_KIND_MULTIPLEXER
	}
	pSig.Kind = pKind

	return pSig
}

func (s *saver) saveStandardSignal(stdSig *StandardSignal) *acmelibv1.StandardSignal {
	pStdSig := new(acmelibv1.StandardSignal)

	if stdSig.typ.ReferenceCount() > 1 {
		entID := stdSig.typ.entityID
		pStdSig.Type = &acmelibv1.StandardSignal_TypeEntityId{
			TypeEntityId: entID.String(),
		}
		s.refSigTypes[entID] = stdSig.typ

	} else {
		pStdSig.Type = &acmelibv1.StandardSignal_EmbeddedType{
			EmbeddedType: s.saveSignalType(stdSig.typ),
		}
	}

	if stdSig.unit == nil {
		return pStdSig
	}

	if stdSig.unit.ReferenceCount() > 1 {
		entID := stdSig.unit.entityID
		pStdSig.Type = &acmelibv1.StandardSignal_TypeEntityId{
			TypeEntityId: entID.String(),
		}
		s.refSigUnits[entID] = stdSig.unit

		return pStdSig
	}

	pStdSig.Unit = &acmelibv1.StandardSignal_EmbeddedUnit{
		EmbeddedUnit: s.saveSignalUnit(stdSig.unit),
	}

	return pStdSig
}

func (s *saver) saveEnumSignal(enumSig *EnumSignal) {

}

func (s *saver) saveMultiplexerSignal(muxSig *MultiplexerSignal) {

}

func (s *saver) saveSignalType(sigType *SignalType) *acmelibv1.SignalType {
	pSigType := new(acmelibv1.SignalType)

	pSigType.Entity = s.saveEntity(sigType.entity)

	pKind := acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_UNSPECIFIED
	switch sigType.kind {
	case SignalTypeKindCustom:
		pKind = acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_CUSTOM
	case SignalTypeKindFlag:
		pKind = acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_FLAG
	case SignalTypeKindInteger:
		pKind = acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_INTEGER
	case SignalTypeKindDecimal:
		pKind = acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_DECIMAL
	}
	pSigType.Kind = pKind

	pSigType.Size = uint32(sigType.size)
	pSigType.Signed = sigType.signed
	pSigType.Min = sigType.min
	pSigType.Max = sigType.max
	pSigType.Scale = sigType.scale
	pSigType.Offset = sigType.offset

	return pSigType
}

func (s *saver) saveSignalUnit(sigUnit *SignalUnit) *acmelibv1.SignalUnit {
	pSigUnit := new(acmelibv1.SignalUnit)

	pSigUnit.Entity = s.saveEntity(sigUnit.entity)

	pKind := acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_UNSPECIFIED
	switch sigUnit.kind {
	case SignalUnitKindCustom:
		pKind = acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_CUSTOM
	case SignalUnitKindTemperature:
		pKind = acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_TEMPERATURE
	case SignalUnitKindElectrical:
		pKind = acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_ELECTRICAL
	case SignalUnitKindPower:
		pKind = acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_POWER
	}
	pSigUnit.Kind = pKind

	pSigUnit.Symbol = sigUnit.symbol

	return pSigUnit
}

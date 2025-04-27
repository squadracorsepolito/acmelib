package acmelib

import (
	"io"
	"slices"
	"strings"

	acmelibv1 "github.com/squadracorsepolito/acmelib/proto/gen/go/acmelib/v1"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SaveEncoding defines the encoding used to save a [Network].
type SaveEncoding uint

const (
	// SaveEncodingWire defines a wire encoding.
	SaveEncodingWire SaveEncoding = 1 << iota
	// SaveEncodingJSON defines a JSON encoding.
	SaveEncodingJSON
	// SaveEncodingText defines a text encoding.
	SaveEncodingText
)

// SaveNetwork saves the given [Network] to the [io.Writer] specified
// by the encoding. It is possible to select more than one encoding by
// using the "|" operator.
//
// It returns an [ArgumentError] that wraps an [ErrIsNil] if the selected
// writer is nil, or a proto/protojson/prototext error if marshal function fails.
func SaveNetwork(network *Network, encoding SaveEncoding, wWire, wJSON, wText io.Writer) error {
	saver := newSaver()
	protoNet := saver.saveNetwork(network)

	if encoding&SaveEncodingWire == SaveEncodingWire {
		if wWire == nil {
			return &ArgumentError{
				Name: "wWire",
				Err:  ErrIsNil,
			}
		}

		data, err := proto.Marshal(protoNet)
		if err != nil {
			return err
		}
		_, err = wWire.Write(data)
		if err != nil {
			return err
		}
	}

	if encoding&SaveEncodingJSON == SaveEncodingJSON {
		if wJSON == nil {
			return &ArgumentError{
				Name: "wJSON",
				Err:  ErrIsNil,
			}
		}

		data, err := protojson.MarshalOptions{Multiline: true}.Marshal(protoNet)
		if err != nil {
			return err
		}
		_, err = wJSON.Write(data)
		if err != nil {
			return err
		}
	}

	if encoding&SaveEncodingText == SaveEncodingText {
		if wText == nil {
			return &ArgumentError{
				Name: "wText",
				Err:  ErrIsNil,
			}
		}

		data, err := prototext.MarshalOptions{Multiline: true}.Marshal(protoNet)
		if err != nil {
			return err
		}
		_, err = wText.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}

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

	canIDBuilders := maps.Values(s.refCANIDBuilders)
	slices.SortFunc(canIDBuilders, func(a, b *CANIDBuilder) int { return strings.Compare(a.name, b.name) })
	for _, canIDBuilder := range canIDBuilders {
		pNet.CanidBuilders = append(pNet.CanidBuilders, s.saveCANIDBuilder(canIDBuilder))
	}

	nodes := maps.Values(s.refNodes)
	slices.SortFunc(nodes, func(a, b *Node) int { return int(a.id - b.id) })
	for _, node := range nodes {
		pNet.Nodes = append(pNet.Nodes, s.saveNode(node))
	}

	sigTypes := maps.Values(s.refSigTypes)
	slices.SortFunc(sigTypes, func(a, b *SignalType) int { return strings.Compare(a.name, b.name) })
	for _, sigType := range sigTypes {
		pNet.SignalTypes = append(pNet.SignalTypes, s.saveSignalType(sigType))
	}

	sigUnits := maps.Values(s.refSigUnits)
	slices.SortFunc(sigUnits, func(a, b *SignalUnit) int { return strings.Compare(a.name, b.name) })
	for _, sigUnit := range sigUnits {
		pNet.SignalUnits = append(pNet.SignalUnits, s.saveSignalUnit(sigUnit))
	}

	sigEnums := maps.Values(s.refSigEnums)
	slices.SortFunc(sigEnums, func(a, b *SignalEnum) int { return strings.Compare(a.name, b.name) })
	for _, sigEnum := range sigEnums {
		pNet.SignalEnums = append(pNet.SignalEnums, s.saveSignalEnum(sigEnum))
	}

	attributes := maps.Values(s.refAttributes)
	slices.SortFunc(attributes, func(a, b Attribute) int { return strings.Compare(a.Name(), b.Name()) })
	for _, att := range attributes {
		pNet.Attributes = append(pNet.Attributes, s.saveAttribute(att))
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

		switch tmpAtt.Type() {
		case AttributeTypeString, AttributeTypeEnum:
			pTmpAttAss.Value = &acmelibv1.AttributeAssignment_ValueString{
				ValueString: tmpAttAss.value.(string),
			}
		case AttributeTypeInteger:
			pTmpAttAss.Value = &acmelibv1.AttributeAssignment_ValueInt{
				ValueInt: int32(tmpAttAss.value.(int)),
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

	entID := bus.canIDBuilder.entityID
	s.refCANIDBuilders[entID] = bus.canIDBuilder
	pBus.CanidBuilderEntityId = string(entID)

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

	pNodeint.Number = int32(nodeInt.number)

	nodeEntID := nodeInt.node.entityID
	s.refNodes[nodeEntID] = nodeInt.node
	pNodeint.NodeEntityId = nodeEntID.String()

	for _, msg := range nodeInt.SentMessages() {
		pNodeint.Messages = append(pNodeint.Messages, s.saveMessage(msg))
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

	pMsg.Payload = s.saveSignalLayout(msg.layout)

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
	case EndiannessLittleEndian:
		pByteOrder = acmelibv1.MessageByteOrder_MESSAGE_BYTE_ORDER_LITTLE_ENDIAN
	case EndiannessBigEndian:
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
		pMsg.Receivers = append(pMsg.Receivers, &acmelibv1.MessageReceiver{
			NodeEntityId:        rec.node.entityID.String(),
			NodeInterfaceNumber: uint32(rec.number),
		})
	}

	return pMsg
}

func (s *saver) saveSignalLayout(layout *SL) *acmelibv1.SignalPayload {
	pPayload := new(acmelibv1.SignalPayload)

	for _, sig := range layout.Signals() {
		pPayload.Refs = append(pPayload.Refs, &acmelibv1.SignalPayloadRef{
			SignalEntityId: sig.EntityID().String(),
			RelStartBit:    uint32(sig.GetStartPos()),
		})
	}

	return pPayload
}

// func (s *saver) saveSignalPayload(payload *signalPayload) *acmelibv1.SignalPayload {
// 	pPayload := new(acmelibv1.SignalPayload)

// 	for _, sig := range payload.signals {
// 		pPayload.Refs = append(pPayload.Refs, &acmelibv1.SignalPayloadRef{
// 			SignalEntityId: sig.EntityID().String(),
// 			RelStartBit:    uint32(sig.GetRelativeStartPos()),
// 		})
// 	}

// 	return pPayload
// }

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

	pSig.StartValue = sig.StartValue()

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

		pKind = acmelibv1.SignalKind_SIGNAL_KIND_ENUM
		pSig.Entity = s.saveEntity(enumSig.entity)
		pSig.Signal = &acmelibv1.Signal_Enum{
			Enum: s.saveEnumSignal(enumSig),
		}

		// case SignalKindMultiplexer:
		// 	muxSig, err := sig.ToMultiplexer()
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	pKind = acmelibv1.SignalKind_SIGNAL_KIND_MULTIPLEXER
		// 	pSig.Entity = s.saveEntity(muxSig.entity)
		// 	pSig.Signal = &acmelibv1.Signal_Multiplexer{
		// 		Multiplexer: s.saveMultiplexerSignal(muxSig),
		// 	}
	}
	pSig.Kind = pKind

	return pSig
}

func (s *saver) saveStandardSignal(stdSig *StandardSignal) *acmelibv1.StandardSignal {
	pStdSig := new(acmelibv1.StandardSignal)

	typeEntID := stdSig.typ.entityID
	s.refSigTypes[typeEntID] = stdSig.typ
	pStdSig.TypeEntityId = string(typeEntID)

	if stdSig.unit == nil {
		return pStdSig
	}

	unitEntID := stdSig.unit.entityID
	s.refSigUnits[unitEntID] = stdSig.unit
	pStdSig.UnitEntityId = string(unitEntID)

	return pStdSig
}

func (s *saver) saveEnumSignal(enumSig *EnumSignal) *acmelibv1.EnumSignal {
	pEnumSig := new(acmelibv1.EnumSignal)

	entID := enumSig.enum.entityID
	s.refSigEnums[entID] = enumSig.enum
	pEnumSig.EnumEntityId = string(entID)

	return pEnumSig
}

// func (s *saver) saveMultiplexerSignal(muxSig *MultiplexerSignal) *acmelibv1.MultiplexerSignal {
// 	// TODO!
// 	pMuxSig := new(acmelibv1.MultiplexerSignal)

// 	// pMuxSig.GroupCount = uint32(muxSig.groupCount)
// 	// pMuxSig.GroupSize = uint32(muxSig.groupSize)

// 	// insFixedSignal := make(map[EntityID]bool)
// 	// for groupID, group := range muxSig.GetSignalGroups() {
// 	// 	for _, muxedSig := range group {
// 	// 		entID := muxedSig.EntityID()

// 	// 		if muxSig.fixedSignals.hasKey(entID) {
// 	// 			if _, ok := insFixedSignal[entID]; ok {
// 	// 				continue
// 	// 			}

// 	// 			insFixedSignal[entID] = true
// 	// 			pMuxSig.FixedSignalEntityIds = append(pMuxSig.FixedSignalEntityIds, entID.String())
// 	// 		}

// 	// 		pMuxSig.Signals = append(pMuxSig.Signals, s.saveSignal(muxedSig))
// 	// 	}

// 	// 	pMuxSig.Groups = append(pMuxSig.Groups, s.saveSignalLayout(muxSig.groups[groupID]))
// 	// }

// 	return pMuxSig
// }

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

func (s *saver) saveSignalEnum(sigEnum *SignalEnum) *acmelibv1.SignalEnum {
	pSigEnum := new(acmelibv1.SignalEnum)

	pSigEnum.Entity = s.saveEntity(sigEnum.entity)

	for _, val := range sigEnum.Values() {
		pSigEnum.Values = append(pSigEnum.Values, s.saveSignalENumValue(val))
	}

	if sigEnum.minSize != 0 {
		pSigEnum.MinSize = uint32(sigEnum.minSize)
	}

	return pSigEnum
}

func (s *saver) saveSignalENumValue(val *SignalEnumValue) *acmelibv1.SignalEnumValue {
	pVal := new(acmelibv1.SignalEnumValue)

	pVal.Entity = s.saveEntity(val.entity)
	pVal.Index = uint32(val.index)

	return pVal
}

func (s *saver) saveAttribute(att Attribute) *acmelibv1.Attribute {
	pAtt := new(acmelibv1.Attribute)

	switch att.Type() {
	case AttributeTypeString:
		strAtt, err := att.ToString()
		if err != nil {
			panic(err)
		}

		pAtt.Type = acmelibv1.AttributeType_ATTRIBUTE_TYPE_STRING
		pAtt.Entity = s.saveEntity(strAtt.entity)
		pAtt.Attribute = &acmelibv1.Attribute_StringAttribute{
			StringAttribute: s.saveStringAttribute(strAtt),
		}

	case AttributeTypeInteger:
		intAtt, err := att.ToInteger()
		if err != nil {
			panic(err)
		}

		pAtt.Type = acmelibv1.AttributeType_ATTRIBUTE_TYPE_INTEGER
		pAtt.Entity = s.saveEntity(intAtt.entity)
		pAtt.Attribute = &acmelibv1.Attribute_IntegerAttribute{
			IntegerAttribute: s.saveIntegerAttribute(intAtt),
		}

	case AttributeTypeFloat:
		floatAtt, err := att.ToFloat()
		if err != nil {
			panic(err)
		}

		pAtt.Type = acmelibv1.AttributeType_ATTRIBUTE_TYPE_FLOAT
		pAtt.Entity = s.saveEntity(floatAtt.entity)
		pAtt.Attribute = &acmelibv1.Attribute_FloatAttribute{
			FloatAttribute: s.saveFloatAttribute(floatAtt),
		}

	case AttributeTypeEnum:
		enumAtt, err := att.ToEnum()
		if err != nil {
			panic(err)
		}

		pAtt.Type = acmelibv1.AttributeType_ATTRIBUTE_TYPE_ENUM
		pAtt.Entity = s.saveEntity(enumAtt.entity)
		pAtt.Attribute = &acmelibv1.Attribute_EnumAttribute{
			EnumAttribute: s.saveEnumAttribute(enumAtt),
		}
	}

	return pAtt
}

func (s *saver) saveStringAttribute(strAtt *StringAttribute) *acmelibv1.StringAttribute {
	pStrAtt := new(acmelibv1.StringAttribute)
	pStrAtt.DefValue = strAtt.defValue
	return pStrAtt
}

func (s *saver) saveIntegerAttribute(intAtt *IntegerAttribute) *acmelibv1.IntegerAttribute {
	pIntAtt := new(acmelibv1.IntegerAttribute)

	pIntAtt.DefValue = int32(intAtt.defValue)
	pIntAtt.Min = int32(intAtt.min)
	pIntAtt.Max = int32(intAtt.max)
	pIntAtt.IsHexFormat = intAtt.isHexFormat

	return pIntAtt
}

func (s *saver) saveFloatAttribute(floatAtt *FloatAttribute) *acmelibv1.FloatAttribute {
	pFloatAtt := new(acmelibv1.FloatAttribute)

	pFloatAtt.DefValue = floatAtt.defValue
	pFloatAtt.Min = floatAtt.min
	pFloatAtt.Max = floatAtt.max

	return pFloatAtt
}

func (s *saver) saveEnumAttribute(enumAtt *EnumAttribute) *acmelibv1.EnumAttribute {
	pEnumAtt := new(acmelibv1.EnumAttribute)

	pEnumAtt.DefValue = enumAtt.defValue

	for _, val := range enumAtt.Values() {
		pEnumAtt.Values = append(pEnumAtt.Values, val)
	}

	return pEnumAtt
}

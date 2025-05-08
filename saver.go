package acmelib

import (
	"io"
	"slices"
	"strings"

	acmelibv2 "github.com/squadracorsepolito/acmelib/gen/acmelib/v2"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SaveNetworkOptions defines the options used to save a [Network].
// There is a field for each supported [SaveEncoding].
type SaveNetworkOptions struct {
	WireWriter, JSONWriter, TextWriter io.Writer
}

// SaveNetwork saves the given [Network] to the [io.Writer] specified
// by the encoding. It is possible to select more than one encoding by
// using the "|" operator.
//
// It returns an [ArgError] that wraps an [ErrIsNil] if the selected
// writer is nil, or a proto/protojson/prototext error if marshal function fails.

// SaveNetwork saves the given [Network] to the writers of the given [SaveNetworkOptions].
//
// It returns [ArgError] if all writers are nil.
func SaveNetwork(network *Network, opts *SaveNetworkOptions) error {
	saver := newSaver()
	protoNet := saver.saveNetwork(network)

	wWire := opts.WireWriter
	wJSON := opts.JSONWriter
	wText := opts.TextWriter

	if wWire == nil && wJSON == nil && wText == nil {
		return newArgError("wWire, wJSON, wText", ErrIsNil)
	}

	if wWire != nil {
		data, err := proto.Marshal(protoNet)
		if err != nil {
			return err
		}
		_, err = wWire.Write(data)
		if err != nil {
			return err
		}
	}

	if wJSON != nil {
		data, err := protojson.MarshalOptions{Multiline: true}.Marshal(protoNet)
		if err != nil {
			return err
		}
		_, err = wJSON.Write(data)
		if err != nil {
			return err
		}
	}

	if wText != nil {
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

func (s *saver) getEntityKind(ek EntityKind) acmelibv2.EntityKind {
	switch ek {
	case EntityKindNetwork:
		return acmelibv2.EntityKind_ENTITY_KIND_NETWORK
	case EntityKindBus:
		return acmelibv2.EntityKind_ENTITY_KIND_BUS
	case EntityKindNode:
		return acmelibv2.EntityKind_ENTITY_KIND_NODE
	case EntityKindMessage:
		return acmelibv2.EntityKind_ENTITY_KIND_MESSAGE
	case EntityKindSignal:
		return acmelibv2.EntityKind_ENTITY_KIND_SIGNAL
	case EntityKindSignalType:
		return acmelibv2.EntityKind_ENTITY_KIND_SIGNAL_TYPE
	case EntityKindSignalUnit:
		return acmelibv2.EntityKind_ENTITY_KIND_SIGNAL_UNIT
	case EntityKindSignalEnum:
		return acmelibv2.EntityKind_ENTITY_KIND_SIGNAL_ENUM
	case EntityKindAttribute:
		return acmelibv2.EntityKind_ENTITY_KIND_ATTRIBUTE
	case EntityKindCANIDBuilder:
		return acmelibv2.EntityKind_ENTITY_KIND_CANID_BUILDER
	default:
		return acmelibv2.EntityKind_ENTITY_KIND_UNSPECIFIED
	}
}

func (s *saver) saveEntity(e *entity) *acmelibv2.Entity {
	pEnt := new(acmelibv2.Entity)

	pEnt.EntityId = e.entityID.String()
	pEnt.Desc = e.desc
	pEnt.Name = e.name
	pEnt.CreateTime = timestamppb.New(e.createTime)
	pEnt.EntityKind = s.getEntityKind(e.entityKind)

	return pEnt
}

func (s *saver) saveNetwork(net *Network) *acmelibv2.Network {
	pNet := new(acmelibv2.Network)

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

func (s *saver) saveAttributeAssignments(attAss []*AttributeAssignment) []*acmelibv2.AttributeAssignment {
	pAttAss := []*acmelibv2.AttributeAssignment{}

	for _, tmpAttAss := range attAss {
		tmpAtt := tmpAttAss.attribute

		pTmpAttAss := new(acmelibv2.AttributeAssignment)

		pTmpAttAss.EntityId = tmpAttAss.EntityID().String()
		pTmpAttAss.AttributeEntityId = tmpAtt.EntityID().String()

		switch tmpAtt.Type() {
		case AttributeTypeString, AttributeTypeEnum:
			pTmpAttAss.Value = &acmelibv2.AttributeAssignment_ValueString{
				ValueString: tmpAttAss.value.(string),
			}
		case AttributeTypeInteger:
			pTmpAttAss.Value = &acmelibv2.AttributeAssignment_ValueInt{
				ValueInt: int32(tmpAttAss.value.(int)),
			}
		case AttributeTypeFloat:
			pTmpAttAss.Value = &acmelibv2.AttributeAssignment_ValueDouble{
				ValueDouble: tmpAttAss.value.(float64),
			}
		}

		pAttAss = append(pAttAss, pTmpAttAss)

		s.refAttributes[tmpAtt.EntityID()] = tmpAtt
	}

	return pAttAss
}

func (s *saver) saveBus(bus *Bus) *acmelibv2.Bus {
	pBus := new(acmelibv2.Bus)

	pBus.Entity = s.saveEntity(bus.entity)
	pBus.AttributeAssignments = s.saveAttributeAssignments(bus.AttributeAssignments())

	pBus.Baudrate = uint32(bus.baudrate)

	pBusType := acmelibv2.BusType_BUS_TYPE_UNSPECIFIED
	switch bus.typ {
	case BusTypeCAN2A:
		pBusType = acmelibv2.BusType_BUS_TYPE_CAN_2A
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

func (s *saver) saveCANIDBuilder(builder *CANIDBuilder) *acmelibv2.CANIDBuilder {
	pBuilder := new(acmelibv2.CANIDBuilder)

	pBuilder.Entity = s.saveEntity(builder.entity)

	for _, tmpOp := range builder.operations {
		pBuilder.Operations = append(pBuilder.Operations, s.saveCANIDBuilderOp(tmpOp))
	}

	return pBuilder
}

func (s *saver) saveCANIDBuilderOp(builderOp *CANIDBuilderOp) *acmelibv2.CANIDBuilderOp {
	pBuilderOp := new(acmelibv2.CANIDBuilderOp)

	pOpKind := acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_UNSPECIFIED
	switch builderOp.kind {
	case CANIDBuilderOpKindMessagePriority:
		pOpKind = acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_PRIORITY
	case CANIDBuilderOpKindMessageID:
		pOpKind = acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_ID
	case CANIDBuilderOpKindNodeID:
		pOpKind = acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_NODE_ID
	case CANIDBuilderOpKindBitMask:
		pOpKind = acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_BIT_MASK
	}
	pBuilderOp.Kind = pOpKind

	pBuilderOp.From = uint32(builderOp.from)
	pBuilderOp.Len = uint32(builderOp.len)

	return pBuilderOp
}

func (s *saver) saveNode(node *Node) *acmelibv2.Node {
	pNode := new(acmelibv2.Node)

	pNode.Entity = s.saveEntity(node.entity)
	pNode.AttributeAssignments = s.saveAttributeAssignments(node.AttributeAssignments())

	pNode.NodeId = uint32(node.id)
	pNode.InterfaceCount = uint32(node.interfaceCount)

	return pNode
}

func (s *saver) saveNodeInterface(nodeInt *NodeInterface) *acmelibv2.NodeInterface {
	pNodeint := new(acmelibv2.NodeInterface)

	pNodeint.Number = int32(nodeInt.number)

	nodeEntID := nodeInt.node.entityID
	s.refNodes[nodeEntID] = nodeInt.node
	pNodeint.NodeEntityId = nodeEntID.String()

	for _, msg := range nodeInt.SentMessages() {
		pNodeint.Messages = append(pNodeint.Messages, s.saveMessage(msg))
	}

	return pNodeint
}

func (s *saver) saveMessage(msg *Message) *acmelibv2.Message {
	pMsg := new(acmelibv2.Message)

	pMsg.Entity = s.saveEntity(msg.entity)
	pMsg.AttributeAssignments = s.saveAttributeAssignments(msg.AttributeAssignments())

	pMsg.SizeByte = uint32(msg.sizeByte)
	pMsg.MessageId = uint32(msg.id)

	pMsg.Layout = s.saveSignalLayout(msg.layout)

	pMsg.StaticCanId = uint32(msg.staticCANID)
	pMsg.HasStaticCanId = msg.hasStaticCANID

	pPriority := acmelibv2.MessagePriority_MESSAGE_PRIORITY_UNSPECIFIED
	switch msg.priority {
	case MessagePriorityVeryHigh:
		pPriority = acmelibv2.MessagePriority_MESSAGE_PRIORITY_VERY_HIGH
	case MessagePriorityHigh:
		pPriority = acmelibv2.MessagePriority_MESSAGE_PRIORITY_HIGH
	case MessagePriorityMedium:
		pPriority = acmelibv2.MessagePriority_MESSAGE_PRIORITY_MEDIUM
	case MessagePriorityLow:
		pPriority = acmelibv2.MessagePriority_MESSAGE_PRIORITY_LOW
	}
	pMsg.Priority = pPriority

	pMsg.CycleTime = uint32(msg.cycleTime)

	pSendType := acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_UNSPECIFIED
	switch msg.sendType {
	case MessageSendTypeCyclic:
		pSendType = acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC
	case MessageSendTypeCyclicIfActive:
		pSendType = acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE
	case MessageSendTypeCyclicAndTriggered:
		pSendType = acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_AND_TRIGGERED
	case MessageSendTypeCyclicIfActiveAndTriggered:
		pSendType = acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE_AND_TRIGGERED
	}
	pMsg.SendType = pSendType

	pMsg.DelayTime = uint32(msg.delayTime)
	pMsg.StartDelayTime = uint32(msg.startDelayTime)

	for _, rec := range msg.Receivers() {
		pMsg.Receivers = append(pMsg.Receivers, &acmelibv2.MessageReceiver{
			NodeEntityId:        rec.node.entityID.String(),
			NodeInterfaceNumber: uint32(rec.number),
		})
	}

	return pMsg
}

func (s *saver) saveSignalLayout(layout *SignalLayout) *acmelibv2.SignalLayout {
	pLayout := new(acmelibv2.SignalLayout)

	pLayout.Id = uint32(layout.id)
	pLayout.SizeByte = uint32(layout.sizeByte)

	for sig := range layout.ibst.InOrder() {
		if sig.Kind() == SignalKindMuxor {
			continue
		}

		pLayout.Signals = append(pLayout.Signals, s.saveSignal(sig))
	}

	for muxLayer := range layout.muxLayers.Values() {
		pLayout.MultiplexedLayers = append(pLayout.MultiplexedLayers, s.saveMultiplexedLayer(muxLayer))
	}

	return pLayout
}

func (s *saver) saveMultiplexedLayer(muxLayer *MultiplexedLayer) *acmelibv2.MultiplexedLayer {
	pMuxLayer := new(acmelibv2.MultiplexedLayer)

	muxor := s.saveSignal(muxLayer.muxor)
	pMuxLayer.Muxor = muxor

	for _, layout := range muxLayer.layouts {
		if layout.ibst.Size() == 0 {
			continue
		}

		pMuxLayer.Layouts = append(pMuxLayer.Layouts, s.saveSignalLayout(layout))
	}

	return pMuxLayer
}

func (s *saver) saveSignal(sig Signal) *acmelibv2.Signal {
	pSig := new(acmelibv2.Signal)

	pSig.AttributeAssignments = s.saveAttributeAssignments(sig.AttributeAssignments())

	pSig.StartPos = uint32(sig.StartPos())

	pEndianness := acmelibv2.Endianness_ENDIANNESS_UNSPECIFIED
	switch sig.Endianness() {
	case EndiannessLittleEndian:
		pEndianness = acmelibv2.Endianness_ENDIANNESS_LITTLE_ENDIAN
	case EndiannessBigEndian:
		pEndianness = acmelibv2.Endianness_ENDIANNESS_BIG_ENDIAN
	}
	pSig.Endianness = pEndianness

	pSendType := acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_UNSPECIFIED
	switch sig.SendType() {
	case SignalSendTypeCyclic:
		pSendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_CYCLIC
	case SignalSendTypeOnWrite:
		pSendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_WRITE
	case SignalSendTypeOnWriteWithRepetition:
		pSendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_WRITE_WITH_REPETITION
	case SignalSendTypeOnChange:
		pSendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_CHANGE
	case SignalSendTypeOnChangeWithRepetition:
		pSendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_CHANGE_WITH_REPETITION
	case SignalSendTypeIfActive:
		pSendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_IF_ACTIVE
	case SignalSendTypeIfActiveWithRepetition:
		pSendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_IF_ACTIVE_WITH_REPETITION
	}
	pSig.SendType = pSendType

	pSig.StartValue = sig.StartValue()

	pKind := acmelibv2.SignalKind_SIGNAL_KIND_UNSPECIFIED
	switch sig.Kind() {
	case SignalKindStandard:
		stdSig, err := sig.ToStandard()
		if err != nil {
			panic(err)
		}

		pKind = acmelibv2.SignalKind_SIGNAL_KIND_STANDARD
		pSig.Entity = s.saveEntity(stdSig.entity)
		pSig.Signal = &acmelibv2.Signal_Standard{
			Standard: s.saveStandardSignal(stdSig),
		}

	case SignalKindEnum:
		enumSig, err := sig.ToEnum()
		if err != nil {
			panic(err)
		}

		pKind = acmelibv2.SignalKind_SIGNAL_KIND_ENUM
		pSig.Entity = s.saveEntity(enumSig.entity)
		pSig.Signal = &acmelibv2.Signal_Enum{
			Enum: s.saveEnumSignal(enumSig),
		}

	case SignalKindMuxor:
		muxorSig, err := sig.ToMuxor()
		if err != nil {
			panic(err)
		}

		pKind = acmelibv2.SignalKind_SIGNAL_KIND_MUXOR
		pSig.Entity = s.saveEntity(muxorSig.entity)
		pSig.Signal = &acmelibv2.Signal_Muxor{
			Muxor: s.saveMuxorSignal(muxorSig),
		}
	}
	pSig.Kind = pKind

	return pSig
}

func (s *saver) saveStandardSignal(stdSig *StandardSignal) *acmelibv2.StandardSignal {
	pStdSig := new(acmelibv2.StandardSignal)

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

func (s *saver) saveEnumSignal(enumSig *EnumSignal) *acmelibv2.EnumSignal {
	pEnumSig := new(acmelibv2.EnumSignal)

	entID := enumSig.enum.entityID
	s.refSigEnums[entID] = enumSig.enum
	pEnumSig.EnumEntityId = string(entID)

	return pEnumSig
}

func (s *saver) saveMuxorSignal(muxorSig *MuxorSignal) *acmelibv2.MuxorSignal {
	pMuxorSig := new(acmelibv2.MuxorSignal)
	pMuxorSig.LayoutCount = uint32(muxorSig.layoutCount)
	return pMuxorSig
}

func (s *saver) saveSignalType(sigType *SignalType) *acmelibv2.SignalType {
	pSigType := new(acmelibv2.SignalType)

	pSigType.Entity = s.saveEntity(sigType.entity)

	pKind := acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_UNSPECIFIED
	switch sigType.kind {
	case SignalTypeKindFlag:
		pKind = acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_FLAG
	case SignalTypeKindInteger:
		pKind = acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_INTEGER
	case SignalTypeKindDecimal:
		pKind = acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_DECIMAL
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

func (s *saver) saveSignalUnit(sigUnit *SignalUnit) *acmelibv2.SignalUnit {
	pSigUnit := new(acmelibv2.SignalUnit)

	pSigUnit.Entity = s.saveEntity(sigUnit.entity)

	pKind := acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_UNSPECIFIED
	switch sigUnit.kind {
	case SignalUnitKindCustom:
		pKind = acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_CUSTOM
	case SignalUnitKindTemperature:
		pKind = acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_TEMPERATURE
	case SignalUnitKindElectrical:
		pKind = acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_ELECTRICAL
	case SignalUnitKindPower:
		pKind = acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_POWER
	}
	pSigUnit.Kind = pKind

	pSigUnit.Symbol = sigUnit.symbol

	return pSigUnit
}

func (s *saver) saveSignalEnum(sigEnum *SignalEnum) *acmelibv2.SignalEnum {
	pSigEnum := new(acmelibv2.SignalEnum)

	pSigEnum.Entity = s.saveEntity(sigEnum.entity)

	pSigEnum.Size = uint32(sigEnum.size)
	pSigEnum.FixedSize = sigEnum.fixedSize

	for _, val := range sigEnum.Values() {
		pSigEnum.Values = append(pSigEnum.Values, s.saveSignalENumValue(val))
	}

	return pSigEnum
}

func (s *saver) saveSignalENumValue(val *SignalEnumValue) *acmelibv2.SignalEnumValue {
	pVal := new(acmelibv2.SignalEnumValue)

	pVal.Name = val.name
	pVal.Index = uint32(val.index)

	if val.desc != "" {
		pVal.Desc = val.desc
	}

	return pVal
}

func (s *saver) saveAttribute(att Attribute) *acmelibv2.Attribute {
	pAtt := new(acmelibv2.Attribute)

	switch att.Type() {
	case AttributeTypeString:
		strAtt, err := att.ToString()
		if err != nil {
			panic(err)
		}

		pAtt.Type = acmelibv2.AttributeType_ATTRIBUTE_TYPE_STRING
		pAtt.Entity = s.saveEntity(strAtt.entity)
		pAtt.Attribute = &acmelibv2.Attribute_StringAttribute{
			StringAttribute: s.saveStringAttribute(strAtt),
		}

	case AttributeTypeInteger:
		intAtt, err := att.ToInteger()
		if err != nil {
			panic(err)
		}

		pAtt.Type = acmelibv2.AttributeType_ATTRIBUTE_TYPE_INTEGER
		pAtt.Entity = s.saveEntity(intAtt.entity)
		pAtt.Attribute = &acmelibv2.Attribute_IntegerAttribute{
			IntegerAttribute: s.saveIntegerAttribute(intAtt),
		}

	case AttributeTypeFloat:
		floatAtt, err := att.ToFloat()
		if err != nil {
			panic(err)
		}

		pAtt.Type = acmelibv2.AttributeType_ATTRIBUTE_TYPE_FLOAT
		pAtt.Entity = s.saveEntity(floatAtt.entity)
		pAtt.Attribute = &acmelibv2.Attribute_FloatAttribute{
			FloatAttribute: s.saveFloatAttribute(floatAtt),
		}

	case AttributeTypeEnum:
		enumAtt, err := att.ToEnum()
		if err != nil {
			panic(err)
		}

		pAtt.Type = acmelibv2.AttributeType_ATTRIBUTE_TYPE_ENUM
		pAtt.Entity = s.saveEntity(enumAtt.entity)
		pAtt.Attribute = &acmelibv2.Attribute_EnumAttribute{
			EnumAttribute: s.saveEnumAttribute(enumAtt),
		}
	}

	return pAtt
}

func (s *saver) saveStringAttribute(strAtt *StringAttribute) *acmelibv2.StringAttribute {
	pStrAtt := new(acmelibv2.StringAttribute)
	pStrAtt.DefValue = strAtt.defValue
	return pStrAtt
}

func (s *saver) saveIntegerAttribute(intAtt *IntegerAttribute) *acmelibv2.IntegerAttribute {
	pIntAtt := new(acmelibv2.IntegerAttribute)

	pIntAtt.DefValue = int32(intAtt.defValue)
	pIntAtt.Min = int32(intAtt.min)
	pIntAtt.Max = int32(intAtt.max)
	pIntAtt.IsHexFormat = intAtt.isHexFormat

	return pIntAtt
}

func (s *saver) saveFloatAttribute(floatAtt *FloatAttribute) *acmelibv2.FloatAttribute {
	pFloatAtt := new(acmelibv2.FloatAttribute)

	pFloatAtt.DefValue = floatAtt.defValue
	pFloatAtt.Min = floatAtt.min
	pFloatAtt.Max = floatAtt.max

	return pFloatAtt
}

func (s *saver) saveEnumAttribute(enumAtt *EnumAttribute) *acmelibv2.EnumAttribute {
	pEnumAtt := new(acmelibv2.EnumAttribute)

	pEnumAtt.DefValue = enumAtt.defValue

	for _, val := range enumAtt.Values() {
		pEnumAtt.Values = append(pEnumAtt.Values, val)
	}

	return pEnumAtt
}

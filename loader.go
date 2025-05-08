package acmelib

import (
	"bytes"
	"io"
	"log"
	"time"

	acmelibv2 "github.com/squadracorsepolito/acmelib/gen/acmelib/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
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

// LoadNetwork loads the content of the [io.Reader] and returns a [Network].
// The encoding parameter specifies the encoding of the content of the reader.
func LoadNetwork(r io.Reader, encoding SaveEncoding) (*Network, error) {
	pNetwork := &acmelibv2.Network{}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	switch encoding {
	case SaveEncodingWire:
		if err := proto.Unmarshal(buf.Bytes(), pNetwork); err != nil {
			return nil, err
		}

	case SaveEncodingJSON:
		if err := protojson.Unmarshal(buf.Bytes(), pNetwork); err != nil {
			return nil, err
		}

	case SaveEncodingText:
		if err := prototext.Unmarshal(buf.Bytes(), pNetwork); err != nil {
			return nil, err
		}
	}

	loader := newLoader()

	return loader.loadNetwork(pNetwork)
}

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

func (l *loader) loadEntity(pEnt *acmelibv2.Entity, entKind EntityKind) *entity {
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

func (l *loader) loadNetwork(pNet *acmelibv2.Network) (*Network, error) {
	net := newNetworkFromEntity(l.loadEntity(pNet.Entity, EntityKindNetwork))

	for _, pBuilder := range pNet.CanidBuilders {
		l.refCANIDBuilders[pBuilder.Entity.EntityId] = l.loadCANIDBuilder(pBuilder)
	}

	for _, pAtt := range pNet.Attributes {
		att, err := l.loadAttribute(pAtt)
		if err != nil {
			return nil, err
		}
		l.refAttributes[pAtt.Entity.EntityId] = att
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

func (l *loader) loadCANIDBuilder(pBuilder *acmelibv2.CANIDBuilder) *CANIDBuilder {
	builder := newCANIDBuilderFromEntity(l.loadEntity(pBuilder.Entity, EntityKindCANIDBuilder))

	for _, pBuilderOp := range pBuilder.Operations {
		builder.operations = append(builder.operations, l.loadCANIDBuilderOp(pBuilderOp))
	}

	return builder
}

func (l *loader) loadCANIDBuilderOp(pBuilderOp *acmelibv2.CANIDBuilderOp) *CANIDBuilderOp {
	var kind CANIDBuilderOpKind
	switch pBuilderOp.Kind {
	case acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_PRIORITY:
		kind = CANIDBuilderOpKindMessagePriority
	case acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_ID:
		kind = CANIDBuilderOpKindMessageID
	case acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_NODE_ID:
		kind = CANIDBuilderOpKindNodeID
	case acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_BIT_MASK:
		kind = CANIDBuilderOpKindBitMask
	}
	return newCANIDBuilderOp(kind, int(pBuilderOp.From), int(pBuilderOp.Len))
}

func (l *loader) loadNode(pNode *acmelibv2.Node) (*Node, error) {
	node := newNodeFromEntity(l.loadEntity(pNode.Entity, EntityKindNode), NodeID(pNode.NodeId), int(pNode.InterfaceCount))

	for _, pAttAss := range pNode.AttributeAssignments {
		if err := l.loadAttributeAssignment(node, pAttAss); err != nil {
			return nil, err
		}
	}

	return node, nil
}

func (l *loader) loadBus(pBus *acmelibv2.Bus) (*Bus, error) {
	bus := newBusFromEntity(l.loadEntity(pBus.Entity, EntityKindBus))

	var typ BusType
	switch pBus.Type {
	case acmelibv2.BusType_BUS_TYPE_CAN_2A:
		typ = BusTypeCAN2A
	}
	bus.SetType(typ)

	bus.SetBaudrate(int(pBus.Baudrate))

	for _, pNodeInt := range pBus.NodeInterfaces {
		nodeInt, err := l.loadNodeInterface(pNodeInt)
		if err != nil {
			return nil, err
		}

		if err := bus.AddNodeInterface(nodeInt); err != nil {
			return nil, err
		}
	}

	for _, pAttAss := range pBus.AttributeAssignments {
		if err := l.loadAttributeAssignment(bus, pAttAss); err != nil {
			return nil, err
		}
	}

	return bus, nil
}

func (l *loader) loadNodeInterface(pNodeInt *acmelibv2.NodeInterface) (*NodeInterface, error) {
	node, ok := l.refNodes[pNodeInt.NodeEntityId]
	if !ok {
		return nil, &EntityIDError{
			EntityID: EntityID(pNodeInt.NodeEntityId),
			Err:      ErrNotFound,
		}
	}

	nodeInt := node.GetInterface(int(pNodeInt.Number))
	if nodeInt == nil {
		return nil, ErrNotFound
	}

	for _, pMsg := range pNodeInt.Messages {
		msg, err := l.loadMessage(pMsg)
		if err != nil {
			return nil, err
		}

		if err := nodeInt.AddSentMessage(msg); err != nil {
			return nil, err
		}
	}

	return nodeInt, nil
}

// func (l *loader) loadSignalPayload(pSigPayload *acmelibv2.SignalPayload) map[string]int {
// 	sigMap := make(map[string]int)
// 	for _, pRef := range pSigPayload.Refs {
// 		sigMap[pRef.SignalEntityId] = int(pRef.RelStartBit)
// 	}
// 	return sigMap
// }

func (l *loader) loadMessage(pMsg *acmelibv2.Message) (*Message, error) {
	msg := newMessageFromEntity(l.loadEntity(pMsg.Entity, EntityKindMessage), MessageID(pMsg.MessageId), int(pMsg.SizeByte))

	if pMsg.HasStaticCanId {
		if err := msg.SetStaticCANID(CANID(pMsg.StaticCanId)); err != nil {
			return nil, err
		}
	}

	switch pMsg.Priority {
	case acmelibv2.MessagePriority_MESSAGE_PRIORITY_VERY_HIGH:
		msg.SetPriority(MessagePriorityVeryHigh)
	case acmelibv2.MessagePriority_MESSAGE_PRIORITY_HIGH:
		msg.SetPriority(MessagePriorityHigh)
	case acmelibv2.MessagePriority_MESSAGE_PRIORITY_MEDIUM:
		msg.SetPriority(MessagePriorityMedium)
	case acmelibv2.MessagePriority_MESSAGE_PRIORITY_LOW:
		msg.SetPriority(MessagePriorityLow)
	}

	if pMsg.CycleTime != 0 {
		msg.SetCycleTime(int(pMsg.CycleTime))
	}

	switch pMsg.SendType {
	case acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC:
		msg.SetSendType(MessageSendTypeCyclic)
	case acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE:
		msg.SetSendType(MessageSendTypeCyclicIfActive)
	case acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_AND_TRIGGERED:
		msg.SetSendType(MessageSendTypeCyclicAndTriggered)
	case acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE_AND_TRIGGERED:
		msg.SetSendType(MessageSendTypeCyclicIfActiveAndTriggered)
	}

	if pMsg.DelayTime != 0 {
		msg.SetDelayTime(int(pMsg.DelayTime))
	}

	if pMsg.StartDelayTime != 0 {
		msg.SetStartDelayTime(int(pMsg.StartDelayTime))
	}

	for _, pRec := range pMsg.Receivers {
		recNode, ok := l.refNodes[pRec.NodeEntityId]
		if !ok {
			return nil, &EntityIDError{
				EntityID: EntityID(pRec.NodeEntityId),
				Err:      ErrNotFound,
			}
		}

		recNodeInt := recNode.GetInterface(int(pRec.NodeInterfaceNumber))
		if recNodeInt == nil {
			return nil, ErrNotFound
		}

		msg.AddReceiver(recNodeInt)
	}

	if err := l.loadSignalLayout(msg.layout, pMsg.Layout); err != nil {
		return nil, err
	}

	for _, pAttAss := range pMsg.AttributeAssignments {
		if err := l.loadAttributeAssignment(msg, pAttAss); err != nil {
			return nil, err
		}
	}

	return msg, nil
}

func (l *loader) loadSignal(pSig *acmelibv2.Signal) (Signal, error) {
	var kind SignalKind
	switch pSig.Kind {
	case acmelibv2.SignalKind_SIGNAL_KIND_STANDARD:
		kind = SignalKindStandard
	case acmelibv2.SignalKind_SIGNAL_KIND_ENUM:
		kind = SignalKindEnum
	case acmelibv2.SignalKind_SIGNAL_KIND_MUXOR:
		kind = SignalKindMuxor
	}

	baseSig := newSignalFromEntity(l.loadEntity(pSig.Entity, EntityKindSignal), kind)

	if pSig.Endianness == acmelibv2.Endianness_ENDIANNESS_BIG_ENDIAN {
		baseSig.SetEndianness(EndiannessBigEndian)
	}

	var sig Signal
	switch tmpPSig := pSig.Signal.(type) {
	case *acmelibv2.Signal_Standard:
		if kind != SignalKindStandard {
			return nil, &ErrInvalidOneof{
				KindTypeField: acmelibv2.SignalKind_SIGNAL_KIND_STANDARD.String(),
			}
		}

		stdSig, err := l.loadStandardSignal(baseSig, tmpPSig.Standard)
		if err != nil {
			return nil, err
		}
		sig = stdSig

	case *acmelibv2.Signal_Enum:
		if kind != SignalKindEnum {
			return nil, &ErrInvalidOneof{
				KindTypeField: acmelibv2.SignalKind_SIGNAL_KIND_ENUM.String(),
			}
		}

		enumSig, err := l.loadEnumSignal(baseSig, tmpPSig.Enum)
		if err != nil {
			return nil, err
		}
		sig = enumSig

	case *acmelibv2.Signal_Muxor:
		if kind != SignalKindMuxor {
			return nil, &ErrInvalidOneof{
				KindTypeField: acmelibv2.SignalKind_SIGNAL_KIND_MUXOR.String(),
			}
		}

		muxorSig, err := l.loadMuxorSignal(baseSig, tmpPSig.Muxor)
		if err != nil {
			return nil, err
		}
		sig = muxorSig
	}

	switch pSig.SendType {
	case acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_CYCLIC:
		sig.SetSendType(SignalSendTypeCyclic)
	case acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_WRITE:
		sig.SetSendType(SignalSendTypeOnWrite)
	case acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_WRITE_WITH_REPETITION:
		sig.SetSendType(SignalSendTypeOnWriteWithRepetition)
	case acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_CHANGE:
		sig.SetSendType(SignalSendTypeOnChange)
	case acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_CHANGE_WITH_REPETITION:
		sig.SetSendType(SignalSendTypeOnChangeWithRepetition)
	case acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_IF_ACTIVE:
		sig.SetSendType(SignalSendTypeIfActive)
	case acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_IF_ACTIVE_WITH_REPETITION:
		sig.SetSendType(SignalSendTypeIfActiveWithRepetition)
	}

	if pSig.StartValue != 0 {
		sig.SetStartValue(pSig.StartValue)
	}

	for _, pAttAss := range pSig.AttributeAssignments {
		if err := l.loadAttributeAssignment(sig, pAttAss); err != nil {
			return nil, err
		}
	}

	return sig, nil
}

func (l *loader) loadStandardSignal(baseSig *signal, pStdSig *acmelibv2.StandardSignal) (*StandardSignal, error) {
	sigTyp, ok := l.refSigTypes[pStdSig.TypeEntityId]
	if !ok {
		return nil, &EntityIDError{
			EntityID: EntityID(pStdSig.TypeEntityId),
			Err:      ErrNotFound,
		}
	}

	stdSig, err := newStandardSignalFromBase(baseSig, sigTyp)
	if err != nil {
		return nil, err
	}

	if len(pStdSig.UnitEntityId) > 0 {
		sigUnit, ok := l.refSigUnits[pStdSig.UnitEntityId]
		if !ok {
			return nil, &EntityIDError{
				EntityID: EntityID(pStdSig.UnitEntityId),
				Err:      ErrNotFound,
			}
		}
		stdSig.SetUnit(sigUnit)
	}

	return stdSig, nil
}

func (l *loader) loadEnumSignal(baseSig *signal, pEnumSig *acmelibv2.EnumSignal) (*EnumSignal, error) {
	sigEnum, ok := l.refSigEnums[pEnumSig.EnumEntityId]
	if !ok {
		return nil, &EntityIDError{
			EntityID: EntityID(pEnumSig.EnumEntityId),
			Err:      ErrNotFound,
		}
	}
	return newEnumSignalFromBase(baseSig, sigEnum)
}

func (l *loader) loadMuxorSignal(baseSig *signal, pMuxorSig *acmelibv2.MuxorSignal) (*MuxorSignal, error) {
	return newMuxorSignalFromBase(baseSig, int(pMuxorSig.LayoutCount))
}

func (l *loader) loadSignalLayout(layout *SignalLayout, pLayout *acmelibv2.SignalLayout) error {
	for _, pSig := range pLayout.Signals {
		sig, err := l.loadSignal(pSig)
		if err != nil {
			return err
		}

		startPos := int(pSig.StartPos)

		if layout.fromMessage() {
			if err := layout.parentMsg.InsertSignal(sig, startPos); err != nil {
				return err
			}
		}

		if layout.fromMultiplexedLayer() {
			layoutID := int(pLayout.Id)
			sig.setStartPos(startPos)

			if err := layout.parentMuxLayer.InsertSignal(sig, startPos, layoutID); err != nil {
				return err
			}
		}
	}

	for _, pMuxLayer := range pLayout.MultiplexedLayers {
		sig, err := l.loadSignal(pMuxLayer.Muxor)
		if err != nil {
			return err
		}
		muxor, err := sig.ToMuxor()
		if err != nil {
			return err
		}

		muxLayer, err := layout.AddMultiplexedLayer(muxor, int(pMuxLayer.Muxor.StartPos))
		if err != nil {
			return err
		}

		for _, pInnerLayout := range pMuxLayer.Layouts {
			innerLayout := muxLayer.GetLayout(int(pInnerLayout.Id))
			if innerLayout == nil {
				continue
			}

			if err := l.loadSignalLayout(innerLayout, pInnerLayout); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *loader) loadSignalType(pSigType *acmelibv2.SignalType) (*SignalType, error) {
	var kind SignalTypeKind
	switch pSigType.Kind {
	case acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_FLAG:
		kind = SignalTypeKindFlag
	case acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_INTEGER:
		kind = SignalTypeKindInteger
	case acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_DECIMAL:
		kind = SignalTypeKindDecimal
	}

	ent := l.loadEntity(pSigType.Entity, EntityKindSignalType)
	return newSignalTypeFromEntity(ent, kind, int(pSigType.Size), pSigType.Signed, pSigType.Min, pSigType.Max, pSigType.Scale, pSigType.Offset)
}

func (l *loader) loadSignalUnit(pSigUnit *acmelibv2.SignalUnit) *SignalUnit {
	var kind SignalUnitKind
	switch pSigUnit.Kind {
	case acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_CUSTOM:
		kind = SignalUnitKindCustom
	case acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_TEMPERATURE:
		kind = SignalUnitKindTemperature
	case acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_ELECTRICAL:
		kind = SignalUnitKindElectrical
	case acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_POWER:
		kind = SignalUnitKindPower
	}
	return newSignalUnitFromEntity(l.loadEntity(pSigUnit.Entity, EntityKindSignalUnit), kind, pSigUnit.Symbol)
}

func (l *loader) loadSignalEnum(pSigEnum *acmelibv2.SignalEnum) (*SignalEnum, error) {
	sigEnum := newSignalEnumFromEntity(l.loadEntity(pSigEnum.Entity, EntityKindSignalEnum))

	for _, pVal := range pSigEnum.Values {
		_, err := sigEnum.AddValue(int(pVal.Index), pVal.Name)
		if err != nil {
			return nil, err
		}
	}

	if pSigEnum.FixedSize {
		sigEnum.SetFixedSize(true)
	}

	if err := sigEnum.UpdateSize(int(pSigEnum.Size)); err != nil {
		return nil, err
	}

	return sigEnum, nil
}

func (l *loader) loadAttribute(pAtt *acmelibv2.Attribute) (Attribute, error) {
	var typ AttributeType
	switch pAtt.Type {
	case acmelibv2.AttributeType_ATTRIBUTE_TYPE_STRING:
		typ = AttributeTypeString
	case acmelibv2.AttributeType_ATTRIBUTE_TYPE_INTEGER:
		typ = AttributeTypeInteger
	case acmelibv2.AttributeType_ATTRIBUTE_TYPE_FLOAT:
		typ = AttributeTypeFloat
	case acmelibv2.AttributeType_ATTRIBUTE_TYPE_ENUM:
		typ = AttributeTypeEnum
	}

	baseAtt := newAttributeFromEntity(l.loadEntity(pAtt.Entity, EntityKindAttribute), typ)

	var att Attribute
	switch tmpPAtt := pAtt.Attribute.(type) {
	case *acmelibv2.Attribute_StringAttribute:
		if typ != AttributeTypeString {
			return nil, &ErrInvalidOneof{
				KindTypeField: acmelibv2.AttributeType_ATTRIBUTE_TYPE_STRING.String(),
			}
		}

		strAtt := l.loadStringAttribute(baseAtt, tmpPAtt.StringAttribute)
		att = strAtt

	case *acmelibv2.Attribute_IntegerAttribute:
		if typ != AttributeTypeInteger {
			return nil, &ErrInvalidOneof{
				KindTypeField: acmelibv2.AttributeType_ATTRIBUTE_TYPE_INTEGER.String(),
			}
		}

		intAtt, err := l.loadIntegerAttribute(baseAtt, tmpPAtt.IntegerAttribute)
		if err != nil {
			return nil, err
		}
		att = intAtt

	case *acmelibv2.Attribute_FloatAttribute:
		if typ != AttributeTypeFloat {
			return nil, &ErrInvalidOneof{
				KindTypeField: acmelibv2.AttributeType_ATTRIBUTE_TYPE_FLOAT.String(),
			}
		}

		floatAtt, err := l.loadFloatAttribute(baseAtt, tmpPAtt.FloatAttribute)
		if err != nil {
			return nil, err
		}
		att = floatAtt

	case *acmelibv2.Attribute_EnumAttribute:
		if typ != AttributeTypeEnum {
			return nil, &ErrInvalidOneof{
				KindTypeField: acmelibv2.AttributeType_ATTRIBUTE_TYPE_ENUM.String(),
			}
		}

		enumAtt, err := l.loadEnumAttribute(baseAtt, tmpPAtt.EnumAttribute)
		if err != nil {
			return nil, err
		}
		att = enumAtt

	default:
		return nil, &ErrMissingOneofField{OneofField: "attribute"}
	}

	return att, nil
}

func (l *loader) loadStringAttribute(baseAtt *attribute, pStrAtt *acmelibv2.StringAttribute) *StringAttribute {
	return newStringAttributeFromBase(baseAtt, pStrAtt.DefValue)
}

func (l *loader) loadIntegerAttribute(baseAtt *attribute, pIntAtt *acmelibv2.IntegerAttribute) (*IntegerAttribute, error) {
	intAtt, err := newIntegerAttributeFromBase(baseAtt, int(pIntAtt.DefValue), int(pIntAtt.Min), int(pIntAtt.Max))
	if err != nil {
		return nil, err
	}

	if pIntAtt.IsHexFormat {
		intAtt.SetFormatHex()
	}

	return intAtt, nil
}

func (l *loader) loadFloatAttribute(baseAtt *attribute, pFloatAtt *acmelibv2.FloatAttribute) (*FloatAttribute, error) {
	return newFloatAttributeFromBase(baseAtt, pFloatAtt.DefValue, pFloatAtt.Min, pFloatAtt.Max)
}

func (l *loader) loadEnumAttribute(baseAtt *attribute, pEnumAtt *acmelibv2.EnumAttribute) (*EnumAttribute, error) {
	values := make([]string, len(pEnumAtt.Values))
	values[0] = pEnumAtt.DefValue
	idx := 1
	for _, val := range pEnumAtt.Values {
		if val == pEnumAtt.DefValue {
			continue
		}
		values[idx] = val
		idx++
	}
	return newEnumAttributeFromBase(baseAtt, values...)
}

func (l *loader) loadAttributeAssignment(attEnt AttributableEntity, pAttAss *acmelibv2.AttributeAssignment) error {
	att, ok := l.refAttributes[pAttAss.AttributeEntityId]
	if !ok {
		return &EntityIDError{
			EntityID: EntityID(pAttAss.AttributeEntityId),
			Err:      ErrNotFound,
		}
	}

	switch tmpVal := pAttAss.Value.(type) {
	case *acmelibv2.AttributeAssignment_ValueString:
		if err := attEnt.AssignAttribute(att, tmpVal.ValueString); err != nil {
			log.Print("string: ", tmpVal)
			return err
		}

	case *acmelibv2.AttributeAssignment_ValueInt:
		if err := attEnt.AssignAttribute(att, int(tmpVal.ValueInt)); err != nil {
			return err
		}

	case *acmelibv2.AttributeAssignment_ValueDouble:
		if err := attEnt.AssignAttribute(att, tmpVal.ValueDouble); err != nil {
			log.Print("double: ", tmpVal)
			return err
		}
	}

	return nil
}

// Package migrate implements migration tools.
package migrate

import (
	"log"

	acmelibv1 "github.com/squadracorsepolito/acmelib/gen/acmelib/v1"
	acmelibv2 "github.com/squadracorsepolito/acmelib/gen/acmelib/v2"
)

// FromV1Model migrates a v1 model to a v2 model.
func FromV1Model(v1 *acmelibv1.Network) *acmelibv2.Network {
	m := &v1Migration{}
	return m.migrateNetwork(v1)
}

type v1Migration struct {
}

func (m *v1Migration) migrateEntity(v1 *acmelibv1.Entity) *acmelibv2.Entity {
	return &acmelibv2.Entity{
		EntityId:   v1.EntityId,
		EntityKind: m.migrateEntityKind(v1.EntityKind),
		Name:       v1.Name,
		Desc:       v1.Desc,
		CreateTime: v1.CreateTime,
	}
}

func (m *v1Migration) migrateEntityKind(v1 acmelibv1.EntityKind) acmelibv2.EntityKind {
	switch v1 {
	case acmelibv1.EntityKind_ENTITY_KIND_UNSPECIFIED:
		return acmelibv2.EntityKind_ENTITY_KIND_UNSPECIFIED
	case acmelibv1.EntityKind_ENTITY_KIND_NETWORK:
		return acmelibv2.EntityKind_ENTITY_KIND_NETWORK
	case acmelibv1.EntityKind_ENTITY_KIND_BUS:
		return acmelibv2.EntityKind_ENTITY_KIND_BUS
	case acmelibv1.EntityKind_ENTITY_KIND_NODE:
		return acmelibv2.EntityKind_ENTITY_KIND_NODE
	case acmelibv1.EntityKind_ENTITY_KIND_MESSAGE:
		return acmelibv2.EntityKind_ENTITY_KIND_MESSAGE
	case acmelibv1.EntityKind_ENTITY_KIND_SIGNAL:
		return acmelibv2.EntityKind_ENTITY_KIND_SIGNAL
	case acmelibv1.EntityKind_ENTITY_KIND_SIGNAL_TYPE:
		return acmelibv2.EntityKind_ENTITY_KIND_SIGNAL_TYPE
	case acmelibv1.EntityKind_ENTITY_KIND_SIGNAL_UNIT:
		return acmelibv2.EntityKind_ENTITY_KIND_SIGNAL_UNIT
	case acmelibv1.EntityKind_ENTITY_KIND_SIGNAL_ENUM:
		return acmelibv2.EntityKind_ENTITY_KIND_SIGNAL_ENUM
	case acmelibv1.EntityKind_ENTITY_KIND_CANID_BUILDER:
		return acmelibv2.EntityKind_ENTITY_KIND_CANID_BUILDER
	case acmelibv1.EntityKind_ENTITY_KIND_ATTRIBUTE:
		return acmelibv2.EntityKind_ENTITY_KIND_ATTRIBUTE
	default:
		return acmelibv2.EntityKind_ENTITY_KIND_UNSPECIFIED
	}
}

func (m *v1Migration) migrateAttributeAssignment(v1 *acmelibv1.AttributeAssignment) *acmelibv2.AttributeAssignment {
	v2 := &acmelibv2.AttributeAssignment{
		EntityId:          v1.EntityId,
		AttributeEntityId: v1.AttributeEntityId,
	}

	switch v1Val := v1.Value.(type) {
	case *acmelibv1.AttributeAssignment_ValueString:
		v2.Value = &acmelibv2.AttributeAssignment_ValueString{
			ValueString: v1Val.ValueString,
		}

	case *acmelibv1.AttributeAssignment_ValueInt:
		v2.Value = &acmelibv2.AttributeAssignment_ValueInt{
			ValueInt: v1Val.ValueInt,
		}

	case *acmelibv1.AttributeAssignment_ValueDouble:
		v2.Value = &acmelibv2.AttributeAssignment_ValueDouble{
			ValueDouble: v1Val.ValueDouble,
		}
	}

	return v2
}

func (m *v1Migration) migrateNetwork(v1 *acmelibv1.Network) *acmelibv2.Network {
	v2 := &acmelibv2.Network{
		Entity: m.migrateEntity(v1.Entity),
	}

	for _, v1CANIDBuilder := range v1.CanidBuilders {
		v2.CanidBuilders = append(v2.CanidBuilders, m.migrateCANIDBuilder(v1CANIDBuilder))
	}

	for _, v1Bus := range v1.Buses {
		v2.Buses = append(v2.Buses, m.migrateBus(v1Bus))
	}

	for _, v1Node := range v1.Nodes {
		v2.Nodes = append(v2.Nodes, m.migrateNode(v1Node))
	}

	for _, v1 := range v1.SignalTypes {
		v2.SignalTypes = append(v2.SignalTypes, m.migrateSignalType(v1))
	}

	for _, v1SigUnit := range v1.SignalUnits {
		v2.SignalUnits = append(v2.SignalUnits, m.migrateSignalUnit(v1SigUnit))
	}

	for _, v1SigEnum := range v1.SignalEnums {
		v2.SignalEnums = append(v2.SignalEnums, m.migrateSignalEnum(v1SigEnum))
	}

	for _, v1Att := range v1.Attributes {
		v2.Attributes = append(v2.Attributes, m.migrateAttribute(v1Att))
	}

	return v2
}

func (m *v1Migration) migrateBus(v1 *acmelibv1.Bus) *acmelibv2.Bus {
	v2 := &acmelibv2.Bus{
		Entity: m.migrateEntity(v1.Entity),

		Baudrate:             v1.Baudrate,
		CanidBuilderEntityId: v1.CanidBuilderEntityId,
	}

	for _, v1NodeInterface := range v1.NodeInterfaces {
		v2.NodeInterfaces = append(v2.NodeInterfaces, m.migrateNodeInterface(v1NodeInterface))
	}

	switch v1.Type {
	case acmelibv1.BusType_BUS_TYPE_UNSPECIFIED:
		v2.Type = acmelibv2.BusType_BUS_TYPE_UNSPECIFIED
	case acmelibv1.BusType_BUS_TYPE_CAN_2A:
		v2.Type = acmelibv2.BusType_BUS_TYPE_CAN_2A
	default:
		v2.Type = acmelibv2.BusType_BUS_TYPE_UNSPECIFIED
	}

	for _, v1AttAss := range v1.AttributeAssignments {
		v2.AttributeAssignments = append(v2.AttributeAssignments, m.migrateAttributeAssignment(v1AttAss))
	}

	return v2
}

func (m *v1Migration) migrateNodeInterface(v1 *acmelibv1.NodeInterface) *acmelibv2.NodeInterface {
	v2 := &acmelibv2.NodeInterface{
		Number:       v1.Number,
		NodeEntityId: v1.NodeEntityId,
	}

	for _, v1Msg := range v1.Messages {
		v2.Messages = append(v2.Messages, m.migrateMessage(v1Msg))
	}

	return v2
}

func (m *v1Migration) migrateMessage(v1 *acmelibv1.Message) *acmelibv2.Message {
	v2 := &acmelibv2.Message{
		Entity: m.migrateEntity(v1.Entity),

		SizeByte: v1.SizeByte,

		MessageId:      v1.MessageId,
		StaticCanId:    v1.StaticCanId,
		HasStaticCanId: v1.HasStaticCanId,

		CycleTime:      v1.CycleTime,
		DelayTime:      v1.DelayTime,
		StartDelayTime: v1.StartDelayTime,
	}

	switch v1.Priority {
	case acmelibv1.MessagePriority_MESSAGE_PRIORITY_UNSPECIFIED:
		v2.Priority = acmelibv2.MessagePriority_MESSAGE_PRIORITY_UNSPECIFIED
	case acmelibv1.MessagePriority_MESSAGE_PRIORITY_VERY_HIGH:
		v2.Priority = acmelibv2.MessagePriority_MESSAGE_PRIORITY_VERY_HIGH
	case acmelibv1.MessagePriority_MESSAGE_PRIORITY_HIGH:
		v2.Priority = acmelibv2.MessagePriority_MESSAGE_PRIORITY_HIGH
	case acmelibv1.MessagePriority_MESSAGE_PRIORITY_MEDIUM:
		v2.Priority = acmelibv2.MessagePriority_MESSAGE_PRIORITY_MEDIUM
	case acmelibv1.MessagePriority_MESSAGE_PRIORITY_LOW:
		v2.Priority = acmelibv2.MessagePriority_MESSAGE_PRIORITY_LOW
	}

	switch v1.SendType {
	case acmelibv1.MessageSendType_MESSAGE_SEND_TYPE_UNSPECIFIED:
		v2.SendType = acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_UNSPECIFIED
	case acmelibv1.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC:
		v2.SendType = acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC
	case acmelibv1.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE:
		v2.SendType = acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE
	case acmelibv1.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_AND_TRIGGERED:
		v2.SendType = acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_AND_TRIGGERED
	case acmelibv1.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE_AND_TRIGGERED:
		v2.SendType = acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE_AND_TRIGGERED
	default:
		v2.SendType = acmelibv2.MessageSendType_MESSAGE_SEND_TYPE_UNSPECIFIED
	}

	v2.Layout = &acmelibv2.SignalLayout{
		Id:       0,
		SizeByte: v1.SizeByte,
	}

	endianness := acmelibv2.Endianness_ENDIANNESS_UNSPECIFIED
	switch v1.ByteOrder {
	case acmelibv1.MessageByteOrder_MESSAGE_BYTE_ORDER_UNSPECIFIED:
		endianness = acmelibv2.Endianness_ENDIANNESS_UNSPECIFIED
	case acmelibv1.MessageByteOrder_MESSAGE_BYTE_ORDER_LITTLE_ENDIAN:
		endianness = acmelibv2.Endianness_ENDIANNESS_LITTLE_ENDIAN
	case acmelibv1.MessageByteOrder_MESSAGE_BYTE_ORDER_BIG_ENDIAN:
		endianness = acmelibv2.Endianness_ENDIANNESS_BIG_ENDIAN
	}

	sigStartPos := make(map[string]uint32)
	for _, v1Ref := range v1.Payload.Refs {
		sigStartPos[v1Ref.SignalEntityId] = uint32(v1Ref.RelStartBit)
	}

	for _, v1Sig := range v1.Signals {
		if v1Sig.Kind == acmelibv1.SignalKind_SIGNAL_KIND_MULTIPLEXER {
			log.Printf("skipping multiplexer signal %s", v1Sig.Entity.Name)
			continue
		}

		sig := m.migrateSignal(v1Sig)

		sig.Endianness = endianness
		sig.StartPos = sigStartPos[v1Sig.Entity.EntityId]

		v2.Layout.Signals = append(v2.Layout.Signals, sig)
	}

	for _, v1AttAss := range v1.AttributeAssignments {
		v2.AttributeAssignments = append(v2.AttributeAssignments, m.migrateAttributeAssignment(v1AttAss))
	}

	return v2
}

func (m *v1Migration) migrateSignal(v1 *acmelibv1.Signal) *acmelibv2.Signal {
	v2 := &acmelibv2.Signal{
		Entity:     m.migrateEntity(v1.Entity),
		StartValue: v1.StartValue,
	}

	switch v1.Kind {
	case acmelibv1.SignalKind_SIGNAL_KIND_UNSPECIFIED:
		v2.Kind = acmelibv2.SignalKind_SIGNAL_KIND_UNSPECIFIED
	case acmelibv1.SignalKind_SIGNAL_KIND_STANDARD:
		v2.Kind = acmelibv2.SignalKind_SIGNAL_KIND_STANDARD
	case acmelibv1.SignalKind_SIGNAL_KIND_ENUM:
		v2.Kind = acmelibv2.SignalKind_SIGNAL_KIND_ENUM
	default:
		v2.Kind = acmelibv2.SignalKind_SIGNAL_KIND_UNSPECIFIED
	}

	switch v1.SendType {
	case acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_UNSPECIFIED:
		v2.SendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_UNSPECIFIED
	case acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_CYCLIC:
		v2.SendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_CYCLIC
	case acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_ON_WRITE:
		v2.SendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_WRITE
	case acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_ON_WRITE_WITH_REPETITION:
		v2.SendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_WRITE_WITH_REPETITION
	case acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_ON_CHANGE:
		v2.SendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_CHANGE
	case acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_ON_CHANGE_WITH_REPETITION:
		v2.SendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_ON_CHANGE_WITH_REPETITION
	case acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_IF_ACTIVE:
		v2.SendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_IF_ACTIVE
	case acmelibv1.SignalSendType_SIGNAL_SEND_TYPE_IF_ACTIVE_WITH_REPETITION:
		v2.SendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_IF_ACTIVE_WITH_REPETITION
	default:
		v2.SendType = acmelibv2.SignalSendType_SIGNAL_SEND_TYPE_UNSPECIFIED
	}

	switch v1Sig := v1.Signal.(type) {
	case *acmelibv1.Signal_Standard:
		v2.Signal = &acmelibv2.Signal_Standard{
			Standard: m.migrateStandardSignal(v1Sig.Standard),
		}

	case *acmelibv1.Signal_Enum:
		v2.Signal = &acmelibv2.Signal_Enum{
			Enum: m.migrateEnumSignal(v1Sig.Enum),
		}
	}

	for _, v1AttAss := range v1.AttributeAssignments {
		v2.AttributeAssignments = append(v2.AttributeAssignments, m.migrateAttributeAssignment(v1AttAss))
	}

	return v2
}

func (m *v1Migration) migrateStandardSignal(v1 *acmelibv1.StandardSignal) *acmelibv2.StandardSignal {
	return &acmelibv2.StandardSignal{
		TypeEntityId: v1.TypeEntityId,
		UnitEntityId: v1.UnitEntityId,
	}
}

func (m *v1Migration) migrateEnumSignal(v1 *acmelibv1.EnumSignal) *acmelibv2.EnumSignal {
	return &acmelibv2.EnumSignal{
		EnumEntityId: v1.EnumEntityId,
	}
}

func (m *v1Migration) migrateNode(v1 *acmelibv1.Node) *acmelibv2.Node {
	v2 := &acmelibv2.Node{
		Entity: m.migrateEntity(v1.Entity),

		NodeId:         v1.NodeId,
		InterfaceCount: v1.InterfaceCount,
	}

	for _, v1AttAss := range v1.AttributeAssignments {
		v2.AttributeAssignments = append(v2.AttributeAssignments, m.migrateAttributeAssignment(v1AttAss))
	}

	return v2
}

func (m *v1Migration) migrateCANIDBuilder(v1 *acmelibv1.CANIDBuilder) *acmelibv2.CANIDBuilder {
	v2 := &acmelibv2.CANIDBuilder{
		Entity: m.migrateEntity(v1.Entity),
	}

	for _, v1Op := range v1.Operations {
		v2.Operations = append(v2.Operations, m.migrateCANIDBuilderOp(v1Op))
	}

	return v2
}

func (m *v1Migration) migrateCANIDBuilderOp(v1 *acmelibv1.CANIDBuilderOp) *acmelibv2.CANIDBuilderOp {
	v2 := &acmelibv2.CANIDBuilderOp{
		From: v1.From,
		Len:  v1.Len,
	}

	switch v1.Kind {
	case acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_UNSPECIFIED:
		v2.Kind = acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_UNSPECIFIED
	case acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_PRIORITY:
		v2.Kind = acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_PRIORITY
	case acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_ID:
		v2.Kind = acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_ID
	case acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_NODE_ID:
		v2.Kind = acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_NODE_ID
	case acmelibv1.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_BIT_MASK:
		v2.Kind = acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_BIT_MASK
	default:
		v2.Kind = acmelibv2.CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_UNSPECIFIED
	}

	return v2
}

func (m *v1Migration) migrateSignalType(v1 *acmelibv1.SignalType) *acmelibv2.SignalType {
	v2 := &acmelibv2.SignalType{
		Entity: m.migrateEntity(v1.Entity),

		Size:   v1.Size,
		Signed: v1.Signed,
		Min:    v1.Min,
		Max:    v1.Max,
		Scale:  v1.Scale,
		Offset: v1.Offset,
	}

	switch v1.Kind {
	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_UNSPECIFIED:
		v2.Kind = acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_UNSPECIFIED
	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_CUSTOM:
		v2.Kind = acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_DECIMAL
	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_FLAG:
		v2.Kind = acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_FLAG
	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_INTEGER:
		v2.Kind = acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_INTEGER
	case acmelibv1.SignalTypeKind_SIGNAL_TYPE_KIND_DECIMAL:
		v2.Kind = acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_DECIMAL
	default:
		v2.Kind = acmelibv2.SignalTypeKind_SIGNAL_TYPE_KIND_UNSPECIFIED
	}

	return v2
}

func (m *v1Migration) migrateSignalUnit(v1 *acmelibv1.SignalUnit) *acmelibv2.SignalUnit {
	v2 := &acmelibv2.SignalUnit{
		Entity: m.migrateEntity(v1.Entity),
		Symbol: v1.Symbol,
	}

	switch v1.Kind {
	case acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_UNSPECIFIED:
		v2.Kind = acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_UNSPECIFIED
	case acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_CUSTOM:
		v2.Kind = acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_CUSTOM
	case acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_TEMPERATURE:
		v2.Kind = acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_TEMPERATURE
	case acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_ELECTRICAL:
		v2.Kind = acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_ELECTRICAL
	case acmelibv1.SignalUnitKind_SIGNAL_UNIT_KIND_POWER:
		v2.Kind = acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_POWER
	default:
		v2.Kind = acmelibv2.SignalUnitKind_SIGNAL_UNIT_KIND_UNSPECIFIED
	}

	return v2
}

func (m *v1Migration) migrateSignalEnum(v1 *acmelibv1.SignalEnum) *acmelibv2.SignalEnum {
	v2 := &acmelibv2.SignalEnum{
		Entity: m.migrateEntity(v1.Entity),
	}

	if v1.MinSize != 0 {
		v2.Size = v1.MinSize
		v2.FixedSize = true
	}

	for _, v1Val := range v1.Values {
		v2.Values = append(v2.Values, m.migrateSignalEnumValue(v1Val))
	}

	return v2
}

func (m *v1Migration) migrateSignalEnumValue(v1 *acmelibv1.SignalEnumValue) *acmelibv2.SignalEnumValue {
	return &acmelibv2.SignalEnumValue{
		Index: v1.Index,
		Name:  v1.Entity.Name,
		Desc:  v1.Entity.Desc,
	}
}

func (m *v1Migration) migrateAttribute(v1 *acmelibv1.Attribute) *acmelibv2.Attribute {
	v2 := &acmelibv2.Attribute{
		Entity: m.migrateEntity(v1.Entity),
	}

	switch v1.Type {
	case acmelibv1.AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED:
		v2.Type = acmelibv2.AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED
	case acmelibv1.AttributeType_ATTRIBUTE_TYPE_STRING:
		v2.Type = acmelibv2.AttributeType_ATTRIBUTE_TYPE_STRING
	case acmelibv1.AttributeType_ATTRIBUTE_TYPE_INTEGER:
		v2.Type = acmelibv2.AttributeType_ATTRIBUTE_TYPE_INTEGER
	case acmelibv1.AttributeType_ATTRIBUTE_TYPE_FLOAT:
		v2.Type = acmelibv2.AttributeType_ATTRIBUTE_TYPE_FLOAT
	case acmelibv1.AttributeType_ATTRIBUTE_TYPE_ENUM:
		v2.Type = acmelibv2.AttributeType_ATTRIBUTE_TYPE_ENUM
	default:
		v2.Type = acmelibv2.AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED
	}

	switch v1Att := v1.Attribute.(type) {
	case *acmelibv1.Attribute_StringAttribute:
		v2.Attribute = &acmelibv2.Attribute_StringAttribute{
			StringAttribute: &acmelibv2.StringAttribute{
				DefValue: v1Att.StringAttribute.DefValue,
			},
		}

	case *acmelibv1.Attribute_IntegerAttribute:
		v2.Attribute = &acmelibv2.Attribute_IntegerAttribute{
			IntegerAttribute: &acmelibv2.IntegerAttribute{
				DefValue:    v1Att.IntegerAttribute.DefValue,
				Min:         v1Att.IntegerAttribute.Min,
				Max:         v1Att.IntegerAttribute.Max,
				IsHexFormat: v1Att.IntegerAttribute.IsHexFormat,
			},
		}

	case *acmelibv1.Attribute_FloatAttribute:
		v2.Attribute = &acmelibv2.Attribute_FloatAttribute{
			FloatAttribute: &acmelibv2.FloatAttribute{
				DefValue: v1Att.FloatAttribute.DefValue,
				Min:      v1Att.FloatAttribute.Min,
				Max:      v1Att.FloatAttribute.Max,
			},
		}

	case *acmelibv1.Attribute_EnumAttribute:
		v2.Attribute = &acmelibv2.Attribute_EnumAttribute{
			EnumAttribute: &acmelibv2.EnumAttribute{
				DefValue: v1Att.EnumAttribute.DefValue,
				Values:   v1Att.EnumAttribute.Values,
			},
		}
	}

	return v2
}

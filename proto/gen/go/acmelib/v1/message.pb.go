// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: acmelib/v1/message.proto

package acmelibv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MessagePriority int32

const (
	MessagePriority_MESSAGE_PRIORITY_UNSPECIFIED MessagePriority = 0
	MessagePriority_MESSAGE_PRIORITY_VERY_HIGH   MessagePriority = 1
	MessagePriority_MESSAGE_PRIORITY_HIGH        MessagePriority = 2
	MessagePriority_MESSAGE_PRIORITY_MEDIUM      MessagePriority = 3
	MessagePriority_MESSAGE_PRIORITY_LOW         MessagePriority = 4
)

// Enum value maps for MessagePriority.
var (
	MessagePriority_name = map[int32]string{
		0: "MESSAGE_PRIORITY_UNSPECIFIED",
		1: "MESSAGE_PRIORITY_VERY_HIGH",
		2: "MESSAGE_PRIORITY_HIGH",
		3: "MESSAGE_PRIORITY_MEDIUM",
		4: "MESSAGE_PRIORITY_LOW",
	}
	MessagePriority_value = map[string]int32{
		"MESSAGE_PRIORITY_UNSPECIFIED": 0,
		"MESSAGE_PRIORITY_VERY_HIGH":   1,
		"MESSAGE_PRIORITY_HIGH":        2,
		"MESSAGE_PRIORITY_MEDIUM":      3,
		"MESSAGE_PRIORITY_LOW":         4,
	}
)

func (x MessagePriority) Enum() *MessagePriority {
	p := new(MessagePriority)
	*p = x
	return p
}

func (x MessagePriority) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessagePriority) Descriptor() protoreflect.EnumDescriptor {
	return file_acmelib_v1_message_proto_enumTypes[0].Descriptor()
}

func (MessagePriority) Type() protoreflect.EnumType {
	return &file_acmelib_v1_message_proto_enumTypes[0]
}

func (x MessagePriority) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessagePriority.Descriptor instead.
func (MessagePriority) EnumDescriptor() ([]byte, []int) {
	return file_acmelib_v1_message_proto_rawDescGZIP(), []int{0}
}

type MessageSendType int32

const (
	MessageSendType_MESSAGE_SEND_TYPE_UNSPECIFIED                    MessageSendType = 0
	MessageSendType_MESSAGE_SEND_TYPE_CYCLIC                         MessageSendType = 1
	MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE               MessageSendType = 2
	MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_AND_TRIGGERED           MessageSendType = 3
	MessageSendType_MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE_AND_TRIGGERED MessageSendType = 4
)

// Enum value maps for MessageSendType.
var (
	MessageSendType_name = map[int32]string{
		0: "MESSAGE_SEND_TYPE_UNSPECIFIED",
		1: "MESSAGE_SEND_TYPE_CYCLIC",
		2: "MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE",
		3: "MESSAGE_SEND_TYPE_CYCLIC_AND_TRIGGERED",
		4: "MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE_AND_TRIGGERED",
	}
	MessageSendType_value = map[string]int32{
		"MESSAGE_SEND_TYPE_UNSPECIFIED":                    0,
		"MESSAGE_SEND_TYPE_CYCLIC":                         1,
		"MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE":               2,
		"MESSAGE_SEND_TYPE_CYCLIC_AND_TRIGGERED":           3,
		"MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE_AND_TRIGGERED": 4,
	}
)

func (x MessageSendType) Enum() *MessageSendType {
	p := new(MessageSendType)
	*p = x
	return p
}

func (x MessageSendType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageSendType) Descriptor() protoreflect.EnumDescriptor {
	return file_acmelib_v1_message_proto_enumTypes[1].Descriptor()
}

func (MessageSendType) Type() protoreflect.EnumType {
	return &file_acmelib_v1_message_proto_enumTypes[1]
}

func (x MessageSendType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageSendType.Descriptor instead.
func (MessageSendType) EnumDescriptor() ([]byte, []int) {
	return file_acmelib_v1_message_proto_rawDescGZIP(), []int{1}
}

type MessageByteOrder int32

const (
	MessageByteOrder_MESSAGE_BYTE_ORDER_UNSPECIFIED   MessageByteOrder = 0
	MessageByteOrder_MESSAGE_BYTE_ORDER_LITTLE_ENDIAN MessageByteOrder = 1
	MessageByteOrder_MESSAGE_BYTE_ORDER_BIG_ENDIAN    MessageByteOrder = 2
)

// Enum value maps for MessageByteOrder.
var (
	MessageByteOrder_name = map[int32]string{
		0: "MESSAGE_BYTE_ORDER_UNSPECIFIED",
		1: "MESSAGE_BYTE_ORDER_LITTLE_ENDIAN",
		2: "MESSAGE_BYTE_ORDER_BIG_ENDIAN",
	}
	MessageByteOrder_value = map[string]int32{
		"MESSAGE_BYTE_ORDER_UNSPECIFIED":   0,
		"MESSAGE_BYTE_ORDER_LITTLE_ENDIAN": 1,
		"MESSAGE_BYTE_ORDER_BIG_ENDIAN":    2,
	}
)

func (x MessageByteOrder) Enum() *MessageByteOrder {
	p := new(MessageByteOrder)
	*p = x
	return p
}

func (x MessageByteOrder) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageByteOrder) Descriptor() protoreflect.EnumDescriptor {
	return file_acmelib_v1_message_proto_enumTypes[2].Descriptor()
}

func (MessageByteOrder) Type() protoreflect.EnumType {
	return &file_acmelib_v1_message_proto_enumTypes[2]
}

func (x MessageByteOrder) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageByteOrder.Descriptor instead.
func (MessageByteOrder) EnumDescriptor() ([]byte, []int) {
	return file_acmelib_v1_message_proto_rawDescGZIP(), []int{2}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entity               *Entity                `protobuf:"bytes,1,opt,name=entity,proto3" json:"entity,omitempty"`
	Signals              []*Signal              `protobuf:"bytes,2,rep,name=signals,proto3" json:"signals,omitempty"`
	Payload              *SignalPayload         `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	SizeByte             uint32                 `protobuf:"varint,4,opt,name=size_byte,json=sizeByte,proto3" json:"size_byte,omitempty"`
	MessageId            uint32                 `protobuf:"varint,5,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	StaticCanId          uint32                 `protobuf:"varint,6,opt,name=static_can_id,json=staticCanId,proto3" json:"static_can_id,omitempty"`
	HasStaticCanId       bool                   `protobuf:"varint,7,opt,name=has_static_can_id,json=hasStaticCanId,proto3" json:"has_static_can_id,omitempty"`
	Priority             MessagePriority        `protobuf:"varint,8,opt,name=priority,proto3,enum=acmelib.v1.MessagePriority" json:"priority,omitempty"`
	ByteOrder            MessageByteOrder       `protobuf:"varint,9,opt,name=byte_order,json=byteOrder,proto3,enum=acmelib.v1.MessageByteOrder" json:"byte_order,omitempty"`
	CycleTime            uint32                 `protobuf:"varint,10,opt,name=cycle_time,json=cycleTime,proto3" json:"cycle_time,omitempty"`
	SendType             MessageSendType        `protobuf:"varint,11,opt,name=send_type,json=sendType,proto3,enum=acmelib.v1.MessageSendType" json:"send_type,omitempty"`
	DelayTime            uint32                 `protobuf:"varint,12,opt,name=delay_time,json=delayTime,proto3" json:"delay_time,omitempty"`
	StartDelayTime       uint32                 `protobuf:"varint,13,opt,name=start_delay_time,json=startDelayTime,proto3" json:"start_delay_time,omitempty"`
	Receivers            []*MessageReceiver     `protobuf:"bytes,14,rep,name=receivers,proto3" json:"receivers,omitempty"`
	AttributeAssignments []*AttributeAssignment `protobuf:"bytes,15,rep,name=attribute_assignments,json=attributeAssignments,proto3" json:"attribute_assignments,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acmelib_v1_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v1_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_acmelib_v1_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetEntity() *Entity {
	if x != nil {
		return x.Entity
	}
	return nil
}

func (x *Message) GetSignals() []*Signal {
	if x != nil {
		return x.Signals
	}
	return nil
}

func (x *Message) GetPayload() *SignalPayload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Message) GetSizeByte() uint32 {
	if x != nil {
		return x.SizeByte
	}
	return 0
}

func (x *Message) GetMessageId() uint32 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

func (x *Message) GetStaticCanId() uint32 {
	if x != nil {
		return x.StaticCanId
	}
	return 0
}

func (x *Message) GetHasStaticCanId() bool {
	if x != nil {
		return x.HasStaticCanId
	}
	return false
}

func (x *Message) GetPriority() MessagePriority {
	if x != nil {
		return x.Priority
	}
	return MessagePriority_MESSAGE_PRIORITY_UNSPECIFIED
}

func (x *Message) GetByteOrder() MessageByteOrder {
	if x != nil {
		return x.ByteOrder
	}
	return MessageByteOrder_MESSAGE_BYTE_ORDER_UNSPECIFIED
}

func (x *Message) GetCycleTime() uint32 {
	if x != nil {
		return x.CycleTime
	}
	return 0
}

func (x *Message) GetSendType() MessageSendType {
	if x != nil {
		return x.SendType
	}
	return MessageSendType_MESSAGE_SEND_TYPE_UNSPECIFIED
}

func (x *Message) GetDelayTime() uint32 {
	if x != nil {
		return x.DelayTime
	}
	return 0
}

func (x *Message) GetStartDelayTime() uint32 {
	if x != nil {
		return x.StartDelayTime
	}
	return 0
}

func (x *Message) GetReceivers() []*MessageReceiver {
	if x != nil {
		return x.Receivers
	}
	return nil
}

func (x *Message) GetAttributeAssignments() []*AttributeAssignment {
	if x != nil {
		return x.AttributeAssignments
	}
	return nil
}

type MessageReceiver struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeEntityId        string `protobuf:"bytes,1,opt,name=node_entity_id,json=nodeEntityId,proto3" json:"node_entity_id,omitempty"`
	NodeInterfaceNumber uint32 `protobuf:"varint,2,opt,name=node_interface_number,json=nodeInterfaceNumber,proto3" json:"node_interface_number,omitempty"`
}

func (x *MessageReceiver) Reset() {
	*x = MessageReceiver{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acmelib_v1_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageReceiver) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageReceiver) ProtoMessage() {}

func (x *MessageReceiver) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v1_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageReceiver.ProtoReflect.Descriptor instead.
func (*MessageReceiver) Descriptor() ([]byte, []int) {
	return file_acmelib_v1_message_proto_rawDescGZIP(), []int{1}
}

func (x *MessageReceiver) GetNodeEntityId() string {
	if x != nil {
		return x.NodeEntityId
	}
	return ""
}

func (x *MessageReceiver) GetNodeInterfaceNumber() uint32 {
	if x != nil {
		return x.NodeInterfaceNumber
	}
	return 0
}

var File_acmelib_v1_message_proto protoreflect.FileDescriptor

var file_acmelib_v1_message_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x61, 0x63, 0x6d, 0x65,
	0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x17, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2f,
	0x76, 0x31, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x17, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69,
	0x62, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcc, 0x05, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x2a, 0x0a, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x52, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x2c, 0x0a, 0x07,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61,
	0x6c, 0x52, 0x07, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x73, 0x12, 0x33, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x63,
	0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x69, 0x7a, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x08, 0x73, 0x69, 0x7a, 0x65, 0x42, 0x79, 0x74, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x63, 0x61, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x43, 0x61, 0x6e, 0x49, 0x64, 0x12,
	0x29, 0x0a, 0x11, 0x68, 0x61, 0x73, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x63, 0x61,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x68, 0x61, 0x73, 0x53,
	0x74, 0x61, 0x74, 0x69, 0x63, 0x43, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x37, 0x0a, 0x08, 0x70, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x61,
	0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72,
	0x69, 0x74, 0x79, 0x12, 0x3b, 0x0a, 0x0a, 0x62, 0x79, 0x74, 0x65, 0x5f, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69,
	0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x79, 0x74, 0x65,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x09, 0x62, 0x79, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x38, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x08, 0x73, 0x65, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x65, 0x6c,
	0x61, 0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x64,
	0x65, 0x6c, 0x61, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x10, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0e, 0x73, 0x74, 0x61, 0x72, 0x74, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x39, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x73, 0x18,
	0x0e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e,
	0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x72, 0x52, 0x09, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x73, 0x12, 0x54, 0x0a,
	0x15, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x5f, 0x61, 0x73, 0x73, 0x69, 0x67,
	0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x61,
	0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x14, 0x61,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x22, 0x6b, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x0e, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x6e, 0x6f, 0x64, 0x65, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x32, 0x0a, 0x15,
	0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x5f, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x13, 0x6e, 0x6f, 0x64,
	0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x2a, 0xa5, 0x01, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x72, 0x69, 0x6f,
	0x72, 0x69, 0x74, 0x79, 0x12, 0x20, 0x0a, 0x1c, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f,
	0x50, 0x52, 0x49, 0x4f, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1e, 0x0a, 0x1a, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47,
	0x45, 0x5f, 0x50, 0x52, 0x49, 0x4f, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x56, 0x45, 0x52, 0x59, 0x5f,
	0x48, 0x49, 0x47, 0x48, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47,
	0x45, 0x5f, 0x50, 0x52, 0x49, 0x4f, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x48, 0x49, 0x47, 0x48, 0x10,
	0x02, 0x12, 0x1b, 0x0a, 0x17, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x50, 0x52, 0x49,
	0x4f, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4d, 0x45, 0x44, 0x49, 0x55, 0x4d, 0x10, 0x03, 0x12, 0x18,
	0x0a, 0x14, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x50, 0x52, 0x49, 0x4f, 0x52, 0x49,
	0x54, 0x59, 0x5f, 0x4c, 0x4f, 0x57, 0x10, 0x04, 0x2a, 0xdc, 0x01, 0x0a, 0x0f, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x1d,
	0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x53, 0x45, 0x4e, 0x44, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x1c, 0x0a, 0x18, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x53, 0x45, 0x4e, 0x44, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x59, 0x43, 0x4c, 0x49, 0x43, 0x10, 0x01, 0x12, 0x26, 0x0a,
	0x22, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x53, 0x45, 0x4e, 0x44, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x43, 0x59, 0x43, 0x4c, 0x49, 0x43, 0x5f, 0x49, 0x46, 0x5f, 0x41, 0x43, 0x54,
	0x49, 0x56, 0x45, 0x10, 0x02, 0x12, 0x2a, 0x0a, 0x26, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45,
	0x5f, 0x53, 0x45, 0x4e, 0x44, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x59, 0x43, 0x4c, 0x49,
	0x43, 0x5f, 0x41, 0x4e, 0x44, 0x5f, 0x54, 0x52, 0x49, 0x47, 0x47, 0x45, 0x52, 0x45, 0x44, 0x10,
	0x03, 0x12, 0x34, 0x0a, 0x30, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x53, 0x45, 0x4e,
	0x44, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x59, 0x43, 0x4c, 0x49, 0x43, 0x5f, 0x49, 0x46,
	0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x5f, 0x41, 0x4e, 0x44, 0x5f, 0x54, 0x52, 0x49, 0x47,
	0x47, 0x45, 0x52, 0x45, 0x44, 0x10, 0x04, 0x2a, 0x7f, 0x0a, 0x10, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x42, 0x79, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x1e, 0x4d,
	0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x42, 0x59, 0x54, 0x45, 0x5f, 0x4f, 0x52, 0x44, 0x45,
	0x52, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x24, 0x0a, 0x20, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x42, 0x59, 0x54, 0x45, 0x5f,
	0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x4c, 0x49, 0x54, 0x54, 0x4c, 0x45, 0x5f, 0x45, 0x4e, 0x44,
	0x49, 0x41, 0x4e, 0x10, 0x01, 0x12, 0x21, 0x0a, 0x1d, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45,
	0x5f, 0x42, 0x59, 0x54, 0x45, 0x5f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x42, 0x49, 0x47, 0x5f,
	0x45, 0x4e, 0x44, 0x49, 0x41, 0x4e, 0x10, 0x02, 0x42, 0x7d, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e,
	0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x14, 0x61, 0x63, 0x6d, 0x65,
	0x6c, 0x69, 0x62, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x76, 0x31,
	0xa2, 0x02, 0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x0a, 0x41, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x0a, 0x41, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x16, 0x41, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x41, 0x63, 0x6d, 0x65,
	0x6c, 0x69, 0x62, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_acmelib_v1_message_proto_rawDescOnce sync.Once
	file_acmelib_v1_message_proto_rawDescData = file_acmelib_v1_message_proto_rawDesc
)

func file_acmelib_v1_message_proto_rawDescGZIP() []byte {
	file_acmelib_v1_message_proto_rawDescOnce.Do(func() {
		file_acmelib_v1_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_acmelib_v1_message_proto_rawDescData)
	})
	return file_acmelib_v1_message_proto_rawDescData
}

var file_acmelib_v1_message_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_acmelib_v1_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_acmelib_v1_message_proto_goTypes = []interface{}{
	(MessagePriority)(0),        // 0: acmelib.v1.MessagePriority
	(MessageSendType)(0),        // 1: acmelib.v1.MessageSendType
	(MessageByteOrder)(0),       // 2: acmelib.v1.MessageByteOrder
	(*Message)(nil),             // 3: acmelib.v1.Message
	(*MessageReceiver)(nil),     // 4: acmelib.v1.MessageReceiver
	(*Entity)(nil),              // 5: acmelib.v1.Entity
	(*Signal)(nil),              // 6: acmelib.v1.Signal
	(*SignalPayload)(nil),       // 7: acmelib.v1.SignalPayload
	(*AttributeAssignment)(nil), // 8: acmelib.v1.AttributeAssignment
}
var file_acmelib_v1_message_proto_depIdxs = []int32{
	5, // 0: acmelib.v1.Message.entity:type_name -> acmelib.v1.Entity
	6, // 1: acmelib.v1.Message.signals:type_name -> acmelib.v1.Signal
	7, // 2: acmelib.v1.Message.payload:type_name -> acmelib.v1.SignalPayload
	0, // 3: acmelib.v1.Message.priority:type_name -> acmelib.v1.MessagePriority
	2, // 4: acmelib.v1.Message.byte_order:type_name -> acmelib.v1.MessageByteOrder
	1, // 5: acmelib.v1.Message.send_type:type_name -> acmelib.v1.MessageSendType
	4, // 6: acmelib.v1.Message.receivers:type_name -> acmelib.v1.MessageReceiver
	8, // 7: acmelib.v1.Message.attribute_assignments:type_name -> acmelib.v1.AttributeAssignment
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_acmelib_v1_message_proto_init() }
func file_acmelib_v1_message_proto_init() {
	if File_acmelib_v1_message_proto != nil {
		return
	}
	file_acmelib_v1_entity_proto_init()
	file_acmelib_v1_signal_proto_init()
	file_acmelib_v1_attribute_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_acmelib_v1_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_acmelib_v1_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageReceiver); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_acmelib_v1_message_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_acmelib_v1_message_proto_goTypes,
		DependencyIndexes: file_acmelib_v1_message_proto_depIdxs,
		EnumInfos:         file_acmelib_v1_message_proto_enumTypes,
		MessageInfos:      file_acmelib_v1_message_proto_msgTypes,
	}.Build()
	File_acmelib_v1_message_proto = out.File
	file_acmelib_v1_message_proto_rawDesc = nil
	file_acmelib_v1_message_proto_goTypes = nil
	file_acmelib_v1_message_proto_depIdxs = nil
}

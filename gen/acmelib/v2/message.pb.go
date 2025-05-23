// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: acmelib/v2/message.proto

package acmelibv2

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
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
	return file_acmelib_v2_message_proto_enumTypes[0].Descriptor()
}

func (MessagePriority) Type() protoreflect.EnumType {
	return &file_acmelib_v2_message_proto_enumTypes[0]
}

func (x MessagePriority) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessagePriority.Descriptor instead.
func (MessagePriority) EnumDescriptor() ([]byte, []int) {
	return file_acmelib_v2_message_proto_rawDescGZIP(), []int{0}
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
	return file_acmelib_v2_message_proto_enumTypes[1].Descriptor()
}

func (MessageSendType) Type() protoreflect.EnumType {
	return &file_acmelib_v2_message_proto_enumTypes[1]
}

func (x MessageSendType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageSendType.Descriptor instead.
func (MessageSendType) EnumDescriptor() ([]byte, []int) {
	return file_acmelib_v2_message_proto_rawDescGZIP(), []int{1}
}

type Message struct {
	state                protoimpl.MessageState `protogen:"open.v1"`
	Entity               *Entity                `protobuf:"bytes,1,opt,name=entity,proto3" json:"entity,omitempty"`
	Layout               *SignalLayout          `protobuf:"bytes,2,opt,name=layout,proto3" json:"layout,omitempty"`
	SizeByte             uint32                 `protobuf:"varint,3,opt,name=size_byte,json=sizeByte,proto3" json:"size_byte,omitempty"`
	MessageId            uint32                 `protobuf:"varint,4,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	StaticCanId          uint32                 `protobuf:"varint,5,opt,name=static_can_id,json=staticCanId,proto3" json:"static_can_id,omitempty"`
	HasStaticCanId       bool                   `protobuf:"varint,6,opt,name=has_static_can_id,json=hasStaticCanId,proto3" json:"has_static_can_id,omitempty"`
	Priority             MessagePriority        `protobuf:"varint,7,opt,name=priority,proto3,enum=acmelib.v2.MessagePriority" json:"priority,omitempty"`
	CycleTime            uint32                 `protobuf:"varint,8,opt,name=cycle_time,json=cycleTime,proto3" json:"cycle_time,omitempty"`
	SendType             MessageSendType        `protobuf:"varint,9,opt,name=send_type,json=sendType,proto3,enum=acmelib.v2.MessageSendType" json:"send_type,omitempty"`
	DelayTime            uint32                 `protobuf:"varint,10,opt,name=delay_time,json=delayTime,proto3" json:"delay_time,omitempty"`
	StartDelayTime       uint32                 `protobuf:"varint,11,opt,name=start_delay_time,json=startDelayTime,proto3" json:"start_delay_time,omitempty"`
	Receivers            []*MessageReceiver     `protobuf:"bytes,12,rep,name=receivers,proto3" json:"receivers,omitempty"`
	AttributeAssignments []*AttributeAssignment `protobuf:"bytes,13,rep,name=attribute_assignments,json=attributeAssignments,proto3" json:"attribute_assignments,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_acmelib_v2_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v2_message_proto_msgTypes[0]
	if x != nil {
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
	return file_acmelib_v2_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetEntity() *Entity {
	if x != nil {
		return x.Entity
	}
	return nil
}

func (x *Message) GetLayout() *SignalLayout {
	if x != nil {
		return x.Layout
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
	state               protoimpl.MessageState `protogen:"open.v1"`
	NodeEntityId        string                 `protobuf:"bytes,1,opt,name=node_entity_id,json=nodeEntityId,proto3" json:"node_entity_id,omitempty"`
	NodeInterfaceNumber uint32                 `protobuf:"varint,2,opt,name=node_interface_number,json=nodeInterfaceNumber,proto3" json:"node_interface_number,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *MessageReceiver) Reset() {
	*x = MessageReceiver{}
	mi := &file_acmelib_v2_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageReceiver) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageReceiver) ProtoMessage() {}

func (x *MessageReceiver) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v2_message_proto_msgTypes[1]
	if x != nil {
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
	return file_acmelib_v2_message_proto_rawDescGZIP(), []int{1}
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

var File_acmelib_v2_message_proto protoreflect.FileDescriptor

const file_acmelib_v2_message_proto_rawDesc = "" +
	"\n" +
	"\x18acmelib/v2/message.proto\x12\n" +
	"acmelib.v2\x1a\x17acmelib/v2/entity.proto\x1a\x17acmelib/v2/signal.proto\x1a\x1aacmelib/v2/attribute.proto\"\xde\x04\n" +
	"\aMessage\x12*\n" +
	"\x06entity\x18\x01 \x01(\v2\x12.acmelib.v2.EntityR\x06entity\x120\n" +
	"\x06layout\x18\x02 \x01(\v2\x18.acmelib.v2.SignalLayoutR\x06layout\x12\x1b\n" +
	"\tsize_byte\x18\x03 \x01(\rR\bsizeByte\x12\x1d\n" +
	"\n" +
	"message_id\x18\x04 \x01(\rR\tmessageId\x12\"\n" +
	"\rstatic_can_id\x18\x05 \x01(\rR\vstaticCanId\x12)\n" +
	"\x11has_static_can_id\x18\x06 \x01(\bR\x0ehasStaticCanId\x127\n" +
	"\bpriority\x18\a \x01(\x0e2\x1b.acmelib.v2.MessagePriorityR\bpriority\x12\x1d\n" +
	"\n" +
	"cycle_time\x18\b \x01(\rR\tcycleTime\x128\n" +
	"\tsend_type\x18\t \x01(\x0e2\x1b.acmelib.v2.MessageSendTypeR\bsendType\x12\x1d\n" +
	"\n" +
	"delay_time\x18\n" +
	" \x01(\rR\tdelayTime\x12(\n" +
	"\x10start_delay_time\x18\v \x01(\rR\x0estartDelayTime\x129\n" +
	"\treceivers\x18\f \x03(\v2\x1b.acmelib.v2.MessageReceiverR\treceivers\x12T\n" +
	"\x15attribute_assignments\x18\r \x03(\v2\x1f.acmelib.v2.AttributeAssignmentR\x14attributeAssignments\"k\n" +
	"\x0fMessageReceiver\x12$\n" +
	"\x0enode_entity_id\x18\x01 \x01(\tR\fnodeEntityId\x122\n" +
	"\x15node_interface_number\x18\x02 \x01(\rR\x13nodeInterfaceNumber*\xa5\x01\n" +
	"\x0fMessagePriority\x12 \n" +
	"\x1cMESSAGE_PRIORITY_UNSPECIFIED\x10\x00\x12\x1e\n" +
	"\x1aMESSAGE_PRIORITY_VERY_HIGH\x10\x01\x12\x19\n" +
	"\x15MESSAGE_PRIORITY_HIGH\x10\x02\x12\x1b\n" +
	"\x17MESSAGE_PRIORITY_MEDIUM\x10\x03\x12\x18\n" +
	"\x14MESSAGE_PRIORITY_LOW\x10\x04*\xdc\x01\n" +
	"\x0fMessageSendType\x12!\n" +
	"\x1dMESSAGE_SEND_TYPE_UNSPECIFIED\x10\x00\x12\x1c\n" +
	"\x18MESSAGE_SEND_TYPE_CYCLIC\x10\x01\x12&\n" +
	"\"MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE\x10\x02\x12*\n" +
	"&MESSAGE_SEND_TYPE_CYCLIC_AND_TRIGGERED\x10\x03\x124\n" +
	"0MESSAGE_SEND_TYPE_CYCLIC_IF_ACTIVE_AND_TRIGGERED\x10\x04B}\n" +
	"\x0ecom.acmelib.v2B\fMessageProtoP\x01Z\x14acmelib/v2;acmelibv2\xa2\x02\x03AXX\xaa\x02\n" +
	"Acmelib.V2\xca\x02\n" +
	"Acmelib\\V2\xe2\x02\x16Acmelib\\V2\\GPBMetadata\xea\x02\vAcmelib::V2b\x06proto3"

var (
	file_acmelib_v2_message_proto_rawDescOnce sync.Once
	file_acmelib_v2_message_proto_rawDescData []byte
)

func file_acmelib_v2_message_proto_rawDescGZIP() []byte {
	file_acmelib_v2_message_proto_rawDescOnce.Do(func() {
		file_acmelib_v2_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_acmelib_v2_message_proto_rawDesc), len(file_acmelib_v2_message_proto_rawDesc)))
	})
	return file_acmelib_v2_message_proto_rawDescData
}

var file_acmelib_v2_message_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_acmelib_v2_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_acmelib_v2_message_proto_goTypes = []any{
	(MessagePriority)(0),        // 0: acmelib.v2.MessagePriority
	(MessageSendType)(0),        // 1: acmelib.v2.MessageSendType
	(*Message)(nil),             // 2: acmelib.v2.Message
	(*MessageReceiver)(nil),     // 3: acmelib.v2.MessageReceiver
	(*Entity)(nil),              // 4: acmelib.v2.Entity
	(*SignalLayout)(nil),        // 5: acmelib.v2.SignalLayout
	(*AttributeAssignment)(nil), // 6: acmelib.v2.AttributeAssignment
}
var file_acmelib_v2_message_proto_depIdxs = []int32{
	4, // 0: acmelib.v2.Message.entity:type_name -> acmelib.v2.Entity
	5, // 1: acmelib.v2.Message.layout:type_name -> acmelib.v2.SignalLayout
	0, // 2: acmelib.v2.Message.priority:type_name -> acmelib.v2.MessagePriority
	1, // 3: acmelib.v2.Message.send_type:type_name -> acmelib.v2.MessageSendType
	3, // 4: acmelib.v2.Message.receivers:type_name -> acmelib.v2.MessageReceiver
	6, // 5: acmelib.v2.Message.attribute_assignments:type_name -> acmelib.v2.AttributeAssignment
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_acmelib_v2_message_proto_init() }
func file_acmelib_v2_message_proto_init() {
	if File_acmelib_v2_message_proto != nil {
		return
	}
	file_acmelib_v2_entity_proto_init()
	file_acmelib_v2_signal_proto_init()
	file_acmelib_v2_attribute_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_acmelib_v2_message_proto_rawDesc), len(file_acmelib_v2_message_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_acmelib_v2_message_proto_goTypes,
		DependencyIndexes: file_acmelib_v2_message_proto_depIdxs,
		EnumInfos:         file_acmelib_v2_message_proto_enumTypes,
		MessageInfos:      file_acmelib_v2_message_proto_msgTypes,
	}.Build()
	File_acmelib_v2_message_proto = out.File
	file_acmelib_v2_message_proto_goTypes = nil
	file_acmelib_v2_message_proto_depIdxs = nil
}

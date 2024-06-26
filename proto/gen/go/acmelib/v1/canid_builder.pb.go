// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: acmelib/v1/canid_builder.proto

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

type CANIDBuilderOpKind int32

const (
	CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_UNSPECIFIED      CANIDBuilderOpKind = 0
	CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_PRIORITY CANIDBuilderOpKind = 1
	CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_MESSAGE_ID       CANIDBuilderOpKind = 2
	CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_NODE_ID          CANIDBuilderOpKind = 3
	CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_BIT_MASK         CANIDBuilderOpKind = 4
)

// Enum value maps for CANIDBuilderOpKind.
var (
	CANIDBuilderOpKind_name = map[int32]string{
		0: "CANID_BUILDER_OP_KIND_UNSPECIFIED",
		1: "CANID_BUILDER_OP_KIND_MESSAGE_PRIORITY",
		2: "CANID_BUILDER_OP_KIND_MESSAGE_ID",
		3: "CANID_BUILDER_OP_KIND_NODE_ID",
		4: "CANID_BUILDER_OP_KIND_BIT_MASK",
	}
	CANIDBuilderOpKind_value = map[string]int32{
		"CANID_BUILDER_OP_KIND_UNSPECIFIED":      0,
		"CANID_BUILDER_OP_KIND_MESSAGE_PRIORITY": 1,
		"CANID_BUILDER_OP_KIND_MESSAGE_ID":       2,
		"CANID_BUILDER_OP_KIND_NODE_ID":          3,
		"CANID_BUILDER_OP_KIND_BIT_MASK":         4,
	}
)

func (x CANIDBuilderOpKind) Enum() *CANIDBuilderOpKind {
	p := new(CANIDBuilderOpKind)
	*p = x
	return p
}

func (x CANIDBuilderOpKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CANIDBuilderOpKind) Descriptor() protoreflect.EnumDescriptor {
	return file_acmelib_v1_canid_builder_proto_enumTypes[0].Descriptor()
}

func (CANIDBuilderOpKind) Type() protoreflect.EnumType {
	return &file_acmelib_v1_canid_builder_proto_enumTypes[0]
}

func (x CANIDBuilderOpKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CANIDBuilderOpKind.Descriptor instead.
func (CANIDBuilderOpKind) EnumDescriptor() ([]byte, []int) {
	return file_acmelib_v1_canid_builder_proto_rawDescGZIP(), []int{0}
}

type CANIDBuilderOp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind CANIDBuilderOpKind `protobuf:"varint,1,opt,name=kind,proto3,enum=acmelib.v1.CANIDBuilderOpKind" json:"kind,omitempty"`
	From uint32             `protobuf:"varint,2,opt,name=from,proto3" json:"from,omitempty"`
	Len  uint32             `protobuf:"varint,3,opt,name=len,proto3" json:"len,omitempty"`
}

func (x *CANIDBuilderOp) Reset() {
	*x = CANIDBuilderOp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acmelib_v1_canid_builder_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CANIDBuilderOp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CANIDBuilderOp) ProtoMessage() {}

func (x *CANIDBuilderOp) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v1_canid_builder_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CANIDBuilderOp.ProtoReflect.Descriptor instead.
func (*CANIDBuilderOp) Descriptor() ([]byte, []int) {
	return file_acmelib_v1_canid_builder_proto_rawDescGZIP(), []int{0}
}

func (x *CANIDBuilderOp) GetKind() CANIDBuilderOpKind {
	if x != nil {
		return x.Kind
	}
	return CANIDBuilderOpKind_CANID_BUILDER_OP_KIND_UNSPECIFIED
}

func (x *CANIDBuilderOp) GetFrom() uint32 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *CANIDBuilderOp) GetLen() uint32 {
	if x != nil {
		return x.Len
	}
	return 0
}

type CANIDBuilder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entity     *Entity           `protobuf:"bytes,1,opt,name=entity,proto3" json:"entity,omitempty"`
	Operations []*CANIDBuilderOp `protobuf:"bytes,2,rep,name=operations,proto3" json:"operations,omitempty"`
}

func (x *CANIDBuilder) Reset() {
	*x = CANIDBuilder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acmelib_v1_canid_builder_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CANIDBuilder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CANIDBuilder) ProtoMessage() {}

func (x *CANIDBuilder) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v1_canid_builder_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CANIDBuilder.ProtoReflect.Descriptor instead.
func (*CANIDBuilder) Descriptor() ([]byte, []int) {
	return file_acmelib_v1_canid_builder_proto_rawDescGZIP(), []int{1}
}

func (x *CANIDBuilder) GetEntity() *Entity {
	if x != nil {
		return x.Entity
	}
	return nil
}

func (x *CANIDBuilder) GetOperations() []*CANIDBuilderOp {
	if x != nil {
		return x.Operations
	}
	return nil
}

var File_acmelib_v1_canid_builder_proto protoreflect.FileDescriptor

var file_acmelib_v1_canid_builder_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x6e,
	0x69, 0x64, 0x5f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0a, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x17, 0x61, 0x63,
	0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6a, 0x0a, 0x0e, 0x43, 0x41, 0x4e, 0x49, 0x44, 0x42, 0x75,
	0x69, 0x6c, 0x64, 0x65, 0x72, 0x4f, 0x70, 0x12, 0x32, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x41, 0x4e, 0x49, 0x44, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x4f,
	0x70, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x66,
	0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12,
	0x10, 0x0a, 0x03, 0x6c, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6c, 0x65,
	0x6e, 0x22, 0x76, 0x0a, 0x0c, 0x43, 0x41, 0x4e, 0x49, 0x44, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x65,
	0x72, 0x12, 0x2a, 0x0a, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x3a, 0x0a,
	0x0a, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x41, 0x4e, 0x49, 0x44, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x4f, 0x70, 0x52, 0x0a, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2a, 0xd4, 0x01, 0x0a, 0x12, 0x43, 0x41,
	0x4e, 0x49, 0x44, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x65, 0x72, 0x4f, 0x70, 0x4b, 0x69, 0x6e, 0x64,
	0x12, 0x25, 0x0a, 0x21, 0x43, 0x41, 0x4e, 0x49, 0x44, 0x5f, 0x42, 0x55, 0x49, 0x4c, 0x44, 0x45,
	0x52, 0x5f, 0x4f, 0x50, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x2a, 0x0a, 0x26, 0x43, 0x41, 0x4e, 0x49, 0x44,
	0x5f, 0x42, 0x55, 0x49, 0x4c, 0x44, 0x45, 0x52, 0x5f, 0x4f, 0x50, 0x5f, 0x4b, 0x49, 0x4e, 0x44,
	0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x50, 0x52, 0x49, 0x4f, 0x52, 0x49, 0x54,
	0x59, 0x10, 0x01, 0x12, 0x24, 0x0a, 0x20, 0x43, 0x41, 0x4e, 0x49, 0x44, 0x5f, 0x42, 0x55, 0x49,
	0x4c, 0x44, 0x45, 0x52, 0x5f, 0x4f, 0x50, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x4d, 0x45, 0x53,
	0x53, 0x41, 0x47, 0x45, 0x5f, 0x49, 0x44, 0x10, 0x02, 0x12, 0x21, 0x0a, 0x1d, 0x43, 0x41, 0x4e,
	0x49, 0x44, 0x5f, 0x42, 0x55, 0x49, 0x4c, 0x44, 0x45, 0x52, 0x5f, 0x4f, 0x50, 0x5f, 0x4b, 0x49,
	0x4e, 0x44, 0x5f, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x49, 0x44, 0x10, 0x03, 0x12, 0x22, 0x0a, 0x1e,
	0x43, 0x41, 0x4e, 0x49, 0x44, 0x5f, 0x42, 0x55, 0x49, 0x4c, 0x44, 0x45, 0x52, 0x5f, 0x4f, 0x50,
	0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x42, 0x49, 0x54, 0x5f, 0x4d, 0x41, 0x53, 0x4b, 0x10, 0x04,
	0x42, 0x82, 0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62,
	0x2e, 0x76, 0x31, 0x42, 0x11, 0x43, 0x61, 0x6e, 0x69, 0x64, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x65,
	0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x14, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69,
	0x62, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x0a, 0x41, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x56,
	0x31, 0xca, 0x02, 0x0a, 0x41, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x16, 0x41, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x41, 0x63, 0x6d, 0x65, 0x6c, 0x69,
	0x62, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_acmelib_v1_canid_builder_proto_rawDescOnce sync.Once
	file_acmelib_v1_canid_builder_proto_rawDescData = file_acmelib_v1_canid_builder_proto_rawDesc
)

func file_acmelib_v1_canid_builder_proto_rawDescGZIP() []byte {
	file_acmelib_v1_canid_builder_proto_rawDescOnce.Do(func() {
		file_acmelib_v1_canid_builder_proto_rawDescData = protoimpl.X.CompressGZIP(file_acmelib_v1_canid_builder_proto_rawDescData)
	})
	return file_acmelib_v1_canid_builder_proto_rawDescData
}

var file_acmelib_v1_canid_builder_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_acmelib_v1_canid_builder_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_acmelib_v1_canid_builder_proto_goTypes = []interface{}{
	(CANIDBuilderOpKind)(0), // 0: acmelib.v1.CANIDBuilderOpKind
	(*CANIDBuilderOp)(nil),  // 1: acmelib.v1.CANIDBuilderOp
	(*CANIDBuilder)(nil),    // 2: acmelib.v1.CANIDBuilder
	(*Entity)(nil),          // 3: acmelib.v1.Entity
}
var file_acmelib_v1_canid_builder_proto_depIdxs = []int32{
	0, // 0: acmelib.v1.CANIDBuilderOp.kind:type_name -> acmelib.v1.CANIDBuilderOpKind
	3, // 1: acmelib.v1.CANIDBuilder.entity:type_name -> acmelib.v1.Entity
	1, // 2: acmelib.v1.CANIDBuilder.operations:type_name -> acmelib.v1.CANIDBuilderOp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_acmelib_v1_canid_builder_proto_init() }
func file_acmelib_v1_canid_builder_proto_init() {
	if File_acmelib_v1_canid_builder_proto != nil {
		return
	}
	file_acmelib_v1_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_acmelib_v1_canid_builder_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CANIDBuilderOp); i {
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
		file_acmelib_v1_canid_builder_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CANIDBuilder); i {
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
			RawDescriptor: file_acmelib_v1_canid_builder_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_acmelib_v1_canid_builder_proto_goTypes,
		DependencyIndexes: file_acmelib_v1_canid_builder_proto_depIdxs,
		EnumInfos:         file_acmelib_v1_canid_builder_proto_enumTypes,
		MessageInfos:      file_acmelib_v1_canid_builder_proto_msgTypes,
	}.Build()
	File_acmelib_v1_canid_builder_proto = out.File
	file_acmelib_v1_canid_builder_proto_rawDesc = nil
	file_acmelib_v1_canid_builder_proto_goTypes = nil
	file_acmelib_v1_canid_builder_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: acmelib/v1/attribute.proto

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

type AttributeType int32

const (
	AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED AttributeType = 0
	AttributeType_ATTRIBUTE_TYPE_STRING      AttributeType = 1
	AttributeType_ATTRIBUTE_TYPE_INTEGER     AttributeType = 2
	AttributeType_ATTRIBUTE_TYPE_FLOAT       AttributeType = 3
	AttributeType_ATTRIBUTE_TYPE_ENUM        AttributeType = 4
)

// Enum value maps for AttributeType.
var (
	AttributeType_name = map[int32]string{
		0: "ATTRIBUTE_TYPE_UNSPECIFIED",
		1: "ATTRIBUTE_TYPE_STRING",
		2: "ATTRIBUTE_TYPE_INTEGER",
		3: "ATTRIBUTE_TYPE_FLOAT",
		4: "ATTRIBUTE_TYPE_ENUM",
	}
	AttributeType_value = map[string]int32{
		"ATTRIBUTE_TYPE_UNSPECIFIED": 0,
		"ATTRIBUTE_TYPE_STRING":      1,
		"ATTRIBUTE_TYPE_INTEGER":     2,
		"ATTRIBUTE_TYPE_FLOAT":       3,
		"ATTRIBUTE_TYPE_ENUM":        4,
	}
)

func (x AttributeType) Enum() *AttributeType {
	p := new(AttributeType)
	*p = x
	return p
}

func (x AttributeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AttributeType) Descriptor() protoreflect.EnumDescriptor {
	return file_acmelib_v1_attribute_proto_enumTypes[0].Descriptor()
}

func (AttributeType) Type() protoreflect.EnumType {
	return &file_acmelib_v1_attribute_proto_enumTypes[0]
}

func (x AttributeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AttributeType.Descriptor instead.
func (AttributeType) EnumDescriptor() ([]byte, []int) {
	return file_acmelib_v1_attribute_proto_rawDescGZIP(), []int{0}
}

type Attribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entity *Entity       `protobuf:"bytes,1,opt,name=entity,proto3" json:"entity,omitempty"`
	Type   AttributeType `protobuf:"varint,2,opt,name=type,proto3,enum=acmelib.v1.AttributeType" json:"type,omitempty"`
	// Types that are assignable to Attribute:
	//
	//	*Attribute_StringAttribute
	//	*Attribute_IntegerAttribute
	//	*Attribute_FloatAttribute
	//	*Attribute_EnumAttribute
	Attribute isAttribute_Attribute `protobuf_oneof:"attribute"`
}

func (x *Attribute) Reset() {
	*x = Attribute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acmelib_v1_attribute_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attribute) ProtoMessage() {}

func (x *Attribute) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v1_attribute_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attribute.ProtoReflect.Descriptor instead.
func (*Attribute) Descriptor() ([]byte, []int) {
	return file_acmelib_v1_attribute_proto_rawDescGZIP(), []int{0}
}

func (x *Attribute) GetEntity() *Entity {
	if x != nil {
		return x.Entity
	}
	return nil
}

func (x *Attribute) GetType() AttributeType {
	if x != nil {
		return x.Type
	}
	return AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED
}

func (m *Attribute) GetAttribute() isAttribute_Attribute {
	if m != nil {
		return m.Attribute
	}
	return nil
}

func (x *Attribute) GetStringAttribute() *StringAttribute {
	if x, ok := x.GetAttribute().(*Attribute_StringAttribute); ok {
		return x.StringAttribute
	}
	return nil
}

func (x *Attribute) GetIntegerAttribute() *IntegerAttribute {
	if x, ok := x.GetAttribute().(*Attribute_IntegerAttribute); ok {
		return x.IntegerAttribute
	}
	return nil
}

func (x *Attribute) GetFloatAttribute() *FloatAttribute {
	if x, ok := x.GetAttribute().(*Attribute_FloatAttribute); ok {
		return x.FloatAttribute
	}
	return nil
}

func (x *Attribute) GetEnumAttribute() *EnumAttribute {
	if x, ok := x.GetAttribute().(*Attribute_EnumAttribute); ok {
		return x.EnumAttribute
	}
	return nil
}

type isAttribute_Attribute interface {
	isAttribute_Attribute()
}

type Attribute_StringAttribute struct {
	StringAttribute *StringAttribute `protobuf:"bytes,3,opt,name=string_attribute,json=stringAttribute,proto3,oneof"`
}

type Attribute_IntegerAttribute struct {
	IntegerAttribute *IntegerAttribute `protobuf:"bytes,4,opt,name=integer_attribute,json=integerAttribute,proto3,oneof"`
}

type Attribute_FloatAttribute struct {
	FloatAttribute *FloatAttribute `protobuf:"bytes,5,opt,name=float_attribute,json=floatAttribute,proto3,oneof"`
}

type Attribute_EnumAttribute struct {
	EnumAttribute *EnumAttribute `protobuf:"bytes,6,opt,name=enum_attribute,json=enumAttribute,proto3,oneof"`
}

func (*Attribute_StringAttribute) isAttribute_Attribute() {}

func (*Attribute_IntegerAttribute) isAttribute_Attribute() {}

func (*Attribute_FloatAttribute) isAttribute_Attribute() {}

func (*Attribute_EnumAttribute) isAttribute_Attribute() {}

type StringAttribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DefValue string `protobuf:"bytes,1,opt,name=def_value,json=defValue,proto3" json:"def_value,omitempty"`
}

func (x *StringAttribute) Reset() {
	*x = StringAttribute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acmelib_v1_attribute_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringAttribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringAttribute) ProtoMessage() {}

func (x *StringAttribute) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v1_attribute_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringAttribute.ProtoReflect.Descriptor instead.
func (*StringAttribute) Descriptor() ([]byte, []int) {
	return file_acmelib_v1_attribute_proto_rawDescGZIP(), []int{1}
}

func (x *StringAttribute) GetDefValue() string {
	if x != nil {
		return x.DefValue
	}
	return ""
}

type IntegerAttribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DefValue    int32 `protobuf:"varint,1,opt,name=def_value,json=defValue,proto3" json:"def_value,omitempty"`
	Min         int32 `protobuf:"varint,2,opt,name=min,proto3" json:"min,omitempty"`
	Max         int32 `protobuf:"varint,3,opt,name=max,proto3" json:"max,omitempty"`
	IsHexFormat bool  `protobuf:"varint,4,opt,name=is_hex_format,json=isHexFormat,proto3" json:"is_hex_format,omitempty"`
}

func (x *IntegerAttribute) Reset() {
	*x = IntegerAttribute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acmelib_v1_attribute_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntegerAttribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntegerAttribute) ProtoMessage() {}

func (x *IntegerAttribute) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v1_attribute_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntegerAttribute.ProtoReflect.Descriptor instead.
func (*IntegerAttribute) Descriptor() ([]byte, []int) {
	return file_acmelib_v1_attribute_proto_rawDescGZIP(), []int{2}
}

func (x *IntegerAttribute) GetDefValue() int32 {
	if x != nil {
		return x.DefValue
	}
	return 0
}

func (x *IntegerAttribute) GetMin() int32 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *IntegerAttribute) GetMax() int32 {
	if x != nil {
		return x.Max
	}
	return 0
}

func (x *IntegerAttribute) GetIsHexFormat() bool {
	if x != nil {
		return x.IsHexFormat
	}
	return false
}

type FloatAttribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DefValue float64 `protobuf:"fixed64,1,opt,name=def_value,json=defValue,proto3" json:"def_value,omitempty"`
	Min      float64 `protobuf:"fixed64,2,opt,name=min,proto3" json:"min,omitempty"`
	Max      float64 `protobuf:"fixed64,3,opt,name=max,proto3" json:"max,omitempty"`
}

func (x *FloatAttribute) Reset() {
	*x = FloatAttribute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acmelib_v1_attribute_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FloatAttribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FloatAttribute) ProtoMessage() {}

func (x *FloatAttribute) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v1_attribute_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FloatAttribute.ProtoReflect.Descriptor instead.
func (*FloatAttribute) Descriptor() ([]byte, []int) {
	return file_acmelib_v1_attribute_proto_rawDescGZIP(), []int{3}
}

func (x *FloatAttribute) GetDefValue() float64 {
	if x != nil {
		return x.DefValue
	}
	return 0
}

func (x *FloatAttribute) GetMin() float64 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *FloatAttribute) GetMax() float64 {
	if x != nil {
		return x.Max
	}
	return 0
}

type EnumAttribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DefValue string   `protobuf:"bytes,1,opt,name=def_value,json=defValue,proto3" json:"def_value,omitempty"`
	Values   []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *EnumAttribute) Reset() {
	*x = EnumAttribute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acmelib_v1_attribute_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnumAttribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnumAttribute) ProtoMessage() {}

func (x *EnumAttribute) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v1_attribute_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnumAttribute.ProtoReflect.Descriptor instead.
func (*EnumAttribute) Descriptor() ([]byte, []int) {
	return file_acmelib_v1_attribute_proto_rawDescGZIP(), []int{4}
}

func (x *EnumAttribute) GetDefValue() string {
	if x != nil {
		return x.DefValue
	}
	return ""
}

func (x *EnumAttribute) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type AttributeAssignment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntityId          string     `protobuf:"bytes,1,opt,name=entity_id,json=entityId,proto3" json:"entity_id,omitempty"`
	EntityKind        EntityKind `protobuf:"varint,2,opt,name=entity_kind,json=entityKind,proto3,enum=acmelib.v1.EntityKind" json:"entity_kind,omitempty"`
	AttributeEntityId string     `protobuf:"bytes,3,opt,name=attribute_entity_id,json=attributeEntityId,proto3" json:"attribute_entity_id,omitempty"`
	// Types that are assignable to Value:
	//
	//	*AttributeAssignment_ValueString
	//	*AttributeAssignment_ValueInt
	//	*AttributeAssignment_ValueDouble
	Value isAttributeAssignment_Value `protobuf_oneof:"value"`
}

func (x *AttributeAssignment) Reset() {
	*x = AttributeAssignment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_acmelib_v1_attribute_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttributeAssignment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttributeAssignment) ProtoMessage() {}

func (x *AttributeAssignment) ProtoReflect() protoreflect.Message {
	mi := &file_acmelib_v1_attribute_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttributeAssignment.ProtoReflect.Descriptor instead.
func (*AttributeAssignment) Descriptor() ([]byte, []int) {
	return file_acmelib_v1_attribute_proto_rawDescGZIP(), []int{5}
}

func (x *AttributeAssignment) GetEntityId() string {
	if x != nil {
		return x.EntityId
	}
	return ""
}

func (x *AttributeAssignment) GetEntityKind() EntityKind {
	if x != nil {
		return x.EntityKind
	}
	return EntityKind_ENTITY_KIND_UNSPECIFIED
}

func (x *AttributeAssignment) GetAttributeEntityId() string {
	if x != nil {
		return x.AttributeEntityId
	}
	return ""
}

func (m *AttributeAssignment) GetValue() isAttributeAssignment_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *AttributeAssignment) GetValueString() string {
	if x, ok := x.GetValue().(*AttributeAssignment_ValueString); ok {
		return x.ValueString
	}
	return ""
}

func (x *AttributeAssignment) GetValueInt() int32 {
	if x, ok := x.GetValue().(*AttributeAssignment_ValueInt); ok {
		return x.ValueInt
	}
	return 0
}

func (x *AttributeAssignment) GetValueDouble() float64 {
	if x, ok := x.GetValue().(*AttributeAssignment_ValueDouble); ok {
		return x.ValueDouble
	}
	return 0
}

type isAttributeAssignment_Value interface {
	isAttributeAssignment_Value()
}

type AttributeAssignment_ValueString struct {
	ValueString string `protobuf:"bytes,4,opt,name=value_string,json=valueString,proto3,oneof"`
}

type AttributeAssignment_ValueInt struct {
	ValueInt int32 `protobuf:"varint,5,opt,name=value_int,json=valueInt,proto3,oneof"`
}

type AttributeAssignment_ValueDouble struct {
	ValueDouble float64 `protobuf:"fixed64,6,opt,name=value_double,json=valueDouble,proto3,oneof"`
}

func (*AttributeAssignment_ValueString) isAttributeAssignment_Value() {}

func (*AttributeAssignment_ValueInt) isAttributeAssignment_Value() {}

func (*AttributeAssignment_ValueDouble) isAttributeAssignment_Value() {}

var File_acmelib_v1_attribute_proto protoreflect.FileDescriptor

var file_acmelib_v1_attribute_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x61, 0x63,
	0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x17, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69,
	0x62, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x95, 0x03, 0x0a, 0x09, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12,
	0x2a, 0x0a, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x52, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x2d, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x61, 0x63, 0x6d, 0x65,
	0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x48, 0x0a, 0x10, 0x73, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x48, 0x00, 0x52, 0x0f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x12, 0x4b, 0x0a, 0x11, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x65, 0x72, 0x5f,
	0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x74,
	0x65, 0x67, 0x65, 0x72, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x48, 0x00, 0x52,
	0x10, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x65, 0x72, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x12, 0x45, 0x0a, 0x0f, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x63, 0x6d,
	0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x41, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x48, 0x00, 0x52, 0x0e, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x41,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x42, 0x0a, 0x0e, 0x65, 0x6e, 0x75, 0x6d,
	0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e,
	0x75, 0x6d, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x48, 0x00, 0x52, 0x0d, 0x65,
	0x6e, 0x75, 0x6d, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x42, 0x0b, 0x0a, 0x09,
	0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x22, 0x2e, 0x0a, 0x0f, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x64, 0x65, 0x66, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x64, 0x65, 0x66, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x77, 0x0a, 0x10, 0x49, 0x6e, 0x74,
	0x65, 0x67, 0x65, 0x72, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x64, 0x65, 0x66, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x64, 0x65, 0x66, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03,
	0x6d, 0x61, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x12, 0x22,
	0x0a, 0x0d, 0x69, 0x73, 0x5f, 0x68, 0x65, 0x78, 0x5f, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x48, 0x65, 0x78, 0x46, 0x6f, 0x72, 0x6d,
	0x61, 0x74, 0x22, 0x51, 0x0a, 0x0e, 0x46, 0x6c, 0x6f, 0x61, 0x74, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x5f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x64, 0x65, 0x66, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03,
	0x6d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x03, 0x6d, 0x61, 0x78, 0x22, 0x44, 0x0a, 0x0d, 0x45, 0x6e, 0x75, 0x6d, 0x41, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x66, 0x5f, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x66, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x8d, 0x02, 0x0a, 0x13,
	0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64,
	0x12, 0x37, 0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6b, 0x69, 0x6e, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x0a, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x2e, 0x0a, 0x13, 0x61, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x5f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0c, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x0b, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x1d,
	0x0a, 0x09, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x69, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x48, 0x00, 0x52, 0x08, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x49, 0x6e, 0x74, 0x12, 0x23, 0x0a,
	0x0c, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x0b, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x44, 0x6f, 0x75, 0x62,
	0x6c, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2a, 0x99, 0x01, 0x0a, 0x0d,
	0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a,
	0x1a, 0x41, 0x54, 0x54, 0x52, 0x49, 0x42, 0x55, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x19, 0x0a,
	0x15, 0x41, 0x54, 0x54, 0x52, 0x49, 0x42, 0x55, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x53, 0x54, 0x52, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x41, 0x54, 0x54, 0x52,
	0x49, 0x42, 0x55, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x47,
	0x45, 0x52, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x54, 0x54, 0x52, 0x49, 0x42, 0x55, 0x54,
	0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x46, 0x4c, 0x4f, 0x41, 0x54, 0x10, 0x03, 0x12, 0x17,
	0x0a, 0x13, 0x41, 0x54, 0x54, 0x52, 0x49, 0x42, 0x55, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x45, 0x4e, 0x55, 0x4d, 0x10, 0x04, 0x42, 0x7f, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x61,
	0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x2e, 0x76, 0x31, 0x42, 0x0e, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x14, 0x61, 0x63, 0x6d,
	0x65, 0x6c, 0x69, 0x62, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x0a, 0x41, 0x63, 0x6d, 0x65, 0x6c, 0x69,
	0x62, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0a, 0x41, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x5c, 0x56,
	0x31, 0xe2, 0x02, 0x16, 0x41, 0x63, 0x6d, 0x65, 0x6c, 0x69, 0x62, 0x5c, 0x56, 0x31, 0x5c, 0x47,
	0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x41, 0x63, 0x6d,
	0x65, 0x6c, 0x69, 0x62, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_acmelib_v1_attribute_proto_rawDescOnce sync.Once
	file_acmelib_v1_attribute_proto_rawDescData = file_acmelib_v1_attribute_proto_rawDesc
)

func file_acmelib_v1_attribute_proto_rawDescGZIP() []byte {
	file_acmelib_v1_attribute_proto_rawDescOnce.Do(func() {
		file_acmelib_v1_attribute_proto_rawDescData = protoimpl.X.CompressGZIP(file_acmelib_v1_attribute_proto_rawDescData)
	})
	return file_acmelib_v1_attribute_proto_rawDescData
}

var file_acmelib_v1_attribute_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_acmelib_v1_attribute_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_acmelib_v1_attribute_proto_goTypes = []interface{}{
	(AttributeType)(0),          // 0: acmelib.v1.AttributeType
	(*Attribute)(nil),           // 1: acmelib.v1.Attribute
	(*StringAttribute)(nil),     // 2: acmelib.v1.StringAttribute
	(*IntegerAttribute)(nil),    // 3: acmelib.v1.IntegerAttribute
	(*FloatAttribute)(nil),      // 4: acmelib.v1.FloatAttribute
	(*EnumAttribute)(nil),       // 5: acmelib.v1.EnumAttribute
	(*AttributeAssignment)(nil), // 6: acmelib.v1.AttributeAssignment
	(*Entity)(nil),              // 7: acmelib.v1.Entity
	(EntityKind)(0),             // 8: acmelib.v1.EntityKind
}
var file_acmelib_v1_attribute_proto_depIdxs = []int32{
	7, // 0: acmelib.v1.Attribute.entity:type_name -> acmelib.v1.Entity
	0, // 1: acmelib.v1.Attribute.type:type_name -> acmelib.v1.AttributeType
	2, // 2: acmelib.v1.Attribute.string_attribute:type_name -> acmelib.v1.StringAttribute
	3, // 3: acmelib.v1.Attribute.integer_attribute:type_name -> acmelib.v1.IntegerAttribute
	4, // 4: acmelib.v1.Attribute.float_attribute:type_name -> acmelib.v1.FloatAttribute
	5, // 5: acmelib.v1.Attribute.enum_attribute:type_name -> acmelib.v1.EnumAttribute
	8, // 6: acmelib.v1.AttributeAssignment.entity_kind:type_name -> acmelib.v1.EntityKind
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_acmelib_v1_attribute_proto_init() }
func file_acmelib_v1_attribute_proto_init() {
	if File_acmelib_v1_attribute_proto != nil {
		return
	}
	file_acmelib_v1_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_acmelib_v1_attribute_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attribute); i {
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
		file_acmelib_v1_attribute_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringAttribute); i {
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
		file_acmelib_v1_attribute_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntegerAttribute); i {
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
		file_acmelib_v1_attribute_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FloatAttribute); i {
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
		file_acmelib_v1_attribute_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnumAttribute); i {
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
		file_acmelib_v1_attribute_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttributeAssignment); i {
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
	file_acmelib_v1_attribute_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Attribute_StringAttribute)(nil),
		(*Attribute_IntegerAttribute)(nil),
		(*Attribute_FloatAttribute)(nil),
		(*Attribute_EnumAttribute)(nil),
	}
	file_acmelib_v1_attribute_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*AttributeAssignment_ValueString)(nil),
		(*AttributeAssignment_ValueInt)(nil),
		(*AttributeAssignment_ValueDouble)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_acmelib_v1_attribute_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_acmelib_v1_attribute_proto_goTypes,
		DependencyIndexes: file_acmelib_v1_attribute_proto_depIdxs,
		EnumInfos:         file_acmelib_v1_attribute_proto_enumTypes,
		MessageInfos:      file_acmelib_v1_attribute_proto_msgTypes,
	}.Build()
	File_acmelib_v1_attribute_proto = out.File
	file_acmelib_v1_attribute_proto_rawDesc = nil
	file_acmelib_v1_attribute_proto_goTypes = nil
	file_acmelib_v1_attribute_proto_depIdxs = nil
}
package acmelib

import (
	"fmt"
	"time"
)

type AttributeKind string

const (
	AttributeKindString  AttributeKind = "attribute-string"
	AttributeKindInteger AttributeKind = "attribute-integer"
	AttributeKindFloat   AttributeKind = "attribute-float"
	AttributeKindEnum    AttributeKind = "attribute-enum"
)

type AttributeValue struct {
	attribute Attribute
	value     any
}

func newAttributeValue(att Attribute, val any) *AttributeValue {
	return &AttributeValue{
		attribute: att,
		value:     val,
	}
}

func (av *AttributeValue) Attribute() Attribute {
	return av.attribute
}

func (av *AttributeValue) Value() any {
	return av.value
}

type AttributeRefKind string

const (
	AttributeRefKindBus     AttributeRefKind = "attribute_ref-bus"
	AttributeRefKindNode    AttributeRefKind = "attribute_ref-node"
	AttributeRefKindMessage AttributeRefKind = "attribute_ref-message"
	AttributeRefKindSignal  AttributeRefKind = "attribute_ref-signal"
)

type AttributeRef struct {
	entityID EntityID
	kind     AttributeRefKind
	value    any
}

func newAttributeRef(entID EntityID, kind AttributeRefKind, val any) *AttributeRef {
	return &AttributeRef{
		entityID: entID,
		kind:     kind,
		value:    val,
	}
}

func (af *AttributeRef) EntityID() EntityID {
	return af.entityID
}

func (af *AttributeRef) Kind() AttributeRefKind {
	return af.kind
}

func (af *AttributeRef) Value() any {
	return af.value
}

type Attribute interface {
	EntityID() EntityID
	Name() string
	Desc() string
	CreateTime() time.Time

	Kind() AttributeKind

	addReference(ref *AttributeRef)
	removeReference(refID EntityID)
	References() []*AttributeRef

	ToString() (*StringAttribute, error)
	ToInteger() (*IntegerAttribute, error)
	ToFloat() (*FloatAttribute, error)
	ToEnum() (*EnumAttribute, error)
}

type attribute struct {
	*entity

	kind AttributeKind

	references *set[EntityID, *AttributeRef]
}

func newAttribute(name, desc string, kind AttributeKind) *attribute {
	return &attribute{
		entity: newEntity(name, desc),

		kind: kind,

		references: newSet[EntityID, *AttributeRef]("reference"),
	}
}

func (a *attribute) Kind() AttributeKind {
	return a.kind
}

func (a *attribute) addReference(ref *AttributeRef) {
	a.references.add(ref.entityID, ref)
}

func (a *attribute) removeReference(refID EntityID) {
	a.references.remove(refID)
}

func (a *attribute) References() []*AttributeRef {
	return a.references.getValues()
}

type StringAttribute struct {
	*attribute

	defValue string
}

func NewStringAttribute(name, desc, defValue string) *StringAttribute {
	return &StringAttribute{
		attribute: newAttribute(name, desc, AttributeKindString),

		defValue: defValue,
	}
}

func (sa *StringAttribute) DefValue() string {
	return sa.defValue
}

func (sa *StringAttribute) ToString() (*StringAttribute, error) {
	return sa, nil
}

func (sa *StringAttribute) ToInteger() (*IntegerAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindInteger, AttributeKindString)
}

func (sa *StringAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindFloat, AttributeKindString)
}

func (sa *StringAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindEnum, AttributeKindString)
}

type IntegerAttribute struct {
	*attribute

	defValue int
	min      int
	max      int

	isHexFormat bool
}

func NewIntegerAttribute(name, desc string, defValue, min, max int) (*IntegerAttribute, error) {
	if min > max {
		return nil, fmt.Errorf("min value cannot be greater then max value")
	}

	if defValue < min || defValue > max {
		return nil, fmt.Errorf(`default value "%d" is out of min/max range ("%d" - "%d")`, defValue, min, max)
	}

	return &IntegerAttribute{
		attribute: newAttribute(name, desc, AttributeKindInteger),

		defValue: defValue,
		min:      min,
		max:      max,

		isHexFormat: false,
	}, nil
}

func (ia *IntegerAttribute) DefValue() int {
	return ia.defValue
}

func (ia *IntegerAttribute) Min() int {
	return ia.min
}

func (ia *IntegerAttribute) Max() int {
	return ia.max
}

func (ia *IntegerAttribute) SetFormatHex() {
	ia.isHexFormat = true
}

func (ia *IntegerAttribute) IsHexFormat() bool {
	return ia.isHexFormat
}

func (ia *IntegerAttribute) ToString() (*StringAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindString, AttributeKindInteger)
}

func (ia *IntegerAttribute) ToInteger() (*IntegerAttribute, error) {
	return ia, nil
}

func (ia *IntegerAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindFloat, AttributeKindInteger)
}

func (ia *IntegerAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindEnum, AttributeKindInteger)
}

type FloatAttribute struct {
	*attribute

	defValue float64
	min      float64
	max      float64
}

func NewFloatAttribute(name, desc string, defValue, min, max float64) (*FloatAttribute, error) {
	if min > max {
		return nil, fmt.Errorf("min value cannot be greater then max value")
	}

	if defValue < min || defValue > max {
		return nil, fmt.Errorf(`default value "%f" is out of min/max range ("%f" - "%f")`, defValue, min, max)
	}

	return &FloatAttribute{
		attribute: newAttribute(name, desc, AttributeKindFloat),

		defValue: defValue,
		min:      min,
		max:      max,
	}, nil
}

func (fa *FloatAttribute) DefValue() float64 {
	return fa.defValue
}

func (fa *FloatAttribute) Min() float64 {
	return fa.min
}

func (fa *FloatAttribute) Max() float64 {
	return fa.max
}

func (fa *FloatAttribute) ToString() (*StringAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindString, AttributeKindFloat)
}

func (fa *FloatAttribute) ToInteger() (*IntegerAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindInteger, AttributeKindFloat)
}

func (fa *FloatAttribute) ToFloat() (*FloatAttribute, error) {
	return fa, nil
}

func (fa *FloatAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindEnum, AttributeKindFloat)
}

type EnumAttribute struct {
	*attribute

	defValue string
	values   *set[string, int]
}

func NewEnumAttribute(name, desc string, values ...string) (*EnumAttribute, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("at least 1 value is required")
	}

	valSet := newSet[string, int]("values")
	currIdx := 0
	for _, val := range values {
		if valSet.hasKey(val) {
			continue
		}

		valSet.add(val, currIdx)
		currIdx++
	}

	return &EnumAttribute{
		attribute: newAttribute(name, desc, AttributeKindEnum),

		defValue: values[0],
		values:   valSet,
	}, nil
}

func (ea *EnumAttribute) DefValue() string {
	return ea.defValue
}

func (ea *EnumAttribute) Values() []string {
	valSlice := make([]string, ea.values.size())
	for val, valIdx := range ea.values.entries() {
		valSlice[valIdx] = val
	}
	return valSlice
}

func (ea *EnumAttribute) GetValueAtIndex(valueIndex int) (string, error) {
	if valueIndex < 0 {
		return "", fmt.Errorf("value index cannot be negative")
	}

	if valueIndex >= ea.values.size() {
		return "", fmt.Errorf(`value index "%d" is out of range ("0" - "%d")`, valueIndex, ea.references.size()-1)
	}

	return ea.Values()[valueIndex], nil
}

func (ea *EnumAttribute) ToString() (*StringAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindString, AttributeKindEnum)
}

func (ea *EnumAttribute) ToInteger() (*IntegerAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindInteger, AttributeKindEnum)
}

func (ea *EnumAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindFloat, AttributeKindEnum)
}

func (ea *EnumAttribute) ToEnum() (*EnumAttribute, error) {
	return ea, nil
}

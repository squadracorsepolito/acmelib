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

type AttributeReferenceKind string

const (
	AttributeReferenceKindBus     AttributeReferenceKind = "attribute_ref-bus"
	AttributeReferenceKindNode    AttributeReferenceKind = "attribute_ref-node"
	AttributeReferenceKindMessage AttributeReferenceKind = "attribute_ref-message"
	AttributeReferenceKindSignal  AttributeReferenceKind = "attribute_ref-signal"
)

type AttributeReference struct {
	entityID EntityID
	kind     AttributeReferenceKind
	value    any
}

func newAttributeReference(entID EntityID, kind AttributeReferenceKind, val any) *AttributeReference {
	return &AttributeReference{
		entityID: entID,
		kind:     kind,
		value:    val,
	}
}

func (af *AttributeReference) EntityID() EntityID {
	return af.entityID
}

func (af *AttributeReference) Kind() AttributeReferenceKind {
	return af.kind
}

func (af *AttributeReference) Value() any {
	return af.value
}

type Attribute interface {
	EntityID() EntityID
	Name() string
	Desc() string
	CreateTime() time.Time

	Kind() AttributeKind
	DefValue() any

	addReference(ref *AttributeReference)
	removeReference(refID EntityID)
	References() []*AttributeReference

	ToString() (*StringAttribute, error)
	ToInteger() (*IntegerAttribute, error)
	ToFloat() (*FloatAttribute, error)
	ToEnum() (*EnumAttribute, error)
}

type attribute[T any] struct {
	*entity

	kind     AttributeKind
	defValue T

	references *set[EntityID, *AttributeReference]
}

func newAttribute[T any](name, desc string, kind AttributeKind, defValue T) *attribute[T] {
	return &attribute[T]{
		entity: newEntity(name, desc),

		kind:     kind,
		defValue: defValue,

		references: newSet[EntityID, *AttributeReference]("reference"),
	}
}

func (a *attribute[T]) Kind() AttributeKind {
	return a.kind
}

func (a *attribute[T]) DefValue() any {
	return a.defValue
}

func (a *attribute[T]) addReference(ref *AttributeReference) {
	a.references.add(ref.entityID, ref)
}

func (a *attribute[T]) removeReference(refID EntityID) {
	a.references.remove(refID)
}

func (a *attribute[T]) References() []*AttributeReference {
	return a.references.getValues()
}

type StringAttribute struct {
	*attribute[string]
}

func NewStringAttribute(name, desc, defValue string) *StringAttribute {
	return &StringAttribute{
		attribute: newAttribute(name, desc, AttributeKindString, defValue),
	}
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
	*attribute[int]

	min int
	max int

	isHexFormat bool
}

func NewIntegerAttribute(name, desc string, defVal, min, max int) (*IntegerAttribute, error) {
	if min > max {
		return nil, fmt.Errorf("min value cannot be greater then max value")
	}

	if defVal < min || defVal > max {
		return nil, fmt.Errorf(`default value "%d" is out of min/max range ("%d" - "%d")`, defVal, min, max)
	}

	return &IntegerAttribute{
		attribute: newAttribute(name, desc, AttributeKindInteger, defVal),

		min: min,
		max: max,

		isHexFormat: false,
	}, nil
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
	*attribute[float64]

	min float64
	max float64
}

func NewFloatAttribute(name, desc string, defVal, min, max float64) (*FloatAttribute, error) {
	if min > max {
		return nil, fmt.Errorf("min value cannot be greater then max value")
	}

	if defVal < min || defVal > max {
		return nil, fmt.Errorf(`default value "%f" is out of min/max range ("%f" - "%f")`, defVal, min, max)
	}

	return &FloatAttribute{
		attribute: newAttribute(name, desc, AttributeKindFloat, defVal),

		min: min,
		max: max,
	}, nil
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
	*attribute[string]

	values *set[string, int]
}

func NewEnumAttribute(name, desc, defVal string, values ...string) (*EnumAttribute, error) {
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

	if !valSet.hasKey(defVal) {
		return nil, fmt.Errorf(`default value "%s" is not present in the enum`, defVal)
	}

	return &EnumAttribute{
		attribute: newAttribute(name, desc, AttributeKindEnum, defVal),

		values: valSet,
	}, nil
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

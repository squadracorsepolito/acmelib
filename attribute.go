package acmelib

import (
	"fmt"
	"strings"
	"time"
)

// AttributeValue connects a general [Attribute] to the value associated by an entity.
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

func (av *AttributeValue) stringify(b *strings.Builder, tabs int) {
	b.WriteString(fmt.Sprintf("%sname: %s; type: %s; value: %v\n",
		getTabString(tabs), av.attribute.Name(), av.attribute.Type(), av.value))
}

// Attribute returns the [Attribute] of the [AttributeValue].
func (av *AttributeValue) Attribute() Attribute {
	return av.attribute
}

// Value returns the value of the [AttributeValue].
func (av *AttributeValue) Value() any {
	return av.value
}

// AttributeRefKind defines the kind of an [AttributeRef].
type AttributeRefKind string

const (
	// AttributeRefKindBus defines a bus reference.
	AttributeRefKindBus AttributeRefKind = "attribute_ref-bus"
	// AttributeRefKindNode defines a node reference.
	AttributeRefKindNode AttributeRefKind = "attribute_ref-node"
	// AttributeRefKindMessage defines a message reference.
	AttributeRefKindMessage AttributeRefKind = "attribute_ref-message"
	// AttributeRefKindSignal defines a signal reference.
	AttributeRefKindSignal AttributeRefKind = "attribute_ref-signal"
)

// AttributeRef connects an [Attribute] to an entity and the value
// the latter has associated to the former.
// It is useful to connect an attribute to the entities that are using it.
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

func (af *AttributeRef) stringify(b *strings.Builder, tabs int) {
	b.WriteString(fmt.Sprintf("%skind: %s; entity_id: %s;value: %v\n", getTabString(tabs), af.kind, af.entityID, af.value))
}

// EntityID returns the entity id of the [AttributeRef]
func (af *AttributeRef) EntityID() EntityID {
	return af.entityID
}

// Kind returns the kind of the [AttributeRef]
func (af *AttributeRef) Kind() AttributeRefKind {
	return af.kind
}

// Value returns the value of the [AttributeRef]
func (af *AttributeRef) Value() any {
	return af.value
}

// AttributeType defines the type of an [Attribute].
type AttributeType int

const (
	// AttributeTypeString defines a string attribute.
	AttributeTypeString AttributeType = iota
	// AttributeTypeInteger defines an integer attribute.
	AttributeTypeInteger
	// AttributeTypeFloat defines a float attribute.
	AttributeTypeFloat
	// AttributeTypeEnum defines an enum attribute.
	AttributeTypeEnum
)

func (at AttributeType) String() string {
	switch at {
	case AttributeTypeString:
		return "string"
	case AttributeTypeInteger:
		return "integer"
	case AttributeTypeFloat:
		return "float"
	case AttributeTypeEnum:
		return "enum"
	default:
		return "unknown"
	}
}

// Attribute interface specifies all common methods of
// [StringAttribute], [IntegerAttribute], [FloatAttribute], and
// [EnumAttribute].
type Attribute interface {
	// EntityID returns the entity id of an attribute.
	EntityID() EntityID
	// Name returns the name of an attribute.
	Name() string
	// Desc returns the description of an attribute.
	Desc() string
	// CreateTime returns the time of creation of an attribute.
	CreateTime() time.Time

	// Type returns the kind of an attribute.
	Type() AttributeType

	addReference(ref *AttributeRef)
	removeReference(refID EntityID)
	// References returns a slice of references of an attribute.
	References() []*AttributeRef

	String() string

	// ToString converts the attribute to a string attribute.
	ToString() (*StringAttribute, error)
	// ToInteger converts the attribute to a integer attribute.
	ToInteger() (*IntegerAttribute, error)
	// ToFloat converts the attribute to a float attribute.
	ToFloat() (*FloatAttribute, error)
	// ToEnum converts the attribute to a enum attribute.
	ToEnum() (*EnumAttribute, error)
}

type attribute struct {
	*entity
	*withRefs[*AttributeRef]

	typ AttributeType
}

func newAttribute(name string, typ AttributeType) *attribute {
	return &attribute{
		entity:   newEntity(name),
		withRefs: newWithRefs[*AttributeRef](),

		typ: typ,
	}
}

func (a *attribute) errorf(err error) error {
	return &EntityError{
		Kind:     EntityKindAttribute,
		EntityID: a.entityID,
		Name:     a.name,
		Err:      err,
	}
}

func (a *attribute) stringify(b *strings.Builder, tabs int) {
	a.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)
	b.WriteString(fmt.Sprintf("%skind: %s\n", tabStr, a.typ))

	if a.refs.size() == 0 {
		return
	}

	b.WriteString(fmt.Sprintf("%sreferences:\n", tabStr))
	for _, ref := range a.References() {
		ref.stringify(b, tabs+1)
	}
}

func (a *attribute) Type() AttributeType {
	return a.typ
}

func (a *attribute) addReference(ref *AttributeRef) {
	a.refs.add(ref.entityID, ref)
}

func (a *attribute) removeReference(refID EntityID) {
	a.refs.remove(refID)
}

// StringAttribute is an [Attribute] that holds a string value.
type StringAttribute struct {
	*attribute

	defValue string
}

// NewStringAttribute creates a new [StringAttribute] with the given name,
// and default value.
func NewStringAttribute(name, defValue string) *StringAttribute {
	return &StringAttribute{
		attribute: newAttribute(name, AttributeTypeString),

		defValue: defValue,
	}
}

func (sa *StringAttribute) stringify(b *strings.Builder, tabs int) {
	sa.attribute.stringify(b, tabs)
	b.WriteString(fmt.Sprintf("%sdefault_value: %s\n", getTabString(tabs), sa.defValue))
}

func (sa *StringAttribute) String() string {
	builder := new(strings.Builder)
	sa.stringify(builder, 0)
	return builder.String()
}

// DefValue returns the default value of the [StringAttribute].
func (sa *StringAttribute) DefValue() string {
	return sa.defValue
}

// ToString returns the [StringAttribute] itself.
func (sa *StringAttribute) ToString() (*StringAttribute, error) {
	return sa, nil
}

// ToInteger always returns an error.
func (sa *StringAttribute) ToInteger() (*IntegerAttribute, error) {
	return nil, sa.errorf(&ConversionError{
		From: AttributeTypeString.String(),
		To:   AttributeTypeInteger.String(),
	})
}

// ToFloat always returns an error.
func (sa *StringAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, sa.errorf(&ConversionError{
		From: AttributeTypeString.String(),
		To:   AttributeTypeFloat.String(),
	})
}

// ToEnum always returns an error.
func (sa *StringAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, sa.errorf(&ConversionError{
		From: AttributeTypeString.String(),
		To:   AttributeTypeEnum.String(),
	})
}

// IntegerAttribute is an [Attribute] that holds an integer value.
type IntegerAttribute struct {
	*attribute

	defValue int
	min      int
	max      int

	isHexFormat bool
}

// NewIntegerAttribute creates a new [IntegerAttribute] with the given name,
// default value, min, and max.
// It may return an error if the default value is out of the min/max range,
// or if the min value is greater then the max value.
func NewIntegerAttribute(name string, defValue, min, max int) (*IntegerAttribute, error) {
	if min > max {
		return nil, &ArgumentError{
			Name: "min",
			Err:  &ErrGreaterThen{Target: "max"},
		}
	}

	if defValue > max {
		return nil, &ArgumentError{
			Name: "defValue",
			Err:  &ErrGreaterThen{Target: "max"},
		}
	}

	if defValue < min {
		return nil, &ArgumentError{
			Name: "defValue",
			Err:  &ErrLowerThen{Target: "min"},
		}
	}

	return &IntegerAttribute{
		attribute: newAttribute(name, AttributeTypeInteger),

		defValue: defValue,
		min:      min,
		max:      max,

		isHexFormat: false,
	}, nil
}

func (ia *IntegerAttribute) stringify(b *strings.Builder, tabs int) {
	ia.attribute.stringify(b, tabs)

	tabStr := getTabString(tabs)
	b.WriteString(fmt.Sprintf("%smin: %d; max: %d; hex_format: %t\n", tabStr, ia.min, ia.max, ia.isHexFormat))
	b.WriteString(fmt.Sprintf("%sdefault_value: %d\n", tabStr, ia.defValue))
}

func (ia *IntegerAttribute) String() string {
	builder := new(strings.Builder)
	ia.stringify(builder, 0)
	return builder.String()
}

// DefValue returns the default value of the [IntegerAttribute].
func (ia *IntegerAttribute) DefValue() int {
	return ia.defValue
}

// Min returns the min value of the [IntegerAttribute].
func (ia *IntegerAttribute) Min() int {
	return ia.min
}

// Max returns the max value of the [IntegerAttribute].
func (ia *IntegerAttribute) Max() int {
	return ia.max
}

// SetFormatHex sets the format of the [IntegerAttribute] to hex.
func (ia *IntegerAttribute) SetFormatHex() {
	ia.isHexFormat = true
}

// IsHexFormat reports whether the [IntegerAttribute] is in hex format.
func (ia *IntegerAttribute) IsHexFormat() bool {
	return ia.isHexFormat
}

// ToString always returns an error.
func (ia *IntegerAttribute) ToString() (*StringAttribute, error) {
	return nil, ia.errorf(&ConversionError{
		From: AttributeTypeInteger.String(),
		To:   AttributeTypeString.String(),
	})
}

// ToInteger returns the [IntegerAttribute] itself.
func (ia *IntegerAttribute) ToInteger() (*IntegerAttribute, error) {
	return ia, nil
}

// ToFloat always returns an error.
func (ia *IntegerAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, ia.errorf(&ConversionError{
		From: AttributeTypeInteger.String(),
		To:   AttributeTypeFloat.String(),
	})
}

// ToEnum always returns an error.
func (ia *IntegerAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, ia.errorf(&ConversionError{
		From: AttributeTypeInteger.String(),
		To:   AttributeTypeEnum.String(),
	})
}

// FloatAttribute is an [Attribute] that holds a float value.
type FloatAttribute struct {
	*attribute

	defValue float64
	min      float64
	max      float64
}

// NewFloatAttribute creates a new [FloatAttribute] with the given name,
// default value, min, and max.
// It may return an error if the default value is out of the min/max range,
// or if the min value is greater then the max value.
func NewFloatAttribute(name string, defValue, min, max float64) (*FloatAttribute, error) {
	if min > max {
		return nil, &ArgumentError{
			Name: "min",
			Err:  &ErrGreaterThen{Target: "max"},
		}
	}

	if defValue > max {
		return nil, &ArgumentError{
			Name: "defValue",
			Err:  &ErrGreaterThen{Target: "max"},
		}
	}

	if defValue < min {
		return nil, &ArgumentError{
			Name: "defValue",
			Err:  &ErrLowerThen{Target: "min"},
		}
	}

	return &FloatAttribute{
		attribute: newAttribute(name, AttributeTypeFloat),

		defValue: defValue,
		min:      min,
		max:      max,
	}, nil
}

func (fa *FloatAttribute) stringify(b *strings.Builder, tabs int) {
	fa.attribute.stringify(b, tabs)

	tabStr := getTabString(tabs)
	b.WriteString(fmt.Sprintf("%smin: %g; max: %g\n", tabStr, fa.min, fa.max))
	b.WriteString(fmt.Sprintf("%sdefault_value: %g\n", tabStr, fa.defValue))
}

func (fa *FloatAttribute) String() string {
	builder := new(strings.Builder)
	fa.stringify(builder, 0)
	return builder.String()
}

// DefValue returns the default value of the [FloatAttribute].
func (fa *FloatAttribute) DefValue() float64 {
	return fa.defValue
}

// Min returns the min value of the [FloatAttribute].
func (fa *FloatAttribute) Min() float64 {
	return fa.min
}

// Max returns the max value of the [FloatAttribute].
func (fa *FloatAttribute) Max() float64 {
	return fa.max
}

// ToString always returns an error.
func (fa *FloatAttribute) ToString() (*StringAttribute, error) {
	return nil, fa.errorf(&ConversionError{
		From: AttributeTypeFloat.String(),
		To:   AttributeTypeString.String(),
	})
}

// ToInteger always returns an error.
func (fa *FloatAttribute) ToInteger() (*IntegerAttribute, error) {
	return nil, fa.errorf(&ConversionError{
		From: AttributeTypeFloat.String(),
		To:   AttributeTypeInteger.String(),
	})
}

// ToFloat returns the [FloatAttribute] itself.
func (fa *FloatAttribute) ToFloat() (*FloatAttribute, error) {
	return fa, nil
}

// ToEnum always returns an error.
func (fa *FloatAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, fa.errorf(&ConversionError{
		From: AttributeTypeFloat.String(),
		To:   AttributeTypeEnum.String(),
	})
}

// EnumAttribute is an [Attribute] that holds an enum as value.
type EnumAttribute struct {
	*attribute

	defValue string
	values   *set[string, int]
}

// NewEnumAttribute creates a new [EnumAttribute] with the given name and values.
// The first value is always selected as the default one.
// It may return an error if no values are passed.
func NewEnumAttribute(name string, values ...string) (*EnumAttribute, error) {
	if len(values) == 0 {
		return nil, &ArgumentError{
			Name: "values",
			Err:  ErrIsNil,
		}
	}

	valSet := newSet[string, int]()
	currIdx := 0
	for _, val := range values {
		if valSet.hasKey(val) {
			continue
		}

		valSet.add(val, currIdx)
		currIdx++
	}

	return &EnumAttribute{
		attribute: newAttribute(name, AttributeTypeEnum),

		defValue: values[0],
		values:   valSet,
	}, nil
}

func (ea *EnumAttribute) stringify(b *strings.Builder, tabs int) {
	ea.attribute.stringify(b, tabs)

	tabStr := getTabString(tabs)

	for idx, val := range ea.Values() {
		b.WriteString(fmt.Sprintf("%value: %s; index: %d\n", tabStr, val, idx))
	}

	b.WriteString(fmt.Sprintf("%sdefault_value: %s\n", tabStr, ea.defValue))
}

func (ea *EnumAttribute) String() string {
	builder := new(strings.Builder)
	ea.stringify(builder, 0)
	return builder.String()
}

// DefValue returns the default value of the [EnumAttribute].
func (ea *EnumAttribute) DefValue() string {
	return ea.defValue
}

// Values returns the values of the [EnumAttribute] in the order specified in the factory method.
func (ea *EnumAttribute) Values() []string {
	valSlice := make([]string, ea.values.size())
	for val, valIdx := range ea.values.entries() {
		valSlice[valIdx] = val
	}
	return valSlice
}

// GetValueAtIndex returns the value at the given index.
// The index refers to the order of the values in the factory method.
// It may return an error if the index is out of range.
func (ea *EnumAttribute) GetValueAtIndex(valueIndex int) (string, error) {
	if valueIndex < 0 {
		return "", ea.errorf(&GetEntityError{
			Err: &ValueIndexError{
				Index: valueIndex,
				Err:   ErrIsNegative,
			},
		})
	}

	if valueIndex >= ea.values.size() {
		return "", ea.errorf(&GetEntityError{
			Err: &ValueIndexError{
				Index: valueIndex,
				Err:   ErrOutOfBounds,
			},
		})
	}

	return ea.Values()[valueIndex], nil
}

// ToString always returns an error.
func (ea *EnumAttribute) ToString() (*StringAttribute, error) {
	return nil, ea.errorf(&ConversionError{
		From: AttributeTypeEnum.String(),
		To:   AttributeTypeString.String(),
	})
}

// ToInteger always returns an error.
func (ea *EnumAttribute) ToInteger() (*IntegerAttribute, error) {
	return nil, ea.errorf(&ConversionError{
		From: AttributeTypeEnum.String(),
		To:   AttributeTypeInteger.String(),
	})
}

// ToFloat always returns an error.
func (ea *EnumAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, ea.errorf(&ConversionError{
		From: AttributeTypeEnum.String(),
		To:   AttributeTypeFloat.String(),
	})
}

// ToEnum returns the [EnumAttribute] itself.
func (ea *EnumAttribute) ToEnum() (*EnumAttribute, error) {
	return ea, nil
}

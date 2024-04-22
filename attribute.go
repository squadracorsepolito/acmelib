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
		getTabString(tabs), av.attribute.Name(), av.attribute.Kind(), av.value))
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

// AttributeKind defines the kind of an [Attribute].
type AttributeKind string

const (
	// AttributeKindString defines a string attribute.
	AttributeKindString AttributeKind = "attribute-string"
	// AttributeKindInteger defines an integer attribute.
	AttributeKindInteger AttributeKind = "attribute-integer"
	// AttributeKindFloat defines a float attribute.
	AttributeKindFloat AttributeKind = "attribute-float"
	// AttributeKindEnum defines an enum attribute.
	AttributeKindEnum AttributeKind = "attribute-enum"
)

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

	// Kind returns the kind of an attribute.
	Kind() AttributeKind

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

	kind AttributeKind

	references *set[EntityID, *AttributeRef]
}

func newAttribute(name string, kind AttributeKind) *attribute {
	return &attribute{
		entity: newEntity(name),

		kind: kind,

		references: newSet[EntityID, *AttributeRef](),
	}
}

func (a *attribute) errorf(err error) error {
	return &AttributeError{
		EntityID: a.entityID,
		Name:     a.name,
		Err:      err,
	}
}

func (a *attribute) stringify(b *strings.Builder, tabs int) {
	a.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)
	b.WriteString(fmt.Sprintf("%skind: %s\n", tabStr, a.kind))

	if a.references.size() == 0 {
		return
	}

	b.WriteString(fmt.Sprintf("%sreferences:\n", tabStr))
	for _, ref := range a.References() {
		ref.stringify(b, tabs+1)
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

// StringAttribute is an [Attribute] that holds a string value.
type StringAttribute struct {
	*attribute

	defValue string
}

// NewStringAttribute creates a new [StringAttribute] with the given name,
// and default value.
func NewStringAttribute(name, defValue string) *StringAttribute {
	return &StringAttribute{
		attribute: newAttribute(name, AttributeKindString),

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
	return nil, sa.errorf(&ConvertionError{
		From: string(AttributeKindString),
		To:   string(AttributeKindInteger),
	})
}

// ToFloat always returns an error.
func (sa *StringAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, sa.errorf(&ConvertionError{
		From: string(AttributeKindString),
		To:   string(AttributeKindFloat),
	})
}

// ToEnum always returns an error.
func (sa *StringAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, sa.errorf(&ConvertionError{
		From: string(AttributeKindString),
		To:   string(AttributeKindEnum),
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
			Err:  &ErrGraterThen{Target: "max"},
		}
	}

	if defValue > max {
		return nil, &ArgumentError{
			Name: "defValue",
			Err:  &ErrGraterThen{Target: "max"},
		}
	}

	if defValue < min {
		return nil, &ArgumentError{
			Name: "defValue",
			Err:  &ErrLowerThen{Target: "min"},
		}
	}

	return &IntegerAttribute{
		attribute: newAttribute(name, AttributeKindInteger),

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
	return nil, ia.errorf(&ConvertionError{
		From: string(AttributeKindInteger),
		To:   string(AttributeKindString),
	})
}

// ToInteger returns the [IntegerAttribute] itself.
func (ia *IntegerAttribute) ToInteger() (*IntegerAttribute, error) {
	return ia, nil
}

// ToFloat always returns an error.
func (ia *IntegerAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, ia.errorf(&ConvertionError{
		From: string(AttributeKindInteger),
		To:   string(AttributeKindFloat),
	})
}

// ToEnum always returns an error.
func (ia *IntegerAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, ia.errorf(&ConvertionError{
		From: string(AttributeKindInteger),
		To:   string(AttributeKindEnum),
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
			Err:  &ErrGraterThen{Target: "max"},
		}
	}

	if defValue > max {
		return nil, &ArgumentError{
			Name: "defValue",
			Err:  &ErrGraterThen{Target: "max"},
		}
	}

	if defValue < min {
		return nil, &ArgumentError{
			Name: "defValue",
			Err:  &ErrLowerThen{Target: "min"},
		}
	}

	return &FloatAttribute{
		attribute: newAttribute(name, AttributeKindFloat),

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
	return nil, fa.errorf(&ConvertionError{
		From: string(AttributeKindFloat),
		To:   string(AttributeKindString),
	})
}

// ToInteger always returns an error.
func (fa *FloatAttribute) ToInteger() (*IntegerAttribute, error) {
	return nil, fa.errorf(&ConvertionError{
		From: string(AttributeKindFloat),
		To:   string(AttributeKindInteger),
	})
}

// ToFloat returns the [FloatAttribute] itself.
func (fa *FloatAttribute) ToFloat() (*FloatAttribute, error) {
	return fa, nil
}

// ToEnum always returns an error.
func (fa *FloatAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, fa.errorf(&ConvertionError{
		From: string(AttributeKindFloat),
		To:   string(AttributeKindEnum),
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
		attribute: newAttribute(name, AttributeKindEnum),

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
	return nil, ea.errorf(&ConvertionError{
		From: string(AttributeKindEnum),
		To:   string(AttributeKindString),
	})
}

// ToInteger always returns an error.
func (ea *EnumAttribute) ToInteger() (*IntegerAttribute, error) {
	return nil, ea.errorf(&ConvertionError{
		From: string(AttributeKindEnum),
		To:   string(AttributeKindInteger),
	})
}

// ToFloat always returns an error.
func (ea *EnumAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, ea.errorf(&ConvertionError{
		From: string(AttributeKindEnum),
		To:   string(AttributeKindFloat),
	})
}

// ToEnum returns the [EnumAttribute] itself.
func (ea *EnumAttribute) ToEnum() (*EnumAttribute, error) {
	return ea, nil
}

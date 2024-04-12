package acmelib

import (
	"fmt"
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

// StringAttribute is an [Attribute] that holds a string value.
type StringAttribute struct {
	*attribute

	defValue string
}

// NewStringAttribute creates a new [StringAttribute] with the given name, description,
// and default value.
func NewStringAttribute(name, desc, defValue string) *StringAttribute {
	return &StringAttribute{
		attribute: newAttribute(name, desc, AttributeKindString),

		defValue: defValue,
	}
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
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindInteger, AttributeKindString)
}

// ToFloat always returns an error.
func (sa *StringAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindFloat, AttributeKindString)
}

// ToEnum always returns an error.
func (sa *StringAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindEnum, AttributeKindString)
}

// IntegerAttribute is an [Attribute] that holds an integer value.
type IntegerAttribute struct {
	*attribute

	defValue int
	min      int
	max      int

	isHexFormat bool
}

// NewIntegerAttribute creates a new [IntegerAttribute] with the given name, description,
// default value, min, and max.
// It may return an error if the default value is out of the min/max range,
// or if the min value is greater then the max value.
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
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindString, AttributeKindInteger)
}

// ToInteger returns the [IntegerAttribute] itself.
func (ia *IntegerAttribute) ToInteger() (*IntegerAttribute, error) {
	return ia, nil
}

// ToFloat always returns an error.
func (ia *IntegerAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindFloat, AttributeKindInteger)
}

// ToEnum always returns an error.
func (ia *IntegerAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindEnum, AttributeKindInteger)
}

// FloatAttribute is an [Attribute] that holds a float value.
type FloatAttribute struct {
	*attribute

	defValue float64
	min      float64
	max      float64
}

// NewFloatAttribute creates a new [FloatAttribute] with the given name, description,
// default value, min, and max.
// It may return an error if the default value is out of the min/max range,
// or if the min value is greater then the max value.
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
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindString, AttributeKindFloat)
}

// ToInteger always returns an error.
func (fa *FloatAttribute) ToInteger() (*IntegerAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindInteger, AttributeKindFloat)
}

// ToFloat returns the [FloatAttribute] itself.
func (fa *FloatAttribute) ToFloat() (*FloatAttribute, error) {
	return fa, nil
}

// ToEnum always returns an error.
func (fa *FloatAttribute) ToEnum() (*EnumAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindEnum, AttributeKindFloat)
}

// EnumAttribute is an [Attribute] that holds an enum as value.
type EnumAttribute struct {
	*attribute

	defValue string
	values   *set[string, int]
}

// NewEnumAttribute creates a new [EnumAttribute] with the given name, description,
// and values. The first value is always selected as the default one.
// It may return an error if no values are passed.
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
		return "", fmt.Errorf("value index cannot be negative")
	}

	if valueIndex >= ea.values.size() {
		return "", fmt.Errorf(`value index "%d" is out of range ("0" - "%d")`, valueIndex, ea.references.size()-1)
	}

	return ea.Values()[valueIndex], nil
}

// ToString always returns an error.
func (ea *EnumAttribute) ToString() (*StringAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindString, AttributeKindEnum)
}

// ToInteger always returns an error.
func (ea *EnumAttribute) ToInteger() (*IntegerAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindInteger, AttributeKindEnum)
}

// ToFloat always returns an error.
func (ea *EnumAttribute) ToFloat() (*FloatAttribute, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the attribute is of kind "%s"`, AttributeKindFloat, AttributeKindEnum)
}

// ToEnum returns the [EnumAttribute] itself.
func (ea *EnumAttribute) ToEnum() (*EnumAttribute, error) {
	return ea, nil
}

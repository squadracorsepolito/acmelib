package acmelib

import (
	"time"

	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

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
	// EntityID returns the entity id of the attribute.
	EntityID() EntityID
	// Name returns the name of the attribute.
	Name() string
	// Desc returns the description of the attribute.
	Desc() string
	// CreateTime returns the time of creation of the attribute.
	CreateTime() time.Time

	// Clone creates a new attribute with the same properties as the current one.
	Clone() (Attribute, error)

	// Type returns the kind of the attribute.
	Type() AttributeType

	addRef(*AttributeAssignment)
	removeRef(EntityID)

	// References returns a slice of references of the attribute.
	References() []*AttributeAssignment

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
	*withRefs[*AttributeAssignment]

	typ AttributeType
}

func newAttributeFromEntity(ent *entity, typ AttributeType) *attribute {
	return &attribute{
		entity:   ent,
		withRefs: newWithRefs[*AttributeAssignment](),

		typ: typ,
	}
}

func newAttribute(name string, typ AttributeType) *attribute {
	return newAttributeFromEntity(newEntity(name, EntityKindAttribute), typ)
}

func (a *attribute) clone() *attribute {
	return newAttributeFromEntity(a.entity.clone(), a.typ)
}

func (a *attribute) errorf(err error) error {
	return &EntityError{
		Kind:     EntityKindAttribute,
		EntityID: a.entityID,
		Name:     a.name,
		Err:      err,
	}
}

func (a *attribute) stringify(s *stringer.Stringer) {
	a.entity.stringify(s)
	s.Write("type: %s\n", a.typ)
	a.withRefs.stringify(s)
}

// Type returns the type of the attribute.
func (a *attribute) Type() AttributeType {
	return a.typ
}

// ToString returns a [ConversionError].
func (a *attribute) ToString() (*StringAttribute, error) {
	return nil, a.errorf(newConversionError(a.typ.String(), AttributeTypeString.String()))
}

// ToInteger returns a [ConversionError].
func (a *attribute) ToInteger() (*IntegerAttribute, error) {
	return nil, a.errorf(newConversionError(a.typ.String(), AttributeTypeInteger.String()))
}

// ToFloat returns a [ConversionError].
func (a *attribute) ToFloat() (*FloatAttribute, error) {
	return nil, a.errorf(newConversionError(a.typ.String(), AttributeTypeFloat.String()))
}

// ToEnum returns a [ConversionError].
func (a *attribute) ToEnum() (*EnumAttribute, error) {
	return nil, a.errorf(newConversionError(a.typ.String(), AttributeTypeEnum.String()))
}

var _ Attribute = (*StringAttribute)(nil)

// StringAttribute is an [Attribute] that holds a string value.
type StringAttribute struct {
	*attribute

	defValue string
}

func newStringAttributeFromBase(base *attribute, defValue string) *StringAttribute {
	return &StringAttribute{
		attribute: base,

		defValue: defValue,
	}
}

// NewStringAttribute creates a new [StringAttribute] with the given name,
// and default value.
func NewStringAttribute(name, defValue string) *StringAttribute {
	return newStringAttributeFromBase(newAttribute(name, AttributeTypeString), defValue)
}

// Clone creates a new [StringAttribute] with the same properties as the current one.
func (sa *StringAttribute) Clone() (Attribute, error) {
	return newStringAttributeFromBase(sa.attribute.clone(), sa.defValue), nil
}

func (sa *StringAttribute) stringify(s *stringer.Stringer) {
	sa.attribute.stringify(s)
	s.Write("default_value: %s\n", sa.defValue)
}

func (sa *StringAttribute) String() string {
	s := stringer.New()
	s.Write("string_attribute:\n")
	sa.stringify(s)
	return s.String()
}

// DefValue returns the default value of the [StringAttribute].
func (sa *StringAttribute) DefValue() string {
	return sa.defValue
}

// ToString returns the [StringAttribute] itself.
func (sa *StringAttribute) ToString() (*StringAttribute, error) {
	return sa, nil
}

// ToAttribute returns the attribute itself.
func (sa *StringAttribute) ToAttribute() (Attribute, error) {
	return sa, nil
}

var _ Attribute = (*IntegerAttribute)(nil)

// IntegerAttribute is an [Attribute] that holds an integer value.
type IntegerAttribute struct {
	*attribute

	defValue int
	min      int
	max      int

	isHexFormat bool
}

func newIntegerAttributeFromBase(base *attribute, defValue, min, max int) (*IntegerAttribute, error) {
	if min > max {
		return nil, newArgError("min", newGreaterError("max"))
	}

	if defValue > max {
		return nil, newArgError("defValue", newGreaterError("max"))
	}

	if defValue < min {
		return nil, newArgError("defValue", newLowerError("min"))
	}

	return &IntegerAttribute{
		attribute: base,

		defValue: defValue,
		min:      min,
		max:      max,

		isHexFormat: false,
	}, nil
}

// NewIntegerAttribute creates a new [IntegerAttribute] with the given name,
// default value, min, and max.
// It may return an error if the default value is out of the min/max range,
// or if the min value is greater then the max value.
func NewIntegerAttribute(name string, defValue, min, max int) (*IntegerAttribute, error) {
	return newIntegerAttributeFromBase(newAttribute(name, AttributeTypeInteger), defValue, min, max)
}

// Clone creates a new [IntegerAttribute] with the same properties as the current one.
func (ia *IntegerAttribute) Clone() (Attribute, error) {
	cloned, err := newIntegerAttributeFromBase(ia.attribute.clone(), ia.defValue, ia.min, ia.max)
	if err != nil {
		return nil, err
	}

	cloned.isHexFormat = ia.isHexFormat

	return cloned, nil
}

func (ia *IntegerAttribute) stringify(s *stringer.Stringer) {
	ia.attribute.stringify(s)

	s.Write("min: %d; max: %d; hex_format: %t\n", ia.min, ia.max, ia.isHexFormat)
	s.Write("default_value: %d\n", ia.defValue)
}

func (ia *IntegerAttribute) String() string {
	s := stringer.New()
	s.Write("integer_attribute:\n")
	ia.stringify(s)
	return s.String()
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

// ToInteger returns the [IntegerAttribute] itself.
func (ia *IntegerAttribute) ToInteger() (*IntegerAttribute, error) {
	return ia, nil
}

// ToAttribute returns the attribute itself.
func (ia *IntegerAttribute) ToAttribute() (Attribute, error) {
	return ia, nil
}

var _ Attribute = (*FloatAttribute)(nil)

// FloatAttribute is an [Attribute] that holds a float value.
type FloatAttribute struct {
	*attribute

	defValue float64
	min      float64
	max      float64
}

func newFloatAttributeFromBase(base *attribute, defValue, min, max float64) (*FloatAttribute, error) {
	if min > max {
		return nil, newArgError("min", newGreaterError("max"))
	}

	if defValue > max {
		return nil, newArgError("defValue", newGreaterError("max"))
	}

	if defValue < min {
		return nil, newArgError("defValue", newLowerError("min"))
	}

	return &FloatAttribute{
		attribute: base,

		defValue: defValue,
		min:      min,
		max:      max,
	}, nil
}

// NewFloatAttribute creates a new [FloatAttribute] with the given name,
// default value, min, and max.
// It may return an error if the default value is out of the min/max range,
// or if the min value is greater then the max value.
func NewFloatAttribute(name string, defValue, min, max float64) (*FloatAttribute, error) {
	return newFloatAttributeFromBase(newAttribute(name, AttributeTypeFloat), defValue, min, max)
}

// Clone creates a new [FloatAttribute] with the same properties as the current one.
func (fa *FloatAttribute) Clone() (Attribute, error) {
	return newFloatAttributeFromBase(fa.attribute.clone(), fa.defValue, fa.min, fa.max)
}

func (fa *FloatAttribute) stringify(s *stringer.Stringer) {
	fa.attribute.stringify(s)

	s.Write("min: %g; max: %g\n", fa.min, fa.max)
	s.Write("default_value: %g\n", fa.defValue)
}

func (fa *FloatAttribute) String() string {
	s := stringer.New()
	s.Write("float_attribute:\n")
	fa.stringify(s)
	return s.String()
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

// ToFloat returns the [FloatAttribute] itself.
func (fa *FloatAttribute) ToFloat() (*FloatAttribute, error) {
	return fa, nil
}

// ToAttribute returns the attribute itself.
func (fa *FloatAttribute) ToAttribute() (Attribute, error) {
	return fa, nil
}

var _ Attribute = (*EnumAttribute)(nil)

// EnumAttribute is an [Attribute] that holds an enum as value.
type EnumAttribute struct {
	*attribute

	defValue string
	values   *collection.Map[string, int]
}

func newEnumAttributeFromBase(base *attribute, values ...string) (*EnumAttribute, error) {
	if len(values) == 0 {
		return nil, newArgError("values", ErrIsNil)
	}

	valSet := collection.NewMap[string, int]()
	currIdx := 0
	for _, val := range values {
		if valSet.Has(val) {
			continue
		}

		valSet.Set(val, currIdx)
		currIdx++
	}

	return &EnumAttribute{
		attribute: base,

		defValue: values[0],
		values:   valSet,
	}, nil
}

// NewEnumAttribute creates a new [EnumAttribute] with the given name and values.
// The first value is always selected as the default one.
// It may return an error if no values are passed.
func NewEnumAttribute(name string, values ...string) (*EnumAttribute, error) {
	return newEnumAttributeFromBase(newAttribute(name, AttributeTypeEnum), values...)
}

// Clone creates a new [EnumAttribute] with the same properties as the current one.
func (ea *EnumAttribute) Clone() (Attribute, error) {
	cloned, err := newEnumAttributeFromBase(ea.attribute.clone(), ea.Values()...)
	if err != nil {
		return nil, err
	}

	cloned.defValue = ea.defValue

	return cloned, nil
}

func (ea *EnumAttribute) stringify(s *stringer.Stringer) {
	ea.attribute.stringify(s)

	if len(ea.Values()) > 0 {
		s.Write("values:\n")
		s.Indent()
		for idx, val := range ea.Values() {
			s.Write("value: %s; index: %d\n", val, idx)
		}
		s.Unindent()
	}

	s.Write("default_value: %s\n", ea.defValue)
}

func (ea *EnumAttribute) String() string {
	s := stringer.New()
	s.Write("enum_attribute:\n")
	ea.stringify(s)
	return s.String()
}

// DefValue returns the default value of the [EnumAttribute].
func (ea *EnumAttribute) DefValue() string {
	return ea.defValue
}

// Values returns the values of the [EnumAttribute] in the order specified in the factory method.
func (ea *EnumAttribute) Values() []string {
	valSlice := make([]string, ea.values.Size())
	for val, valIdx := range ea.values.Entries() {
		valSlice[valIdx] = val
	}
	return valSlice
}

// GetValueAtIndex returns the value at the given index.
// The index refers to the order of the values in the factory method.
// It may return an error if the index is out of range.
func (ea *EnumAttribute) GetValueAtIndex(valueIndex int) (string, error) {
	if valueIndex < 0 {
		return "", ea.errorf(newArgError("valueIndex", ErrIsNegative))
	}

	if valueIndex >= ea.values.Size() {
		return "", ea.errorf(newArgError("valueIndex", ErrOutOfBounds))
	}

	return ea.Values()[valueIndex], nil
}

// ToEnum returns the [EnumAttribute] itself.
func (ea *EnumAttribute) ToEnum() (*EnumAttribute, error) {
	return ea, nil
}

// ToAttribute returns the attribute itself.
func (ea *EnumAttribute) ToAttribute() (Attribute, error) {
	return ea, nil
}

// AttributableEntity represents an entity that can hold attributes.
type AttributableEntity interface {
	errorf(err error) error

	// EntityID returns the unique identifier of the entity.
	EntityID() EntityID
	// EntityKind returns the kind of the entity.
	EntityKind() EntityKind
	// Name returns the name of the entity.
	Name() string

	// AssignAttribute assigns the given attribute/value pair to the entity.
	AssignAttribute(attribute Attribute, value any) error
	// RemoveAttributeAssignment removes the attribute assignment
	// with the given attribute entity id from the entity.
	RemoveAttributeAssignment(attributeEntityID EntityID) error
	// RemoveAllAttributeAssignments removes all the attribute assignments from the entity.
	RemoveAllAttributeAssignments()
	// AttributeAssignments returns a slice of all attribute assignments of the entity.
	AttributeAssignments() []*AttributeAssignment
	// GetAttributeAssignment returns the attribute assignment
	// with the given attribute entity id from the entity.
	GetAttributeAssignment(attributeEntityID EntityID) (*AttributeAssignment, error)
}

// AttributeAssignment represents a link between an [Attribute] and an [AttributableEntity]
// with an assigned value.
type AttributeAssignment struct {
	attribute Attribute
	entity    AttributableEntity
	value     any
}

func newAttributeAssignment(att Attribute, ent AttributableEntity, val any) *AttributeAssignment {
	return &AttributeAssignment{
		attribute: att,
		entity:    ent,
		value:     val,
	}
}

func (aa *AttributeAssignment) stringify(s *stringer.Stringer) {
	s.Write("entity_id: %s; entity_kind: %s; name: %s; value: %v;\n",
		aa.EntityID(), aa.entity.EntityKind(), aa.entity.Name(), aa.value)
}

// EntityID returns the entity id of the [AttributableEntity] of the [AttributeAssignment].
func (aa *AttributeAssignment) EntityID() EntityID {
	return aa.entity.EntityID()
}

// Attribute returns the [Attribute] of the [AttributeAssignment].
func (aa *AttributeAssignment) Attribute() Attribute {
	return aa.attribute
}

// Value returns the value of the [AttributeAssignment].
func (aa *AttributeAssignment) Value() any {
	return aa.value
}

// Entity returns the [AttributableEntity] of the [AttributeAssignment].
func (aa *AttributeAssignment) Entity() AttributableEntity {
	return aa.entity
}

// ToBusEntity returns the [AttributableEntity] as a [Bus].
//
// It returns a [ConversionError] if the kind of the entity is not equal to
// [EntityKindBus].
func (aa *AttributeAssignment) ToBusEntity() (*Bus, error) {
	if aa.entity.EntityKind() == EntityKindBus {
		return aa.entity.(*Bus), nil
	}

	return nil, aa.entity.errorf(&ConversionError{
		From: aa.entity.EntityKind().String(),
		To:   EntityKindBus.String(),
	})
}

// ToNodeEntity returns the [AttributableEntity] as a [Node].
//
// It returns a [ConversionError] if the kind of the entity is not equal to
// [EntityKindNode].
func (aa *AttributeAssignment) ToNodeEntity() (*Node, error) {
	if aa.entity.EntityKind() == EntityKindNode {
		return aa.entity.(*Node), nil
	}

	return nil, aa.entity.errorf(&ConversionError{
		From: aa.entity.EntityKind().String(),
		To:   EntityKindNode.String(),
	})
}

// ToMessageEntity returns the [AttributableEntity] as a [Message].
//
// It returns a [ConversionError] if the kind of the entity is not equal to
// [EntityKindMessage].
func (aa *AttributeAssignment) ToMessageEntity() (*Message, error) {
	if aa.entity.EntityKind() == EntityKindMessage {
		return aa.entity.(*Message), nil
	}

	return nil, aa.entity.errorf(&ConversionError{
		From: aa.entity.EntityKind().String(),
		To:   EntityKindMessage.String(),
	})
}

// ToSignalEntity returns the [AttributableEntity] as a [Signal].
//
// It returns a [ConversionError] if the kind of the entity is not equal to
// [EntityKindSignal].
func (aa *AttributeAssignment) ToSignalEntity() (Signal, error) {
	if aa.entity.EntityKind() == EntityKindSignal {
		return aa.entity.(Signal), nil
	}

	return nil, aa.entity.errorf(&ConversionError{
		From: aa.entity.EntityKind().String(),
		To:   EntityKindSignal.String(),
	})
}

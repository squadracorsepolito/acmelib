package acmelib

import (
	"slices"
	"strings"
	"time"

	"github.com/jaevor/go-nanoid"
	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// EntityKind is the kind of an entity.
type EntityKind int

const (
	// EntityKindNetwork represents a [Network] entity.
	EntityKindNetwork EntityKind = iota
	// EntityKindBus represents a [Bus] entity.
	EntityKindBus
	// EntityKindNode represents a [Node] entity.
	EntityKindNode
	// EntityKindMessage represents a [Message] entity.
	EntityKindMessage
	// EntityKindSignal represents a [Signal] entity.
	EntityKindSignal
	// EntityKindSignalType represents a [SignalType] entity.
	EntityKindSignalType
	// EntityKindSignalUnit represents a [SignalUnit] entity.
	EntityKindSignalUnit
	// EntityKindSignalEnum represents a [SignalEnum] entity.
	EntityKindSignalEnum
	// EntityKindAttribute represents a [Attribute] entity.
	EntityKindAttribute
	// EntityKindCANIDBuilder represents a [CANIDBuilder] entity.
	EntityKindCANIDBuilder
)

func (ek EntityKind) String() string {
	switch ek {
	case EntityKindNetwork:
		return "network"
	case EntityKindBus:
		return "bus"
	case EntityKindNode:
		return "node"
	case EntityKindMessage:
		return "message"
	case EntityKindSignal:
		return "signal"
	case EntityKindSignalType:
		return "signal-type"
	case EntityKindSignalUnit:
		return "signal-unit"
	case EntityKindSignalEnum:
		return "signal-enum"
	case EntityKindAttribute:
		return "attribute"
	case EntityKindCANIDBuilder:
		return "canid-builder"
	default:
		return "unknown"
	}
}

// EntityID is the unique identifier of an entity.
type EntityID string

func newEntityID() EntityID {
	gen, err := nanoid.Standard(21)
	if err != nil {
		panic(err)
	}
	return EntityID(gen())
}

func (id EntityID) String() string {
	return string(id)
}

// Entity interface represents an entity.
type Entity interface {
	// EntityID returns the unique identifier of the entity.
	EntityID() EntityID

	// EntityKind returns the kind of the entity.
	EntityKind() EntityKind

	// Name returns the name of the entity.
	Name() string
	// UpdateName updates the name of the entity.
	UpdateName(newName string) error

	// Desc returns the description of the entity.
	Desc() string
	// SetDesc updates the description of the entity.
	SetDesc(desc string)

	// CreateTime returns the time of creation of the entity.
	CreateTime() time.Time

	ToNetwork() (*Network, error)
	ToBus() (*Bus, error)
	ToNode() (*Node, error)
	ToMessage() (*Message, error)
	ToSignal() (Signal, error)
	ToSignalType() (*SignalType, error)
	ToSignalUnit() (*SignalUnit, error)
	ToSignalEnum() (*SignalEnum, error)
	ToAttribute() (Attribute, error)
	ToCANIDBuilder() (*CANIDBuilder, error)
}

var _ Entity = (*entity)(nil)

type entity struct {
	entityID   EntityID
	entityKind EntityKind
	name       string
	desc       string
	createTime time.Time
}

func newEntity(name string, kind EntityKind) *entity {
	id := newEntityID()
	createTime := time.Now()

	return &entity{
		entityID:   id,
		entityKind: kind,
		name:       name,
		desc:       "",
		createTime: createTime,
	}
}

// EntityID returns the unique identifier of the entity.
func (e *entity) EntityID() EntityID {
	return e.entityID
}

// EntityKind returns the kind of the entity.
func (e *entity) EntityKind() EntityKind {
	return e.entityKind
}

// Name returns the name of the entity.
func (e *entity) Name() string {
	return e.name
}

// UpdateName updates the name of the entity.
func (e *entity) UpdateName(newName string) error {
	if e.name == newName {
		return nil
	}

	e.name = newName
	return nil
}

// Desc returns the description of the entity.
func (e *entity) Desc() string {
	return e.desc
}

// CreateTime returns the time when the entity was created.
func (e *entity) CreateTime() time.Time {
	return e.createTime
}

// SetDesc sets the description of the entity.
func (e *entity) SetDesc(desc string) {
	e.desc = desc
}

func (e *entity) stringify(s *stringer.Stringer) {
	s.Write("entity_id: %s; entity_kind: %s\n", e.entityID, e.entityKind)
	s.Write("name: %s\n", e.name)

	if len(e.desc) > 0 {
		s.Write("desc: %s\n", e.desc)
	}

	s.Write("create_time: %s\n", e.createTime.Format(time.RFC3339))
}

func (e *entity) clone() *entity {
	return &entity{
		entityID:   newEntityID(),
		entityKind: e.entityKind,
		name:       e.name,
		desc:       e.desc,
		createTime: e.createTime,
	}
}

func (e *entity) ToNetwork() (*Network, error) {
	return nil, newConversionError(e.entityKind.String(), EntityKindNetwork.String())
}

func (e *entity) ToBus() (*Bus, error) {
	return nil, newConversionError(e.entityKind.String(), EntityKindBus.String())
}

func (e *entity) ToNode() (*Node, error) {
	return nil, newConversionError(e.entityKind.String(), EntityKindNode.String())
}

func (e *entity) ToMessage() (*Message, error) {
	return nil, newConversionError(e.entityKind.String(), EntityKindMessage.String())
}

func (e *entity) ToSignal() (Signal, error) {
	return nil, newConversionError(e.entityKind.String(), EntityKindSignal.String())
}

func (e *entity) ToSignalType() (*SignalType, error) {
	return nil, newConversionError(e.entityKind.String(), EntityKindSignalType.String())
}

func (e *entity) ToSignalUnit() (*SignalUnit, error) {
	return nil, newConversionError(e.entityKind.String(), EntityKindSignalUnit.String())
}

func (e *entity) ToSignalEnum() (*SignalEnum, error) {
	return nil, newConversionError(e.entityKind.String(), EntityKindSignalEnum.String())
}

func (e *entity) ToAttribute() (Attribute, error) {
	return nil, newConversionError(e.entityKind.String(), EntityKindAttribute.String())
}

func (e *entity) ToCANIDBuilder() (*CANIDBuilder, error) {
	return nil, newConversionError(e.entityKind.String(), EntityKindCANIDBuilder.String())
}

type withAttributes struct {
	attAssignments *collection.Map[EntityID, *AttributeAssignment]
}

func newWithAttributes() *withAttributes {
	return &withAttributes{
		attAssignments: collection.NewMap[EntityID, *AttributeAssignment](),
	}
}

func (wa *withAttributes) stringify(s *stringer.Stringer) {
	if wa.attAssignments.Size() > 0 {
		s.Write("attribute_assignments:\n")
		s.Indent()
		for _, attAss := range wa.AttributeAssignments() {
			attAss.stringify(s)
		}
		s.Unindent()
	}
}

func (wa *withAttributes) addAttributeAssignment(attribute Attribute, ent AttributableEntity, val any) error {
	if attribute == nil {
		return &ArgError{
			Name: "attribute",
			Err:  ErrIsNil,
		}
	}

	switch v := val.(type) {
	case int:
		if attribute.Type() != AttributeTypeInteger {
			return &AttributeValueError{Err: ErrInvalidType}
		}

		intAtt, err := attribute.ToInteger()
		if err != nil {
			panic(err)
		}
		if v < intAtt.min || v > intAtt.max {
			return &AttributeValueError{Err: ErrOutOfBounds}
		}

	case float64:
		if attribute.Type() != AttributeTypeFloat {
			return &AttributeValueError{Err: ErrInvalidType}
		}

		floatAtt, err := attribute.ToFloat()
		if err != nil {
			panic(err)
		}
		if v < floatAtt.min || v > floatAtt.max {
			return &AttributeValueError{Err: ErrOutOfBounds}
		}

	case string:
		switch attribute.Type() {
		case AttributeTypeString:
		case AttributeTypeEnum:
			enumAtt, err := attribute.ToEnum()
			if err != nil {
				panic(err)
			}
			if !enumAtt.values.Has(v) {
				return &AttributeValueError{Err: ErrNotFound}
			}

		default:
			return &AttributeValueError{Err: ErrInvalidType}
		}

	default:
		return &AttributeValueError{Err: ErrInvalidType}
	}

	attAss := newAttributeAssignment(attribute, ent, val)

	wa.attAssignments.Set(attribute.EntityID(), attAss)
	attribute.addRef(attAss)

	return nil
}

func (wa *withAttributes) removeAttributeAssignment(attEntID EntityID) error {
	attAss, ok := wa.attAssignments.Get(attEntID)
	if !ok {
		return ErrNotFound
	}

	wa.attAssignments.Delete(attEntID)
	attAss.attribute.removeRef(attAss.EntityID())

	return nil
}

// RemoveAllAttributeAssignments removes all the attribute assignments from the entity.
func (wa *withAttributes) RemoveAllAttributeAssignments() {
	for attVal := range wa.attAssignments.Values() {
		attVal.attribute.removeRef(attVal.EntityID())
	}
	wa.attAssignments.Clear()
}

// AttributeAssignments returns a slice of all attribute assignments of the entity.
func (wa *withAttributes) AttributeAssignments() []*AttributeAssignment {
	attSlice := slices.Collect(wa.attAssignments.Values())
	slices.SortFunc(attSlice, func(a, b *AttributeAssignment) int {
		return strings.Compare(a.attribute.Name(), b.attribute.Name())
	})
	return attSlice
}

func (wa *withAttributes) getAttributeAssignment(attributeEntityID EntityID) (*AttributeAssignment, error) {
	attVal, ok := wa.attAssignments.Get(attributeEntityID)
	if !ok {
		return nil, ErrNotFound
	}
	return attVal, nil
}

type referenceableEntity interface {
	EntityID() EntityID
}

type withRefs[R referenceableEntity] struct {
	refs *collection.Map[EntityID, R]
}

func newWithRefs[R referenceableEntity]() *withRefs[R] {
	return &withRefs[R]{
		refs: collection.NewMap[EntityID, R](),
	}
}

func (t *withRefs[R]) stringify(s *stringer.Stringer) {
	refCount := t.ReferenceCount()
	if refCount > 0 {
		s.Write("reference_count: %d\n", refCount)
	}
}

func (t *withRefs[R]) addRef(ref R) {
	t.refs.Set(ref.EntityID(), ref)
}

func (t *withRefs[R]) removeRef(refID EntityID) {
	t.refs.Delete(refID)
}

func (t *withRefs[R]) ReferenceCount() int {
	return t.refs.Size()
}

func (t *withRefs[R]) References() []R {
	return slices.Collect(t.refs.Values())
}

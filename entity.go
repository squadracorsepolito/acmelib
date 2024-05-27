package acmelib

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/jaevor/go-nanoid"
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
	// EntityKindNodeInterface represents a [NodeInterface] entity.
	EntityKindNodeInterface
	// EntityKindMessage represents a [Message] entity.
	EntityKindMessage
	// EntityKindSignal represents a [Signal] entity.
	EntityKindSignal
	// EntityKindSignalType represents a [SignalType] entity.
	EntityKindSignalType
	// EntityKindSignalEnum represents a [SignalEnum] entity.
	EntityKindSignalEnum
	// EntityKindSignalEnumValue represents a [SignalEnumValue] entity.
	EntityKindSignalEnumValue
	// EntityKindSignalUnit represents a [SignalUnit] entity.
	EntityKindSignalUnit
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
	case EntityKindNodeInterface:
		return "node-interface"
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
	case EntityKindSignalEnumValue:
		return "signal-enum-value"
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

func (e *entity) stringify(b *strings.Builder, tabs int) {
	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%sentity_id: %s; entity_kind: %s\n", tabStr, e.entityID, e.entityKind))
	b.WriteString(fmt.Sprintf("%sname: %s\n", tabStr, e.name))

	if len(e.desc) > 0 {
		b.WriteString(fmt.Sprintf("%sdesc: %s\n", tabStr, e.desc))
	}

	b.WriteString(fmt.Sprintf("%screate_time: %s\n", tabStr, e.createTime.Format(time.RFC3339)))
}

type withAttributes struct {
	attributes *set[EntityID, *AttributeAssignment]
}

func newWithAttributes() *withAttributes {
	return &withAttributes{
		attributes: newSet[EntityID, *AttributeAssignment](),
	}
}

func (wa *withAttributes) addAttributeAssignment(attribute Attribute, ent AttributableEntity, val any) error {
	if attribute == nil {
		return &ArgumentError{
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
			if !enumAtt.values.hasKey(v) {
				return &AttributeValueError{Err: ErrNotFound}
			}

		default:
			return &AttributeValueError{Err: ErrInvalidType}
		}

	default:
		return &AttributeValueError{Err: ErrInvalidType}
	}

	attAss := newAttributeAssignment(attribute, ent, val)

	wa.attributes.add(attribute.EntityID(), attAss)
	attribute.addRef(attAss)

	return nil
}

func (wa *withAttributes) removeAttributeAssignment(attEntID EntityID) error {
	attAss, err := wa.attributes.getValue(attEntID)
	if err != nil {
		return err
	}

	wa.attributes.remove(attEntID)
	attAss.attribute.removeRef(attAss.EntityID())

	return nil
}

func (wa *withAttributes) RemoveAllAttributeAssignments() {
	for _, attVal := range wa.attributes.entries() {
		attVal.attribute.removeRef(attVal.EntityID())
	}
	wa.attributes.clear()
}

func (wa *withAttributes) AttributeAssignments() []*AttributeAssignment {
	attSlice := wa.attributes.getValues()
	slices.SortFunc(attSlice, func(a, b *AttributeAssignment) int {
		return strings.Compare(a.attribute.Name(), b.attribute.Name())
	})
	return attSlice
}

func (wa *withAttributes) getAttributeAssignment(attributeEntityID EntityID) (*AttributeAssignment, error) {
	attVal, err := wa.attributes.getValue(attributeEntityID)
	if err != nil {
		return nil, &GetEntityError{
			EntityID: attributeEntityID,
			Err:      err,
		}
	}
	return attVal, nil
}

type referenceableEntity interface {
	EntityID() EntityID
}

type withRefs[R referenceableEntity] struct {
	refs *set[EntityID, R]
}

func newWithRefs[R referenceableEntity]() *withRefs[R] {
	return &withRefs[R]{
		refs: newSet[EntityID, R](),
	}
}

func (t *withRefs[R]) addRef(ref R) {
	t.refs.add(ref.EntityID(), ref)
}

func (t *withRefs[R]) removeRef(refID EntityID) {
	t.refs.remove(refID)
}

func (t *withRefs[R]) ReferenceCount() int {
	return t.refs.size()
}

func (t *withRefs[R]) References() []R {
	return t.refs.getValues()
}

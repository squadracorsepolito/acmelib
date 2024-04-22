package acmelib

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/jaevor/go-nanoid"
)

// EntityID is the unique identifier of an entity.
// Entities are:
//   - networks
//   - buses
//   - nodes
//   - messages
//   - signals
//   - signal types
//   - signal enums
//   - signal enum values
//   - signal units
//   - attributes
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
	name       string
	desc       string
	createTime time.Time
}

func newEntity(name string) *entity {
	id := newEntityID()
	createTime := time.Now()

	return &entity{
		entityID:   id,
		name:       name,
		desc:       "",
		createTime: createTime,
	}
}

// EntityID returns the unique identifier of the entity.
func (e *entity) EntityID() EntityID {
	return e.entityID
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

	b.WriteString(fmt.Sprintf("%sentity_id: %s\n", tabStr, e.entityID.String()))
	b.WriteString(fmt.Sprintf("%sname: %s\n", tabStr, e.name))

	if len(e.desc) > 0 {
		b.WriteString(fmt.Sprintf("%sdesc: %s\n", tabStr, e.desc))
	}

	b.WriteString(fmt.Sprintf("%screate_time: %s\n", tabStr, e.createTime.Format(time.RFC3339)))
}

type attributeEntity struct {
	*entity

	attributeValues *set[EntityID, *AttributeValue]
	attRefKind      AttributeRefKind
}

func newAttributeEntity(name string, attRefKind AttributeRefKind) *attributeEntity {
	return &attributeEntity{
		entity: newEntity(name),

		attributeValues: newSet[EntityID, *AttributeValue](),
		attRefKind:      attRefKind,
	}
}

func (ae *attributeEntity) stringify(b *strings.Builder, tabs int) {
	ae.entity.stringify(b, tabs)

	if ae.attributeValues.size() == 0 {
		return
	}

	tabStr := getTabString(tabs)
	b.WriteString(fmt.Sprintf("%sattribute values:\n", tabStr))
	for _, attVal := range ae.AttributeValues() {
		attVal.stringify(b, tabs+1)
	}
}

// AddAttributeValue adds an [Attribute] to the entity and it assign
// the given value to it.
// It may return an error if the given value is not valid for the given
// [Attribute].
func (ae *attributeEntity) AddAttributeValue(attribute Attribute, value any) error {
	switch v := value.(type) {
	case int:
		if attribute.Kind() != AttributeKindInteger {
			return &ArgumentError{
				Name: "value",
				Err:  ErrInvalidType,
			}
		}
		intAtt, err := attribute.ToInteger()
		if err != nil {
			panic(err)
		}
		if v < intAtt.min || v > intAtt.max {
			return &ArgumentError{
				Name: "value",
				Err:  ErrOutOfBounds,
			}
		}

	case float64:
		if attribute.Kind() != AttributeKindFloat {
			return &ArgumentError{
				Name: "value",
				Err:  ErrInvalidType,
			}
		}
		floatAtt, err := attribute.ToFloat()
		if err != nil {
			panic(err)
		}
		if v < floatAtt.min || v > floatAtt.max {
			return &ArgumentError{
				Name: "value",
				Err:  ErrOutOfBounds,
			}
		}

	case string:
		switch attribute.Kind() {
		case AttributeKindString:
		case AttributeKindEnum:
			enumAtt, err := attribute.ToEnum()
			if err != nil {
				panic(err)
			}
			if !enumAtt.values.hasKey(v) {
				return &ArgumentError{
					Name: "value",
					Err:  ErrNotFound,
				}
			}

		default:
			return &ArgumentError{
				Name: "value",
				Err:  ErrInvalidType,
			}
		}

	default:
		return &ArgumentError{
			Name: "value",
			Err:  ErrInvalidType,
		}
	}

	ae.attributeValues.add(attribute.EntityID(), newAttributeValue(attribute, value))
	attribute.addReference(newAttributeRef(ae.entityID, ae.attRefKind, value))

	return nil
}

// RemoveAttributeValue removes an [Attribute] with the given entity id from the entity.
// It also removes the reference to the entity from the attribute.
// It may return an error if the attribute with the given entity id does not exist
// in the entity.
func (ae *attributeEntity) RemoveAttributeValue(attributeEntityID EntityID) error {
	att, err := ae.attributeValues.getValue(attributeEntityID)
	if err != nil {
		return &RemoveEntityError{
			EntityID: attributeEntityID,
			Err:      err,
		}
	}

	ae.attributeValues.remove(attributeEntityID)
	att.attribute.removeReference(ae.entityID)

	return nil
}

// RemoveAllAttributeValues removes all [Attributes] from the entity.
func (ae *attributeEntity) RemoveAllAttributeValues() {
	for _, attVal := range ae.attributeValues.entries() {
		attVal.attribute.removeReference(ae.entityID)
	}

	ae.attributeValues.clear()
}

// AttributeValues returns slice of all the attributes of the entity.
func (ae *attributeEntity) AttributeValues() []*AttributeValue {
	attValSlice := ae.attributeValues.getValues()
	slices.SortFunc(attValSlice, func(a, b *AttributeValue) int {
		return strings.Compare(a.attribute.Name(), b.attribute.Name())
	})
	return attValSlice
}

// GetAttributeValue returns the [Attribute] with the given entity id from the entity.
// It may return an error if the attribute with the given entity id does not exist
// in the entity.
func (ae *attributeEntity) GetAttributeValue(attributeEntityID EntityID) (*AttributeValue, error) {
	attVal, err := ae.attributeValues.getValue(attributeEntityID)
	if err != nil {
		return nil, &GetEntityError{
			EntityID: attributeEntityID,
			Err:      err,
		}
	}
	return attVal, nil
}

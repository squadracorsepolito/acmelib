package acmelib

import (
	"fmt"
	"strings"
	"time"

	"github.com/jaevor/go-nanoid"
)

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

func newEntity(name, desc string) *entity {
	id := newEntityID()
	createTime := time.Now()

	return &entity{
		entityID:   id,
		name:       name,
		desc:       desc,
		createTime: createTime,
	}
}

func (e *entity) EntityID() EntityID {
	return e.entityID
}

func (e *entity) Name() string {
	return e.name
}

func (e *entity) Desc() string {
	return e.desc
}

func (e *entity) CreateTime() time.Time {
	return e.createTime
}

func (e *entity) UpdateDesc(desc string) error {
	if e.desc != desc {
		e.desc = desc
	}

	return nil
}

func (e *entity) UpdateName(name string) error {
	if e.name != name {
		e.name = name
	}

	return nil
}

func (e *entity) toString() string {
	var builder strings.Builder

	builder.WriteString("entity_id: " + e.entityID.String() + "\n")
	builder.WriteString("name: " + e.name + "\n")

	if len(e.desc) > 0 {
		builder.WriteString(fmt.Sprintf("desc: %s\n", e.desc))
	}

	builder.WriteString(fmt.Sprintf("create_time: %s\n", e.createTime.Format(time.RFC3339)))

	return builder.String()
}

type entityWithAttributes struct {
	*entity

	attributeValues *set[EntityID, *AttributeValue]
	attRefKind      AttributeReferenceKind
}

func newEntityWithAttributes(name, desc string, attRefKind AttributeReferenceKind) *entityWithAttributes {
	return &entityWithAttributes{
		entity: newEntity(name, desc),

		attributeValues: newSet[EntityID, *AttributeValue]("attribute value"),
		attRefKind:      attRefKind,
	}
}

func (ewa *entityWithAttributes) AddAttributeValue(attribute Attribute, value any) error {
	switch v := value.(type) {
	case int:
		if attribute.Kind() != AttributeKindInteger {
			return fmt.Errorf(`cannot assign an int value to attribute "%s" of type "%s"`, attribute.Name(), attribute.Kind())
		}
		intAtt, err := attribute.ToInteger()
		if err != nil {
			return fmt.Errorf(`cannot assign value "%d" : %w`, v, err)
		}
		if v < intAtt.min || v > intAtt.max {
			return fmt.Errorf(`cannot assign value "%d" because it is out of min/max range ("%d" - "%d")`, v, intAtt.min, intAtt.max)
		}

	case float64:
		if attribute.Kind() != AttributeKindFloat {
			return fmt.Errorf(`cannot assign a float64 value to attribute "%s" of type "%s"`, attribute.Name(), attribute.Kind())
		}
		floatAtt, err := attribute.ToFloat()
		if err != nil {
			return fmt.Errorf(`cannot assign value "%f" : %w`, v, err)
		}
		if v < floatAtt.min || v > floatAtt.max {
			return fmt.Errorf(`cannot assign value "%f" because it is out of min/max range ("%f" - "%f")`, v, floatAtt.min, floatAtt.max)
		}

	case string:
		switch attribute.Kind() {
		case AttributeKindString:
		case AttributeKindEnum:
			enumAtt, err := attribute.ToEnum()
			if err != nil {
				return fmt.Errorf(`cannot assign value "%s" : %w`, v, err)
			}
			if !enumAtt.values.hasKey(v) {
				return fmt.Errorf(`cannot assign value "%s" beacuse it is not present in the enum`, v)
			}

		default:
			return fmt.Errorf(`cannot assign a string value to attribute "%s" of type "%s"`, attribute.Name(), attribute.Kind())
		}
	}

	ewa.attributeValues.add(attribute.EntityID(), newAttributeValue(attribute, value))
	attribute.addReference(newAttributeReference(ewa.entityID, ewa.attRefKind, value))

	return nil
}

func (ewa *entityWithAttributes) RemoveAttributeValue(attributeEntityID EntityID) error {
	att, err := ewa.attributeValues.getValue(attributeEntityID)
	if err != nil {
		return fmt.Errorf(`cannot remove attribute with entity id "%s" : %w`, attributeEntityID, err)
	}

	ewa.attributeValues.remove(attributeEntityID)
	att.attribute.removeReference(ewa.entityID)

	return nil
}

func (ewa *entityWithAttributes) RemoveAllAttributeValues() {
	for _, attVal := range ewa.attributeValues.entries() {
		attVal.attribute.removeReference(ewa.entityID)
	}

	ewa.attributeValues.clear()
}

func (ewa *entityWithAttributes) AttributeValues() []*AttributeValue {
	return ewa.attributeValues.getValues()
}

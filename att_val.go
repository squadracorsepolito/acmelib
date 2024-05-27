package acmelib

type AttributableEntity interface {
	EntityID() EntityID
	Name() string

	AssignAttribute(attribute Attribute, value any) error
}

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

func (aa *AttributeAssignment) EntityID() EntityID {
	return aa.entity.EntityID()
}

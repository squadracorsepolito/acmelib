package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_attributeEntity_AddAttributeValue(t *testing.T) {
	assert := assert.New(t)

	dummyEnt := NewBus("dummy_bus")
	e := newWithAttributes()

	intAtt, err := NewIntegerAttribute("int_att", 0, 0, 100)
	assert.NoError(err)

	// should not return an error
	assert.NoError(e.addAttributeAssignment(intAtt, dummyEnt, 10))

	// should return an error beacause the value is not an integer
	assert.Error(e.addAttributeAssignment(intAtt, dummyEnt, "string"))
	assert.Error(e.addAttributeAssignment(intAtt, dummyEnt, 10.0))

	// should return an error beacause the value is out of range
	assert.Error(e.addAttributeAssignment(intAtt, dummyEnt, 1000))

	floatAtt, err := NewFloatAttribute("float_att", 0, 0, 100)
	assert.NoError(err)

	// should not return an error
	assert.NoError(e.addAttributeAssignment(floatAtt, dummyEnt, 10.0))

	// should return an error beacause the value is not a float
	assert.Error(e.addAttributeAssignment(floatAtt, dummyEnt, "string"))
	assert.Error(e.addAttributeAssignment(floatAtt, dummyEnt, 10))

	// should return an error beacause the value is out of range
	assert.Error(e.addAttributeAssignment(floatAtt, dummyEnt, 1000.0))

	strAtt := NewStringAttribute("str_att", "")

	// should not return an error
	assert.NoError(e.addAttributeAssignment(strAtt, dummyEnt, "string"))

	// should return an error beacause the value is not a string
	assert.Error(e.addAttributeAssignment(strAtt, dummyEnt, 10))
	assert.Error(e.addAttributeAssignment(strAtt, dummyEnt, 10.0))

	enumAtt, err := NewEnumAttribute("enum_att", "", "val_0", "val_1", "val_2")
	assert.NoError(err)

	// should not return an error
	assert.NoError(e.addAttributeAssignment(enumAtt, dummyEnt, "val_1"))

	// should return an error beacause the value is not a string
	assert.Error(e.addAttributeAssignment(enumAtt, dummyEnt, 10))
	assert.Error(e.addAttributeAssignment(enumAtt, dummyEnt, 10.0))

	// should return an error beacause the value is not present in the enum
	assert.Error(e.addAttributeAssignment(enumAtt, dummyEnt, "val_3"))

	expectedNames := []string{"enum_att", "float_att", "int_att", "str_att"}
	expectedValues := []any{"val_1", 10.0, 10, "string"}
	for idx, attVal := range e.AttributeAssignments() {
		assert.Equal(expectedNames[idx], attVal.Attribute().Name())
		assert.Equal(dummyEnt.EntityID(), attVal.Attribute().References()[0].EntityID())
		assert.Equal(expectedValues[idx], attVal.Value())
	}
}

func Test_attributeEntity_RemoveAttributeValue(t *testing.T) {
	assert := assert.New(t)

	dummyEnt := NewBus("dummy_bus")
	e := newWithAttributes()

	intAtt0, err := NewIntegerAttribute("int_att_0", 0, 0, 100)
	assert.NoError(err)
	intAtt1, err := NewIntegerAttribute("int_att_1", 0, 0, 100)
	assert.NoError(err)
	intAtt2, err := NewIntegerAttribute("int_att_2", 0, 0, 100)
	assert.NoError(err)

	assert.NoError(e.addAttributeAssignment(intAtt0, dummyEnt, 10))
	assert.NoError(e.addAttributeAssignment(intAtt1, dummyEnt, 10))
	assert.NoError(e.addAttributeAssignment(intAtt2, dummyEnt, 10))

	assert.NoError(e.removeAttributeAssignment(intAtt1.EntityID()))

	expectedNames := []string{"int_att_0", "int_att_2"}
	expectedValues := []any{10, 10}
	for idx, attVal := range e.AttributeAssignments() {
		assert.Equal(expectedNames[idx], attVal.Attribute().Name())
		assert.Equal(dummyEnt.EntityID(), attVal.Attribute().References()[0].EntityID())
		assert.Equal(expectedValues[idx], attVal.Attribute().References()[0].Value())
	}

	assert.Error(e.removeAttributeAssignment("dummy-id"))
}

func Test_attributeEntity_RemoveAllAttributeValues(t *testing.T) {
	assert := assert.New(t)

	dummyEnt := NewBus("dummy_bus")
	e := newWithAttributes()

	intAtt0, err := NewIntegerAttribute("int_att_0", 0, 0, 100)
	assert.NoError(err)
	intAtt1, err := NewIntegerAttribute("int_att_1", 0, 0, 100)
	assert.NoError(err)
	intAtt2, err := NewIntegerAttribute("int_att_2", 0, 0, 100)
	assert.NoError(err)

	assert.NoError(e.addAttributeAssignment(intAtt0, dummyEnt, 10))
	assert.NoError(e.addAttributeAssignment(intAtt1, dummyEnt, 10))
	assert.NoError(e.addAttributeAssignment(intAtt2, dummyEnt, 10))

	e.RemoveAllAttributeAssignments()

	assert.Equal(0, len(e.AttributeAssignments()))
	assert.Equal(0, len(intAtt0.References()))
	assert.Equal(0, len(intAtt1.References()))
	assert.Equal(0, len(intAtt2.References()))
}

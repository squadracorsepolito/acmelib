package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringAttribute(t *testing.T) {
	assert := assert.New(t)

	strAtt := NewStringAttribute("str_att", "def_val")
	assert.Equal(AttributeTypeString, strAtt.Type())
	assert.Equal("str_att", strAtt.Name())
	assert.Equal("def_val", strAtt.DefValue())

	_, err := strAtt.ToString()
	assert.NoError(err)

	_, err = strAtt.ToInteger()
	assert.Error(err)
	_, err = strAtt.ToFloat()
	assert.Error(err)
	_, err = strAtt.ToEnum()
	assert.Error(err)
}

func Test_IntegerAttribute(t *testing.T) {
	assert := assert.New(t)

	intAtt, err := NewIntegerAttribute("int_att", 10, 0, 100)
	assert.NoError(err)
	assert.Equal(AttributeTypeInteger, intAtt.Type())
	assert.Equal("int_att", intAtt.Name())
	assert.Equal(10, intAtt.DefValue())
	assert.Equal(0, intAtt.Min())
	assert.Equal(100, intAtt.Max())

	// should return an error because min is greater then max
	_, err = NewIntegerAttribute("int_att", 0, 100, 0)
	assert.Error(err)

	// should return an error because the default value is out of range
	_, err = NewIntegerAttribute("int_att", 1000, 0, 100)
	assert.Error(err)

	_, err = intAtt.ToInteger()
	assert.NoError(err)

	_, err = intAtt.ToString()
	assert.Error(err)
	_, err = intAtt.ToFloat()
	assert.Error(err)
	_, err = intAtt.ToEnum()
	assert.Error(err)
}

func Test_FloatAttribute(t *testing.T) {
	assert := assert.New(t)

	floatAtt, err := NewFloatAttribute("float_att", 10, 0, 100)
	assert.NoError(err)
	assert.Equal(AttributeTypeFloat, floatAtt.Type())
	assert.Equal("float_att", floatAtt.Name())
	assert.Equal(10.0, floatAtt.DefValue())
	assert.Equal(0.0, floatAtt.Min())
	assert.Equal(100.0, floatAtt.Max())

	// should return an error because min is greater then max
	_, err = NewFloatAttribute("float_att", 0, 100, 0)
	assert.Error(err)

	// should return an error because the default value is out of range
	_, err = NewFloatAttribute("float_att", 1000, 0, 100)
	assert.Error(err)

	_, err = floatAtt.ToFloat()
	assert.NoError(err)

	_, err = floatAtt.ToString()
	assert.Error(err)
	_, err = floatAtt.ToInteger()
	assert.Error(err)
	_, err = floatAtt.ToEnum()
	assert.Error(err)
}

func Test_EnumAttribute(t *testing.T) {
	assert := assert.New(t)

	enumAtt, err := NewEnumAttribute("enum_att", "val_0", "val_1", "val_2")
	assert.NoError(err)
	assert.Equal(AttributeTypeEnum, enumAtt.Type())
	assert.Equal("enum_att", enumAtt.Name())
	assert.Equal("val_0", enumAtt.DefValue())
	assert.Equal(3, len(enumAtt.Values()))
	expectedValues := []string{"val_0", "val_1", "val_2"}
	for idx, val := range enumAtt.Values() {
		assert.Equal(expectedValues[idx], val)
	}

	// should return an error because there are no values defined
	_, err = NewEnumAttribute("enum_att")
	assert.Error(err)

	// should compact the values because val_1 is duplicated
	compEnumAtt, err := NewEnumAttribute("enum_att", "val_0", "val_1", "val_2", "val_1")
	assert.NoError(err)
	assert.Equal(3, len(compEnumAtt.Values()))

	val, err := enumAtt.GetValueAtIndex(2)
	assert.NoError(err)
	assert.Equal("val_2", val)

	// should return an error because index is negative
	_, err = enumAtt.GetValueAtIndex(-1)
	assert.Error(err)

	// should return an error because index is out of range
	_, err = enumAtt.GetValueAtIndex(5)
	assert.Error(err)

	_, err = enumAtt.ToEnum()
	assert.NoError(err)

	_, err = enumAtt.ToString()
	assert.Error(err)
	_, err = enumAtt.ToInteger()
	assert.Error(err)
	_, err = enumAtt.ToFloat()
	assert.Error(err)
}

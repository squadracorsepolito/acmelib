package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SignalEnum_AddDeleteValue(t *testing.T) {
	assert := assert.New(t)

	tdEnumMsg := initEnumMessage(assert)

	enum4 := tdEnumMsg.signals.with4Values.enum
	val, err := enum4.AddValue(4, "enum_value_4")
	assert.NoError(err)
	assert.Equal(4, val.Index())
	assert.Equal("enum_value_4", val.Name())
	assert.Equal(4, enum4.MaxIndex())
	assert.Equal(3, enum4.Size())
	enum4.DeleteValue(4)
	assert.Equal(3, enum4.MaxIndex())
	assert.Equal(2, enum4.Size())

	_, err = enum4.AddValue(256, "invalid")
	assert.Error(err)
	_, err = enum4.AddValue(-1, "invalid")
	assert.Error(err)

	enumFixed := tdEnumMsg.signals.fixedSize.enum
	assert.Equal(8, enumFixed.Size())
	assert.Equal(127, enumFixed.MaxIndex())

	_, err = enumFixed.AddValue(0, "enum_value_0")
	assert.Error(err)
	_, err = enumFixed.AddValue(1, "duplicated_name")
	assert.NoError(err)
	_, err = enumFixed.AddValue(2, "duplicated_name")
	assert.NoError(err)
	enumFixed.DeleteValue(1)
	enumFixed.DeleteValue(2)

	_, err = enumFixed.AddValue(256, "invalid")
	assert.Error(err)
}

func Test_SignalEnum_UpdateSize(t *testing.T) {
	assert := assert.New(t)

	tdEnumMsg := initEnumMessage(assert)

	enumFixed := tdEnumMsg.signals.fixedSize.enum
	enumFixed.SetFixedSize(false)
	assert.Equal(7, enumFixed.Size())

	enumFixed.SetFixedSize(true)

	assert.NoError(enumFixed.UpdateSize(16))
	assert.Equal(16, enumFixed.Size())
	assert.Equal(127, enumFixed.MaxIndex())

	assert.NoError(enumFixed.UpdateSize(8))
	assert.Equal(8, enumFixed.Size())

	assert.Error(enumFixed.UpdateSize(17))
}

func Test_SignalEnumValue_UpdateIndex(t *testing.T) {
	assert := assert.New(t)

	tdEnumMsg := initEnumMessage(assert)

	enum4 := tdEnumMsg.signals.with4Values.enum
	val0 := enum4.GetValue0(0)
	assert.NoError(val0.UpdateIndex(255))
	assert.Equal(255, val0.Index())
	assert.Equal(255, enum4.MaxIndex())
	assert.Equal(8, enum4.Size())
	assert.NoError(val0.UpdateIndex(0))

	assert.Error(val0.UpdateIndex(-1))
	assert.Error(val0.UpdateIndex(1))
	assert.Error(val0.UpdateIndex(256))
}

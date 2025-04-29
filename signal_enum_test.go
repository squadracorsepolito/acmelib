package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SignalEnum_AddDeleteValue(t *testing.T) {
	assert := assert.New(t)

	tdEnumMsg := initEnumMessage(assert)

	enum4 := tdEnumMsg.signals.with4Values.enum
	val, err := enum4.AddValue0(4, "enum_value_4")
	assert.NoError(err)
	assert.Equal(4, val.Index())
	assert.Equal("enum_value_4", val.Name())
	assert.Equal(4, enum4.MaxIndex())
	assert.Equal(3, enum4.Size())
	enum4.DeleteValue(4)
	assert.Equal(3, enum4.MaxIndex())
	assert.Equal(2, enum4.Size())

	_, err = enum4.AddValue0(256, "invalid")
	assert.Error(err)
	_, err = enum4.AddValue0(-1, "invalid")
	assert.Error(err)

	enumFixed := tdEnumMsg.signals.fixedSize.enum
	assert.Equal(8, enumFixed.Size())
	assert.Equal(127, enumFixed.MaxIndex())

	_, err = enumFixed.AddValue0(0, "enum_value_0")
	assert.Error(err)
	_, err = enumFixed.AddValue0(1, "duplicated_name")
	assert.NoError(err)
	_, err = enumFixed.AddValue0(2, "duplicated_name")
	assert.NoError(err)
	enumFixed.DeleteValue(1)
	enumFixed.DeleteValue(2)

	_, err = enumFixed.AddValue0(256, "invalid")
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

// func Test_SignalEnum(t *testing.T) {
// 	assert := assert.New(t)

// 	enum := NewSignalEnum("sig_enum")
// 	assert.Equal("sig_enum", enum.Name())

// 	enum.SetMinSize(8)
// 	assert.Equal(8, enum.GetSize())
// }

// func Test_SignalEnum_AddValue(t *testing.T) {
// 	assert := assert.New(t)

// 	msg := NewMessage("msg", 1, 1)

// 	enum := NewSignalEnum("enum")

// 	sig, err := NewEnumSignal("sig", enum)
// 	assert.NoError(err)
// 	assert.NoError(msg.AppendSignal(sig))

// 	enumVal0 := NewSignalEnumValue("enum_val_0", 0)
// 	enumVal1 := NewSignalEnumValue("enum_val_1", 1)
// 	enumVal2 := NewSignalEnumValue("enum_val_2", 255)

// 	// should insert enumVal0, enumVal1, and enumVal2 without returning errors
// 	assert.NoError(enum.AddValue(enumVal0))
// 	assert.NoError(enum.AddValue(enumVal1))
// 	assert.NoError(enum.AddValue(enumVal2))

// 	assert.Equal(255, enum.MaxIndex())

// 	// should return an error because index 256 cannot fit in 8 bits
// 	enumVal3 := NewSignalEnumValue("enum_val_3", 256)
// 	assert.Error(enum.AddValue(enumVal3))

// 	// should return an error because enumVal4 has a duplicated name
// 	enumVal4 := NewSignalEnumValue("enum_val_0", 2)
// 	assert.Error(enum.AddValue(enumVal4))

// 	// should return an error because enumVal5 has a duplicated index
// 	enumVal5 := NewSignalEnumValue("enum_val_5", 0)
// 	assert.Error(enum.AddValue(enumVal5))
// }

// func Test_SignalEnum_RemoveValue(t *testing.T) {
// 	assert := assert.New(t)

// 	msg := NewMessage("msg", 1, 1)

// 	enum := NewSignalEnum("enum")

// 	sig, err := NewEnumSignal("sig", enum)
// 	assert.NoError(err)
// 	assert.NoError(msg.AppendSignal(sig))

// 	enumVal0 := NewSignalEnumValue("enum_val_0", 0)
// 	enumVal1 := NewSignalEnumValue("enum_val_1", 1)
// 	enumVal2 := NewSignalEnumValue("enum_val_2", 255)

// 	assert.NoError(enum.AddValue(enumVal0))
// 	assert.NoError(enum.AddValue(enumVal1))
// 	assert.NoError(enum.AddValue(enumVal2))

// 	// should remove enumVal0
// 	assert.NoError(enum.RemoveValue(enumVal0.EntityID()))

// 	expectedNames := []string{"enum_val_1", "enum_val_2"}
// 	for idx, val := range enum.Values() {
// 		assert.Equal(expectedNames[idx], val.Name())
// 	}

// 	// should remove enumVal2 and set the max index to 1
// 	assert.NoError(enum.RemoveValue(enumVal2.EntityID()))
// 	assert.Equal(1, len(enum.Values()))
// 	assert.Equal(1, enum.MaxIndex())

// 	// should return an error because enumVal0 is not part of the enum
// 	assert.Error(enum.RemoveValue(enumVal0.EntityID()))
// }

// func Test_SignalEnumValue_UpdateName(t *testing.T) {
// 	assert := assert.New(t)

// 	enum := NewSignalEnum("enum")

// 	enumVal0 := NewSignalEnumValue("enum_val_0", 0)
// 	enumVal1 := NewSignalEnumValue("enum_val_1", 1)
// 	enumVal2 := NewSignalEnumValue("enum_val_2", 2)

// 	assert.NoError(enum.AddValue(enumVal0))
// 	assert.NoError(enum.AddValue(enumVal1))

// 	// should rename enumVal0 to my_new_enum_name
// 	assert.NoError(enumVal0.UpdateName("my_new_enum_name"))
// 	assert.Equal("my_new_enum_name", enumVal0.Name())

// 	// should rename enumVal2 to my_new_enum_name
// 	assert.NoError(enumVal2.UpdateName("my_new_enum_name"))
// 	assert.Equal("my_new_enum_name", enumVal2.Name())

// 	// should return an error because my_new_enum_name is already taken
// 	assert.Error(enumVal1.UpdateName("my_new_enum_name"))
// }

// func Test_SignalEnumValue_UpdateIndex(t *testing.T) {
// 	assert := assert.New(t)

// 	msg0 := NewMessage("msg_0", 1, 1)
// 	msg1 := NewMessage("msg_1", 2, 2)

// 	enum := NewSignalEnum("enum")

// 	enumVal := NewSignalEnumValue("enum_val", 0)

// 	assert.NoError(enum.AddValue(enumVal))

// 	sig0, err := NewEnumSignal("sig_0", enum)
// 	assert.NoError(err)
// 	assert.NoError(msg0.AppendSignal(sig0))

// 	sig1, err := NewEnumSignal("sig_1", enum)
// 	assert.NoError(err)
// 	assert.NoError(msg1.AppendSignal(sig1))

// 	// should not return error because there is no change in the index
// 	assert.NoError(enumVal.UpdateIndex(0))
// 	assert.Equal(0, enumVal.Index())

// 	// should set the index to 8
// 	assert.NoError(enumVal.UpdateIndex(8))
// 	assert.Equal(8, enumVal.Index())

// 	// should set the index to 255
// 	assert.NoError(enumVal.UpdateIndex(255))
// 	assert.Equal(255, enumVal.Index())

// 	// should return an error because msg0 has a payload of 8 bits
// 	assert.Error(enumVal.UpdateIndex(256))
// 	assert.Equal(255, enumVal.Index())

// 	// should set the index to 8
// 	assert.NoError(enumVal.UpdateIndex(8))
// 	assert.Equal(8, enumVal.Index())
// }

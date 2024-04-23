package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_signal_UpdateName(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg", 2)

	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	stdSig0, err := NewStandardSignal("std_sig_0", size8Type)
	assert.NoError(err)
	stdSig1, err := NewStandardSignal("std_sig_1", size8Type)
	assert.NoError(err)

	assert.NoError(msg.AppendSignal(stdSig0))
	assert.NoError(msg.AppendSignal(stdSig1))

	// should be possible because the new name is not duplicated
	assert.NoError(stdSig0.UpdateName("std_sig_00"))
	assert.Equal("std_sig_00", stdSig0.Name())

	// should not return an error, but the name should not change
	assert.NoError(stdSig1.UpdateName("std_sig_1"))
	assert.Equal("std_sig_1", stdSig1.Name())

	//should return an error because the name is duplicated
	assert.Error(stdSig1.UpdateName("std_sig_00"))
}

func Test_StandardSignal(t *testing.T) {
	assert := assert.New(t)

	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	voltUnit := NewSignalUnit("volt", SignalUnitKindElectrical, "V")
	currUnit := NewSignalUnit("ampere", SignalUnitKindElectrical, "A")

	sig, err := NewStandardSignal("sig", size8Type)
	assert.NoError(err)
	sig.SetUnit(voltUnit)

	_, err = NewStandardSignal("sig", nil)
	assert.Error(err)

	assert.Equal(size8Type.EntityID(), sig.Type().EntityID())
	assert.Equal(size8Type.Min(), sig.Min())
	assert.Equal(size8Type.Max(), sig.Max())
	assert.Equal(float64(0), sig.Offset())
	assert.Equal(float64(1), sig.Scale())

	assert.Equal(voltUnit.EntityID(), sig.Unit().EntityID())
	sig.SetUnit(currUnit)
	assert.Equal(currUnit.EntityID(), sig.Unit().EntityID())

	assert.NoError(sig.SetPhysicalValues(-10, 10, 1, 2))
	assert.Equal(float64(-10), sig.Min())
	assert.Equal(float64(10), sig.Max())
	assert.Equal(float64(1), sig.Offset())
	assert.Equal(float64(2), sig.Scale())

	assert.Error(sig.SetPhysicalValues(-10, 10, 1, 0))

	_, err = sig.ToStandard()
	assert.NoError(err)
	_, err = sig.ToEnum()
	assert.Error(err)
	_, err = sig.ToMultiplexer()
	assert.Error(err)
}

func Test_StandardSignal_SetType(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg", 2)

	size16Type, err := NewIntegerSignalType("16_bits", 16, false)
	assert.NoError(err)
	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", size16Type)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size8Type)
	assert.NoError(err)

	// starting from this message payload
	// 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
	assert.NoError(msg.AppendSignal(sig0))

	// should get this message payload after setting to the new type
	// 0 0 0 0 0 0 0 0 - - - - - - - -
	assert.NoError(sig0.SetType(size8Type))
	assert.Equal(0, sig0.GetStartBit())
	assert.Equal(8, sig0.GetSize())

	// starting from this message payload
	// 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1
	assert.NoError(msg.AppendSignal(sig1))

	// should get an error because a 16 bit signal cannot fit in the message payload
	assert.Error(sig1.SetType(size16Type))

	// should get an error because the signal type is nil
	assert.Error(sig1.SetType(nil))
}

func Test_EnumSignal_SetEnum(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg", 2)

	size8Enum := NewSignalEnum("8_bits")
	assert.NoError(size8Enum.AddValue(NewSignalEnumValue("val_255", 255)))
	size9Enum := NewSignalEnum("9_bits")
	assert.NoError(size9Enum.AddValue(NewSignalEnumValue("val_256", 256)))

	sig0, err := NewEnumSignal("sig_0", size9Enum)
	assert.NoError(err)
	sig1, err := NewEnumSignal("sig_1", size8Enum)
	assert.NoError(err)

	// starting from this message payload
	// 0 0 0 0 0 0 0 0 0 - - - - - - -
	assert.NoError(msg.AppendSignal(sig0))

	// should get this message payload after setting to the new enum
	// 0 0 0 0 0 0 0 0 - - - - - - - -
	assert.NoError(sig0.SetEnum(size8Enum))
	assert.Equal(0, sig0.GetStartBit())
	assert.Equal(8, sig0.GetSize())

	// starting from this message payload
	// 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1
	assert.NoError(msg.AppendSignal(sig1))

	// should get an error because the message payload is full
	assert.Error(sig0.SetEnum(size9Enum))

	// should get an error because the enum is nil
	assert.Error(sig0.SetEnum(nil))
}

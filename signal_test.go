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

	stdSig0, err := NewStandardSignal("std_sig_0", size8Type, size8Type.Min(), size8Type.Max(), 0, 1, nil)
	assert.NoError(err)
	stdSig1, err := NewStandardSignal("std_sig_1", size8Type, size8Type.Min(), size8Type.Max(), 0, 1, nil)
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

	voltUnit := NewSignalUnit("volt", SignalUnitKindVoltage, "V")
	currUnit := NewSignalUnit("ampere", SignalUnitKindCurrent, "A")

	sig, err := NewStandardSignal("sig", size8Type, size8Type.Min(), size8Type.Max(), 0, 1, voltUnit)
	assert.NoError(err)

	_, err = NewStandardSignal("sig", nil, 0, 0, 0, 1, nil)
	assert.Error(err)

	assert.Equal(size8Type.EntityID(), sig.Type().EntityID())
	assert.Equal(size8Type.Min(), sig.Min())
	assert.Equal(size8Type.Max(), sig.Max())
	assert.Equal(float64(0), sig.Offset())
	assert.Equal(float64(1), sig.Scale())

	assert.Equal(voltUnit.EntityID(), sig.Unit().EntityID())
	sig.SetUnit(currUnit)
	assert.Equal(currUnit.EntityID(), sig.Unit().EntityID())

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

	sig0, err := NewStandardSignal("sig_0", size16Type, size16Type.Min(), size16Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size8Type, size8Type.Min(), size8Type.Max(), 0, 1, nil)
	assert.NoError(err)

	// starting from this message payload
	// 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
	assert.NoError(msg.AppendSignal(sig0))

	// should get this message payload after setting to the new type
	// 0 0 0 0 0 0 0 0 - - - - - - - -
	assert.NoError(sig0.SetType(size8Type, size8Type.Min(), size8Type.Max(), 0, 1))
	assert.Equal(0, sig0.GetStartBit())
	assert.Equal(8, sig0.GetSize())

	// starting from this message payload
	// 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1
	assert.NoError(msg.AppendSignal(sig1))

	// should get an error because a 16 bit signal cannot fit in the message payload
	assert.Error(sig1.SetType(size16Type, size16Type.Min(), size16Type.Max(), 0, 1))

	// should get an error because the signal type is nil
	assert.Error(sig1.SetType(nil, 0, 0, 0, 1))
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

func Test_MultiplexerSignal_AppendMuxSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg", 2)

	muxSig0, err := NewMultiplexerSignal("mux_sig_0", 16, 1)
	assert.NoError(err)

	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)
	size4Type, err := NewIntegerSignalType("4_bits", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", size8Type, size8Type.Min(), size8Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig3, err := NewStandardSignal("sig_3", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig4, err := NewStandardSignal("sig_4", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)

	// starting with a multiplexer signal that spreads over the entire message payload
	assert.NoError(msg.AppendSignal(muxSig0))

	// multiplexed signals at select value 0
	// sv0 0 0 0 0 0 0 0 0 - - - - - - -
	assert.NoError(muxSig0.AppendMuxSignal(0, sig0))

	expectedStartBits := []int{1}
	for idx, tmpSig := range muxSig0.GetSelectedMuxSignals(0) {
		assert.Equal(expectedStartBits[idx], tmpSig.GetStartBit())
	}

	// multiplexed signals at select value 1
	// sv1 1 1 1 1 2 2 2 2 3 3 3 3 - - -
	assert.NoError(muxSig0.AppendMuxSignal(1, sig1))
	assert.NoError(muxSig0.AppendMuxSignal(1, sig2))
	assert.NoError(muxSig0.AppendMuxSignal(1, sig3))

	expectedStartBits = []int{1, 5, 9}
	for idx, tmpSig := range muxSig0.GetSelectedMuxSignals(1) {
		assert.Equal(expectedStartBits[idx], tmpSig.GetStartBit())
	}

	// should return an error because select value 5 cannot be written in 1 bit
	assert.Error(muxSig0.AppendMuxSignal(5, sig4))

	// should return an error because sig4 of size 4 cannot fit in the payload at select value 1
	assert.Error(muxSig0.AppendMuxSignal(1, sig4))

	muxSig1, err := NewMultiplexerSignal("mux_sig_1", 7, 1)
	assert.NoError(err)

	// muxSig0 at select value 0 should nest muxSig1 after sig0
	// muxSig0 : ( sv0 0 0 0 0 0 0 0 0, muxSig1 : ( sl1 - - - - - - ) )
	assert.NoError(muxSig0.AppendMuxSignal(0, muxSig1))
	assert.Equal(9, muxSig1.GetStartBit())

	// should return no error for adding sig4 inside muxSig1 at select value 1
	// muxSig0 : ( sv0 0 0 0 0 0 0 0 0, muxSig1 : ( sl1 4 4 4 4 - - ) )
	assert.NoError(muxSig1.AppendMuxSignal(1, sig4))
	assert.Equal(10, sig4.GetStartBit())

	sig5, err := NewStandardSignal("sig_5", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)

	// should return an error because select value is negative
	assert.Error(muxSig1.AppendMuxSignal(-1, sig5))

	// should return an error because signal name sig_4 is duplicated
	assert.NoError(sig5.UpdateName("sig_4"))
	assert.Error(muxSig1.AppendMuxSignal(0, sig5))

	// should return an error because signal name sig_0 is duplicated
	assert.NoError(sig5.UpdateName("sig_0"))
	assert.Error(muxSig1.AppendMuxSignal(0, sig5))

	// should return an error because there is no space left in the payload of muxSig1 at select value 1
	assert.Error(sig4.SetType(size8Type, size8Type.Min(), size8Type.Max(), 0, 1))

	// should be possible to update sig4 name
	assert.NoError(sig4.UpdateName("sig_44"))
	tmpSig4, err := msg.GetSignal(sig4.EntityID())
	assert.NoError(err)
	assert.Equal("sig_44", tmpSig4.Name())

	// should not be possible to update sig4 name
	assert.Error(sig4.UpdateName("sig_0"))

	// should return an error because there is no more space
	assert.Error(sig4.SetType(size8Type, size8Type.Min(), size8Type.Max(), 0, 1))

	size2Type, err := NewIntegerSignalType("2_bits", 2, false)
	assert.NoError(err)

	// this time it should be possible because it is shrinking
	assert.NoError(sig4.SetType(size2Type, size2Type.Min(), size2Type.Max(), 0, 1))

	// also this time because it returns to the original size
	assert.NoError(sig4.SetType(size4Type, size4Type.Min(), size4Type.Max(), 0, 1))
}

func Test_MultiplexerSignal_RemoveMuxSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg", 2)

	muxSig0, err := NewMultiplexerSignal("mux_sig_0", 16, 1)
	assert.NoError(err)
	muxSig1, err := NewMultiplexerSignal("mux_sig_1", 7, 1)
	assert.NoError(err)

	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)
	size4Type, err := NewIntegerSignalType("4_bits", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", size8Type, size8Type.Min(), size8Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)

	// starting with this message payload
	// muxSig0 : ( sv0 0 0 0 0 0 0 0 0, muxSig1 : ( sl0 1 1 1 1 - - ) )
	assert.NoError(msg.AppendSignal(muxSig0))
	assert.NoError(muxSig0.AppendMuxSignal(0, sig0))
	assert.NoError(muxSig0.AppendMuxSignal(0, muxSig1))
	assert.NoError(muxSig1.AppendMuxSignal(0, sig1))

	// remove sig1, then sig0, then muxSig1
	assert.NoError(muxSig1.RemoveMuxSignal(sig1.EntityID()))
	assert.NoError(muxSig0.RemoveMuxSignal(sig0.EntityID()))
	assert.NoError(muxSig0.RemoveMuxSignal(muxSig1.EntityID()))

	assert.Equal(1, len(msg.Signals()))
	assert.Equal(0, len(muxSig0.GetSelectedMuxSignals(0)))

}

func Test_MultiplexerSignal_RemoveAllMuxSignals(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg", 2)

	muxSig0, err := NewMultiplexerSignal("mux_sig_0", 16, 1)
	assert.NoError(err)
	muxSig1, err := NewMultiplexerSignal("mux_sig_1", 7, 1)
	assert.NoError(err)

	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)
	size4Type, err := NewIntegerSignalType("4_bits", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", size8Type, size8Type.Min(), size8Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)

	// starting with this message payload
	// muxSig0 : ( sv0 0 0 0 0 0 0 0 0, muxSig1 : ( sl0 1 1 1 1 - - ) )
	assert.NoError(msg.AppendSignal(muxSig0))
	assert.NoError(muxSig0.AppendMuxSignal(0, sig0))
	assert.NoError(muxSig0.AppendMuxSignal(0, muxSig1))
	assert.NoError(muxSig1.AppendMuxSignal(0, sig1))

	// remove muxSig1 and sig1
	muxSig0.RemoveAllMuxSignals()

	assert.Equal(1, len(msg.Signals()))
	assert.Equal(0, len(muxSig0.GetSelectedMuxSignals(0)))
}

func Test_MultiplexerSignal_AddSelectValueRange(t *testing.T) {
	assert := assert.New(t)

	muxSig, err := NewMultiplexerSignal("mux_sig", 16, 4)
	assert.NoError(err)

	size4Type, err := NewIntegerSignalType("4_bits", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig4, err := NewStandardSignal("sig_4", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)
	sig5, err := NewStandardSignal("sig_5", size4Type, size4Type.Min(), size4Type.Max(), 0, 1, nil)
	assert.NoError(err)

	// setting a value range between 0 and 4
	assert.NoError(muxSig.AddSelectValueRange(0, 4))

	// appending sig0 and sig1 in the value range created before
	assert.NoError(muxSig.InsertMuxSignal(2, sig0, 0))
	assert.NoError(muxSig.InsertMuxSignal(0, sig1, 4))

	// should get sig0 and sig1 because 3 is in the range
	expectedNames := []string{"sig_0", "sig_1"}
	for idx, tmpSig := range muxSig.GetSelectedMuxSignals(3) {
		assert.Equal(expectedNames[idx], tmpSig.Name())
	}

	// inserting a signal with select value of 8, and creating a value range between 5 and 10
	assert.NoError(muxSig.InsertMuxSignal(8, sig2, 0))
	assert.NoError(muxSig.AddSelectValueRange(5, 10))
	assert.Equal(1, len(muxSig.GetSelectedMuxSignals(10)))

	// should return an error because from is greater then to
	assert.Error(muxSig.AddSelectValueRange(20, 5))

	// should return an error because from is invalid
	assert.Error(muxSig.AddSelectValueRange(1024, 1025))

	// should return an error because to is invalid
	assert.Error(muxSig.AddSelectValueRange(11, 1024))

	// should return an error because 3 is in another range
	assert.Error(muxSig.AddSelectValueRange(3, 12))

	// inserting sig4 and sig5 in different signal values
	assert.NoError(muxSig.InsertMuxSignal(11, sig4, 0))
	assert.NoError(muxSig.InsertMuxSignal(12, sig5, 0))

	// should return an error because in the range between 11 and 12 there are 2 different payloads
	assert.Error(muxSig.AddSelectValueRange(11, 12))
}

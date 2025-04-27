package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_signal_UpdateName(t *testing.T) {
	assert := assert.New(t)

	muxMsg := initMultiplexedMessage(assert)

	sigBase := muxMsg.signals.base
	assert.NoError(sigBase.UpdateName("new_base_signal"))
	assert.Equal("new_base_signal", sigBase.Name())
	assert.NoError(sigBase.UpdateName("base_signal"))
	assert.Error(sigBase.UpdateName("top_signal_in_0"))
	assert.Error(sigBase.UpdateName("top_inner_signal_in_0"))
	assert.Error(sigBase.UpdateName("bottom_signal_in_0"))
	assert.Error(sigBase.UpdateName("bottom_inner_signal_in_0"))

	sigTopIn0 := muxMsg.layers.top.signals.in0
	assert.NoError(sigTopIn0.UpdateName("new_top_signal_in_0"))
	assert.Equal("new_top_signal_in_0", sigTopIn0.Name())
	assert.NoError(sigTopIn0.UpdateName("top_signal_in_0"))
	assert.Error(sigTopIn0.UpdateName("base_signal"))
	assert.Error(sigTopIn0.UpdateName("top_signal_in_255"))
	assert.Error(sigTopIn0.UpdateName("top_inner_signal_in_0"))
	assert.Error(sigTopIn0.UpdateName("bottom_signal_in_0"))
	assert.Error(sigTopIn0.UpdateName("bottom_inner_signal_in_0"))

	sigToInnerIn0 := muxMsg.layers.top.inner.signals.in0
	assert.NoError(sigToInnerIn0.UpdateName("new_top_inner_signal_in_0"))
	assert.Equal("new_top_inner_signal_in_0", sigToInnerIn0.Name())
	assert.NoError(sigToInnerIn0.UpdateName("top_inner_signal_in_0"))
	assert.Error(sigToInnerIn0.UpdateName("base_signal"))
	assert.Error(sigToInnerIn0.UpdateName("top_signal_in_0"))
	assert.Error(sigToInnerIn0.UpdateName("top_inner_signal_in_255"))
	assert.Error(sigToInnerIn0.UpdateName("bottom_signal_in_0"))
	assert.Error(sigToInnerIn0.UpdateName("bottom_inner_signal_in_0"))

	sigTopMuxor := muxMsg.layers.top.layer.Muxor()
	assert.NoError(sigTopMuxor.UpdateName("new_top_muxor"))
	assert.Equal("new_top_muxor", muxMsg.layers.top.layer.Muxor().Name())
	assert.NoError(sigTopMuxor.UpdateName("top_muxor"))
	assert.Error(sigTopMuxor.UpdateName("top_inner_muxor"))
	assert.Error(sigTopMuxor.UpdateName("bottom_muxor"))
	assert.Error(sigTopMuxor.UpdateName("bottom_inner_muxor"))

	sigTopInnerMuxor := muxMsg.layers.top.inner.layer.Muxor()
	assert.NoError(sigTopInnerMuxor.UpdateName("new_top_inner_muxor"))
	assert.Equal("new_top_inner_muxor", sigTopInnerMuxor.Name())
	assert.NoError(sigTopInnerMuxor.UpdateName("top_inner_muxor"))
	assert.Error(sigTopInnerMuxor.UpdateName("top_muxor"))
	assert.Error(sigTopInnerMuxor.UpdateName("bottom_muxor"))
	assert.Error(sigTopInnerMuxor.UpdateName("bottom_inner_muxor"))
}

func Test_signal_UpdateStartPos(t *testing.T) {
	assert := assert.New(t)

	muxMsg := initMultiplexedMessage(assert)

	sigBase := muxMsg.signals.base
	assert.Error(sigBase.UpdateStartPos(0))
	assert.Error(sigBase.UpdateStartPos(8))

	assert.NoError(muxMsg.layers.top.signals.in255.UpdateStartPos(16))
	assert.Equal(16, muxMsg.layers.top.signals.in255.GetStartPos())
	assert.NoError(muxMsg.layers.top.signals.in255.UpdateStartPos(8))

	sigTopIn02 := muxMsg.layers.top.signals.in02
	assert.Error(sigTopIn02.UpdateStartPos(8))

	assert.Error(muxMsg.layers.top.inner.signals.in0.UpdateStartPos(58))
}

func Test_signal_updateSize(t *testing.T) {
	assert := assert.New(t)

	muxMsg := initMultiplexedMessage(assert)

	sigBase := muxMsg.signals.base
	assert.NoError(sigBase.updateSize(15))
	assert.Equal(15, sigBase.GetSize())
	assert.Error(sigBase.updateSize(17))
	assert.NoError(sigBase.updateSize(16))

	sigTopIn0 := muxMsg.layers.top.signals.in0
	assert.NoError(sigTopIn0.updateSize(7))
	assert.Equal(7, sigTopIn0.GetSize())
	assert.Error(sigTopIn0.updateSize(9))
	assert.NoError(sigTopIn0.updateSize(8))

	sigTopIn255 := muxMsg.layers.top.signals.in255
	assert.NoError(sigTopIn255.updateSize(16))
	assert.Equal(16, sigTopIn255.GetSize())
	assert.Error(sigTopIn255.updateSize(17))

	sigTopInnerIn0 := muxMsg.layers.top.inner.signals.in0
	assert.NoError(sigTopInnerIn0.updateSize(7))
	assert.Equal(7, sigTopInnerIn0.GetSize())
	assert.Error(sigTopInnerIn0.updateSize(9))
	assert.NoError(sigTopInnerIn0.updateSize(8))

	sigTopMuxor := muxMsg.layers.top.layer.Muxor()
	assert.NoError(sigTopMuxor.updateSize(3))
	assert.Equal(3, sigTopMuxor.GetSize())
	assert.NoError(sigTopMuxor.updateSize(8))
}

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func Test_signal_UpdateName(t *testing.T) {
// 	assert := assert.New(t)

// 	msg := NewMessage("msg", 1, 2)

// 	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
// 	assert.NoError(err)

// 	stdSig0, err := NewStandardSignal("std_sig_0", size8Type)
// 	assert.NoError(err)
// 	stdSig1, err := NewStandardSignal("std_sig_1", size8Type)
// 	assert.NoError(err)

// 	assert.NoError(msg.AppendSignal(stdSig0))
// 	assert.NoError(msg.AppendSignal(stdSig1))

// 	// should be possible because the new name is not duplicated
// 	assert.NoError(stdSig0.UpdateName("std_sig_00"))
// 	assert.Equal("std_sig_00", stdSig0.Name())

// 	// should not return an error, but the name should not change
// 	assert.NoError(stdSig1.UpdateName("std_sig_1"))
// 	assert.Equal("std_sig_1", stdSig1.Name())

// 	//should return an error because the name is duplicated
// 	assert.Error(stdSig1.UpdateName("std_sig_00"))
// }

// func Test_StandardSignal(t *testing.T) {
// 	assert := assert.New(t)

// 	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
// 	assert.NoError(err)

// 	voltUnit := NewSignalUnit("volt", SignalUnitKindElectrical, "V")
// 	currUnit := NewSignalUnit("ampere", SignalUnitKindElectrical, "A")

// 	sig, err := NewStandardSignal("sig", size8Type)
// 	assert.NoError(err)
// 	sig.SetUnit(voltUnit)

// 	_, err = NewStandardSignal("sig", nil)
// 	assert.Error(err)

// 	assert.Equal(size8Type.EntityID(), sig.Type().EntityID())

// 	assert.Equal(voltUnit.EntityID(), sig.Unit().EntityID())
// 	sig.SetUnit(currUnit)
// 	assert.Equal(currUnit.EntityID(), sig.Unit().EntityID())

// 	_, err = sig.ToStandard()
// 	assert.NoError(err)
// 	_, err = sig.ToEnum()
// 	assert.Error(err)
// 	_, err = sig.ToMultiplexer()
// 	assert.Error(err)
// }

// func Test_StandardSignal_SetType(t *testing.T) {
// 	assert := assert.New(t)

// 	msg := NewMessage("msg", 1, 2)

// 	size16Type, err := NewIntegerSignalType("16_bits", 16, false)
// 	assert.NoError(err)
// 	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
// 	assert.NoError(err)

// 	sig0, err := NewStandardSignal("sig_0", size16Type)
// 	assert.NoError(err)
// 	sig1, err := NewStandardSignal("sig_1", size8Type)
// 	assert.NoError(err)

// 	// starting from this message payload
// 	// 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
// 	assert.NoError(msg.AppendSignal(sig0))

// 	// should get this message payload after setting to the new type
// 	// 0 0 0 0 0 0 0 0 - - - - - - - -
// 	assert.NoError(sig0.SetType(size8Type))
// 	assert.Equal(0, sig0.GetStartBit())
// 	assert.Equal(8, sig0.GetSize())

// 	// starting from this message payload
// 	// 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1
// 	assert.NoError(msg.AppendSignal(sig1))

// 	// should get an error because a 16 bit signal cannot fit in the message payload
// 	assert.Error(sig1.SetType(size16Type))

// 	// should get an error because the signal type is nil
// 	assert.Error(sig1.SetType(nil))
// }

// func Test_EnumSignal_SetEnum(t *testing.T) {
// 	assert := assert.New(t)

// 	msg := NewMessage("msg", 1, 2)

// 	size8Enum := NewSignalEnum("8_bits")
// 	assert.NoError(size8Enum.AddValue(NewSignalEnumValue("val_255", 255)))
// 	size9Enum := NewSignalEnum("9_bits")
// 	assert.NoError(size9Enum.AddValue(NewSignalEnumValue("val_256", 256)))

// 	sig0, err := NewEnumSignal("sig_0", size9Enum)
// 	assert.NoError(err)
// 	sig1, err := NewEnumSignal("sig_1", size8Enum)
// 	assert.NoError(err)

// 	// starting from this message payload
// 	// 0 0 0 0 0 0 0 0 0 - - - - - - -
// 	assert.NoError(msg.AppendSignal(sig0))

// 	// should get this message payload after setting to the new enum
// 	// 0 0 0 0 0 0 0 0 - - - - - - - -
// 	assert.NoError(sig0.SetEnum(size8Enum))
// 	assert.Equal(0, sig0.GetStartBit())
// 	assert.Equal(8, sig0.GetSize())

// 	// starting from this message payload
// 	// 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1
// 	assert.NoError(msg.AppendSignal(sig1))

// 	// should get an error because the message payload is full
// 	assert.Error(sig0.SetEnum(size9Enum))

// 	// should get an error because the enum is nil
// 	assert.Error(sig0.SetEnum(nil))
// }

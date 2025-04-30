package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SignalLayout_insert(t *testing.T) {
	assert := assert.New(t)

	sl := newSL(8)

	// All signals are 8 bits long
	sigType, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", sigType)
	assert.NoError(err)

	// Should not return an error
	assert.NoError(sl.verifyAndInsert(sig0, 0))

	sig1, err := NewStandardSignal("sig_1", sigType)
	assert.NoError(err)

	// Should return an error
	assert.Error(sl.verifyAndInsert(sig1, 0))

	// Should not return an error
	assert.NoError(sl.verifyAndInsert(sig1, 16))

	sig2, err := NewStandardSignal("sig_2", sigType)
	assert.NoError(err)

	// Should return an error in each case
	assert.Error(sl.verifyAndInsert(sig2, -1))
	assert.Error(sl.verifyAndInsert(sig2, 63))
	assert.Error(sl.verifyAndInsert(sig2, 57))
	assert.Error(sl.verifyAndInsert(sig2, 7))
	assert.Error(sl.verifyAndInsert(sig2, 23))

	// Should not return an error
	assert.NoError(sl.verifyAndInsert(sig2, 8))

	// Should have 3 signals and filters
	assert.Len(sl.Signals(), 3)
	assert.Len(sl.Filters(), 3)
}

func Test_SignalLayout_delete(t *testing.T) {
	assert := assert.New(t)

	sl := newSL(8)

	// All signals are 8 bits long
	sigType, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", sigType)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", sigType)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", sigType)
	assert.NoError(err)

	assert.NoError(sl.verifyAndInsert(sig0, 0))
	assert.NoError(sl.verifyAndInsert(sig1, 8))
	assert.NoError(sl.verifyAndInsert(sig2, 16))

	// Should have 3 signals
	assert.Len(sl.Signals(), 3)

	// Should delete signal 1
	sl.delete(sig1)

	// Should have 2 signals and filters
	assert.Len(sl.Signals(), 2)
	assert.Len(sl.Filters(), 2)
}

func Test_SignalLayout_clear(t *testing.T) {
	assert := assert.New(t)

	sl := newSL(8)

	// All signals are 8 bits long
	sigType, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", sigType)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", sigType)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", sigType)
	assert.NoError(err)

	assert.NoError(sl.verifyAndInsert(sig0, 0))
	assert.NoError(sl.verifyAndInsert(sig1, 8))
	assert.NoError(sl.verifyAndInsert(sig2, 16))

	// Should have 3 signals
	assert.Len(sl.Signals(), 3)

	// Should have 0 signals and filters
	sl.clear()
	assert.Len(sl.Signals(), 0)
	assert.Len(sl.Filters(), 0)
}

func Test_SignalLayout_updateStartPos(t *testing.T) {
	assert := assert.New(t)

	sl := newSL(8)

	// All signals are 8 bits long
	sigType, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", sigType)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", sigType)
	assert.NoError(err)

	assert.NoError(sl.verifyAndInsert(sig0, 0))
	assert.NoError(sl.verifyAndInsert(sig1, 8))

	// Should not return an error
	assert.NoError(sl.verifyAndUpdateStartPos(sig1, 16))

	// Should return an error
	assert.Error(sl.verifyAndUpdateStartPos(sig1, -1))
	assert.Error(sl.verifyAndUpdateStartPos(sig1, 63))
	assert.Error(sl.verifyAndUpdateStartPos(sig1, 57))
	assert.Error(sl.verifyAndUpdateStartPos(sig1, 7))

	// Should not return an error
	assert.NoError(sl.verifyAndUpdateStartPos(sig0, 8))

	// Should return an error
	assert.Error(sl.verifyAndUpdateStartPos(sig0, 15))

	// Should have 2 signals and filters
	assert.Len(sl.Signals(), 2)
	assert.Len(sl.Filters(), 2)
}

func Test_SignalLayout_updateSize(t *testing.T) {
	assert := assert.New(t)

	sl := newSL(8)

	// All signals are 8 bits long at start
	sigType, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", sigType)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", sigType)
	assert.NoError(err)

	assert.NoError(sl.verifyAndInsert(sig0, 0))
	assert.NoError(sl.verifyAndInsert(sig1, 8))

	// Should not return an error
	assert.NoError(sl.verifyAndUpdateSize(sig1, 32))

	assert.Equal(32, sig1.Size())

	// Should return an error
	assert.Error(sl.verifyAndUpdateSize(sig0, -1))
	assert.Error(sl.verifyAndUpdateSize(sig0, 9))
	assert.Error(sl.verifyAndUpdateSize(sig0, 65))

	// Should have 2 signals and 5 filters
	assert.Len(sl.Signals(), 2)
	assert.Len(sl.Filters(), 5)
}

func Test_SignalLayout_Compact(t *testing.T) {
	assert := assert.New(t)

	tdBasicMsg := initBasicMessage(assert)

	msgBasic := tdBasicMsg.message
	assert.NoError(msgBasic.InsertSignal(dummySignal, 56))

	tdBasicMsg.layout.Compact()
	assert.Equal(48, dummySignal.StartPos())

	assert.NoError(msgBasic.DeleteSignal(dummySignal.EntityID()))
}

func Test_SignalLayout_Decode(t *testing.T) {
	assert := assert.New(t)

	tdBasicMsg := initBasicMessage(assert)
	layout := tdBasicMsg.layout

	// Every raw value will be 11xxx1 where x are all 0s depending on the size of the signal
	data := []byte{
		0b00000001,
		0b00011100,
		0b11000000,
		0b00000001,
		0b00011100,
		0b11000000,
		0b00000000,
		0b00000000,
	}

	expectedRaw := []any{uint64(3073), uint64(3073), uint64(3073), uint64(3073)}
	decodings := layout.Decode(data)
	assert.Len(decodings, 4)
	for idx, dec := range decodings {
		assert.Equal(expectedRaw[idx], dec.RawValue)
	}

	// Sets the first signal to use a decimal signal type
	decSigType, err := NewDecimalSignalType("decimal", 12, false)
	assert.NoError(err)
	decSigType.SetScale(0.5)
	decSigType.SetOffset(1000)
	basic0Sig, err := tdBasicMsg.signals.basic0.ToStandard()
	assert.NoError(err)
	assert.NoError(basic0Sig.UpdateType(decSigType))
	expectedVal := []any{float64(2536.5), uint64(3073), uint64(3073), uint64(3073)}
	decodings = layout.Decode(data)
	assert.Len(decodings, 4)
	for idx, dec := range decodings {
		assert.Equal(expectedVal[idx], dec.Value)
	}
}

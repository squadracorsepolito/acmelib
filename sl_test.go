package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SL_insert(t *testing.T) {
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

func Test_SL_delete(t *testing.T) {
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

func Test_SL_clear(t *testing.T) {
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

func Test_SL_updateStartPos(t *testing.T) {
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

func Test_SL_updateSize(t *testing.T) {
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

	assert.Equal(32, sig1.GetSize())

	// Should return an error
	assert.Error(sl.verifyAndUpdateSize(sig0, -1))
	assert.Error(sl.verifyAndUpdateSize(sig0, 9))
	assert.Error(sl.verifyAndUpdateSize(sig0, 65))

	// Should have 2 signals and 5 filters
	assert.Len(sl.Signals(), 2)
	assert.Len(sl.Filters(), 5)
}

func Test_SL_ApplyMultiplexedLayer(t *testing.T) {
	assert := assert.New(t)

	sl := newSL(4)

	// All signals are 8 bits long at start
	sigType, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	ml0 := NewMultiplexedLayer(4, calcValueFromSize(8)-1, "ml_0_muxor")
	assert.Equal(8, ml0.Muxor().GetSize())

	assert.NoError(sl.ApplyMultiplexedLayer(ml0))

	// Define 2 signals to be inserted into the multiplexed layer 0
	ml0Sig0, err := NewStandardSignal("ml_0_sig_0", sigType)
	assert.NoError(err)
	ml0Sig1, err := NewStandardSignal("ml_0_sig_1", sigType)
	assert.NoError(err)

	// Should fail because muxor is at start position 0
	assert.Error(ml0.InsertSignal(ml0Sig0, 0, 0))
	assert.Error(ml0.InsertSignal(ml0Sig1, 0, 1))

	assert.NoError(ml0.InsertSignal(ml0Sig0, 8, 0))
	assert.NoError(ml0.InsertSignal(ml0Sig1, 8, 1))

	ml1 := NewMultiplexedLayer(4, calcValueFromSize(8)-1, "ml_1_muxor")

	// Should return an error because the muxor start position is set to 0
	assert.Error(sl.ApplyMultiplexedLayer(ml1))

	ml1sig0, err := NewStandardSignal("ml_1_sig_0", sigType)
	assert.NoError(err)
	ml1sig1, err := NewStandardSignal("ml_1_sig_1", sigType)
	assert.NoError(err)

	// Update the muxor start position and insert 2 signals
	assert.NoError(ml1.UpdateMuxorStartPos(16))
	assert.NoError(ml1.InsertSignal(ml1sig0, 0, 0))
	assert.NoError(ml1.InsertSignal(ml1sig1, 24, 1))

	// Should return an error because the first signal start position is set to 0
	assert.Error(sl.ApplyMultiplexedLayer(ml1))

	// Update the first signal start position to an invalid position
	assert.NoError(ml1.DeleteSignal(ml1sig0))
	assert.NoError(ml1.InsertSignal(ml1sig0, 8, 0))
	assert.Error(sl.ApplyMultiplexedLayer(ml1))

	// Should not return an error
	assert.NoError(ml1.DeleteSignal(ml1sig0))
	assert.NoError(ml1.InsertSignal(ml1sig0, 24, 0))
	assert.NoError(sl.ApplyMultiplexedLayer(ml1))

	// Should return nil because ml1 is already in the signal layout
	assert.Nil(sl.ApplyMultiplexedLayer(ml1))

	// Should return an error because ml2 size is defferent than the signal layout size
	ml2 := NewMultiplexedLayer(5, calcValueFromSize(8)-1, "ml_2_muxor")
	assert.Error(sl.ApplyMultiplexedLayer(ml2))

	sig0, err := NewStandardSignal("sig_0", sigType)
	assert.NoError(err)

	// Should not return an error
	assert.Error(sl.verifyAndInsert(sig0, 24))

	ml1.Clear()

	// Should not return an error
	assert.NoError(sl.verifyAndInsert(sig0, 24))
}

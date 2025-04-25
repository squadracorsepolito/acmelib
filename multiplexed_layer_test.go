package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MultiplexedLayer_InsertSignal(t *testing.T) {
	assert := assert.New(t)

	sigType, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	sl := newSL(5)

	ml0 := NewMultiplexedLayer(5, calcValueFromSize(8)-1, "ml_0_muxor")
	assert.NoError(sl.AttachMultiplexedLayer(ml0))

	// Add 2 signals in the same layout
	ml0sig00, err := NewStandardSignal("ml0_sig_0_0", sigType)
	assert.NoError(err)
	ml0sig01, err := NewStandardSignal("ml0_sig_0_1", sigType)
	assert.NoError(err)

	assert.NoError(ml0.InsertSignal(ml0sig00, 8, 0))
	assert.NoError(ml0.InsertSignal(ml0sig01, 16, 0, 1))

	// Should return an error because ml0_sig_0_0
	// is already inserted in layout 0 at position 8
	assert.Error(ml0.InsertSignal(ml0sig00, 16, 1))

	// Add a signal in all layout at position 24
	fixedSig, err := NewStandardSignal("fixed_sig", sigType)
	assert.NoError(err)
	assert.NoError(ml0.InsertSignal(fixedSig, 24))

	assert.Len(ml0.GetSignals(0), 3)
	assert.Len(ml0.GetSignals(1), 2)
	assert.Len(ml0.GetSignals(2), 1)
	assert.Equal(24, ml0.GetSignals(2)[0].GetRelativeStartPos())

	// In layer 2 a new multiplexed layer is added
	ml1 := NewMultiplexedLayer(5, calcValueFromSize(8)-1, "ml_1_muxor")
	assert.NoError(ml1.Muxor().UpdateStartPos(8))

	ml0l2 := ml0.GetLayout(2)
	assert.NoError(ml0l2.AttachMultiplexedLayer(ml1))

	// Add 2 signals in multiplexed layer 1
	ml1Sig00, err := NewStandardSignal("ml1_sig_0_0", sigType)
	assert.NoError(err)
	ml1Sig10, err := NewStandardSignal("ml1_sig_1_0", sigType)
	assert.NoError(err)

	assert.NoError(ml1.InsertSignal(ml1Sig00, 16, 0))

	assert.Error(ml1.InsertSignal(ml1Sig10, 0, 1))
	assert.Error(ml1.InsertSignal(ml1Sig10, 24, 1))
	assert.NoError(ml1.InsertSignal(ml1Sig10, 16, 1))

	// Add a signal at the end of the base layout
	sig0, err := NewStandardSignal("sig_0", sigType)
	assert.NoError(err)
	assert.NoError(sl.verifyAndInsert(sig0, 32))

	ml1sig01, err := NewStandardSignal("ml1_sig_0_1", sigType)
	assert.NoError(err)

	// Should return an error because of the signal in the base layout
	assert.Error(ml1.InsertSignal(ml1sig01, 32, 2))

	initMultiplexedMessage(assert)
}

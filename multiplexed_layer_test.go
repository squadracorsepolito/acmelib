package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var byteSigType, _ = NewIntegerSignalType("uint8_t", 8, false)
var wordSigType, _ = NewIntegerSignalType("uint16_t", 16, false)

type testdataMuxLayout struct {
	layout                             *SL
	top, topInner, bottom, bottomInner *MultiplexedLayer
}

func initMultiplexedLayout(assert *assert.Assertions) *testdataMuxLayout {
	base := newSL(8)
	baseSignal, err := NewStandardSignal("base_signal", wordSigType)
	assert.NoError(err)
	assert.NoError(base.verifyAndInsert(baseSignal, 24))

	// top multiplexed layer
	top := NewMultiplexedLayer(8, 256, "top_muxor")
	assert.NoError(base.AttachMultiplexedLayer(top))

	topSig0, err := NewStandardSignal("top_signal_in_0", byteSigType)
	assert.NoError(err)
	assert.NoError(top.InsertSignal(topSig0, 8, 0))

	topSig255, err := NewStandardSignal("top_signal_in_255", byteSigType)
	assert.NoError(err)
	assert.NoError(top.InsertSignal(topSig255, 8, 255))

	topSig02, err := NewStandardSignal("top_signal_in_0_2", byteSigType)
	assert.NoError(err)
	assert.NoError(top.InsertSignal(topSig02, 16, 0, 2))

	// top inner multiplexed layer
	topInn := NewMultiplexedLayer(8, 256, "top_inner_muxor")
	assert.NoError(topInn.Muxor().UpdateStartPos(8))
	assert.NoError(top.GetLayout(1).AttachMultiplexedLayer(topInn))

	topInnSig0, err := NewStandardSignal("top_inner_signal_in_0", byteSigType)
	assert.NoError(err)
	assert.NoError(topInn.InsertSignal(topInnSig0, 16, 0))

	topInnSig255, err := NewStandardSignal("top_inner_signal_in_255", byteSigType)
	assert.NoError(err)
	assert.NoError(topInn.InsertSignal(topInnSig255, 16, 255))

	// bottom multiplexed layer
	bottom := NewMultiplexedLayer(8, 256, "bottom_muxor")
	assert.NoError(bottom.Muxor().UpdateStartPos(56))
	assert.NoError(base.AttachMultiplexedLayer(bottom))

	bottomSig0, err := NewStandardSignal("bottom_signal_in_0", byteSigType)
	assert.NoError(err)
	assert.NoError(bottom.InsertSignal(bottomSig0, 48, 0))

	bottomSig255, err := NewStandardSignal("bottom_signal_in_255", byteSigType)
	assert.NoError(err)
	assert.NoError(bottom.InsertSignal(bottomSig255, 48, 255))

	bottomSig02, err := NewStandardSignal("bottom_signal_in_0_2", byteSigType)
	assert.NoError(err)
	assert.NoError(bottom.InsertSignal(bottomSig02, 40, 0, 2))

	// bottom inner multiplexed layer
	bottomInn := NewMultiplexedLayer(8, 256, "bottom_inner_muxor")
	assert.NoError(bottomInn.Muxor().UpdateStartPos(48))
	assert.NoError(bottom.GetLayout(1).AttachMultiplexedLayer(bottomInn))

	bottomInnSig0, err := NewStandardSignal("bottom_inner_signal_in_0", byteSigType)
	assert.NoError(err)
	assert.NoError(bottomInn.InsertSignal(bottomInnSig0, 40, 0))

	bottomInnSig255, err := NewStandardSignal("bottom_inner_signal_in_255", byteSigType)
	assert.NoError(err)
	assert.NoError(bottomInn.InsertSignal(bottomInnSig255, 40, 255))

	return &testdataMuxLayout{
		layout:      base,
		top:         top,
		topInner:    topInn,
		bottom:      bottom,
		bottomInner: bottomInn,
	}
}

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

	initMultiplexedLayout(assert)
}

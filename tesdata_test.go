package acmelib

import (
	"github.com/stretchr/testify/assert"
)

var byteSigType, _ = NewIntegerSignalType("uint8_t", 8, false)
var wordSigType, _ = NewIntegerSignalType("uint16_t", 16, false)

var dummySignal, _ = NewStandardSignal("dummy_signal", byteSigType)

type testdataMuxMsgLayerInner struct {
	layer   *MultiplexedLayer
	signals struct {
		in0, in255 Signal
	}
}

type testdataMuxMsgLayer struct {
	layer   *MultiplexedLayer
	signals struct {
		in0, in255, in02 Signal
	}

	inner testdataMuxMsgLayerInner
}

type testdataMuxMsg struct {
	message *Message
	layout  *SL
	signals struct {
		base Signal
	}

	layers struct {
		top, bottom testdataMuxMsgLayer
	}
}

func initMultiplexedMessage(assert *assert.Assertions) *testdataMuxMsg {
	msg := NewMessage("mux_message", 16, 8)

	baseSignal, err := NewStandardSignal("base_signal", wordSigType)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(baseSignal, 24))

	layout := msg.SignalLayout()

	// top multiplexed layer
	top, err := layout.AddMultiplexedLayer("top_muxor", 0, 256)
	assert.NoError(err)

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
	topInn, err := top.GetLayout(1).AddMultiplexedLayer("top_inner_muxor", 8, 256)
	assert.NoError(err)

	topInnSig0, err := NewStandardSignal("top_inner_signal_in_0", byteSigType)
	assert.NoError(err)
	assert.NoError(topInn.InsertSignal(topInnSig0, 16, 0))

	topInnSig255, err := NewStandardSignal("top_inner_signal_in_255", byteSigType)
	assert.NoError(err)
	assert.NoError(topInn.InsertSignal(topInnSig255, 16, 255))

	// bottom multiplexed layer
	bottom, err := layout.AddMultiplexedLayer("bottom_muxor", 56, 256)
	assert.NoError(err)

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
	bottomInn, err := bottom.GetLayout(1).AddMultiplexedLayer("bottom_inner_muxor", 48, 256)
	assert.NoError(err)

	bottomInnSig0, err := NewStandardSignal("bottom_inner_signal_in_0", byteSigType)
	assert.NoError(err)
	assert.NoError(bottomInn.InsertSignal(bottomInnSig0, 40, 0))

	bottomInnSig255, err := NewStandardSignal("bottom_inner_signal_in_255", byteSigType)
	assert.NoError(err)
	assert.NoError(bottomInn.InsertSignal(bottomInnSig255, 40, 255))

	return &testdataMuxMsg{
		message: msg,
		layout:  layout,
		signals: struct{ base Signal }{baseSignal},

		layers: struct {
			top, bottom testdataMuxMsgLayer
		}{
			top: testdataMuxMsgLayer{
				layer:   top,
				signals: struct{ in0, in255, in02 Signal }{topSig0, topSig255, topSig02},

				inner: testdataMuxMsgLayerInner{
					layer:   topInn,
					signals: struct{ in0, in255 Signal }{topInnSig0, topInnSig255},
				},
			},

			bottom: testdataMuxMsgLayer{
				layer:   bottom,
				signals: struct{ in0, in255, in02 Signal }{bottomSig0, bottomSig255, bottomSig02},

				inner: testdataMuxMsgLayerInner{
					layer:   bottomInn,
					signals: struct{ in0, in255 Signal }{bottomInnSig0, bottomInnSig255},
				},
			},
		},
	}
}

type testdataSimpleMuxMsg struct {
	message *Message
	layout  *SL

	layer        *MultiplexedLayer
	layerSignals struct {
		in0, in1, in2 Signal
	}
}

func initSimpleMuxMessage(assert *assert.Assertions) *testdataSimpleMuxMsg {
	msg := NewMessage("simple_mux_message", 32, 4)

	layout := msg.SignalLayout()

	muxLayer, err := layout.AddMultiplexedLayer("muxor", 0, 256)
	assert.NoError(err)

	sigIn0, err := NewStandardSignal("signal_in_0", byteSigType)
	assert.NoError(err)
	assert.NoError(muxLayer.InsertSignal(sigIn0, 8, 0))

	sigIn1, err := NewStandardSignal("signal_in_1", byteSigType)
	assert.NoError(err)
	assert.NoError(muxLayer.InsertSignal(sigIn1, 8, 1))

	sigIn2, err := NewStandardSignal("signal_in_2", byteSigType)
	assert.NoError(err)
	assert.NoError(muxLayer.InsertSignal(sigIn2, 8, 2))

	return &testdataSimpleMuxMsg{
		message: msg,
		layout:  layout,

		layer: muxLayer,
		layerSignals: struct {
			in0 Signal
			in1 Signal
			in2 Signal
		}{sigIn0, sigIn1, sigIn2},
	}
}

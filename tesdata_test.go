package acmelib

import (
	"fmt"

	"github.com/stretchr/testify/assert"
)

var byteSigType, _ = NewIntegerSignalType("uint8_t", 8, false)
var sigType12, _ = NewIntegerSignalType("uint12_t", 12, false)
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

type testdataBasicMessage struct {
	message *Message
	layout  *SL

	signals struct {
		basic0, basic1, basic2, basic3 Signal
	}
}

func initBasicMessage(assert *assert.Assertions) *testdataBasicMessage {
	msg := NewMessage("basic_message", 1, 8)

	sigBasic0, err := NewStandardSignal("basic_signal_0", sigType12)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sigBasic0, 0))

	sigBasic1, err := NewStandardSignal("basic_signal_1", sigType12)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sigBasic1, 12))

	sigBasic2, err := NewStandardSignal("basic_signal_2", sigType12)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sigBasic2, 24))

	sigBasic3, err := NewStandardSignal("basic_signal_3", sigType12)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sigBasic3, 36))

	return &testdataBasicMessage{
		message: msg,
		layout:  msg.SignalLayout(),

		signals: struct {
			basic0 Signal
			basic1 Signal
			basic2 Signal
			basic3 Signal
		}{sigBasic0, sigBasic1, sigBasic2, sigBasic3},
	}
}

type testdataMessage[S any] struct {
	message *Message
	layout  *SL
	signals S
}

type testdataEnumMessageSignal struct {
	signal *EnumSignal
	enum   *SignalEnum
}

type testdataEnumMessage = testdataMessage[struct {
	with4Values, with8Values, fixedSize testdataEnumMessageSignal
}]

func initEnumMessage(assert *assert.Assertions) *testdataEnumMessage {
	msg := NewMessage("enum_message", 4, 4)

	enum4 := NewSignalEnum("enum_with_4_values")
	for i := range 4 {
		_, err := enum4.AddValue0(i, fmt.Sprintf("enum_value_%d", i))
		assert.NoError(err)
	}
	sig4, err := NewEnumSignal("enum_signal_4_values", enum4)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig4, 0))

	enum8 := NewSignalEnum("enum_signal_8_values")
	for i := range 8 {
		_, err := enum8.AddValue0(i, fmt.Sprintf("enum_value_%d", i))
		assert.NoError(err)
	}
	sig8, err := NewEnumSignal("signal_with_8_values", enum8)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig8, 8))

	enumFixed := NewSignalEnum("enum_signal_fixed_size")
	enumFixed.SetFixedSize(true)
	assert.NoError(enumFixed.UpdateSize(8))
	_, err = enumFixed.AddValue0(0, "enum_value_0")
	assert.NoError(err)
	_, err = enumFixed.AddValue0(127, "enum_value_127")
	assert.NoError(err)
	sigFixed, err := NewEnumSignal("enum_signal_fixed_size", enumFixed)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sigFixed, 16))

	return &testdataEnumMessage{
		message: msg,
		layout:  msg.SignalLayout(),

		signals: struct {
			with4Values testdataEnumMessageSignal
			with8Values testdataEnumMessageSignal
			fixedSize   testdataEnumMessageSignal
		}{
			with4Values: testdataEnumMessageSignal{signal: sig4, enum: enum4},
			with8Values: testdataEnumMessageSignal{signal: sig8, enum: enum8},
			fixedSize:   testdataEnumMessageSignal{signal: sigFixed, enum: enumFixed},
		},
	}
}

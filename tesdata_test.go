package acmelib

import (
	"fmt"

	"github.com/stretchr/testify/assert"
)

var byteSigType, _ = NewIntegerSignalType("uint8_t", 8, false)
var sigType12, _ = NewIntegerSignalType("uint12_t", 12, false)
var wordSigType, _ = NewIntegerSignalType("uint16_t", 16, false)

var dummySignal, _ = NewStandardSignal("dummy_signal", byteSigType)

type testdataMessage[S any] struct {
	message *Message
	layout  *SignalLayout
	signals S
}

type testdataBasicMessage testdataMessage[struct {
	basic0, basic1, basic2, basic3 Signal
}]

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

type testdataTypedMessage = testdataMessage[struct {
	flag, intUnsigned, intSigned, decUnsigned, decSigned testdataTypedMessageSignal
}]

type testdataTypedMessageSignal struct {
	signal *StandardSignal
	typ    *SignalType
}

func initTypedMessage(assert *assert.Assertions) *testdataTypedMessage {
	msg := NewMessage("typed_message", 2, 7)

	unit := NewSignalUnit("voltage", SignalUnitKindElectrical, "V")

	flagType := NewFlagSignalType("flag_type")
	flagSig, err := NewStandardSignal("flag_signal", flagType)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(flagSig, 0))

	intUnsignedType, err := NewIntegerSignalType("int_unsigned_type", 8, false)
	assert.NoError(err)
	intUnsignedSig, err := NewStandardSignal("int_unsigned_signal", intUnsignedType)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(intUnsignedSig, 8))

	intSignedType, err := NewIntegerSignalType("int_signed_type", 8, true)
	assert.NoError(err)
	intSignedSig, err := NewStandardSignal("int_signed_signal", intSignedType)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(intSignedSig, 16))

	decUnsignedType, err := NewDecimalSignalType("dec_unsigned_type", 16, false)
	assert.NoError(err)
	decUnsignedType.SetScale(0.5)
	decUnsignedType.SetOffset(100.5)
	decUnsignedSig, err := NewStandardSignal("dec_unsigned_signal", decUnsignedType)
	assert.NoError(err)
	decUnsignedSig.SetUnit(unit)
	assert.NoError(msg.InsertSignal(decUnsignedSig, 24))

	decSignedType, err := NewDecimalSignalType("dec_signed_type", 16, true)
	assert.NoError(err)
	decSignedType.SetScale(0.5)
	decSignedType.SetOffset(100.5)
	decSignedSig, err := NewStandardSignal("dec_signed_signal", decSignedType)
	assert.NoError(err)
	decSignedSig.SetUnit(unit)
	assert.NoError(msg.InsertSignal(decSignedSig, 40))

	return &testdataTypedMessage{
		message: msg,
		layout:  msg.SignalLayout(),

		signals: struct {
			flag        testdataTypedMessageSignal
			intUnsigned testdataTypedMessageSignal
			intSigned   testdataTypedMessageSignal
			decUnsigned testdataTypedMessageSignal
			decSigned   testdataTypedMessageSignal
		}{
			flag:        testdataTypedMessageSignal{flagSig, flagType},
			intUnsigned: testdataTypedMessageSignal{intUnsignedSig, intUnsignedType},
			intSigned:   testdataTypedMessageSignal{intSignedSig, intSignedType},
			decUnsigned: testdataTypedMessageSignal{decUnsignedSig, decUnsignedType},
			decSigned:   testdataTypedMessageSignal{decSignedSig, decSignedType},
		},
	}
}

type testdataBigEndianMessage = testdataMessage[struct {
	big0, big1, big2, big3 Signal
}]

func initBigEndianMessage(assert *assert.Assertions) *testdataBigEndianMessage {
	msg := NewMessage("big_endian_message", 4, 7)

	sigBig0, err := NewStandardSignal("big_endian_signal_0", sigType12)
	assert.NoError(err)
	sigBig0.SetEndianness(EndiannessBigEndian)
	assert.NoError(msg.InsertSignal(sigBig0, StartPosFromBigEndian(7)))

	sigBig1, err := NewStandardSignal("big_endian_signal_1", sigType12)
	assert.NoError(err)
	sigBig1.SetEndianness(EndiannessBigEndian)
	assert.NoError(msg.InsertSignal(sigBig1, StartPosFromBigEndian(11)))

	sigBig2, err := NewStandardSignal("big_endian_signal_2", sigType12)
	assert.NoError(err)
	sigBig2.SetEndianness(EndiannessBigEndian)
	assert.NoError(msg.InsertSignal(sigBig2, StartPosFromBigEndian(31)))

	sigBig3, err := NewStandardSignal("big_endian_signal_3", sigType12)
	assert.NoError(err)
	sigBig3.SetEndianness(EndiannessBigEndian)
	assert.NoError(msg.InsertSignal(sigBig3, StartPosFromBigEndian(35)))

	return &testdataBigEndianMessage{
		message: msg,
		layout:  msg.SignalLayout(),

		signals: struct {
			big0 Signal
			big1 Signal
			big2 Signal
			big3 Signal
		}{sigBig0, sigBig1, sigBig2, sigBig3},
	}
}

type testdataEnumMessageSignal struct {
	signal *EnumSignal
	enum   *SignalEnum
}

type testdataEnumMessage = testdataMessage[struct {
	with4Values, with8Values, fixedSize testdataEnumMessageSignal
}]

func initEnumMessage(assert *assert.Assertions) *testdataEnumMessage {
	msg := NewMessage("enum_message", 8, 4)

	enum4 := NewSignalEnum("enum_with_4_values")
	for i := range 4 {
		_, err := enum4.AddValue(i, fmt.Sprintf("enum_value_%d", i))
		assert.NoError(err)
	}
	sig4, err := NewEnumSignal("enum_signal_4_values", enum4)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig4, 0))

	enum8 := NewSignalEnum("enum_signal_8_values")
	for i := range 8 {
		_, err := enum8.AddValue(i, fmt.Sprintf("enum_value_%d", i))
		assert.NoError(err)
	}
	sig8, err := NewEnumSignal("signal_with_8_values", enum8)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig8, 8))

	enumFixed := NewSignalEnum("enum_signal_fixed_size")
	enumFixed.SetFixedSize(true)
	assert.NoError(enumFixed.UpdateSize(8))
	_, err = enumFixed.AddValue(0, "enum_value_0")
	assert.NoError(err)
	_, err = enumFixed.AddValue(127, "enum_value_127")
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

type testdataSimpleMuxMessage struct {
	message *Message
	layout  *SignalLayout

	layer        *MultiplexedLayer
	layerSignals struct {
		in0, in1, in2 Signal
	}
}

func initSimpleMuxMessage(assert *assert.Assertions) *testdataSimpleMuxMessage {
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

	return &testdataSimpleMuxMessage{
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

type testdataMuxMessageLayerInner struct {
	layer   *MultiplexedLayer
	signals struct {
		in0, in255 Signal
	}
}

type testdataMuxMessageLayer struct {
	layer   *MultiplexedLayer
	signals struct {
		in0, in255, in02 Signal
	}

	inner testdataMuxMessageLayerInner
}

type testdataMuxMessage struct {
	testdataMessage[struct {
		base Signal
	}]

	layers struct {
		top, bottom testdataMuxMessageLayer
	}
}

func initMuxMessage(assert *assert.Assertions) *testdataMuxMessage {
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

	return &testdataMuxMessage{
		testdataMessage: testdataMessage[struct {
			base Signal
		}]{msg, layout, struct{ base Signal }{baseSignal}},

		layers: struct {
			top, bottom testdataMuxMessageLayer
		}{
			top: testdataMuxMessageLayer{
				layer:   top,
				signals: struct{ in0, in255, in02 Signal }{topSig0, topSig255, topSig02},

				inner: testdataMuxMessageLayerInner{
					layer:   topInn,
					signals: struct{ in0, in255 Signal }{topInnSig0, topInnSig255},
				},
			},

			bottom: testdataMuxMessageLayer{
				layer:   bottom,
				signals: struct{ in0, in255, in02 Signal }{bottomSig0, bottomSig255, bottomSig02},

				inner: testdataMuxMessageLayerInner{
					layer:   bottomInn,
					signals: struct{ in0, in255 Signal }{bottomInnSig0, bottomInnSig255},
				},
			},
		},
	}
}

type testdataNetwork struct {
	net             *Network
	bus             *Bus
	node0, recNode0 *Node

	messages struct {
		basic     *testdataBasicMessage
		typed     *testdataTypedMessage
		bigEndian *testdataBigEndianMessage
		enum      *testdataEnumMessage
		mux       *testdataMuxMessage
		simpleMux *testdataSimpleMuxMessage
	}
}

func initNetwork(assert *assert.Assertions) *testdataNetwork {
	net := NewNetwork("network")

	bus := NewBus("bus")
	assert.NoError(net.AddBus(bus))

	node0 := NewNode("node_0", 0, 1)
	node0Int := node0.GetInterface(0)
	assert.NoError(bus.AddNodeInterface(node0Int))

	recNode0 := NewNode("rec_node_0", 0, 1)
	recNode0Int := recNode0.GetInterface(0)
	assert.NoError(bus.AddNodeInterface(recNode0Int))

	basicMsg := initBasicMessage(assert)
	assert.NoError(node0Int.AddSentMessage(basicMsg.message))
	assert.NoError(basicMsg.message.AddReceiver(recNode0Int))

	typedMsg := initTypedMessage(assert)
	assert.NoError(node0Int.AddSentMessage(typedMsg.message))
	assert.NoError(typedMsg.message.AddReceiver(recNode0Int))

	bigEndianMsg := initBigEndianMessage(assert)
	assert.NoError(node0Int.AddSentMessage(bigEndianMsg.message))
	assert.NoError(bigEndianMsg.message.AddReceiver(recNode0Int))

	enumMsg := initEnumMessage(assert)
	assert.NoError(node0Int.AddSentMessage(enumMsg.message))
	assert.NoError(enumMsg.message.AddReceiver(recNode0Int))

	muxMsg := initMuxMessage(assert)
	assert.NoError(node0Int.AddSentMessage(muxMsg.message))
	assert.NoError(muxMsg.message.AddReceiver(recNode0Int))

	simpleMuxMsg := initSimpleMuxMessage(assert)
	assert.NoError(node0Int.AddSentMessage(simpleMuxMsg.message))
	assert.NoError(simpleMuxMsg.message.AddReceiver(recNode0Int))

	return &testdataNetwork{
		net:      net,
		bus:      bus,
		node0:    node0,
		recNode0: recNode0,
		messages: struct {
			basic     *testdataBasicMessage
			typed     *testdataTypedMessage
			bigEndian *testdataBigEndianMessage
			enum      *testdataEnumMessage
			mux       *testdataMuxMessage
			simpleMux *testdataSimpleMuxMessage
		}{
			basic:     basicMsg,
			typed:     typedMsg,
			bigEndian: bigEndianMsg,
			enum:      enumMsg,
			mux:       muxMsg,
			simpleMux: simpleMuxMsg,
		},
	}
}

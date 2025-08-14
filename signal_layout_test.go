package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SignalLayout_insert(t *testing.T) {
	assert := assert.New(t)

	sl := newSignalLayout(8)

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

	sl := newSignalLayout(8)

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

	sl := newSignalLayout(8)

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

	sl := newSignalLayout(8)

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

	sl := newSignalLayout(8)

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

	// Basic message
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

	// Typed message
	tdTypedMsg := initTypedMessage(assert)
	layout = tdTypedMsg.layout

	data = []byte{
		0b00000001,
		0b11000001,
		0b11000001,
		0b00000001,
		0b11000000,
		0b00000001,
		0b11000000,
	}

	expectedDec := []SignalDecoding{
		{ValueType: SignalValueTypeFlag, RawValue: 1, Value: true},
		{ValueType: SignalValueTypeUint, RawValue: 193, Value: uint64(193)},
		{ValueType: SignalValueTypeInt, RawValue: ^uint64(63 - 1), Value: int64(-63)},
		{ValueType: SignalValueTypeFloat, RawValue: 49153, Value: float64(24677), Unit: "V"},              // scale: 0.5; offset: 100.5
		{ValueType: SignalValueTypeFloat, RawValue: ^uint64(16383 - 1), Value: float64(-8091), Unit: "V"}, // scale: 0.5; offset: 100.5
	}

	decodings = layout.Decode(data)
	assert.Len(decodings, 5)
	for idx, dec := range decodings {
		assert.Equal(expectedDec[idx].ValueType, dec.ValueType)
		assert.Equal(expectedDec[idx].RawValue, dec.RawValue)
		assert.Equal(expectedDec[idx].Value, dec.Value)
		assert.Equal(expectedDec[idx].Unit, dec.Unit)
	}

	// Big endian message
	tdBigEndianMsg := initBigEndianMessage(assert)
	layout = tdBigEndianMsg.layout

	data = []byte{
		0b11000000,
		0b00011100,
		0b00000001,
		0b11000000,
		0b00011100,
		0b00000001,
		0b00000000,
		0b00000000,
	}

	decodings = layout.Decode(data)
	assert.Len(decodings, 4)
	for idx, dec := range decodings {
		assert.Equal(expectedRaw[idx], dec.RawValue)
	}

	// Enum message
	tdEnumMsg := initEnumMessage(assert)
	layout = tdEnumMsg.layout

	data = []byte{
		0b00000001,
		0b00000010,
		0b01111111,
		0b00000000,
	}

	expectedDec = []SignalDecoding{
		{ValueType: SignalValueTypeEnum, RawValue: 1, Value: "enum_value_1"},
		{ValueType: SignalValueTypeEnum, RawValue: 2, Value: "enum_value_2"},
		{ValueType: SignalValueTypeEnum, RawValue: 127, Value: "enum_value_127"},
	}

	decodings = layout.Decode(data)
	assert.Len(decodings, 3)
	for idx, dec := range decodings {
		assert.Equal(expectedDec[idx].ValueType, dec.ValueType)
		assert.Equal(expectedDec[idx].RawValue, dec.RawValue)
		assert.Equal(expectedDec[idx].Value, dec.Value)
	}

	// Multiplexed message
	tdMuxMsg := initMuxMessage(assert)
	layout = tdMuxMsg.layout

	data = []byte{
		0b00000001, // selecting layout 1 from top
		0b11111111, // selecting layout 255 from top inner
		0b11000001, // top inner in 255
		0b00000001,
		0b11000000,
		0b11000001, // bottom in 0,2
		0b11000001, // bottom in 0
		0b00000000, // selecting layout 0 from bottom
	}

	expectedDec = []SignalDecoding{
		{
			Signal:    tdMuxMsg.layers.top.layer.muxor,
			ValueType: SignalValueTypeUint,
			RawValue:  1,
			Value:     uint64(1),
		},
		{
			Signal:    tdMuxMsg.layers.top.inner.layer.muxor,
			ValueType: SignalValueTypeUint,
			RawValue:  255,
			Value:     uint64(255),
		},
		{
			Signal:    tdMuxMsg.layers.top.inner.signals.in255,
			ValueType: SignalValueTypeUint,
			RawValue:  193,
			Value:     uint64(193),
		},
		{
			Signal:    tdMuxMsg.signals.base,
			ValueType: SignalValueTypeUint,
			RawValue:  49153,
			Value:     uint64(49153),
		},
		{
			Signal:    tdMuxMsg.layers.bottom.layer.muxor,
			ValueType: SignalValueTypeUint,
			RawValue:  0,
			Value:     uint64(0),
		},
		{
			Signal:    tdMuxMsg.layers.bottom.signals.in02,
			ValueType: SignalValueTypeUint,
			RawValue:  193,
			Value:     uint64(193),
		},
		{
			Signal:    tdMuxMsg.layers.bottom.signals.in0,
			ValueType: SignalValueTypeUint,
			RawValue:  193,
			Value:     uint64(193),
		},
	}

	decodings = layout.Decode(data)
	assert.Len(decodings, 7)
	for idx, dec := range decodings {
		assert.Equal(expectedDec[idx].Signal.Name(), dec.Signal.Name())
		assert.Equal(expectedDec[idx].ValueType, dec.ValueType)
		assert.Equal(expectedDec[idx].RawValue, dec.RawValue)
		assert.Equal(expectedDec[idx].Value, dec.Value)
	}
}

func Test_SignalLayout_Encode(t *testing.T) {
	assert := assert.New(t)

	// Basic message
	tdBasicMsg := initBasicMessage(assert)
	layout := tdBasicMsg.layout

	assert.NoError(tdBasicMsg.signals.basic0.UpdateEncodedValue(3073))
	assert.NoError(tdBasicMsg.signals.basic1.UpdateEncodedValue(3073))
	assert.NoError(tdBasicMsg.signals.basic2.UpdateEncodedValue(3073))
	assert.NoError(tdBasicMsg.signals.basic3.UpdateEncodedValue(3073))

	expectedEncData := []byte{
		0b00000001,
		0b00011100,
		0b11000000,
		0b00000001,
		0b00011100,
		0b11000000,
		0b00000000,
		0b00000000,
	}

	assert.Equal(expectedEncData, layout.Encode())

	// Typed message
	tdTypedMsg := initTypedMessage(assert)
	layout = tdTypedMsg.layout

	assert.NoError(tdTypedMsg.signals.flag.signal.UpdateEncodedValue(1))
	assert.NoError(tdTypedMsg.signals.intUnsigned.signal.UpdateEncodedValue(193))
	assert.NoError(tdTypedMsg.signals.intSigned.signal.UpdateEncodedValue(-63))
	assert.NoError(tdTypedMsg.signals.decUnsigned.signal.UpdateEncodedValue(24677))
	assert.NoError(tdTypedMsg.signals.decSigned.signal.UpdateEncodedValue(-8091))

	expectedEncData = []byte{
		0b00000001,
		0b11000001,
		0b11000001,
		0b00000001,
		0b11000000,
		0b00000001,
		0b11000000,
	}

	assert.Equal(expectedEncData, layout.Encode())

	// Big endian message
	tdBigEndianMsg := initBigEndianMessage(assert)
	layout = tdBigEndianMsg.layout

	assert.NoError(tdBigEndianMsg.signals.big0.UpdateEncodedValue(3073))
	assert.NoError(tdBigEndianMsg.signals.big1.UpdateEncodedValue(3073))
	assert.NoError(tdBigEndianMsg.signals.big2.UpdateEncodedValue(3073))
	assert.NoError(tdBigEndianMsg.signals.big3.UpdateEncodedValue(3073))

	expectedEncData = []byte{
		0b11000000,
		0b00011100,
		0b00000001,
		0b11000000,
		0b00011100,
		0b00000001,
		0b00000000,
		0b00000000,
	}

	assert.Equal(expectedEncData, layout.Encode())

	// Enum message
	tdEnumMsg := initEnumMessage(assert)
	layout = tdEnumMsg.layout

	assert.NoError(tdEnumMsg.signals.with4Values.signal.UpdateEncodedValue(1))
	assert.NoError(tdEnumMsg.signals.with8Values.signal.UpdateEncodedValue(2))
	assert.NoError(tdEnumMsg.signals.fixedSize.signal.UpdateEncodedValue(127))

	expectedEncData = []byte{
		0b00000001,
		0b00000010,
		0b01111111,
		0b00000000,
	}

	assert.Equal(expectedEncData, layout.Encode())

	// Multiplexed message
	tdMuxMsg := initMuxMessage(assert)
	layout = tdMuxMsg.layout

	assert.NoError(tdMuxMsg.layers.top.layer.muxor.UpdateEncodedValue(1))
	assert.NoError(tdMuxMsg.layers.top.inner.layer.muxor.UpdateEncodedValue(255))
	assert.NoError(tdMuxMsg.layers.top.inner.signals.in255.UpdateEncodedValue(193))
	assert.NoError(tdMuxMsg.signals.base.UpdateEncodedValue(49153))
	assert.NoError(tdMuxMsg.layers.bottom.layer.muxor.UpdateEncodedValue(0))
	assert.NoError(tdMuxMsg.layers.bottom.signals.in02.UpdateEncodedValue(193))
	assert.NoError(tdMuxMsg.layers.bottom.signals.in0.UpdateEncodedValue(193))

	expectedEncData = []byte{
		0b00000001, // selecting layout 1 from top
		0b11111111, // selecting layout 255 from top inner
		0b11000001, // top inner in 255
		0b00000001,
		0b11000000,
		0b11000001, // bottom in 0,2
		0b11000001, // bottom in 0
		0b00000000, // selecting layout 0 from bottom
	}

	assert.Equal(expectedEncData, layout.Encode())
}

func Benchmark_SignalLayout_Decode(b *testing.B) {
	assert := assert.New(b)

	tdEnumMsg := initEnumMessage(assert)
	layout := tdEnumMsg.layout

	data := []byte{
		0b00000001,
		0b00000010,
		0b01111111,
		0b00000000,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		layout.Decode(data)
	}
}

func Benchmark_SignalLayout_DecodeMultiplexed(b *testing.B) {
	assert := assert.New(b)

	tdMuxMsg := initMuxMessage(assert)
	layout := tdMuxMsg.layout

	data := []byte{
		0b00000001, // selecting layout 1 from top
		0b11111111, // selecting layout 255 from top inner
		0b11000001, // top inner in 255
		0b00000001,
		0b11000000,
		0b11000001, // bottom in 0,2
		0b11000001, // bottom in 0
		0b00000000, // selecting layout 0 from bottom
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		layout.Decode(data)
	}
}

func Benchmark_SignalLayout_Encode(b *testing.B) {
	assert := assert.New(b)

	tdEnumMsg := initEnumMessage(assert)
	layout := tdEnumMsg.layout

	assert.NoError(tdEnumMsg.signals.with4Values.signal.UpdateEncodedValue(1))
	assert.NoError(tdEnumMsg.signals.with8Values.signal.UpdateEncodedValue(2))
	assert.NoError(tdEnumMsg.signals.fixedSize.signal.UpdateEncodedValue(127))

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		layout.Encode()
	}
}

func Benchmark_SignalLayout_EncodeMultiplexed(b *testing.B) {
	assert := assert.New(b)

	tdMuxMsg := initMuxMessage(assert)
	layout := tdMuxMsg.layout

	assert.NoError(tdMuxMsg.layers.top.layer.muxor.UpdateEncodedValue(1))
	assert.NoError(tdMuxMsg.layers.top.inner.layer.muxor.UpdateEncodedValue(255))
	assert.NoError(tdMuxMsg.layers.top.inner.signals.in255.UpdateEncodedValue(193))
	assert.NoError(tdMuxMsg.signals.base.UpdateEncodedValue(49153))
	assert.NoError(tdMuxMsg.layers.bottom.layer.muxor.UpdateEncodedValue(0))
	assert.NoError(tdMuxMsg.layers.bottom.signals.in02.UpdateEncodedValue(193))
	assert.NoError(tdMuxMsg.layers.bottom.signals.in0.UpdateEncodedValue(193))

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		layout.Encode()
	}
}

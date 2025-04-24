package acmelib

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func Test_SignalLayout_Unpack(t *testing.T) {
// 	assert := assert.New(t)

// 	// Define all signal types
// 	uint12, err := NewIntegerSignalType("uint12_t", 12, false)
// 	assert.NoError(err)
// 	int4, err := NewIntegerSignalType("int4_t", 4, true)
// 	assert.NoError(err)
// 	flag := NewFlagSignalType("flag_t")
// 	int20, err := NewIntegerSignalType("int20_t", 20, true)
// 	assert.NoError(err)

// 	// Define the signal enum
// 	enum := NewSignalEnum("enum_e")
// 	enum.SetMinSize(4)
// 	enum.AddValue(NewSignalEnumValue("OK", 1))

// 	// Define the little endian signals
// 	sig0le, err := NewStandardSignal("sig_0_le", uint12)
// 	assert.NoError(err)
// 	sig1le, err := NewStandardSignal("sig_1_le", int4)
// 	assert.NoError(err)
// 	sig2le, err := NewStandardSignal("sig_2_le", int20)
// 	assert.NoError(err)
// 	sig3le, err := NewEnumSignal("sig_3_le", enum)
// 	assert.NoError(err)
// 	sig4le, err := NewStandardSignal("sig_4_le", flag)
// 	assert.NoError(err)

// 	// Define the little endian message
// 	msgle := NewMessage("msg_le", 1, 8)
// 	assert.NoError(msgle.AppendSignal(sig0le))
// 	assert.NoError(msgle.InsertSignal(sig1le, 28))
// 	assert.NoError(msgle.InsertSignal(sig2le, 34))
// 	assert.NoError(msgle.InsertSignal(sig3le, 58))
// 	assert.NoError(msgle.InsertSignal(sig4le, 63))

// 	// Get the signal layout for the little endian message
// 	slle := msgle.SignalLayout()

// 	// Define the big endian signals
// 	sig0be, err := NewStandardSignal("sig_0_be", uint12)
// 	assert.NoError(err)
// 	sig1be, err := NewStandardSignal("sig_1_be", int4)
// 	assert.NoError(err)
// 	sig2be, err := NewStandardSignal("sig_2_be", int20)
// 	assert.NoError(err)
// 	sig3be, err := NewEnumSignal("sig_3_be", enum)
// 	assert.NoError(err)
// 	sig4be, err := NewStandardSignal("sig_4_be", flag)
// 	assert.NoError(err)

// 	// Define the big endian message
// 	msgbe := NewMessage("msg_be", 1, 8)
// 	assert.NoError(msgbe.AppendSignal(sig0be))
// 	assert.NoError(msgbe.InsertSignal(sig1be, 28))
// 	assert.NoError(msgbe.InsertSignal(sig2be, 34))
// 	assert.NoError(msgbe.InsertSignal(sig3be, 58))
// 	assert.NoError(msgbe.InsertSignal(sig4be, 63))
// 	msgbe.SetByteOrder(MessageByteOrderBigEndian)

// 	// Get the signal layout for the big endian message
// 	slbe := msgbe.SignalLayout()

// 	// Test the value of unpacked signals
// 	expectedValues := []any{uint64(2049), int64(-7), int64(-524287), "OK", true}

// 	data0 := []byte{
// 		0b00000001,
// 		0b00001000,
// 		0b00000000,
// 		0b10010000,
// 		0b00000100,
// 		0b00000000,
// 		0b00100000,
// 		0b10000100,
// 	}

// 	for idx, dec := range slle.Decode(data0) {
// 		assert.Equal(expectedValues[idx], dec.Value)
// 	}

// 	data1 := []byte{
// 		0b10000000,
// 		0b00010000,
// 		0b00000000,
// 		0b10010000,
// 		0b00100000,
// 		0b00000000,
// 		0b00000100,
// 		0b10000100,
// 	}

// 	for idx, dec := range slbe.Decode(data1) {
// 		assert.Equal(expectedValues[idx], dec.Value)
// 	}

// 	// Testing multiplexed signals

// 	// sigMux0Little, err := NewStandardSignal("sig_mux_0_le", int4)
// 	// assert.NoError(err)
// 	// sigMux1Little, err := NewStandardSignal("sig_mux_1_le", int20)
// 	// assert.NoError(err)
// 	// sigMux2Little, err := NewEnumSignal("sig_mux_2_le", enum)
// 	// assert.NoError(err)
// 	// sigMux3Little, err := NewStandardSignal("sig_mux_3_le", flag)
// 	// assert.NoError(err)

// 	// sigMuxorLittle, err := NewMultiplexerSignal("sig_muxor_le", 4096, 52)
// 	// assert.NoError(err)
// 	// assert.Equal(12, sigMuxorLittle.GetGroupCountSize())

// 	// assert.NoError(sigMuxorLittle.InsertSignal(sigMux0Little, 16, 0))
// 	// assert.NoError(sigMuxorLittle.InsertSignal(sigMux1Little, 22, 1))
// 	// assert.NoError(sigMuxorLittle.InsertSignal(sigMux2Little, 46, 2))
// 	// assert.NoError(sigMuxorLittle.InsertSignal(sigMux3Little, 51, 2049))

// 	// msgMuxLittle := NewMessage("msg_mux_le", 1, 8)
// 	// assert.NoError(msgMuxLittle.AppendSignal(sigMuxorLittle))

// 	// slMuxLittle := msgMuxLittle.SignalLayout()

// 	// log.Print(slMuxLittle)

// 	// dataMuxLittle := []byte{
// 	// 	0b00000001,
// 	// 	0b00001000,
// 	// 	0b00000000,
// 	// 	0b10010000,
// 	// 	0b00000100,
// 	// 	0b00000000,
// 	// 	0b00100000,
// 	// 	0b10000100,
// 	// }

// 	// for _, dec := range slMuxLittle.Decode(dataMuxLittle) {
// 	// 	log.Print(dec)
// 	// }
// }

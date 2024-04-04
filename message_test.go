package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Message(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)
	assert.Equal(msg.name, "msg_0")
	assert.Equal(msg.desc, "msg_0_desc")
	assert.Equal(msg.SizeByte(), 8)

}

func Test_Message_AppendSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)

	sigTypInt8, _ := newSignalType("int8", "int8_desc", SignalTypeKindInteger, 8, true, -128, 127)
	sigTypInt32, _ := newSignalType("int32", "int32_desc", SignalTypeKindInteger, 32, true, -128, 127)

	sigNames := []string{"sig_0", "sig_1", "sig_2", "sig_3", "sig_4"}

	sig0, _ := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig1, _ := NewStandardSignal(sigNames[1], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig2, _ := NewStandardSignal(sigNames[2], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig3, _ := NewStandardSignal(sigNames[3], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig4, _ := NewStandardSignal(sigNames[4], "", sigTypInt32, 0, 100, 0, 1, nil)

	assert.NoError(msg.AppendSignal(sig0))

	duplicatedSigName, _ := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.AppendSignal(duplicatedSigName))

	assert.NoError(msg.AppendSignal(sig1))
	assert.NoError(msg.AppendSignal(sig2))
	assert.NoError(msg.AppendSignal(sig3))

	assert.NoError(msg.AppendSignal(sig4))

	sigTypMassive, _ := newSignalType("massive", "", SignalTypeKindInteger, 128, true, -128, 127)
	massiveSig, _ := NewStandardSignal("massive_sig", "", sigTypMassive, 0, 100, 0, 1, nil)
	assert.Error(msg.AppendSignal(massiveSig))

	exidingSig, _ := NewStandardSignal("exiding_sig", "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.AppendSignal(exidingSig))

	results := msg.Signals()
	assert.Equal(len(results), 5)
	for idx, sig := range results {
		assert.Equal(sigNames[idx], sig.Name())
	}
}

func Test_Message_InsertSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)

	sigTypInt8, _ := newSignalType("int8", "", SignalTypeKindInteger, 8, true, -128, 127)
	sigTypInt32, _ := newSignalType("int32", "int32_desc", SignalTypeKindInteger, 32, true, -128, 127)

	sigNames := []string{"sig_0", "sig_1", "sig_2", "sig_3", "sig_4"}

	sig0, _ := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig1, _ := NewStandardSignal(sigNames[1], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig2, _ := NewStandardSignal(sigNames[2], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig3, _ := NewStandardSignal(sigNames[3], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig4, _ := NewStandardSignal(sigNames[4], "", sigTypInt32, 0, 100, 0, 1, nil)

	assert.NoError(msg.InsertSignal(sig0, 0))

	assert.NoError(msg.InsertSignal(sig1, 24))

	duplicatedSigName, _ := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignal(duplicatedSigName, 16))

	overlappingSig, _ := NewStandardSignal("overlapping_sig", "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignal(overlappingSig, 0))
	assert.Error(msg.InsertSignal(overlappingSig, 7))
	assert.Error(msg.InsertSignal(overlappingSig, 23))

	assert.NoError(msg.InsertSignal(sig2, 16))
	assert.NoError(msg.InsertSignal(sig3, 8))
	assert.NoError(msg.InsertSignal(sig4, 32))

	sigTypMassive, _ := newSignalType("massive", "", SignalTypeKindInteger, 128, true, -128, 127)
	massiveSig, _ := NewStandardSignal("massive_sig", "", sigTypMassive, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignal(massiveSig, 0))

	exidingSig, _ := NewStandardSignal("exiding_sig", "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignal(exidingSig, 0))
	assert.Error(msg.InsertSignal(exidingSig, 64))

	correctOrder := []string{"sig_0", "sig_3", "sig_2", "sig_1", "sig_4"}

	results := msg.Signals()
	assert.Equal(len(results), 5)
	for idx, sig := range results {
		assert.Equal(correctOrder[idx], sig.Name())
	}

}

func Test_Message_RemoveSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)

	sigTypInt8, _ := newSignalType("int8", "int8_desc", SignalTypeKindInteger, 8, true, -128, 127)
	sigTypInt32, _ := newSignalType("int32", "int32_desc", SignalTypeKindInteger, 32, true, -128, 127)

	sigNames := []string{"sig_0", "sig_1", "sig_2", "sig_3", "sig_4"}

	sig0, _ := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig1, _ := NewStandardSignal(sigNames[1], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig2, _ := NewStandardSignal(sigNames[2], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig3, _ := NewStandardSignal(sigNames[3], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig4, _ := NewStandardSignal(sigNames[4], "", sigTypInt32, 0, 100, 0, 1, nil)

	assert.NoError(msg.AppendSignal(sig0))
	assert.NoError(msg.AppendSignal(sig1))
	assert.NoError(msg.AppendSignal(sig2))
	assert.NoError(msg.AppendSignal(sig3))
	assert.NoError(msg.AppendSignal(sig4))

	assert.NoError(msg.RemoveSignal(sig2.entityID))
	assert.Error(msg.RemoveSignal(EntityID("invalid_entity_id")))

	correctOrder := []string{"sig_0", "sig_1", "sig_3", "sig_4"}

	results := msg.Signals()
	assert.Equal(len(results), 4)
	for idx, sig := range results {
		assert.Equal(correctOrder[idx], sig.Name())
	}

}

func Test_Message_CompactSignals(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)

	sigTypInt8, _ := newSignalType("int8", "int8_desc", SignalTypeKindInteger, 8, true, -128, 127)

	sig0, _ := NewStandardSignal("sig_0", "", sigTypInt8, 0, 100, 0, 1, nil)
	sig1, _ := NewStandardSignal("sig_1", "", sigTypInt8, 0, 100, 0, 1, nil)
	sig2, _ := NewStandardSignal("sig_2", "", sigTypInt8, 0, 100, 0, 1, nil)

	assert.NoError(msg.InsertSignal(sig0, 2))
	assert.NoError(msg.InsertSignal(sig1, 18))
	assert.NoError(msg.InsertSignal(sig2, 26))

	msg.CompactSignals()

	correctStartBits := []int{0, 8, 16}

	for idx, sig := range msg.Signals() {
		assert.Equal(correctStartBits[idx], sig.GetStartBit())
	}

}

func Test_Message_ShiftSignalLeft(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("message", "", 2)

	sigTypeInt4, err := newSignalType("int4", "", SignalTypeKindInteger, 4, true, -8, 7)
	assert.NoError(err)

	sig0, err := NewStandardSignal("signal_0", "", sigTypeInt4, -8, 7, 0, 1, nil)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig0, 12))

	assert.Equal(0, msg.ShiftSignalLeft(sig0.EntityID(), 0))

	assert.Equal(1, msg.ShiftSignalLeft(sig0.EntityID(), 1))
	assert.Equal(11, msg.ShiftSignalLeft(sig0.EntityID(), 16))

	sig1, err := NewStandardSignal("signal_1", "", sigTypeInt4, -8, 7, 0, 1, nil)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig1, 12))

	assert.Equal(8, msg.ShiftSignalLeft(sig1.EntityID(), 16))

	sig2, err := NewStandardSignal("signal_2", "", sigTypeInt4, -8, 7, 0, 1, nil)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig2, 12))

	assert.Equal(4, msg.ShiftSignalLeft(sig2.EntityID(), 16))

	finalStartBits := []int{0, 4, 8}
	for idx, sig := range msg.Signals() {
		assert.Equal(finalStartBits[idx], sig.GetStartBit())
	}

}

func Test_Message_ShiftSignalRight(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("message", "", 2)

	sigTypeInt4, _ := newSignalType("int4", "", SignalTypeKindInteger, 4, true, -8, 7)

	sig0, err := NewStandardSignal("signal_0", "", sigTypeInt4, -8, 7, 0, 1, nil)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig0, 0))

	assert.Equal(0, msg.ShiftSignalRight(sig0.EntityID(), 0))

	assert.Equal(1, msg.ShiftSignalRight(sig0.EntityID(), 1))
	assert.Equal(11, msg.ShiftSignalRight(sig0.EntityID(), 16))

	sig1, err := NewStandardSignal("signal_1", "", sigTypeInt4, -8, 7, 0, 1, nil)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig1, 0))

	assert.Equal(8, msg.ShiftSignalRight(sig1.EntityID(), 16))

	sig2, err := NewStandardSignal("signal_2", "", sigTypeInt4, -8, 7, 0, 1, nil)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig2, 0))

	assert.Equal(4, msg.ShiftSignalRight(sig2.EntityID(), 16))

	finalStartBits := []int{4, 8, 12}
	for idx, sig := range msg.Signals() {
		assert.Equal(finalStartBits[idx], sig.GetStartBit())
	}

}

package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Message(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)
	assert.Equal(msg.Name, "msg_0")
	assert.Equal(msg.Desc, "msg_0_desc")
	assert.Equal(msg.Size, 8)

	t.Log(msg.String())
}

func Test_Message_AppendSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)

	sigTypInt8 := NewSignalType("int8", "int8_desc", SignalTypeKindInteger, 8, true, -128, 127)
	sigTypInt32 := NewSignalType("int32", "int32_desc", SignalTypeKindInteger, 32, true, -128, 127)

	sigNames := []string{"sig_0", "sig_1", "sig_2", "sig_3", "sig_4"}

	sig0 := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig1 := NewStandardSignal(sigNames[1], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig2 := NewStandardSignal(sigNames[2], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig3 := NewStandardSignal(sigNames[3], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig4 := NewStandardSignal(sigNames[4], "", sigTypInt32, 0, 100, 0, 1, nil)

	assert.NoError(msg.AppendSignal(sig0))

	duplicatedSigName := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.AppendSignal(duplicatedSigName))

	assert.NoError(msg.AppendSignal(sig1))
	assert.NoError(msg.AppendSignal(sig2))
	assert.NoError(msg.AppendSignal(sig3))
	assert.NoError(msg.AppendSignal(sig4))

	sigTypMassive := NewSignalType("massive", "", SignalTypeKindInteger, 128, true, -128, 127)
	massiveSig := NewStandardSignal("massive_sig", "", sigTypMassive, 0, 100, 0, 1, nil)
	assert.Error(msg.AppendSignal(massiveSig))

	exidingSig := NewStandardSignal("exiding_sig", "", sigTypInt8, 0, 100, 0, 1, nil)
	t.Log(msg.String())
	assert.Error(msg.AppendSignal(exidingSig))

	results := msg.SignalsByPosition()
	assert.Equal(len(results), 5)
	for idx, sig := range results {
		assert.Equal(sig.Name, sigNames[idx])
		assert.Equal(sig.Position, idx)
	}

	t.Log(msg.String())
}

func Test_Message_InsertSignalAtPosition(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)

	sigTypInt8 := NewSignalType("int8", "", SignalTypeKindInteger, 8, true, -128, 127)
	sigTypInt32 := NewSignalType("int32", "int32_desc", SignalTypeKindInteger, 32, true, -128, 127)

	sigNames := []string{"sig_0", "sig_1", "sig_2", "sig_3", "sig_4"}

	sig0 := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig1 := NewStandardSignal(sigNames[1], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig2 := NewStandardSignal(sigNames[2], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig3 := NewStandardSignal(sigNames[3], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig4 := NewStandardSignal(sigNames[4], "", sigTypInt32, 0, 100, 0, 1, nil)

	assert.Error(msg.InsertSignalAtPosition(sig0, 1))

	assert.NoError(msg.InsertSignalAtPosition(sig0, 0))

	duplicatedSigName := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignalAtPosition(duplicatedSigName, 0))

	assert.NoError(msg.InsertSignalAtPosition(sig1, 1))
	assert.NoError(msg.InsertSignalAtPosition(sig2, 1))
	assert.NoError(msg.InsertSignalAtPosition(sig3, 1))
	assert.NoError(msg.InsertSignalAtPosition(sig4, 4))

	sigTypMassive := NewSignalType("massive", "", SignalTypeKindInteger, 128, true, -128, 127)
	massiveSig := NewStandardSignal("massive_sig", "", sigTypMassive, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignalAtPosition(massiveSig, 0))

	exidingSig := NewStandardSignal("exiding_sig", "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignalAtPosition(exidingSig, 0))

	correctOrder := []string{"sig_0", "sig_3", "sig_2", "sig_1", "sig_4"}

	results := msg.SignalsByPosition()
	assert.Equal(len(results), 5)
	for idx, sig := range results {
		assert.Equal(sig.Name, correctOrder[idx])
		assert.Equal(sig.Position, idx)
	}

	t.Log(msg.String())
}

func Test_Message_InsertSignalAtStartBit(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)

	sigTypInt8 := NewSignalType("int8", "", SignalTypeKindInteger, 8, true, -128, 127)
	sigTypInt32 := NewSignalType("int32", "int32_desc", SignalTypeKindInteger, 32, true, -128, 127)

	sigNames := []string{"sig_0", "sig_1", "sig_2", "sig_3", "sig_4"}

	sig0 := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig1 := NewStandardSignal(sigNames[1], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig2 := NewStandardSignal(sigNames[2], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig3 := NewStandardSignal(sigNames[3], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig4 := NewStandardSignal(sigNames[4], "", sigTypInt32, 0, 100, 0, 1, nil)

	assert.NoError(msg.InsertSignalAtStartBit(sig0, 0))

	assert.NoError(msg.InsertSignalAtStartBit(sig1, 24))

	duplicatedSigName := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignalAtStartBit(duplicatedSigName, 16))

	overlappingSig := NewStandardSignal("overlapping_sig", "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignalAtStartBit(overlappingSig, 0))
	assert.Error(msg.InsertSignalAtStartBit(overlappingSig, 7))
	assert.Error(msg.InsertSignalAtStartBit(overlappingSig, 23))

	assert.NoError(msg.InsertSignalAtStartBit(sig2, 16))
	assert.NoError(msg.InsertSignalAtStartBit(sig3, 8))
	assert.NoError(msg.InsertSignalAtStartBit(sig4, 32))

	sigTypMassive := NewSignalType("massive", "", SignalTypeKindInteger, 128, true, -128, 127)
	massiveSig := NewStandardSignal("massive_sig", "", sigTypMassive, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignalAtStartBit(massiveSig, 0))

	exidingSig := NewStandardSignal("exiding_sig", "", sigTypInt8, 0, 100, 0, 1, nil)
	assert.Error(msg.InsertSignalAtStartBit(exidingSig, 0))
	assert.Error(msg.InsertSignalAtStartBit(exidingSig, 64))

	correctOrder := []string{"sig_0", "sig_3", "sig_2", "sig_1", "sig_4"}

	results := msg.SignalsByPosition()
	assert.Equal(len(results), 5)
	for idx, sig := range results {
		assert.Equal(sig.Name, correctOrder[idx])
		assert.Equal(sig.Position, idx)
	}

	t.Log(msg.String())
}

func Test_Message_RemoveSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)

	sigTypInt8 := NewSignalType("int8", "int8_desc", SignalTypeKindInteger, 8, true, -128, 127)
	sigTypInt32 := NewSignalType("int32", "int32_desc", SignalTypeKindInteger, 32, true, -128, 127)

	sigNames := []string{"sig_0", "sig_1", "sig_2", "sig_3", "sig_4"}

	sig0 := NewStandardSignal(sigNames[0], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig1 := NewStandardSignal(sigNames[1], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig2 := NewStandardSignal(sigNames[2], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig3 := NewStandardSignal(sigNames[3], "", sigTypInt8, 0, 100, 0, 1, nil)
	sig4 := NewStandardSignal(sigNames[4], "", sigTypInt32, 0, 100, 0, 1, nil)

	assert.NoError(msg.AppendSignal(sig0))
	assert.NoError(msg.AppendSignal(sig1))
	assert.NoError(msg.AppendSignal(sig2))
	assert.NoError(msg.AppendSignal(sig3))
	assert.NoError(msg.AppendSignal(sig4))

	assert.NoError(msg.RemoveSignal(sig2.EntityID))
	assert.Error(msg.RemoveSignal(EntityID("invalid_entity_id")))

	correctOrder := []string{"sig_0", "sig_1", "sig_3", "sig_4"}

	results := msg.SignalsByPosition()
	assert.Equal(len(results), 4)
	for idx, sig := range results {
		assert.Equal(sig.Name, correctOrder[idx])
		assert.Equal(sig.Position, idx)
	}

	t.Log(msg.String())
}

func Test_Message_CompactSignals(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)

	sigTypInt8 := NewSignalType("int8", "int8_desc", SignalTypeKindInteger, 8, true, -128, 127)

	sig0 := NewStandardSignal("sig_0", "", sigTypInt8, 0, 100, 0, 1, nil)
	sig1 := NewStandardSignal("sig_1", "", sigTypInt8, 0, 100, 0, 1, nil)
	sig2 := NewStandardSignal("sig_2", "", sigTypInt8, 0, 100, 0, 1, nil)

	assert.NoError(msg.InsertSignalAtStartBit(sig0, 2))
	assert.NoError(msg.InsertSignalAtStartBit(sig1, 18))
	assert.NoError(msg.InsertSignalAtStartBit(sig2, 26))

	msg.CompactSignals()

	correctStartBits := []int{0, 8, 16}

	for idx, sig := range msg.SignalsByPosition() {
		assert.Equal(sig.StartBit, correctStartBits[idx])
	}

	t.Log(msg.String())
}

func Test_Message_GetAvailableSignalPositions(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", "msg_0_desc", 8)

	sigTypInt8 := NewSignalType("int8", "int8_desc", SignalTypeKindInteger, 8, true, -128, 127)

	sig0 := NewStandardSignal("sig_0", "", sigTypInt8, 0, 100, 0, 1, nil)
	sig1 := NewStandardSignal("sig_1", "", sigTypInt8, 0, 100, 0, 1, nil)
	sig2 := NewStandardSignal("sig_2", "", sigTypInt8, 0, 100, 0, 1, nil)

	assert.NoError(msg.InsertSignalAtStartBit(sig0, 2))
	assert.NoError(msg.InsertSignalAtStartBit(sig1, 18))
	assert.NoError(msg.InsertSignalAtStartBit(sig2, 26))

	positions := msg.GetAvailableSignalPositions()

	correctPositions := []*SignalPosition{NewSignalPosition(0, 1), NewSignalPosition(10, 17), NewSignalPosition(34, 63)}

	assert.Equal(len(positions), 3)
	for idx, pos := range positions {
		assert.Equal(pos.From, correctPositions[idx].From)
		assert.Equal(pos.To, correctPositions[idx].To)
	}

	t.Log(msg.String())
}

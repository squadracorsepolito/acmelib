package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MultiplexerSignal(t *testing.T) {
	assert := assert.New(t)

	muxSig, err := NewMultiplexerSignal("mux_sig", 4, 16)
	assert.NoError(err)
	assert.Equal(18, muxSig.GetSize())
	assert.Equal(2, muxSig.GetGroupCountSize())

	_, err = NewMultiplexerSignal("invalid_mux_sig", -1, 16)
	assert.ErrorAs(err, &ErrIsNegative)
	_, err = NewMultiplexerSignal("invalid_mux_sig", 0, 16)
	assert.ErrorAs(err, &ErrIsZero)
	_, err = NewMultiplexerSignal("invalid_mux_sig", 4, -1)
	assert.ErrorAs(err, &ErrIsNegative)
	_, err = NewMultiplexerSignal("invalid_mux_sig", 4, 0)
	assert.ErrorAs(err, &ErrIsZero)

	msg := NewMessage("msg", 1, 8)
	assert.NoError(msg.AppendSignal(muxSig))

	// testing multiplexed signal size change
	size4Type, err := NewIntegerSignalType("4_bits", 4, false)
	assert.NoError(err)
	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)
	size16Type, err := NewIntegerSignalType("16_bits", 16, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", size4Type)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size4Type)
	assert.NoError(err)

	assert.NoError(muxSig.InsertSignal(sig0, 0))
	assert.NoError(muxSig.InsertSignal(sig1, 8, 0))

	assert.NoError(sig0.SetType(size8Type))

	assert.Error(sig0.SetType(size16Type))
	assert.Error(sig1.SetType(size16Type))

	assert.NoError(sig1.SetType(size8Type))
	assert.NoError(sig1.SetType(size4Type))

	group := muxSig.GetSignalGroup(0)
	assert.Equal(0, group[0].GetRelativeStartPos())
	assert.Equal(8, group[1].GetRelativeStartPos())

	enum := NewSignalEnum("enum")
	enumVal := NewSignalEnumValue("val_0", 255)
	assert.NoError(enum.AddValue(enumVal))

	enumSig, err := NewEnumSignal("enum_sig", enum)
	assert.NoError(err)

	assert.NoError(muxSig.InsertSignal(enumSig, 8, 1))

	assert.NoError(enumVal.UpdateIndex(0))
	assert.NoError(enumVal.UpdateIndex(127))
	assert.Error(enumVal.UpdateIndex(256))

	// testing multiplexed signal name change
	msgSig, err := NewStandardSignal("msg_sig", size4Type)
	assert.NoError(err)

	assert.NoError(msg.AppendSignal(msgSig))

	assert.NoError(sig0.UpdateName("sig_00"))

	assert.Error(sig1.UpdateName("msg_sig"))
	assert.Error(sig1.UpdateName("sig_00"))

	group = muxSig.GetSignalGroup(0)
	assert.Equal("sig_00", group[0].Name())
	assert.Equal("sig_1", group[1].Name())

	nestedMuxSig, err := NewMultiplexerSignal("nested_mux_sig", 2, 4)
	assert.NoError(err)

	sig2, err := NewStandardSignal("sig_2", size4Type)
	assert.NoError(err)

	assert.NoError(nestedMuxSig.InsertSignal(sig2, 0))
	assert.NoError(muxSig.InsertSignal(nestedMuxSig, 8, 2))

	assert.NoError(sig2.UpdateName("sig_22"))
	assert.Equal("sig_22", nestedMuxSig.GetSignalGroup(0)[0].Name())
}

func Test_MultiplexerSignal_InsertSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg", 1, 8)

	muxSig, err := NewMultiplexerSignal("mux_sig", 4, 16)
	assert.NoError(err)

	assert.NoError(msg.AppendSignal(muxSig))

	size4Type, err := NewIntegerSignalType("4_bits", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", size4Type)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size4Type)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", size4Type)
	assert.NoError(err)
	sig3, err := NewStandardSignal("sig_3", size4Type)
	assert.NoError(err)

	// insert sig0 as a fixed signal
	// gID gID | - - - - - - - - 0 0 0 0 - - - -
	assert.NoError(muxSig.InsertSignal(sig0, 8))
	groups := muxSig.GetSignalGroups()
	assert.Len(groups, 4)
	for _, group := range groups {
		assert.Len(group, 1)
		assert.Equal("sig_0", group[0].Name())
	}

	// insert sig1 in groupID 0 and 2
	// gID gID | 1 1 1 1 - - - - 0 0 0 0 - - - -
	assert.NoError(muxSig.InsertSignal(sig1, 0, 0, 2))
	group := muxSig.GetSignalGroup(0)
	assert.Len(group, 2)
	assert.Equal("sig_1", group[0].Name())
	group = muxSig.GetSignalGroup(2)
	assert.Len(group, 2)
	assert.Equal("sig_1", group[0].Name())

	// insert sig2 and sig3 in groupID 0
	// gID gID | 1 1 1 1 2 2 2 2 0 0 0 0 3 3 3 3
	assert.NoError(muxSig.InsertSignal(sig2, 4, 0))
	assert.NoError(muxSig.InsertSignal(sig3, 12, 0))
	assert.Len(muxSig.GetSignalGroup(0), 4)

	// insert sig2 in groupID 2
	// gID gID | 1 1 1 1 2 2 2 2 0 0 0 0 - - - -
	assert.NoError(muxSig.InsertSignal(sig2, 4, 2))

	// should return an error because sig2 is already in groupID 2
	assert.Error(muxSig.InsertSignal(sig2, 4, 2))

	muxSig.ClearAllSignalGroups()

	nestedMuxSig, err := NewMultiplexerSignal("nested_mux_sig", 2, 8)
	assert.NoError(err)

	// insert nestedMuxSig in groupID 0 of muxSig
	assert.NoError(muxSig.InsertSignal(nestedMuxSig, 0, 0))

	// insert sig0 and sig1 as fixed signals in nestedMuxSig
	assert.NoError(nestedMuxSig.InsertSignal(sig0, 0))
	assert.NoError(nestedMuxSig.InsertSignal(sig1, 4))

	group = muxSig.GetSignalGroup(0)
	assert.Len(group, 1)
	assert.Equal("nested_mux_sig", group[0].Name())

	tmpNestedMuxSig, err := group[0].ToMultiplexer()
	assert.NoError(err)
	nestedGroups := tmpNestedMuxSig.GetSignalGroups()
	assert.Len(nestedGroups, 2)
	for _, group := range nestedGroups {
		assert.Equal("sig_0", group[0].Name())
		assert.Equal("sig_1", group[1].Name())
	}

	// should return an error because sig0 is already in nestedMuxSig
	assert.Error(muxSig.InsertSignal(sig0, 8, 1))

	// should return an error because sig1 is at start bit 4
	assert.Error(nestedMuxSig.InsertSignal(sig2, 4))
	assert.Error(nestedMuxSig.InsertSignal(sig2, 4, 0))

	// should return an error because groupID is invalid
	assert.Error(nestedMuxSig.InsertSignal(sig2, 8, 512))
}

func Test_MultiplexerSignal_RemoveSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg", 1, 8)

	muxSig, err := NewMultiplexerSignal("mux_sig", 4, 16)
	assert.NoError(err)

	assert.NoError(msg.AppendSignal(muxSig))

	size4Type, err := NewIntegerSignalType("4_bits", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", size4Type)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size4Type)
	assert.NoError(err)

	// insert sig0 in groupID 0 and 2
	// gID gID | 0 0 0 0 - - - - - - - - - - - -
	assert.NoError(muxSig.InsertSignal(sig0, 0, 0, 2))

	// insert sig1 as fixed
	// gID gID | 0 0 0 0 1 1 1 1 - - - - - - - -
	assert.NoError(muxSig.InsertSignal(sig1, 4))

	assert.NoError(muxSig.RemoveSignal(sig0.EntityID()))
	assert.Len(muxSig.GetSignalGroup(0), 1)

	assert.NoError(muxSig.RemoveSignal(sig1.EntityID()))
	for _, group := range muxSig.GetSignalGroups() {
		assert.Len(group, 0)
	}

	// should return an error because it is given an invalid id
	assert.Error(muxSig.RemoveSignal("invalid-id"))

	muxSig.ClearAllSignalGroups()

	nestedMuxSig, err := NewMultiplexerSignal("nested_mux_sig", 2, 4)
	assert.NoError(err)

	// insert a nested multiplexer signal in groupID 0
	assert.NoError(muxSig.InsertSignal(nestedMuxSig, 0, 0))

	// insert sig0 and sig1 to the nestedMuxSig
	assert.NoError(nestedMuxSig.InsertSignal(sig0, 0, 0))
	assert.NoError(nestedMuxSig.InsertSignal(sig1, 0, 1))

	assert.NoError(nestedMuxSig.RemoveSignal(sig0.EntityID()))
	assert.Len(nestedMuxSig.GetSignalGroup(0), 0)

	// remove nestedMuxSig
	assert.NoError(muxSig.RemoveSignal(nestedMuxSig.EntityID()))
	assert.Equal(1, msg.signalNames.size())
}

func Test_MultiplexerSignal_ClearSignalGroup(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg", 1, 8)

	muxSig, err := NewMultiplexerSignal("mux_sig", 4, 16)
	assert.NoError(err)

	assert.NoError(msg.AppendSignal(muxSig))

	size4Type, err := NewIntegerSignalType("4_bits", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", size4Type)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size4Type)
	assert.NoError(err)

	// insert sig0 in groupID 0 and 2
	// gID gID | 0 0 0 0 - - - - - - - - - - - -
	assert.NoError(muxSig.InsertSignal(sig0, 0, 0, 2))

	// insert sig1 as fixed
	// gID gID | 0 0 0 0 1 1 1 1 - - - - - - - -
	assert.NoError(muxSig.InsertSignal(sig1, 4))

	assert.NoError(muxSig.ClearSignalGroup(2))
	assert.Len(muxSig.GetSignalGroup(2), 1)

	assert.NoError(muxSig.ClearSignalGroup(0))
	assert.Len(muxSig.GetSignalGroup(0), 1)

	// should return an error beacause groupID is invalid
	assert.Error(muxSig.ClearSignalGroup(512))
}

func Test_MultiplexerSignal_ClearAllSignalGroups(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg", 1, 8)

	muxSig, err := NewMultiplexerSignal("mux_sig", 4, 16)
	assert.NoError(err)

	assert.NoError(msg.AppendSignal(muxSig))

	size4Type, err := NewIntegerSignalType("4_bits", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", size4Type)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", size4Type)
	assert.NoError(err)

	// insert sig0 in groupID 0 and 2
	// gID gID | 0 0 0 0 - - - - - - - - - - - -
	assert.NoError(muxSig.InsertSignal(sig0, 0, 0, 2))

	// insert sig1 as fixed
	// gID gID | 0 0 0 0 1 1 1 1 - - - - - - - -
	assert.NoError(muxSig.InsertSignal(sig1, 4))

	muxSig.ClearAllSignalGroups()

	for _, group := range muxSig.GetSignalGroups() {
		assert.Len(group, 0)
	}
}

package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StandardSignal_UpdateType(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("message", "", 2)

	sigTypeInt4 := NewSignalType("int4", "", SignalTypeKindInteger, 4, true, -8, 7)

	sig0, err := NewStandardSignal("signal_0", "", sigTypeInt4, -8, 7, 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("signal_1", "", sigTypeInt4, -8, 7, 0, 1, nil)
	assert.NoError(err)
	sig2, err := NewStandardSignal("signal_2", "", sigTypeInt4, -8, 7, 0, 1, nil)
	assert.NoError(err)

	assert.NoError(msg.AppendSignal(sig0))
	assert.NoError(msg.AppendSignal(sig1))
	assert.NoError(msg.AppendSignal(sig2))

	sigTypeInt8 := NewSignalType("int8", "", SignalTypeKindInteger, 8, true, -128, 127)

	assert.NoError(sig0.UpdateType(sigTypeInt8))

	assert.Error(sig1.UpdateType(sigTypeInt8))

	assert.NoError(sig0.UpdateType(sigTypeInt4))

	assert.NoError(sig2.UpdateType(sigTypeInt8))
	assert.NoError(sig2.UpdateType(sigTypeInt8))

	assert.Error(sig1.UpdateType(sigTypeInt8))

	assert.NoError(sig2.UpdateType(sigTypeInt4))

	assert.NoError(sig1.UpdateType(sigTypeInt8))
	assert.Error(sig2.UpdateType(sigTypeInt8))

	t.Log(msg.String())
}

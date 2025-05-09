package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SignalType(t *testing.T) {
	assert := assert.New(t)

	sigType, err := NewDecimalSignalType("signal_type", 16, true)
	assert.NoError(err)
	sigType.SetScale(0.5)
	assert.Equal(0.5, sigType.Scale())
	assert.Equal(float64(-16384), sigType.Min())
	assert.Equal(float64(16383.5), sigType.Max())

	sigType.SetOffset(100.5)
	assert.Equal(100.5, sigType.Offset())
	assert.Equal(float64(-16283.5), sigType.Min())
	assert.Equal(float64(16484), sigType.Max())
}

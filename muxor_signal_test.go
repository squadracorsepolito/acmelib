package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MuxorSignal_UpdateLayoutCount(t *testing.T) {
	assert := assert.New(t)

	simpleMuxMsg := initSimpleMuxMessage(assert)

	sigMuxor := simpleMuxMsg.layer.Muxor()
	assert.NoError(sigMuxor.UpdateLayoutCount(200))
	assert.Len(simpleMuxMsg.layer.Layouts(), 200)

	assert.Error(sigMuxor.UpdateLayoutCount(257))
	assert.Error(sigMuxor.UpdateLayoutCount(2))
}

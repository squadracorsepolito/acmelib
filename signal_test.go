package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_signal_UpdateName(t *testing.T) {
	assert := assert.New(t)

	muxMsg := initMultiplexedMessage(assert)

	sigBase := muxMsg.signals.base
	assert.NoError(sigBase.UpdateName("new_base_signal"))
	assert.Equal("new_base_signal", sigBase.Name())
	assert.NoError(sigBase.UpdateName("base_signal"))
	assert.Error(sigBase.UpdateName("top_signal_in_0"))
	assert.Error(sigBase.UpdateName("top_inner_signal_in_0"))
	assert.Error(sigBase.UpdateName("bottom_signal_in_0"))
	assert.Error(sigBase.UpdateName("bottom_inner_signal_in_0"))

	sigTopIn0 := muxMsg.layers.top.signals.in0
	assert.NoError(sigTopIn0.UpdateName("new_top_signal_in_0"))
	assert.Equal("new_top_signal_in_0", sigTopIn0.Name())
	assert.NoError(sigTopIn0.UpdateName("top_signal_in_0"))
	assert.Error(sigTopIn0.UpdateName("base_signal"))
	assert.Error(sigTopIn0.UpdateName("top_signal_in_255"))
	assert.Error(sigTopIn0.UpdateName("top_inner_signal_in_0"))
	assert.Error(sigTopIn0.UpdateName("bottom_signal_in_0"))
	assert.Error(sigTopIn0.UpdateName("bottom_inner_signal_in_0"))

	sigToInnerIn0 := muxMsg.layers.top.inner.signals.in0
	assert.NoError(sigToInnerIn0.UpdateName("new_top_inner_signal_in_0"))
	assert.Equal("new_top_inner_signal_in_0", sigToInnerIn0.Name())
	assert.NoError(sigToInnerIn0.UpdateName("top_inner_signal_in_0"))
	assert.Error(sigToInnerIn0.UpdateName("base_signal"))
	assert.Error(sigToInnerIn0.UpdateName("top_signal_in_0"))
	assert.Error(sigToInnerIn0.UpdateName("top_inner_signal_in_255"))
	assert.Error(sigToInnerIn0.UpdateName("bottom_signal_in_0"))
	assert.Error(sigToInnerIn0.UpdateName("bottom_inner_signal_in_0"))

	sigTopMuxor := muxMsg.layers.top.layer.Muxor()
	assert.NoError(sigTopMuxor.UpdateName("new_top_muxor"))
	assert.Equal("new_top_muxor", muxMsg.layers.top.layer.Muxor().Name())
	assert.NoError(sigTopMuxor.UpdateName("top_muxor"))
	assert.Error(sigTopMuxor.UpdateName("top_inner_muxor"))
	assert.Error(sigTopMuxor.UpdateName("bottom_muxor"))
	assert.Error(sigTopMuxor.UpdateName("bottom_inner_muxor"))

	sigTopInnerMuxor := muxMsg.layers.top.inner.layer.Muxor()
	assert.NoError(sigTopInnerMuxor.UpdateName("new_top_inner_muxor"))
	assert.Equal("new_top_inner_muxor", sigTopInnerMuxor.Name())
	assert.NoError(sigTopInnerMuxor.UpdateName("top_inner_muxor"))
	assert.Error(sigTopInnerMuxor.UpdateName("top_muxor"))
	assert.Error(sigTopInnerMuxor.UpdateName("bottom_muxor"))
	assert.Error(sigTopInnerMuxor.UpdateName("bottom_inner_muxor"))
}

func Test_signal_UpdateStartPos(t *testing.T) {
	assert := assert.New(t)

	muxMsg := initMultiplexedMessage(assert)

	sigBase := muxMsg.signals.base
	assert.Error(sigBase.UpdateStartPos(0))
	assert.Error(sigBase.UpdateStartPos(8))

	assert.NoError(muxMsg.layers.top.signals.in255.UpdateStartPos(16))
	assert.Equal(16, muxMsg.layers.top.signals.in255.GetStartPos())
	assert.NoError(muxMsg.layers.top.signals.in255.UpdateStartPos(8))

	sigTopIn02 := muxMsg.layers.top.signals.in02
	assert.Error(sigTopIn02.UpdateStartPos(8))

	assert.Error(muxMsg.layers.top.inner.signals.in0.UpdateStartPos(58))

	sigTopMuxor := muxMsg.layers.top.layer.Muxor()
	assert.Error(sigTopMuxor.UpdateStartPos(8))
	assert.Error(sigTopMuxor.UpdateStartPos(16))
	assert.Error(sigTopMuxor.UpdateStartPos(24))

	simpleMuxMsg := initSimpleMuxMessage(assert)
	sigMuxor := simpleMuxMsg.layer.Muxor()
	assert.NoError(sigMuxor.UpdateStartPos(16))
	assert.Equal(16, sigMuxor.GetStartPos())
	assert.NoError(sigMuxor.UpdateStartPos(0))
	assert.Error(sigMuxor.UpdateStartPos(8))
}

func Test_signal_updateSize(t *testing.T) {
	assert := assert.New(t)

	muxMsg := initMultiplexedMessage(assert)

	sigBase := muxMsg.signals.base
	assert.NoError(sigBase.verifyAndUpdateSize(15))
	assert.Equal(15, sigBase.GetSize())
	assert.Error(sigBase.verifyAndUpdateSize(17))
	assert.NoError(sigBase.verifyAndUpdateSize(16))

	sigTopIn0 := muxMsg.layers.top.signals.in0
	assert.NoError(sigTopIn0.verifyAndUpdateSize(7))
	assert.Equal(7, sigTopIn0.GetSize())
	assert.Error(sigTopIn0.verifyAndUpdateSize(9))
	assert.NoError(sigTopIn0.verifyAndUpdateSize(8))

	sigTopIn255 := muxMsg.layers.top.signals.in255
	assert.NoError(sigTopIn255.verifyAndUpdateSize(16))
	assert.Equal(16, sigTopIn255.GetSize())
	assert.Error(sigTopIn255.verifyAndUpdateSize(17))

	sigTopInnerIn0 := muxMsg.layers.top.inner.signals.in0
	assert.NoError(sigTopInnerIn0.verifyAndUpdateSize(7))
	assert.Equal(7, sigTopInnerIn0.GetSize())
	assert.Error(sigTopInnerIn0.verifyAndUpdateSize(9))
	assert.NoError(sigTopInnerIn0.verifyAndUpdateSize(8))

	sigTopMuxor := muxMsg.layers.top.layer.Muxor()
	assert.NoError(sigTopMuxor.verifyAndUpdateSize(3))
	assert.Equal(3, sigTopMuxor.GetSize())
	assert.NoError(sigTopMuxor.verifyAndUpdateSize(8))
}

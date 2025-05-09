package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MultiplexedLayer_InsertDeleteSignal(t *testing.T) {
	assert := assert.New(t)

	muxMsg := initMuxMessage(assert)

	mlTop := muxMsg.layers.top.layer
	layoutIDs := []int{3, 4, 5}
	assert.NoError(mlTop.InsertSignal(dummySignal, 8, layoutIDs...))
	for _, lID := range layoutIDs {
		assert.Len(mlTop.GetLayout(lID).Signals(), 1)
	}
	assert.NoError(mlTop.DeleteSignal(dummySignal.EntityID()))
	for _, lID := range layoutIDs {
		assert.Len(mlTop.GetLayout(lID).Signals(), 0)
	}

	// Wrong layout IDs
	assert.Error(mlTop.InsertSignal(dummySignal, 8, -1))
	assert.Error(mlTop.InsertSignal(dummySignal, 8, 256))

	// Wrong signal name
	dummySignal.UpdateName("top_signal_in_0")
	assert.Error(mlTop.InsertSignal(dummySignal, 8, 3))
	dummySignal.UpdateName("dummy_signal")

	// Intersects with existing signal in the same layer
	assert.Error(mlTop.InsertSignal(dummySignal, 8, 0))
	assert.Error(mlTop.InsertSignal(dummySignal, 8, 1))

	// Intersects with existing signal in parent layer (message)
	assert.Error(mlTop.InsertSignal(dummySignal, 0, 3))
	assert.Error(mlTop.InsertSignal(dummySignal, 24, 3))
	assert.Error(mlTop.InsertSignal(dummySignal, 56, 3))

	// Intersects with existing signal in child layer (top inner)
	assert.Error(mlTop.InsertSignal(dummySignal, 16, 1))

	// Intersects with existing signal in sibling layer (bottom)
	assert.Error(mlTop.InsertSignal(dummySignal, 48, 3))

	// Intersects with existing signal in sibling child layer (bottom inner)
	assert.Error(mlTop.InsertSignal(dummySignal, 40, 3))

}

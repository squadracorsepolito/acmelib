package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalculateBusLoad(t *testing.T) {
	assert := assert.New(t)

	bus := NewBus("bus")
	bus.SetBaudrate(250_000)

	node := NewNode("node", 1, 1)
	nodeInt := node.Interfaces()[0]
	assert.NoError(bus.AddNodeInterface(nodeInt))

	msg0 := NewMessage("msg_0", 1, 8)
	msg0.SetCycleTime(100)
	assert.NoError(nodeInt.AddSentMessage(msg0))
	msg1 := NewMessage("msg_1", 2, 8)
	msg1.SetCycleTime(10)
	assert.NoError(nodeInt.AddSentMessage(msg1))

	load, msgLoads, err := CalculateBusLoad(bus, 500)
	assert.NoError(err)
	assert.Greater(load, 5.0)
	assert.Less(load, 6.0)
	assert.Len(msgLoads, 2)

	expectedNames := []string{"msg_1", "msg_0"}
	expectedBitsPerSec := []float64{13200, 1320}
	for idx, msgLoad := range msgLoads {
		assert.Equal(expectedNames[idx], msgLoad.Message.Name())
		assert.Equal(expectedBitsPerSec[idx], msgLoad.BitsPerSec)
	}

	argErr := &ArgumentError{}

	_, _, err = CalculateBusLoad(bus, 0)
	assert.ErrorAs(err, &argErr)
	assert.ErrorAs(err, &ErrIsZero)

	_, _, err = CalculateBusLoad(bus, -1)
	assert.ErrorAs(err, &argErr)
	assert.ErrorAs(err, &ErrIsNegative)
}

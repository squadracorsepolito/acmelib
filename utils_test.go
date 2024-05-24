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
	msg0.SetCycleTime(10)
	assert.NoError(nodeInt.AddMessage(msg0))
	msg1 := NewMessage("msg_1", 2, 8)
	msg1.SetCycleTime(10)
	assert.NoError(nodeInt.AddMessage(msg1))

	load, err := CalculateBusLoad(bus, 500)
	assert.NoError(err)
	assert.Greater(load, 10.0)
	assert.Less(load, 11.0)

	argErr := &ArgumentError{}

	_, err = CalculateBusLoad(bus, 0)
	assert.ErrorAs(err, &argErr)
	assert.ErrorAs(err, &ErrIsZero)

	_, err = CalculateBusLoad(bus, -1)
	assert.ErrorAs(err, &argErr)
	assert.ErrorAs(err, &ErrIsNegative)
}

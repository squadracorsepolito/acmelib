package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CANIDBuilder(t *testing.T) {
	assert := assert.New(t)

	b := NewCANIDBuilder("canid_builder")
	b.UseMessagePriority(30).UseMessageID(4, 10).UseNodeID(0, 4)

	msgPriority := MessagePriorityLow
	msgID := MessageID(0b1111111111)
	nodeID := NodeID(0b11)

	expected := uint32(msgPriority << 30)
	expected |= uint32(msgID << 4)
	expected |= uint32(nodeID)

	res := b.Calculate(msgPriority, msgID, nodeID)

	assert.Equal(expected, uint32(res))
}

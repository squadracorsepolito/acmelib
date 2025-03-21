package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CANIDBuilder(t *testing.T) {
	assert := assert.New(t)

	b0 := NewCANIDBuilder("canid_builder_0")
	b0.UseMessagePriority(30).UseMessageID(4, 10).UseNodeID(0, 4)

	msgPriority := MessagePriorityLow
	msgID := MessageID(0b1111111111)
	nodeID := NodeID(0b11)

	expected := uint32(msgPriority << 30)
	expected |= uint32(msgID << 4)
	expected |= uint32(nodeID)

	res0 := b0.Calculate(msgPriority, msgID, nodeID)

	assert.Equal(expected, uint32(res0))

	b1 := NewCANIDBuilder("canid_builder_1")
	b1.UseMessageID(4, 10).UseCAN2A()

	expected = uint32(msgID<<4) & 0b11111111111

	res1 := b1.Calculate(msgPriority, msgID, nodeID)

	assert.Equal(expected, uint32(res1))
}

func Test_CANIDBuilder_RemoveOperation(t *testing.T) {
	assert := assert.New(t)

	b := NewCANIDBuilder("canid_builder")
	b.UseMessagePriority(30).UseMessageID(4, 10).UseNodeID(0, 4)

	assert.Len(b.Operations(), 3)

	assert.NoError(b.RemoveOperation(2))

	ops := b.Operations()
	assert.Len(ops, 2)
	assert.Equal(CANIDBuilderOpKindMessagePriority, ops[0].kind)
	assert.Equal(CANIDBuilderOpKindMessageID, ops[1].kind)

	assert.Error(b.RemoveOperation(-1))
	assert.Len(b.Operations(), 2)
}

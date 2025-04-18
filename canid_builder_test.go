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

func Test_CANIDBuilder_CalculatePartials(t *testing.T) {
	assert := assert.New(t)

	b := NewCANIDBuilder("canid_builder")
	b.UseMessagePriority(30).UseMessageID(4, 10).UseNodeID(0, 4)

	msgPriority := MessagePriorityLow
	msgID := MessageID(0b1111111111)
	nodeID := NodeID(0b11)

	expected0 := uint32(msgPriority << 30)
	expected1 := expected0 | uint32(msgID<<4)
	expected2 := expected1 | uint32(nodeID)

	expected := []uint32{expected0, expected1, expected2}
	for idx, canID := range b.CalculatePartials(msgPriority, msgID, nodeID) {
		assert.Equal(expected[idx], uint32(canID))
	}
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

func Test_CANIDBuilder_InsertOperation(t *testing.T) {
	assert := assert.New(t)

	b := NewCANIDBuilder("canid_builder")

	assert.NoError(b.InsertOperation(CANIDBuilderOpKindBitMask, 0, 11, 0))
	assert.NoError(b.InsertOperation(CANIDBuilderOpKindNodeID, 0, 4, 0))
	assert.NoError(b.InsertOperation(CANIDBuilderOpKindMessageID, 3, 7, 1))

	expectedKinds := []CANIDBuilderOpKind{CANIDBuilderOpKindNodeID, CANIDBuilderOpKindMessageID, CANIDBuilderOpKindBitMask}

	operations := b.Operations()
	assert.Len(operations, 3)
	for idx, op := range operations {
		assert.Equal(expectedKinds[idx], op.Kind())
	}

	assert.Error(b.InsertOperation(CANIDBuilderOpKindBitMask, -1, 1, 0))
	assert.Error(b.InsertOperation(CANIDBuilderOpKindBitMask, 32, 1, 0))

	assert.Error(b.InsertOperation(CANIDBuilderOpKindBitMask, 0, -1, 0))
	assert.Error(b.InsertOperation(CANIDBuilderOpKindBitMask, 31, 2, 0))

	assert.Error(b.InsertOperation(CANIDBuilderOpKindBitMask, 0, 1, -1))
	assert.Error(b.InsertOperation(CANIDBuilderOpKindBitMask, 0, 1, 4))
}

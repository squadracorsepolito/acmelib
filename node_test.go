package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Node_AddMessage(t *testing.T) {
	assert := assert.New(t)

	node := NewNode("node", 0).AddInterface()

	msg0 := NewMessage("msg_0", 1, 1)
	msg1 := NewMessage("msg_1", 2, 1)
	msg2 := NewMessage("msg_2", 3, 1)

	// should add msg0, msg1, and msg2 without errors
	assert.NoError(node.AddMessage(msg0))
	assert.NoError(node.AddMessage(msg1))
	assert.NoError(node.AddMessage(msg2))
	expectedIDs := []MessageID{1, 2, 3}
	expectedNames := []string{"msg_0", "msg_1", "msg_2"}
	for idx, tmpMsg := range node.Messages() {
		assert.Equal(expectedIDs[idx], tmpMsg.ID())
		assert.Equal(expectedNames[idx], tmpMsg.Name())
	}

	// should return an error because id 3 is already taken
	dupIDMsg := NewMessage("", 3, 1)
	assert.Error(node.AddMessage(dupIDMsg))

	// should return an error because name msg_2 is already taken
	dupNameMsg := NewMessage("msg_2", 4, 1)
	assert.Error(node.AddMessage(dupNameMsg))
}

func Test_Node_RemoveMessage(t *testing.T) {
	assert := assert.New(t)

	node := NewNode("node", 0).AddInterface()

	msg0 := NewMessage("msg_0", 1, 1)
	msg1 := NewMessage("msg_1", 2, 1)
	msg2 := NewMessage("msg_2", 3, 1)
	msg3 := NewMessage("msg_3", 4, 1)

	assert.NoError(node.AddMessage(msg0))
	assert.NoError(node.AddMessage(msg1))
	assert.NoError(node.AddMessage(msg2))
	assert.NoError(node.AddMessage(msg3))

	// should be able to remove msg1 and to cause the other ids to re-generate with the exeption of msg3
	assert.NoError(node.RemoveMessage(msg1.EntityID()))
	expectedIDs := []MessageID{1, 3, 4}
	expectedNames := []string{"msg_0", "msg_2", "msg_3"}
	for idx, tmpMsg := range node.Messages() {
		assert.Equal(expectedIDs[idx], tmpMsg.ID())
		assert.Equal(expectedNames[idx], tmpMsg.Name())
	}

	// should return an error because the entity id is invalid
	assert.Error(node.RemoveMessage("dummy-id"))
}

func Test_Node_RemoveAllMessages(t *testing.T) {
	assert := assert.New(t)

	node := NewNode("node", 0).AddInterface()

	msg0 := NewMessage("msg_0", 1, 1)
	msg1 := NewMessage("msg_1", 2, 1)
	msg2 := NewMessage("msg_2", 3, 1)
	msg3 := NewMessage("msg_3", 4, 1)

	assert.NoError(node.AddMessage(msg0))
	assert.NoError(node.AddMessage(msg1))
	assert.NoError(node.AddMessage(msg2))
	assert.NoError(node.AddMessage(msg3))

	node.RemoveAllMessages()

	assert.Equal(0, len(node.Messages()))
}

// func Test_Node_UpdateName(t *testing.T) {
// 	assert := assert.New(t)

// 	bus := NewBus("bus")

// 	node0 := NewNode("node_0", 0).AddInterface()
// 	node1 := NewNode("node_1", 1).AddInterface()

// 	assert.NoError(bus.AddNode(node0))
// 	assert.NoError(bus.AddNode(node1))

// 	// should change the name to node_00
// 	assert.NoError(node0.UpdateName("node_00"))
// 	assert.Equal("node_00", node0.Name())

// 	// should not change the name
// 	assert.NoError(node1.UpdateName("node_1"))
// 	assert.Equal("node_1", node1.Name())

// 	// should return an error because node_00 is already taken
// 	assert.Error(node1.UpdateName("node_00"))
// }

package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Node_AddSentMessage(t *testing.T) {
	assert := assert.New(t)

	node := NewNode("node", 0, 1).Interfaces()[0]

	msg0 := NewMessage("msg_0", 1, 1)
	msg1 := NewMessage("msg_1", 2, 1)
	msg2 := NewMessage("msg_2", 3, 1)

	// should add msg0, msg1, and msg2 without errors
	assert.NoError(node.AddSentMessage(msg0))
	assert.NoError(node.AddSentMessage(msg1))
	assert.NoError(node.AddSentMessage(msg2))
	expectedIDs := []MessageID{1, 2, 3}
	expectedNames := []string{"msg_0", "msg_1", "msg_2"}
	for idx, tmpMsg := range node.SentMessages() {
		assert.Equal(expectedIDs[idx], tmpMsg.ID())
		assert.Equal(expectedNames[idx], tmpMsg.Name())
	}

	// should return an error because id 3 is already taken
	dupIDMsg := NewMessage("", 3, 1)
	assert.Error(node.AddSentMessage(dupIDMsg))

	// should return an error because name msg_2 is already taken
	dupNameMsg := NewMessage("msg_2", 4, 1)
	assert.Error(node.AddSentMessage(dupNameMsg))

	// create a new bus of type CAN 2.0A
	bus := NewBus("bus")
	bus.SetType(BusTypeCAN2A)

	// attach the node interface to the bus
	assert.NoError(bus.AddNodeInterface(node))

	// create a big message
	bigMsg := NewMessage("big_msg", 5, 9)

	// should return an error because the message cannot be sent over a CAN 2.0A bus
	assert.Error(node.AddSentMessage(bigMsg))
}

func Test_Node_RemoveSentMessage(t *testing.T) {
	assert := assert.New(t)

	node := NewNode("node", 0, 1).Interfaces()[0]

	msg0 := NewMessage("msg_0", 1, 1)
	msg1 := NewMessage("msg_1", 2, 1)
	msg2 := NewMessage("msg_2", 3, 1)
	msg3 := NewMessage("msg_3", 4, 1)

	assert.NoError(node.AddSentMessage(msg0))
	assert.NoError(node.AddSentMessage(msg1))
	assert.NoError(node.AddSentMessage(msg2))
	assert.NoError(node.AddSentMessage(msg3))

	// should be able to remove msg1 and to cause the other ids to re-generate with the exeption of msg3
	assert.NoError(node.RemoveSentMessage(msg1.EntityID()))
	expectedIDs := []MessageID{1, 3, 4}
	expectedNames := []string{"msg_0", "msg_2", "msg_3"}
	for idx, tmpMsg := range node.SentMessages() {
		assert.Equal(expectedIDs[idx], tmpMsg.ID())
		assert.Equal(expectedNames[idx], tmpMsg.Name())
	}

	// should return an error because the entity id is invalid
	assert.Error(node.RemoveSentMessage("dummy-id"))
}

func Test_Node_RemoveSentMessages(t *testing.T) {
	assert := assert.New(t)

	node := NewNode("node", 0, 1).Interfaces()[0]

	msg0 := NewMessage("msg_0", 1, 1)
	msg1 := NewMessage("msg_1", 2, 1)
	msg2 := NewMessage("msg_2", 3, 1)
	msg3 := NewMessage("msg_3", 4, 1)

	assert.NoError(node.AddSentMessage(msg0))
	assert.NoError(node.AddSentMessage(msg1))
	assert.NoError(node.AddSentMessage(msg2))
	assert.NoError(node.AddSentMessage(msg3))

	node.RemoveAllSentMessages()

	assert.Equal(0, len(node.SentMessages()))
}

func Test_Node_UpdateName(t *testing.T) {
	assert := assert.New(t)

	bus := NewBus("bus")

	node0 := NewNode("node_0", 0, 1)
	node1 := NewNode("node_1", 1, 1)

	assert.NoError(bus.AddNodeInterface(node0.Interfaces()[0]))
	assert.NoError(bus.AddNodeInterface(node1.Interfaces()[0]))

	// should change the name to node_00
	assert.NoError(node0.UpdateName("node_00"))
	assert.Equal("node_00", node0.Name())

	// should not change the name
	assert.NoError(node1.UpdateName("node_1"))
	assert.Equal("node_1", node1.Name())

	// should return an error because node_00 is already taken
	assert.Error(node1.UpdateName("node_00"))
}

func Test_Node_RemoveInterface(t *testing.T) {
	assert := assert.New(t)

	bus := NewBus("bus")

	// create a node with one interface
	node := NewNode("node", 1, 1)
	nodeInt0 := node.GetInterface(0)
	assert.NotNil(nodeInt0)

	// attach the first interface to the bus
	assert.NoError(bus.AddNodeInterface(nodeInt0))

	// add another interface
	node.AddInterface()
	nodeInt1 := node.GetInterface(1)
	assert.NotNil(nodeInt1)
	assert.Equal(1, nodeInt1.Number())
	assert.Len(node.Interfaces(), 2)

	// remove the first interface and check that the second interface is now the first
	assert.NoError(node.RemoveInterface(0))
	assert.Equal(0, nodeInt1.Number())
	assert.Len(bus.NodeInterfaces(), 0)
}

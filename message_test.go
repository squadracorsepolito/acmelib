package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Message(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", 1, 8)
	assert.Equal(msg.name, "msg_0")
	assert.Equal(msg.SizeByte(), 8)
}

func Test_Message_InsertDeleteSignal(t *testing.T) {
	assert := assert.New(t)

	tdBasicMsg := initBasicMessage(assert)

	msgBasic := tdBasicMsg.message
	assert.NoError(msgBasic.InsertSignal(dummySignal, 48))
	assert.Len(msgBasic.Signals(), 5)
	assert.NoError(msgBasic.DeleteSignal(dummySignal.EntityID()))
	assert.Len(msgBasic.Signals(), 4)

	assert.Error(msgBasic.InsertSignal(dummySignal, 0))
	assert.Error(msgBasic.InsertSignal(dummySignal, 47))
	assert.Error(msgBasic.InsertSignal(dummySignal, 63))

	assert.NoError(dummySignal.UpdateName("basic_signal_0"))
	assert.Error(msgBasic.InsertSignal(dummySignal, 48))
	assert.NoError(dummySignal.UpdateName("dummy_signal"))
}

func Test_Message_SetStaticCANID(t *testing.T) {
	assert := assert.New(t)

	bus := NewBus("bus")
	node1 := NewNode("node_1", 1, 1)
	nodeInt1 := node1.Interfaces()[0]
	assert.NoError(bus.AddNodeInterface(nodeInt1))

	msg1 := NewMessage("msg_1", 1, 1)
	assert.NoError(msg1.SetStaticCANID(500))
	assert.Equal(CANID(500), msg1.GetCANID())
	assert.NoError(nodeInt1.AddSentMessage(msg1))

	msg2 := NewMessage("msg_2", 2, 1)
	assert.NoError(msg2.SetStaticCANID(500))
	assert.Error(nodeInt1.AddSentMessage(msg2))
	assert.NoError(msg2.SetStaticCANID(600))
	assert.NoError(nodeInt1.AddSentMessage(msg2))

	node2 := NewNode("node_2", 2, 1)
	nodeInt2 := node2.Interfaces()[0]
	assert.NoError(bus.AddNodeInterface(nodeInt2))

	msg3 := NewMessage("msg_3", 3, 1)
	assert.NoError(msg3.SetStaticCANID(600))
	assert.Error(nodeInt2.AddSentMessage(msg3))
	assert.NoError(msg3.SetStaticCANID(700))
	assert.NoError(nodeInt1.AddSentMessage(msg3))

	node3 := NewNode("node_3", 3, 1)
	nodeInt3 := node3.Interfaces()[0]
	msg4 := NewMessage("msg_4", 4, 1)
	assert.NoError(msg4.SetStaticCANID(700))
	assert.NoError(nodeInt3.AddSentMessage(msg4))

	assert.Error(bus.AddNodeInterface(nodeInt3))
}

func Test_Message_UpdateID(t *testing.T) {
	assert := assert.New(t)

	node := NewNode("node", 1, 1)
	nodeInt := node.Interfaces()[0]

	msg0 := NewMessage("msg_0", 1, 1)
	msg1 := NewMessage("msg_1", 2, 1)

	assert.NoError(nodeInt.AddSentMessage(msg0))
	assert.NoError(nodeInt.AddSentMessage(msg1))

	// update msg0 to id 3
	assert.NoError(msg0.UpdateID(3))

	// update msg1 to id 3, this should return an error
	assert.Error(msg1.UpdateID(3))

	// update msg0 to id 1
	assert.NoError(msg0.UpdateID(1))

	// update msg0 to id 2, this should return an error
	assert.Error(msg0.UpdateID(2))

	msg2 := NewMessage("msg_2", 3, 1)
	assert.NoError(nodeInt.AddSentMessage(msg2))

	// setting msg2 with a static id
	assert.NoError(msg2.SetStaticCANID(100))

	// should update msg2 to id 3 and remove the static id 100 from the node interface
	assert.NoError(msg2.UpdateID(3))

	msg3 := NewMessage("msg_3", 4, 1)
	assert.NoError(nodeInt.AddSentMessage(msg3))

	// setting msg3 with a static id
	assert.NoError(msg3.SetStaticCANID(100))
}

func Test_Message_AddReceiver(t *testing.T) {
	assert := assert.New(t)

	node0 := NewNode("node_0", 0, 2)
	nodeInt00 := node0.Interfaces()[0]
	nodeInt01 := node0.Interfaces()[1]

	node1 := NewNode("node_1", 1, 1)
	nodeInt1 := node1.Interfaces()[0]

	node2 := NewNode("node_2", 2, 1)
	nodeInt2 := node2.Interfaces()[0]

	msg := NewMessage("msg", 1, 1)
	assert.NoError(nodeInt00.AddSentMessage(msg))

	assert.Error(msg.AddReceiver(nodeInt00))
	assert.NoError(msg.AddReceiver(nodeInt1))

	assert.Len(msg.Receivers(), 1)
	assert.Len(nodeInt1.ReceivedMessages(), 1)

	assert.NoError(msg.AddReceiver(nodeInt01))
	assert.NoError(msg.AddReceiver(nodeInt2))

	assert.Len(msg.Receivers(), 3)
	assert.Len(nodeInt01.ReceivedMessages(), 1)
	assert.Len(nodeInt2.ReceivedMessages(), 1)
}

func Test_Message_UpdateSizeByte(t *testing.T) {
	assert := assert.New(t)

	// create a message of size 1 and update to 8 -> 1 -> 8
	msg := NewMessage("msg", 1, 1)
	assert.NoError(msg.UpdateSizeByte(8))
	assert.Equal(msg.SizeByte(), 8)
	assert.NoError(msg.UpdateSizeByte(1))
	assert.Equal(msg.SizeByte(), 1)
	assert.NoError(msg.UpdateSizeByte(8))
	assert.Equal(msg.SizeByte(), 8)

	// add a 64 bits (8 bytes) signal to the message
	bigSigType, err := NewIntegerSignalType("64_bits", 64, false)
	assert.NoError(err)
	bigSig, err := NewStandardSignal("big_sig", bigSigType)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(bigSig, 0))

	// it should return an error because the 8 bytes signal
	// prevents the message size to be shrinked
	assert.Error(msg.UpdateSizeByte(1))

	// remove the 64 bits signal
	assert.NoError(msg.DeleteSignal(bigSig.EntityID()))

	// add a 32 bits (4 bytes) signal
	smallSigType, err := NewIntegerSignalType("32_bits", 32, false)
	assert.NoError(err)
	smallSig, err := NewStandardSignal("small_sig", smallSigType)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(smallSig, 0))

	// update the message size to 4 bytes
	assert.NoError(msg.UpdateSizeByte(4))

	// create a CAN2.0A bus and a node
	bus := NewBus("bus")
	bus.SetType(BusTypeCAN2A)
	node := NewNode("node", 1, 1)

	// add the message to the node
	nodeInt := node.GetInterface(0)
	assert.NotNil(nodeInt)
	assert.NoError(bus.AddNodeInterface(nodeInt))
	assert.NoError(nodeInt.AddSentMessage(msg))

	// should return an error because the message size is too big
	// for a CAN2.0A bus
	assert.Error(msg.UpdateSizeByte(9))
}

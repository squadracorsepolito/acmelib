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

func Test_Message_AppendSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", 1, 8)

	size8Type, _ := NewIntegerSignalType("8_bits", 8, false)
	size32Type, _ := NewIntegerSignalType("32_bits", 32, false)

	sigNames := []string{"sig_0", "sig_1", "sig_2", "sig_3", "sig_4"}

	sig0, _ := NewStandardSignal(sigNames[0], size8Type)
	sig1, _ := NewStandardSignal(sigNames[1], size8Type)
	sig2, _ := NewStandardSignal(sigNames[2], size8Type)
	sig3, _ := NewStandardSignal(sigNames[3], size8Type)
	sig4, _ := NewStandardSignal(sigNames[4], size32Type)

	assert.NoError(msg.AppendSignal(sig0))

	duplicatedSigName, _ := NewStandardSignal(sigNames[0], size8Type)
	assert.Error(msg.AppendSignal(duplicatedSigName))

	assert.NoError(msg.AppendSignal(sig1))
	assert.NoError(msg.AppendSignal(sig2))
	assert.NoError(msg.AppendSignal(sig3))

	assert.NoError(msg.AppendSignal(sig4))

	sigTypMassive, err := NewIntegerSignalType("massive", 128, false)
	assert.NoError(err)
	massiveSig, _ := NewStandardSignal("massive_sig", sigTypMassive)
	assert.Error(msg.AppendSignal(massiveSig))

	exidingSig, _ := NewStandardSignal("exiding_sig", size8Type)
	assert.Error(msg.AppendSignal(exidingSig))

	results := msg.Signals()
	assert.Equal(len(results), 5)
	for idx, sig := range results {
		assert.Equal(sigNames[idx], sig.Name())
	}
}

func Test_Message_InsertSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", 1, 8)

	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)
	size32Type, err := NewIntegerSignalType("32_bits", 32, false)
	assert.NoError(err)

	sigNames := []string{"sig_0", "sig_1", "sig_2", "sig_3", "sig_4"}

	sig0, _ := NewStandardSignal(sigNames[0], size8Type)
	sig1, _ := NewStandardSignal(sigNames[1], size8Type)
	sig2, _ := NewStandardSignal(sigNames[2], size8Type)
	sig3, _ := NewStandardSignal(sigNames[3], size8Type)
	sig4, _ := NewStandardSignal(sigNames[4], size32Type)

	assert.NoError(msg.InsertSignal(sig0, 0))

	assert.NoError(msg.InsertSignal(sig1, 24))

	duplicatedSigName, _ := NewStandardSignal(sigNames[0], size8Type)
	assert.Error(msg.InsertSignal(duplicatedSigName, 16))

	overlappingSig, _ := NewStandardSignal("overlapping_sig", size8Type)
	assert.Error(msg.InsertSignal(overlappingSig, 0))
	assert.Error(msg.InsertSignal(overlappingSig, 7))
	assert.Error(msg.InsertSignal(overlappingSig, 23))

	assert.NoError(msg.InsertSignal(sig2, 16))
	assert.NoError(msg.InsertSignal(sig3, 8))
	assert.NoError(msg.InsertSignal(sig4, 32))

	sigTypMassive, err := NewIntegerSignalType("massive", 128, false)
	assert.NoError(err)
	massiveSig, _ := NewStandardSignal("massive_sig", sigTypMassive)
	assert.Error(msg.InsertSignal(massiveSig, 0))

	exidingSig, _ := NewStandardSignal("exiding_sig", size8Type)
	assert.Error(msg.InsertSignal(exidingSig, 0))
	assert.Error(msg.InsertSignal(exidingSig, 64))

	correctOrder := []string{"sig_0", "sig_3", "sig_2", "sig_1", "sig_4"}

	results := msg.Signals()
	assert.Equal(len(results), 5)
	for idx, sig := range results {
		assert.Equal(correctOrder[idx], sig.Name())
	}

}

func Test_Message_RemoveSignal(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", 1, 8)

	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)
	size32Type, err := NewIntegerSignalType("32_bits", 32, false)
	assert.NoError(err)

	sigNames := []string{"sig_0", "sig_1", "sig_2", "sig_3", "sig_4"}

	sig0, _ := NewStandardSignal(sigNames[0], size8Type)
	sig1, _ := NewStandardSignal(sigNames[1], size8Type)
	sig2, _ := NewStandardSignal(sigNames[2], size8Type)
	sig3, _ := NewStandardSignal(sigNames[3], size8Type)
	sig4, _ := NewStandardSignal(sigNames[4], size32Type)

	assert.NoError(msg.AppendSignal(sig0))
	assert.NoError(msg.AppendSignal(sig1))
	assert.NoError(msg.AppendSignal(sig2))
	assert.NoError(msg.AppendSignal(sig3))
	assert.NoError(msg.AppendSignal(sig4))

	assert.NoError(msg.RemoveSignal(sig2.entityID))
	assert.Error(msg.RemoveSignal(EntityID("invalid_entity_id")))

	correctOrder := []string{"sig_0", "sig_1", "sig_3", "sig_4"}

	results := msg.Signals()
	assert.Equal(len(results), 4)
	for idx, sig := range results {
		assert.Equal(correctOrder[idx], sig.Name())
	}

}

func Test_Message_CompactSignals(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("msg_0", 1, 8)

	size8Type, err := NewIntegerSignalType("8_bits", 8, false)
	assert.NoError(err)

	sig0, _ := NewStandardSignal("sig_0", size8Type)
	sig1, _ := NewStandardSignal("sig_1", size8Type)
	sig2, _ := NewStandardSignal("sig_2", size8Type)

	assert.NoError(msg.InsertSignal(sig0, 2))
	assert.NoError(msg.InsertSignal(sig1, 18))
	assert.NoError(msg.InsertSignal(sig2, 26))

	msg.CompactSignals()

	correctStartBits := []int{0, 8, 16}

	for idx, sig := range msg.Signals() {
		assert.Equal(correctStartBits[idx], sig.GetStartBit())
	}

}

func Test_Message_ShiftSignalLeft(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("message", 1, 2)

	size4Type, err := NewIntegerSignalType("32_bits", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("signal_0", size4Type)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig0, 12))

	assert.Equal(0, msg.ShiftSignalLeft(sig0.EntityID(), 0))

	assert.Equal(1, msg.ShiftSignalLeft(sig0.EntityID(), 1))
	assert.Equal(11, msg.ShiftSignalLeft(sig0.EntityID(), 16))

	sig1, err := NewStandardSignal("signal_1", size4Type)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig1, 12))

	assert.Equal(8, msg.ShiftSignalLeft(sig1.EntityID(), 16))

	sig2, err := NewStandardSignal("signal_2", size4Type)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig2, 12))

	assert.Equal(4, msg.ShiftSignalLeft(sig2.EntityID(), 16))

	finalStartBits := []int{0, 4, 8}
	for idx, sig := range msg.Signals() {
		assert.Equal(finalStartBits[idx], sig.GetStartBit())
	}

}

func Test_Message_ShiftSignalRight(t *testing.T) {
	assert := assert.New(t)

	msg := NewMessage("message", 1, 2)

	size4Type, err := NewIntegerSignalType("32_bits", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("signal_0", size4Type)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig0, 0))

	assert.Equal(0, msg.ShiftSignalRight(sig0.EntityID(), 0))

	assert.Equal(1, msg.ShiftSignalRight(sig0.EntityID(), 1))
	assert.Equal(11, msg.ShiftSignalRight(sig0.EntityID(), 16))

	sig1, err := NewStandardSignal("signal_1", size4Type)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig1, 0))

	assert.Equal(8, msg.ShiftSignalRight(sig1.EntityID(), 16))

	sig2, err := NewStandardSignal("signal_2", size4Type)
	assert.NoError(err)
	assert.NoError(msg.InsertSignal(sig2, 0))

	assert.Equal(4, msg.ShiftSignalRight(sig2.EntityID(), 16))

	finalStartBits := []int{4, 8, 12}
	for idx, sig := range msg.Signals() {
		assert.Equal(finalStartBits[idx], sig.GetStartBit())
	}
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
	assert.NoError(msg.RemoveSignal(bigSig.EntityID()))

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
	nodeInt, err := node.GetInterface(0)
	assert.NoError(err)
	assert.NoError(bus.AddNodeInterface(nodeInt))
	assert.NoError(nodeInt.AddSentMessage(msg))

	// should return an error because the message size is too big
	// for a CAN2.0A bus
	assert.Error(msg.UpdateSizeByte(9))
}

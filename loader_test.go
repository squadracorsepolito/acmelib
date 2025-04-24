package acmelib

// import (
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func compareNetworks(assert *assert.Assertions, expected, curr *Network) {
// 	assert.Equal(expected.Name(), curr.Name())

// 	expectedBuses := expected.Buses()
// 	currBuses := curr.Buses()
// 	assert.Len(currBuses, len(expectedBuses))
// 	for i, currBus := range currBuses {
// 		expectedBus := expectedBuses[i]
// 		compareBuses(assert, expectedBus, currBus)
// 	}
// }

// func compareBuses(assert *assert.Assertions, expected, curr *Bus) {
// 	assert.Equal(expected.Name(), curr.Name())
// 	assert.Equal(expected.Type(), curr.Type())

// 	compareCANIDBuilders(assert, expected.CANIDBuilder(), curr.CANIDBuilder())

// 	expectedNodeInterfaces := expected.NodeInterfaces()
// 	currNodeInterfaces := curr.NodeInterfaces()
// 	assert.Len(currNodeInterfaces, len(expectedNodeInterfaces))
// 	for i, currNodeInterface := range currNodeInterfaces {
// 		expectedNodeInterface := expectedNodeInterfaces[i]
// 		compareNodeInterfaces(assert, expectedNodeInterface, currNodeInterface)
// 	}
// }

// func compareCANIDBuilders(assert *assert.Assertions, expected, curr *CANIDBuilder) {
// 	assert.Equal(expected.Name(), curr.Name())
// 	assert.Len(curr.Operations(), len(expected.Operations()))
// }

// func compareNodeInterfaces(assert *assert.Assertions, expected, curr *NodeInterface) {
// 	assert.Equal(expected.Number(), curr.Number())

// 	compareNode(assert, expected.Node(), curr.Node())

// 	expectedMessages := expected.SentMessages()
// 	currMessages := curr.SentMessages()
// 	assert.Len(currMessages, len(expectedMessages))
// 	for i, currMessage := range currMessages {
// 		expectedMessage := expectedMessages[i]
// 		compareMessage(assert, expectedMessage, currMessage)
// 	}
// }

// func compareNode(assert *assert.Assertions, expected, curr *Node) {
// 	assert.Equal(expected.Name(), curr.Name())
// 	assert.Equal(expected.ID(), curr.ID())
// }

// func compareMessage(assert *assert.Assertions, expected, curr *Message) {
// 	assert.Equal(expected.Name(), curr.Name())
// 	assert.Equal(expected.ID(), curr.ID())
// 	assert.Equal(expected.SizeByte(), curr.SizeByte())
// 	assert.Equal(expected.SendType(), curr.SendType())
// 	assert.Equal(expected.CycleTime(), curr.CycleTime())

// 	expectedSignals := expected.Signals()
// 	currSignals := curr.Signals()
// 	assert.Len(currSignals, len(expectedSignals))
// 	for i, currSignal := range currSignals {
// 		expectedSignal := expectedSignals[i]
// 		compareSignal(assert, expectedSignal, currSignal)
// 	}
// }

// func compareSignal(assert *assert.Assertions, expected, curr Signal) {
// 	assert.Equal(expected.Name(), curr.Name())
// 	assert.Equal(expected.GetStartBit(), curr.GetStartBit())
// 	assert.Equal(expected.GetSize(), curr.GetSize())
// }

// const expectedWireFile = "testdata/expected.binpb"
// const expectedJSONFile = "testdata/expected.json"
// const expectedTextFile = "testdata/expected.txtpb"

// func Test_LoadNetwork(t *testing.T) {
// 	assert := assert.New(t)

// 	wireFile, err := os.Open(expectedWireFile)
// 	assert.NoError(err)
// 	defer wireFile.Close()

// 	jsonFile, err := os.Open(expectedJSONFile)
// 	assert.NoError(err)
// 	defer jsonFile.Close()

// 	textFile, err := os.Open(expectedTextFile)
// 	assert.NoError(err)
// 	defer textFile.Close()

// 	expectedNet := initTestNetwork(assert)

// 	net, err := LoadNetwork(wireFile, SaveEncodingWire)
// 	assert.NoError(err)
// 	compareNetworks(assert, expectedNet, net)

// 	net, err = LoadNetwork(jsonFile, SaveEncodingJSON)
// 	assert.NoError(err)
// 	compareNetworks(assert, expectedNet, net)

// 	net, err = LoadNetwork(textFile, SaveEncodingText)
// 	assert.NoError(err)
// 	compareNetworks(assert, expectedNet, net)
// }

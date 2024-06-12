package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Bus_AddNodeInterface(t *testing.T) {
	assert := assert.New(t)

	bus := NewBus("bus")

	node0 := NewNode("node_0", 0, 1).Interfaces()[0]
	node1 := NewNode("node_1", 1, 1).Interfaces()[0]
	node2 := NewNode("node_2", 2, 1).Interfaces()[0]

	// should add node0, node1, and node2 without errors
	assert.NoError(bus.AddNodeInterface(node0))
	assert.NoError(bus.AddNodeInterface(node1))
	assert.NoError(bus.AddNodeInterface(node2))
	expectedIDs := []NodeID{0, 1, 2}
	expectedNames := []string{"node_0", "node_1", "node_2"}
	for idx, tmpNode := range bus.NodeInterfaces() {
		assert.Equal(expectedIDs[idx], tmpNode.node.ID())
		assert.Equal(expectedNames[idx], tmpNode.node.Name())
	}

	// should return an error because id 1 is already taken
	dupIDNode := NewNode("", 2, 1).Interfaces()[0]
	assert.Error(bus.AddNodeInterface(dupIDNode))

	// should return an error because name node_1 is already taken
	dupNameNode := NewNode("node_1", 3, 1).Interfaces()[0]
	assert.Error(bus.AddNodeInterface(dupNameNode))
}

func Test_Bus_RemoveNodeInterface(t *testing.T) {
	assert := assert.New(t)

	bus := NewBus("bus")

	node0 := NewNode("node_0", 0, 1).Interfaces()[0]
	node1 := NewNode("node_1", 1, 1).Interfaces()[0]
	node2 := NewNode("node_2", 2, 1).Interfaces()[0]
	node3 := NewNode("node_3", 3, 1).Interfaces()[0]

	assert.NoError(bus.AddNodeInterface(node0))
	assert.NoError(bus.AddNodeInterface(node1))
	assert.NoError(bus.AddNodeInterface(node2))
	assert.NoError(bus.AddNodeInterface(node3))

	// should remove without problems node2
	assert.NoError(bus.RemoveNodeInterface(node2.Node().EntityID()))
	expectedIDs := []NodeID{0, 1, 3}
	expectedNames := []string{"node_0", "node_1", "node_3"}
	for idx, tmpNode := range bus.NodeInterfaces() {
		assert.Equal(expectedIDs[idx], tmpNode.node.ID())
		assert.Equal(expectedNames[idx], tmpNode.node.Name())
	}

	// should return an error because the entity id is invalid
	assert.Error(bus.RemoveNodeInterface("dummy-id"))
}

func Test_Bus_RemoveAllNodeInterfaces(t *testing.T) {
	assert := assert.New(t)

	bus := NewBus("bus")

	node0 := NewNode("node_0", 0, 1).Interfaces()[0]
	node1 := NewNode("node_1", 1, 1).Interfaces()[0]
	node2 := NewNode("node_2", 2, 1).Interfaces()[0]
	node3 := NewNode("node_3", 3, 1).Interfaces()[0]

	assert.NoError(bus.AddNodeInterface(node0))
	assert.NoError(bus.AddNodeInterface(node1))
	assert.NoError(bus.AddNodeInterface(node2))
	assert.NoError(bus.AddNodeInterface(node3))

	bus.RemoveAllNodeInterfaces()

	assert.Equal(0, len(bus.NodeInterfaces()))
}

func Test_Bus_UpdateName(t *testing.T) {
	assert := assert.New(t)

	net := NewNetwork("net")

	bus0 := NewBus("bus_0")
	bus1 := NewBus("bus_1")

	assert.NoError(net.AddBus(bus0))
	assert.NoError(net.AddBus(bus1))

	// should change the name to bus_00
	assert.NoError(bus0.UpdateName("bus_00"))
	assert.Equal("bus_00", bus0.Name())

	// should not change the name
	assert.NoError(bus1.UpdateName("bus_1"))
	assert.Equal("bus_1", bus1.Name())

	// should return an error because bus_00 is already taken
	assert.Error(bus1.UpdateName("bus_00"))
}

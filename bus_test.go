package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Bus_AddNode(t *testing.T) {
	assert := assert.New(t)

	bus := NewBus("bus")

	node0 := NewNode("node_0", 0)
	node1 := NewNode("node_1", 1)
	node2 := NewNode("node_2", 2)

	// should add node0, node1, and node2 without errors
	assert.NoError(bus.AddNode(node0))
	assert.NoError(bus.AddNode(node1))
	assert.NoError(bus.AddNode(node2))
	expectedIDs := []NodeID{0, 1, 2}
	expectedNames := []string{"node_0", "node_1", "node_2"}
	for idx, tmpNode := range bus.Nodes() {
		assert.Equal(expectedIDs[idx], tmpNode.ID())
		assert.Equal(expectedNames[idx], tmpNode.Name())
	}

	// should return an error because id 1 is already taken
	dupIDNode := NewNode("", 2)
	assert.Error(bus.AddNode(dupIDNode))

	// should return an error because name node_1 is already taken
	dupNameNode := NewNode("node_1", 3)
	assert.Error(bus.AddNode(dupNameNode))
}

func Test_Bus_RemoveNode(t *testing.T) {
	assert := assert.New(t)

	bus := NewBus("bus")

	node0 := NewNode("node_0", 0)
	node1 := NewNode("node_1", 1)
	node2 := NewNode("node_2", 2)
	node3 := NewNode("node_3", 3)

	assert.NoError(bus.AddNode(node0))
	assert.NoError(bus.AddNode(node1))
	assert.NoError(bus.AddNode(node2))
	assert.NoError(bus.AddNode(node3))

	// should remove without problems node2
	assert.NoError(bus.RemoveNode(node2.EntityID()))
	expectedIDs := []NodeID{0, 1, 3}
	expectedNames := []string{"node_0", "node_1", "node_3"}
	for idx, tmpNode := range bus.Nodes() {
		assert.Equal(expectedIDs[idx], tmpNode.ID())
		assert.Equal(expectedNames[idx], tmpNode.Name())
	}

	// should return an error because the entity id is invalid
	assert.Error(bus.RemoveNode("dummy-id"))
}

func Test_Bus_RemoveAllNodes(t *testing.T) {
	assert := assert.New(t)

	bus := NewBus("bus")

	node0 := NewNode("node_0", 0)
	node1 := NewNode("node_1", 1)
	node2 := NewNode("node_2", 2)
	node3 := NewNode("node_3", 3)

	assert.NoError(bus.AddNode(node0))
	assert.NoError(bus.AddNode(node1))
	assert.NoError(bus.AddNode(node2))
	assert.NoError(bus.AddNode(node3))

	bus.RemoveAllNodes()

	assert.Equal(0, len(bus.Nodes()))
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

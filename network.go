package acmelib

import (
	"slices"
	"strings"

	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// Network is the highest level entity in the package.
// Its main purpose is to hold all buses belonging to the same network.
// For example, a car can be seen as a network with multiple buses that
// are serving different areas or ECUs in the vehicle.
type Network struct {
	*entity

	buses    *collection.Map[EntityID, *Bus]
	busNames *collection.Map[string, EntityID]
}

func newNetworkFromEntity(ent *entity) *Network {
	return &Network{
		entity: ent,

		buses:    collection.NewMap[EntityID, *Bus](),
		busNames: collection.NewMap[string, EntityID](),
	}
}

// NewNetwork returns a new [Network] with the given name.
func NewNetwork(name string) *Network {
	return newNetworkFromEntity(newEntity(name, EntityKindNetwork))
}

func (n *Network) errorf(err error) error {
	return &EntityError{
		Kind:     EntityKindNetwork,
		EntityID: n.entityID,
		Name:     n.name,
		Err:      err,
	}
}

func (n *Network) verifyBusName(name string) error {
	if n.busNames.Has(name) {
		return newNameError(name, ErrIsDuplicated)
	}
	return nil
}

func (n *Network) String() string {
	s := stringer.New()

	s.Write("network:\n")

	n.entity.stringify(s)

	if n.buses.Size() == 0 {
		return s.String()
	}

	s.Write("buses:\n")
	s.Indent()
	for _, bus := range n.Buses() {
		bus.stringify(s)
	}
	s.Unindent()

	return s.String()
}

// AddBus adds a [Bus] to the [Network].
// It may return an error if the bus name is already taken.
func (n *Network) AddBus(bus *Bus) error {
	if bus == nil {
		return newArgError("bus", ErrIsNil)
	}

	if err := n.verifyBusName(bus.name); err != nil {
		return n.errorf(err)
	}

	n.buses.Set(bus.entityID, bus)
	n.busNames.Set(bus.name, bus.entityID)

	bus.setParentNetwork(n)

	return nil
}

// RemoveBus removes a [Bus] that matches the given entity id from the [Network].
// It may return an error if the bus with the given entity id is not part of the network.
func (n *Network) RemoveBus(busEntityID EntityID) error {
	bus, ok := n.buses.Get(busEntityID)
	if !ok {
		return ErrNotFound
	}

	bus.setParentNetwork(nil)

	n.buses.Delete(busEntityID)
	n.busNames.Delete(bus.name)

	return nil
}

// RemoveAllBuses removes all [Bus]es from the [Network].
func (n *Network) RemoveAllBuses() {
	for tmpBus := range n.buses.Values() {
		tmpBus.setParentNetwork(nil)
	}

	n.buses.Clear()
	n.busNames.Clear()
}

// Buses returns a slice of all [Bus]es in the [Network] sorted by name.
func (n *Network) Buses() []*Bus {
	busSlice := slices.Collect(n.buses.Values())
	slices.SortFunc(busSlice, func(a, b *Bus) int {
		return strings.Compare(a.name, b.name)
	})
	return busSlice
}

// ToNetwork returns the network itself.
func (n *Network) ToNetwork() (*Network, error) {
	return n, nil
}

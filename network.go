package acmelib

import (
	"fmt"
	"slices"
	"strings"
)

type Network struct {
	*entity

	buses    *set[EntityID, *Bus]
	busNames *set[string, EntityID]
}

// NewNetwork returns a new [Network] with the given name and description.
func NewNetwork(name, desc string) *Network {
	return &Network{
		entity: newEntity(name, desc),

		buses:    newSet[EntityID, *Bus]("bus"),
		busNames: newSet[string, EntityID]("bus name"),
	}
}

func (p *Network) errorf(err error) error {
	return fmt.Errorf(`project "%s" : %w`, p.name, err)
}

func (p *Network) modifyBusName(busEntID EntityID, newName string) {
	bus, err := p.buses.getValue(busEntID)
	if err != nil {
		panic(err)
	}

	oldName := bus.name
	p.busNames.modifyKey(oldName, newName, busEntID)
}

// UpdateName updates the name of the [Network].
func (p *Network) UpdateName(newName string) {
	p.name = newName
}

// AddBus adds a [Bus] to the [Network].
// It may return an error if the bus name is already taken.
func (p *Network) AddBus(bus *Bus) error {
	if err := p.busNames.verifyKey(bus.name); err != nil {
		return p.errorf(fmt.Errorf(`cannot add bus "%s" : %w`, bus.name, err))
	}

	p.buses.add(bus.entityID, bus)
	p.busNames.add(bus.name, bus.entityID)

	bus.setParentNetwork(p)

	return nil
}

// RemoveBus removes a [Bus] that matches the given entity id from the [Network].
// It may return an error if the bus with the given entity id is not part of the network.
func (p *Network) RemoveBus(busEntityID EntityID) error {
	bus, err := p.buses.getValue(busEntityID)
	if err != nil {
		return p.errorf(fmt.Errorf(`cannot remove bus with entity id "%s" : %w`, busEntityID, err))
	}

	bus.setParentNetwork(nil)

	p.buses.remove(busEntityID)
	p.busNames.remove(bus.name)

	return nil
}

// RemoveAllBuses removes all [Bus]es from the [Network].
func (p *Network) RemoveAllBuses() {
	for _, tmpBus := range p.buses.entries() {
		tmpBus.setParentNetwork(nil)
	}

	p.buses.clear()
	p.busNames.clear()
}

// Buses returns a slice of all [Bus]es in the [Network] sorted by name.
func (p *Network) Buses() []*Bus {
	busSlice := p.buses.getValues()
	slices.SortFunc(busSlice, func(a, b *Bus) int {
		return strings.Compare(a.name, b.name)
	})
	return busSlice
}

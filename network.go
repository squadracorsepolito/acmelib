package acmelib

import "fmt"

type Network struct {
	*entity

	buses    *set[EntityID, *Bus]
	busNames *set[string, EntityID]
}

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

func (p *Network) UpdateName(newName string) {
	p.name = newName
}

func (p *Network) AddBus(bus *Bus) error {
	if err := p.busNames.verifyKey(bus.name); err != nil {
		return p.errorf(fmt.Errorf(`cannot add bus "%s" : %w`, bus.name, err))
	}

	p.buses.add(bus.entityID, bus)
	p.busNames.add(bus.name, bus.entityID)

	bus.setParentNetwork(p)

	return nil
}

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

func (p *Network) RemoveAllBuses() {
	for _, tmpBus := range p.buses.entries() {
		tmpBus.setParentNetwork(nil)
	}

	p.buses.clear()
	p.busNames.clear()
}

func (p *Network) Buses() []*Bus {
	return p.buses.getValues()
}

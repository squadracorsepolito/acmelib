package acmelib

import "fmt"

type Project struct {
	*entity

	buses *entityCollection[*Bus, BusSortingMethod]
}

func NewProject(name, desc string) *Project {
	return &Project{
		entity: newEntity(name, desc),

		buses: newEntityCollection(newBusSorter()),
	}
}

func (p *Project) errorf(err error) error {
	return fmt.Errorf("project %s: %v", p.Name, err)
}

func (p *Project) AddBus(bus *Bus) error {
	if err := p.buses.addEntity(bus); err != nil {
		return p.errorf(err)
	}

	bus.ParentProject = p
	p.setUpdateTimeNow()

	return nil
}

func (p *Project) ListBuses(sortingMethod BusSortingMethod) []*Bus {
	return p.buses.listEntities(sortingMethod)
}

func (p *Project) RemoveBus(busID EntityID) error {
	if err := p.buses.removeEntity(busID); err != nil {
		return p.errorf(err)
	}

	p.setUpdateTimeNow()

	return nil
}

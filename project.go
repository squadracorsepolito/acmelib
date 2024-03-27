package acmelib

import "fmt"

type Project struct {
	*entity

	buses *entityCollection[*Bus]
}

func NewProject(name, desc string) *Project {
	return &Project{
		entity: newEntity(name, desc),

		buses: newEntityCollection[*Bus](),
	}
}

func (p *Project) errorf(err error) error {
	return fmt.Errorf(`project "%s": %v`, p.name, err)
}

func (p *Project) AddBus(bus *Bus) error {
	if err := p.buses.addEntity(bus); err != nil {
		return p.errorf(err)
	}

	bus.ParentProject = p

	return nil
}

func (p *Project) ListBuses() []*Bus {
	return p.buses.listEntities()
}

func (p *Project) RemoveBus(busID EntityID) error {
	if err := p.buses.removeEntity(busID); err != nil {
		return p.errorf(err)
	}

	return nil
}

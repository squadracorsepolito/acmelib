package acmelib

import (
	"slices"
	"strings"
)

// Network is the highest level entity in the package.
// Its main purpose is to hold all buses belonging to the same network.
// For example, a car can be seen as a network with multiple buses that
// are serving different areas or ECUs in the vehicle.
type Network struct {
	*entity

	buses    *set[EntityID, *Bus]
	busNames *set[string, EntityID]

	templates *set[EntityID, Template]
}

// NewNetwork returns a new [Network] with the given name.
func NewNetwork(name string) *Network {
	return &Network{
		entity: newEntity(name),

		buses:    newSet[EntityID, *Bus](),
		busNames: newSet[string, EntityID](),

		templates: newSet[EntityID, Template](),
	}
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
	err := n.busNames.verifyKeyUnique(name)
	if err != nil {
		return &NameError{
			Name: name,
			Err:  err,
		}
	}
	return nil
}

func (n *Network) String() string {
	var builder strings.Builder

	n.entity.stringify(&builder, 0)

	if n.buses.size() == 0 {
		return builder.String()
	}

	builder.WriteString("buses:\n")
	for _, bus := range n.Buses() {
		bus.stringify(&builder, 1)
		builder.WriteRune('\n')
	}

	return builder.String()
}

// UpdateName updates the name of the [Network].
func (n *Network) UpdateName(newName string) {
	n.name = newName
}

// AddBus adds a [Bus] to the [Network].
// It may return an error if the bus name is already taken.
func (n *Network) AddBus(bus *Bus) error {
	if bus == nil {
		return &ArgumentError{
			Name: "bus",
			Err:  ErrIsNil,
		}
	}

	if err := n.verifyBusName(bus.name); err != nil {
		return n.errorf(&AddEntityError{
			EntityID: bus.entityID,
			Name:     bus.name,
			Err:      err,
		})
	}

	n.buses.add(bus.entityID, bus)
	n.busNames.add(bus.name, bus.entityID)

	bus.setParentNetwork(n)

	return nil
}

// RemoveBus removes a [Bus] that matches the given entity id from the [Network].
// It may return an error if the bus with the given entity id is not part of the network.
func (n *Network) RemoveBus(busEntityID EntityID) error {
	bus, err := n.buses.getValue(busEntityID)
	if err != nil {
		return n.errorf(&RemoveEntityError{
			EntityID: busEntityID,
			Err:      err,
		})
	}

	bus.setParentNetwork(nil)

	n.buses.remove(busEntityID)
	n.busNames.remove(bus.name)

	return nil
}

// RemoveAllBuses removes all [Bus]es from the [Network].
func (n *Network) RemoveAllBuses() {
	for _, tmpBus := range n.buses.entries() {
		tmpBus.setParentNetwork(nil)
	}

	n.buses.clear()
	n.busNames.clear()
}

// Buses returns a slice of all [Bus]es in the [Network] sorted by name.
func (n *Network) Buses() []*Bus {
	busSlice := n.buses.getValues()
	slices.SortFunc(busSlice, func(a, b *Bus) int {
		return strings.Compare(a.name, b.name)
	})
	return busSlice
}

func (n *Network) AddTemplate(template Template) error {
	if template == nil {
		return &ArgumentError{
			Name: "template",
			Err:  ErrIsNil,
		}
	}

	n.templates.add(template.EntityID(), template)
	template.setParentNetwork(n)

	return nil
}

func (n *Network) RemoveTemplate(templateEntityID EntityID) error {
	tpl, err := n.templates.getValue(templateEntityID)
	if err != nil {
		return n.errorf(&RemoveEntityError{
			EntityID: templateEntityID,
			Err:      err,
		})
	}

	tpl.setParentNetwork(nil)
	n.templates.remove(templateEntityID)

	return nil
}

func (n *Network) Templates() []Template {
	templates := n.templates.getValues()
	slices.SortFunc(templates, func(a, b Template) int {
		kindA := a.TemplateKind()
		kindB := b.TemplateKind()

		if kindA == kindB {
			return strings.Compare(a.Name(), b.Name())
		}

		return int(kindA) - int(kindB)
	})
	return templates
}

func (n *Network) CANIDBuilderTemplates() []*CANIDBuilder {
	canIDBuilder := []*CANIDBuilder{}

	for _, tpl := range n.templates.getValues() {
		if tpl.TemplateKind() == TemplateKindCANIDBuilder {
			tmpCANIDBuilder, err := tpl.ToCANIDBuilder()
			if err != nil {
				panic(err)
			}
			canIDBuilder = append(canIDBuilder, tmpCANIDBuilder)
		}
	}

	return canIDBuilder
}

func (n *Network) SignaTypeTemplates() []*SignalType {
	sigTypes := []*SignalType{}

	for _, tpl := range n.templates.getValues() {
		if tpl.TemplateKind() == TemplateKindSignalType {
			tmpSigType, err := tpl.ToSignalType()
			if err != nil {
				panic(err)
			}
			sigTypes = append(sigTypes, tmpSigType)
		}
	}

	return sigTypes
}

func (n *Network) SignaUnitTemplates() []*SignalUnit {
	sigUnits := []*SignalUnit{}

	for _, tpl := range n.templates.getValues() {
		if tpl.TemplateKind() == TemplateKindSignalUnit {
			tmpSigUnit, err := tpl.ToSignalUnit()
			if err != nil {
				panic(err)
			}
			sigUnits = append(sigUnits, tmpSigUnit)
		}
	}

	return sigUnits
}

func (n *Network) SignaEnumTemplates() []*SignalEnum {
	sigEnums := []*SignalEnum{}

	for _, tpl := range n.templates.getValues() {
		if tpl.TemplateKind() == TemplateKindSignalEnum {
			tmpSigEnum, err := tpl.ToSignalEnum()
			if err != nil {
				panic(err)
			}
			sigEnums = append(sigEnums, tmpSigEnum)
		}
	}

	return sigEnums
}

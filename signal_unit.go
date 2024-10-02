package acmelib

import (
	"fmt"
	"strings"
)

// SignalUnitKind defines the kind of a [SignalUnit].
type SignalUnitKind int

const (
	// SignalUnitKindCustom defines a custom unit.
	SignalUnitKindCustom SignalUnitKind = iota
	// SignalUnitKindTemperature defines a temperature unit.
	SignalUnitKindTemperature
	// SignalUnitKindElectrical defines an electrical unit.
	SignalUnitKindElectrical
	// SignalUnitKindPower defines a power unit.
	SignalUnitKindPower
)

func (suk SignalUnitKind) String() string {
	switch suk {
	case SignalUnitKindCustom:
		return "custom"
	case SignalUnitKindTemperature:
		return "temperature"
	case SignalUnitKindElectrical:
		return "electrical"
	case SignalUnitKindPower:
		return "power"
	default:
		return "unknown"
	}
}

// SignalUnit is an entity that defines the physical unit of a [Signal].
type SignalUnit struct {
	*entity
	*withRefs[*StandardSignal]

	kind   SignalUnitKind
	symbol string
}

func newSignalUnitFromEntity(ent *entity, kind SignalUnitKind, symbol string) *SignalUnit {
	return &SignalUnit{
		entity:   ent,
		withRefs: newWithRefs[*StandardSignal](),

		kind:   kind,
		symbol: symbol,
	}
}

// Clone creates a new [SignalUnit] with the same values as the current one.
func (su *SignalUnit) Clone() *SignalUnit {
	return &SignalUnit{
		entity:   su.entity.clone(),
		withRefs: newWithRefs[*StandardSignal](),

		kind:   su.kind,
		symbol: su.symbol,
	}
}

// NewSignalUnit creates a new [SignalUnit] with the given name,
// kind, and symbol.
func NewSignalUnit(name string, kind SignalUnitKind, symbol string) *SignalUnit {
	return newSignalUnitFromEntity(newEntity(name, EntityKindSignalUnit), kind, symbol)
}

func (su *SignalUnit) stringify(b *strings.Builder, tabs int) {
	su.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%skind: %s\n", tabStr, su.kind))
	b.WriteString(fmt.Sprintf("%ssymbol: %s\n", tabStr, su.symbol))

	refCount := su.ReferenceCount()
	if refCount > 0 {
		b.WriteString(fmt.Sprintf("%sreference_count: %d\n", tabStr, refCount))
	}
}

func (su *SignalUnit) String() string {
	builder := new(strings.Builder)
	su.stringify(builder, 0)
	return builder.String()
}

// Kind returns the kind of the [SignalUnit].
func (su *SignalUnit) Kind() SignalUnitKind {
	return su.kind
}

// Symbol returns the symbol of the [SignalUnit].
func (su *SignalUnit) Symbol() string {
	return su.symbol
}

// SetName sets the name of the [SignalUnit] to the given one.
func (su *SignalUnit) SetName(name string) {
	su.name = name
}

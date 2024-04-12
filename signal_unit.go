package acmelib

import (
	"fmt"
	"strings"
)

// SignalUnitKind defines the kind of a [SignalUnit].
type SignalUnitKind string

const (
	// SignalUnitKindCustom defines a custom unit.
	SignalUnitKindCustom SignalUnitKind = "signal_unit-custom"
	// SignalUnitKindTemperature defines a temperature unit.
	SignalUnitKindTemperature SignalUnitKind = "signal_unit-temperature"
	// SignalUnitKindVoltage defines a voltage unit.
	SignalUnitKindVoltage SignalUnitKind = "signal_unit-voltage"
	// SignalUnitKindCurrent defines a current unit.
	SignalUnitKindCurrent SignalUnitKind = "signal_unit-current"
)

// SignalUnit is an entity that defines the physical unit of a [Signal].
type SignalUnit struct {
	*entity

	kind   SignalUnitKind
	symbol string
}

// NewSignalUnit creates a new [SignalUnit] with the given name,
// kind, and symbol.
func NewSignalUnit(name string, kind SignalUnitKind, symbol string) *SignalUnit {
	return &SignalUnit{
		entity: newEntity(name),

		kind:   kind,
		symbol: symbol,
	}
}

func (su *SignalUnit) stringify(b *strings.Builder, tabs int) {
	su.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%skind: %s\n", tabStr, su.kind))
	b.WriteString(fmt.Sprintf("%ssymbol: %s\n", tabStr, su.symbol))
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

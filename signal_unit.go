package acmelib

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

// NewSignalUnit creates a new [SignalUnit] with the given name, description,
// kind, and symbol.
func NewSignalUnit(name, desc string, kind SignalUnitKind, symbol string) *SignalUnit {
	return &SignalUnit{
		entity: newEntity(name, desc),

		kind:   kind,
		symbol: symbol,
	}
}

// Kind returns the kind of the [SignalUnit].
func (su *SignalUnit) Kind() SignalUnitKind {
	return su.kind
}

// Symbol returns the symbol of the [SignalUnit].
func (su *SignalUnit) Symbol() string {
	return su.symbol
}

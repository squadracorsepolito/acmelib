package acmelib

type SignalUnitKind string

const (
	SignalUnitKindCustom      SignalUnitKind = "custom"
	SignalUnitKindTemperature SignalUnitKind = "temperature"
	SignalUnitKindVoltage     SignalUnitKind = "voltage"
	SignalUnitKindCurrent     SignalUnitKind = "current"
)

type SignalUnit struct {
	*entity

	kind   SignalUnitKind
	symbol string
}

func NewSignalUnit(name, desc string, kind SignalUnitKind, symbol string) *SignalUnit {
	return &SignalUnit{
		entity: newEntity(name, desc),

		kind:   kind,
		symbol: symbol,
	}
}

func (su *SignalUnit) Kind() SignalUnitKind {
	return su.kind
}

func (su *SignalUnit) Symbol() string {
	return su.symbol
}

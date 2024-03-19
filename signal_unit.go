package acmelib

type SignalUnitKind string

const (
	SignalUnitKindCustom      SignalUnitKind = "signal_unit_custom"
	SignalUnitKindTemperature SignalUnitKind = "signal_unit_temperature"
	SignalUnitKindVoltage     SignalUnitKind = "signal_unit_voltage"
	SignalUnitKindCurrent     SignalUnitKind = "signal_unit_current"
)

type SIgnalUnit struct {
	*entity

	Kind   SignalUnitKind
	Symbol string
}

func NewSignalUnit(name, desc string, kind SignalUnitKind, symbol string) *SIgnalUnit {
	return &SIgnalUnit{
		entity: newEntity(name, desc),

		Kind:   kind,
		Symbol: symbol,
	}
}

package acmelib

type SignalTypeKind string

const (
	SignalTypeKindFlag    SignalTypeKind = "flag"
	SignalTypeKindInteger SignalTypeKind = "integer"
	SignalTypeKindFloat   SignalTypeKind = "float"
	SignalTypeKindCustom  SignalTypeKind = "custom"
)

type SignalType struct {
	*entity

	Kind   SignalTypeKind
	Size   int
	Signed bool
	Min    float64
	Max    float64
}

func NewSignalType(name, desc string, kind SignalTypeKind, size int, signed bool, min, max float64) *SignalType {
	return &SignalType{
		entity: newEntity(name, desc),

		Kind:   kind,
		Size:   size,
		Signed: signed,
		Min:    min,
		Max:    max,
	}
}

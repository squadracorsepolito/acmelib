package acmelib

import "fmt"

type SignalTypeKind string

const (
	SignalTypeKindCustom  SignalTypeKind = "custom"
	SignalTypeKindFlag    SignalTypeKind = "flag"
	SignalTypeKindInteger SignalTypeKind = "integer"
	SignalTypeKindFloat   SignalTypeKind = "float"
)

type SignalType struct {
	*entity

	kind   SignalTypeKind
	size   int
	signed bool
	min    float64
	max    float64
}

func newSignalType(name, desc string, kind SignalTypeKind, size int, signed bool, min, max float64) (*SignalType, error) {
	if size <= 0 {
		return nil, fmt.Errorf("signal type size cannot be negative")
	}

	return &SignalType{
		entity: newEntity(name, desc),

		kind:   kind,
		size:   size,
		signed: signed,
		min:    min,
		max:    max,
	}, nil
}

func NewIntegerSignalType(name, desc string, size int, signed bool) (*SignalType, error) {
	var min float64
	var max float64

	if signed {
		tmpMax := (1<<size - 1) - 1
		tmpMin := -(1<<size - 1)
		min = float64(tmpMin)
		max = float64(tmpMax)
	} else {
		tmp := (1 << size) - 1
		min = 0
		max = float64(tmp)
	}

	return newSignalType(name, desc, SignalTypeKindInteger, size, signed, min, max)
}

func (st *SignalType) Kind() SignalTypeKind {
	return st.kind
}

func (st *SignalType) Size() int {
	return st.size
}

func (st *SignalType) Signed() bool {
	return st.signed
}

func (st *SignalType) Min() float64 {
	return st.min
}

func (st *SignalType) Max() float64 {
	return st.max
}

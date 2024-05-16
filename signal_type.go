package acmelib

import (
	"fmt"
	"strings"
)

// SignalTypeKind represents the kind of a [SignalType].
type SignalTypeKind string

const (
	// SignalTypeKindCustom defines a signal of type custom.
	SignalTypeKindCustom SignalTypeKind = "signal_type-custom"
	// SignalTypeKindFlag defines a signal of type flag (1 bit).
	SignalTypeKindFlag SignalTypeKind = "signal_type-flag"
	// SignalTypeKindInteger defines a signal of type integer.
	SignalTypeKindInteger SignalTypeKind = "signal_type-integer"
	// SignalTypeKindFloat defines a signal of type float.
	SignalTypeKindFloat SignalTypeKind = "signal_type-float"
)

// SignalType is the representation of a signal type.
type SignalType struct {
	*entity
	*withTemplateRefs[*StandardSignal]

	kind   SignalTypeKind
	size   int
	signed bool
	min    float64
	max    float64
}

func newSignalType(name string, kind SignalTypeKind, size int, signed bool, min, max float64) (*SignalType, error) {
	if size < 0 {
		return nil, &ArgumentError{
			Name: "size",
			Err:  ErrIsNegative,
		}
	}

	if size == 0 {
		return nil, &ArgumentError{
			Name: "size",
			Err:  ErrIsZero,
		}
	}

	return &SignalType{
		entity:           newEntity(name),
		withTemplateRefs: newWithTemplateRefs[*StandardSignal](),

		kind:   kind,
		size:   size,
		signed: signed,
		min:    min,
		max:    max,
	}, nil
}

// NewCustomSignalType creates a new [SignalType] of kind [SignalTypeKindCustom]
// with the given name, size, signed, order, and min/max values.
// It may return an error if the size is negative.
func NewCustomSignalType(name string, size int, signed bool, min, max float64) (*SignalType, error) {
	return newSignalType(name, SignalTypeKindCustom, size, signed, min, max)
}

// NewFlagSignalType creates a new [SignalType] of kind [SignalTypeKindFlag]
// with the given name. The order is set to little endian.
func NewFlagSignalType(name string) *SignalType {
	sig, err := newSignalType(name, SignalTypeKindFlag, 1, false, 0, 1)
	if err != nil {
		panic(err)
	}
	return sig
}

// NewIntegerSignalType creates a new [SignalType] of kind [SignalTypeKindInteger]
// with the given name, size, and signed. The order is set to little endian.
// It may return an error if the size is negative.
func NewIntegerSignalType(name string, size int, signed bool) (*SignalType, error) {
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

	return newSignalType(name, SignalTypeKindInteger, size, signed, min, max)
}

// NewFloatSignalType creates a new [SignalType] of kind [SignalTypeKindFloat]
// with the given name and size. The order is set to little endian.
// It may return an error if the size is negative.
func NewFloatSignalType(name string, size int) (*SignalType, error) {
	min := (1<<size - 1) - 1
	max := -(1<<size - 1)
	return newSignalType(name, SignalTypeKindFloat, size, true, float64(min), float64(max))
}

func (st *SignalType) stringify(b *strings.Builder, tabs int) {
	st.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%skind: %s\n", tabStr, st.kind))
	b.WriteString(fmt.Sprintf("%ssize: %d; min: %g; max: %g; signed: %t\n", tabStr, st.size, st.min, st.max, st.signed))
}

func (st *SignalType) String() string {
	builder := new(strings.Builder)
	st.stringify(builder, 0)
	return builder.String()
}

// Kind returns the kind of the [SignalType].
func (st *SignalType) Kind() SignalTypeKind {
	return st.kind
}

// Size returns the size of the [SignalType].
func (st *SignalType) Size() int {
	return st.size
}

// Signed returns whether the [SignalType] is signed.
func (st *SignalType) Signed() bool {
	return st.signed
}

// Min returns the minimum value of the [SignalType].
func (st *SignalType) Min() float64 {
	return st.min
}

// Max returns the maximum value of the [SignalType].
func (st *SignalType) Max() float64 {
	return st.max
}

func (st *SignalType) ToSignalType() (*SignalType, error) {
	return st, nil
}

package acmelib

import (
	"fmt"
	"strings"
)

// SignalTypeKind represents the kind of a [SignalType].
type SignalTypeKind int

const (
	// SignalTypeKindCustom defines a signal of type custom.
	SignalTypeKindCustom SignalTypeKind = iota
	// SignalTypeKindFlag defines a signal of type flag (1 bit).
	SignalTypeKindFlag
	// SignalTypeKindInteger defines a signal of type integer.
	SignalTypeKindInteger
	// SignalTypeKindDecimal defines a signal of type float.
	SignalTypeKindDecimal
)

func (stk SignalTypeKind) String() string {
	switch stk {
	case SignalTypeKindCustom:
		return "custom"
	case SignalTypeKindFlag:
		return "flag"
	case SignalTypeKindInteger:
		return "integer"
	case SignalTypeKindDecimal:
		return "decimal"
	default:
		return "unknown"
	}
}

// SignalType is the representation of a signal type.
type SignalType struct {
	*entity
	*withRefs[*StandardSignal]

	kind   SignalTypeKind
	size   int
	signed bool
	min    float64
	max    float64
	scale  float64
	offset float64
}

func newSignalTypeFromEntity(ent *entity, kind SignalTypeKind, size int, signed bool, min, max, scale, offset float64) (*SignalType, error) {
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
		entity:   ent,
		withRefs: newWithRefs[*StandardSignal](),

		kind:   kind,
		size:   size,
		signed: signed,
		min:    min,
		max:    max,
		scale:  scale,
		offset: offset,
	}, nil
}

func newSignalType(name string, kind SignalTypeKind, size int, signed bool, min, max, scale, offset float64) (*SignalType, error) {
	return newSignalTypeFromEntity(newEntity(name, EntityKindSignalType), kind, size, signed, min, max, scale, offset)
}

// NewCustomSignalType creates a new [SignalType] of kind [SignalTypeKindCustom]
// with the given name, size, signed, order, min/max values, scale, and offset.
// It may return an error if the size is negative.
func NewCustomSignalType(name string, size int, signed bool, min, max, scale, offset float64) (*SignalType, error) {
	return newSignalType(name, SignalTypeKindCustom, size, signed, min, max, scale, offset)
}

// NewFlagSignalType creates a new [SignalType] of kind [SignalTypeKindFlag]
// with the given name.
func NewFlagSignalType(name string) *SignalType {
	sig, err := newSignalType(name, SignalTypeKindFlag, 1, false, 0, 1, 1, 0)
	if err != nil {
		panic(err)
	}
	return sig
}

// NewIntegerSignalType creates a new [SignalType] of kind [SignalTypeKindInteger]
// with the given name, size, and signed.
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

	return newSignalType(name, SignalTypeKindInteger, size, signed, min, max, 1, 0)
}

// NewDecimalSignalType creates a new [SignalType] of kind [SignalTypeKindDecimal]
// with the given name, size and signed.
// It may return an error if the size is negative.
func NewDecimalSignalType(name string, size int, signed bool) (*SignalType, error) {
	min := (1<<size - 1) - 1
	max := -(1<<size - 1)
	return newSignalType(name, SignalTypeKindDecimal, size, signed, float64(min), float64(max), 1, 0)
}

func (st *SignalType) stringify(b *strings.Builder, tabs int) {
	st.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%skind: %s\n", tabStr, st.kind))
	b.WriteString(fmt.Sprintf("%ssize: %d; signed: %t; min: %g; max: %g; scale: %g; offset: %g\n", tabStr, st.size, st.signed, st.min, st.max, st.scale, st.offset))

	refCount := st.ReferenceCount()
	if refCount > 0 {
		b.WriteString(fmt.Sprintf("%sreference_count: %d\n", tabStr, refCount))
	}
}

func (st *SignalType) String() string {
	builder := new(strings.Builder)
	st.stringify(builder, 0)
	return builder.String()
}

// SetName sets the [SignalType] name to the given one.
func (st *SignalType) SetName(name string) {
	st.name = name
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

// SetMin sets the minimum value of the [SignalType].
func (st *SignalType) SetMin(min float64) {
	st.min = min
}

// Min returns the minimum value of the [SignalType].
func (st *SignalType) Min() float64 {
	return st.min
}

// SetMax sets the maximum value of the [SignalType].
func (st *SignalType) SetMax(max float64) {
	st.max = max
}

// Max returns the maximum value of the [SignalType].
func (st *SignalType) Max() float64 {
	return st.max
}

// SetScale sets the scale of the [SignalType].
func (st *SignalType) SetScale(scale float64) {
	st.scale = scale
}

// Scale returns the scale of the [SignalType].
func (st *SignalType) Scale() float64 {
	return st.scale
}

// SetOffset sets the offset of the [SignalType].
func (st *SignalType) SetOffset(offset float64) {
	st.offset = offset
}

// Offset returns the offset of the [SignalType].
func (st *SignalType) Offset() float64 {
	return st.offset
}

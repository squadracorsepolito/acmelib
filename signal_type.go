package acmelib

import (
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// SignalTypeKind represents the kind of a [SignalType].
type SignalTypeKind int

const (
	// SignalTypeKindFlag defines a signal of type flag (1 bit).
	SignalTypeKindFlag SignalTypeKind = iota
	// SignalTypeKindInteger defines a signal of type integer.
	SignalTypeKindInteger
	// SignalTypeKindDecimal defines a signal of type float.
	SignalTypeKindDecimal
)

func (stk SignalTypeKind) String() string {
	switch stk {
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
		return nil, newArgError("size", ErrIsNegative)
	}

	if size == 0 {
		return nil, newArgError("size", ErrIsZero)
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

// Clone creates a new [SignalType] with the same values as the current one.
func (st *SignalType) Clone() *SignalType {
	return &SignalType{
		entity:   st.entity.clone(),
		withRefs: newWithRefs[*StandardSignal](),

		kind:   st.kind,
		size:   st.size,
		signed: st.signed,
		min:    st.min,
		max:    st.max,
		scale:  st.scale,
		offset: st.offset,
	}
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

func (st *SignalType) stringify(s *stringer.Stringer) {
	st.entity.stringify(s)

	s.Write("kind: %s\n", st.kind)
	s.Write("size: %d; signed: %t; min: %g; max: %g; scale: %g; offset: %g\n", st.size, st.signed, st.min, st.max, st.scale, st.offset)

	refCount := st.ReferenceCount()
	if refCount > 0 {
		s.Write("reference_count: %d\n", refCount)
	}
}

func (st *SignalType) String() string {
	s := stringer.New()
	s.Write("signal_type:\n")
	st.stringify(s)
	return s.String()
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

// UpdateSigned updates the signed flag of the [SignalType].
func (st *SignalType) UpdateSigned(signed bool) {
	if st.kind == SignalTypeKindFlag {
		return
	}

	st.signed = signed
}

// SetMin sets the minimum value of the [SignalType].
func (st *SignalType) SetMin(min float64) {
	if st.kind == SignalTypeKindFlag {
		return
	}

	st.min = min
}

// Min returns the minimum value of the [SignalType].
func (st *SignalType) Min() float64 {
	return st.min
}

// SetMax sets the maximum value of the [SignalType].
func (st *SignalType) SetMax(max float64) {
	if st.kind == SignalTypeKindFlag {
		return
	}

	st.max = max
}

// Max returns the maximum value of the [SignalType].
func (st *SignalType) Max() float64 {
	return st.max
}

// SetScale sets the scale of the [SignalType].
func (st *SignalType) SetScale(scale float64) {
	if st.kind == SignalTypeKindFlag {
		return
	}

	st.scale = scale
}

// Scale returns the scale of the [SignalType].
func (st *SignalType) Scale() float64 {
	return st.scale
}

// SetOffset sets the offset of the [SignalType].
func (st *SignalType) SetOffset(offset float64) {
	if st.kind == SignalTypeKindFlag {
		return
	}

	st.offset = offset
}

// Offset returns the offset of the [SignalType].
func (st *SignalType) Offset() float64 {
	return st.offset
}

// ToSignalType returns the type itself.
func (st *SignalType) ToSignalType() (*SignalType, error) {
	return st, nil
}

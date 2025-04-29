package acmelib

import (
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// StandardSignal is the representation of a normal signal that has a [SignalType],
// a min, a max, an offset, a scale, and can have a [SignalUnit].
type StandardSignal struct {
	*signal

	typ  *SignalType
	unit *SignalUnit
}

func newStandardSignalFromBase(base *signal, typ *SignalType) (*StandardSignal, error) {
	if typ == nil {
		return nil, newArgError("typ", ErrIsNil)
	}

	ss := &StandardSignal{
		signal: base,

		unit: nil,
	}

	ss.addType(typ)
	ss.setSize(typ.size)

	return ss, nil
}

// NewStandardSignal creates a new [StandardSignal] with the given name and [SignalType].
// It may return an error if the given [SignalType] is nil.
func NewStandardSignal(name string, typ *SignalType) (*StandardSignal, error) {
	return newStandardSignalFromBase(newSignal(name, SignalKindStandard), typ)
}

// ToStandard returns the [StandardSignal] itself.
func (ss *StandardSignal) ToStandard() (*StandardSignal, error) {
	return ss, nil
}

func (ss *StandardSignal) stringify(s *stringer.Stringer) {
	ss.signal.stringify(s)

	s.Write("type:\n")
	s.Indent()
	ss.typ.stringify(s)
	s.Unindent()

	if ss.unit != nil {
		s.Write("unit:\n")
		s.Indent()
		ss.unit.stringify(s)
		s.Unindent()
	}
}

func (ss *StandardSignal) String() string {
	s := stringer.New()
	s.Write("standard_signal:\n")
	ss.stringify(s)
	return s.String()
}

func (ss *StandardSignal) addType(typ *SignalType) {
	typ.addRef(ss)
	ss.typ = typ
}

func (ss *StandardSignal) removeType() {
	ss.typ.removeRef(ss.entityID)
	ss.typ = nil
}

// Type returns the [SignalType] of the [StandardSignal].
func (ss *StandardSignal) Type() *SignalType {
	return ss.typ
}

// UpdateType updates the [SignalType] of the signal.
//
// It returns:
//   - [ArgError] if the given type is nil.
//   - [SizeError] if the new type size cannot fit in the layout.
func (ss *StandardSignal) UpdateType(newType *SignalType) error {
	if newType == nil {
		return ss.errorf(newArgError("newType", ErrIsNil))
	}

	// Check if the new type can fit in the layout
	if err := ss.verifyAndUpdateSize(newType.size); err != nil {
		return ss.errorf(err)
	}

	// Swap the types
	ss.removeType()
	ss.addType(newType)

	return nil
}

// SetUnit sets the [SignalUnit] of the [StandardSignal] to the given one.
func (ss *StandardSignal) SetUnit(unit *SignalUnit) {
	if ss.unit != nil {
		ss.unit.removeRef(ss.entityID)
	}

	if unit == nil {
		ss.unit = nil
		return
	}

	unit.addRef(ss)
	ss.unit = unit
}

// Unit returns the [SignalUnit] of the [StandardSignal].
func (ss *StandardSignal) Unit() *SignalUnit {
	return ss.unit
}

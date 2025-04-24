package acmelib

import (
	"fmt"
	"strings"
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
		return nil, &ArgumentError{
			Name: "typ",
			Err:  ErrIsNil,
		}
	}

	sig := &StandardSignal{
		signal: base,

		typ:  typ,
		unit: nil,
	}

	typ.addRef(sig)
	sig.setSize(typ.size)

	return sig, nil
}

// NewStandardSignal creates a new [StandardSignal] with the given name and [SignalType].
// It may return an error if the given [SignalType] is nil.
func NewStandardSignal(name string, typ *SignalType) (*StandardSignal, error) {
	return newStandardSignalFromBase(newSignal(name, SignalKindStandard), typ)
}

// // GetSize returns the size of the [StandardSignal].
// func (ss *StandardSignal) GetSize() int {
// 	return ss.typ.size
// }

// ToStandard returns the [StandardSignal] itself.
func (ss *StandardSignal) ToStandard() (*StandardSignal, error) {
	return ss, nil
}

func (ss *StandardSignal) stringify(b *strings.Builder, tabs int) {
	ss.signal.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("size: %d\n", ss.GetSize()))

	b.WriteString(fmt.Sprintf("%stype:\n", tabStr))
	ss.typ.stringify(b, tabs+1)

	if ss.unit != nil {
		b.WriteString(fmt.Sprintf("%sunit:\n", tabStr))
		ss.unit.stringify(b, tabs+1)
	}
}

func (ss *StandardSignal) String() string {
	builder := new(strings.Builder)
	ss.stringify(builder, 0)
	return builder.String()
}

// Type returns the [SignalType] of the [StandardSignal].
func (ss *StandardSignal) Type() *SignalType {
	return ss.typ
}

// SetType sets the [SignalType] of the [StandardSignal].
// It resets the physical values.
// It may return an error if the given [SignalType] is nil, or if the new signal type
// size cannot fit in the message payload.
func (ss *StandardSignal) SetType(typ *SignalType) error {
	if typ == nil {
		return ss.errorf(&ArgumentError{
			Name: "typ",
			Err:  ErrIsNil,
		})
	}

	if err := ss.modifySize(typ.size - ss.typ.size); err != nil {
		return ss.errorf(err)
	}

	ss.typ.removeRef(ss.entityID)

	ss.typ = typ

	typ.addRef(ss)
	ss.setSize(typ.size)

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

func (ss *StandardSignal) GetHigh() int {
	return ss.GetStartBit() + ss.GetSize() - 1
}

// AssignAttribute assigns the given attribute/value pair to the [StandardSignal].
//
// It returns an [ArgumentError] if the attribute is nil,
// or an [AttributeValueError] if the value does not conform to the attribute.
func (ss *StandardSignal) AssignAttribute(attribute Attribute, value any) error {
	if err := ss.addAttributeAssignment(attribute, ss, value); err != nil {
		return ss.errorf(err)
	}
	return nil
}

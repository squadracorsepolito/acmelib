package acmelib

import (
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

var _ Signal = (*EnumSignal)(nil)

// EnumSignal is a signal that holds a [SignalEnum].
type EnumSignal struct {
	*signal

	enum *SignalEnum
}

func newEnumSignalFromBase(base *signal, enum *SignalEnum) (*EnumSignal, error) {
	if enum == nil {
		return nil, newArgError("enum", ErrIsNil)
	}

	es := &EnumSignal{
		signal: base,
	}

	es.addEnum(enum)
	es.setSize(enum.size)

	return es, nil
}

// NewEnumSignal creates a new [EnumSignal] with the given name and [SignalEnum].
// It may return an error if the given [SignalEnum] is nil.
func NewEnumSignal(name string, enum *SignalEnum) (*EnumSignal, error) {
	return newEnumSignalFromBase(newSignal(name, SignalKindEnum), enum)
}

// ToEnum returns the [EnumSignal] itself.
func (es *EnumSignal) ToEnum() (*EnumSignal, error) {
	return es, nil
}

func (es *EnumSignal) stringify(s *stringer.Stringer) {
	es.signal.stringify(s)

	s.Write("enum:\n")
	s.Indent()
	es.enum.stringify(s)
	s.Unindent()
}

func (es *EnumSignal) String() string {
	s := stringer.New()
	s.Write("enum_signal:\n")
	es.stringify(s)
	return s.String()
}

func (es *EnumSignal) addEnum(enum *SignalEnum) {
	enum.addRef(es)
	es.enum = enum
}

func (es *EnumSignal) removeEnum() {
	es.enum.removeRef(es.entityID)
	es.enum = nil
}

// Enum returns the [SignalEnum] of the [EnumSignal].
func (es *EnumSignal) Enum() *SignalEnum {
	return es.enum
}

// UpdateEnum updates the [EnumSignal] to use the signal.
//
// It returns:
//   - [ArgError] if the given enum is nil.
//   - [SizeError] if the new enum size cannot fit in the layout.
func (es *EnumSignal) UpdateEnum(newEnum *SignalEnum) error {
	if newEnum == nil {
		return es.errorf(newArgError("newEnum", ErrIsNil))
	}

	// Check if the new enum can fit in the layout
	if err := es.verifyAndUpdateSize(es, newEnum.size); err != nil {
		return es.errorf(err)
	}

	// Swap the enums
	es.removeEnum()
	es.addEnum(newEnum)

	return nil
}

// UpdateStartPos updates the start position of the signal.
//
// It returns a [StartPosError] if the new start position is invalid.
func (es *EnumSignal) UpdateStartPos(newStartPos int) error {
	return es.signal.updateStartPos(es, newStartPos)
}

// ToSignal returns the signal itself.
func (es *EnumSignal) ToSignal() (Signal, error) {
	return es, nil
}

// AssignAttribute assigns the given attribute/value pair to the signal.
func (es *EnumSignal) AssignAttribute(attribute Attribute, value any) error {
	return es.signal.assignAttribute(es, attribute, value)
}

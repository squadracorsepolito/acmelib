package acmelib

import (
	"fmt"
	"strings"
)

// EnumSignal is a signal that holds a [SignalEnum].
type EnumSignal struct {
	*signal

	enum *SignalEnum
}

func newEnumSignalFromBase(base *signal, enum *SignalEnum) (*EnumSignal, error) {
	if enum == nil {
		return nil, &ArgError{
			Name: "enum",
			Err:  ErrIsNil,
		}
	}

	sig := &EnumSignal{
		signal: base,

		enum: enum,
	}

	enum.addRef(sig)

	return sig, nil
}

// NewEnumSignal creates a new [EnumSignal] with the given name and [SignalEnum].
// It may return an error if the given [SignalEnum] is nil.
func NewEnumSignal(name string, enum *SignalEnum) (*EnumSignal, error) {
	return newEnumSignalFromBase(newSignal(name, SignalKindEnum), enum)
}

// GetSize returns the size of the [EnumSignal].
func (es *EnumSignal) GetSize() int {
	return es.enum.GetSize()
}

// ToEnum returns the [EnumSignal] itself.
func (es *EnumSignal) ToEnum() (*EnumSignal, error) {
	return es, nil
}

func (es *EnumSignal) stringifyOld(b *strings.Builder, tabs int) {
	es.signal.stringifyOld(b, tabs)
	b.WriteString(fmt.Sprintf("size: %d\n", es.GetSize()))

	tabStr := getTabString(tabs)
	b.WriteString(fmt.Sprintf("%senum:\n", tabStr))

	es.enum.stringifyOld(b, tabs+1)
}

func (es *EnumSignal) String() string {
	builder := new(strings.Builder)
	es.stringifyOld(builder, 0)
	return builder.String()
}

// Enum returns the [SignalEnum] of the [EnumSignal].
func (es *EnumSignal) Enum() *SignalEnum {
	return es.enum
}

// SetEnum sets the [SignalEnum] of the [EnumSignal] to the given one.
// It may return an error if the given [SignalEnum] is nil, or if the new enum
// size cannot fit in the message payload.
func (es *EnumSignal) SetEnum(enum *SignalEnum) error {
	if enum == nil {
		return es.errorf(&ArgError{
			Name: "enum",
			Err:  ErrIsNil,
		})
	}

	if err := es.modifySize(enum.GetSize() - es.GetSize()); err != nil {
		return es.errorf(err)
	}

	es.enum.removeRef(es.entityID)

	es.enum = enum

	enum.addRef(es)

	return nil
}

// AssignAttribute assigns the given attribute/value pair to the [EnumSignal].
//
// It returns an [ArgError] if the attribute is nil,
// or an [AttributeValueError] if the value does not conform to the attribute.
func (es *EnumSignal) AssignAttribute(attribute Attribute, value any) error {
	if err := es.addAttributeAssignment(attribute, es, value); err != nil {
		return es.errorf(err)
	}
	return nil
}

func (es *EnumSignal) GetHigh() int {
	return es.GetStartBit() + es.GetSize() - 1
}

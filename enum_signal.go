package acmelib

import (
	"errors"
	"fmt"
	"strings"
)

type EnumSignal struct {
	*signal

	enum *SignalEnum
}

func NewEnumSignal(name, desc string, enum *SignalEnum) (*EnumSignal, error) {
	if enum == nil {
		return nil, errors.New("signal enum cannot be nil")
	}

	sig := &EnumSignal{
		signal: newSignal(name, desc, SignalKindEnum),

		enum: enum,
	}

	enum.addSignalRef(sig)

	return sig, nil
}

func (es *EnumSignal) String() string {
	var builder strings.Builder

	builder.WriteString(es.signal.String())
	builder.WriteString(fmt.Sprintf("size: %d\n", es.GetSize()))

	builder.WriteString("\n+++END SIGNAL+++\n\n")

	return builder.String()
}

func (es *EnumSignal) GetSize() int {
	return es.enum.GetSize()
}

func (es *EnumSignal) ToStandard() (*StandardSignal, error) {
	return nil, es.errorf(errors.New(`cannot covert to "standard", the signal is of kind "enum"`))
}

func (es *EnumSignal) ToEnum() (*EnumSignal, error) {
	return es, nil
}

func (es *EnumSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return nil, es.errorf(errors.New(`cannot covert to "multiplexer", the signal is of kind "enum"`))
}

func (es *EnumSignal) GetEnum() *SignalEnum {
	return es.enum
}

func (es *EnumSignal) UpdateEnum(newEnum *SignalEnum) error {
	if err := es.modifySize(newEnum.GetSize() - es.GetSize()); err != nil {
		return es.errorf(err)
	}

	es.enum.removeSignalRef(es.GetEntityID())

	es.enum = newEnum
	es.setUpdateTimeNow()

	return nil
}

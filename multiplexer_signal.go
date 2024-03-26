package acmelib

import (
	"errors"
	"fmt"
	"strings"
)

type MultiplexerSignal struct {
	*signal

	multiplexedSignals *signalPayload
	totalSize          int
	selectSize         int
}

func NewMultiplexerSignal(name, desc string, totalSize, selectSize int) (*MultiplexerSignal, error) {
	ms := &MultiplexerSignal{
		signal: newSignal(name, desc, SignalKindMultiplexer),

		multiplexedSignals: newSignalPayload(totalSize - selectSize),
		totalSize:          totalSize,
		selectSize:         selectSize,
	}

	return ms, nil
}

func (ms *MultiplexerSignal) String() string {
	var builder strings.Builder

	builder.WriteString(ms.signal.String())
	builder.WriteString(fmt.Sprintf("size: %d\n", ms.GetSize()))

	builder.WriteString("\n+++END SIGNAL+++\n\n")

	return builder.String()
}

func (ms *MultiplexerSignal) GetSize() int {
	return ms.totalSize
}

func (ms *MultiplexerSignal) ToStandard() (*StandardSignal, error) {
	return nil, ms.errorf(errors.New(`cannot covert to "standard", the signal is of kind "multiplexer"`))
}

func (ms *MultiplexerSignal) ToEnum() (*EnumSignal, error) {
	return nil, ms.errorf(errors.New(`cannot covert to "enum", the signal is of kind "multiplexer"`))
}

func (ms *MultiplexerSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return ms, nil
}

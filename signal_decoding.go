package acmelib

import "github.com/squadracorsepolito/acmelib/internal/stringer"

// SignalValueType defines the value type of a [Signal] when decoded.
type SignalValueType string

const (
	// SignalValueTypeFlag defines a flag signal value type.
	SignalValueTypeFlag SignalValueType = "flag"
	// SignalValueTypeInt defines an integer signal value type.
	SignalValueTypeInt SignalValueType = "int"
	// SignalValueTypeUint defines an unsigned integer signal value type.
	SignalValueTypeUint SignalValueType = "uint"
	// SignalValueTypeFloat defines a float signal value type.
	SignalValueTypeFloat SignalValueType = "float"
	// SignalValueTypeEnum defines an enum signal value type.
	SignalValueTypeEnum SignalValueType = "enum"
)

func (svt SignalValueType) String() string {
	switch svt {
	case SignalValueTypeFlag:
		return "flag"
	case SignalValueTypeInt:
		return "int"
	case SignalValueTypeUint:
		return "uint"
	case SignalValueTypeFloat:
		return "float"
	case SignalValueTypeEnum:
		return "enum"
	default:
		return "unknown"
	}
}

// SignalDecoding represents a signal when decoded.
type SignalDecoding struct {
	Signal    Signal
	RawValue  uint64
	ValueType SignalValueType
	Value     any
	Unit      string
}

func newSignalDecoding(sig Signal, rawValue uint64, valueType SignalValueType, value any, unit string) *SignalDecoding {
	return &SignalDecoding{
		Signal:    sig,
		RawValue:  rawValue,
		ValueType: valueType,
		Value:     value,
		Unit:      unit,
	}
}

// ValueAsFlag returns the decoded value as a flag.
// Returns false if the value type is not a flag.
func (sd *SignalDecoding) ValueAsFlag() bool {
	if sd.ValueType != SignalValueTypeFlag {
		return false
	}
	return sd.Value.(bool)
}

// ValueAsInt returns the decoded value as an integer.
// Returns 0 if the value type is not an integer.
func (sd *SignalDecoding) ValueAsInt() int64 {
	if sd.ValueType != SignalValueTypeInt {
		return 0
	}
	return sd.Value.(int64)
}

// ValueAsUint returns the decoded value as an unsigned integer.
// Returns 0 if the value type is not an unsigned integer.
func (sd *SignalDecoding) ValueAsUint() uint64 {
	if sd.ValueType != SignalValueTypeUint {
		return 0
	}
	return sd.Value.(uint64)
}

// ValueAsFloat returns the decoded value as a float.
// Returns 0 if the value type is not a float.
func (sd *SignalDecoding) ValueAsFloat() float64 {
	if sd.ValueType != SignalValueTypeFloat {
		return 0
	}
	return sd.Value.(float64)
}

// ValueAsEnum returns the decoded value as an enum.
// Returns an empty string if the value type is not an enum.
func (sd *SignalDecoding) ValueAsEnum() string {
	if sd.ValueType != SignalValueTypeEnum {
		return ""
	}
	return sd.Value.(string)
}

func (sd *SignalDecoding) String() string {
	s := stringer.New()

	s.Write("signal_name: %s; value_type: %s\n", sd.Signal.Name(), sd.ValueType)
	s.Write("raw_value: %x; value: %v\n", sd.RawValue, sd.Value)

	if sd.Unit != "" {
		s.Write("unit: %s\n", sd.Unit)
	}

	return s.String()
}

package acmelib

import (
	"fmt"
	"strings"
	"time"
)

// SignalKind rappresents the kind of a [Signal].
// It can be standard, enum, or multiplexer
type SignalKind int

const (
	// SignalKindStandard defines a standard signal.
	SignalKindStandard SignalKind = iota
	// SignalKindEnum defines a enum signal.
	SignalKindEnum
	// SignalKindMultiplexer defines a multiplexer signal.
	SignalKindMultiplexer
)

func (sk SignalKind) String() string {
	switch sk {
	case SignalKindStandard:
		return "standard"
	case SignalKindEnum:
		return "enum"
	case SignalKindMultiplexer:
		return "multiplexer"
	default:
		return "unknown"
	}
}

// SignalSendType rappresents the send type of a [Signal].
type SignalSendType int

const (
	// SignalSendTypeUnset defines an unset transmission type.
	SignalSendTypeUnset SignalSendType = iota
	// SignalSendTypeCyclic defines a cyclic transmission type.
	SignalSendTypeCyclic
	// SignalSendTypeOnWrite defines an on write transmission type.
	SignalSendTypeOnWrite
	// SignalSendTypeOnWriteWithRepetition defines an on write with repetition transmission type.
	SignalSendTypeOnWriteWithRepetition
	// SignalSendTypeOnChange defines an on change transmission type.
	SignalSendTypeOnChange
	// SignalSendTypeOnChangeWithRepetition defines an on change with repetition transmission type.
	SignalSendTypeOnChangeWithRepetition
	// SignalSendTypeIfActive defines an if active transmission type.
	SignalSendTypeIfActive
	// SignalSendTypeIfActiveWithRepetition defines an if active with repetition transmission type.
	SignalSendTypeIfActiveWithRepetition
)

func (sst SignalSendType) String() string {
	switch sst {
	case SignalSendTypeUnset:
		return "unset"
	case SignalSendTypeCyclic:
		return "cyclic"
	case SignalSendTypeOnWrite:
		return "on-write"
	case SignalSendTypeOnWriteWithRepetition:
		return "on-write_with_repetition"
	case SignalSendTypeOnChange:
		return "on-change"
	case SignalSendTypeOnChangeWithRepetition:
		return "on-change_with_repetition"
	case SignalSendTypeIfActive:
		return "if_active"
	case SignalSendTypeIfActiveWithRepetition:
		return "if_active_with_repetition"
	default:
		return "unknown"
	}
}

// Signal interface specifies all common methods of
// [StandardSignal], [EnumSignal], and [MultiplexerSignal1].
type Signal interface {
	// EntityID returns the entity id of the signal.
	EntityID() EntityID
	// Name returns the name of the signal.
	Name() string
	// SetDesc stes the description of the signal.
	SetDesc(desc string)
	// Desc returns the description of the signal.
	Desc() string
	// CreateTime returns the creation time of the signal.
	CreateTime() time.Time

	// AddAttributeValue adds to the signal an attribute and its value.
	AddAttributeValue(attribute Attribute, value any) error
	// RemoveAttributeValue removes from the signal an attribute and its value.
	RemoveAttributeValue(attributeEntityID EntityID) error
	// RemoveAllAttributeValues removes all attribute and their values from the signal.
	RemoveAllAttributeValues()
	// AttributeValues returns all attribute and their values from the signal.
	AttributeValues() []*AttributeValue
	// GetAttributeValue returns the value of an attribute and its value from the signal.
	GetAttributeValue(attributeEntityID EntityID) (*AttributeValue, error)

	stringify(b *strings.Builder, tabs int)
	String() string

	// Kind returns the kind of the signal.
	Kind() SignalKind

	// ParentMessage returns the parent message of the signal or nil if not set.
	ParentMessage() *Message
	// ParentMultiplexerSignal returns the parent multiplexer signal of the signal
	// or nil if not set.
	ParentMultiplexerSignal() *MultiplexerSignal

	setParentMsg(parentMsg *Message)
	setParentMuxSig(parentMuxSig *MultiplexerSignal)

	// SetStartValue sets the initial raw value of the signal.
	SetStartValue(startValue int)
	// StartValue returns the initial raw value of the signal.
	StartValue() int
	// SetSendType sets the send type of the signal.
	SetSendType(sendType SignalSendType)
	// SendType returns the send type of the signal.
	SendType() SignalSendType

	// GetStartBit returns the start bit of the signal.
	GetStartBit() int
	getRelStartBit() int
	setRelStartBit(startBit int)

	// GetSize returns the size of the signal.
	GetSize() int

	// ToStandard returns the signal as a standard signal.
	ToStandard() (*StandardSignal, error)
	// ToEnum returns the signal as a enum signal.
	ToEnum() (*EnumSignal, error)
	// ToMultiplexer returns the signal as a multiplexer signal.
	ToMultiplexer() (*MultiplexerSignal, error)
}

type signal struct {
	*attributeEntity

	parentMsg    *Message
	parentMuxSig *MultiplexerSignal

	kind SignalKind

	startValue int
	sendType   SignalSendType

	relStartBit int
}

func newSignal(name string, kind SignalKind) *signal {
	return &signal{
		attributeEntity: newAttributeEntity(name, AttributeRefKindSignal),

		parentMsg:    nil,
		parentMuxSig: nil,

		kind: kind,

		startValue: 0,
		sendType:   SignalSendTypeUnset,

		relStartBit: 0,
	}
}

func (s *signal) hasParentMsg() bool {
	return s.parentMsg != nil
}

func (s *signal) hasParentMuxSig() bool {
	return s.parentMuxSig != nil
}

func (s *signal) modifySize(amount int) error {
	if s.hasParentMuxSig() {
		return s.parentMuxSig.modifySignalSize(s.EntityID(), amount)
	}

	if s.hasParentMsg() {
		return s.parentMsg.modifySignalSize(s.EntityID(), amount)
	}

	return nil
}

func (s *signal) errorf(err error) error {
	sigErr := &EntityError{
		Kind:     EntityKindSignal,
		EntityID: s.entityID,
		Name:     s.name,
		Err:      err,
	}

	if s.hasParentMsg() {
		return s.parentMsg.errorf(sigErr)
	}

	return sigErr
}

func (s *signal) Kind() SignalKind {
	return s.kind
}

func (s *signal) ParentMessage() *Message {
	return s.parentMsg
}

func (s *signal) ParentMultiplexerSignal() *MultiplexerSignal {
	return s.parentMuxSig
}

func (s *signal) setParentMsg(parentMsg *Message) {
	s.parentMsg = parentMsg
}

func (s *signal) setParentMuxSig(parentMuxSig *MultiplexerSignal) {
	s.parentMuxSig = parentMuxSig
}

func (s *signal) getRelStartBit() int {
	return s.relStartBit
}

func (s *signal) setRelStartBit(startBit int) {
	s.relStartBit = startBit
}

func (s *signal) stringify(b *strings.Builder, tabs int) {
	s.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%skind: %s\n", tabStr, s.kind))

	if s.sendType != SignalSendTypeUnset {
		b.WriteString(fmt.Sprintf("%ssend_type: %q\n", tabStr, s.sendType))
	}

	b.WriteString(fmt.Sprintf("%sstart_bit: %d; ", tabStr, s.relStartBit))
}

func (s *signal) SetStartValue(startValue int) {
	s.startValue = startValue
}

func (s *signal) StartValue() int {
	return s.startValue
}

func (s *signal) SetSendType(sendType SignalSendType) {
	s.sendType = sendType
}

func (s *signal) SendType() SignalSendType {
	return s.sendType
}

func (s *signal) GetStartBit() int {
	if s.hasParentMuxSig() {
		return s.parentMuxSig.GetStartBit() + s.parentMuxSig.GetGroupCountSize() + s.relStartBit
	}
	return s.relStartBit
}

func (s *signal) UpdateName(newName string) error {
	sigID := s.entityID
	oldName := s.name

	if oldName == newName {
		return nil
	}

	canUpdMuxSig := false
	if s.hasParentMuxSig() {
		if err := s.parentMuxSig.verifySignalName(sigID, newName); err != nil {
			return s.errorf(&UpdateNameError{Err: err})
		}
		canUpdMuxSig = true
	}

	if s.hasParentMsg() {
		if err := s.parentMsg.verifySignalName(newName); err != nil {
			return s.errorf(&UpdateNameError{Err: err})
		}

		s.parentMsg.signalNames.remove(oldName)
		s.parentMsg.signalNames.add(newName, sigID)
	}

	if canUpdMuxSig {
		s.parentMuxSig.signalNames.remove(oldName)
		s.parentMuxSig.signalNames.add(newName, sigID)
	}

	s.name = newName

	return nil
}

// StandardSignal is the representation of a normal signal that has a [SignalType],
// a min, a max, an offset, a scale, and can have a [SignalUnit].
type StandardSignal struct {
	*signal

	typ    *SignalType
	min    float64
	max    float64
	offset float64
	scale  float64
	unit   *SignalUnit
}

// NewStandardSignal creates a new [StandardSignal] with the given name and [SignalType].
// It may return an error if the given [SignalType] is nil.
func NewStandardSignal(name string, typ *SignalType) (*StandardSignal, error) {
	if typ == nil {
		return nil, &ArgumentError{
			Name: "typ",
			Err:  ErrIsNil,
		}
	}

	return &StandardSignal{
		signal: newSignal(name, SignalKindStandard),

		typ:    typ,
		min:    typ.min,
		max:    typ.max,
		offset: 0,
		scale:  1,
		unit:   nil,
	}, nil
}

// GetSize returns the size of the [StandardSignal].
func (ss *StandardSignal) GetSize() int {
	return ss.typ.size
}

// ToStandard returns the [StandardSignal] itself.
func (ss *StandardSignal) ToStandard() (*StandardSignal, error) {
	return ss, nil
}

// ToEnum always returns an error, because a [StandardSignal] cannot be converted to an [EnumSignal].
func (ss *StandardSignal) ToEnum() (*EnumSignal, error) {
	return nil, ss.errorf(&ConversionError{
		From: SignalKindStandard.String(),
		To:   SignalKindEnum.String(),
	})
}

// ToMultiplexer always returns an error, because a [StandardSigna] cannot be converted to a [MultiplexerSignal].
func (ss *StandardSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return nil, ss.errorf(&ConversionError{
		From: SignalKindStandard.String(),
		To:   SignalKindMultiplexer.String(),
	})
}

func (ss *StandardSignal) stringify(b *strings.Builder, tabs int) {
	ss.signal.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("size: %d\n", ss.GetSize()))
	b.WriteString(fmt.Sprintf("%smin: %g; max: %g; offset: %g; scale: %g\n", tabStr, ss.min, ss.max, ss.offset, ss.scale))

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

	ss.typ = typ
	ss.min = typ.min
	ss.max = typ.max
	ss.offset = 0
	ss.scale = 1

	return nil
}

// SetPhysicalValues sets the physical min, max, offset, and scale of the [StandardSignal].
// It returns an error if the scale is equal to 0.
func (ss *StandardSignal) SetPhysicalValues(min, max, offset, scale float64) error {
	if scale == 0 {
		return ss.errorf(&ArgumentError{
			Name: "scale",
			Err:  ErrIsZero,
		})
	}

	ss.min = min
	ss.max = max
	ss.offset = offset
	ss.scale = scale

	return nil
}

// Min returns the minimum value of the [StandardSignal].
// It may differ from the minimum value of the signal type associated
// with the [StandardSignal].
func (ss *StandardSignal) Min() float64 {
	return ss.min
}

// Max returns the maximum value of the [StandardSignal].
// It may differ from the maximum value of the signal type associated
// with the [StandardSignal].
func (ss *StandardSignal) Max() float64 {
	return ss.max
}

// Offset returns the offset of the [StandardSignal].
func (ss *StandardSignal) Offset() float64 {
	return ss.offset
}

// Scale returns the scale of the [StandardSignal].
func (ss *StandardSignal) Scale() float64 {
	return ss.scale
}

// SetUnit sets the [SignalUnit] of the [StandardSignal] to the given one.
func (ss *StandardSignal) SetUnit(unit *SignalUnit) {
	ss.unit = unit
}

// Unit returns the [SignalUnit] of the [StandardSignal].
func (ss *StandardSignal) Unit() *SignalUnit {
	return ss.unit
}

// EnumSignal is a signal that holds a [SignalEnum].
type EnumSignal struct {
	*signal

	enum *SignalEnum
}

// NewEnumSignal creates a new [EnumSignal] with the given name and [SignalEnum].
// It may return an error if the given [SignalEnum] is nil.
func NewEnumSignal(name string, enum *SignalEnum) (*EnumSignal, error) {
	if enum == nil {
		return nil, &ArgumentError{
			Name: "enum",
			Err:  ErrIsNil,
		}
	}

	sig := &EnumSignal{
		signal: newSignal(name, SignalKindEnum),

		enum: enum,
	}

	enum.parentSignals.add(sig.entityID, sig)

	return sig, nil
}

// GetSize returns the size of the [EnumSignal].
func (es *EnumSignal) GetSize() int {
	return es.enum.GetSize()
}

// ToStandard always returns an error, because an [EnumSignal] cannot be converted to a [StandardSignal].
func (es *EnumSignal) ToStandard() (*StandardSignal, error) {
	return nil, es.errorf(&ConversionError{
		From: SignalKindEnum.String(),
		To:   SignalKindStandard.String(),
	})
}

// ToEnum returns the [EnumSignal] itself.
func (es *EnumSignal) ToEnum() (*EnumSignal, error) {
	return es, nil
}

// ToMultiplexer always returns an error, because an [EnumSignal] cannot be converted to a [MultiplexerSignal].
func (es *EnumSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return nil, es.errorf(&ConversionError{
		From: SignalKindEnum.String(),
		To:   SignalKindMultiplexer.String(),
	})
}

func (es *EnumSignal) stringify(b *strings.Builder, tabs int) {
	es.signal.stringify(b, tabs)
	b.WriteString(fmt.Sprintf("size: %d\n", es.GetSize()))

	tabStr := getTabString(tabs)
	b.WriteString(fmt.Sprintf("%senum:\n", tabStr))

	es.enum.stringify(b, tabs+1)
}

func (es *EnumSignal) String() string {
	builder := new(strings.Builder)
	es.stringify(builder, 0)
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
		return es.errorf(&ArgumentError{
			Name: "enum",
			Err:  ErrIsNil,
		})
	}

	if err := es.modifySize(enum.GetSize() - es.GetSize()); err != nil {
		return es.errorf(err)
	}

	es.enum = enum

	es.enum.parentSignals.remove(es.entityID)

	enum.parentSignals.add(es.entityID, es)

	return nil
}

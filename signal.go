package acmelib

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type SignalKind string

const (
	SignalKindStandard    SignalKind = "standard"
	SignalKindEnum        SignalKind = "enum"
	SignalKindMultiplexer SignalKind = "multiplexer"
)

type Signal interface {
	GetEntityID() EntityID
	GetName() string
	GetDesc() string
	GetCreateTime() time.Time
	GetUpdateTime() time.Time

	String() string

	GetKind() SignalKind

	GetParentMessage() *Message
	setParentMessage(parentMessage *Message)

	GetStartBit() int
	setStartBit(startBit int)

	GetSize() int

	ToStandard() (*StandardSignal, error)
	ToEnum() (*EnumSignal, error)
	ToMultiplexer() (*MultiplexerSignal, error)
}

type signal struct {
	*entity

	kind          SignalKind
	parentMessage *Message
	startBit      int
}

func newSignal(name, desc string, kind SignalKind) *signal {
	return &signal{
		entity: newEntity(name, desc),

		kind:          kind,
		parentMessage: nil,
		startBit:      0,
	}
}

func (s *signal) String() string {
	var builder strings.Builder

	builder.WriteString("\n+++START SIGNAL+++\n\n")
	builder.WriteString(s.toString())
	builder.WriteString(fmt.Sprintf("kind: %s\n", s.kind))
	builder.WriteString(fmt.Sprintf("start_bit: %d; ", s.startBit))

	return builder.String()
}

func (s *signal) GetKind() SignalKind {
	return s.kind
}

func (s *signal) GetParentMessage() *Message {
	return s.parentMessage
}

func (s *signal) setParentMessage(parentMessage *Message) {
	s.parentMessage = parentMessage
}

func (s *signal) GetStartBit() int {
	return s.startBit
}

func (s *signal) setStartBit(startBit int) {
	s.startBit = startBit
}

func (s *signal) hasParent() bool {
	return s.parentMessage != nil
}

func (s *signal) errorf(err error) error {
	sigErr := fmt.Errorf(`signal "%s": %v`, s.Name, err)
	if s.hasParent() {
		return s.parentMessage.errorf(sigErr)
	}
	return sigErr
}

func (s *signal) UpdateName(newName string) error {
	if s.hasParent() {
		if err := s.parentMessage.signals.updateEntityName(s.EntityID, s.Name, newName); err != nil {
			return err
		}
	}
	return s.entity.UpdateName(newName)
}

func (s *signal) modifySize(amount int) error {
	if s.hasParent() {
		return s.parentMessage.modifySignalSize(s.GetEntityID(), amount)
	}
	return nil
}

type StandardSignal struct {
	*signal

	typ    *SignalType
	min    float64
	max    float64
	offset float64
	scale  float64
	unit   *SIgnalUnit
}

func NewStandardSignal(name, desc string, typ *SignalType, min, max, offset, scale float64, unit *SIgnalUnit) (*StandardSignal, error) {
	if typ == nil {
		return nil, errors.New("signal type cannot be nil")
	}

	return &StandardSignal{
		signal: newSignal(name, desc, SignalKindStandard),

		typ:    typ,
		min:    min,
		max:    max,
		offset: offset,
		scale:  scale,
		unit:   unit,
	}, nil
}

func (ss *StandardSignal) String() string {
	var builder strings.Builder

	builder.WriteString(ss.signal.String())
	builder.WriteString(fmt.Sprintf("size: %d\n", ss.GetSize()))
	builder.WriteString(fmt.Sprintf("min: %f; max: %f; offset: %f; scale: %f\n", ss.min, ss.max, ss.offset, ss.scale))

	builder.WriteString("\n+++END SIGNAL+++\n\n")

	return builder.String()
}

func (ss *StandardSignal) GetSize() int {
	return ss.typ.Size
}

func (ss *StandardSignal) ToStandard() (*StandardSignal, error) {
	return ss, nil
}

func (ss *StandardSignal) ToEnum() (*EnumSignal, error) {
	return nil, ss.errorf(errors.New(`cannot covert to "enum", the signal is of kind "standard"`))
}

func (ss *StandardSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return nil, ss.errorf(errors.New(`cannot covert to "multiplexer", the signal is of kind "standard"`))
}

func (ss *StandardSignal) GetType() *SignalType {
	return ss.typ
}

func (ss *StandardSignal) GetMin() float64 {
	return ss.min
}

func (ss *StandardSignal) GetMax() float64 {
	return ss.max
}

func (ss *StandardSignal) GetOffset() float64 {
	return ss.offset
}

func (ss *StandardSignal) GetScale() float64 {
	return ss.scale
}

func (ss *StandardSignal) GetUnit() *SIgnalUnit {
	return ss.unit
}

func (ss *StandardSignal) UpdateType(newType *SignalType) error {
	if err := ss.modifySize(newType.Size - ss.typ.Size); err != nil {
		return ss.errorf(err)
	}

	ss.typ = newType
	ss.setUpdateTimeNow()

	return nil
}

package acmelib

import (
	"fmt"
	"time"
)

type SignalKind int

const (
	SignalKindStandard SignalKind = iota
	SignalKindEnum
)

type Signal interface {
	GetEntityID() EntityID
	GetName() string
	GetDesc() string
	GetCreateTime() time.Time
	GetUpdateTime() time.Time

	GetKind() SignalKind

	GetParentMessage() *Message
	setParentMessage(parentMessage *Message)

	GetStartBit() int
	setStartBit(startBit int)

	GetSize() int

	ToStandard() (*standardSignal, error)
	ToEnum() (*EnumSignal, error)
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

func (s *signal) errorf(err error) error {
	sigErr := fmt.Errorf(`signal "%s": %v`, s.Name, err)
	if s.parentMessage != nil {
		return s.parentMessage.errorf(sigErr)
	}
	return sigErr
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

type standardSignal struct {
	*entity
	ParentMessage *Message

	Kind     SignalKind
	Type     *SignalType
	StartBit int
	Index    int
	Min      float64
	Max      float64
	Offset   float64
	Scale    float64
	Unit     *SIgnalUnit
}

func NewStandardSignal(name, desc string, typ *SignalType, min, max, offset, scale float64, unit *SIgnalUnit) *standardSignal {
	return &standardSignal{
		entity: newEntity(name, desc),

		Kind:   SignalKindStandard,
		Type:   typ,
		Min:    min,
		Max:    max,
		Offset: offset,
		Scale:  scale,
		Unit:   unit,
	}
}

func (s *standardSignal) errorf(err error) error {
	sigErr := fmt.Errorf(`signal "%s": %v`, s.Name, err)
	if s.ParentMessage != nil {
		return s.ParentMessage.errorf(sigErr)
	}
	return sigErr
}

func (s *standardSignal) BitSize() int {
	return s.Type.Size
}

func (s *standardSignal) UpdateName(name string) error {
	if s.ParentMessage != nil {
		if err := s.ParentMessage.signals.updateEntityName(s.EntityID, s.Name, name); err != nil {
			return s.errorf(err)
		}
	}

	return s.entity.UpdateName(name)
}

func (s *standardSignal) UpdatePosition(pos int) error {
	if s.Index == pos {
		return nil
	}

	signals := s.ParentMessage.SignalsByStartBit()
	sigCount := len(signals)

	if pos >= sigCount {
		return s.errorf(fmt.Errorf(`position "%d" is out of bounds`, pos))
	}

	for idx, sig := range signals {
		if sig.EntityID == s.EntityID {
			for i := idx + 1; i < sigCount; i++ {
				tmpSig := signals[i]
				if tmpSig.Index <= pos {
					tmpSig.Index--
				}
			}

			sig.Index = pos
			s.setUpdateTimeNow()

			break
		}
	}

	return nil
}

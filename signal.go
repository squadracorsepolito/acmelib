package acmelib

import (
	"fmt"
)

type SignalKind string

const (
	SignalKindStandard    SignalKind = "signal_standard"
	SignalKindEnum        SignalKind = "signal_enum"
	SignalKindMultiplexed SignalKind = "signal_multiplexed"
)

type SignalPosition struct {
	From int
	To   int

	size int
}

func NewSignalPosition(from, to int) *SignalPosition {
	return &SignalPosition{
		From: from,
		To:   to,

		size: to - from,
	}
}

type Signal struct {
	*entity
	ParentMessage *Message

	Kind     SignalKind
	Type     *SignalType
	StartBit int
	Position int
	Min      float64
	Max      float64
	Offset   float64
	Scale    float64
	Unit     *SIgnalUnit
}

func NewStandardSignal(name, desc string, typ *SignalType, min, max, offset, scale float64, unit *SIgnalUnit) *Signal {
	return &Signal{
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

func (s *Signal) errorf(err error) error {
	sigErr := fmt.Errorf(`signal "%s": %v`, s.Name, err)
	if s.ParentMessage != nil {
		return s.ParentMessage.errorf(sigErr)
	}
	return sigErr
}

func (s *Signal) UpdateName(name string) error {
	if err := s.ParentMessage.signals.updateEntityName(s.EntityID, s.Name, name); err != nil {
		return s.errorf(err)
	}

	return s.entity.UpdateName(name)
}

func (s *Signal) UpdatePosition(pos int) error {
	if s.Position == pos {
		return nil
	}

	signals := s.ParentMessage.SignalsByPosition()
	sigCount := len(signals)

	if pos >= sigCount {
		return s.errorf(fmt.Errorf(`position "%d" is out of bounds`, pos))
	}

	for idx, sig := range signals {
		if sig.EntityID == s.EntityID {
			for i := idx + 1; i < sigCount; i++ {
				tmpSig := signals[i]
				if tmpSig.Position <= pos {
					tmpSig.Position--
				}
			}

			sig.Position = pos
			s.setUpdateTimeNow()

			break
		}
	}

	return nil
}

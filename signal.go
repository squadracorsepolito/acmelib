package acmelib

import (
	"fmt"
	"slices"
)

type SignalSortingMethod string

const (
	SignalsByName       SignalSortingMethod = "signals_by_name"
	SignalsByCreateTime SignalSortingMethod = "signals_by_create_time"
	SignalsByUpdateTime SignalSortingMethod = "signals_by_update_time"
	SignalsByPosition   SignalSortingMethod = "signals_by_position"
)

func sortSignalsByPosition(signals []*Signal) []*Signal {
	slices.SortFunc(signals, func(a, b *Signal) int { return a.Position - b.Position })
	return signals
}

var signalsSorter = newEntitySorter(
	newEntitySorterMethod(SignalsByName, func(signals []*Signal) []*Signal { return sortByName(signals) }),
	newEntitySorterMethod(SignalsByCreateTime, func(signals []*Signal) []*Signal { return sortByCreateTime(signals) }),
	newEntitySorterMethod(SignalsByUpdateTime, func(signals []*Signal) []*Signal { return sortByUpdateTime(signals) }),
	newEntitySorterMethod(SignalsByPosition, func(signals []*Signal) []*Signal { return sortSignalsByPosition(signals) }),
)

func SelectSignalsSortingMethod(method SignalSortingMethod) {
	signalsSorter.selectSortingMethod(method)
}

type SignalKind string

const (
	SignalKindStandard    SignalKind = "signal_standard"
	SignalKindEnum        SignalKind = "signal_enum"
	SignalKindMultiplexed SignalKind = "signal_multiplexed"
)

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
	return s.ParentMessage.errorf(fmt.Errorf("signal %s: %v", s.Name, err))
}

func (s *Signal) UpdateName(name string) error {
	if err := s.ParentMessage.signals.updateName(s.ID, s.Name, name); err != nil {
		return s.errorf(err)
	}

	return s.entity.UpdateName(name)
}

func (s *Signal) UpdatePosition(pos int) error {
	if s.Position == pos {
		return nil
	}

	signals := sortSignalsByPosition(s.ParentMessage.signals.listEntities())
	sigCount := len(signals)

	if pos >= sigCount {
		return s.errorf(fmt.Errorf(`position "%d" is out of bounds`, pos))
	}

	for idx, sig := range signals {
		if sig.ID == s.ID {
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

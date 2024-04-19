package acmelib

import (
	"errors"
	"slices"
)

type MultiplexerSignal0 struct {
	*signal

	muxSignals     *set[EntityID, Signal]
	muxSigNames    *set[string, EntityID]
	muxSigGroupIDs *set[EntityID, []int]

	fixedSignals *set[EntityID, bool]

	groupCount int
	groupSize  int
	groups     []*signalPayload
}

func NewMultiplexerSignal0(name string, groupCount, groupSize int) *MultiplexerSignal0 {
	groups := make([]*signalPayload, groupCount)
	for i := 0; i < groupCount; i++ {
		groups[i] = newSignalPayload(groupSize)
	}

	return &MultiplexerSignal0{
		signal: newSignal(name, SignalKindMultiplexer),

		muxSignals:     newSet[EntityID, Signal]("multiplexed signal"),
		muxSigNames:    newSet[string, EntityID]("multiplexed signal name"),
		muxSigGroupIDs: newSet[EntityID, []int]("multiplexer signal group id"),

		fixedSignals: newSet[EntityID, bool]("fixed signals"),

		groupCount: groupCount,
		groupSize:  groupSize,
		groups:     groups,
	}
}

func (ms *MultiplexerSignal0) addMuxSignal(sig Signal) {
	id := sig.EntityID()
	name := sig.Name()

	ms.muxSignals.add(id, sig)
	ms.muxSigNames.add(name, id)

	if ms.hasParent() {
		parent := ms.Parent()
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case SignalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case SignalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				msgParent.signals.add(id, sig)
				msgParent.signalNames.add(name, id)
				return
			}
		}
	}
}

func (ms *MultiplexerSignal0) GetSize() int {
	return ms.groupSize + calcSizeFromValue(ms.groupCount)
}

func (ms *MultiplexerSignal0) InsertMuxSignal(signal Signal, startBit int, groupIDs ...int) error {
	// TODO! checks on groupIDs and signal name

	groupIDs = slices.Compact(groupIDs)

	prevGroupIDs := []int{}
	if ms.muxSigGroupIDs.hasKey(signal.EntityID()) {
		tmpIDs, err := ms.muxSigGroupIDs.getValue(signal.EntityID())
		if err != nil {
			panic(err)
		}
		prevGroupIDs = tmpIDs
	}

	for _, groupID := range groupIDs {
		if slices.Contains(prevGroupIDs, groupID) {
			return ms.errorf(errors.New("dupl group id"))
		}

		if err := ms.groups[groupID].insert(signal, startBit); err != nil {
			return err
		}
	}

	ms.addMuxSignal(signal)

	signal.setParent(nil)

	return nil
}

func (ms *MultiplexerSignal0) InsertFixedSignal(signal Signal, startBit int) error {
	// TODO! check signal name

	for i := 0; i < ms.groupCount; i++ {
		if err := ms.groups[i].insert(signal, startBit); err != nil {
			return err
		}
	}

	ms.addMuxSignal(signal)
	ms.fixedSignals.add(signal.EntityID(), true)

	signal.setParent(nil)

	return nil
}

func (ms *MultiplexerSignal0) RemoveSignal(muxSignalEntityID EntityID) error {
	return nil
}

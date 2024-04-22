package acmelib

import (
	"fmt"
	"slices"
	"strings"
)

type MultiplexerSignal struct {
	*signal

	signals        *set[EntityID, Signal]
	signalNames    *set[string, EntityID]
	signalGroupIDs *set[EntityID, []int]

	fixedSignals *set[EntityID, bool]

	groupCount int
	groupSize  int
	groups     []*signalPayload
}

func NewMultiplexerSignal(name string, groupCount, groupSize int) (*MultiplexerSignal, error) {
	if groupCount <= 0 {
		err := &ArgumentError{
			Name: "groupCount",
		}

		if groupCount == 0 {
			err.Err = ErrIsZero
			return nil, err
		}

		err.Err = ErrIsNegative
		return nil, err
	}

	if groupSize <= 0 {
		err := &ArgumentError{
			Name: "groupSize",
		}

		if groupSize == 0 {
			err.Err = ErrIsZero
			return nil, err
		}

		err.Err = ErrIsNegative
		return nil, err
	}

	groups := make([]*signalPayload, groupCount)
	for i := 0; i < groupCount; i++ {
		groups[i] = newSignalPayload(groupSize)
	}

	return &MultiplexerSignal{
		signal: newSignal(name, SignalKindMultiplexer),

		signals:        newSet[EntityID, Signal](),
		signalNames:    newSet[string, EntityID](),
		signalGroupIDs: newSet[EntityID, []int](),

		fixedSignals: newSet[EntityID, bool](),

		groupCount: groupCount,
		groupSize:  groupSize,
		groups:     groups,
	}, nil
}

func (ms *MultiplexerSignal) addSignal(sig Signal) {
	id := sig.EntityID()
	name := sig.Name()

	ms.signals.add(id, sig)
	ms.signalNames.add(name, id)

	// sig.setParent(ms)
	sig.setParentMuxSig(ms)

	if ms.hasParentMsg() {
		if sig.Kind() == SignalKindMultiplexer {
			muxSig, err := sig.ToMultiplexer()
			if err != nil {
				panic(err)
			}

			for tmpSigID, tmpSig := range muxSig.signals.entries() {
				ms.parentMsg.signals.add(tmpSigID, tmpSig)
			}

			for tmpName, tmpSigID := range muxSig.signalNames.entries() {
				ms.parentMsg.signalNames.add(tmpName, tmpSigID)
			}
		}

		ms.parentMsg.signals.add(id, sig)
		ms.parentMsg.signalNames.add(name, id)

		sig.setParentMsg(ms.parentMsg)
	}

	// if ms.hasParent() {
	// 	parent := ms.Parent()
	// 	for parent != nil {
	// 		switch parent.GetSignalParentKind() {
	// 		case SignalParentKindMultiplexerSignal:
	// 			muxParent, err := parent.ToParentMultiplexerSignal()
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			parent = muxParent.Parent()

	// 		case SignalParentKindMessage:
	// 			msgParent, err := parent.ToParentMessage()
	// 			if err != nil {
	// 				panic(err)
	// 			}

	// 			if sig.Kind() == SignalKindMultiplexer {
	// 				muxSig, err := sig.ToMultiplexer()
	// 				if err != nil {
	// 					panic(err)
	// 				}

	// 				for tmpSigID, tmpSig := range muxSig.signals.entries() {
	// 					msgParent.signals.add(tmpSigID, tmpSig)
	// 				}

	// 				for tmpName, tmpSigID := range muxSig.signalNames.entries() {
	// 					msgParent.signalNames.add(tmpName, tmpSigID)
	// 				}

	// 			}

	// 			msgParent.signals.add(id, sig)
	// 			msgParent.signalNames.add(name, id)

	// 			return
	// 		}
	// 	}
	// }
}

func (ms *MultiplexerSignal) removeSignal(sig Signal) {
	id := sig.EntityID()
	name := sig.Name()

	ms.signals.remove(id)
	ms.signalNames.remove(name)

	// sig.setParent(nil)
	sig.setParentMuxSig(nil)

	if ms.hasParentMsg() {
		if sig.Kind() == SignalKindMultiplexer {
			muxSig, err := sig.ToMultiplexer()
			if err != nil {
				panic(err)
			}

			for _, tmpSigID := range muxSig.signals.getKeys() {
				ms.parentMsg.signals.remove(tmpSigID)
			}

			for _, tmpName := range muxSig.signalNames.getKeys() {
				ms.parentMsg.signalNames.remove(tmpName)
			}
		}

		ms.parentMsg.signals.remove(id)
		ms.parentMsg.signalNames.remove(name)

		sig.setParentMsg(nil)
	}

	// if ms.hasParent() {
	// 	parent := ms.Parent()
	// 	for parent != nil {
	// 		switch parent.GetSignalParentKind() {
	// 		case SignalParentKindMultiplexerSignal:
	// 			muxParent, err := parent.ToParentMultiplexerSignal()
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			parent = muxParent.Parent()

	// 		case SignalParentKindMessage:
	// 			msgParent, err := parent.ToParentMessage()
	// 			if err != nil {
	// 				panic(err)
	// 			}

	// 			if sig.Kind() == SignalKindMultiplexer {
	// 				muxSig, err := sig.ToMultiplexer()
	// 				if err != nil {
	// 					panic(err)
	// 				}

	// 				for tmpSigID := range muxSig.signals.entries() {
	// 					msgParent.signals.remove(tmpSigID)
	// 				}

	// 				for tmpName := range muxSig.signalNames.entries() {
	// 					msgParent.signalNames.remove(tmpName)
	// 				}

	// 			}

	// 			msgParent.signals.remove(id)
	// 			msgParent.signalNames.remove(name)
	// 			return
	// 		}
	// 	}
	// }
}

// func (ms *MultiplexerSignal) modifySignalName(sigID EntityID, newName string) error {
// 	if ms.hasParent() {
// 		parent := ms.Parent()

// 	loop:
// 		for parent != nil {
// 			switch parent.GetSignalParentKind() {
// 			case SignalParentKindMultiplexerSignal:
// 				muxParent, err := parent.ToParentMultiplexerSignal()
// 				if err != nil {
// 					panic(err)
// 				}
// 				parent = muxParent.Parent()

// 			case SignalParentKindMessage:
// 				msgParent, err := parent.ToParentMessage()
// 				if err != nil {
// 					panic(err)
// 				}

// 				if err := msgParent.modifySignalName(sigID, newName); err != nil {
// 					return err
// 				}
// 				break loop
// 			}
// 		}
// 	}

// 	sig, err := ms.signals.getValue(sigID)
// 	if err != nil {
// 		return err
// 	}

// 	oldName := sig.Name()

// 	ms.signalNames.remove(oldName)
// 	ms.signalNames.add(newName, sigID)

// 	return nil
// }

func (ms *MultiplexerSignal) verifySignalName(sigID EntityID, name string) error {

	// if ms.hasParent() {
	// 	return ms.parent.verifySignalName(sigID, name)
	// }

	if ms.signalNames.hasKey(name) {
		tmpSigID, err := ms.signalNames.getValue(name)
		if err != nil {
			panic(err)
		}

		if sigID != tmpSigID {
			return &NameError{
				Name: name,
				Err:  ErrIsDuplicated,
			}
		}

		return nil
	}

	if ms.hasParentMsg() {
		if err := ms.parentMsg.verifySignalName(name); err != nil {
			return &NameError{
				Name: name,
				Err:  err,
			}
		}
	}

	return nil
}

func (ms *MultiplexerSignal) verifySignalSizeAmount(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	sig, err := ms.signals.getValue(sigID)
	if err != nil {
		return err
	}

	groupIDs := []int{}

	if ms.fixedSignals.hasKey(sigID) {
		for i := 0; i < ms.groupCount; i++ {
			groupIDs = append(groupIDs, i)
		}
	} else {
		groupIDs, err = ms.signalGroupIDs.getValue(sigID)
		if err != nil {
			panic(err)
		}
	}

	for _, groupID := range groupIDs {
		if amount > 0 {
			if err := ms.groups[groupID].verifyBeforeGrow(sig, amount); err != nil {
				return err
			}

			continue
		}

		if err := ms.groups[groupID].verifyBeforeShrink(sig, -amount); err != nil {
			return err
		}
	}

	return nil
}

func (ms *MultiplexerSignal) modifySignalSize(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	sig, err := ms.signals.getValue(sigID)
	if err != nil {
		return err
	}

	groupIDs := []int{}

	if ms.fixedSignals.hasKey(sigID) {
		for i := 0; i < ms.groupCount; i++ {
			groupIDs = append(groupIDs, i)
		}
	} else {
		groupIDs, err = ms.signalGroupIDs.getValue(sigID)
		if err != nil {
			panic(err)
		}
	}

	for _, groupID := range groupIDs {
		if amount > 0 {
			if err := ms.groups[groupID].modifyStartBitsOnGrow(sig, amount); err != nil {
				return err
			}

			continue
		}

		if err := ms.groups[groupID].modifyStartBitsOnShrink(sig, -amount); err != nil {
			return err
		}
	}

	return nil
}

func (ms *MultiplexerSignal) verifyGroupID(groupID int) error {
	err := &GroupIDError{
		GroupID: groupID,
	}

	if groupID < 0 {
		err.Err = ErrIsNegative
		return err
	}

	if groupID >= ms.groupCount {
		err.Err = ErrOutOfBounds
		return err
	}

	return nil
}

func (ms *MultiplexerSignal) stringify(b *strings.Builder, tabs int) {
	ms.signal.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("size: %d\n", ms.GetSize()))

	if ms.signals.size() == 0 {
		return
	}

	for groupID, group := range ms.GetSignalGroups() {
		b.WriteString(fmt.Sprintf("%sgroup id: %d\n", tabStr, groupID))

		b.WriteString(fmt.Sprintf("%smultiplexed signals:\n", tabStr))
		for _, muxSig := range group {
			muxSig.stringify(b, tabs+1)
			b.WriteRune('\n')
		}
	}
}

func (ms *MultiplexerSignal) String() string {
	builder := new(strings.Builder)
	ms.stringify(builder, 0)
	return builder.String()
}

// GetSignalParentKind always returns [SignalParentKindMultiplexerSignal].
func (ms *MultiplexerSignal) GetSignalParentKind() SignalParentKind {
	return SignalParentKindMultiplexerSignal
}

// ToParentMessage always returns an error, since [MultiplexerSignal1] cannot be converted to [Message].
func (ms *MultiplexerSignal) ToParentMessage() (*Message, error) {
	return nil, fmt.Errorf(`cannot convert to "%s" signal parent is of kind "%s"`,
		SignalParentKindMessage, SignalParentKindMultiplexerSignal)
}

// ToParentMultiplexerSignal returns the [MultiplexerSignal] itself.
func (ms *MultiplexerSignal) ToParentMultiplexerSignal() (*MultiplexerSignal, error) {
	return ms, nil
}

func (ms *MultiplexerSignal) GetSize() int {
	return ms.groupSize + ms.GetGroupCountSize()
}

// ToStandard always returns an error, since [MultiplexerSignal1] cannot be converted to [StandardSignal].
func (ms *MultiplexerSignal) ToStandard() (*StandardSignal, error) {
	return nil, ms.errorf(fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindStandard, SignalKindMultiplexer))
}

// ToEnum always returns an error, since [MultiplexerSignal1] cannot be converted to [EnumSignal].
func (ms *MultiplexerSignal) ToEnum() (*EnumSignal, error) {
	return nil, ms.errorf(fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindEnum, SignalKindMultiplexer))
}

// ToMultiplexer always returns the [MultiplexerSignal] itself.
func (ms *MultiplexerSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return ms, nil
}

func (ms *MultiplexerSignal) InsertSignal(signal Signal, startBit int, groupIDs ...int) error {
	insErr := &InsertSignalError{
		EntityID: signal.EntityID(),
		Name:     signal.Name(),
		StartBit: startBit,
	}

	if err := ms.verifySignalName(signal.EntityID(), signal.Name()); err != nil {
		insErr.Err = err
		return ms.errorf(insErr)
	}

	if len(groupIDs) == 0 {
		for i := 0; i < ms.groupCount; i++ {
			if err := ms.groups[i].verifyBeforeInsert(signal, startBit); err != nil {
				insErr.Err = err
				return ms.errorf(insErr)
			}
		}

		for i := 0; i < ms.groupCount; i++ {
			ms.groups[i].insert(signal, startBit)
		}

		ms.fixedSignals.add(signal.EntityID(), true)

	} else {
		groupIDs = slices.Compact(groupIDs)

		prevGroupIDs := []int{}
		if ms.signalGroupIDs.hasKey(signal.EntityID()) {
			tmpIDs, err := ms.signalGroupIDs.getValue(signal.EntityID())
			if err != nil {
				panic(err)
			}
			prevGroupIDs = tmpIDs
		}

		for _, groupID := range groupIDs {
			if err := ms.verifyGroupID(groupID); err != nil {
				insErr.Err = err
				return ms.errorf(insErr)
			}

			if slices.Contains(prevGroupIDs, groupID) {
				insErr.Err = &GroupIDError{
					GroupID: groupID,
					Err:     ErrIsDuplicated,
				}
				return ms.errorf(insErr)
			}

			if err := ms.groups[groupID].verifyBeforeInsert(signal, startBit); err != nil {
				insErr.Err = err
				return ms.errorf(insErr)
			}
		}

		for _, groupID := range groupIDs {
			ms.groups[groupID].insert(signal, startBit)
		}

		groupIDs = slices.Concat(prevGroupIDs, groupIDs)
		slices.Sort(groupIDs)

		ms.signalGroupIDs.add(signal.EntityID(), groupIDs)
	}

	ms.addSignal(signal)

	return nil
}

func (ms *MultiplexerSignal) RemoveSignal(signalEntityID EntityID) error {
	sig, err := ms.signals.getValue(signalEntityID)
	if err != nil {
		return ms.errorf(&RemoveSignalError{
			EntityID: signalEntityID,
			Err:      err,
		})
	}

	if ms.fixedSignals.hasKey(signalEntityID) {
		for i := 0; i < ms.groupCount; i++ {
			ms.groups[i].remove(signalEntityID)
		}

		ms.removeSignal(sig)
		ms.fixedSignals.remove(signalEntityID)

		return nil
	}

	groupIDs, err := ms.signalGroupIDs.getValue(signalEntityID)
	if err != nil {
		panic(err)
	}

	for _, groupID := range groupIDs {
		ms.groups[groupID].remove(signalEntityID)
	}

	ms.removeSignal(sig)
	ms.signalGroupIDs.remove(signalEntityID)

	return nil
}

func (ms *MultiplexerSignal) ClearSignalGroup(groupID int) error {
	if err := ms.verifyGroupID(groupID); err != nil {
		return ms.errorf(&ClearSignalGroupError{
			Err: err,
		})
	}

	group := ms.groups[groupID]
	signals := slices.Clone(group.signals)

	for _, sig := range signals {
		sigID := sig.EntityID()

		if ms.fixedSignals.hasKey(sigID) {
			continue
		}

		group.remove(sigID)

		groupIDs, err := ms.signalGroupIDs.getValue(sigID)
		if err != nil {
			panic(err)
		}

		if len(groupIDs) == 1 {
			ms.removeSignal(sig)
			ms.signalGroupIDs.remove(sigID)

			continue
		}

		groupIDs = slices.DeleteFunc(groupIDs, func(id int) bool { return id == groupID })
		ms.signalGroupIDs.add(sigID, groupIDs)
	}

	return nil
}

func (ms *MultiplexerSignal) ClearAllSignalGroups() {
	for _, sig := range ms.signals.getValues() {
		ms.removeSignal(sig)
	}

	for i := 0; i < ms.groupCount; i++ {
		ms.groups[i].removeAll()
	}

	ms.signalGroupIDs.clear()
	ms.fixedSignals.clear()
}

func (ms *MultiplexerSignal) ShiftSignalLeft(signalEntityID EntityID, amount int) int {
	groupIDs, err := ms.signalGroupIDs.getValue(signalEntityID)
	if err != nil || len(groupIDs) > 1 {
		return 0
	}

	return ms.groups[groupIDs[0]].shiftLeft(signalEntityID, amount)
}

func (ms *MultiplexerSignal) ShiftSignalRight(signalEntityID EntityID, amount int) int {
	groupIDs, err := ms.signalGroupIDs.getValue(signalEntityID)
	if err != nil || len(groupIDs) > 1 {
		return 0
	}

	return ms.groups[groupIDs[0]].shiftRight(signalEntityID, amount)
}

func (ms *MultiplexerSignal) GetSignalGroup(groupID int) []Signal {
	if err := ms.verifyGroupID(groupID); err != nil {
		return []Signal{}
	}
	return ms.groups[groupID].signals
}

func (ms *MultiplexerSignal) GetSignalGroups() [][]Signal {
	res := make([][]Signal, ms.groupCount)

	for i := 0; i < ms.groupCount; i++ {
		res[i] = ms.groups[i].signals
	}

	return res
}

func (ms *MultiplexerSignal) GroupCount() int {
	return ms.groupCount
}

func (ms *MultiplexerSignal) GetGroupCountSize() int {
	return calcSizeFromValue(ms.groupCount - 1)
}

func (ms *MultiplexerSignal) GroupSize() int {
	return ms.groupSize
}

package acmelib

// import (
// 	"fmt"
// 	"slices"
// 	"strings"
// )

// // MultiplexerSignal is a signal that holds groups of other signals
// // that are selected/multiplexed by the value of the group id.
// // It can multiplex all the kinds of signals ([StandardSignal], [EnumSignal],
// // [MultiplexerSignal]), so it is possible to create multiple levels of multiplexing.
// type MultiplexerSignal struct {
// 	*signal

// 	signals        *set[EntityID, Signal]
// 	signalNames    *set[string, EntityID]
// 	signalGroupIDs *set[EntityID, []int]

// 	fixedSignals *set[EntityID, bool]

// 	groupCount int
// 	groupSize  int
// 	groups     []*SignalLayout //[]*signalPayload
// }

// func newMultiplexerSignalFromBase(base *signal, groupCount, groupSize int) (*MultiplexerSignal, error) {
// 	if groupCount <= 0 {
// 		err := &ArgumentError{
// 			Name: "groupCount",
// 		}

// 		if groupCount == 0 {
// 			err.Err = ErrIsZero
// 			return nil, err
// 		}

// 		err.Err = ErrIsNegative
// 		return nil, err
// 	}

// 	if groupSize <= 0 {
// 		err := &ArgumentError{
// 			Name: "groupSize",
// 		}

// 		if groupSize == 0 {
// 			err.Err = ErrIsZero
// 			return nil, err
// 		}

// 		err.Err = ErrIsNegative
// 		return nil, err
// 	}

// 	groups := make([]*SignalLayout, groupCount)
// 	for i := 0; i < groupCount; i++ {
// 		groups[i] = newSignalLayout(groupSize)
// 	}

// 	// groups := make([]*signalPayload, groupCount)
// 	// for i := 0; i < groupCount; i++ {
// 	// 	groups[i] = newSignalPayload(groupSize)
// 	// }

// 	return &MultiplexerSignal{
// 		signal: base,

// 		signals:        newSet[EntityID, Signal](),
// 		signalNames:    newSet[string, EntityID](),
// 		signalGroupIDs: newSet[EntityID, []int](),

// 		fixedSignals: newSet[EntityID, bool](),

// 		groupCount: groupCount,
// 		groupSize:  groupSize,
// 		groups:     groups,
// 	}, nil
// }

// // NewMultiplexerSignal creates a new [MultiplexerSignal] with the given name,
// // group count and group size.
// // The group count defines the number of groups that the signal will hold
// // and the group size defines the dimension in bits of each group.
// //
// // It will return an [ArgumentError] if group count or group size is invalid.
// func NewMultiplexerSignal(name string, groupCount, groupSize int) (*MultiplexerSignal, error) {
// 	return newMultiplexerSignalFromBase(newSignal(name, SignalKindMultiplexer), groupCount, groupSize)
// }

// func (ms *MultiplexerSignal) addSignal(sig Signal) {
// 	id := sig.EntityID()
// 	name := sig.Name()

// 	ms.signals.add(id, sig)
// 	ms.signalNames.add(name, id)

// 	sig.setParentMuxSig(ms)

// 	if ms.hasParentMsg() {
// 		if sig.Kind() == SignalKindMultiplexer {
// 			muxSig, err := sig.ToMultiplexer()
// 			if err != nil {
// 				panic(err)
// 			}

// 			for tmpSigID, tmpSig := range muxSig.signals.entries() {
// 				ms.parentMsg.signals.add(tmpSigID, tmpSig)
// 			}

// 			for tmpName, tmpSigID := range muxSig.signalNames.entries() {
// 				ms.parentMsg.signalNames.add(tmpName, tmpSigID)
// 			}
// 		}

// 		ms.parentMsg.signals.add(id, sig)
// 		ms.parentMsg.signalNames.add(name, id)

// 		sig.setParentMsg(ms.parentMsg)
// 	}
// }

// func (ms *MultiplexerSignal) removeSignal(sig Signal) {
// 	id := sig.EntityID()
// 	name := sig.Name()

// 	ms.signals.remove(id)
// 	ms.signalNames.remove(name)

// 	sig.setParentMuxSig(nil)

// 	if ms.hasParentMsg() {
// 		if sig.Kind() == SignalKindMultiplexer {
// 			muxSig, err := sig.ToMultiplexer()
// 			if err != nil {
// 				panic(err)
// 			}

// 			for _, tmpSigID := range muxSig.signals.getKeys() {
// 				ms.parentMsg.signals.remove(tmpSigID)
// 			}

// 			for _, tmpName := range muxSig.signalNames.getKeys() {
// 				ms.parentMsg.signalNames.remove(tmpName)
// 			}
// 		}

// 		ms.parentMsg.signals.remove(id)
// 		ms.parentMsg.signalNames.remove(name)

// 		sig.setParentMsg(nil)
// 	}
// }

// func (ms *MultiplexerSignal) verifySignalName(sigID EntityID, name string) error {
// 	if ms.signalNames.hasKey(name) {
// 		tmpSigID, err := ms.signalNames.getValue(name)
// 		if err != nil {
// 			panic(err)
// 		}

// 		if sigID != tmpSigID {
// 			return &NameError{
// 				Name: name,
// 				Err:  ErrIsDuplicated,
// 			}
// 		}

// 		return nil
// 	}

// 	if ms.hasParentMsg() {
// 		if err := ms.parentMsg.verifySignalName(name); err != nil {
// 			return &NameError{
// 				Name: name,
// 				Err:  err,
// 			}
// 		}
// 	}

// 	return nil
// }

// func (ms *MultiplexerSignal) verifySignalSizeAmount(sigID EntityID, amount int) error {
// 	if amount == 0 {
// 		return nil
// 	}

// 	sig, err := ms.signals.getValue(sigID)
// 	if err != nil {
// 		return err
// 	}

// 	groupIDs := []int{}

// 	if ms.fixedSignals.hasKey(sigID) {
// 		for i := 0; i < ms.groupCount; i++ {
// 			groupIDs = append(groupIDs, i)
// 		}
// 	} else {
// 		groupIDs, err = ms.signalGroupIDs.getValue(sigID)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	for _, groupID := range groupIDs {
// 		if amount > 0 {
// 			if err := ms.groups[groupID].verifyBeforeGrow(sig, amount); err != nil {
// 				return &SignalSizeError{
// 					Size: sig.GetSize() + amount,
// 					Err:  err,
// 				}
// 			}

// 			continue
// 		}

// 		if err := ms.groups[groupID].verifyBeforeShrink(sig, -amount); err != nil {
// 			return &SignalSizeError{
// 				Size: sig.GetSize() + amount,
// 				Err:  err,
// 			}
// 		}
// 	}

// 	return nil
// }

// func (ms *MultiplexerSignal) modifySignalSize(sigID EntityID, amount int) error {
// 	if amount == 0 {
// 		return nil
// 	}

// 	sig, err := ms.signals.getValue(sigID)
// 	if err != nil {
// 		panic(err)
// 	}

// 	groupIDs := []int{}

// 	if ms.fixedSignals.hasKey(sigID) {
// 		for i := 0; i < ms.groupCount; i++ {
// 			groupIDs = append(groupIDs, i)
// 		}
// 	} else {
// 		groupIDs, err = ms.signalGroupIDs.getValue(sigID)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	for _, groupID := range groupIDs {
// 		if amount > 0 {
// 			if err := ms.groups[groupID].modifyStartBitsOnGrow(sig, amount); err != nil {
// 				return err
// 			}

// 			continue
// 		}

// 		if err := ms.groups[groupID].modifyStartBitsOnShrink(sig, -amount); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (ms *MultiplexerSignal) verifyGroupID(groupID int) error {
// 	err := &GroupIDError{
// 		GroupID: groupID,
// 	}

// 	if groupID < 0 {
// 		err.Err = ErrIsNegative
// 		return err
// 	}

// 	if groupID >= ms.groupCount {
// 		err.Err = ErrOutOfBounds
// 		return err
// 	}

// 	return nil
// }

// func (ms *MultiplexerSignal) stringify(b *strings.Builder, tabs int) {
// 	ms.signal.stringify(b, tabs)

// 	tabStr := getTabString(tabs)

// 	b.WriteString(fmt.Sprintf("size: %d\n", ms.GetSize()))

// 	if ms.signals.size() == 0 {
// 		return
// 	}

// 	for groupID, group := range ms.GetSignalGroups() {
// 		b.WriteString(fmt.Sprintf("%sgroup id: %d\n", tabStr, groupID))

// 		b.WriteString(fmt.Sprintf("%smultiplexed signals:\n", tabStr))
// 		for _, muxSig := range group {
// 			muxSig.stringify(b, tabs+1)
// 			b.WriteRune('\n')
// 		}
// 	}
// }

// func (ms *MultiplexerSignal) String() string {
// 	builder := new(strings.Builder)
// 	ms.stringify(builder, 0)
// 	return builder.String()
// }

// // GetSize returns the total size of the [MultiplexerSignal].
// // The returned value is the sum of the size of the groups and
// // the number of bits needed to select the right group.
// // e.g. with group count = 8 and group size = 16, the total size
// // will be 3 + 16 = 19 bits, since 8 groups can be selected by 3 bits.
// func (ms *MultiplexerSignal) GetSize() int {
// 	return ms.groupSize + ms.GetGroupCountSize()
// }

// // ToStandard always returns a [ConversionError], since [MultiplexerSignal] cannot be converted to [StandardSignal].
// func (ms *MultiplexerSignal) ToStandard() (*StandardSignal, error) {
// 	return nil, ms.errorf(&ConversionError{
// 		From: SignalKindMultiplexer.String(),
// 		To:   SignalKindStandard.String(),
// 	})
// }

// // ToEnum always returns a [ConversionError], since [MultiplexerSignal] cannot be converted to [EnumSignal].
// func (ms *MultiplexerSignal) ToEnum() (*EnumSignal, error) {
// 	return nil, ms.errorf(&ConversionError{
// 		From: SignalKindMultiplexer.String(),
// 		To:   SignalKindEnum.String(),
// 	})
// }

// // ToMultiplexer always returns the [MultiplexerSignal] itself.
// func (ms *MultiplexerSignal) ToMultiplexer() (*MultiplexerSignal, error) {
// 	return ms, nil
// }

// // InsertSignal inserts a [Signal] at the given start bit.
// // If no group IDs are given, the signal will be considered as fixed
// // and it will be inserted into all groups.
// // If group IDs are given, the signal will be inserted into the given groups.
// //
// // It will return an [InsertSignalError] if the signal cannot be inserted
// // at the given start bit into the given group. This error can wrap:
// //   - [NameError] in case of an invalid signal name
// //   - [StartPosError] in case of an invalid start bit
// //   - [GroupIDError] in case of an invalid group ID
// func (ms *MultiplexerSignal) InsertSignal(signal Signal, startBit int, groupIDs ...int) error {
// 	if signal == nil {
// 		return &ArgumentError{
// 			Name: "signal",
// 			Err:  ErrIsNil,
// 		}
// 	}

// 	insErr := &InsertSignalError{
// 		EntityID: signal.EntityID(),
// 		Name:     signal.Name(),
// 		StartBit: startBit,
// 	}

// 	if err := ms.verifySignalName(signal.EntityID(), signal.Name()); err != nil {
// 		insErr.Err = err
// 		return ms.errorf(insErr)
// 	}

// 	if len(groupIDs) == 0 {
// 		for i := 0; i < ms.groupCount; i++ {
// 			if err := ms.groups[i].verifyBeforeInsert(signal, startBit); err != nil {
// 				insErr.Err = err
// 				return ms.errorf(insErr)
// 			}
// 		}

// 		for i := 0; i < ms.groupCount; i++ {
// 			ms.groups[i].insert(signal, startBit)
// 		}

// 		ms.fixedSignals.add(signal.EntityID(), true)

// 	} else {
// 		groupIDs = slices.Compact(groupIDs)

// 		prevGroupIDs := []int{}
// 		if ms.signalGroupIDs.hasKey(signal.EntityID()) {
// 			tmpIDs, err := ms.signalGroupIDs.getValue(signal.EntityID())
// 			if err != nil {
// 				panic(err)
// 			}
// 			prevGroupIDs = tmpIDs
// 		}

// 		for _, groupID := range groupIDs {
// 			if err := ms.verifyGroupID(groupID); err != nil {
// 				insErr.Err = err
// 				return ms.errorf(insErr)
// 			}

// 			if slices.Contains(prevGroupIDs, groupID) {
// 				insErr.Err = &GroupIDError{
// 					GroupID: groupID,
// 					Err:     ErrIsDuplicated,
// 				}
// 				return ms.errorf(insErr)
// 			}

// 			if err := ms.groups[groupID].verifyBeforeInsert(signal, startBit); err != nil {
// 				insErr.Err = err
// 				return ms.errorf(insErr)
// 			}
// 		}

// 		for _, groupID := range groupIDs {
// 			ms.groups[groupID].insert(signal, startBit)
// 		}

// 		groupIDs = slices.Concat(prevGroupIDs, groupIDs)
// 		slices.Sort(groupIDs)

// 		ms.signalGroupIDs.add(signal.EntityID(), groupIDs)
// 	}

// 	ms.addSignal(signal)

// 	return nil
// }

// // RemoveSignal removes the [Signal] with the given entity ID.
// //
// // It will return an [RemoveEntityError] if the signal cannot be removed.
// func (ms *MultiplexerSignal) RemoveSignal(signalEntityID EntityID) error {
// 	sig, err := ms.signals.getValue(signalEntityID)
// 	if err != nil {
// 		return ms.errorf(&RemoveEntityError{
// 			EntityID: signalEntityID,
// 			Err:      err,
// 		})
// 	}

// 	if ms.fixedSignals.hasKey(signalEntityID) {
// 		for i := 0; i < ms.groupCount; i++ {
// 			ms.groups[i].remove(signalEntityID)
// 		}

// 		ms.removeSignal(sig)
// 		ms.fixedSignals.remove(signalEntityID)

// 		return nil
// 	}

// 	groupIDs, err := ms.signalGroupIDs.getValue(signalEntityID)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, groupID := range groupIDs {
// 		ms.groups[groupID].remove(signalEntityID)
// 	}

// 	ms.removeSignal(sig)
// 	ms.signalGroupIDs.remove(signalEntityID)

// 	return nil
// }

// // ClearSignalGroup removes all signals from a group with the given group ID.
// //
// // It will return a [GroupIDError] if the given group ID is invalid.
// func (ms *MultiplexerSignal) ClearSignalGroup(groupID int) error {
// 	if err := ms.verifyGroupID(groupID); err != nil {
// 		return ms.errorf(err)
// 	}

// 	group := ms.groups[groupID]
// 	signals := slices.Clone(group.signals)

// 	for _, sig := range signals {
// 		sigID := sig.EntityID()

// 		if ms.fixedSignals.hasKey(sigID) {
// 			continue
// 		}

// 		group.remove(sigID)

// 		groupIDs, err := ms.signalGroupIDs.getValue(sigID)
// 		if err != nil {
// 			panic(err)
// 		}

// 		if len(groupIDs) == 1 {
// 			ms.removeSignal(sig)
// 			ms.signalGroupIDs.remove(sigID)

// 			continue
// 		}

// 		groupIDs = slices.DeleteFunc(groupIDs, func(id int) bool { return id == groupID })
// 		ms.signalGroupIDs.add(sigID, groupIDs)
// 	}

// 	return nil
// }

// // ClearAllSignalGroups removes all signals from all groups.
// func (ms *MultiplexerSignal) ClearAllSignalGroups() {
// 	for _, sig := range ms.signals.getValues() {
// 		ms.removeSignal(sig)
// 	}

// 	for i := 0; i < ms.groupCount; i++ {
// 		ms.groups[i].removeAll()
// 	}

// 	ms.signalGroupIDs.clear()
// 	ms.fixedSignals.clear()
// }

// // ShiftSignalLeft shifts the [Signal] with the given entity ID left by the given amount
// // and it returns the number of bits shifted.
// // It will not shift signals that are fixed or assigned to more then one group.
// func (ms *MultiplexerSignal) ShiftSignalLeft(signalEntityID EntityID, amount int) int {
// 	groupIDs, err := ms.signalGroupIDs.getValue(signalEntityID)
// 	if err != nil || len(groupIDs) > 1 {
// 		return 0
// 	}

// 	return ms.groups[groupIDs[0]].shiftLeft(signalEntityID, amount)
// }

// // ShiftSignalRight shifts the [Signal] with the given entity ID right by the given amount
// // and it returns the number of bits shifted.
// // It will not shift signals that are fixed or assigned to more then one group.
// func (ms *MultiplexerSignal) ShiftSignalRight(signalEntityID EntityID, amount int) int {
// 	groupIDs, err := ms.signalGroupIDs.getValue(signalEntityID)
// 	if err != nil || len(groupIDs) > 1 {
// 		return 0
// 	}

// 	return ms.groups[groupIDs[0]].shiftRight(signalEntityID, amount)
// }

// // GetSignalGroup returns a slice of signals present in the group
// // selected by the given group ID.
// // The signals are sorted by their start bit.
// func (ms *MultiplexerSignal) GetSignalGroup(groupID int) []Signal {
// 	if err := ms.verifyGroupID(groupID); err != nil {
// 		return []Signal{}
// 	}
// 	return ms.groups[groupID].signals
// }

// // GetSignalGroups returns a slice of groups sorted by their group ID.
// // Each group contains a slice of signals sorted by their start bit.
// func (ms *MultiplexerSignal) GetSignalGroups() [][]Signal {
// 	res := make([][]Signal, ms.groupCount)

// 	for i := 0; i < ms.groupCount; i++ {
// 		res[i] = ms.groups[i].signals
// 	}

// 	return res
// }

// // GroupCount returns the number of groups.
// func (ms *MultiplexerSignal) GroupCount() int {
// 	return ms.groupCount
// }

// // GetGroupCountSize returns the number of bits needed to select
// // the right group.
// func (ms *MultiplexerSignal) GetGroupCountSize() int {
// 	return calcSizeFromValue(ms.groupCount - 1)
// }

// // GroupSize returns the size of a group.
// func (ms *MultiplexerSignal) GroupSize() int {
// 	return ms.groupSize
// }

// // AssignAttribute assigns the given attribute/value pair to the [MultiplexerSignal].
// //
// // It returns an [ArgumentError] if the attribute is nil,
// // or an [AttributeValueError] if the value does not conform to the attribute.
// func (ms *MultiplexerSignal) AssignAttribute(attribute Attribute, value any) error {
// 	if err := ms.addAttributeAssignment(attribute, ms, value); err != nil {
// 		return ms.errorf(err)
// 	}
// 	return nil
// }

// func (ms *MultiplexerSignal) GetHigh() int {
// 	return ms.GetStartBit() + ms.GetSize() - 1
// }

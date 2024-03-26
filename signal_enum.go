package acmelib

import (
	"fmt"
	"slices"
)

type SignalEnum struct {
	*entity

	signalRefs   map[EntityID]*EnumSignal
	values       *entityCollection[*SignalEnumValue]
	valueIndexes map[int]EntityID
	maxIndex     int
}

func NewSignalEnum(name, desc string) *SignalEnum {
	return &SignalEnum{
		entity: newEntity(name, desc),

		signalRefs:   make(map[EntityID]*EnumSignal),
		values:       newEntityCollection[*SignalEnumValue](),
		valueIndexes: make(map[int]EntityID),
		maxIndex:     0,
	}
}

func (se *SignalEnum) errorf(err error) error {
	enumErr := fmt.Errorf(`signal enum "%s": %v`, se.GetName(), err)
	return enumErr
}

func (se *SignalEnum) verifyValueIndex(index int) error {
	if _, ok := se.valueIndexes[index]; ok {
		return fmt.Errorf(`index "%d" is already used`, index)
	}

	if index > se.maxIndex {
		prevSize := se.GetSize()
		newSize := calcSizeFromValue(index)

		for _, tmpSig := range se.signalRefs {
			if !tmpSig.hasParent() {
				continue
			}

			if err := tmpSig.parentMessage.signalPayload.verifyBeforeGrow(tmpSig, newSize-prevSize); err != nil {
				return fmt.Errorf(`index "%d" is invalid : %v`, index, err)
			}
		}
	}

	return nil
}

func (se *SignalEnum) modifyValueIndex(value *SignalEnumValue, newIndex int) {
	updateMaxIdx := false
	if value.index == se.maxIndex && newIndex < se.maxIndex {
		updateMaxIdx = true
	}

	delete(se.valueIndexes, value.index)
	se.valueIndexes[newIndex] = value.GetEntityID()

	if updateMaxIdx {
		currMax := 0

		for _, tmpVal := range se.GetValues() {
			tmpIdx := tmpVal.GetIndex()

			if tmpIdx > currMax {
				currMax = tmpIdx
			}
		}

		se.maxIndex = currMax
		se.setUpdateTimeNow()
	}
}

func (se *SignalEnum) addSignalRef(sig *EnumSignal) {
	se.signalRefs[sig.GetEntityID()] = sig
}

func (se *SignalEnum) removeSignalRef(sigID EntityID) {
	delete(se.signalRefs, sigID)
}

func (se *SignalEnum) AddValue(value *SignalEnumValue) error {
	if err := se.verifyValueIndex(value.index); err != nil {
		return se.errorf(err)
	}

	if err := se.values.addEntity(value); err != nil {
		return se.errorf(err)
	}

	se.valueIndexes[value.GetIndex()] = value.GetEntityID()

	value.parentEnum = se
	se.setUpdateTimeNow()

	return nil
}

func (se *SignalEnum) GetValues() []*SignalEnumValue {
	values := se.values.listEntities()
	slices.SortFunc(values, func(a *SignalEnumValue, b *SignalEnumValue) int { return a.index - b.index })
	return values
}

func (se *SignalEnum) GetSize() int {
	return calcSizeFromValue(se.maxIndex)
}

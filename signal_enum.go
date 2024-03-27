package acmelib

import (
	"fmt"
	"slices"
)

type SignalEnum struct {
	*entity

	signalRefs map[EntityID]*EnumSignal

	values       *entityCollection[*SignalEnumValue]
	valueIndexes map[int]EntityID
	maxIndex     int
}

func NewSignalEnum(name, desc string) *SignalEnum {
	return &SignalEnum{
		entity: newEntity(name, desc),

		signalRefs: make(map[EntityID]*EnumSignal),

		values:       newEntityCollection[*SignalEnumValue](),
		valueIndexes: make(map[int]EntityID),
		maxIndex:     0,
	}
}

func (se *SignalEnum) errorf(err error) error {
	enumErr := fmt.Errorf(`signal enum "%s": %v`, se.Name(), err)
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

			switch tmpSig.parent.getSignalParentKind() {
			case signalParentKindMessage:
				msgParent, err := tmpSig.parent.toParentMessage()
				if err != nil {
					panic(err)
				}
				if err := msgParent.signalPayload.verifyBeforeGrow(tmpSig, newSize-prevSize); err != nil {
					return fmt.Errorf(`index "%d" is invalid : %v`, index, err)
				}

			}
		}
	}

	return nil
}

func (se *SignalEnum) modifyValueIndex(value *SignalEnumValue, newIndex int) error {
	if err := se.verifyValueIndex(newIndex); err != nil {
		return err
	}

	if newIndex > se.maxIndex {
		se.maxIndex = newIndex

		return nil
	}

	updateMaxIdx := false
	if value.index == se.maxIndex && newIndex < se.maxIndex {
		updateMaxIdx = true
	}

	delete(se.valueIndexes, value.index)
	se.valueIndexes[newIndex] = value.EntityID()

	if updateMaxIdx {
		amount := calcSizeFromValue(newIndex) - se.GetSize()
		for _, tmpSig := range se.signalRefs {
			if err := tmpSig.modifySize(amount); err != nil {
				return err
			}
		}

		se.setMaxIndex()
	}

	return nil
}

func (se *SignalEnum) setMaxIndex() {
	currMax := 0

	for _, tmpVal := range se.GetValues() {
		tmpIdx := tmpVal.GetIndex()

		if tmpIdx > currMax {
			currMax = tmpIdx
		}
	}

	se.maxIndex = currMax
}

func (se *SignalEnum) addSignalRef(sig *EnumSignal) {
	se.signalRefs[sig.EntityID()] = sig
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

	se.valueIndexes[value.GetIndex()] = value.EntityID()

	index := value.GetIndex()
	if index > se.maxIndex {
		se.maxIndex = index
	}

	value.parentEnum = se

	return nil
}

func (se *SignalEnum) RemoveValue(valueID EntityID) error {
	val, err := se.GetValueByID(valueID)
	if err != nil {
		return se.errorf(err)
	}

	valIdx := val.index
	wasMaxIndex := false
	if valIdx == se.maxIndex {
		wasMaxIndex = true
	}

	if err := se.values.removeEntity(valueID); err != nil {
		return err
	}

	val.parentEnum = nil

	delete(se.valueIndexes, valIdx)

	if wasMaxIndex {
		se.setMaxIndex()
	}

	return nil
}

func (se *SignalEnum) GetValueByID(valueID EntityID) (*SignalEnumValue, error) {
	return se.values.getEntityByID(valueID)
}

func (se *SignalEnum) GetValues() []*SignalEnumValue {
	values := se.values.listEntities()
	slices.SortFunc(values, func(a *SignalEnumValue, b *SignalEnumValue) int { return a.index - b.index })
	return values
}

func (se *SignalEnum) GetSize() int {
	return calcSizeFromValue(se.maxIndex)
}

func (se *SignalEnum) GetMaxIndex() int {
	return se.maxIndex
}

type SignalEnumValue struct {
	*entity

	parentEnum *SignalEnum

	index int
}

func NewSignalEnumValue(name, desc string, index int) *SignalEnumValue {
	return &SignalEnumValue{
		entity: newEntity(name, desc),

		parentEnum: nil,

		index: index,
	}
}

func (sev *SignalEnumValue) errorf(err error) error {
	enumValErr := fmt.Errorf(`signal enum value "%s": %v`, sev.Name(), err)
	if sev.hasParentEnum() {
		return sev.parentEnum.errorf(enumValErr)
	}
	return enumValErr
}

func (sev *SignalEnumValue) hasParentEnum() bool {
	return sev.parentEnum != nil
}

func (sev *SignalEnumValue) GetParentEnum() *SignalEnum {
	return sev.parentEnum
}

func (sev *SignalEnumValue) GetIndex() int {
	return sev.index
}

func (sev *SignalEnumValue) UpdateName(name string) error {
	if sev.parentEnum != nil {
		if err := sev.parentEnum.values.updateEntityName(sev.entityID, sev.name, name); err != nil {
			return sev.errorf(err)
		}
	}

	return sev.entity.UpdateName(name)
}

func (sev *SignalEnumValue) UpdateIndex(newIndex int) error {
	if sev.index == newIndex {
		return nil
	}

	if sev.hasParentEnum() {
		if err := sev.parentEnum.modifyValueIndex(sev, newIndex); err != nil {
			return sev.errorf(err)
		}
	}

	sev.index = newIndex

	return nil
}

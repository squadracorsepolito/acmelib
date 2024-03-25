package acmelib

import (
	"fmt"
	"slices"
)

type SignalEnum struct {
	*entity

	signalRefs    []*EnumSignal
	values        *entityCollection[*SignalEnumValue]
	valuesByIndex map[int]EntityID
	maxIndex      int
}

func NewSignalEnum(name, desc string) *SignalEnum {
	return &SignalEnum{
		entity: newEntity(name, desc),

		signalRefs:    []*EnumSignal{},
		values:        newEntityCollection[*SignalEnumValue](),
		valuesByIndex: make(map[int]EntityID),
		maxIndex:      0,
	}
}

func (se *SignalEnum) errorf(err error) error {
	enumErr := fmt.Errorf(`enum "%s": %v`, se.Name, err)
	return enumErr
}

func (se *SignalEnum) calculateSize(maxIdx int) int {
	maxBitSize := 64

	for i := 0; i < maxBitSize; i++ {
		if maxIdx <= 1<<i {
			return i + 1
		}
	}

	return maxBitSize
}

func (se *SignalEnum) verifyValueIndex(index int) error {
	if _, ok := se.valuesByIndex[index]; ok {
		return fmt.Errorf(`index "%d" is already used`, index)
	}

	// values := se.GetValuesByIndex()
	// prevMaxIdx := se.maxIndex
	// maxIdx := prevMaxIdx

	// for valIdx, val := range values {
	// 	tmpIndex := val.index

	// 	if tmpIndex > index {
	// 		break
	// 	}

	// 	if tmpIndex == index {
	// 		return fmt.Errorf(`index "%d" is already used by "%s"`, index, val.Name)
	// 	}

	// 	if tmpIndex < index && valIdx == len(values)-1 {
	// 		maxIdx = index
	// 	}
	// }
	// if maxIdx == prevMaxIdx {
	// 	return nil
	// }

	// prevSize := se.calculateSize(prevMaxIdx)
	// maxSize := se.calculateSize(maxIdx)

	// updSignals := []*EnumSignal{}
	// for _, sigRef := range se.signalRefs {
	// 	err := sigRef.modifySize(maxSize - prevSize)
	// 	if err != nil {
	// 		for _, updSig := range updSignals {
	// 			if err := updSig.modifySize(prevSize - maxSize); err != nil {
	// 				return err
	// 			}
	// 		}

	// 		se.maxIndex = prevMaxIdx
	// 		return fmt.Errorf(`index "%d" causes signal "%s" to modify its size, but the message cannot handle it`, index, sigRef.GetName())
	// 	}
	// }

	// se.maxIndex = maxIdx

	return nil
}

func (se *SignalEnum) AddValue(value *SignalEnumValue) error {
	if err := se.verifyValueIndex(value.index); err != nil {
		return se.errorf(err)
	}

	if err := se.values.addEntity(value); err != nil {
		return se.errorf(err)
	}

	se.valuesByIndex[value.GetIndex()] = value.GetEntityID()

	value.parentEnum = se
	se.setUpdateTimeNow()

	return nil
}

func (se *SignalEnum) GetValuesByIndex() []*SignalEnumValue {
	values := se.values.listEntities()
	slices.SortFunc(values, func(a *SignalEnumValue, b *SignalEnumValue) int { return (a.index) - (b.index) })
	return values
}

func (se *SignalEnum) GetSize() int {
	return se.calculateSize(se.maxIndex)
}

type SignalEnumValue struct {
	*entity

	parentEnum *SignalEnum
	index      int
}

func NewSignalEnumValue(name, desc string, index int) *SignalEnumValue {
	return &SignalEnumValue{
		entity: newEntity(name, desc),

		index: index,
	}
}

func (sev *SignalEnumValue) errorf(err error) error {
	enumValErr := fmt.Errorf(`enum value "%s": %v`, sev.Name, err)
	if sev.parentEnum != nil {
		return sev.parentEnum.errorf(enumValErr)
	}
	return enumValErr
}

func (sev *SignalEnumValue) hasParent() bool {
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
		if err := sev.parentEnum.values.updateEntityName(sev.EntityID, sev.Name, name); err != nil {
			return sev.errorf(err)
		}
	}

	return sev.entity.UpdateName(name)
}

func (sev *SignalEnumValue) UpdateIndex(newIndex int) error {
	if sev.index == newIndex {
		return nil
	}

	if sev.hasParent() {

	}

	sev.index = newIndex
	sev.setUpdateTimeNow()

	return nil
}

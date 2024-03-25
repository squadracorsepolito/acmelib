package acmelib

import (
	"fmt"
	"slices"
)

type SignalEnumValue struct {
	*entity
	ParentEnum *SignalEnum

	Index int
}

func NewSignalEnumValue(name, desc string, index int) *SignalEnumValue {
	return &SignalEnumValue{
		entity: newEntity(name, desc),

		Index: index,
	}
}

func (sev *SignalEnumValue) errorf(err error) error {
	enumValErr := fmt.Errorf(`enum value "%s": %v`, sev.Name, err)
	if sev.ParentEnum != nil {
		return sev.ParentEnum.errorf(enumValErr)
	}
	return enumValErr
}

func (sev *SignalEnumValue) UpdateName(name string) error {
	if sev.ParentEnum != nil {
		if err := sev.ParentEnum.values.updateEntityName(sev.EntityID, sev.Name, name); err != nil {
			return sev.errorf(err)
		}
	}

	return sev.entity.UpdateName(name)
}

type SignalEnum struct {
	*entity

	signalRefs []*EnumSignal
	values     *entityCollection[*SignalEnumValue]
	maxIndex   int
}

func NewSignalEnum(name, desc string) *SignalEnum {
	return &SignalEnum{
		entity: newEntity(name, desc),

		signalRefs: []*EnumSignal{},
		values:     newEntityCollection[*SignalEnumValue](),
		maxIndex:   0,
	}
}

func (se *SignalEnum) errorf(err error) error {
	enumErr := fmt.Errorf(`enum "%s": %v`, se.Name, err)
	return enumErr
}

func (se *SignalEnum) AddValue(value *SignalEnumValue) error {
	values := se.GetValuesByIndex()

	for _, tmpVal := range values {
		if value.Index == tmpVal.Index {
			return se.errorf(fmt.Errorf(`value "%s" cannot have index "%d" because value "%s" already has`, value.Name, value.Index, tmpVal.Name))
		}
	}

	if err := se.values.addEntity(value); err != nil {
		return se.errorf(err)
	}

	if value.Index > se.maxIndex {
		se.maxIndex = value.Index
	}

	value.ParentEnum = se
	se.setUpdateTimeNow()

	return nil
}

func (se *SignalEnum) GetValuesByIndex() []*SignalEnumValue {
	values := se.values.listEntities()
	slices.SortFunc(values, func(a *SignalEnumValue, b *SignalEnumValue) int { return a.Index - b.Index })
	return values
}

func (se *SignalEnum) GetSize() int {
	maxBitSize := 64

	for i := 0; i < maxBitSize; i++ {
		if se.maxIndex <= 1<<i {
			return i + 1
		}
	}

	return maxBitSize
}

package acmelib

import "fmt"

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
	enumValErr := fmt.Errorf(`signal enum value "%s": %v`, sev.GetName(), err)
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

	if sev.hasParentEnum() {
		if err := sev.parentEnum.verifyValueIndex(newIndex); err != nil {
			return err
		}

		sev.parentEnum.modifyValueIndex(sev, newIndex)
	}

	sev.index = newIndex
	sev.setUpdateTimeNow()

	return nil
}

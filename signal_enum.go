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
	signalRefs []*enumSignal

	values *entityCollection[*SignalEnumValue]

	maxIndex int
}

func NewSignalEnum(name, desc string) *SignalEnum {
	return &SignalEnum{
		entity:     newEntity(name, desc),
		signalRefs: []*enumSignal{},

		values: newEntityCollection[*SignalEnumValue](),

		maxIndex: 0,
	}
}

func (se *SignalEnum) errorf(err error) error {
	enumErr := fmt.Errorf(`enum "%s": %v`, se.Name, err)
	return enumErr
}

func (se *SignalEnum) AddValue(value *SignalEnumValue) error {
	values := se.ValuesByIndex()

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

func (se *SignalEnum) ValuesByIndex() []*SignalEnumValue {
	values := se.values.listEntities()
	slices.SortFunc(values, func(a *SignalEnumValue, b *SignalEnumValue) int { return a.Index - b.Index })
	return values
}

type enumSignal struct {
	*entity
	ParentMessage *Message

	Enum     *SignalEnum
	StartBit int
}

func newEnumSignal(name, desc string, enum *SignalEnum) *enumSignal {
	sig := &enumSignal{
		entity: newEntity(name, desc),

		Enum: enum,
	}

	enum.signalRefs = append(enum.signalRefs, sig)

	return sig
}

func (es *enumSignal) GetBitSize() int {
	maxBitSize := 64
	enumMaxIdx := es.Enum.maxIndex

	for i := 0; i < maxBitSize; i++ {
		if enumMaxIdx <= 1<<i {
			return i + 1
		}
	}

	return maxBitSize
}

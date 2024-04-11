package acmelib

import (
	"fmt"
	"slices"
	"strings"
)

// SignalEnum is the representation of an enum that can be assigned
// to a signal.
type SignalEnum struct {
	*entity

	parentSignals *set[EntityID, *EnumSignal]
	parErrID      EntityID

	values       *set[EntityID, *SignalEnumValue]
	valueNames   *set[string, EntityID]
	valueIndexes *set[int, EntityID]

	maxIndex int
}

// NewSignalEnum creates a new [SignalEnum] with the given name and description.
func NewSignalEnum(name, desc string) *SignalEnum {
	return &SignalEnum{
		entity: newEntity(name, desc),

		parentSignals: newSet[EntityID, *EnumSignal]("parent signal"),
		parErrID:      "",

		values:       newSet[EntityID, *SignalEnumValue]("value"),
		valueNames:   newSet[string, EntityID]("value name"),
		valueIndexes: newSet[int, EntityID]("value index"),

		maxIndex: 0,
	}
}

func (se *SignalEnum) errorf(err error) error {
	enumErr := fmt.Errorf(`signal enum "%s": %w`, se.Name(), err)

	if se.parentSignals.size() > 0 {
		if se.parErrID != "" {
			parSig, err := se.parentSignals.getValue(se.parErrID)
			if err != nil {
				panic(err)
			}

			se.parErrID = ""
			return parSig.errorf(enumErr)
		}

		return se.parentSignals.getValues()[0].errorf(enumErr)
	}

	return enumErr
}

func (se *SignalEnum) modifyValueName(valEntID EntityID, newName string) {
	val, err := se.values.getValue(valEntID)
	if err != nil {
		panic(err)
	}

	oldName := val.name
	se.valueNames.modifyKey(oldName, newName, valEntID)
}

func (se *SignalEnum) verifyValueIndex(index int) error {
	if err := se.valueIndexes.verifyKey(index); err != nil {
		return err
	}

	if index > se.maxIndex {
		prevSize := se.GetSize()
		newSize := calcSizeFromValue(index)

		for _, tmpSig := range se.parentSignals.entries() {
			if !tmpSig.hasParent() {
				continue
			}

			if err := tmpSig.parent.verifySignalSizeAmount(tmpSig.entityID, newSize-prevSize); err != nil {
				se.parErrID = tmpSig.entityID
				return err
			}
		}
	}

	return nil
}

func (se *SignalEnum) modifyValueIndex(value *SignalEnumValue, newIndex int) {
	gtMaxIndex := false
	if maxSize > se.maxIndex {
		gtMaxIndex = true
	}

	updateMaxIdx := false
	if value.index == se.maxIndex && newIndex < se.maxIndex {
		updateMaxIdx = true
	}

	if gtMaxIndex || updateMaxIdx {
		amount := calcSizeFromValue(newIndex) - se.GetSize()

		for _, tmpSig := range se.parentSignals.entries() {
			if err := tmpSig.modifySize(amount); err != nil {
				panic(err)
			}
		}

		if gtMaxIndex {
			se.maxIndex = newIndex
		} else {
			se.setMaxIndex()
		}
	}

	oldIndex := value.index
	se.valueIndexes.modifyKey(oldIndex, newIndex, value.entityID)
}

func (se *SignalEnum) setMaxIndex() {
	currMax := 0

	for _, tmpVal := range se.values.entries() {
		tmpIdx := tmpVal.Index()

		if tmpIdx > currMax {
			currMax = tmpIdx
		}
	}

	se.maxIndex = currMax
}

// UpdateName updates the name of the [SignalEnum] to the given new one.
func (se *SignalEnum) UpdateName(newName string) {
	se.name = newName
}

func (se *SignalEnum) stringify(b *strings.Builder, tabs int) {
	se.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%smax_index: %d\n", tabStr, se.maxIndex))

	if se.values.size() == 0 {
		return
	}

	b.WriteString(fmt.Sprintf("%svalues:\n", tabStr))
	for _, enumVal := range se.Values() {
		enumVal.stringify(b, tabs+1)
		b.WriteRune('\n')
	}
}

func (se *SignalEnum) String() string {
	builder := new(strings.Builder)
	se.stringify(builder, 0)
	return builder.String()
}

// AddValue adds the given [SignalEnumValue] to the [SignalEnum].
// It may return an error if the value name is already in use within
// the signal enum, or if it has an invalid index.
func (se *SignalEnum) AddValue(value *SignalEnumValue) error {
	if err := se.verifyValueIndex(value.index); err != nil {
		return se.errorf(fmt.Errorf(`cannot add value "%s" : %w`, value.name, err))
	}

	if err := se.valueNames.verifyKey(value.name); err != nil {
		return se.errorf(fmt.Errorf(`cannot add value "%s" : %w`, value.name, err))
	}

	index := value.index
	if index > se.maxIndex {
		se.maxIndex = index
	}

	se.values.add(value.entityID, value)
	se.valueNames.add(value.name, value.entityID)
	se.valueIndexes.add(value.index, value.entityID)

	value.setParentEnum(se)

	return nil
}

// RemoveValue removes the [SignalEnumValue] with the given entity id from the [SignalEnum].
// It may return an error if the value with the given entity id is not found.
func (se *SignalEnum) RemoveValue(valueEntityID EntityID) error {
	val, err := se.values.getValue(valueEntityID)
	if err != nil {
		return se.errorf(fmt.Errorf(`cannot remove value with entity id "%s" : %w`, valueEntityID, err))
	}

	valIdx := val.index
	wasMaxIndex := false
	if valIdx == se.maxIndex {
		wasMaxIndex = true
	}

	val.setParentEnum(nil)

	se.values.remove(valueEntityID)
	se.valueNames.remove(val.name)
	se.valueIndexes.remove(val.index)

	if wasMaxIndex {
		se.setMaxIndex()
	}

	return nil
}

// RemoveAllValues removes all enum values from the [SignalEnum].
func (se *SignalEnum) RemoveAllValues() {
	for _, tmpVal := range se.values.entries() {
		tmpVal.setParentEnum(nil)
	}

	se.values.clear()
	se.valueNames.clear()
	se.valueIndexes.clear()
}

// Values returns a slice of all the enum values of the [SignalEnum].
func (se *SignalEnum) Values() []*SignalEnumValue {
	valueSlice := se.values.getValues()
	slices.SortFunc(valueSlice, func(a *SignalEnumValue, b *SignalEnumValue) int { return a.index - b.index })
	return valueSlice
}

// GetSize returns the size of the [SignalEnum] in bits.
func (se *SignalEnum) GetSize() int {
	return calcSizeFromValue(se.maxIndex)
}

// MaxIndex returns the highest index of the enum values of the [SignalEnum].
func (se *SignalEnum) MaxIndex() int {
	return se.maxIndex
}

// SignalEnumValue holds the key (name) and the value (index) of a signal enum
// entry.
type SignalEnumValue struct {
	*entity

	parentEnum *SignalEnum

	index int
}

// NewSignalEnumValue creates a new [SignalEnumValue] with the given name,
// description, and index.
func NewSignalEnumValue(name, desc string, index int) *SignalEnumValue {
	return &SignalEnumValue{
		entity: newEntity(name, desc),

		parentEnum: nil,

		index: index,
	}
}

func (sev *SignalEnumValue) hasParentEnum() bool {
	return sev.parentEnum != nil
}

func (sev *SignalEnumValue) setParentEnum(enum *SignalEnum) {
	sev.parentEnum = enum
}

func (sev *SignalEnumValue) errorf(err error) error {
	enumValErr := fmt.Errorf(`signal enum value "%s": %w`, sev.Name(), err)
	if sev.hasParentEnum() {
		return sev.parentEnum.errorf(enumValErr)
	}
	return enumValErr
}

func (sev *SignalEnumValue) stringify(b *strings.Builder, tabs int) {
	sev.entity.stringify(b, tabs)
	tabStr := getTabString(tabs)
	b.WriteString(fmt.Sprintf("%sindex: %d\n", tabStr, sev.index))
}

func (sev *SignalEnumValue) String() string {
	builder := new(strings.Builder)
	sev.stringify(builder, 0)
	return builder.String()
}

// UpdateName updates the name of the [SignalEnumValue] to the given new one.
// It may return an error if the new name is already in use within the parent enum.
func (sev *SignalEnumValue) UpdateName(newName string) error {
	if sev.name == newName {
		return nil
	}

	if sev.hasParentEnum() {
		if err := sev.parentEnum.valueNames.verifyKey(newName); err != nil {
			return sev.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}

		sev.parentEnum.modifyValueName(sev.entityID, newName)
	}

	sev.name = newName

	return nil
}

// ParentEnum returns the parent [SignalEnum] of the [SignalEnumValue],
// or nil if not set.
func (sev *SignalEnumValue) ParentEnum() *SignalEnum {
	return sev.parentEnum
}

// UpdateIndex updates the index of the [SignalEnumValue] to the given new one.
// It may return an error if the new index is invalid.
func (sev *SignalEnumValue) UpdateIndex(newIndex int) error {
	if sev.index == newIndex {
		return nil
	}

	if sev.hasParentEnum() {
		if err := sev.parentEnum.verifyValueIndex(newIndex); err != nil {
			return sev.errorf(fmt.Errorf(`cannot update index to "%d" : %w`, newIndex, err))
		}

		sev.parentEnum.modifyValueIndex(sev, newIndex)
	}

	sev.index = newIndex

	return nil
}

// Index returns the index of the [SignalEnumValue].
func (sev *SignalEnumValue) Index() int {
	return sev.index
}

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
	*withRefs[*EnumSignal]

	parErrID EntityID

	values       *set[EntityID, *SignalEnumValue]
	valueNames   *set[string, EntityID]
	valueIndexes *set[int, EntityID]

	maxIndex int
	minSize  int
}

func newSignalEnumFromEntity(ent *entity) *SignalEnum {
	return &SignalEnum{
		entity:   ent,
		withRefs: newWithRefs[*EnumSignal](),

		parErrID: "",

		values:       newSet[EntityID, *SignalEnumValue](),
		valueNames:   newSet[string, EntityID](),
		valueIndexes: newSet[int, EntityID](),

		maxIndex: 0,
		minSize:  1,
	}
}

// NewSignalEnum creates a new [SignalEnum] with the given name.
func NewSignalEnum(name string) *SignalEnum {
	return newSignalEnumFromEntity(newEntity(name, EntityKindSignalEnum))
}

// Clone creates a new [SignalEnum] with the same values as the current one.
func (se *SignalEnum) Clone() (*SignalEnum, error) {
	cloned := newSignalEnumFromEntity(se.entity.clone())

	for _, tmpVal := range se.values.getValues() {
		if err := cloned.AddValue(tmpVal.Clone()); err != nil {
			return nil, err
		}
	}

	return cloned, nil
}

func (se *SignalEnum) errorf(err error) error {
	enumErr := &EntityError{
		Kind:     EntityKindSignalEnum,
		EntityID: se.entityID,
		Name:     se.name,
		Err:      err,
	}

	if se.refs.size() > 0 {
		if se.parErrID != "" {
			parSig, err := se.refs.getValue(se.parErrID)
			if err != nil {
				panic(err)
			}

			se.parErrID = ""
			return parSig.errorf(enumErr)
		}

		return se.refs.getValues()[0].errorf(enumErr)
	}

	return enumErr
}

func (se *SignalEnum) verifyValueName(name string) error {
	err := se.valueNames.verifyKeyUnique(name)
	if err != nil {
		return &NameError{
			Name: name,
			Err:  err,
		}
	}
	return nil
}

func (se *SignalEnum) verifyValueIndex(index int) error {
	if err := se.valueIndexes.verifyKeyUnique(index); err != nil {
		return err
	}

	if index > se.maxIndex {
		prevSize := se.GetSize()
		newSize := calcSizeFromValue(index)

		for _, tmpSig := range se.refs.entries() {
			if tmpSig.hasParentMsg() {
				if err := tmpSig.parentMsg.verifySignalSizeAmount(tmpSig.entityID, newSize-prevSize); err != nil {
					se.parErrID = tmpSig.entityID
					return &ValueIndexError{
						Index: index,
						Err:   err,
					}
				}
			}

			if tmpSig.hasParentMuxSig() {
				if err := tmpSig.parentMuxSig.verifySignalSizeAmount(tmpSig.entityID, newSize-prevSize); err != nil {
					se.parErrID = tmpSig.entityID
					return &ValueIndexError{
						Index: index,
						Err:   err,
					}
				}
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

		for _, tmpSig := range se.refs.entries() {
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

	refCount := se.ReferenceCount()
	if refCount > 0 {
		b.WriteString(fmt.Sprintf("%sreference_count: %d\n", tabStr, refCount))
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
	if value == nil {
		return &ArgumentError{
			Name: "value",
			Err:  ErrIsNil,
		}
	}

	addValErr := &AddEntityError{
		EntityID: value.entityID,
		Name:     value.name,
	}

	if err := se.verifyValueIndex(value.index); err != nil {
		addValErr.Err = err
		return se.errorf(addValErr)
	}

	if err := se.verifyValueName(value.name); err != nil {
		addValErr.Err = err
		return se.errorf(addValErr)
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
		return se.errorf(&RemoveEntityError{
			EntityID: valueEntityID,
			Err:      err,
		})
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

// GetValue returns the [SignalEnumValue] with the given entity id.
//
// It returns a [GetEntityError] if the value with the given entity id is not found.
func (se *SignalEnum) GetValue(valueEntityID EntityID) (*SignalEnumValue, error) {
	val, err := se.values.getValue(valueEntityID)
	if err != nil {
		return nil, se.errorf(&GetEntityError{
			EntityID: valueEntityID,
			Err:      err,
		})
	}
	return val, nil
}

// GetSize returns the size of the [SignalEnum] in bits.
func (se *SignalEnum) GetSize() int {
	maxIdxSize := calcSizeFromValue(se.maxIndex)
	if se.minSize > maxIdxSize {
		return se.minSize
	}
	return maxIdxSize
}

// MaxIndex returns the highest index of the enum values of the [SignalEnum].
func (se *SignalEnum) MaxIndex() int {
	return se.maxIndex
}

// SetMinSize sets the minimum size in bit of the [SignalEnum].
// By defaul it is set to 1.
func (se *SignalEnum) SetMinSize(minSize int) {
	se.minSize = minSize
}

// MinSize return the minimum size of the [SignalEnum] in bits.
func (se *SignalEnum) MinSize() int {
	return se.minSize
}

// SignalEnumValue holds the key (name) and the value (index) of a signal enum
// entry.
type SignalEnumValue struct {
	*entity

	parentEnum *SignalEnum

	index int
}

func newSignalEnumValueFromEntity(ent *entity, index int) *SignalEnumValue {
	return &SignalEnumValue{
		entity: ent,

		parentEnum: nil,

		index: index,
	}
}

// NewSignalEnumValue creates a new [SignalEnumValue] with the given name and index.
func NewSignalEnumValue(name string, index int) *SignalEnumValue {
	return newSignalEnumValueFromEntity(newEntity(name, EntityKindSignalEnumValue), index)
}

// Clone creates a new [SignalEnumValue] with the values as the current one.
func (sev *SignalEnumValue) Clone() *SignalEnumValue {
	return newSignalEnumValueFromEntity(sev.entity.clone(), sev.index)
}

func (sev *SignalEnumValue) hasParentEnum() bool {
	return sev.parentEnum != nil
}

func (sev *SignalEnumValue) setParentEnum(enum *SignalEnum) {
	sev.parentEnum = enum
}

func (sev *SignalEnumValue) errorf(err error) error {
	enumValErr := &EntityError{
		Kind:     EntityKindSignalEnumValue,
		EntityID: sev.entityID,
		Name:     sev.name,
		Err:      err,
	}

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
		if err := sev.parentEnum.verifyValueName(newName); err != nil {
			return sev.errorf(&UpdateNameError{Err: err})
		}

		sev.parentEnum.valueNames.modifyKey(sev.name, newName, sev.entityID)
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
			return sev.errorf(&UpdateIndexError{Err: err})
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

package acmelib

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

type SignalEnumValue0 struct {
	index int
	name  string
	desc  string

	parentEnum *SignalEnum
}

func newSignalEnumValue0(index int, name string) *SignalEnumValue0 {
	return &SignalEnumValue0{
		index: index,
		name:  name,
		desc:  "",

		parentEnum: nil,
	}
}

func (v *SignalEnumValue0) setParentEnum(enum *SignalEnum) {
	v.parentEnum = enum
}

func (v *SignalEnumValue0) stringify(s *stringer.Stringer) {
	s.Write("index: %d; name: %q\n", v.index, v.name)

	if len(v.desc) > 0 {
		s.Write(" desc: %q\n", v.desc)
	}
}

func (v *SignalEnumValue0) String() string {
	s := stringer.New()
	s.Write("signal_enum_value:\n")
	v.stringify(s)
	return s.String()
}

func (v *SignalEnumValue0) Index() int {
	return v.index
}

func (v *SignalEnumValue0) UpdateIndex(newIndex int) error {
	if v.parentEnum == nil || v.index == newIndex {
		return nil
	}

	return v.parentEnum.updateValueIndex(v, newIndex)
}

func (v *SignalEnumValue0) Name() string {
	return v.name
}

func (v *SignalEnumValue0) SetName(name string) {
	v.name = name
}

func (v *SignalEnumValue0) Desc() string {
	return v.desc
}

func (v *SignalEnumValue0) SetDesc(desc string) {
	v.desc = desc
}

// SignalEnum is the representation of an enum that can be assigned
// to a signal.
type SignalEnum struct {
	*entity
	*withRefs[*EnumSignal]

	parErrID EntityID

	values       *set[EntityID, *SignalEnumValue]
	valueNames   *set[string, EntityID]
	valueIndexes *set[int, EntityID]

	values0       []*SignalEnumValue0
	valueIndexes0 *collection.Set[int]

	size int

	fixedSize bool

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

		values0:       []*SignalEnumValue0{},
		valueIndexes0: collection.NewSet[int](),

		size: 1,

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

	// if index > se.maxIndex {
	// 	prevSize := se.GetSize()
	// 	newSize := calcSizeFromValue(index)

	// 	for _, tmpSig := range se.refs.entries() {
	// 		if tmpSig.hasParentMsg() {
	// 			if err := tmpSig.parentMsg.verifySignalSizeAmount(tmpSig.entityID, newSize-prevSize); err != nil {
	// 				se.parErrID = tmpSig.entityID
	// 				return &ValueIndexError{
	// 					Index: index,
	// 					Err:   err,
	// 				}
	// 			}
	// 		}

	// 		if tmpSig.hasParentMuxSig() {
	// 			if err := tmpSig.parentMuxSig.verifySignalSizeAmount(tmpSig.entityID, newSize-prevSize); err != nil {
	// 				se.parErrID = tmpSig.entityID
	// 				return &ValueIndexError{
	// 					Index: index,
	// 					Err:   err,
	// 				}
	// 			}
	// 		}
	// 	}
	// }

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
		amount := getSizeFromValue(newIndex) - se.GetSize()

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

// func (se *SignalEnum) stringifyOld(b *strings.Builder, tabs int) {
// 	se.entity.stringifyOld(b, tabs)

// 	tabStr := getTabString(tabs)

// 	b.WriteString(fmt.Sprintf("%smax_index: %d\n", tabStr, se.maxIndex))

// 	if se.values.size() == 0 {
// 		return
// 	}

// 	b.WriteString(fmt.Sprintf("%svalues:\n", tabStr))
// 	for _, enumVal := range se.Values() {
// 		enumVal.stringify(b, tabs+1)
// 		b.WriteRune('\n')
// 	}

// 	refCount := se.ReferenceCount()
// 	if refCount > 0 {
// 		b.WriteString(fmt.Sprintf("%sreference_count: %d\n", tabStr, refCount))
// 	}
// }

// func (se *SignalEnum) String() string {
// 	builder := new(strings.Builder)
// 	se.stringifyOld(builder, 0)
// 	return builder.String()
// }

// AddValue adds the given [SignalEnumValue] to the [SignalEnum].
// It may return an error if the value name is already in use within
// the signal enum, or if it has an invalid index.
func (se *SignalEnum) AddValue(value *SignalEnumValue) error {
	if value == nil {
		return &ArgError{
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

//
//
//
//
//

func (se *SignalEnum) stringify(s *stringer.Stringer) {
	se.entity.stringify(s)

	s.Write("max_index: %d\n", se.maxIndex)
	s.Write("fixed_size: %t\n", se.fixedSize)
	s.Write("size: %d\n", se.size)

	if se.values.size() == 0 {
		s.Write("values:\n")
		s.Indent()
		for _, val := range se.values0 {
			val.stringify(s)
		}
		s.Unindent()
	}

	refCount := se.ReferenceCount()
	if refCount > 0 {
		s.Write("reference_count: %d\n", refCount)
	}
}

func (se *SignalEnum) String() string {
	s := stringer.New()
	s.Write("signal_enum:\n")
	se.stringify(s)
	return s.String()
}

func (se *SignalEnum) sortValues() {
	slices.SortFunc(se.values0, func(a, b *SignalEnumValue0) int {
		return cmp.Compare(a.index, b.index)
	})
}

// updateValueIndex updates the index of the given enum value.
func (se *SignalEnum) updateValueIndex(val *SignalEnumValue0, newIndex int) error {
	// Check if the new index is valid
	if err := se.verifyIndex(newIndex); err != nil {
		return se.errorf(err)
	}

	// Swap the indexes
	oldIndex := val.index
	se.valueIndexes0.Delete(oldIndex)
	se.valueIndexes0.Add(newIndex)
	val.index = newIndex

	se.sortValues()
	se.genMaxIndex()

	return nil
}

// verifyRefNewSize checks if each referenced signal can grow to the new size.
func (se *SignalEnum) verifyRefNewSize(newSize int) error {
	for _, ref := range se.refs.entries() {
		if err := ref.verifyNewSize(newSize); err != nil {
			se.parErrID = ref.entityID
			return err
		}
	}
	return nil
}

// verifyIndex checks if the given index is valid.
func (se *SignalEnum) verifyIndex(index int) error {
	if index < 0 {
		return newIndexError(index, ErrIsNegative)
	}

	if se.valueIndexes0.Has(index) {
		return newIndexError(index, ErrIsDuplicated)
	}

	indexSize := getSizeFromValue(index)

	// If the size is fixed and the new size is larger than the current size, return an error
	if se.fixedSize && indexSize > se.size {
		return newIndexError(index, ErrOutOfBounds)
	}

	if indexSize > se.size {
		// Check if each referenced signal can grow to the new size
		if err := se.verifyRefNewSize(indexSize); err != nil {
			return newIndexError(index, err)
		}
	}

	return nil
}

// updateSize updates the size of the enum and all referenced signals.
func (se *SignalEnum) updateSize(newSize int) {
	for _, ref := range se.refs.entries() {
		ref.updateSize(newSize)
	}

	se.size = newSize
}

// Size returns the size of the enum in bits.
func (se *SignalEnum) Size() int {
	return se.size
}

// SetFixedSize sets whether the size of the enum is fized.
// If it is set to false, it will resize the enum to the size of the largest index.
func (se *SignalEnum) SetFixedSize(fixedSize bool) {
	se.fixedSize = fixedSize

	if !fixedSize {
		se.size = getSizeFromValue(se.maxIndex)
	}
}

// UpdateSize updates the size of the enum, but only when it is set to use a fixed size.
//
// It retruns [SizeError] if the new size is invalid.
func (se *SignalEnum) UpdateSize(newSize int) error {
	if !se.fixedSize {
		return nil
	}

	// Check if the new size is too small
	if newSize < getSizeFromValue(se.maxIndex) {
		return se.errorf(newSizeError(newSize, ErrTooSmall))
	}

	if err := se.verifyRefNewSize(newSize); err != nil {
		return se.errorf(err)
	}

	se.updateSize(newSize)

	return nil
}

// SetName sets the name of the enum.
func (se *SignalEnum) SetName(newName string) {
	se.name = newName
}

func (se *SignalEnum) addValue(val *SignalEnumValue0) {
	se.values0 = append(se.values0, val)
	se.valueIndexes0.Add(val.index)
	val.setParentEnum(se)

	se.sortValues()
}

// AddValue0 creates a new [SignalEnumValue0] with the given index and name and adds it to the enum.
// The index must be unique, but the name can be duplicated.
//
// It returns [IndexError] if the index is invalid.
func (se *SignalEnum) AddValue0(index int, name string) (*SignalEnumValue0, error) {
	// Check if the index is valid
	if err := se.verifyIndex(index); err != nil {
		return nil, se.errorf(err)
	}

	// Create and add the value
	val := newSignalEnumValue0(index, name)
	se.addValue(val)

	// Generate the max index
	se.genMaxIndex()

	return val, nil
}

// DeleteValue deletes the value with the given index.
func (se *SignalEnum) DeleteValue(index int) {
	se.values0 = slices.DeleteFunc(se.values0, func(v *SignalEnumValue0) bool {
		if v.index != index {
			return false
		}

		v.setParentEnum(nil)
		se.valueIndexes0.Delete(index)

		return true
	})

	// Generate the max index
	se.genMaxIndex()
}

// GetValue0 returns the [SignalEnumValue0] with the given index.
// It returns nil if there is no value with the given index.
func (se *SignalEnum) GetValue0(index int) *SignalEnumValue0 {
	if !se.valueIndexes0.Has(index) {
		return nil
	}

	for _, val := range se.values0 {
		if val.index == index {
			return val
		}
	}

	return nil
}

// genMaxIndex generates the max index of the enum
// and if it gets updated and the size is not fixed, it will update the size.
func (se *SignalEnum) genMaxIndex() {
	if len(se.values0) == 0 {
		se.maxIndex = 0
		if !se.fixedSize {
			se.updateSize(1)
		}

		return
	}

	// Get the last value and check if it is the max index
	lastVal := se.values0[len(se.values0)-1]
	if se.maxIndex == lastVal.index {
		return
	}

	// Update the size only if it is not fixed
	newSize := getSizeFromValue(lastVal.index)
	if !se.fixedSize {
		se.updateSize(newSize)
	}

	se.maxIndex = lastVal.index
}

// Clear removes all values from the enum.
func (se *SignalEnum) Clear() {
	se.values0 = []*SignalEnumValue0{}
	se.valueIndexes0.Clear()
	se.genMaxIndex()
}

//
//
//
//
//

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
	maxIdxSize := getSizeFromValue(se.maxIndex)
	if se.minSize > maxIdxSize {
		return se.minSize
	}
	return maxIdxSize
}

// MaxIndex returns the highest index of the enum values of the [SignalEnum].
func (se *SignalEnum) MaxIndex() int {
	return se.maxIndex
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

// SetMinSize sets the minimum size in bit of the [SignalEnum].
// By defaul it is set to 1.
func (se *SignalEnum) SetMinSize(minSize int) {
	se.minSize = minSize
}

func (sev *SignalEnumValue) stringify(b *strings.Builder, tabs int) {
	sev.entity.stringifyOld(b, tabs)
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

package acmelib

import (
	"cmp"
	"slices"

	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// SignalEnumValue represents a value of a [SignalEnum].
type SignalEnumValue struct {
	index int
	name  string
	desc  string

	parentEnum *SignalEnum
}

func newSignalEnumValue(index int, name string) *SignalEnumValue {
	return &SignalEnumValue{
		index: index,
		name:  name,
		desc:  "",

		parentEnum: nil,
	}
}

func (v *SignalEnumValue) setParentEnum(enum *SignalEnum) {
	v.parentEnum = enum
}

// ParentEnum returns the enum that contains this value.
func (v *SignalEnumValue) ParentEnum() *SignalEnum {
	return v.parentEnum
}

func (v *SignalEnumValue) stringify(s *stringer.Stringer) {
	s.Write("index: %d; name: %q\n", v.index, v.name)

	if len(v.desc) > 0 {
		s.Write(" desc: %q\n", v.desc)
	}
}

func (v *SignalEnumValue) String() string {
	s := stringer.New()
	s.Write("signal_enum_value:\n")
	v.stringify(s)
	return s.String()
}

// Index returns the index of the value.
func (v *SignalEnumValue) Index() int {
	return v.index
}

// UpdateIndex updates the index of the value.
//
// It returns [IndexError] if the new index is invalid.
func (v *SignalEnumValue) UpdateIndex(newIndex int) error {
	if v.parentEnum == nil || v.index == newIndex {
		return nil
	}

	return v.parentEnum.updateValueIndex(v, newIndex)
}

// Name returns the name of the value.
func (v *SignalEnumValue) Name() string {
	return v.name
}

// SetName updates the name of the value.
func (v *SignalEnumValue) SetName(name string) {
	v.name = name
}

// Desc returns the description of the value.
func (v *SignalEnumValue) Desc() string {
	return v.desc
}

// SetDesc updates the description of the value.
func (v *SignalEnumValue) SetDesc(desc string) {
	v.desc = desc
}

// SignalEnum is the representation of an enum that can be assigned
// to a signal.
type SignalEnum struct {
	*entity
	*withRefs[*EnumSignal]

	parErrID EntityID

	values       []*SignalEnumValue
	valueIndexes *collection.Set[int]

	maxIndex int

	size      int
	fixedSize bool
}

func newSignalEnumFromEntity(ent *entity) *SignalEnum {
	return &SignalEnum{
		entity:   ent,
		withRefs: newWithRefs[*EnumSignal](),

		parErrID: "",

		values:       []*SignalEnumValue{},
		valueIndexes: collection.NewSet[int](),

		maxIndex: 0,

		size:      1,
		fixedSize: false,
	}
}

// NewSignalEnum creates a new [SignalEnum] with the given name.
func NewSignalEnum(name string) *SignalEnum {
	return newSignalEnumFromEntity(newEntity(name, EntityKindSignalEnum))
}

func (se *SignalEnum) errorf(err error) error {
	enumErr := &EntityError{
		Kind:     EntityKindSignalEnum,
		EntityID: se.entityID,
		Name:     se.name,
		Err:      err,
	}

	if se.refs.Size() > 0 {
		if se.parErrID != "" {
			parSig, ok := se.refs.Get(se.parErrID)
			if !ok {
				return enumErr
			}

			se.parErrID = ""
			return parSig.errorf(enumErr)
		}

		return slices.Collect(se.refs.Values())[0].errorf(enumErr)
	}

	return enumErr
}

func (se *SignalEnum) stringify(s *stringer.Stringer) {
	se.entity.stringify(s)

	s.Write("max_index: %d\n", se.maxIndex)
	s.Write("fixed_size: %t\n", se.fixedSize)
	s.Write("size: %d\n", se.size)

	if len(se.values) > 0 {
		s.Write("values:\n")
		s.Indent()
		for _, val := range se.values {
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

// SetName sets the name of the enum.
func (se *SignalEnum) SetName(newName string) {
	se.name = newName
}

// updateValueIndex updates the index of the given enum value.
func (se *SignalEnum) updateValueIndex(val *SignalEnumValue, newIndex int) error {
	// Check if the new index is valid
	if err := se.verifyIndex(newIndex); err != nil {
		return se.errorf(err)
	}

	// Swap the indexes
	oldIndex := val.index
	se.valueIndexes.Delete(oldIndex)
	se.valueIndexes.Add(newIndex)
	val.index = newIndex

	se.sortValues()
	se.genMaxIndex()

	return nil
}

// verifyRefNewSize checks if each referenced signal can grow to the new size.
func (se *SignalEnum) verifyRefNewSize(newSize int) error {
	for ref := range se.refs.Values() {
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

	if se.valueIndexes.Has(index) {
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
	for ref := range se.refs.Values() {
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

func (se *SignalEnum) sortValues() {
	slices.SortFunc(se.values, func(a, b *SignalEnumValue) int {
		return cmp.Compare(a.index, b.index)
	})
}

func (se *SignalEnum) addValue(val *SignalEnumValue) {
	se.values = append(se.values, val)
	se.valueIndexes.Add(val.index)
	val.setParentEnum(se)

	se.sortValues()
}

// AddValue creates a new [SignalEnumValue] with the given index and name and adds it to the enum.
// The index must be unique, but the name can be duplicated.
//
// It returns [IndexError] if the index is invalid.
func (se *SignalEnum) AddValue(index int, name string) (*SignalEnumValue, error) {
	// Check if the index is valid
	if err := se.verifyIndex(index); err != nil {
		return nil, se.errorf(err)
	}

	// Create and add the value
	val := newSignalEnumValue(index, name)
	se.addValue(val)

	// Generate the max index
	se.genMaxIndex()

	return val, nil
}

// DeleteValue deletes the value with the given index.
func (se *SignalEnum) DeleteValue(index int) {
	se.values = slices.DeleteFunc(se.values, func(v *SignalEnumValue) bool {
		if v.index != index {
			return false
		}

		v.setParentEnum(nil)
		se.valueIndexes.Delete(index)

		return true
	})

	// Generate the max index
	se.genMaxIndex()
}

// GetValue0 returns the [SignalEnumValue] with the given index.
// It returns nil if there is no value with the given index.
func (se *SignalEnum) GetValue0(index int) *SignalEnumValue {
	if !se.valueIndexes.Has(index) {
		return nil
	}

	for _, val := range se.values {
		if val.index == index {
			return val
		}
	}

	return nil
}

// genMaxIndex generates the max index of the enum
// and if it gets updated and the size is not fixed, it will update the size.
func (se *SignalEnum) genMaxIndex() {
	if len(se.values) == 0 {
		se.maxIndex = 0
		if !se.fixedSize {
			se.updateSize(1)
		}

		return
	}

	// Get the last value and check if it is the max index
	lastVal := se.values[len(se.values)-1]
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
	se.values = []*SignalEnumValue{}
	se.valueIndexes.Clear()
	se.genMaxIndex()
}

// MaxIndex returns the maximum index of the enum.
func (se *SignalEnum) MaxIndex() int {
	return se.maxIndex
}

// ToSignalEnum returns the enum itself.
func (se *SignalEnum) ToSignalEnum() (*SignalEnum, error) {
	return se, nil
}

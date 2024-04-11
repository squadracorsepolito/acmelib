package acmelib

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// SignalKind rappresents the kind of a [Signal].
// It can be standard, enum, or multiplexer
type SignalKind string

const (
	// SignalKindStandard defines a standard signal.
	SignalKindStandard SignalKind = "signal-standard"
	// SignalKindEnum defines a enum signal.
	SignalKindEnum SignalKind = "signal-enum"
	// SignalKindMultiplexer defines a multiplexer signal.
	SignalKindMultiplexer SignalKind = "signal-multiplexer"
)

// SignalParentKind rappresents the kind of a [SignalParent].
// It can be message or multiplexer signal.
type SignalParentKind string

const (
	// SignalParentKindMessage defines a message parent.
	SignalParentKindMessage SignalParentKind = "signal_parent-message"
	// SignalParentKindMultiplexerSignal defines a multiplexer signal parent.
	SignalParentKindMultiplexerSignal SignalParentKind = "signal_parent-multiplexer_signal"
)

// SignalParent interface specifies the common methods of
// a parent of a signal.
type SignalParent interface {
	errorf(err error) error

	// GetSignalParentKind return the [SignalParentKind] of the parent.
	GetSignalParentKind() SignalParentKind

	verifySignalName(sigID EntityID, name string) error
	modifySignalName(sigID EntityID, newName string) error

	verifySignalSizeAmount(sigID EntityID, amount int) error
	modifySignalSize(sigID EntityID, amount int) error

	// ToParentMessage returns the signal parent as a [Message].
	ToParentMessage() (*Message, error)
	// ToParentMultiplexerSignal returns the signal parent as a [MultiplexerSignal].
	ToParentMultiplexerSignal() (*MultiplexerSignal, error)
}

// Signal interface specifies all common methods of
// [StandardSignal], [EnumSignal], and [MultiplexerSignal].
type Signal interface {
	// EntityID returns the entity id of the signal.
	EntityID() EntityID
	// Name returns the name of the signal.
	Name() string
	// Desc returns the description of the signal.
	Desc() string
	// CreateTime returns the creation time of the signal.
	CreateTime() time.Time

	// AddAttributeValue adds to the signal an [Attribute] and its value.
	AddAttributeValue(attribute Attribute, value any) error
	// RemoveAttributeValue removes from the signal an [Attribute] and its value.
	RemoveAttributeValue(attributeEntityID EntityID) error
	// RemoveAllAttributeValues removes all [Attribute] and their values from the signal.
	RemoveAllAttributeValues()
	// AttributeValues returns all [Attribute] and their values from the signal.
	AttributeValues() []*AttributeValue
	// GetAttributeValue returns the value of an [Attribute] and its value from the signal.
	GetAttributeValue(attributeEntityID EntityID) (*AttributeValue, error)

	stringify(b *strings.Builder, tabs int)
	String() string

	// Kind returns the kind of the signal.
	Kind() SignalKind

	// Parent returns the parent of the signal.
	Parent() SignalParent
	setParent(parent SignalParent)

	// GetStartBit returns the start bit of the signal.
	GetStartBit() int
	getRelStartBit() int
	setRelStartBit(startBit int)

	// GetSize returns the size of the signal.
	GetSize() int

	// ToStandard returns the signal as a [StandardSignal].
	ToStandard() (*StandardSignal, error)
	// ToEnum returns the signal as a [EnumSignal].
	ToEnum() (*EnumSignal, error)
	// ToMultiplexer returns the signal as a [MultiplexerSignal].
	ToMultiplexer() (*MultiplexerSignal, error)
}

type signal struct {
	*attributeEntity

	parent SignalParent

	kind        SignalKind
	relStartBit int
}

func newSignal(name, desc string, kind SignalKind) *signal {
	return &signal{
		attributeEntity: newAttributeEntity(name, desc, AttributeRefKindSignal),

		parent: nil,

		kind:        kind,
		relStartBit: 0,
	}
}

func (s *signal) hasParent() bool {
	return s.parent != nil
}

func (s *signal) modifySize(amount int) error {
	if s.hasParent() {
		return s.parent.modifySignalSize(s.EntityID(), amount)
	}

	return nil
}

func (s *signal) errorf(err error) error {
	sigErr := fmt.Errorf(`signal "%s": %w`, s.name, err)
	if s.hasParent() {
		return s.parent.errorf(sigErr)
	}
	return sigErr
}

func (s *signal) Kind() SignalKind {
	return s.kind
}

func (s *signal) Parent() SignalParent {
	return s.parent
}

func (s *signal) setParent(parent SignalParent) {
	s.parent = parent
}

func (s *signal) getRelStartBit() int {
	return s.relStartBit
}

func (s *signal) setRelStartBit(startBit int) {
	s.relStartBit = startBit
}

func (s *signal) stringify(b *strings.Builder, tabs int) {
	s.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%skind: %s\n", tabStr, s.kind))
	b.WriteString(fmt.Sprintf("%sstart_bit: %d; ", tabStr, s.relStartBit))
}

func (s *signal) GetStartBit() int {
	if s.hasParent() {
		if s.parent.GetSignalParentKind() == SignalParentKindMultiplexerSignal {
			muxParent, err := s.parent.ToParentMultiplexerSignal()
			if err != nil {
				panic(err)
			}

			return muxParent.GetStartBit() + muxParent.SelectSize() + s.relStartBit
		}
	}

	return s.relStartBit
}

func (s *signal) UpdateName(newName string) error {
	if s.name == newName {
		return nil
	}

	if s.hasParent() {
		if err := s.parent.verifySignalName(s.entityID, newName); err != nil {
			return s.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}

		if err := s.parent.modifySignalName(s.entityID, newName); err != nil {
			return s.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}
	}

	s.name = newName

	return nil
}

// StandardSignal is the representation of a normal signal that has a [SignalType],
// a min, a max, an offset, a scale, and can have a [SignalUnit].
type StandardSignal struct {
	*signal

	typ    *SignalType
	min    float64
	max    float64
	offset float64
	scale  float64
	unit   *SignalUnit
}

// NewStandardSignal creates a new [StandardSignal] with the given name, description,
// [SignalType], min, max, offset, scale, and unit.
// It may return an error if the given [SignalType] is nil.
func NewStandardSignal(name, desc string, typ *SignalType, min, max, offset, scale float64, unit *SignalUnit) (*StandardSignal, error) {
	if typ == nil {
		return nil, errors.New("signal type cannot be nil")
	}

	return &StandardSignal{
		signal: newSignal(name, desc, SignalKindStandard),

		typ:    typ,
		min:    min,
		max:    max,
		offset: offset,
		scale:  scale,
		unit:   unit,
	}, nil
}

// GetSize returns the size of the [StandardSignal].
func (ss *StandardSignal) GetSize() int {
	return ss.typ.size
}

// ToStandard returns the [StandardSignal] itself.
func (ss *StandardSignal) ToStandard() (*StandardSignal, error) {
	return ss, nil
}

// ToEnum always returns an error, because a [StandardSignal] cannot be converted to an [EnumSignal].
func (ss *StandardSignal) ToEnum() (*EnumSignal, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindEnum, SignalKindStandard)
}

// ToMultiplexer always returns an error, because a [StandardSignal] cannot be converted to a [MultiplexerSignal].
func (ss *StandardSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindMultiplexer, SignalKindStandard)
}

func (ss *StandardSignal) stringify(b *strings.Builder, tabs int) {
	ss.signal.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("size: %d\n", ss.GetSize()))
	b.WriteString(fmt.Sprintf("%smin: %g; max: %g; offset: %g; scale: %g\n", tabStr, ss.min, ss.max, ss.offset, ss.scale))

	b.WriteString(fmt.Sprintf("%stype:\n", tabStr))
	ss.typ.stringify(b, tabs+1)

	if ss.unit != nil {
		b.WriteString(fmt.Sprintf("%sunit:\n", tabStr))
		ss.unit.stringify(b, tabs+1)
	}
}

func (ss *StandardSignal) String() string {
	builder := new(strings.Builder)
	ss.stringify(builder, 0)
	return builder.String()
}

// Type returns the [SignalType] of the [StandardSignal].
func (ss *StandardSignal) Type() *SignalType {
	return ss.typ
}

// SetType sets the [SignalType] of the [StandardSignal] to the given [SignalType], min, max,
// scale, and offset.
// It may return an error if the given [SignalType] is nil, or if the new signal type
// size cannot fit in the message payload.
func (ss *StandardSignal) SetType(typ *SignalType, min, max, offset, scale float64) error {
	if typ == nil {
		return errors.New("signal type cannot be nil")
	}

	if err := ss.modifySize(typ.size - ss.typ.size); err != nil {
		return ss.errorf(err)
	}

	ss.typ = typ
	ss.min = min
	ss.max = max
	ss.offset = offset
	ss.scale = scale

	return nil
}

// Min returns the minimum value of the [StandardSignal].
// It may differ from the minimum value of the signal type associated
// with the [StandardSignal].
func (ss *StandardSignal) Min() float64 {
	return ss.min
}

// Max returns the maximum value of the [StandardSignal].
// It may differ from the maximum value of the signal type associated
// with the [StandardSignal].
func (ss *StandardSignal) Max() float64 {
	return ss.max
}

// Offset returns the offset of the [StandardSignal].
func (ss *StandardSignal) Offset() float64 {
	return ss.offset
}

// Scale returns the scale of the [StandardSignal].
func (ss *StandardSignal) Scale() float64 {
	return ss.scale
}

// Unit returns the [SignalUnit] of the [StandardSignal].
func (ss *StandardSignal) Unit() *SignalUnit {
	return ss.unit
}

// SetUnit sets the [SignalUnit] of the [StandardSignal] to the given one.
func (ss *StandardSignal) SetUnit(unit *SignalUnit) {
	ss.unit = unit
}

// EnumSignal is a signal that holds a [SignalEnum].
type EnumSignal struct {
	*signal

	enum *SignalEnum
}

// NewEnumSignal creates a new [EnumSignal] with the given name, description,
// and [SignalEnum].
// It may return an error if the given [SignalEnum] is nil.
func NewEnumSignal(name, desc string, enum *SignalEnum) (*EnumSignal, error) {
	if enum == nil {
		return nil, errors.New("signal enum cannot be nil")
	}

	sig := &EnumSignal{
		signal: newSignal(name, desc, SignalKindEnum),

		enum: enum,
	}

	enum.parentSignals.add(sig.entityID, sig)

	return sig, nil
}

// GetSize returns the size of the [EnumSignal].
func (es *EnumSignal) GetSize() int {
	return es.enum.GetSize()
}

// ToStandard always returns an error, because an [EnumSignal] cannot be converted to a [StandardSignal].
func (es *EnumSignal) ToStandard() (*StandardSignal, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindStandard, SignalKindEnum)
}

// ToEnum returns the [EnumSignal] itself.
func (es *EnumSignal) ToEnum() (*EnumSignal, error) {
	return es, nil
}

// ToMultiplexer always returns an error, because an [EnumSignal] cannot be converted to a [MultiplexerSignal].
func (es *EnumSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindMultiplexer, SignalKindEnum)
}

func (es *EnumSignal) stringify(b *strings.Builder, tabs int) {
	es.signal.stringify(b, tabs)
	b.WriteString(fmt.Sprintf("size: %d\n", es.GetSize()))

	tabStr := getTabString(tabs)
	b.WriteString(fmt.Sprintf("%senum:\n", tabStr))

	es.enum.stringify(b, tabs+1)
}

func (es *EnumSignal) String() string {
	builder := new(strings.Builder)
	es.stringify(builder, 0)
	return builder.String()
}

// Enum returns the [SignalEnum] of the [EnumSignal].
func (es *EnumSignal) Enum() *SignalEnum {
	return es.enum
}

// SetEnum sets the [SignalEnum] of the [EnumSignal] to the given one.
// It may return an error if the given [SignalEnum] is nil, or if the new enum
// size cannot fit in the message payload.
func (es *EnumSignal) SetEnum(enum *SignalEnum) error {
	if enum == nil {
		return errors.New("signal enum cannot be nil")
	}

	if err := es.modifySize(enum.GetSize() - es.GetSize()); err != nil {
		return es.errorf(err)
	}

	es.enum = enum

	es.enum.parentSignals.remove(es.entityID)

	enum.parentSignals.add(es.entityID, es)

	return nil
}

// MultiplexerSignal is a signal that holds other signals that are
// multiplexed by the first n bits specified with the select size attribute.
// It can multiplex all the kinds of signals (Standard, Enum, Multiplexer),
// so it is possible to create multiple levels of multiplexing.
type MultiplexerSignal struct {
	*signal

	muxSignals     *set[EntityID, Signal]
	muxSignalNames *set[string, EntityID]

	muxSignalSelValues map[EntityID]int

	signalPayloads map[int]*signalPayload

	selValRanges map[int]int

	totalSize  int
	selectSize int
}

// NewMultiplexerSignal creates a new [MultiplexerSignal] with the given name, description,
// total size, and select size.
// The select size defines the number bits used for selecting the different groups of signals
// of the multiplexer (select size = log2(number of groups)).
// The total size is the sum of the select and the maximum size of groups.
// Ex. selectSize = 2, totalSize = 10 means that the [MultiplexerSignal] can have
// 4 groups of 8 bits.
// It may return an error if the select size is greater then the total size, or if
// the total and select size are lower or equal to zero.
func NewMultiplexerSignal(name, desc string, totalSize, selectSize int) (*MultiplexerSignal, error) {
	if selectSize <= 0 {
		return nil, fmt.Errorf("the select size cannot be lower or equal to 0")
	}

	if totalSize <= 0 {
		return nil, fmt.Errorf("the total size cannot be lower or equal to 0")
	}

	if selectSize > totalSize {
		return nil, fmt.Errorf("the select size cannot be greater then the total size")
	}

	return &MultiplexerSignal{
		signal: newSignal(name, desc, SignalKindMultiplexer),

		muxSignals:     newSet[EntityID, Signal]("multiplexed signal"),
		muxSignalNames: newSet[string, EntityID]("multiplexed signal name"),

		muxSignalSelValues: make(map[EntityID]int),

		signalPayloads: make(map[int]*signalPayload),

		selValRanges: make(map[int]int),

		totalSize:  totalSize,
		selectSize: selectSize,
	}, nil
}

func (ms *MultiplexerSignal) addSignalPayload(selVal int) *signalPayload {
	payload := newSignalPayload(ms.totalSize - ms.selectSize)
	ms.signalPayloads[selVal] = payload
	return payload
}

func (ms *MultiplexerSignal) getSignalPayload(selVal int) (*signalPayload, int) {
	tmpSelVal := selVal

	if len(ms.selValRanges) > 0 {
		if rangeSelVal, ok := ms.selValRanges[selVal]; ok {
			tmpSelVal = rangeSelVal
		}
	}

	if payload, ok := ms.signalPayloads[tmpSelVal]; ok {
		return payload, tmpSelVal
	}

	return nil, tmpSelVal
}

func (ms *MultiplexerSignal) addMuxSignalName(sigID EntityID, name string) {
	ms.muxSignalNames.add(name, sigID)

	if ms.hasParent() {
		parent := ms.Parent()
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case SignalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case SignalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				msgParent.signalNames.add(name, sigID)
				return
			}
		}
	}
}

func (ms *MultiplexerSignal) removeMuxSignalName(name string) {
	ms.muxSignalNames.remove(name)

	if ms.hasParent() {
		parent := ms.Parent()
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case SignalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case SignalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				msgParent.signalNames.remove(name)
				return
			}
		}
	}
}

func (ms *MultiplexerSignal) addMuxSignal(selValue int, sig Signal) {
	id := sig.EntityID()

	ms.muxSignals.add(id, sig)
	ms.addMuxSignalName(id, sig.Name())
	ms.muxSignalSelValues[id] = selValue

	if ms.hasParent() {
		parent := ms.Parent()
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case SignalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case SignalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				msgParent.signals.add(id, sig)
				return
			}
		}
	}

}

func (ms *MultiplexerSignal) removeMuxSignal(sigID EntityID) {
	delete(ms.muxSignalSelValues, sigID)
	ms.muxSignals.remove(sigID)

	if ms.hasParent() {
		parent := ms.Parent()
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case SignalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case SignalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				msgParent.signals.remove(sigID)
				return
			}
		}
	}
}

func (ms *MultiplexerSignal) getMuxSignalSelValue(sigID EntityID) (int, error) {
	if selVal, ok := ms.muxSignalSelValues[sigID]; ok {
		return selVal, nil
	}
	return -1, fmt.Errorf(`select value for multiplexed signal with id "%s" not found`, sigID)
}

func (ms *MultiplexerSignal) verifySelectValue(selVal int) error {
	if selVal < 0 {
		return errors.New("select value cannot be negative")
	}

	if calcSizeFromValue(selVal) > ms.selectSize {
		return fmt.Errorf(`select value "%d" size exceeds the max select value size ("%d")`, selVal, ms.selectSize)
	}

	return nil
}

// GetSignalParentKind always returns [SignalParentKindMultiplexerSignal].
func (ms *MultiplexerSignal) GetSignalParentKind() SignalParentKind {
	return SignalParentKindMultiplexerSignal
}

func (ms *MultiplexerSignal) modifySignalName(sigID EntityID, newName string) error {
	if ms.hasParent() {
		parent := ms.Parent()

	loop:
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case SignalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case SignalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				if err := msgParent.modifySignalName(sigID, newName); err != nil {
					return err
				}
				break loop
			}
		}
	}

	sig, err := ms.muxSignals.getValue(sigID)
	if err != nil {
		return err
	}

	oldName := sig.Name()

	ms.removeMuxSignalName(oldName)
	ms.addMuxSignalName(sigID, newName)

	return nil
}

func (ms *MultiplexerSignal) verifySignalName(sigID EntityID, name string) error {
	if err := ms.muxSignalNames.verifyKey(name); err != nil {
		return err
	}

	if ms.hasParent() {
		return ms.parent.verifySignalName(sigID, name)
	}

	return nil
}

func (ms *MultiplexerSignal) verifySignalSizeAmount(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	sig, err := ms.muxSignals.getValue(sigID)
	if err != nil {
		return err
	}

	selVal, err := ms.getMuxSignalSelValue(sigID)
	if err != nil {
		return err
	}

	payload, _ := ms.getSignalPayload(selVal)

	if amount > 0 {
		return payload.verifyBeforeGrow(sig, amount)
	}

	return payload.verifyBeforeShrink(sig, -amount)
}

func (ms *MultiplexerSignal) modifySignalSize(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	sig, err := ms.muxSignals.getValue(sigID)
	if err != nil {
		return err
	}

	selVal, err := ms.getMuxSignalSelValue(sigID)
	if err != nil {
		return err
	}

	payload, _ := ms.getSignalPayload(selVal)

	if amount > 0 {
		return payload.modifyStartBitsOnGrow(sig, amount)
	}

	return payload.modifyStartBitsOnShrink(sig, -amount)
}

// ToParentMessage always returns an error, since [MultiplexerSignal] cannot be converted to [Message].
func (ms *MultiplexerSignal) ToParentMessage() (*Message, error) {
	return nil, fmt.Errorf(`cannot convert to "%s" signal parent is of kind "%s"`,
		SignalParentKindMessage, SignalParentKindMultiplexerSignal)
}

// ToParentMultiplexerSignal returns the [MultiplexerSignal] itself.
func (ms *MultiplexerSignal) ToParentMultiplexerSignal() (*MultiplexerSignal, error) {
	return ms, nil
}

// GetSize returns the total size of the [MultiplexerSignal].
func (ms *MultiplexerSignal) GetSize() int {
	return ms.totalSize
}

// ToStandard always returns an error, since [MultiplexerSignal] cannot be converted to [StandardSignal].
func (ms *MultiplexerSignal) ToStandard() (*StandardSignal, error) {
	return nil, ms.errorf(fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindStandard, SignalKindMultiplexer))
}

// ToEnum always returns an error, since [MultiplexerSignal] cannot be converted to [EnumSignal].
func (ms *MultiplexerSignal) ToEnum() (*EnumSignal, error) {
	return nil, ms.errorf(fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindEnum, SignalKindMultiplexer))
}

// ToMultiplexer always returns the [MultiplexerSignal] itself.
func (ms *MultiplexerSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return ms, nil
}

func (ms *MultiplexerSignal) stringify(b *strings.Builder, tabs int) {
	ms.signal.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("size: %d\n", ms.GetSize()))

	if ms.muxSignals.size() == 0 {
		return
	}

	for selVal, muxGroup := range ms.MuxSignals() {
		b.WriteString(fmt.Sprintf("%sselect value: %d\n", tabStr, selVal))

		b.WriteString(fmt.Sprintf("%smultiplexed signals:\n", tabStr))
		for _, muxSig := range muxGroup {
			muxSig.stringify(b, tabs+1)
			b.WriteRune('\n')
		}
	}
}

func (ms *MultiplexerSignal) String() string {
	builder := new(strings.Builder)
	ms.stringify(builder, 0)
	return builder.String()
}

// GetSelectedMuxSignals returns a slice of signals which belong to the selected group.
func (ms *MultiplexerSignal) GetSelectedMuxSignals(selectValue int) []Signal {
	payload, _ := ms.getSignalPayload(selectValue)

	if payload != nil {
		return payload.signals
	}

	return []Signal{}
}

// MuxSignals returns a map of signal slices, with key the selector value and
// the corresponding value is a slice of signals which belong to the selected group.
// Keep in mind that the keys in the map are not sorted, so it is not guaranteed
// that the first key in the map will corresponde to the smaller select value.
func (ms *MultiplexerSignal) MuxSignals() map[int][]Signal {
	res := make(map[int][]Signal)

	for selVal, payload := range ms.signalPayloads {
		res[selVal] = payload.signals
	}

	return res
}

// AppendMuxSignal appends the [Signal] to the group specified by the select value.
// It may return an error if the signal name is already used by the [MultiplexerSignal]
// or by the [Message] that owns the [MultiplexerSignal]. Also, it may return an error
// if the select value is out of bounds, or if the signal cannot fit in the group.
func (ms *MultiplexerSignal) AppendMuxSignal(selectValue int, signal Signal) error {
	if err := ms.verifySignalName(signal.EntityID(), signal.Name()); err != nil {
		return ms.errorf(err)
	}

	if err := ms.verifySelectValue(selectValue); err != nil {
		return ms.errorf(err)
	}

	payload, realSelVal := ms.getSignalPayload(selectValue)
	if payload == nil {
		payload = ms.addSignalPayload(realSelVal)
	}

	if err := payload.append(signal); err != nil {
		return ms.errorf(err)
	}

	ms.addMuxSignal(realSelVal, signal)

	signal.setParent(ms)

	return nil
}

// InsertMuxSignal inserts the [Signal] to the group specified by the select value starting
// from the specified bit.
// It may return an error if the signal name is already used by the [MultiplexerSignal]
// or by the [Message] that owns the [MultiplexerSignal]. Also, it may return an error
// if the select value is out of bounds, or if the signal cannot fit in the group.
func (ms *MultiplexerSignal) InsertMuxSignal(selectValue int, signal Signal, startBit int) error {
	if err := ms.verifySignalName(signal.EntityID(), signal.Name()); err != nil {
		return ms.errorf(err)
	}

	if err := ms.verifySelectValue(selectValue); err != nil {
		return ms.errorf(err)
	}

	payload, realSelVal := ms.getSignalPayload(selectValue)
	if payload == nil {
		payload = ms.addSignalPayload(realSelVal)
	}

	if err := payload.insert(signal, startBit); err != nil {
		return ms.errorf(err)
	}

	ms.addMuxSignal(realSelVal, signal)

	signal.setParent(ms)

	return nil
}

// ShiftMuxSignalLeft shifts the multiplexed signal with the given entity id left by the given amount.
// It returns the amount of bits shifted.
func (ms *MultiplexerSignal) ShiftMuxSignalLeft(muxSignalEntityID EntityID, amount int) int {
	selVal, err := ms.getMuxSignalSelValue(muxSignalEntityID)
	if err != nil {
		return 0
	}

	sig, err := ms.muxSignals.getValue(muxSignalEntityID)
	if err != nil {
		return 0
	}

	payload, _ := ms.getSignalPayload(selVal)
	if payload == nil {
		return 0
	}

	return payload.shiftLeft(sig, amount)
}

// ShiftMuxSignalRight shifts the multiplexed signal with the given entity id right by the given amount.
// It returns the amount of bits shifted.
func (ms *MultiplexerSignal) ShiftMuxSignalRight(muxSignalEntityID EntityID, amount int) int {
	selVal, err := ms.getMuxSignalSelValue(muxSignalEntityID)
	if err != nil {
		return 0
	}

	sig, err := ms.muxSignals.getValue(muxSignalEntityID)
	if err != nil {
		return 0
	}

	payload, _ := ms.getSignalPayload(selVal)
	if payload == nil {
		return 0
	}

	return payload.shiftRight(sig, amount)
}

// RemoveMuxSignal removes the multiplexed signal with the given entity id from the [MultiplexerSignal].
// It may return an error if the multiplied signal with the given entity id
// is not found in the [MultiplexerSignal].
func (ms *MultiplexerSignal) RemoveMuxSignal(muxSignalEntityID EntityID) error {
	selVal, err := ms.getMuxSignalSelValue(muxSignalEntityID)
	if err != nil {
		return ms.errorf(fmt.Errorf(`cannot remove mux signal with id "%s" : %w`, muxSignalEntityID, err))
	}

	sig, err := ms.muxSignals.getValue(muxSignalEntityID)
	if err != nil {
		return ms.errorf(fmt.Errorf(`cannot remove mux signal with id "%s" : %w`, muxSignalEntityID, err))
	}

	payload, _ := ms.getSignalPayload(selVal)
	if payload == nil {
		return nil
	}

	ms.removeMuxSignal(muxSignalEntityID)
	ms.removeMuxSignalName(sig.Name())

	payload.remove(muxSignalEntityID)

	sig.setParent(nil)

	return nil
}

// RemoveAllMuxSignals removes all the multiplexed signals from the [MultiplexerSignal].
func (ms *MultiplexerSignal) RemoveAllMuxSignals() {
	for muxSigID, tmpMuxSig := range ms.muxSignals.entries() {
		tmpMuxSig.setParent(nil)
		ms.removeMuxSignalName(tmpMuxSig.Name())
		ms.removeMuxSignal(muxSigID)
	}

	for _, payload := range ms.signalPayloads {
		payload.removeAll()
	}
}

// SelectSize returns the number of bits of the select value in the [MultiplexerSignal].
func (ms *MultiplexerSignal) SelectSize() int {
	return ms.selectSize
}

// AddSelectValueRange adds a range of select values to the [MultiplexerSignal].
// It is used when a range of select values is used for selecting one group.
// Ex. from = 0, to = 2 means that there is only one group for select value 0, 1 and 2.
// It may return an error if from is greater then to, or if any of the values in the range
// is already used for selecting more then one group (ex. selVal = 0 -> group0,
// selVal = 1 -> group1: cannot use the range from 0 to 1).
func (ms *MultiplexerSignal) AddSelectValueRange(from, to int) error {
	if from > to {
		return ms.errorf(fmt.Errorf(`cannot set select value range because from "%d" is greater then to "%d"`, from, to))
	}

	if err := ms.verifySelectValue(from); err != nil {
		return ms.errorf(fmt.Errorf(`cannot set select value range : %w`, err))
	}

	if err := ms.verifySelectValue(to); err != nil {
		return ms.errorf(fmt.Errorf(`cannot set select value range : %w`, err))
	}

	foundOne := false
	foundSelVal := from
	for i := from; i <= to; i++ {
		if _, ok := ms.signalPayloads[i]; ok {
			if foundOne {
				return ms.errorf(fmt.Errorf(`cannot set select value range because there are more than 1 payloads between "%d" an "%d"`, from, to))
			}

			foundSelVal = i
			foundOne = true
		}

		if _, ok := ms.selValRanges[i]; ok {
			return ms.errorf(fmt.Errorf(`cannot set select value range because value "%d" is already used in another range`, i))
		}
	}

	for i := from; i <= to; i++ {
		ms.selValRanges[i] = foundSelVal
	}

	return nil
}

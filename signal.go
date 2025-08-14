package acmelib

import (
	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// SignalKind rappresents the kind of a [Signal].
// It can be standard, enum, or multiplexer
type SignalKind int

const (
	// SignalKindStandard defines a standard signal.
	SignalKindStandard SignalKind = iota
	// SignalKindEnum defines a enum signal.
	SignalKindEnum
	// SignalKindMuxor defines a muxor signal.
	SignalKindMuxor
)

func (sk SignalKind) String() string {
	switch sk {
	case SignalKindStandard:
		return "standard"
	case SignalKindEnum:
		return "enum"
	case SignalKindMuxor:
		return "muxor"

	default:
		return "unknown"
	}
}

// SignalSendType rappresents the send type of a [Signal].
type SignalSendType int

const (
	// SignalSendTypeUnset defines an unset transmission type.
	SignalSendTypeUnset SignalSendType = iota
	// SignalSendTypeCyclic defines a cyclic transmission type.
	SignalSendTypeCyclic
	// SignalSendTypeOnWrite defines an on write transmission type.
	SignalSendTypeOnWrite
	// SignalSendTypeOnWriteWithRepetition defines an on write with repetition transmission type.
	SignalSendTypeOnWriteWithRepetition
	// SignalSendTypeOnChange defines an on change transmission type.
	SignalSendTypeOnChange
	// SignalSendTypeOnChangeWithRepetition defines an on change with repetition transmission type.
	SignalSendTypeOnChangeWithRepetition
	// SignalSendTypeIfActive defines an if active transmission type.
	SignalSendTypeIfActive
	// SignalSendTypeIfActiveWithRepetition defines an if active with repetition transmission type.
	SignalSendTypeIfActiveWithRepetition
)

func (sst SignalSendType) String() string {
	switch sst {
	case SignalSendTypeUnset:
		return "unset"
	case SignalSendTypeCyclic:
		return "cyclic"
	case SignalSendTypeOnWrite:
		return "on-write"
	case SignalSendTypeOnWriteWithRepetition:
		return "on-write_with_repetition"
	case SignalSendTypeOnChange:
		return "on-change"
	case SignalSendTypeOnChangeWithRepetition:
		return "on-change_with_repetition"
	case SignalSendTypeIfActive:
		return "if_active"
	case SignalSendTypeIfActiveWithRepetition:
		return "if_active_with_repetition"
	default:
		return "unknown"
	}
}

// Endianness rappresents the byte order of a signal.
// By default a [Endianness] of [EndiannessLittleEndian] is used.
type Endianness int

const (
	// EndiannessLittleEndian defines a little endian byte order.
	EndiannessLittleEndian Endianness = iota
	// EndiannessBigEndian defines a big endian byte order.
	EndiannessBigEndian
)

func (e Endianness) String() string {
	switch e {
	case EndiannessLittleEndian:
		return "little-endian"
	case EndiannessBigEndian:
		return "big-endian"
	default:
		return "unknown"
	}
}

// Signal interface specifies all common methods of
// [StandardSignal], [EnumSignal], and [MultiplexerSignal1].
type Signal interface {
	Entity

	errorf(err error) error

	// UpdateName updates the name of the signal.
	//
	// It returns [NameError] if the new name is not valid.
	UpdateName(name string) error

	// SetDesc stes the description of the signal.
	SetDesc(desc string)

	// AssignAttribute assigns the given attribute/value pair to the signal.
	AssignAttribute(attribute Attribute, value any) error
	// RemoveAttributeAssignment removes the attribute assignment
	// with the given attribute entity id from the Signal.
	RemoveAttributeAssignment(attributeEntityID EntityID) error
	// RemoveAllAttributeAssignments removes all the attribute assignments from the signal.
	RemoveAllAttributeAssignments()
	// AttributeAssignments returns a slice of all attribute assignments of the signal.
	AttributeAssignments() []*AttributeAssignment
	// GetAttributeAssignment returns the attribute assignment
	// with the given attribute entity id from the signal.
	GetAttributeAssignment(attributeEntityID EntityID) (*AttributeAssignment, error)

	stringify(s *stringer.Stringer)
	String() string

	// Kind returns the kind of the signal.
	Kind() SignalKind

	// ParentMessage returns the parent message of the signal.
	// If the signal is part of a multiplexed layer, it will traverse back
	// the layout tree until it finds the parent message.
	// If the signal is standalone (not related to any message), it will return nil.
	ParentMessage() *Message
	setParentMsg(parentMsg *Message)

	// ParentMuxLayer returns the parent multiplexed layer of the signal.
	// If the signal is not part of a multiplexed layer, it will return nil.
	ParentMuxLayer() *MultiplexedLayer
	setparentMuxLayer(ml *MultiplexedLayer)
	setLayout(layout *SignalLayout)

	// StartPos returns the start postion of the signal.
	StartPos() int
	setStartPos(startPos int)
	// UpdateStartPos updates the start position of the signal.
	UpdateStartPos(startPos int) error

	// Size returns the size of the signal.
	Size() int
	setSize(size int)
	verifyAndUpdateSize(instance Signal, newSize int) error

	// SetStartValue sets the initial raw value of the signal.
	SetStartValue(startValue float64)
	// StartValue returns the initial raw value of the signal.
	StartValue() float64
	// SetSendType sets the send type of the signal.
	SetSendType(sendType SignalSendType)
	// SendType returns the send type of the signal.
	SendType() SignalSendType

	// Endianness returns the endianness of the signal.
	Endianness() Endianness
	// SetEndianness sets the endianness of the signal.
	SetEndianness(endianness Endianness)

	// EncodedValue returns the encoded value (raw value) of the signal.
	EncodedValue() uint64

	// ToStandard returns the signal as a standard signal.
	ToStandard() (*StandardSignal, error)
	// ToEnum returns the signal as a enum signal.
	ToEnum() (*EnumSignal, error)
	// ToMuxor returns the signal as a muxor signal.
	ToMuxor() (*MuxorSignal, error)

	// GetLow is used for the ibst
	GetLow() int
	// SetLow is used for the ibst
	SetLow(low int)
	// GetHigh is used for the ibst
	GetHigh() int
	// SetHigh is used for the ibst
	SetHigh(high int)
}

type signal struct {
	*entity
	*withAttributes

	parentMsg      *Message
	parentMuxLayer *MultiplexedLayer

	layout *SignalLayout

	kind SignalKind

	startValue float64
	sendType   SignalSendType

	endianness Endianness

	startPos int
	size     int

	encodedValue uint64
}

func newSignalFromEntity(ent *entity, kind SignalKind) *signal {
	return &signal{
		entity:         ent,
		withAttributes: newWithAttributes(),

		parentMsg:      nil,
		parentMuxLayer: nil,

		layout: nil,

		kind: kind,

		startValue: 0,
		sendType:   SignalSendTypeUnset,

		size:     0,
		startPos: 0,
	}
}

func newSignal(name string, kind SignalKind) *signal {
	return newSignalFromEntity(newEntity(name, EntityKindSignal), kind)
}

func (s *signal) hasParentMsg() bool {
	return s.parentMsg != nil
}

func (s *signal) errorf(err error) error {
	sigErr := &EntityError{
		Kind:     EntityKindSignal,
		EntityID: s.entityID,
		Name:     s.name,
		Err:      err,
	}

	if s.hasParentMsg() {
		return s.parentMsg.errorf(sigErr)
	}

	return sigErr
}

func (s *signal) stringify(str *stringer.Stringer) {
	s.entity.stringify(str)

	str.Write("kind: %s\n", s.kind)
	str.Write("start_pos: %d\n", s.startPos)
	str.Write("size: %d\n", s.size)

	if s.sendType != SignalSendTypeUnset {
		str.Write("send_type: %q\n", s.sendType)
	}

	s.withAttributes.stringify(str)
}

func (s *signal) String() string {
	str := stringer.New()
	str.Write("signal:\n")
	s.stringify(str)
	return str.String()
}

func (s *signal) Kind() SignalKind {
	return s.kind
}

func (s *signal) ParentMessage() *Message {
	if s.parentMsg != nil {
		return s.parentMsg
	}

	currMuxLayer := s.parentMuxLayer

	// Check if the signal is standalone
	if currMuxLayer == nil {
		return nil
	}

	// Traverse back the multiplexed layers to find the parent message
	for currMuxLayer != nil {
		currLayout := currMuxLayer.attachedLayout
		if currLayout == nil {
			return nil
		}

		if currLayout.parentMsg != nil {
			return currLayout.parentMsg
		}

		currMuxLayer = currLayout.parentMuxLayer
	}

	return nil
}

func (s *signal) setParentMsg(parentMsg *Message) {
	s.parentMsg = parentMsg
}

func (s *signal) StartPos() int {
	return s.startPos
}

func (s *signal) setStartPos(startPos int) {
	s.startPos = startPos
}

func (s *signal) setLayout(sl *SignalLayout) {
	s.layout = sl
}

func (s *signal) SetStartValue(startValue float64) {
	s.startValue = startValue
}

func (s *signal) StartValue() float64 {
	return s.startValue
}

func (s *signal) SetSendType(sendType SignalSendType) {
	s.sendType = sendType
}

func (s *signal) SendType() SignalSendType {
	return s.sendType
}

func (s *signal) SetEndianness(endianness Endianness) {
	s.endianness = endianness

	if s.layout != nil {
		s.layout.genFilters()
	}
}

func (s *signal) Endianness() Endianness {
	return s.endianness
}

// UpdateName updates the name of the signal.
//
// It returns a [NameError] if the new name is not valid.
func (s *signal) UpdateName(newName string) error {
	if s.name == newName {
		return nil
	}

	// Check if the signal is standalone
	if s.layout == nil {
		s.name = newName
		return nil
	}

	var sigNamesMap *collection.Map[string, EntityID]

	if s.kind == SignalKindMuxor {
		if err := s.parentMuxLayer.verifySignalName(newName); err != nil {
			return err
		}

		sigNamesMap = s.parentMuxLayer.signalNames
		goto updateName
	}

	// Check if the signal is attached to a message
	if s.layout.parentMsg != nil {
		if err := s.layout.parentMsg.verifySignalName(newName); err != nil {
			return err
		}

		sigNamesMap = s.layout.parentMsg.signalNames
		goto updateName
	}

	// Check if the signal is attached to a multiplexed layer
	if s.layout.parentMuxLayer != nil {
		if err := s.layout.parentMuxLayer.verifySignalName(newName); err != nil {
			return err
		}

		sigNamesMap = s.layout.parentMuxLayer.signalNames
	}

updateName:
	sigNamesMap.Delete(s.name)
	s.name = newName
	sigNamesMap.Set(s.name, s.entityID)

	return nil
}

func (s *signal) Size() int {
	return s.size
}

func (s *signal) setSize(size int) {
	s.size = size
}

func (s *signal) ParentMuxLayer() *MultiplexedLayer {
	return s.parentMuxLayer
}

func (s *signal) setparentMuxLayer(ml *MultiplexedLayer) {
	s.parentMuxLayer = ml
}

func (s *signal) hasParentMuxLayer() bool {
	return s.parentMuxLayer != nil
}

// updateStartPos updates the start position of the signal.
// The instance is required because of golang composition,
// otherwise the layout tree will be set update the item
// to *signal instead of StandardSignal/EnumSignal/MuxorSignal.
func (s *signal) updateStartPos(instance Signal, newStartPos int) error {
	if newStartPos == s.startPos {
		return nil
	}

	if s.kind == SignalKindMuxor {
		if err := s.layout.verifyAndUpdateStartPos(instance, newStartPos); err != nil {
			return s.errorf(err)
		}

		goto setStartPos
	}

	if s.hasParentMuxLayer() {
		// Get all IDs of the layouts that contain the signal
		if layoutIDs, ok := s.parentMuxLayer.singalLayoutIDs.Get(s.entityID); ok {
			for _, lID := range layoutIDs {
				// Check if the new start position is valid,
				// this recursively checks until the base layout is reached (message layout)
				if err := s.parentMuxLayer.layouts[lID].verifyNewStartPos(instance, newStartPos); err != nil {
					return s.errorf(err)
				}
			}

			// The new position is valid, you can update it
			for _, lID := range layoutIDs {
				s.parentMuxLayer.layouts[lID].updateStartPos(instance, newStartPos)
			}
		}

		goto setStartPos
	}

	if s.hasParentMsg() {
		// Check if the new start position is valid and update it
		if err := s.parentMsg.layout.verifyAndUpdateStartPos(instance, newStartPos); err != nil {
			return s.errorf(err)
		}
	}

setStartPos:
	s.setStartPos(newStartPos)
	return nil
}

// verifyNewSize checks if setting the signal to the new size does not intersect with another one.
// The instance is required because of golang composition,
// otherwise the layout tree will be set update the item
// to *signal instead of StandardSignal/EnumSignal/MuxorSignal.
func (s *signal) verifyNewSize(instance Signal, newSize int) error {
	// If the new size is smaller than the original, it cannot be invalid
	if newSize < s.size {
		return nil
	}

	if s.kind == SignalKindMuxor {
		if err := s.layout.verifyNewSize(instance, newSize); err != nil {
			return s.errorf(err)
		}

		return nil
	}

	if s.hasParentMuxLayer() {
		// Get all IDs of the layouts that contain the signal
		if layoutIDs, ok := s.parentMuxLayer.singalLayoutIDs.Get(s.entityID); ok {
			for _, lID := range layoutIDs {
				// Check if the new size is valid,
				// this recursively checks until the base layout is reached (message layout)
				if err := s.parentMuxLayer.layouts[lID].verifyNewSize(instance, newSize); err != nil {
					return s.errorf(err)
				}
			}
		}

		return nil
	}

	if s.hasParentMsg() {
		// Check if the new size is valid and update it
		if err := s.parentMsg.layout.verifyNewSize(instance, newSize); err != nil {
			return s.errorf(err)
		}
	}

	return nil
}

// updateSize updates the size of the signal.
// It does not check if the new size is valid.
// The instance is required because of golang composition,
// otherwise the layout tree will be set update the item
// to *signal instead of StandardSignal/EnumSignal/MuxorSignal.
func (s *signal) updateSize(instance Signal, newSize int) {
	if s.kind == SignalKindMuxor {
		s.layout.updateSize(instance, newSize)
		goto setSize
	}

	if s.hasParentMuxLayer() {
		// Get all IDs of the layouts that contain the signal
		if layoutIDs, ok := s.parentMuxLayer.singalLayoutIDs.Get(s.entityID); ok {
			for _, lID := range layoutIDs {
				s.parentMuxLayer.layouts[lID].updateSize(instance, newSize)
			}
		}

		goto setSize
	}

	if s.hasParentMsg() {
		s.parentMsg.layout.updateSize(instance, newSize)
	}

setSize:
	s.setSize(newSize)
}

// verifyAndUpdateSize checks and updates the size of the signal.
// It is a combination of [verifyNewSize] and [updateSize].
// The instance is required because of golang composition,
// otherwise the layout tree will be set update the item
// to *signal instead of StandardSignal/EnumSignal/MuxorSignal.
func (s *signal) verifyAndUpdateSize(instance Signal, newSize int) error {
	if newSize == s.size {
		return nil
	}

	// If the new size is smaller than the original, it cannot be invalid
	if newSize < s.size {
		goto setSize
	}

	if s.kind == SignalKindMuxor {
		if err := s.layout.verifyAndUpdateSize(instance, newSize); err != nil {
			return s.errorf(err)
		}

		goto setSize
	}

	if s.hasParentMuxLayer() {
		// Get all IDs of the layouts that contain the signal
		if layoutIDs, ok := s.parentMuxLayer.singalLayoutIDs.Get(s.entityID); ok {
			for _, lID := range layoutIDs {
				// Check if the new size is valid,
				// this recursively checks until the base layout is reached (message layout)
				if err := s.parentMuxLayer.layouts[lID].verifyNewSize(instance, newSize); err != nil {
					return s.errorf(err)
				}
			}

			// The new size is valid, you can update it
			for _, lID := range layoutIDs {
				s.parentMuxLayer.layouts[lID].updateSize(instance, newSize)
			}
		}

		goto setSize
	}

	if s.hasParentMsg() {
		// Check if the new size is valid and update it
		if err := s.parentMsg.layout.verifyAndUpdateSize(instance, newSize); err != nil {
			return s.errorf(err)
		}
	}

setSize:
	s.setSize(newSize)
	return nil
}

func (s *signal) EncodedValue() uint64 {
	return s.encodedValue
}

func (s *signal) ToStandard() (*StandardSignal, error) {
	return nil, s.errorf(newConversionError(s.kind.String(), SignalKindStandard.String()))
}

func (s *signal) ToEnum() (*EnumSignal, error) {
	return nil, s.errorf(newConversionError(s.kind.String(), SignalKindEnum.String()))
}

func (s *signal) ToMuxor() (*MuxorSignal, error) {
	return nil, s.errorf(newConversionError(s.kind.String(), SignalKindMuxor.String()))
}

func (s *signal) RemoveAttributeAssignment(attributeEntityID EntityID) error {
	if err := s.removeAttributeAssignment(attributeEntityID); err != nil {
		return s.errorf(err)
	}
	return nil
}

func (s *signal) GetAttributeAssignment(attributeEntityID EntityID) (*AttributeAssignment, error) {
	attAss, err := s.getAttributeAssignment(attributeEntityID)
	if err != nil {
		return nil, s.errorf(err)
	}
	return attAss, nil
}

func (s *signal) assignAttribute(instance Signal, att Attribute, value any) error {
	if err := s.addAttributeAssignment(att, instance, value); err != nil {
		return s.errorf(err)
	}
	return nil
}

func (s *signal) GetLow() int {
	return s.StartPos()
}

func (s *signal) SetLow(low int) {
	s.setStartPos(low)
}

func (s *signal) SetHigh(high int) {
	s.setSize(high - s.GetLow() + 1)
}

func (s *signal) GetHigh() int {
	return s.GetLow() + s.size - 1
}

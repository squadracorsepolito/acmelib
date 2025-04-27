package acmelib

import (
	"fmt"
	"strings"

	"github.com/squadracorsepolito/acmelib/internal/collection"
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

	stringify(b *strings.Builder, tabs int)
	String() string

	// Kind returns the kind of the signal.
	Kind() SignalKind

	// ParentMessage returns the parent message of the signal or nil if not set.
	ParentMessage() *Message
	setParentMsg(parentMsg *Message)

	setparentMuxLayer(ml *MultiplexedLayer)
	setLayout(layout *SL)

	// GetStartPos returns the start postion of the signal.
	GetStartPos() int
	setStartPos(startPos int)
	// UpdateStartPos updates the start position of the signal.
	UpdateStartPos(startPos int) error

	// GetSize returns the size of the signal.
	GetSize() int
	setSize(size int)
	updateSize(newSize int) error

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
	SetEndianness(endianness Endianness)

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

var _ Signal = (*signal)(nil)

type signal struct {
	*entity
	*withAttributes

	parentMsg      *Message
	parentMuxLayer *MultiplexedLayer

	layout *SL

	kind SignalKind

	startValue float64
	sendType   SignalSendType

	endianness Endianness

	relStartPos int
	size        int
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

		size:        0,
		relStartPos: 0,
	}
}

func newSignal(name string, kind SignalKind) *signal {
	return newSignalFromEntity(newEntity(name, EntityKindSignal), kind)
}

func (s *signal) hasParentMsg() bool {
	return s.parentMsg != nil
}

func (s *signal) modifySize(amount int) error {
	// if s.hasParentMuxSig() {
	// 	return s.parentMuxSig.modifySignalSize(s.EntityID(), amount)
	// }

	// if s.hasParentMsg() {
	// 	return s.parentMsg.modifySignalSize(s.EntityID(), amount)
	// }

	return nil
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

func (s *signal) Kind() SignalKind {
	return s.kind
}

func (s *signal) ParentMessage() *Message {
	return s.parentMsg
}

func (s *signal) setParentMsg(parentMsg *Message) {
	s.parentMsg = parentMsg

	if parentMsg != nil {
		s.endianness = parentMsg.byteOrder
	}
}

func (s *signal) GetStartPos() int {
	return s.relStartPos
}

func (s *signal) setStartPos(startPos int) {
	s.relStartPos = startPos
}

func (s *signal) stringify(b *strings.Builder, tabs int) {
	s.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%skind: %s\n", tabStr, s.kind))

	if s.sendType != SignalSendTypeUnset {
		b.WriteString(fmt.Sprintf("%ssend_type: %q\n", tabStr, s.sendType))
	}

	b.WriteString(fmt.Sprintf("%sstart_pos: %d; ", tabStr, s.relStartPos))
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
}

func (s *signal) Endianness() Endianness {
	return s.endianness
}

func (s *signal) GetStartBit() int {
	// if s.hasParentMuxSig() {
	// 	return s.parentMuxSig.GetStartBit() + s.parentMuxSig.GetGroupCountSize() + s.relStartPos
	// }
	return s.relStartPos
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

func (s *signal) GetSize() int {
	return s.size
}

func (s *signal) setSize(size int) {
	s.size = size
}

func (s *signal) GetLow() int {
	return s.GetStartPos()
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

func (s *signal) ToStandard() (*StandardSignal, error) {
	return nil, s.errorf(&ConversionError{
		From: s.kind.String(),
		To:   SignalKindStandard.String(),
	})
}

func (s *signal) ToEnum() (*EnumSignal, error) {
	return nil, s.errorf(&ConversionError{
		From: s.kind.String(),
		To:   SignalKindEnum.String(),
	})
}

func (s *signal) ToMuxor() (*MuxorSignal, error) {
	return nil, s.errorf(&ConversionError{
		From: s.kind.String(),
		To:   SignalKindMuxor.String(),
	})
}

func (s *signal) String() string {
	return ""
}

func (s *signal) AssignAttribute(attribute Attribute, value any) error {
	if err := s.addAttributeAssignment(attribute, s, value); err != nil {
		return s.errorf(err)
	}
	return nil
}

func (s *signal) setparentMuxLayer(ml *MultiplexedLayer) {
	s.parentMuxLayer = ml
}

func (s *signal) hasParentMuxLayer() bool {
	return s.parentMuxLayer != nil
}

// UpdateStartPos updates the start position of the signal.
//
// It returns a [StartPosError] if the new start position is invalid.
func (s *signal) UpdateStartPos(newStartPos int) error {
	if s.hasParentMuxLayer() {
		// Get all IDs of the layouts that contain the signal
		if layoutIDs, ok := s.parentMuxLayer.singalLayoutIDs.Get(s.entityID); ok {
			for _, lID := range layoutIDs {
				// Check if the new start position is valid,
				// this recursively checks until the base layout is reached (message layout)
				if err := s.parentMuxLayer.layouts[lID].verifyNewStartPos(s, newStartPos); err != nil {
					return s.errorf(err)
				}
			}

			// The new position is valid, you can update it
			for _, lID := range layoutIDs {
				s.parentMuxLayer.layouts[lID].updateStartPos(s, newStartPos)
			}
		}

		return nil
	}

	if s.hasParentMsg() {
		// Check if the new start position is valid and update it
		if err := s.parentMsg.layout.verifyAndUpdateStartPos(s, newStartPos); err != nil {
			return s.errorf(err)
		}
	}

	// The signal is not attached to anything
	s.setStartPos(newStartPos)

	return nil
}

// updateSize updates the size of the signal.
func (s *signal) updateSize(newSize int) error {
	if s.hasParentMuxLayer() {
		// Get all IDs of the layouts that contain the signal
		if layoutIDs, ok := s.parentMuxLayer.singalLayoutIDs.Get(s.entityID); ok {
			for _, lID := range layoutIDs {
				// Check if the new size is valid,
				// this recursively checks until the base layout is reached (message layout)
				if err := s.parentMuxLayer.layouts[lID].verifyNewSize(s, newSize); err != nil {
					return s.errorf(err)
				}
			}

			// The new size is valid, you can update it
			for _, lID := range layoutIDs {
				s.parentMuxLayer.layouts[lID].updateSize(s, newSize)
			}
		}
	}

	if s.hasParentMsg() {
		// Check if the new size is valid and update it
		if err := s.parentMsg.layout.verifyAndUpdateSize(s, newSize); err != nil {
			return s.errorf(err)
		}
	}

	// The signal is not attached to anything
	s.setSize(newSize)

	return nil
}

func (s *signal) setLayout(sl *SL) {
	s.layout = sl
}

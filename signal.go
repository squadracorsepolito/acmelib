package acmelib

import (
	"fmt"
	"strings"
	"time"
)

// SignalKind rappresents the kind of a [Signal].
// It can be standard, enum, or multiplexer
type SignalKind int

const (
	// SignalKindStandard defines a standard signal.
	SignalKindStandard SignalKind = iota
	// SignalKindEnum defines a enum signal.
	SignalKindEnum
	// SignalKindMultiplexer defines a multiplexer signal.
	SignalKindMultiplexer
)

func (sk SignalKind) String() string {
	switch sk {
	case SignalKindStandard:
		return "standard"
	case SignalKindEnum:
		return "enum"
	case SignalKindMultiplexer:
		return "multiplexer"
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

// Signal interface specifies all common methods of
// [StandardSignal], [EnumSignal], and [MultiplexerSignal1].
type Signal interface {
	errorf(err error) error

	// EntityID returns the entity id of the signal.
	EntityID() EntityID
	// EntityKind returns the entity kind of the signal.
	EntityKind() EntityKind
	// Name returns the name of the signal.
	Name() string

	// UpdateName updates the name of the signal.
	//
	// It returns [NameError] if the new name is not valid.
	UpdateName(name string) error

	// SetDesc stes the description of the signal.
	SetDesc(desc string)
	// Desc returns the description of the signal.
	Desc() string
	// CreateTime returns the creation time of the signal.
	CreateTime() time.Time

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
	// ParentMultiplexerSignal returns the parent multiplexer signal of the signal
	// or nil if not set.
	ParentMultiplexerSignal() *MultiplexerSignal

	setParentMsg(parentMsg *Message)
	setParentMuxSig(parentMuxSig *MultiplexerSignal)

	// SetStartValue sets the initial raw value of the signal.
	SetStartValue(startValue float64)
	// StartValue returns the initial raw value of the signal.
	StartValue() float64
	// SetSendType sets the send type of the signal.
	SetSendType(sendType SignalSendType)
	// SendType returns the send type of the signal.
	SendType() SignalSendType

	// Endianness returns the endianness of the signal.
	Endianness() MessageByteOrder
	setEndianness(endianness MessageByteOrder)

	// GetStartBit returns the start bit of the signal.
	GetStartBit() int

	// ToStandard returns the signal as a standard signal.
	ToStandard() (*StandardSignal, error)
	// ToEnum returns the signal as a enum signal.
	ToEnum() (*EnumSignal, error)
	// ToMultiplexer returns the signal as a multiplexer signal.
	ToMultiplexer() (*MultiplexerSignal, error)

	// GetSize returns the size of the signal.
	GetSize() int
	setSize(size int)

	// GetRelativeStartPos returns the relative start postion of the signal.
	// It is the same as GetStartBit for non-multiplexed signals.
	GetRelativeStartPos() int
	setRelativeStartPos(startPos int)
	resetStartPos()

	GetLow() int
	SetLow(low int)
	GetHigh() int
	SetHigh(high int)
}

type signal struct {
	*entity
	*withAttributes

	parentMsg    *Message
	parentMuxSig *MultiplexerSignal

	kind SignalKind

	startValue float64
	sendType   SignalSendType

	endianness MessageByteOrder

	relStartPos int
	size        int
}

func newSignalFromEntity(ent *entity, kind SignalKind) *signal {
	return &signal{
		entity:         ent,
		withAttributes: newWithAttributes(),

		parentMsg:    nil,
		parentMuxSig: nil,

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

func (s *signal) hasParentMuxSig() bool {
	return s.parentMuxSig != nil
}

func (s *signal) modifySize(amount int) error {
	if s.hasParentMuxSig() {
		return s.parentMuxSig.modifySignalSize(s.EntityID(), amount)
	}

	if s.hasParentMsg() {
		return s.parentMsg.modifySignalSize(s.EntityID(), amount)
	}

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

func (s *signal) ParentMultiplexerSignal() *MultiplexerSignal {
	return s.parentMuxSig
}

func (s *signal) setParentMsg(parentMsg *Message) {
	s.parentMsg = parentMsg

	if parentMsg != nil {
		s.endianness = parentMsg.byteOrder
	}
}

func (s *signal) setParentMuxSig(parentMuxSig *MultiplexerSignal) {
	s.parentMuxSig = parentMuxSig

	if parentMuxSig != nil && parentMuxSig.hasParentMsg() {
		s.endianness = parentMuxSig.parentMsg.byteOrder
	}
}

func (s *signal) GetRelativeStartPos() int {
	return s.relStartPos
}

func (s *signal) setRelativeStartPos(startPos int) {
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

func (s *signal) Endianness() MessageByteOrder {
	return s.endianness
}

func (s *signal) setEndianness(endianness MessageByteOrder) {
	s.endianness = endianness
}

func (s *signal) GetStartBit() int {
	if s.hasParentMuxSig() {
		return s.parentMuxSig.GetStartBit() + s.parentMuxSig.GetGroupCountSize() + s.relStartPos
	}
	return s.relStartPos
}

func (s *signal) UpdateName(newName string) error {
	sigID := s.entityID
	oldName := s.name

	if oldName == newName {
		return nil
	}

	canUpdMuxSig := false
	if s.hasParentMuxSig() {
		if err := s.parentMuxSig.verifySignalName(sigID, newName); err != nil {
			return s.errorf(&UpdateNameError{Err: err})
		}
		canUpdMuxSig = true
	}

	if s.hasParentMsg() {
		if err := s.parentMsg.verifySignalName(newName); err != nil {
			return s.errorf(&UpdateNameError{Err: err})
		}

		s.parentMsg.signalNames.remove(oldName)
		s.parentMsg.signalNames.add(newName, sigID)
	}

	if canUpdMuxSig {
		s.parentMuxSig.signalNames.remove(oldName)
		s.parentMuxSig.signalNames.add(newName, sigID)
	}

	s.name = newName

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
	return s.GetRelativeStartPos()
}

func (s *signal) resetStartPos() {
	s.relStartPos = 0
}

func (s *signal) SetLow(low int) {
	s.setRelativeStartPos(low)
}

func (s *signal) SetHigh(high int) {
	s.size = high - s.GetLow() + 1
}

func (s *signal) GetHigh() int {
	return s.GetLow() + s.size - 1
}

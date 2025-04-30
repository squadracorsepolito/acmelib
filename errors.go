package acmelib

import (
	"errors"
	"fmt"
)

// ErrIsDuplicated is returned when an entity is duplicated.
var ErrIsDuplicated = errors.New("is duplicated")

// ErrNotFound is returned when an entity is not found.
var ErrNotFound = errors.New("not found")

// ErrIsNegative is returned when a value is negative.
var ErrIsNegative = errors.New("is negative")

// ErrOutOfBounds is returned when a value is out of bounds.
var ErrOutOfBounds = errors.New("out of bounds")

// ErrIsZero is returned when a value is zero.
var ErrIsZero = errors.New("is zero")

// ErrIsNil is returned when a value or entity is nil.
var ErrIsNil = errors.New("is nil")

// ErrIntersects is returned when two entities are intersecting.
var ErrIntersects = errors.New("is intersecting")

// ErrInvalidType is returned when an invalid type is used.
var ErrInvalidType = errors.New("invalid type")

// ErrReceiverIsSender is returned when the receiver is the sender.
var ErrReceiverIsSender = errors.New("receiver is sender")

// ErrTooSmall is returned when a value is too small.
var ErrTooSmall = errors.New("too small")

// ErrTooBig is returned when a value is too big.
var ErrTooBig = errors.New("too big")

// ErrNotClear is returned when a value is not clear.
var ErrNotClear = errors.New("not clear")

// ErrInvalidOneof is returned when a oneof field does not match
// a kind/type field.
type ErrInvalidOneof struct {
	KindTypeField string
}

func (e *ErrInvalidOneof) Error() string {
	return fmt.Sprintf("kind/type field must be %q for this oneof field", e.KindTypeField)
}

// ErrMissingOneofField is returned when a oneof field is missing.
type ErrMissingOneofField struct {
	OneofField string
}

func (e *ErrMissingOneofField) Error() string {
	return fmt.Sprintf("oneof field %q is missing", e.OneofField)
}

// ErrIsRequired is returned when something is required.
// The Thing field is what is required.
type ErrIsRequired struct {
	Item string
}

func (e *ErrIsRequired) Error() string {
	return fmt.Sprintf("%q is required", e.Item)
}

// GreaterThenError is returned when a value is greater than a target.
// The Target field is the target.
type GreaterThenError struct {
	Target string
}

func newGreaterError(target string) *GreaterThenError {
	return &GreaterThenError{Target: target}
}

func (e *GreaterThenError) Error() string {
	return fmt.Sprintf("is greater then %q", e.Target)
}

// LowerThenError is returned when a value is lower than a target.
// The Target field is the target.
type LowerThenError struct {
	Target string
}

func newLowerError(target string) *LowerThenError {
	return &LowerThenError{Target: target}
}

func (e *LowerThenError) Error() string {
	return fmt.Sprintf("is lower then %q", e.Target)
}

// EntityError is returned when a method of an entity fails.
// The Kind field is the entity kind, the EntityID field is the ID,
// the Name field is the name, and the Err field is the cause.
type EntityError struct {
	Kind     EntityKind
	EntityID EntityID
	Name     string
	Err      error
}

func (e *EntityError) Error() string {
	return fmt.Sprintf("%s error; entity_id:%q, name:%q : %v", e.Kind, e.EntityID.String(), e.Name, e.Err)
}

func (e *EntityError) Unwrap() error { return e.Err }

// AddEntityError is returned when an entity cannot be added.
// The EntityID field is the ID of the entity and the Name field is the name,
// and the Err field is the cause.
type AddEntityError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *AddEntityError) Error() string {
	return fmt.Sprintf("add entity error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *AddEntityError) Unwrap() error { return e.Err }

// AttributeValueError is returned when an attribute value is invalid.
// The Err field contains the cause.
type AttributeValueError struct {
	Err error
}

func (e *AttributeValueError) Error() string {
	return fmt.Sprintf("attribute value error : %v", e.Err)
}

func (e *AttributeValueError) Unwrap() error { return e.Err }

// EntityIDError is returned when an entity id is invalid.
// The EntityID field is the id of the entity and the Err field is the cause.
type EntityIDError struct {
	EntityID EntityID
	Err      error
}

func (e *EntityIDError) Error() string {
	return fmt.Sprintf("entity id error; entity_id:%q : %v", e.EntityID, e.Err)
}

func (e *EntityIDError) Unwrap() error { return e.Err }

// SizeError is returned when a size is invalid.
// The Size field is the size and the Err field is the cause.
type SizeError struct {
	Size int
	Err  error
}

func newSizeError(size int, err error) *SizeError {
	return &SizeError{Size: size, Err: err}
}

func (e *SizeError) Error() string {
	return fmt.Sprintf("size error; size:%d : %v", e.Size, e.Err)
}

func (e *SizeError) Unwrap() error { return e.Err }

// StartPosError is returned when a start position is invalid.
// The StartPos field is the start bit and the Err field is the cause.
type StartPosError struct {
	StartPos int
	Err      error
}

func newStartPosError(startPos int, err error) *StartPosError {
	return &StartPosError{StartPos: startPos, Err: err}
}

func (e *StartPosError) Error() string {
	return fmt.Sprintf("start pos error; start_pos:%d : %v", e.StartPos, e.Err)
}

func (e *StartPosError) Unwrap() error { return e.Err }

// LayoutIDError is returned when a layout id is invalid.
// The LayoutID field is the id of the layout and the Err field is the cause.
type LayoutIDError struct {
	LayoutID int
	Err      error
}

func newLayoutIDError(layoutID int, err error) *LayoutIDError {
	return &LayoutIDError{LayoutID: layoutID, Err: err}
}

func (e *LayoutIDError) Error() string {
	return fmt.Sprintf("layout id error; layout_id:%d : %v", e.LayoutID, e.Err)
}

func (e *LayoutIDError) Unwrap() error { return e.Err }

// NameError is returned when a name is invalid.
// The Name field is the name and the Err field is the cause.
type NameError struct {
	Name string
	Err  error
}

func newNameError(name string, err error) *NameError {
	return &NameError{Name: name, Err: err}
}

func (e *NameError) Error() string {
	return fmt.Sprintf("name error; name:%q : %v", e.Name, e.Err)
}

func (e *NameError) Unwrap() error { return e.Err }

// ArgError is returned when an argument is invalid.
// The Name field is the name of the argument and the Err field is the cause.
type ArgError struct {
	Name string
	Err  error
}

func newArgError(name string, err error) *ArgError {
	return &ArgError{Name: name, Err: err}
}

func (e *ArgError) Error() string {
	return fmt.Sprintf("argument error; name:%q : %v", e.Name, e.Err)
}

func (e *ArgError) Unwrap() error { return e.Err }

// IntersectionError is returned when two signals intersect.
// The With field is the name of the other signal.
type IntersectionError struct {
	With string
}

func newIntersectionError(with string) *IntersectionError {
	return &IntersectionError{With: with}
}

func (e *IntersectionError) Error() string {
	return fmt.Sprintf("intersect with %q", e.With)
}

// ConversionError is returned when a signal cannot be converted.
type ConversionError struct {
	From string
	To   string
}

func newConversionError(from string, to string) *ConversionError {
	return &ConversionError{From: from, To: to}
}

func (e *ConversionError) Error() string {
	return fmt.Sprintf("conversion error; from:%q, to:%q", e.From, e.To)
}

// IndexError is returned when an index is invalid.
// The Index field is the index and the Err field is the cause.
type IndexError struct {
	Index int
	Err   error
}

func newIndexError(index int, err error) *IndexError {
	return &IndexError{Index: index, Err: err}
}

func (e *IndexError) Error() string {
	return fmt.Sprintf("index error; index:%d : %v", e.Index, e.Err)
}

func (e *IndexError) Unwrap() error { return e.Err }

// NodeIDError is returned when a [NodeID] is invalid.
// The NodeID field is the node ID and the Err field is the cause.
type NodeIDError struct {
	NodeID NodeID
	Err    error
}

func newNodeIDError(nodeID NodeID, err error) *NodeIDError {
	return &NodeIDError{NodeID: nodeID, Err: err}
}

func (e *NodeIDError) Error() string {
	return fmt.Sprintf("node id error; node_id:%d : %v", e.NodeID, e.Err)
}

func (e *NodeIDError) Unwrap() error { return e.Err }

// CANIDError is returned when a [CANID] is invalid.
// The CANID field is the CAN-ID and the Err field is the cause.
type CANIDError struct {
	CANID CANID
	Err   error
}

func newCANIDError(canID CANID, err error) *CANIDError {
	return &CANIDError{CANID: canID, Err: err}
}

func (e *CANIDError) Error() string {
	return fmt.Sprintf("can id error; can_id:%d : %v", e.CANID, e.Err)
}

func (e *CANIDError) Unwrap() error { return e.Err }

// MessageIDError is returned when a [MessageCANID] is invalid.
// The MessageID field is the message ID and the Err field is the cause.
type MessageIDError struct {
	MessageID MessageID
	Err       error
}

func newMessageIDError(msgID MessageID, err error) *MessageIDError {
	return &MessageIDError{MessageID: msgID, Err: err}
}

func (e *MessageIDError) Error() string {
	return fmt.Sprintf("message id error; message_id:%d : %v", e.MessageID, e.Err)
}

func (e *MessageIDError) Unwrap() error { return e.Err }

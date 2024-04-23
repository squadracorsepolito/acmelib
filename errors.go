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

// ErrNoSpaceLeft is returned when there is not enough space left.
var ErrNoSpaceLeft = errors.New("not enough space left")

// ErrIntersect is returned when two entities are intersecting.
var ErrIntersect = errors.New("is intersecting")

// ErrInvalidType is returned when an invalid type is used.
var ErrInvalidType = errors.New("invalid type")

// ErrGreaterThen is returned when a value is greater than a target.
// The Target field is the target.
type ErrGreaterThen struct {
	Target string
}

func (e *ErrGreaterThen) Error() string {
	return fmt.Sprintf("is greater then %q", e.Target)
}

// ErrLowerThen is returned when a value is lower than a target.
// The Target field is the target.
type ErrLowerThen struct {
	Target string
}

func (e *ErrLowerThen) Error() string {
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

// GetEntityError is returned when an entity cannot be retrieved.
// The EntityID field is the ID of the entity and the Err field is the cause.
type GetEntityError struct {
	EntityID EntityID
	Err      error
}

func (e *GetEntityError) Error() string {
	return fmt.Sprintf("get entity error; entity_id:%q : %v", e.EntityID, e.Err)
}

func (e *GetEntityError) Unwrap() error { return e.Err }

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

// RemoveEntityError is returned when an entity cannot be removed.
// The EntityID field is the ID of the entity and the Err field is the cause.
type RemoveEntityError struct {
	EntityID EntityID
	Err      error
}

func (e *RemoveEntityError) Error() string {
	return fmt.Sprintf("remove entity error; entity_id:%q : %v", e.EntityID.String(), e.Err)
}

func (e *RemoveEntityError) Unwrap() error { return e.Err }

// ArgumentError is returned when an argument is invalid.
// The Name field is the name of the argument and the Err field is the cause.
type ArgumentError struct {
	Name string
	Err  error
}

func (e *ArgumentError) Error() string {
	return fmt.Sprintf("argument error; name:%q : %v", e.Name, e.Err)
}

func (e *ArgumentError) Unwrap() error { return e.Err }

// NameError is returned when a name is invalid.
// The Name field is the name and the Err field is the cause.
type NameError struct {
	Name string
	Err  error
}

func (e *NameError) Error() string {
	return fmt.Sprintf("name error; name:%q : %v", e.Name, e.Err)
}

func (e *NameError) Unwrap() error { return e.Err }

// UpdateNameError is returned when a name cannot be updated.
type UpdateNameError struct {
	Err error
}

func (e *UpdateNameError) Error() string {
	return fmt.Sprintf("update name error : %v", e.Err)
}

func (e *UpdateNameError) Unwrap() error { return e.Err }

// NodeIDError is returned when a [NodeID] is invalid.
// The NodeID field is the node ID and the Err field is the cause.
type NodeIDError struct {
	NodeID NodeID
	Err    error
}

func (e *NodeIDError) Error() string {
	return fmt.Sprintf("node id error; node_id:%d : %v", e.NodeID, e.Err)
}

func (e *NodeIDError) Unwrap() error { return e.Err }

// MessageIDError is returned when a [MessageID] is invalid.
// The MessageID field is the message ID and the Err field is the cause.
type MessageIDError struct {
	MessageID MessageID
	Err       error
}

func (e *MessageIDError) Error() string {
	return fmt.Sprintf("message id error; message_id:%d : %v", e.MessageID, e.Err)
}

func (e *MessageIDError) Unwrap() error { return e.Err }

// GroupIDError is returned when a group ID is invalid.
// The GroupID field is the group ID and the Err field is the cause.
type GroupIDError struct {
	GroupID int
	Err     error
}

func (e *GroupIDError) Error() string {
	return fmt.Sprintf("group id error; group_id:%d : %v", e.GroupID, e.Err)
}

func (e *GroupIDError) Unwrap() error { return e.Err }

// InsertSignalError is returned when a signal cannot be inserted.
// The EntityID field is the ID of the signal, the Name field is the name,
// the StartBit field is the start bit, and the Err field is the cause.
type InsertSignalError struct {
	EntityID EntityID
	Name     string
	StartBit int
	Err      error
}

func (e *InsertSignalError) Error() string {
	return fmt.Sprintf("insert signal error; entity_id:%q, name:%q, start_bit:%d : %v", e.EntityID.String(), e.Name, e.StartBit, e.Err)
}

func (e *InsertSignalError) Unwrap() error { return e.Err }

// AppendSignalError is returned when a signal cannot be appended.
// The EntityID field is the ID of the signal, the Name field is the name,
// and the Err field is the cause.
type AppendSignalError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *AppendSignalError) Error() string {
	return fmt.Sprintf("append signal error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *AppendSignalError) Unwrap() error { return e.Err }

// ConvertionError is returned when a signal cannot be converted.
type ConvertionError struct {
	From string
	To   string
}

func (e *ConvertionError) Error() string {
	return fmt.Sprintf("convertion error; from:%q, to:%q", e.From, e.To)
}

// SignalSizeError is returned when a signal size is invalid.
// The Size field is the size and the Err field is the cause.
type SignalSizeError struct {
	Size int
	Err  error
}

func (e *SignalSizeError) Error() string {
	return fmt.Sprintf("signal size error; size:%d : %v", e.Size, e.Err)
}

func (e *SignalSizeError) Unwrap() error { return e.Err }

// StartBitError is returned when a start bit is invalid.
// The StartBit field is the start bit and the Err field is the cause.
type StartBitError struct {
	StartBit int
	Err      error
}

func (e *StartBitError) Error() string {
	return fmt.Sprintf("start bit error; start_bit:%d : %v", e.StartBit, e.Err)
}

func (e *StartBitError) Unwrap() error { return e.Err }

// UpdateIndexError is returned when an index cannot be updated.
// The Err field is the cause.
type UpdateIndexError struct {
	Err error
}

func (e *UpdateIndexError) Error() string {
	return fmt.Sprintf("update index value error : %v", e.Err)
}

func (e *UpdateIndexError) Unwrap() error { return e.Err }

// ValueIndexError is returned when a value index is invalid.
// The Index field is the index and the Err field is the cause.
type ValueIndexError struct {
	Index int
	Err   error
}

func (e *ValueIndexError) Error() string {
	return fmt.Sprintf("value index value error; index:%d : %v", e.Index, e.Err)
}

func (e *ValueIndexError) Unwrap() error { return e.Err }

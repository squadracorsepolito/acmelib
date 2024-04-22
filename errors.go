package acmelib

import (
	"errors"
	"fmt"
)

var (
	ErrIsDuplicated = errors.New("is duplicated")
	ErrNotFound     = errors.New("not found")
	ErrIsNegative   = errors.New("is negative")
	ErrOutOfBounds  = errors.New("out of bounds")
	ErrIsZero       = errors.New("is zero")
	ErrIsNil        = errors.New("is nil")
	ErrNoSpaceLeft  = errors.New("not enough space left")
	ErrIntersect    = errors.New("is intersecting")
	ErrInvalidType  = errors.New("invalid type")
)

type ArgumentError struct {
	Name string
	Err  error
}

func (e *ArgumentError) Error() string {
	return fmt.Sprintf("argument error; name:%q : %v", e.Name, e.Err)
}

func (e *ArgumentError) Unwrap() error { return e.Err }

type SignalError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *SignalError) Error() string {
	return fmt.Sprintf("signal error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *SignalError) Unwrap() error { return e.Err }

type MessageError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *MessageError) Error() string {
	return fmt.Sprintf("message error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *MessageError) Unwrap() error { return e.Err }

type NameError struct {
	Name string
	Err  error
}

func (e *NameError) Error() string {
	return fmt.Sprintf("name error; name:%q : %v", e.Name, e.Err)
}

func (e *NameError) Unwrap() error { return e.Err }

type MessageIDError struct {
	MessageID MessageID
	Err       error
}

func (e *MessageIDError) Error() string {
	return fmt.Sprintf("message id error; message_id:%d : %v", e.MessageID, e.Err)
}

func (e *MessageIDError) Unwrap() error { return e.Err }

type NodeIDError struct {
	NodeID NodeID
	Err    error
}

func (e *NodeIDError) Error() string {
	return fmt.Sprintf("node id error; node_id:%d : %v", e.NodeID, e.Err)
}

func (e *NodeIDError) Unwrap() error { return e.Err }

type UpdateNameError struct {
	Err error
}

func (e *UpdateNameError) Error() string {
	return fmt.Sprintf("update name error : %v", e.Err)
}

func (e *UpdateNameError) Unwrap() error { return e.Err }

type GroupIDError struct {
	GroupID int
	Err     error
}

func (e *GroupIDError) Error() string {
	return fmt.Sprintf("group id error; group_id:%d : %v", e.GroupID, e.Err)
}

func (e *GroupIDError) Unwrap() error { return e.Err }

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

type AppendSignalError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *AppendSignalError) Error() string {
	return fmt.Sprintf("append signal error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *AppendSignalError) Unwrap() error { return e.Err }

type RemoveEntityError struct {
	EntityID EntityID
	Err      error
}

func (e *RemoveEntityError) Error() string {
	return fmt.Sprintf("remove entity error; entity_id:%q : %v", e.EntityID.String(), e.Err)
}

func (e *RemoveEntityError) Unwrap() error { return e.Err }

type ClearSignalGroupError struct {
	Err error
}

func (e *ClearSignalGroupError) Error() string {
	return fmt.Sprintf("clear signal group error : %v", e.Err)
}

func (e *ClearSignalGroupError) Unwrap() error { return e.Err }

type ConvertionError struct {
	From string
	To   string
}

func (e *ConvertionError) Error() string {
	return fmt.Sprintf("convertion error; from:%q, to:%q", e.From, e.To)
}

type SignalSizeError struct {
	Size int
	Err  error
}

func (e *SignalSizeError) Error() string {
	return fmt.Sprintf("signal size error; size:%d : %v", e.Size, e.Err)
}

func (e *SignalSizeError) Unwrap() error { return e.Err }

type SignalStartBitError struct {
	StartBit int
	Err      error
}

func (e *SignalStartBitError) Error() string {
	return fmt.Sprintf("signal start bit error; start_bit:%d : %v", e.StartBit, e.Err)
}

func (e *SignalStartBitError) Unwrap() error { return e.Err }

type GetEntityError struct {
	EntityID EntityID
	Err      error
}

func (e *GetEntityError) Error() string {
	return fmt.Sprintf("get entity error; entity_id:%q : %v", e.EntityID, e.Err)
}

func (e *GetEntityError) Unwrap() error { return e.Err }

type NodeError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *NodeError) Error() string {
	return fmt.Sprintf("node error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *NodeError) Unwrap() error { return e.Err }

type BusError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *BusError) Error() string {
	return fmt.Sprintf("bus error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *BusError) Unwrap() error { return e.Err }

type NetworkError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *NetworkError) Unwrap() error { return e.Err }

type AddEntityError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *AddEntityError) Error() string {
	return fmt.Sprintf("add entity error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *AddEntityError) Unwrap() error { return e.Err }

type SignalEnumError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *SignalEnumError) Error() string {
	return fmt.Sprintf("signal enum error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *SignalEnumError) Unwrap() error { return e.Err }

type SignalEnumValueError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *SignalEnumValueError) Error() string {
	return fmt.Sprintf("signal enum value error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *SignalEnumValueError) Unwrap() error { return e.Err }

type UpdateIndexError struct {
	Err error
}

func (e *UpdateIndexError) Error() string {
	return fmt.Sprintf("update index value error : %v", e.Err)
}

func (e *UpdateIndexError) Unwrap() error { return e.Err }

type ValueIndexError struct {
	Index int
	Err   error
}

func (e *ValueIndexError) Error() string {
	return fmt.Sprintf("value index value error; index:%d : %v", e.Index, e.Err)
}

func (e *ValueIndexError) Unwrap() error { return e.Err }

type AttributeError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *AttributeError) Error() string {
	return fmt.Sprintf("attribute error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *AttributeError) Unwrap() error { return e.Err }

type ErrGraterThen struct {
	Target string
}

func (e *ErrGraterThen) Error() string {
	return fmt.Sprintf("is greater then %q", e.Target)
}

type ErrLowerThen struct {
	Target string
}

func (e *ErrLowerThen) Error() string {
	return fmt.Sprintf("is lower then %q", e.Target)
}

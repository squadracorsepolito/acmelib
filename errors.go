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
	return fmt.Sprintf("multiplexer signal error; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *SignalError) Unwrap() error { return e.Err }

type NameError struct {
	Name string
	Err  error
}

func (e *NameError) Error() string {
	return fmt.Sprintf("name error; name:%q : %v", e.Name, e.Err)
}

func (e *NameError) Unwrap() error {
	return e.Err
}

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

type RemoveSignalError struct {
	EntityID EntityID
	Err      error
}

func (e *RemoveSignalError) Error() string {
	return fmt.Sprintf("remove signal error; entity_id:%q : %v", e.EntityID.String(), e.Err)
}

func (e *RemoveSignalError) Unwrap() error { return e.Err }

type ClearSignalGroupError struct {
	Err error
}

func (e *ClearSignalGroupError) Error() string {
	return fmt.Sprintf("clear signal group error : %v", e.Err)
}

func (e *ClearSignalGroupError) Unwrap() error { return e.Err }

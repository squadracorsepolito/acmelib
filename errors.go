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
	return fmt.Sprintf("argument; name:%q : %v", e.Name, e.Err)
}

func (e *ArgumentError) Unwrap() error { return e.Err }

type SignalError struct {
	EntityID EntityID
	Name     string
	Err      error
}

func (e *SignalError) Error() string {
	return fmt.Sprintf("multiplexer signal; entity_id:%q, name:%q : %v", e.EntityID.String(), e.Name, e.Err)
}

func (e *SignalError) Unwrap() error { return e.Err }

type SignalNameError struct {
	Name string
	Err  error
}

func (e *SignalNameError) Error() string {
	return fmt.Sprintf("signal name %q : %v", e.Name, e.Err)
}

func (e *SignalNameError) Unwrap() error {
	return e.Err
}

type GroupIDError struct {
	GroupID int
	Err     error
}

func (e *GroupIDError) Error() string {
	return fmt.Sprintf("group id %d : %v", e.GroupID, e.Err)
}

func (e *GroupIDError) Unwrap() error { return e.Err }

type InsertSignalError struct {
	EntityID EntityID
	Name     string
	StartBit int
	Err      error
}

func (e *InsertSignalError) Error() string {
	return fmt.Sprintf("insert signal; entity_id:%q, name:%q, start_bit:%d : %v", e.EntityID.String(), e.Name, e.StartBit, e.Err)
}

func (e *InsertSignalError) Unwrap() error { return e.Err }

type RemoveSignalError struct {
	EntityID EntityID
	Err      error
}

func (e *RemoveSignalError) Error() string {
	return fmt.Sprintf("remove signal; entity_id:%q : %v", e.EntityID.String(), e.Err)
}

func (e *RemoveSignalError) Unwrap() error { return e.Err }

type ClearSignalGroupError struct {
	Err error
}

func (e *ClearSignalGroupError) Error() string {
	return fmt.Sprintf("clear signal group : %v", e.Err)
}

func (e *ClearSignalGroupError) Unwrap() error { return e.Err }

package acmelib

import (
	"slices"

	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

func newDefaultCANIDBuilder() *CANIDBuilder {
	return NewCANIDBuilder("default_CANID_builder").UseNodeID(0, 4).UseMessageID(4, 7).UseCAN2A()
}

// CANID is the CAN-ID of a [Message] within a [Bus].
// Every message should have a different CAN-ID.
type CANID uint32

// CANIDBuilderOpKind is the kind of an operation
// perfomed by the [CANIDBuilder].
type CANIDBuilderOpKind int

const (
	// CANIDBuilderOpKindMessagePriority represents an operation
	// that involves the message priority.
	CANIDBuilderOpKindMessagePriority CANIDBuilderOpKind = iota
	// CANIDBuilderOpKindMessageID represents an operation
	// that involves the message id.
	CANIDBuilderOpKindMessageID
	// CANIDBuilderOpKindNodeID represents an operation
	// that involves the node id.
	CANIDBuilderOpKindNodeID
	// CANIDBuilderOpKindBitMask represents a bit masking operation.
	CANIDBuilderOpKindBitMask
)

func (bok CANIDBuilderOpKind) String() string {
	switch bok {
	case CANIDBuilderOpKindMessagePriority:
		return "message-priority"
	case CANIDBuilderOpKindMessageID:
		return "message-id"
	case CANIDBuilderOpKindNodeID:
		return "node-id"
	case CANIDBuilderOpKindBitMask:
		return "bit-mask"
	default:
		return "unknown"
	}
}

// CANIDBuilderOp is an operation performed by the [CANIDBuilder].
type CANIDBuilderOp struct {
	kind CANIDBuilderOpKind
	from int
	len  int
}

func newCANIDBuilderOp(kind CANIDBuilderOpKind, from, len int) *CANIDBuilderOp {
	return &CANIDBuilderOp{
		kind: kind,
		from: from,
		len:  len,
	}
}

func (bo *CANIDBuilderOp) stringify(s *stringer.Stringer) {
	s.Write("kind: %s; from: %d; len: %d\n", bo.kind, bo.from, bo.len)
}

// Kind returns the kind of the operation.
func (bo *CANIDBuilderOp) Kind() CANIDBuilderOpKind {
	return bo.kind
}

// From returns the index of the first bit on which the operation is performed.
func (bo *CANIDBuilderOp) From() int {
	return bo.from
}

// Len returns the number of bits on which the operation is performed.
func (bo *CANIDBuilderOp) Len() int {
	return bo.len
}

// CANIDBuilder is a builder used to describe how to generate
// the CAN-ID of the messages within a [Bus].
type CANIDBuilder struct {
	*entity
	*withRefs[*Bus]

	operations []*CANIDBuilderOp
}

func newCANIDBuilderFromEntity(ent *entity) *CANIDBuilder {
	return &CANIDBuilder{
		entity:   ent,
		withRefs: newWithRefs[*Bus](),

		operations: []*CANIDBuilderOp{},
	}
}

// NewCANIDBuilder creates a new [CANIDBuilder] with the given name.
func NewCANIDBuilder(name string) *CANIDBuilder {
	return newCANIDBuilderFromEntity(newEntity(name, EntityKindCANIDBuilder))
}

func (b *CANIDBuilder) stringify(s *stringer.Stringer) {
	b.entity.stringify(s)

	if len(b.operations) > 0 {
		s.Write("operations:\n")
		s.Indent()
		for _, op := range b.operations {
			op.stringify(s)
		}
		s.Unindent()
	}

	b.withRefs.stringify(s)
}

func (b *CANIDBuilder) String() string {
	s := stringer.New()
	s.Write("can_id_builder:\n")
	b.stringify(s)
	return s.String()
}

// UpdateName updates the name of the [CANIDBuilder].
func (b *CANIDBuilder) UpdateName(name string) error {
	b.entity.name = name
	return nil
}

// Operations returns the operations performed by the [CANIDBuilder].
func (b *CANIDBuilder) Operations() []*CANIDBuilderOp {
	return b.operations
}

func (b *CANIDBuilder) calculateOp(op *CANIDBuilderOp, prev CANID, msgPriority MessagePriority, msgID MessageID, nodeID NodeID) CANID {
	canID := uint32(prev)

	if op.kind == CANIDBuilderOpKindBitMask {
		mask := uint32(0xFFFFFFFF) >> uint32(32-op.len)
		canID &= (mask << uint32(op.from))
		return CANID(canID)
	}

	tmpVal := uint32(0)
	switch op.kind {
	case CANIDBuilderOpKindMessagePriority:
		tmpVal = uint32(msgPriority)
	case CANIDBuilderOpKindMessageID:
		tmpVal = uint32(msgID)
	case CANIDBuilderOpKindNodeID:
		tmpVal = uint32(nodeID)
	}

	mask := uint32(0xFFFFFFFF) >> uint32(32-op.len)
	tmpVal &= mask

	tmpVal = tmpVal << uint32(op.from)
	canID |= tmpVal

	return CANID(canID)
}

// Calculate returns the CAN-ID calculated by applying the operations.
func (b *CANIDBuilder) Calculate(messagePriority MessagePriority, messageID MessageID, nodeID NodeID) CANID {
	canID := CANID(0)

	for _, op := range b.operations {
		canID = b.calculateOp(op, canID, messagePriority, messageID, nodeID)
	}

	return canID
}

// CalculatePartials returns the CAN-IDs calculated by applying the operations.
// The last entry is the final CAN-ID.
func (b *CANIDBuilder) CalculatePartials(messagePriority MessagePriority, messageID MessageID, nodeID NodeID) []CANID {
	canIDs := []CANID{}

	prevCANID := CANID(0)
	for _, op := range b.operations {
		prevCANID = b.calculateOp(op, prevCANID, messagePriority, messageID, nodeID)
		canIDs = append(canIDs, prevCANID)
	}

	return canIDs
}

// UseMessagePriority adds an operation that involves the message priority from the given index.
// The length of the operation is fixed (2 bits).
func (b *CANIDBuilder) UseMessagePriority(from int) *CANIDBuilder {
	b.operations = append(b.operations, newCANIDBuilderOp(CANIDBuilderOpKindMessagePriority, from, 2))
	return b
}

// UseMessageID adds an operation that involves the message id from the given index and length.
func (b *CANIDBuilder) UseMessageID(from, len int) *CANIDBuilder {
	b.operations = append(b.operations, newCANIDBuilderOp(CANIDBuilderOpKindMessageID, from, len))
	return b
}

// UseNodeID adds an operation that involves the node id from the given index and length.
func (b *CANIDBuilder) UseNodeID(from, len int) *CANIDBuilder {
	b.operations = append(b.operations, newCANIDBuilderOp(CANIDBuilderOpKindNodeID, from, len))
	return b
}

// UseCAN2A adds a bit mask from 0 with a length of 11,
// which makes the calculated CAN-ID conformed to the CAN 2.0A.
func (b *CANIDBuilder) UseCAN2A() *CANIDBuilder {
	b.operations = append(b.operations, newCANIDBuilderOp(CANIDBuilderOpKindBitMask, 0, 11))
	return b
}

// UseBitMask adds a bit mask operation from the given index and length.
func (b *CANIDBuilder) UseBitMask(from, len int) *CANIDBuilder {
	b.operations = append(b.operations, newCANIDBuilderOp(CANIDBuilderOpKindBitMask, from, len))
	return b
}

// InsertOperation inserts an operation at the given index in the [CANIDBuilder].
//
// It returns an [ArgError] if one of the arguments is out of bounds.
func (b *CANIDBuilder) InsertOperation(kind CANIDBuilderOpKind, from, length, opIndex int) error {
	if from < 0 || from > 31 {
		return &ArgError{
			Name: "from",
			Err:  ErrOutOfBounds,
		}
	}

	if length < 0 || length > 32-from {
		return &ArgError{
			Name: "length",
			Err:  ErrOutOfBounds,
		}
	}

	if opIndex < 0 || opIndex > len(b.operations) {
		return &ArgError{
			Name: "opIndex",
			Err:  ErrOutOfBounds,
		}
	}

	op := newCANIDBuilderOp(kind, from, length)
	b.operations = slices.Insert(b.operations, opIndex, op)

	return nil
}

// RemoveOperation removes the operation at the given index from the [CANIDBuilder].
//
// It returns an [ArgError] if the operation's index is out of bounds.
func (b *CANIDBuilder) RemoveOperation(opIndex int) error {
	if opIndex < 0 || opIndex >= len(b.operations) {
		return &ArgError{
			Name: "opIndex",
			Err:  ErrOutOfBounds,
		}
	}

	b.operations = slices.Delete(b.operations, opIndex, opIndex+1)

	return nil
}

// RemoveAllOperations removes all operations from the [CANIDBuilder].
func (b *CANIDBuilder) RemoveAllOperations() {
	b.operations = []*CANIDBuilderOp{}
}

package acmelib

import (
	"fmt"
	"strings"
)

var defaulCANIDBuilder = NewCANIDBuilder("default CAN-ID builder").UseNodeID(0, 4).UseMessageID(4, 7)

// CANID is the CAN-ID of a [Message] within a [Bus].
// Every message should have a different CAN-ID.
type CANID uint32

// CANIDBuilderOpKind is the kind of an operation
// perfomed by the [CANIDBuilder].
type CANIDBuilderOpKind int

const (
	// CANIDBuilderOpMessagePriority represents an operation
	// that involves the message priority.
	CANIDBuilderOpMessagePriority CANIDBuilderOpKind = iota
	// CANIDBuilderOpMessageID represents an operation
	// that involves the message id.
	CANIDBuilderOpMessageID
	// CANIDBuilderOpNodeID represents an operation
	// that involves the node id.
	CANIDBuilderOpNodeID
)

func (bok CANIDBuilderOpKind) String() string {
	switch bok {
	case CANIDBuilderOpMessagePriority:
		return "message-priority"
	case CANIDBuilderOpMessageID:
		return "message-id"
	case CANIDBuilderOpNodeID:
		return "node-id"
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

func (bo *CANIDBuilderOp) stringify(b *strings.Builder, tabs int) {
	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%skind: %s\n", tabStr, bo.kind))
	b.WriteString(fmt.Sprintf("%sfrom: %d; len: %d\n", tabStr, bo.from, bo.len))
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
	*withTemplateRefs[*Bus]

	operations []*CANIDBuilderOp
}

// NewCANIDBuilder creates a new [CANIDBuilder] with the given name.
func NewCANIDBuilder(name string) *CANIDBuilder {
	return &CANIDBuilder{
		entity:           newEntity(name),
		withTemplateRefs: newWithTemplateRefs[*Bus](),

		operations: []*CANIDBuilderOp{},
	}
}

func (b *CANIDBuilder) stringify(builder *strings.Builder, tabs int) {
	b.entity.stringify(builder, tabs)

	tabStr := getTabString(tabs)

	builder.WriteString(fmt.Sprintf("%soperations:\n", tabStr))
	for _, op := range b.operations {
		op.stringify(builder, tabs+1)
	}
}

func (b *CANIDBuilder) String() string {
	builder := new(strings.Builder)
	b.stringify(builder, 0)
	return builder.String()
}

// Operations returns the operations performed by the [CANIDBuilder].
func (b *CANIDBuilder) Operations() []*CANIDBuilderOp {
	return b.operations
}

// Calculate returns the CAN-ID calculated by applying the operations.
func (b *CANIDBuilder) Calculate(messagePriority MessagePriority, messageID MessageID, nodeID NodeID) CANID {
	canID := uint32(0)

	for _, op := range b.operations {
		tmpVal := uint32(0)
		switch op.kind {
		case CANIDBuilderOpMessagePriority:
			tmpVal = uint32(messagePriority)
		case CANIDBuilderOpMessageID:
			tmpVal = uint32(messageID)
		case CANIDBuilderOpNodeID:
			tmpVal = uint32(nodeID)
		}

		mask := uint32(0xFFFFFFFF) >> uint32(32-op.len)
		tmpVal &= mask

		tmpVal = tmpVal << uint32(op.from)
		canID |= tmpVal
	}

	return CANID(canID)
}

// UseMessagePriority adds an operation that involves the message priority from the given index.
// The length of the operation is fixed (2 bits).
func (b *CANIDBuilder) UseMessagePriority(from int) *CANIDBuilder {
	b.operations = append(b.operations, &CANIDBuilderOp{
		kind: CANIDBuilderOpMessagePriority,
		from: from,
		len:  2,
	})
	return b
}

// UseMessageID adds an operation that involves the message id from the given index and length.
func (b *CANIDBuilder) UseMessageID(from, len int) *CANIDBuilder {
	b.operations = append(b.operations, &CANIDBuilderOp{
		kind: CANIDBuilderOpMessageID,
		from: from,
		len:  len,
	})
	return b
}

// UseNodeID adds an operation that involves the node id from the given index and length.
func (b *CANIDBuilder) UseNodeID(from, len int) *CANIDBuilder {
	b.operations = append(b.operations, &CANIDBuilderOp{
		kind: CANIDBuilderOpNodeID,
		from: from,
		len:  len,
	})
	return b
}

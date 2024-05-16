package acmelib

import (
	"fmt"
	"strings"
)

var defaulCANIDBuilder = NewCANIDBuilder("default CAN-ID builder").UseNodeID(0, 4).UseMessageID(4, 7)

type CANID uint32

type CANIDBuilderOpKind int

const (
	CANIDBuilderOpMessagePriority CANIDBuilderOpKind = iota
	CANIDBuilderOpMessageID
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

func (bo *CANIDBuilderOp) Kind() CANIDBuilderOpKind {
	return bo.kind
}

func (bo *CANIDBuilderOp) From() int {
	return bo.from
}

func (bo *CANIDBuilderOp) Len() int {
	return bo.len
}

type CANIDBuilder struct {
	*entity
	*withTemplateRefs[*Bus]

	operations []*CANIDBuilderOp
}

func NewCANIDBuilder(name string) *CANIDBuilder {
	return &CANIDBuilder{
		entity:           newEntity(name),
		withTemplateRefs: newWithTemplateRefs[*Bus](TemplateKindCANIDBuilder),

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

func (b *CANIDBuilder) Operations() []*CANIDBuilderOp {
	return b.operations
}

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

func (b *CANIDBuilder) UseMessagePriority(from int) *CANIDBuilder {
	b.operations = append(b.operations, &CANIDBuilderOp{
		kind: CANIDBuilderOpMessagePriority,
		from: from,
		len:  2,
	})
	return b
}

func (b *CANIDBuilder) UseMessageID(from, len int) *CANIDBuilder {
	b.operations = append(b.operations, &CANIDBuilderOp{
		kind: CANIDBuilderOpMessageID,
		from: from,
		len:  len,
	})
	return b
}

func (b *CANIDBuilder) UseNodeID(from, len int) *CANIDBuilder {
	b.operations = append(b.operations, &CANIDBuilderOp{
		kind: CANIDBuilderOpNodeID,
		from: from,
		len:  len,
	})
	return b
}

func (b *CANIDBuilder) ToCANIDBuilder() (*CANIDBuilder, error) {
	return b, nil
}

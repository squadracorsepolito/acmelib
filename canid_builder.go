package acmelib

type CANID uint32

type CANIDBuilderOpKind int

const (
	CANIDBuilderOpMessagePriority CANIDBuilderOpKind = iota
	CANIDBuilderOpMessageID
	CANIDBuilderOpNodeID
)

type CANIDBuilderOp struct {
	kind CANIDBuilderOpKind
	from int
	len  int
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
	operations []*CANIDBuilderOp
}

func NewCANIDBuilder() *CANIDBuilder {
	return &CANIDBuilder{
		operations: []*CANIDBuilderOp{},
	}
}

func (b *CANIDBuilder) Operations() []*CANIDBuilderOp {
	return b.operations
}

func (b *CANIDBuilder) Calculate(messagePriority MessagePriority, messageID int, nodeID NodeID) CANID {
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

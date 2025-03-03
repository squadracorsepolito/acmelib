package acmelib

import (
	"fmt"
	"slices"
	"strings"
)

// MessageID rappresents the ID of a [Message].
// It must be unique within all the messages sended by a [NodeInterface].
type MessageID uint32

// MessagePriority rappresents the priority of a [Message].
// The priorities are very high, high, medium, and low.
// The higher priority has the value 0 and the lower has 3.
type MessagePriority uint32

const (
	// MessagePriorityVeryHigh defines a very high priority.
	MessagePriorityVeryHigh MessagePriority = iota
	// MessagePriorityHigh defines an high priority.
	MessagePriorityHigh
	// MessagePriorityMedium defines a medium priority.
	MessagePriorityMedium
	// MessagePriorityLow defines a low priority.
	MessagePriorityLow
)

// MessageSendType rappresents the transition type of a [Message].
type MessageSendType int

const (
	// MessageSendTypeUnset defines an unset transmission type.
	MessageSendTypeUnset MessageSendType = iota
	// MessageSendTypeCyclic defines a cyclic transmission type.
	MessageSendTypeCyclic
	// MessageSendTypeCyclicIfActive defines a cyclic if active transmission type.
	MessageSendTypeCyclicIfActive
	// MessageSendTypeCyclicAndTriggered defines a cyclic and triggered transmission type.
	MessageSendTypeCyclicAndTriggered
	// MessageSendTypeCyclicIfActiveAndTriggered defines a cyclic if active and triggered transmission type.
	MessageSendTypeCyclicIfActiveAndTriggered
)

func (mst MessageSendType) String() string {
	switch mst {
	case MessageSendTypeUnset:
		return "unset"
	case MessageSendTypeCyclic:
		return "cyclic"
	case MessageSendTypeCyclicIfActive:
		return "cyclic_if_active"
	case MessageSendTypeCyclicAndTriggered:
		return "cyclic_and_triggered"
	case MessageSendTypeCyclicIfActiveAndTriggered:
		return "cyclic_if_active_and_triggered"
	default:
		return "unknown"
	}
}

// MessageByteOrder rappresents the byte order of the payload of a [Message].
// By default a [MessageByteOrder] of [MessageByteOrderLittleEndian] is used.
type MessageByteOrder int

const (
	// MessageByteOrderLittleEndian defines a little endian byte order.
	MessageByteOrderLittleEndian MessageByteOrder = iota
	// MessageByteOrderBigEndian defines a big endian byte order.
	MessageByteOrderBigEndian
)

func (mbo MessageByteOrder) String() string {
	switch mbo {
	case MessageByteOrderLittleEndian:
		return "little-endian"
	case MessageByteOrderBigEndian:
		return "big-endian"
	default:
		return "unknown"
	}
}

// Message is the representation of data sent by a node thought the bus.
// It holds a list of signals that are contained in the message payload.
type Message struct {
	*entity
	*withAttributes

	senderNodeInt *NodeInterface

	signals     *set[EntityID, Signal]
	signalNames *set[string, EntityID]

	signalPayload *signalPayload

	sizeByte int

	id             MessageID
	staticCANID    CANID
	hasStaticCANID bool

	priority  MessagePriority
	byteOrder MessageByteOrder

	cycleTime      int
	sendType       MessageSendType
	delayTime      int
	startDelayTime int

	receivers *set[EntityID, *NodeInterface]
}

func newMessageFromEntity(ent *entity, id MessageID, sizeByte int) *Message {
	return &Message{
		entity:         ent,
		withAttributes: newWithAttributes(),

		senderNodeInt: nil,

		signals:     newSet[EntityID, Signal](),
		signalNames: newSet[string, EntityID](),

		signalPayload: newSignalPayload(sizeByte * 8),

		sizeByte: sizeByte,

		id:             id,
		staticCANID:    0,
		hasStaticCANID: false,

		priority:  MessagePriorityVeryHigh,
		byteOrder: MessageByteOrderLittleEndian,

		cycleTime:      0,
		sendType:       MessageSendTypeUnset,
		delayTime:      0,
		startDelayTime: 0,

		receivers: newSet[EntityID, *NodeInterface](),
	}
}

// NewMessage creates a new [Message] with the given name, id and size in bytes.
// By default a [MessagePriority] of [MessagePriorityVeryHigh] is used.
func NewMessage(name string, id MessageID, sizeByte int) *Message {
	return newMessageFromEntity(newEntity(name, EntityKindMessage), id, sizeByte)
}

func (m *Message) hasSenderNodeInt() bool {
	return m.senderNodeInt != nil
}

func (m *Message) errorf(err error) error {
	msgErr := &EntityError{
		Kind:     EntityKindMessage,
		EntityID: m.entityID,
		Name:     m.name,
		Err:      err,
	}
	if m.hasSenderNodeInt() {
		return m.senderNodeInt.errorf(msgErr)
	}
	return msgErr
}

func (m *Message) verifySignalName(name string) error {
	err := m.signalNames.verifyKeyUnique(name)
	if err != nil {
		return &NameError{
			Name: name,
			Err:  err,
		}
	}
	return nil
}

func (m *Message) verifySignalSizeAmount(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	sig, err := m.signals.getValue(sigID)
	if err != nil {
		return err
	}

	if amount > 0 {
		if err := m.signalPayload.verifyBeforeGrow(sig, amount); err != nil {
			return &SignalSizeError{
				Size: sig.GetSize() + amount,
				Err:  err,
			}
		}

		return nil
	}

	if err := m.signalPayload.verifyBeforeShrink(sig, -amount); err != nil {
		return &SignalSizeError{
			Size: sig.GetSize() + amount,
			Err:  err,
		}
	}

	return nil
}

func (m *Message) modifySignalSize(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	sig, err := m.signals.getValue(sigID)
	if err != nil {
		return err
	}

	if amount > 0 {
		return m.signalPayload.modifyStartBitsOnGrow(sig, amount)
	}

	return m.signalPayload.modifyStartBitsOnShrink(sig, -amount)
}

func (m *Message) stringify(b *strings.Builder, tabs int) {
	m.entity.stringify(b, tabs)

	tabStr := getTabString(tabs)

	if m.id != 0 {
		b.WriteString(fmt.Sprintf("%smessage_id: %d\n", tabStr, m.id))
	}

	b.WriteString(fmt.Sprintf("%spriority: %d (very_high=0; low=3)\n", tabStr, m.priority))
	b.WriteString(fmt.Sprintf("%ssize: %d bytes\n", tabStr, m.sizeByte))

	if m.cycleTime != 0 {
		b.WriteString(fmt.Sprintf("%scycle_time: %d ms\n", tabStr, m.cycleTime))
	}

	if m.delayTime != 0 {
		b.WriteString(fmt.Sprintf("%sdelay_time: %d ms\n", tabStr, m.delayTime))
	}

	if m.startDelayTime != 0 {
		b.WriteString(fmt.Sprintf("%sstart_delay_time: %d ms\n", tabStr, m.startDelayTime))
	}

	if m.sendType != MessageSendTypeUnset {
		b.WriteString(fmt.Sprintf("%ssend_type: %q\n", tabStr, m.sendType))
	}

	if m.receivers.size() > 0 {
		b.WriteString(fmt.Sprintf("%sreceivers:\n", tabStr))
		for _, rec := range m.Receivers() {
			b.WriteString(fmt.Sprintf("%s\tname: %s; node_id: %d; entity_id: %s\n", tabStr, rec.node.name, rec.node.id, rec.node.entityID))
		}
	}

	if m.signals.size() == 0 {
		return
	}

	b.WriteString(fmt.Sprintf("%ssignals:\n", tabStr))
	for _, sig := range m.Signals() {
		sig.stringify(b, tabs+1)
		b.WriteRune('\n')
	}
}

func (m *Message) addSignal(sig Signal) {
	sigID := sig.EntityID()

	m.signals.add(sigID, sig)
	m.signalNames.add(sig.Name(), sigID)

	sig.setParentMsg(m)

	if sig.Kind() != SignalKindMultiplexer {
		return
	}

	muxSigStack := newStack[Signal]()
	muxSigStack.push(sig)

	for muxSigStack.size() > 0 {
		currSig := muxSigStack.pop()

		muxSig, err := currSig.ToMultiplexer()
		if err != nil {
			panic(err)
		}

		for tmpSigID, tmpSig := range muxSig.signals.entries() {
			if tmpSig.Kind() == SignalKindMultiplexer {
				muxSigStack.push(tmpSig)
			}

			m.signals.add(tmpSigID, tmpSig)
			tmpSig.setParentMsg(m)
		}

		for tmpName, tmpSigID := range muxSig.signalNames.entries() {
			m.signalNames.add(tmpName, tmpSigID)
		}
	}
}

func (m *Message) removeSignal(sig Signal) {
	sigID := sig.EntityID()

	m.signals.remove(sigID)
	m.signalNames.remove(sig.Name())

	sig.setParentMsg(nil)

	if sig.Kind() != SignalKindMultiplexer {
		return
	}

	muxSigStack := newStack[Signal]()
	muxSigStack.push(sig)

	for muxSigStack.size() > 0 {
		currSig := muxSigStack.pop()

		muxSig, err := currSig.ToMultiplexer()
		if err != nil {
			panic(err)
		}

		for tmpSigID, tmpSig := range muxSig.signals.entries() {
			if tmpSig.Kind() == SignalKindMultiplexer {
				muxSigStack.push(tmpSig)
			}

			m.signals.remove(tmpSigID)
			tmpSig.setParentMsg(nil)
		}

		for _, tmpName := range muxSig.signalNames.getKeys() {
			m.signalNames.remove(tmpName)
		}
	}
}

func (m *Message) String() string {
	builder := new(strings.Builder)
	m.stringify(builder, 0)
	return builder.String()
}

// UpdateName updates the name of the [Message].
// It may return an error if the new name is already used within a node.
func (m *Message) UpdateName(newName string) error {
	if m.name == newName {
		return nil
	}

	if m.hasSenderNodeInt() {
		if err := m.senderNodeInt.sentMessageNames.verifyKeyUnique(newName); err != nil {
			return m.errorf(&UpdateNameError{
				Err: &NameError{
					Name: newName,
					Err:  err,
				},
			})
		}

		m.senderNodeInt.sentMessageNames.modifyKey(m.name, newName, m.entityID)
	}

	m.name = newName

	return nil
}

// UpdateSizeByte updates the size of the [Message] to the given value in bytes.
//
// It returns:
//   - [ArgumentError] if the new size is negative.
//   - [MessageSizeError] if the new size is invalid.
func (m *Message) UpdateSizeByte(newSizeByte int) error {
	if newSizeByte < 0 {
		return m.errorf(
			&MessageSizeError{
				Size: newSizeByte,
				Err:  ErrIsNegative,
			},
		)
	}

	if m.sizeByte == newSizeByte {
		return nil
	}

	if m.hasSenderNodeInt() {
		if err := m.senderNodeInt.verifyMessageSize(newSizeByte); err != nil {
			return err
		}
	}

	if err := m.signalPayload.resize(newSizeByte * 8); err != nil {
		return m.errorf(err)
	}

	m.sizeByte = newSizeByte

	return nil
}

// SenderNodeInterface returns the [NodeInterface] that is responsible for sending the [Message].
// If the [Message] is not sent by a [NodeInterface], it will return nil.
func (m *Message) SenderNodeInterface() *NodeInterface {
	return m.senderNodeInt
}

// AppendSignal appends a [Signal] to the last position of the [Message] payload.
// It may return an error if the signal name is already used within the message,
// or if the signal cannot fit in the available space left at the end of the message payload.
func (m *Message) AppendSignal(signal Signal) error {
	if signal == nil {
		return &ArgumentError{
			Name: "signal",
			Err:  ErrIsNil,
		}
	}

	if err := m.verifySignalName(signal.Name()); err != nil {
		return m.errorf(&AppendSignalError{
			EntityID: signal.EntityID(),
			Name:     signal.Name(),
			Err:      err,
		})
	}

	if err := m.signalPayload.append(signal); err != nil {
		return m.errorf(err)
	}

	m.addSignal(signal)

	return nil
}

// InsertSignal inserts a [Signal] at the given position of the [Message] payload.
// The start bit defines the index of the message payload where the signal will start.
// It may return an error if the signal name is already used within the message,
// or if the signal cannot fit in the available space left at the start bit.
func (m *Message) InsertSignal(signal Signal, startBit int) error {
	if signal == nil {
		return &ArgumentError{
			Name: "signal",
			Err:  ErrIsNil,
		}
	}

	if err := m.verifySignalName(signal.Name()); err != nil {
		return m.errorf(&InsertSignalError{
			EntityID: signal.EntityID(),
			Name:     signal.Name(),
			StartBit: startBit,
			Err:      err,
		})
	}

	if err := m.signalPayload.verifyAndInsert(signal, startBit); err != nil {
		return m.errorf(err)
	}

	m.addSignal(signal)

	return nil
}

// RemoveSignal removes a [Signal] that matches the given entity id from the [Message].
// It may return an error if the signal with the given entity id is not part of the message payload.
func (m *Message) RemoveSignal(signalEntityID EntityID) error {
	sig, err := m.signals.getValue(signalEntityID)
	if err != nil {
		return m.errorf(&RemoveEntityError{
			EntityID: signalEntityID,
			Err:      err,
		})
	}

	m.removeSignal(sig)

	m.signalPayload.remove(signalEntityID)

	return nil
}

// RemoveAllSignals removes all signals from the [Message].
func (m *Message) RemoveAllSignals() {
	for _, tmpSig := range m.signals.entries() {
		tmpSig.setParentMsg(nil)
	}

	m.signals.clear()
	m.signalNames.clear()

	m.signalPayload.removeAll()
}

// ShiftSignalLeft shifts the signal with the given entity id left by the given amount.
// It returns the amount of bits shifted.
func (m *Message) ShiftSignalLeft(signalEntityID EntityID, amount int) int {
	sig, err := m.signals.getValue(signalEntityID)
	if err != nil {
		return 0
	}

	return m.signalPayload.shiftLeft(sig.EntityID(), amount)
}

// ShiftSignalRight shifts the signal with the given entity id right by the given amount.
// It returns the amount of bits shifted.
func (m *Message) ShiftSignalRight(signalEntityID EntityID, amount int) int {
	sig, err := m.signals.getValue(signalEntityID)
	if err != nil {
		return 0
	}

	return m.signalPayload.shiftRight(sig.EntityID(), amount)
}

// CompactSignals compacts the [Message] payload.
func (m *Message) CompactSignals() {
	m.signalPayload.compact()
}

// Signals returns a slice of all signals in the [Message].
func (m *Message) Signals() []Signal {
	return m.signalPayload.signals
}

// GetSignal returns the [Signal] that matches the given entity id.
func (m *Message) GetSignal(signalEntityID EntityID) (Signal, error) {
	sig, err := m.signals.getValue(signalEntityID)
	if err != nil {
		return nil, m.errorf(&GetEntityError{
			EntityID: signalEntityID,
			Err:      err,
		})
	}
	return sig, nil
}

// GetSignalByName returns the [Signal] with the given name.
//
// It returns an [ErrNotFound] wrapped by a [NameError]
// if the name does not match any signal.
func (m *Message) GetSignalByName(name string) (Signal, error) {
	id, err := m.signalNames.getValue(name)
	if err != nil {
		return nil, m.errorf(&NameError{Err: err})
	}

	sig, err := m.signals.getValue(id)
	if err != nil {
		panic(err)
	}

	return sig, nil
}

// SignalNames returns a slice of all signal names in the [Message].
func (m *Message) SignalNames() []string {
	return m.signalNames.getKeys()
}

// SizeByte returns the size of the [Message] in bytes.
func (m *Message) SizeByte() int {
	return m.sizeByte
}

// SetPriority sets the message priority.
func (m *Message) SetPriority(priority MessagePriority) {
	m.priority = priority
}

// Priority returns the message priority.
func (m *Message) Priority() MessagePriority {
	return m.priority
}

// SetCycleTime sets the message cycle time in ms.
func (m *Message) SetCycleTime(cycleTime int) {
	m.cycleTime = cycleTime
}

// CycleTime returns the message cycle time.
func (m *Message) CycleTime() int {
	return m.cycleTime
}

// SetSendType sets the send type of the [Message].
func (m *Message) SetSendType(sendType MessageSendType) {
	m.sendType = sendType
}

// SendType returns the message send type.
func (m *Message) SendType() MessageSendType {
	return m.sendType
}

// SetDelayTime sets the delay time of the [Message].
func (m *Message) SetDelayTime(delayTime int) {
	m.delayTime = delayTime
}

// DelayTime returns the message delay time.
func (m *Message) DelayTime() int {
	return m.delayTime
}

// SetStartDelayTime sets the start delay time of the [Message].
func (m *Message) SetStartDelayTime(startDelayTime int) {
	m.startDelayTime = startDelayTime
}

// StartDelayTime returns the message start delay time.
func (m *Message) StartDelayTime() int {
	return m.startDelayTime
}

// AddReceiver adds a receiver to the [Message].
//
// It returns an [ArgumentError] if the given receiver is nil or
// a [ErrReceiverIsSender] wrapped by an [AddEntityError]
// if the receiver is the same as the sender.
func (m *Message) AddReceiver(receiver *NodeInterface) error {
	if receiver == nil {
		return m.errorf(&ArgumentError{
			Name: "receiver",
			Err:  ErrIsNil,
		})
	}

	if err := receiver.addReceivedMessage(m); err != nil {
		return m.errorf(&AddEntityError{
			EntityID: receiver.node.entityID,
			Name:     receiver.node.name,
			Err:      err,
		})
	}

	return nil
}

// RemoveReceiver removes a receiver from the [Message].
//
// It returns an [ErrNotFound] wrapped by a [RemoveEntityError]
// if the receiver with the given entity id is not found.
func (m *Message) RemoveReceiver(receiverEntityID EntityID) error {
	receiver, err := m.receivers.getValue(receiverEntityID)
	if err != nil {
		return m.errorf(&RemoveEntityError{
			EntityID: receiverEntityID,
			Err:      err,
		})
	}

	receiver.removeReceivedMessage(m)

	return nil
}

// Receivers returns a slice of all receivers of the [Message].
func (m *Message) Receivers() []*NodeInterface {
	recSlice := m.receivers.getValues()
	slices.SortFunc(recSlice, func(a, b *NodeInterface) int {
		return strings.Compare(a.node.name, b.node.name)
	})
	return recSlice
}

// SetByteOrder sets the byte order of the [Message].
func (m *Message) SetByteOrder(byteOrder MessageByteOrder) {
	m.byteOrder = byteOrder
}

// ByteOrder returns the byte order of the [Message].
func (m *Message) ByteOrder() MessageByteOrder {
	return m.byteOrder
}

// UpdateID updates the id of the [Message].
// It will also reset the static CAN-ID of the message.
//
// It may return an error if the new message id is already used within a [NodeInterface].
func (m *Message) UpdateID(newID MessageID) error {
	if m.id == newID && !m.hasStaticCANID {
		return nil
	}

	if m.hasSenderNodeInt() {
		if err := m.senderNodeInt.verifyMessageID(newID); err != nil {
			return m.errorf(err)
		}
	}

	m.staticCANID = 0
	m.hasStaticCANID = false
	m.id = newID

	return nil
}

// ID returns the id of the [Message].
func (m *Message) ID() MessageID {
	return m.id
}

// GetCANID returns the [GetCANID] associated to the [Message].
// If the message has a static CAN-ID, it will be returned.
// If the message does not have a sender [NodeInterface], it will return the message id.
// Otherwise, it will calculate the CAN-ID based on the [CANIDBuilder]
// provided by the [Bus] which owns the node interface.
func (m *Message) GetCANID() CANID {
	if m.hasStaticCANID {
		return m.staticCANID
	}

	if !m.hasSenderNodeInt() {
		return CANID(m.id)
	}

	nodeInt := m.senderNodeInt
	if !nodeInt.hasParentBus() {
		return CANID(m.id)
	}

	return nodeInt.parentBus.canIDBuilder.Calculate(m.priority, m.id, nodeInt.node.id)
}

// SetStaticCANID sets the static CAN-ID of the [Message].
//
// It returns a [CANIDError] if the given static CAN-ID is already used.
func (m *Message) SetStaticCANID(staticCANID CANID) error {
	if m.hasSenderNodeInt() {
		if err := m.senderNodeInt.verifyStaticCANID(staticCANID); err != nil {
			return err
		}
	}

	m.hasStaticCANID = true
	m.staticCANID = staticCANID
	m.id = MessageID(staticCANID)

	return nil
}

// HasStaticCANID returns whether the [Message] has a static CAN-ID.
func (m *Message) HasStaticCANID() bool {
	return m.hasStaticCANID
}

// AssignAttribute assigns the given attribute/value pair to the [Message].
//
// It returns an [ArgumentError] if the attribute is nil,
// or an [AttributeValueError] if the value does not conform to the attribute.
func (m *Message) AssignAttribute(attribute Attribute, value any) error {
	if err := m.addAttributeAssignment(attribute, m, value); err != nil {
		return m.errorf(err)
	}
	return nil
}

// RemoveAttributeAssignment removes the [AttributeAssignment]
// with the given attribute entity id from the [Message].
//
// It returns an [ErrNotFound] if the provided attribute entity id is not found.
func (m *Message) RemoveAttributeAssignment(attributeEntityID EntityID) error {
	if err := m.removeAttributeAssignment(attributeEntityID); err != nil {
		return m.errorf(err)
	}
	return nil
}

// GetAttributeAssignment returns the [AttributeAssignment]
// with the given attribute entity id from the [Message].
//
// It returns an [ErrNotFound] if the provided attribute entity id is not found.
func (m *Message) GetAttributeAssignment(attributeEntityID EntityID) (*AttributeAssignment, error) {
	attAss, err := m.getAttributeAssignment(attributeEntityID)
	if err != nil {
		return nil, m.errorf(err)
	}
	return attAss, nil
}

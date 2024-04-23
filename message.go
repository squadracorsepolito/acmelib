package acmelib

import (
	"fmt"
	"strings"
)

// MessageID is the bus unique identifier of a [Message].
// By default 11 bit message ids are used.
type MessageID uint32

func (mid MessageID) String() string {
	return fmt.Sprintf("%d", mid)
}

// MessageIDGeneratorFn is callback used for generating automatically
// the [MessageID] of a [Message]. It is triggered when a [Message] is added to
// a [Node] or when the former is removed. It takes as prameters the priority
// of the message, the number of messages sended by the node, and the node id,
// then it returns the computed message id.
// By default the messages calculate their 11 bit ids by putting the node id
// in the 4 lsb, the message count (number of messages sended by the node) from
// bit 4 to 9, and the priority in the 2 msb.
type MessageIDGeneratorFn func(priority MessagePriority, messageCount int, nodeID NodeID) (messageID MessageID)

var defMsgIDGenFn = func(priority MessagePriority, messageCount int, nodeID NodeID) (messageID MessageID) {
	messageID = (MessageID(messageCount) & 0b11111) << 4
	messageID |= MessageID(priority) << 9
	messageID |= MessageID(nodeID) & 0b1111
	return messageID
}

// MessagePriority rappresents the priority of a [Message].
// The priorities are very high, high, medium, and low.
// The higher priority has the value 0 and the lower has 3.
type MessagePriority uint

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
type MessageSendType string

const (
	// MessageSendTypeUnset defines an unset transmission type.
	MessageSendTypeUnset MessageSendType = "unset"
	// MessageSendTypeCyclic defines a cyclic transmission type.
	MessageSendTypeCyclic MessageSendType = "cyclic"
	// MessageSendTypeCyclicIfActive defines a cyclic if active transmission type.
	MessageSendTypeCyclicIfActive MessageSendType = "cyclicIfActive"
	// MessageSendTypeCyclicIfTriggered defines a cyclic if triggered transmission type.
	MessageSendTypeCyclicIfTriggered MessageSendType = "cyclicAndTriggered"
	// MessageSendTypeCyclicIfActiveAndTriggered defines a cyclic if active and triggered transmission type.
	MessageSendTypeCyclicIfActiveAndTriggered MessageSendType = "cyclicIfActiveAndTriggered"
)

// Message is the representation of data sent by a node thought the bus.
// It holds a list of signals that are contained in the message payload.
type Message struct {
	*attributeEntity

	senderNode *Node

	signals     *set[EntityID, Signal]
	signalNames *set[string, EntityID]

	signalPayload *signalPayload

	sizeByte int

	id         MessageID
	isStaticID bool
	idGenFn    MessageIDGeneratorFn
	priority   MessagePriority

	cycleTime      int
	sendType       MessageSendType
	delayTime      int
	startDelayTime int

	receivers *set[EntityID, *Node]
}

// NewMessage creates a new [Message] with the given name and size in bytes.
// By default a [MessagePriority] of [MessagePriorityVeryHigh] is used.
func NewMessage(name string, sizeByte int) *Message {
	return &Message{
		attributeEntity: newAttributeEntity(name, AttributeRefKindMessage),

		senderNode: nil,

		signals:     newSet[EntityID, Signal](),
		signalNames: newSet[string, EntityID](),

		signalPayload: newSignalPayload(sizeByte * 8),

		sizeByte: sizeByte,

		id:         0,
		isStaticID: false,
		idGenFn:    defMsgIDGenFn,
		priority:   MessagePriorityVeryHigh,

		cycleTime:      0,
		sendType:       MessageSendTypeUnset,
		delayTime:      0,
		startDelayTime: 0,

		receivers: newSet[EntityID, *Node](),
	}
}

func (m *Message) generateID(msgCount int, nodeID NodeID) {
	if m.isStaticID {
		return
	}
	m.id = m.idGenFn(m.priority, msgCount, nodeID)
}

func (m *Message) resetID() {
	m.isStaticID = false
	m.id = 0
}

func (m *Message) hasSenderNode() bool {
	return m.senderNode != nil
}

func (m *Message) errorf(err error) error {
	msgErr := &EntityError{
		Kind:     EntityKindMessage,
		EntityID: m.entityID,
		Name:     m.name,
		Err:      err,
	}
	if m.hasSenderNode() {
		return m.senderNode.errorf(msgErr)
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
		b.WriteString(fmt.Sprintf("%smessage_id: %d; is_static_id: %t\n", tabStr, m.id, m.isStaticID))
	}

	b.WriteString(fmt.Sprintf("%spriority: %d (very_high=0; low=3)\n", tabStr, m.priority))
	b.WriteString(fmt.Sprintf("%ssize: %d bytes\n", tabStr, m.sizeByte))

	if m.cycleTime != 0 {
		b.WriteString(fmt.Sprintf("%scycle_time: %d ms\n", tabStr, m.cycleTime))
	}

	if m.receivers.size() > 0 {
		b.WriteString(fmt.Sprintf("%sreceivers:\n", tabStr))
		for _, rec := range m.Receivers() {
			b.WriteString(fmt.Sprintf("%s\tname: %s; node_id: %d; entity_id: %s\n", tabStr, rec.name, rec.id, rec.entityID))
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

	if m.hasSenderNode() {
		if err := m.senderNode.messageNames.verifyKeyUnique(newName); err != nil {
			return m.errorf(&UpdateNameError{
				Err: &NameError{
					Name: newName,
					Err:  err,
				},
			})
		}

		m.senderNode.messageNames.modifyKey(m.name, newName, m.entityID)
	}

	m.name = newName

	return nil
}

// SenderNode returns the [Node] that sends the [Message].
// It returns nil if the message is not added to a node.
func (m *Message) SenderNode() *Node {
	return m.senderNode
}

func (m *Message) addSignal(sig Signal) {
	sigID := sig.EntityID()

	m.signals.add(sigID, sig)
	m.signalNames.add(sig.Name(), sigID)

	sig.setParentMsg(m)

	if sig.Kind() == SignalKindMultiplexer {
		muxSig, err := sig.ToMultiplexer()
		if err != nil {
			panic(err)
		}

		for tmpSigID, tmpSig := range muxSig.signals.entries() {
			m.signals.add(tmpSigID, tmpSig)
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

	if sig.Kind() == SignalKindMultiplexer {
		muxSig, err := sig.ToMultiplexer()
		if err != nil {
			panic(err)
		}

		for _, tmpSigID := range muxSig.signals.getKeys() {
			m.signals.remove(tmpSigID)
		}

		for _, tmpName := range muxSig.signalNames.getKeys() {
			m.signalNames.remove(tmpName)
		}
	}
}

// AppendSignal appends a [Signal] to the last position of the [Message] payload.
// It may return an error if the signal name is already used within the message,
// or if the signal cannot fit in the available space left at the end of the message payload.
func (m *Message) AppendSignal(signal Signal) error {
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

// SignalNames returns a slice of all signal names in the [Message].
func (m *Message) SignalNames() []string {
	return m.signalNames.getKeys()
}

// SizeByte returns the size of the [Message] in bytes.
func (m *Message) SizeByte() int {
	return m.sizeByte
}

// ID returns the message id.
func (m *Message) ID() MessageID {
	return m.id
}

// SetIDGeneratorFn sets the message id generator function.
func (m *Message) SetIDGeneratorFn(idGeneratorFn MessageIDGeneratorFn) {
	m.isStaticID = false
	m.idGenFn = idGeneratorFn
}

// SetID sets the message id.
// When the message id is set in this way, the message id generator function is not used anymore.
func (m *Message) SetID(messageID MessageID) {
	m.isStaticID = true
	m.id = messageID
}

// SetPriority sets the message priority.
func (m *Message) SetPriority(priority MessagePriority) {
	m.priority = priority
}

// Priority returns the message priority.
func (m *Message) Priority() MessagePriority {
	return m.priority
}

// SetCycleTime sets the message cycle time.
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

// AddReceiver adds a receiver [Node] to the [Message].
func (m *Message) AddReceiver(receiver *Node) {
	m.receivers.add(receiver.entityID, receiver)
}

// RemoveReceiver removes a receiver [Node] of the [Message].
func (m *Message) RemoveReceiver(receiverEntityID EntityID) {
	m.receivers.remove(receiverEntityID)
}

// Receivers returns a slice of all receiver nodes of the [Message].
func (m *Message) Receivers() []*Node {
	return m.receivers.getValues()
}

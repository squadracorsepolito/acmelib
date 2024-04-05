package acmelib

import (
	"fmt"
	"strings"
)

type MessageID uint32

type MessagePriority uint

const (
	MessagePriorityVeryHigh MessagePriority = iota
	MessagePriorityHigh
	MessagePriorityMedium
	MessagePriorityLow
)

type MessageIDGeneratorFn func(priority MessagePriority, messageCount int, nodeID NodeID) (messageID MessageID)

var defMsgIDGenFn = func(priority MessagePriority, messageCount int, nodeID NodeID) (messageID MessageID) {
	messageID = (MessageID(messageCount) & 0b11111) << 4
	messageID |= MessageID(priority) << 9
	messageID |= MessageID(nodeID) & 0b1111
	return messageID
}

type Message struct {
	*attributeEntity

	parentNodes *set[EntityID, *Node]
	parErrID    EntityID

	signals     *set[EntityID, Signal]
	signalNames *set[string, EntityID]

	signalPayload *signalPayload

	sizeByte int

	id         MessageID
	isStaticID bool
	idGenFn    MessageIDGeneratorFn
	priority   MessagePriority

	cycleTime uint

	receivers *set[EntityID, *Node]
}

func NewMessage(name, desc string, sizeByte int) *Message {
	return &Message{
		attributeEntity: newAttributeEntity(name, desc, AttributeRefKindMessage),

		parentNodes: newSet[EntityID, *Node]("parent node"),
		parErrID:    "",

		signals:     newSet[EntityID, Signal]("signal"),
		signalNames: newSet[string, EntityID]("signal name"),

		signalPayload: newSignalPayload(sizeByte * 8),

		sizeByte: sizeByte,

		id:         0,
		isStaticID: false,
		idGenFn:    defMsgIDGenFn,
		priority:   MessagePriorityVeryHigh,

		cycleTime: 0,

		receivers: newSet[EntityID, *Node]("receiver"),
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

// ---------------------------------------------------
// +++ START signalParent interface implementation +++
// ---------------------------------------------------

func (m *Message) errorf(err error) error {
	msgErr := fmt.Errorf(`message "%s": %w`, m.name, err)

	if m.parentNodes.size() > 0 {
		if m.parErrID != "" {
			parNode, err := m.parentNodes.getValue(m.parErrID)
			if err != nil {
				panic(err)
			}

			m.parErrID = ""
			return parNode.errorf(msgErr)
		}

		return m.parentNodes.getValues()[0].errorf(msgErr)
	}

	return msgErr
}

func (m *Message) GetSignalParentKind() SignalParentKind {
	return SignalParentKindMessage
}

func (m *Message) verifySignalName(_ EntityID, name string) error {
	return m.signalNames.verifyKey(name)
}

func (m *Message) modifySignalName(sigID EntityID, newName string) error {
	sig, err := m.signals.getValue(sigID)
	if err != nil {
		return err
	}

	oldName := sig.Name()
	m.signalNames.modifyKey(oldName, newName, sigID)

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
		return m.signalPayload.verifyBeforeGrow(sig, amount)
	}

	return m.signalPayload.verifyBeforeShrink(sig, -amount)
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

func (m *Message) ToParentMessage() (*Message, error) {
	return m, nil
}

func (m *Message) ToParentMultiplexerSignal() (*MultiplexerSignal, error) {
	return nil, fmt.Errorf(`cannot convert to "%s" signal parent is of kind "%s"`,
		SignalParentKindMultiplexerSignal, SignalParentKindMessage)
}

// -------------------------------------------------
// +++ END signalParent interface implementation +++
// -------------------------------------------------

func (m *Message) String() string {
	var builder strings.Builder

	builder.WriteString("\n+++START MESSAGE+++\n\n")
	builder.WriteString(m.toString())
	builder.WriteString(fmt.Sprintf("size: %d\n", m.sizeByte))

	signalsByPos := m.Signals()
	if len(signalsByPos) == 0 {
		return builder.String()
	}

	builder.WriteString("signals:\n")
	for _, sig := range signalsByPos {
		builder.WriteString(sig.String())
	}

	builder.WriteString("\n+++END MESSAGE+++\n")

	return builder.String()
}

func (m *Message) UpdateName(newName string) error {
	if m.name == newName {
		return nil
	}

	for _, tmpNode := range m.parentNodes.entries() {
		if err := tmpNode.messageNames.verifyKey(newName); err != nil {
			m.parErrID = tmpNode.entityID
			return m.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}

		tmpNode.modifyMessageName(m.entityID, newName)
	}

	m.name = newName

	return nil
}

func (m *Message) ParentNodes() []*Node {
	return m.parentNodes.getValues()
}

func (m *Message) AppendSignal(signal Signal) error {
	if err := m.verifySignalName(signal.EntityID(), signal.Name()); err != nil {
		return m.errorf(fmt.Errorf(`cannot append signal "%s" : %w`, signal.Name(), err))
	}

	if err := m.signalPayload.append(signal); err != nil {
		return m.errorf(err)
	}

	m.signals.add(signal.EntityID(), signal)
	m.signalNames.add(signal.Name(), signal.EntityID())

	signal.setParent(m)

	return nil
}

func (m *Message) InsertSignal(signal Signal, startBit int) error {
	if err := m.verifySignalName(signal.EntityID(), signal.Name()); err != nil {
		return m.errorf(fmt.Errorf(`cannot insert signal "%s" : %w`, signal.Name(), err))
	}

	if err := m.signalPayload.insert(signal, startBit); err != nil {
		return m.errorf(err)
	}

	m.signals.add(signal.EntityID(), signal)
	m.signalNames.add(signal.Name(), signal.EntityID())

	signal.setParent(m)

	return nil
}

func (m *Message) RemoveSignal(signalEntityID EntityID) error {
	sig, err := m.signals.getValue(signalEntityID)
	if err != nil {
		return m.errorf(fmt.Errorf(`cannot remove signal with entity id "%s" : %w`, signalEntityID, err))
	}

	if sig.Kind() == SignalKindMultiplexer {
		muxSig, err := sig.ToMultiplexer()
		if err != nil {
			panic(err)
		}

		for _, muxSignals := range muxSig.MuxSignals() {
			for _, tmpSig := range muxSignals {
				m.signals.remove(tmpSig.EntityID())
				m.signalNames.remove(tmpSig.Name())
			}
		}
	}

	sig.setParent(nil)

	m.signals.remove(signalEntityID)
	m.signalNames.remove(sig.Name())

	m.signalPayload.remove(signalEntityID)

	return nil
}

func (m *Message) RemoveAllSignals() {
	for _, tmpSig := range m.signals.entries() {
		tmpSig.setParent(nil)
	}

	m.signals.clear()
	m.signalNames.clear()

	m.signalPayload.removeAll()
}

func (m *Message) ShiftSignalLeft(signalEntityID EntityID, amount int) int {
	sig, err := m.signals.getValue(signalEntityID)
	if err != nil {
		return 0
	}

	return m.signalPayload.shiftLeft(sig, amount)
}

func (m *Message) ShiftSignalRight(signalEntityID EntityID, amount int) int {
	sig, err := m.signals.getValue(signalEntityID)
	if err != nil {
		return 0
	}

	return m.signalPayload.shiftRight(sig, amount)
}

func (m *Message) CompactSignals() {
	m.signalPayload.compact()
}

func (m *Message) Signals() []Signal {
	return m.signalPayload.signals
}

func (m *Message) GetSignal(signalEntityID EntityID) (Signal, error) {
	sig, err := m.signals.getValue(signalEntityID)
	if err != nil {
		return nil, m.errorf(fmt.Errorf(`cannot get signal with entity id "%s" : %w`, signalEntityID, err))
	}
	return sig, nil
}

func (m *Message) SignalNames() []string {
	return m.signalNames.getKeys()
}

func (m *Message) SizeByte() int {
	return m.sizeByte
}

func (m *Message) ID() MessageID {
	return m.id
}

func (m *Message) SetIDGeneratorFn(idGeneratorFn MessageIDGeneratorFn) {
	m.isStaticID = false
	m.idGenFn = idGeneratorFn
}

func (m *Message) SetID(messageID MessageID) {
	m.isStaticID = true
	m.id = messageID
}

func (m *Message) SetPriority(priority MessagePriority) {
	m.priority = priority
}

func (m *Message) Priority() MessagePriority {
	return m.priority
}

func (m *Message) SetCycleTime(cycleTime uint) {
	m.cycleTime = cycleTime
}

func (m *Message) CycleTime() uint {
	return m.cycleTime
}

func (m *Message) AddReceiver(receiver *Node) {
	m.receivers.add(receiver.entityID, receiver)
}

func (m *Message) RemoveReceiver(receiverEntityID EntityID) {
	m.receivers.remove(receiverEntityID)
}

func (m *Message) Receivers() []*Node {
	return m.receivers.getValues()
}

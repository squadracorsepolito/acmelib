package acmelib

import (
	"fmt"
	"strings"
)

type Message struct {
	*entity

	parentNode    *Node
	signals       *entityCollection[Signal]
	signalPayload *signalPayload
	sizeByte      int
	sizeBit       int
}

func NewMessage(name, desc string, sizeByte int) *Message {
	return &Message{
		entity: newEntity(name, desc),

		parentNode:    nil,
		signals:       newEntityCollection[Signal](),
		signalPayload: newSignalPayload(sizeByte * 8),
		sizeByte:      sizeByte,
		sizeBit:       sizeByte * 8,
	}
}

func (m *Message) hasParent() bool {
	return m.parentNode != nil
}

func (m *Message) errorf(err error) error {
	msgErr := fmt.Errorf(`message "%s": %v`, m.Name, err)
	if m.hasParent() {
		return m.parentNode.errorf(msgErr)
	}
	return msgErr
}

func (m *Message) String() string {
	var builder strings.Builder

	builder.WriteString("\n+++START MESSAGE+++\n\n")
	builder.WriteString(m.toString())
	builder.WriteString(fmt.Sprintf("size: %d\n", m.sizeByte))

	signalsByPos := m.GetSignalsByStartBit()
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

func (m *Message) GetSize() int {
	return m.sizeByte
}

func (m *Message) UpdateName(name string) error {
	if m.hasParent() {
		if err := m.parentNode.messages.updateEntityName(m.EntityID, m.Name, name); err != nil {
			return m.errorf(err)
		}
	}
	return m.entity.UpdateName(name)
}

func (m *Message) addSignal(sig Signal) error {
	if err := m.signals.addEntity(sig); err != nil {
		return m.errorf(err)
	}

	sig.setParentMessage(m)
	m.setUpdateTimeNow()

	return nil
}

func (m *Message) modifySignalSize(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	sig, err := m.GetSignalByEntityID(sigID)
	if err != nil {
		return err
	}

	if amount > 0 {
		return m.signalPayload.modifyStartBitsOnGrow(sig, amount)
	}

	m.signalPayload.modifyStartBitsOnShrink(sig, -amount)

	return nil
}

func (m *Message) AppendSignal(signal Signal) error {
	if err := m.signals.verifyEntityName(signal.GetName()); err != nil {
		return m.errorf(err)
	}

	if err := m.signalPayload.append(signal); err != nil {
		return m.errorf(err)
	}

	return m.addSignal(signal)
}

func (m *Message) InsertSignal(signal Signal, startBit int) error {
	if err := m.signals.verifyEntityName(signal.GetName()); err != nil {
		return m.errorf(err)
	}

	if err := m.signalPayload.insert(signal, startBit); err != nil {
		return m.errorf(err)
	}

	return m.addSignal(signal)
}

func (m *Message) RemoveSignal(signalID EntityID) error {
	m.signalPayload.remove(signalID)

	sig, err := m.GetSignalByEntityID(signalID)
	if err != nil {
		return err
	}

	sig.setParentMessage(nil)

	if err := m.signals.removeEntity(signalID); err != nil {
		return m.errorf(err)
	}

	m.setUpdateTimeNow()

	return nil
}

func (m *Message) CompactSignals() {
	m.signalPayload.compact()
}

func (m *Message) GetAvailableSignalSpaces() [][]int {
	positions := [][]int{}
	signals := m.GetSignalsByStartBit()

	from := 0
	for _, sig := range signals {
		sigStartBit := sig.GetStartBit()

		if from > sigStartBit {
			continue
		}

		if from < sigStartBit {
			positions = append(positions, []int{from, sigStartBit - 1})
		}

		from = sigStartBit + sig.GetSize()
	}

	if from < m.sizeBit {
		positions = append(positions, []int{from, m.sizeBit - 1})
	}

	return positions
}

func (m *Message) GetSignalsByName() []Signal {
	return sortByName(m.signals.listEntities())
}

func (m *Message) GetSignalsByCreateTime() []Signal {
	return sortByCreateTime(m.signals.listEntities())
}

func (m *Message) GetSignalsByUpdateTime() []Signal {
	return sortByUpdateTime(m.signals.listEntities())
}

func (m *Message) GetSignalsByStartBit() []Signal {
	return m.signalPayload.signals
}

func (m *Message) GetSignalByEntityID(id EntityID) (Signal, error) {
	return m.signals.getEntityByID(id)
}

func (m *Message) GetSignalByName(name string) (Signal, error) {
	return m.signals.getEntityByName(name)
}

func (m *Message) RemoveAllSignals() {
	for _, tmpSig := range m.signals.listEntities() {
		tmpSig.setParentMessage(nil)
	}

	m.signals.removeAllEntities()
	m.signalPayload.removeAll()
}

func (m *Message) ShiftSignalLeft(signalEntityID EntityID, amount int) int {
	sig, err := m.GetSignalByEntityID(signalEntityID)
	if err != nil {
		return 0
	}

	return m.signalPayload.shiftLeft(sig, amount)
}

func (m *Message) ShiftSignalRight(signalEntityID EntityID, amount int) int {
	sig, err := m.GetSignalByEntityID(signalEntityID)
	if err != nil {
		return 0
	}

	return m.signalPayload.shiftRight(sig, amount)
}

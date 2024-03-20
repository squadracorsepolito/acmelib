package acmelib

import (
	"fmt"
	"slices"
	"strings"
)

type Message struct {
	*entity
	ParentMessage *Node

	signals *entityCollection[*Signal]

	Size int

	bitSize int
}

func NewMessage(name, desc string, size int) *Message {
	return &Message{
		entity: newEntity(name, desc),

		signals: newEntityCollection[*Signal](),

		Size: size,

		bitSize: size * 8,
	}
}

func (m *Message) errorf(err error) error {
	msgErr := fmt.Errorf(`message "%s": %v`, m.Name, err)
	if m.ParentMessage != nil {
		return m.ParentMessage.errorf(msgErr)
	}
	return msgErr
}

func (m *Message) String() string {
	var builder strings.Builder

	builder.WriteString("\nMESSAGE\n")
	builder.WriteString(m.toString())
	builder.WriteString(fmt.Sprintf("size: %d\n", m.Size))

	signalsByPos := m.SignalsByPosition()
	if len(signalsByPos) == 0 {
		return builder.String()
	}

	builder.WriteString("signals:\n")
	for pos, sig := range signalsByPos {
		builder.WriteString(fmt.Sprintf("\t- %s position: %d, from start_bit: %d\n", sig.Name, pos, sig.StartBit))
	}

	return builder.String()
}

func (m *Message) UpdateName(name string) error {
	if err := m.ParentMessage.messages.updateEntityName(m.EntityID, m.Name, name); err != nil {
		return m.errorf(err)
	}

	return m.entity.UpdateName(name)
}

func (m *Message) addSignal(sig *Signal) error {
	if err := m.signals.addEntity(sig); err != nil {
		return m.errorf(err)
	}

	sig.ParentMessage = m
	m.setUpdateTimeNow()

	return nil
}

func (m *Message) verifySignalSize(sig *Signal) error {
	sigSize := sig.Type.Size
	if sigSize > m.Size*8 {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bytes`, sig.Name, sigSize, m.Size))
	}

	return nil
}

func (m *Message) AppendSignal(signal *Signal) error {
	if err := m.verifySignalSize(signal); err != nil {
		return err
	}

	if m.signals.size() == 0 {
		signal.StartBit = 0
		signal.Position = 0

		return m.addSignal(signal)
	}

	signalsByPos := m.SignalsByPosition()

	lastSig := signalsByPos[len(signalsByPos)-1]
	startBit := (lastSig.StartBit + lastSig.Type.Size)
	leftSpace := m.Size*8 - startBit
	sigSize := signal.Type.Size

	if leftSpace < sigSize {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bits left in the message`, signal.Name, sigSize, leftSpace))
	}

	if err := m.addSignal(signal); err != nil {
		return err
	}

	signal.StartBit = startBit
	signal.Position = lastSig.Position + 1

	return nil
}

func (m *Message) InsertSignalAtPosition(signal *Signal, pos int) error {
	if err := m.verifySignalSize(signal); err != nil {
		return err
	}

	signalsByPos := m.SignalsByPosition()

	sigCount := len(signalsByPos)
	if pos > sigCount {
		return m.errorf(fmt.Errorf(`signal "%s" position "%d" is out of bound, valid values are from "0" to "%d"`, signal.Name, pos, sigCount))
	}

	if sigCount == 0 {
		signal.StartBit = 0
		signal.Position = 0

		return m.addSignal(signal)
	}

	lastSig := signalsByPos[len(signalsByPos)-1]
	leftSpace := m.Size*8 - (lastSig.StartBit + lastSig.Type.Size)
	sigSize := signal.Type.Size

	if leftSpace < sigSize {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bits left in the message`, signal.Name, sigSize, leftSpace))
	}

	if sigCount == pos {
		signal.StartBit = lastSig.StartBit + lastSig.Type.Size
		signal.Position = pos

		return m.addSignal(signal)
	}

	inserted := false
	for _, tmpSig := range signalsByPos {
		if inserted {
			tmpSig.Position++
			tmpSig.StartBit += sigSize
			continue
		}

		if tmpSig.Position == pos {
			inserted = true

			if err := m.addSignal(signal); err != nil {
				return err
			}

			signal.Position = pos
			signal.StartBit = tmpSig.StartBit

			tmpSig.Position++
			tmpSig.StartBit += sigSize
		}
	}

	return nil
}

func (m *Message) InsertSignalAtStartBit(signal *Signal, startBit int) error {
	if err := m.verifySignalSize(signal); err != nil {
		return err
	}

	sigSize := signal.Type.Size
	if startBit+sigSize > m.Size*8 {
		return m.errorf(fmt.Errorf(`signal "%s" starting at bit "%d" of size "%d" bits cannot fit in "%d" bytes`, signal.Name, startBit, sigSize, m.Size))
	}

	signalsByPos := m.SignalsByPosition()
	sigCount := len(signalsByPos)

	if sigCount == 0 {
		signal.StartBit = startBit
		signal.Position = 0

		return m.addSignal(signal)
	}

	lastSig := signalsByPos[sigCount-1]
	leftSpace := m.Size*8 - (lastSig.StartBit + lastSig.Type.Size)

	if leftSpace < sigSize {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bits left in the message`, signal.Name, sigSize, leftSpace))
	}

	inserted := false
	for pos, tmpSig := range signalsByPos {
		if inserted {
			tmpSig.Position++
			continue
		}

		if startBit == tmpSig.StartBit {
			return m.errorf(fmt.Errorf(`signal "%s" cannot start at bit "%d" because signal "%s" alreay does`, signal.Name, startBit, tmpSig.Name))
		}

		if startBit > tmpSig.StartBit {
			tmpSigSpan := tmpSig.StartBit + tmpSig.Type.Size
			if startBit < tmpSigSpan {
				return m.errorf(fmt.Errorf(`signal "%s" cannot start at bit "%d" because signal "%s" spans from bit "%d" to "%d"`,
					signal.Name, startBit, tmpSig.Name, tmpSig.StartBit, tmpSigSpan-1))
			}

			continue
		}

		if startBit+sigSize > tmpSig.StartBit {
			return m.errorf(fmt.Errorf(`signal "%s" cannot start at bit "%d" because it will span over signal "%s"`, signal.Name, startBit, tmpSig.Name))
		}

		if err := m.addSignal(signal); err != nil {
			return err
		}

		signal.Position = pos
		signal.StartBit = startBit
		tmpSig.Position++
		inserted = true
	}

	if !inserted {
		signal.StartBit = startBit
		signal.Position = sigCount

		return m.addSignal(signal)
	}

	return nil
}

func (m *Message) RemoveSignal(signalID EntityID) error {
	removed := false
	for _, sig := range m.SignalsByPosition() {
		if removed {
			sig.Position--
			continue
		}

		if sig.EntityID == signalID {
			removed = true
		}
	}

	if err := m.signals.removeEntity(signalID); err != nil {
		return m.errorf(err)
	}

	m.setUpdateTimeNow()

	return nil
}

func (m *Message) CompactSignals() {
	lastStartBit := 0
	for _, sig := range m.SignalsByPosition() {
		if lastStartBit < sig.StartBit {
			sig.StartBit = lastStartBit
			lastStartBit += sig.Type.Size
		}
	}
}

func (m *Message) GetAvailableSignalPositions() []*SignalPosition {
	positions := []*SignalPosition{}

	from := 0
	for _, sig := range m.SignalsByPosition() {
		sigStartBit := sig.StartBit

		if from > sigStartBit {
			continue
		}

		if from < sigStartBit {
			positions = append(positions, NewSignalPosition(from, sigStartBit-1))
		}

		from = sigStartBit + sig.Type.Size
	}

	if from < m.bitSize {
		positions = append(positions, NewSignalPosition(from, m.bitSize-1))
	}

	return positions
}

func (m *Message) SignalsByName() []*Signal {
	return sortByName(m.signals.listEntities())
}

func (m *Message) SignalsByCreateTime() []*Signal {
	return sortByCreateTime(m.signals.listEntities())
}

func (m *Message) SignalsByUpdateTime() []*Signal {
	return sortByUpdateTime(m.signals.listEntities())
}

func (m *Message) SignalsByPosition() []*Signal {
	signals := m.signals.listEntities()
	slices.SortFunc(signals, func(a, b *Signal) int { return a.StartBit - b.StartBit })
	return signals
}

func (m *Message) GetSignalByEntityID(id EntityID) (*Signal, error) {
	return m.signals.getEntityByID(id)
}

func (m *Message) GetSignalByName(name string) (*Signal, error) {
	return m.signals.getEntityByName(name)
}

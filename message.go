package acmelib

import (
	"fmt"
	"slices"
	"strings"
)

type Message struct {
	*entity
	ParentNode *Node

	signals *entityCollection[*standardSignal]

	Size int

	bitSize int
}

func NewMessage(name, desc string, size int) *Message {
	return &Message{
		entity: newEntity(name, desc),

		signals: newEntityCollection[*standardSignal](),

		Size: size,

		bitSize: size * 8,
	}
}

func (m *Message) errorf(err error) error {
	msgErr := fmt.Errorf(`message "%s": %v`, m.Name, err)
	if m.ParentNode != nil {
		return m.ParentNode.errorf(msgErr)
	}
	return msgErr
}

func (m *Message) String() string {
	var builder strings.Builder

	builder.WriteString("\nMESSAGE\n")
	builder.WriteString(m.toString())
	builder.WriteString(fmt.Sprintf("size: %d\n", m.Size))

	signalsByPos := m.SignalsByStartBit()
	if len(signalsByPos) == 0 {
		return builder.String()
	}

	builder.WriteString("signals:\n")
	for _, sig := range signalsByPos {
		builder.WriteString(fmt.Sprintf("\t- %s index: %d, from start_bit: %d\n", sig.Name, sig.Index, sig.StartBit))
	}

	return builder.String()
}

func (m *Message) UpdateName(name string) error {
	if m.ParentNode != nil {
		if err := m.ParentNode.messages.updateEntityName(m.EntityID, m.Name, name); err != nil {
			return m.errorf(err)
		}
	}

	return m.entity.UpdateName(name)
}

func (m *Message) addSignal(sig *standardSignal) error {
	if err := m.signals.addEntity(sig); err != nil {
		return m.errorf(err)
	}

	sig.ParentMessage = m
	m.setUpdateTimeNow()

	return nil
}

func (m *Message) verifySignalSize(sig *standardSignal) error {
	sigSize := sig.BitSize()
	if sigSize > m.bitSize {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bytes`, sig.Name, sigSize, m.Size))
	}
	return nil
}

func (m *Message) AppendSignal(signal *standardSignal) error {
	if err := m.verifySignalSize(signal); err != nil {
		return err
	}

	if m.signals.getSize() == 0 {
		signal.StartBit = 0
		signal.Index = 0

		return m.addSignal(signal)
	}

	signals := m.SignalsByStartBit()
	sigCount := len(signals) - 1

	lastSig := signals[sigCount]
	startBit := lastSig.StartBit + lastSig.BitSize()
	leftSpace := m.bitSize - startBit
	sigSize := signal.BitSize()

	if sigSize > leftSpace {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bits left in the message`, signal.Name, sigSize, leftSpace))
	}

	if err := m.addSignal(signal); err != nil {
		return err
	}

	signal.StartBit = startBit
	signal.Index = sigCount + 1

	return nil
}

func (m *Message) InsertSignalAtIndex(signal *standardSignal, index int) error {
	if err := m.verifySignalSize(signal); err != nil {
		return err
	}

	signals := m.SignalsByStartBit()
	sigCount := len(signals)

	if index > sigCount {
		return m.errorf(fmt.Errorf(`signal "%s" index "%d" is out of range, valid values are from "0" to "%d"`, signal.Name, index, sigCount))
	}

	if sigCount == 0 {
		signal.StartBit = 0
		signal.Index = 0

		return m.addSignal(signal)
	}

	leftSpace := m.getMaxSignalBitSize()
	sigSize := signal.BitSize()

	if sigSize > leftSpace {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bits left in the message`, signal.Name, sigSize, leftSpace))
	}

	lastSig := signals[sigCount-1]
	if sigCount == index {
		signal.StartBit = lastSig.StartBit + lastSig.BitSize()
		signal.Index = index

		return m.addSignal(signal)
	}

	inserted := false
	for _, tmpSig := range signals {
		if inserted {
			tmpSig.Index++
			tmpSig.StartBit += sigSize
			continue
		}

		if tmpSig.Index == index {
			inserted = true

			if err := m.addSignal(signal); err != nil {
				return err
			}

			signal.Index = index
			signal.StartBit = tmpSig.StartBit

			tmpSig.Index++
			tmpSig.StartBit += sigSize
		}
	}

	return nil
}

func (m *Message) InsertSignalAtStartBit(signal *standardSignal, startBit int) error {
	if err := m.verifySignalSize(signal); err != nil {
		return err
	}

	sigSize := signal.BitSize()
	if startBit+sigSize > m.bitSize {
		return m.errorf(fmt.Errorf(`signal "%s" starting at bit "%d" of size "%d" bits cannot fit in "%d" bytes`, signal.Name, startBit, sigSize, m.Size))
	}

	signalsByPos := m.SignalsByStartBit()
	sigCount := len(signalsByPos)

	if sigCount == 0 {
		signal.StartBit = startBit
		signal.Index = 0

		return m.addSignal(signal)
	}

	leftSpace := m.getMaxSignalBitSize()
	if leftSpace < sigSize {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bits left in the message`, signal.Name, sigSize, leftSpace))
	}

	inserted := false
	for idx, tmpSig := range signalsByPos {
		if inserted {
			tmpSig.Index++
			continue
		}

		if startBit == tmpSig.StartBit {
			return m.errorf(fmt.Errorf(`signal "%s" cannot start at bit "%d" because signal "%s" alreay does`, signal.Name, startBit, tmpSig.Name))
		}

		if startBit > tmpSig.StartBit {
			tmpSigSpan := tmpSig.StartBit + tmpSig.BitSize()
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

		signal.Index = idx
		signal.StartBit = startBit
		tmpSig.Index++
		inserted = true
	}

	if !inserted {
		signal.StartBit = startBit
		signal.Index = sigCount

		return m.addSignal(signal)
	}

	return nil
}

func (m *Message) RemoveSignal(signalID EntityID) error {
	removed := false
	for _, sig := range m.SignalsByStartBit() {
		if removed {
			sig.Index--
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
	for _, sig := range m.SignalsByStartBit() {
		if lastStartBit < sig.StartBit {
			sig.StartBit = lastStartBit
			lastStartBit += sig.BitSize()
		}
	}
}

func (m *Message) getMaxSignalBitSize() int {
	max := 0

	positions := m.GetAvailableSignalSpaces()
	for _, pos := range positions {
		tmpSize := pos[1] - pos[0] + 1
		if tmpSize > max {
			max = tmpSize
		}
	}

	return max
}

func (m *Message) GetAvailableSignalSpaces() [][]int {
	positions := [][]int{}
	signals := m.SignalsByStartBit()

	from := 0
	for _, sig := range signals {
		sigStartBit := sig.StartBit

		if from > sigStartBit {
			continue
		}

		if from < sigStartBit {
			positions = append(positions, []int{from, sigStartBit - 1})
		}

		from = sigStartBit + sig.BitSize()
	}

	if from < m.bitSize {
		positions = append(positions, []int{from, m.bitSize - 1})
	}

	return positions
}

func (m *Message) SignalsByName() []*standardSignal {
	return sortByName(m.signals.listEntities())
}

func (m *Message) SignalsByCreateTime() []*standardSignal {
	return sortByCreateTime(m.signals.listEntities())
}

func (m *Message) SignalsByUpdateTime() []*standardSignal {
	return sortByUpdateTime(m.signals.listEntities())
}

func (m *Message) SignalsByStartBit() []*standardSignal {
	signals := m.signals.listEntities()
	slices.SortFunc(signals, func(a, b *standardSignal) int { return a.StartBit - b.StartBit })
	return signals
}

func (m *Message) GetSignalByEntityID(id EntityID) (*standardSignal, error) {
	return m.signals.getEntityByID(id)
}

func (m *Message) GetSignalByName(name string) (*standardSignal, error) {
	return m.signals.getEntityByName(name)
}

func (m *Message) shiftSignalsLeft(index, offset int) error {
	signals := m.SignalsByStartBit()
	sigCount := len(signals)

	newStartingBits := []int{}
	lastSigBitSize := 0
	i := index
	for {
		tmpSig := signals[i]
		newStartingBits = append(newStartingBits, tmpSig.StartBit+offset)

		i++
		if i == sigCount {
			lastSigBitSize = tmpSig.BitSize()
			break
		}
	}

	exceedingBits := newStartingBits[sigCount-1] + lastSigBitSize - m.bitSize
	if exceedingBits > 0 {
		return m.errorf(fmt.Errorf(`cannot shift signals left because it will exceed the message payload of "%d" bits`, exceedingBits))
	}

	for idx, startBit := range newStartingBits {
		signals[idx+index].StartBit = startBit
	}

	return nil
}

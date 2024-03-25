package acmelib

import (
	"fmt"
	"slices"
	"strings"
)

type Message struct {
	*entity

	parentNode *Node
	signals    *entityCollection[Signal]
	sizeByte   int
	sizeBit    int
}

func NewMessage(name, desc string, sizeByte int) *Message {
	return &Message{
		entity: newEntity(name, desc),

		parentNode: nil,
		signals:    newEntityCollection[Signal](),
		sizeByte:   sizeByte,
		sizeBit:    sizeByte * 8,
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

func (m *Message) verifySignalSize(sig Signal) error {
	sigSize := sig.GetSize()
	if sigSize > m.sizeBit {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bytes`, sig.GetName(), sigSize, m.sizeByte))
	}
	return nil
}

func (m *Message) AppendSignal(signal Signal) error {
	if err := m.verifySignalSize(signal); err != nil {
		return err
	}

	if m.signals.getSize() == 0 {
		signal.setStartBit(0)
		return m.addSignal(signal)
	}

	signals := m.GetSignalsByStartBit()
	sigCount := len(signals) - 1

	lastSig := signals[sigCount]
	startBit := lastSig.GetStartBit() + lastSig.GetSize()
	leftSpace := m.sizeBit - startBit
	sigSize := signal.GetSize()

	if sigSize > leftSpace {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bits left in the message`, signal.GetName(), sigSize, leftSpace))
	}

	if err := m.addSignal(signal); err != nil {
		return err
	}

	signal.setStartBit(startBit)

	return nil
}

func (m *Message) InsertSignal(signal Signal, startBit int) error {
	if err := m.verifySignalSize(signal); err != nil {
		return err
	}

	sigSize := signal.GetSize()
	if startBit+sigSize > m.sizeBit {
		return m.errorf(fmt.Errorf(`signal "%s" starting at bit "%d" of size "%d" bits cannot fit in "%d" bytes`, signal.GetName(), startBit, sigSize, m.sizeByte))
	}

	signalsByPos := m.GetSignalsByStartBit()
	sigCount := len(signalsByPos)

	if sigCount == 0 {
		signal.setStartBit(startBit)
		return m.addSignal(signal)
	}

	leftSpace := m.getMaxSignalBitSize()
	if leftSpace < sigSize {
		return m.errorf(fmt.Errorf(`signal "%s" of size "%d" bits cannot fit in "%d" bits left in the message`, signal.GetName(), sigSize, leftSpace))
	}

	inserted := false
	for _, tmpSig := range signalsByPos {
		if inserted {
			continue
		}

		if startBit == tmpSig.GetStartBit() {
			return m.errorf(fmt.Errorf(`signal "%s" cannot start at bit "%d" because signal "%s" alreay does`, signal.GetName(), startBit, tmpSig.GetName()))
		}

		if startBit > tmpSig.GetStartBit() {
			tmpSigSpan := tmpSig.GetStartBit() + tmpSig.GetSize()
			if startBit < tmpSigSpan {
				return m.errorf(fmt.Errorf(`signal "%s" cannot start at bit "%d" because signal "%s" spans from bit "%d" to "%d"`,
					signal.GetName(), startBit, tmpSig.GetName(), tmpSig.GetStartBit(), tmpSigSpan-1))
			}

			continue
		}

		if startBit+sigSize > tmpSig.GetStartBit() {
			return m.errorf(fmt.Errorf(`signal "%s" cannot start at bit "%d" because it will span over signal "%s"`, signal.GetName(), startBit, tmpSig.GetName()))
		}

		if err := m.addSignal(signal); err != nil {
			return err
		}

		signal.setStartBit(startBit)
		inserted = true
	}

	if !inserted {
		signal.setStartBit(startBit)
		return m.addSignal(signal)
	}

	return nil
}

func (m *Message) RemoveSignal(signalID EntityID) error {
	removed := false
	for _, sig := range m.GetSignalsByStartBit() {
		if removed {
			continue
		}

		if sig.GetEntityID() == signalID {
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
	for _, sig := range m.GetSignalsByStartBit() {
		if lastStartBit < sig.GetStartBit() {
			sig.setStartBit(lastStartBit)
			lastStartBit += sig.GetSize()
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
	signals := m.signals.listEntities()
	slices.SortFunc(signals, func(a, b Signal) int { return a.GetStartBit() - b.GetStartBit() })
	return signals
}

func (m *Message) GetSignalByEntityID(id EntityID) (Signal, error) {
	return m.signals.getEntityByID(id)
}

func (m *Message) GetSignalByName(name string) (Signal, error) {
	return m.signals.getEntityByName(name)
}

func (m *Message) modifySignalSize(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	signals := m.GetSignalsByStartBit()

	if amount < 0 {
		found := false
		for _, tmpSig := range signals {
			if found {
				tmpSig.setStartBit(tmpSig.GetStartBit() + amount)
				continue
			}

			if tmpSig.GetEntityID() == sigID {
				found = true
			}
		}

		return nil
	}

	index := 0
	for idx, tmpSig := range signals {
		if tmpSig.GetEntityID() == sigID {
			index = idx + 1
			break
		}
	}

	sigCount := len(signals)
	if index == sigCount {
		sig := signals[index-1]
		exceedingBits := sig.GetStartBit() + sig.GetSize() + amount - m.sizeBit
		if exceedingBits > 0 {
			return fmt.Errorf(`cannot grow signal size because it will cause an overflow of "%d" bits in the message payload`, exceedingBits)
		}

		return nil
	}

	newStartingBits := []int{}
	lastSigBitSize := 0
	i := index
	for {
		tmpSig := signals[i]
		newStartingBits = append(newStartingBits, tmpSig.GetStartBit()+amount)

		i++
		if i == sigCount {
			lastSigBitSize = tmpSig.GetSize()
			break
		}
	}

	exceedingBits := newStartingBits[len(newStartingBits)-1] + lastSigBitSize - m.sizeBit
	if exceedingBits > 0 {
		return fmt.Errorf(`cannot grow signal size because it will cause an overflow of "%d" bits in the message payload`, exceedingBits)
	}

	for idx, startBit := range newStartingBits {
		signals[idx+index].setStartBit(startBit)
	}

	return nil
}

func (m *Message) RemoveAllSignals() {
	m.signals.removeAllEntities()
}

func (m *Message) ShiftSignalLeft(signalEntityID EntityID, amount int) int {
	if amount <= 0 {
		return 0
	}

	signals := m.GetSignalsByStartBit()
	perfShift := amount

	var prevSig Signal
	for idx, tmpSig := range signals {
		if idx > 0 {
			prevSig = signals[idx-1]
		}

		if tmpSig.GetEntityID() == signalEntityID {
			tmpStartBit := tmpSig.GetStartBit()
			targetStartBit := tmpStartBit - amount

			if targetStartBit < 0 {
				targetStartBit = 0
			}

			if prevSig != nil {
				prevEndBit := prevSig.GetStartBit() + prevSig.GetSize()

				if targetStartBit < prevEndBit {
					targetStartBit = prevEndBit
				}
			}

			tmpSig.setStartBit(targetStartBit)
			perfShift = tmpStartBit - targetStartBit

			break
		}
	}

	return perfShift
}

func (m *Message) ShiftSignalRight(signalEntityID EntityID, amount int) int {
	if amount <= 0 {
		return 0
	}

	signals := m.GetSignalsByStartBit()
	perfShift := amount

	var nextSig Signal
	for idx, tmpSig := range signals {
		if idx == len(signals)-1 {
			nextSig = nil
		} else {
			nextSig = signals[idx+1]
		}

		if tmpSig.GetEntityID() == signalEntityID {
			tmpStartBit := tmpSig.GetStartBit()
			targetStartBit := tmpStartBit + amount
			targetEndBit := targetStartBit + tmpSig.GetSize()

			if targetEndBit > m.sizeBit {
				targetStartBit = m.sizeBit - tmpSig.GetSize()
			}

			if nextSig != nil {
				nextStartBit := nextSig.GetStartBit()

				if targetEndBit > nextStartBit {
					targetStartBit = nextStartBit - tmpSig.GetSize()
				}
			}

			tmpSig.setStartBit(targetStartBit)
			perfShift = targetStartBit - tmpStartBit

			break
		}
	}

	return perfShift
}

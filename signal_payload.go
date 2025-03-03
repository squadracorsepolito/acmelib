package acmelib

import (
	"golang.org/x/exp/slices"
)

type signalPayload struct {
	size    int
	signals []Signal
}

func newSignalPayload(size int) *signalPayload {
	return &signalPayload{
		size:    size,
		signals: []Signal{},
	}
}

func (sp *signalPayload) verifyBeforeAppend(sig Signal) error {
	sigSize := sig.GetSize()
	sigCount := len(sp.signals)
	if sigCount == 0 {
		if sigSize > sp.size {
			return &SignalSizeError{
				Size: sigSize,
				Err:  ErrOutOfBounds,
			}
		}

		return nil
	}

	lastSig := sp.signals[sigCount-1]
	trailingSpace := sp.size - (lastSig.getRelStartBit() + lastSig.GetSize())

	if sigSize > trailingSpace {
		return &SignalSizeError{
			Size: sigSize,
			Err:  ErrNoSpaceLeft,
		}
	}

	return nil
}

func (sp *signalPayload) append(sig Signal) error {
	if err := sp.verifyBeforeAppend(sig); err != nil {
		return &AppendSignalError{
			EntityID: sig.EntityID(),
			Name:     sig.Name(),
			Err:      err,
		}
	}

	if len(sp.signals) == 0 {
		sig.setRelStartBit(0)
	} else {
		lastSig := sp.signals[len(sp.signals)-1]
		sig.setRelStartBit(lastSig.getRelStartBit() + lastSig.GetSize())
	}

	sp.signals = append(sp.signals, sig)

	return nil
}

func (sp *signalPayload) verifyBeforeInsert(sig Signal, startBit int) error {
	if startBit < 0 {
		return &StartBitError{
			StartBit: startBit,
			Err:      ErrIsNegative,
		}
	}

	sigSize := sig.GetSize()
	endBit := startBit + sigSize

	if sigSize > sp.size {
		return &SignalSizeError{
			Size: sigSize,
			Err:  ErrOutOfBounds,
		}
	}

	if endBit > sp.size {
		return &SignalSizeError{
			Size: sigSize,
			Err:  ErrNoSpaceLeft,
		}
	}

	for _, tmpSig := range sp.signals {
		tmpStartBit := tmpSig.getRelStartBit()
		tmpEndBit := tmpStartBit + tmpSig.GetSize()

		if endBit <= tmpStartBit {
			break
		}

		if startBit >= tmpEndBit {
			continue
		}

		if startBit >= tmpStartBit || endBit > tmpStartBit {
			return &StartBitError{
				StartBit: startBit,
				Err:      ErrIntersect,
			}
		}
	}

	return nil
}

func (sp *signalPayload) insert(sig Signal, startBit int) {
	if len(sp.signals) == 0 {
		sig.setRelStartBit(startBit)
		sp.signals = append(sp.signals, sig)
		return
	}

	inserted := false
	for idx, tmpSig := range sp.signals {
		tmpStartBit := tmpSig.getRelStartBit()

		if tmpStartBit > startBit {
			inserted = true
			sp.signals = slices.Insert(sp.signals, idx, sig)
			break
		}
	}

	if !inserted {
		sp.signals = append(sp.signals, sig)
	}

	sig.setRelStartBit(startBit)
}

func (sp *signalPayload) verifyAndInsert(sig Signal, startBit int) error {
	if err := sp.verifyBeforeInsert(sig, startBit); err != nil {
		return &InsertSignalError{
			EntityID: sig.EntityID(),
			Name:     sig.Name(),
			StartBit: startBit,
			Err:      err,
		}
	}

	sp.insert(sig, startBit)

	return nil
}

func (sp *signalPayload) remove(sigID EntityID) {
	sp.signals = slices.DeleteFunc(sp.signals, func(s Signal) bool { return s.EntityID() == sigID })
}

func (sp *signalPayload) removeAll() {
	sp.signals = []Signal{}
}

func (sp *signalPayload) compact() {
	lastStartBit := 0
	for _, sig := range sp.signals {
		tmpStartBit := sig.getRelStartBit()

		if tmpStartBit == lastStartBit {
			lastStartBit += sig.GetSize()
			continue
		}

		if lastStartBit < tmpStartBit {
			sig.setRelStartBit(lastStartBit)
			lastStartBit += sig.GetSize()
		}
	}
}

func (sp *signalPayload) verifyBeforeShrink(sig Signal, amount int) error {
	if amount < 0 {
		return ErrIsNegative
	}

	sizeDiff := sig.GetSize() - amount
	if sizeDiff < 0 {
		return ErrIsNegative
	}

	if sizeDiff == 0 {
		return ErrIsZero
	}

	return nil
}

func (sp *signalPayload) modifyStartBitsOnShrink(sig Signal, amount int) error {
	if amount == 0 {
		return nil
	}

	if err := sp.verifyBeforeShrink(sig, amount); err != nil {
		return &SignalSizeError{
			Size: sig.GetSize() - amount,
			Err:  err,
		}
	}

	found := false
	for _, tmpSig := range sp.signals {
		if found {
			tmpSig.setRelStartBit(tmpSig.getRelStartBit() - amount)
			continue
		}

		if sig.EntityID() == tmpSig.EntityID() {
			found = true
		}
	}

	return nil
}

func (sp *signalPayload) verifyBeforeGrow(sig Signal, amount int) error {
	if amount < 0 {
		return ErrIsNegative
	}

	availableSpace := 0
	prevEndBit := 0
	found := false

	for _, tmpSig := range sp.signals {
		tmpStartBit := tmpSig.getRelStartBit()

		if found {
			availableSpace += tmpStartBit - prevEndBit
		} else if tmpSig.EntityID() == sig.EntityID() {
			found = true
		}

		prevEndBit = tmpStartBit + tmpSig.GetSize()
	}

	availableSpace += sp.size - prevEndBit

	if amount > availableSpace {
		return ErrNoSpaceLeft
	}

	return nil
}

func (sp *signalPayload) modifyStartBitsOnGrow(sig Signal, amount int) error {
	if amount == 0 {
		return nil
	}

	if err := sp.verifyBeforeGrow(sig, amount); err != nil {
		return &SignalSizeError{
			Size: sig.GetSize() + amount,
			Err:  err,
		}
	}

	prevEndBit := 0
	spaces := []int{}
	nextSigIdx := 0
	found := false

	for idx, tmpSig := range sp.signals {
		tmpStartBit := tmpSig.getRelStartBit()

		if found {
			space := tmpStartBit - prevEndBit
			spaces = append(spaces, space)

		} else if sig.EntityID() == tmpSig.EntityID() {
			if idx == len(sp.signals)-1 {
				return nil
			}

			found = true
			nextSigIdx = idx + 1
		}

		prevEndBit = tmpStartBit + tmpSig.GetSize()
	}

	spaces = append(spaces, sp.size-prevEndBit)

	spaceIdx := 0
	acc := amount
	for i := nextSigIdx; i < len(sp.signals); i++ {
		tmpSpace := spaces[spaceIdx]

		if tmpSpace >= acc {
			break
		}

		acc -= tmpSpace
		tmpSig := sp.signals[i]
		tmpSig.setRelStartBit(tmpSig.getRelStartBit() + acc)
		spaceIdx++
	}

	return nil
}

func (sp *signalPayload) verifyBeforeResize(newSize int) error {
	if newSize > sp.size {
		return nil
	}

	lastSig := sp.signals[len(sp.signals)-1]
	if lastSig.getRelStartBit()+lastSig.GetSize() > newSize {
		return ErrTooSmall
	}

	return nil
}

func (sp *signalPayload) resize(newSize int) error {
	if err := sp.verifyBeforeResize(newSize); err != nil {
		return &MessageSizeError{
			Size: newSize,
			Err:  err,
		}
	}

	sp.size = newSize

	return nil
}

func (sp *signalPayload) shiftLeft(sigID EntityID, amount int) int {
	if amount <= 0 {
		return 0
	}

	perfShift := amount
	var prevSig Signal

	for idx, tmpSig := range sp.signals {
		if idx > 0 {
			prevSig = sp.signals[idx-1]
		}

		if sigID == tmpSig.EntityID() {
			tmpStartBit := tmpSig.getRelStartBit()
			targetStartBit := tmpStartBit - amount

			if targetStartBit < 0 {
				targetStartBit = 0
			}

			if prevSig != nil {
				prevEndBit := prevSig.getRelStartBit() + prevSig.GetSize()

				if targetStartBit < prevEndBit {
					targetStartBit = prevEndBit
				}
			}

			tmpSig.setRelStartBit(targetStartBit)
			perfShift = tmpStartBit - targetStartBit

			break
		}
	}

	return perfShift
}

func (sp *signalPayload) shiftRight(sigID EntityID, amount int) int {
	if amount <= 0 {
		return 0
	}

	perfShift := amount
	var nextSig Signal

	for idx, tmpSig := range sp.signals {
		if idx == len(sp.signals)-1 {
			nextSig = nil
		} else {
			nextSig = sp.signals[idx+1]
		}

		if sigID == tmpSig.EntityID() {
			tmpStartBit := tmpSig.getRelStartBit()
			targetStartBit := tmpStartBit + amount
			targetEndBit := targetStartBit + tmpSig.GetSize()

			if targetEndBit > sp.size {
				targetStartBit = sp.size - tmpSig.GetSize()
			}

			if nextSig != nil {
				nextStartBit := nextSig.getRelStartBit()

				if targetEndBit > nextStartBit {
					targetStartBit = nextStartBit - tmpSig.GetSize()
				}
			}

			tmpSig.setRelStartBit(targetStartBit)
			perfShift = targetStartBit - tmpStartBit

			break
		}
	}

	return perfShift
}

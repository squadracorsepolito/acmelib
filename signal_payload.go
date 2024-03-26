package acmelib

import (
	"errors"
	"fmt"
	"sync"

	"golang.org/x/exp/slices"
)

type signalPayload struct {
	size    int
	signals []Signal
	mux     *sync.RWMutex
}

func newSignalPayload(size int) *signalPayload {
	return &signalPayload{
		size:    size,
		signals: []Signal{},
		mux:     &sync.RWMutex{},
	}
}

func (sp *signalPayload) verifyBeforeAppend(sig Signal) error {
	sp.mux.RLock()
	defer sp.mux.RUnlock()

	sigSize := sig.GetSize()
	sigCount := len(sp.signals)
	if sigCount == 0 {
		if sigSize > sp.size {
			return fmt.Errorf(`signal of size "%d" exceeds the max payload size ("%d")`, sigSize, sp.size)
		}

		return nil
	}

	lastSig := sp.signals[sigCount-1]
	trailingSpace := sp.size - (lastSig.GetStartBit() + lastSig.GetSize())

	if sigSize > trailingSpace {
		return fmt.Errorf(`signal of size "%d" exceeds the available space ("%d") at the end of the payload`, sigSize, trailingSpace)
	}

	return nil
}

func (sp *signalPayload) append(sig Signal) error {
	if err := sp.verifyBeforeAppend(sig); err != nil {
		return fmt.Errorf(`cannot append signal "%s" : %v`, sig.GetName(), err)
	}

	sp.mux.Lock()
	defer sp.mux.Unlock()

	if len(sp.signals) == 0 {
		sig.setStartBit(0)
	} else {
		lastSig := sp.signals[len(sp.signals)-1]
		sig.setStartBit(lastSig.GetStartBit() + lastSig.GetSize())
	}

	sp.signals = append(sp.signals, sig)

	return nil
}

func (sp *signalPayload) verifyBeforeInsert(sig Signal, startBit int) error {
	if startBit < 0 {
		return errors.New("start bit cannot be negative")
	}

	sigSize := sig.GetSize()
	endBit := startBit + sigSize

	if sigSize > sp.size {
		return fmt.Errorf(`signal of size "%d" exceeds the max payload size of "%d"`, sigSize, sp.size)
	}

	if endBit > sp.size {
		return fmt.Errorf(`signal of size "%d" starting at "%d" exceeds the max payload size ("%d")`, sigSize, startBit, sp.size)
	}

	sp.mux.RLock()
	defer sp.mux.RUnlock()

	for _, tmpSig := range sp.signals {
		tmpStartBit := tmpSig.GetStartBit()
		tmpEndBit := tmpStartBit + tmpSig.GetSize()

		if endBit <= tmpStartBit {
			break
		}

		if startBit >= tmpEndBit {
			continue
		}

		if startBit >= tmpStartBit || endBit > tmpStartBit {
			return fmt.Errorf(`signal of size "%d" starting at "%d" intersects signal "%s" (start bit "%d", size "%d")`,
				sigSize, startBit, tmpSig.GetName(), tmpStartBit, tmpSig.GetSize())
		}
	}

	return nil
}

func (sp *signalPayload) insert(sig Signal, startBit int) error {
	if err := sp.verifyBeforeInsert(sig, startBit); err != nil {
		return fmt.Errorf(`cannot insert signal "%s" : %v`, sig.GetName(), err)
	}

	sp.mux.Lock()
	defer sp.mux.Unlock()

	if len(sp.signals) == 0 {
		sig.setStartBit(startBit)
		sp.signals = append(sp.signals, sig)

		return nil
	}

	inserted := false
	for idx, tmpSig := range sp.signals {
		tmpStartBit := tmpSig.GetStartBit()

		if tmpStartBit > startBit {
			inserted = true
			sp.signals = slices.Insert(sp.signals, idx, sig)
			break
		}
	}

	if !inserted {
		sp.signals = append(sp.signals, sig)
	}

	sig.setStartBit(startBit)

	return nil
}

func (sp *signalPayload) remove(sigID EntityID) {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	sp.signals = slices.DeleteFunc(sp.signals, func(s Signal) bool { return s.GetEntityID() == sigID })
}

func (sp *signalPayload) removeAll() {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	sp.signals = []Signal{}
}

func (sp *signalPayload) compact() {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	lastStartBit := 0
	for _, sig := range sp.signals {
		tmpStartBit := sig.GetStartBit()

		if tmpStartBit == lastStartBit {
			lastStartBit += sig.GetSize()
			continue
		}

		if lastStartBit < tmpStartBit {
			sig.setStartBit(lastStartBit)
			lastStartBit += sig.GetSize()
		}
	}
}

func (sp *signalPayload) modifyStartBitsOnShrink(sig Signal, amount int) {
	if amount <= 0 {
		return
	}

	sp.mux.Lock()
	defer sp.mux.Unlock()

	found := false
	for _, tmpSig := range sp.signals {
		if found {
			tmpSig.setStartBit(tmpSig.GetStartBit() - amount)
			continue
		}

		if sig.GetEntityID() == tmpSig.GetEntityID() {
			found = true
		}
	}
}

func (sp *signalPayload) verifyBeforeGrow(sig Signal, amount int) error {
	if amount < 0 {
		return errors.New("amount cannot be negative")
	}

	availableSpace := 0
	prevEndBit := 0
	found := false

	sp.mux.RLock()
	defer sp.mux.RUnlock()

	for _, tmpSig := range sp.signals {
		tmpStartBit := tmpSig.GetStartBit()

		if found {
			availableSpace += tmpStartBit - prevEndBit
		} else if tmpSig.GetEntityID() == sig.GetEntityID() {
			found = true
		}

		prevEndBit = tmpStartBit + tmpSig.GetSize()
	}

	availableSpace += sp.size - prevEndBit

	if amount > availableSpace {
		return fmt.Errorf(`amount "%d" exceeds the available space left at the right of the signal ("%d")`, amount, availableSpace)
	}

	return nil

	// return amount <= availableSpace
}

func (sp *signalPayload) modifyStartBitsOnGrow(sig Signal, amount int) error {
	if amount == 0 {
		return nil
	}

	if err := sp.verifyBeforeGrow(sig, amount); err != nil {
		return fmt.Errorf(`cannot grow signal "%s" : %v`, sig.GetName(), err)
	}

	sp.mux.Lock()
	defer sp.mux.Unlock()

	prevEndBit := 0
	spaces := []int{}
	nextSigIdx := 0
	found := false

	for idx, tmpSig := range sp.signals {
		tmpStartBit := tmpSig.GetStartBit()

		if found {
			space := tmpStartBit - prevEndBit
			spaces = append(spaces, space)

		} else if sig.GetEntityID() == tmpSig.GetEntityID() {
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
		tmpSig.setStartBit(tmpSig.GetStartBit() + acc)
		spaceIdx++
	}

	return nil
}

func (sp *signalPayload) shiftLeft(sig Signal, amount int) int {
	if amount <= 0 {
		return 0
	}

	perfShift := amount
	var prevSig Signal

	sp.mux.Lock()
	defer sp.mux.Unlock()

	for idx, tmpSig := range sp.signals {
		if idx > 0 {
			prevSig = sp.signals[idx-1]
		}

		if sig.GetEntityID() == tmpSig.GetEntityID() {
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

func (sp *signalPayload) shiftRight(sig Signal, amount int) int {
	if amount <= 0 {
		return 0
	}

	perfShift := amount
	var nextSig Signal

	sp.mux.Lock()
	defer sp.mux.Unlock()

	for idx, tmpSig := range sp.signals {
		if idx == len(sp.signals)-1 {
			nextSig = nil
		} else {
			nextSig = sp.signals[idx+1]
		}

		if sig.GetEntityID() == tmpSig.GetEntityID() {
			tmpStartBit := tmpSig.GetStartBit()
			targetStartBit := tmpStartBit + amount
			targetEndBit := targetStartBit + tmpSig.GetSize()

			if targetEndBit > sp.size {
				targetStartBit = sp.size - tmpSig.GetSize()
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

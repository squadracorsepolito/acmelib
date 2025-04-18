package acmelib

import "github.com/squadracorsepolito/acmelib/internal/ibst"

type SL struct {
	sizeByte int
	tree     *ibst.Tree[Signal]
	filters  []*SignalLayoutFilter
}

func newSL(sizeByte int) *SL {
	return &SL{
		sizeByte: sizeByte,
		tree:     ibst.NewTree[Signal](),
		filters:  []*SignalLayoutFilter{},
	}
}

// genFilters generates the signal layout filters.
// It must be called every time the layout is changed.
func (sl *SL) genFilters() {
	sl.filters = []*SignalLayoutFilter{}

	for _, sig := range sl.tree.GetAllIntervals() {
		sigSize := sig.GetSize()
		startPos := sig.GetRelativeStartPos()
		endianness := sig.Endianness()

		firstIdx := startPos / 8
		lastIdx := (startPos + sigSize - 1) / 8

		// Check if signal fit in a single row
		if firstIdx == lastIdx {
			// Calculate left offset
			leftOffset := startPos % 8

			// Make the mask of the size of the signal and shift it to the left offset
			mask := 1<<sigSize - 1
			mask <<= leftOffset

			sl.filters = append(sl.filters, &SignalLayoutFilter{
				signal:     sig,
				byteIdx:    firstIdx,
				mask:       uint8(mask),
				length:     sigSize,
				leftOffset: leftOffset,
			})

			continue
		}

		remainingBits := sigSize
		for i := firstIdx; i <= lastIdx; i++ {
			// Set mask, length and left offset to default
			mask := 0xff
			length := 8
			leftOffset := 0

			// Check if it is not the first or last byte
			if i != firstIdx && i != lastIdx {
				goto appendFilter
			}

			// Check if it is the first byte
			if i == firstIdx {
				tmpOffset := startPos % 8
				length = 8 - tmpOffset

				if endianness == MessageByteOrderBigEndian {
					// In case of big endian, shift the mask right
					// because startPos always refers to the first bit
					// in little endian
					mask >>= tmpOffset
				} else {
					mask <<= tmpOffset
					leftOffset = tmpOffset
				}

				goto appendFilter
			}

			// Last byte
			length = remainingBits
			mask = 1<<remainingBits - 1

			if endianness == MessageByteOrderBigEndian {
				// For the same reason as above, shift the mask left
				leftOffset = 8 - remainingBits
				mask <<= leftOffset
			}

		appendFilter:
			sl.filters = append(sl.filters, &SignalLayoutFilter{
				signal:     sig,
				byteIdx:    i,
				mask:       uint8(mask),
				length:     length,
				leftOffset: leftOffset,
			})

			remainingBits -= length
		}
	}
}

// verifyStartPos checks if the start position is valid.
func (sl *SL) verifyStartPos(startPos int) error {
	if startPos < 0 {
		return &StartPosError{
			StartPos: startPos,
			Err:      ErrIsNegative,
		}
	}

	if startPos > sl.sizeByte*8 {
		return &StartPosError{
			StartPos: startPos,
			Err:      ErrOutOfBounds,
		}
	}

	return nil
}

// verifySize checks if the size is valid.
func (sl *SL) verifySize(size int) error {
	if size < 0 {
		return &SignalSizeError{
			Size: size,
			Err:  ErrIsNegative,
		}
	}

	if size > sl.sizeByte*8 {
		return &SignalSizeError{
			Size: size,
			Err:  ErrOutOfBounds,
		}
	}

	return nil
}

// verifyStartPosPlusSize checks if the start position plus the size is valid.
// It doesn't check if the size or the start position are valid, it only checks if the sum is valid.
func (sl *SL) verifyStartPosPlusSize(startPos, size int) error {
	if startPos+size > sl.sizeByte*8 {
		return &SignalSizeError{
			Size: size,
			Err:  ErrOutOfBounds,
		}
	}

	return nil
}

// verifyInsert checks if the signal does not intersect with another signal.
func (sl *SL) verifyInsert(sig Signal, startPos int) error {
	if err := sl.verifyStartPos(startPos); err != nil {
		return err
	}

	size := sig.GetSize()
	if err := sl.verifySize(size); err != nil {
		return err
	}

	if err := sl.verifyStartPosPlusSize(startPos, size); err != nil {
		return err
	}

	// Set the start position as the low interval
	sig.SetLow(startPos)

	// Reset the low interval
	defer sig.SetLow(0)

	// Check if the signal intersects with another
	if sl.tree.Intersects(sig) {
		return &StartPosError{
			StartPos: startPos,
			Err:      ErrIntersects,
		}
	}

	return nil
}

// insert inserts the signal into the signal layout.
// It must be called after the verify function since it assumes the signal is valid.
func (sl *SL) insert(sig Signal, startPos int) {
	sig.setRelativeStartPos(startPos)
	sl.tree.Insert(sig)

	// Regenerate the filters
	sl.genFilters()
}

// verifyAndInsert checks and inserts the signal into the signal layout.
func (sl *SL) verifyAndInsert(sig Signal, startPos int) error {
	if err := sl.verifyInsert(sig, startPos); err != nil {
		return err
	}

	sl.insert(sig, startPos)

	return nil
}

// delete removes the signal from the signal layout.
func (sl *SL) delete(sig Signal) {
	sl.tree.Delete(sig)

	// Regenerate the filters
	sl.genFilters()
}

// clear removes all signals from the signal layout.
func (sl *SL) clear() {
	sl.tree.Clear()

	// Reset the filters
	sl.filters = []*SignalLayoutFilter{}
}

func (sl *SL) getIntervalFromNewStartPos(sig Signal, newStartPos int) (int, int) {
	return newStartPos, newStartPos + sig.GetSize() - 1
}

// verifyNewStartPos checks if setting the signal to the new start position
// does not intersect with another one.
func (sl *SL) verifyNewStartPos(sig Signal, newStartPos int) error {
	if err := sl.verifyStartPos(newStartPos); err != nil {
		return err
	}

	if err := sl.verifyStartPosPlusSize(newStartPos, sig.GetSize()); err != nil {
		return err
	}

	newLow, newHigh := sl.getIntervalFromNewStartPos(sig, newStartPos)

	if !sl.tree.CanUpdate(sig, newLow, newHigh) {
		return &StartPosError{
			StartPos: newStartPos,
			Err:      ErrIntersects,
		}
	}

	return nil
}

// updateStartPos updates the start position of the signal in the signal layout.
// It must be called after the verify function since it assumes the
// new start position is valid.
func (sl *SL) updateStartPos(sig Signal, newStartPos int) {
	newLow, newHigh := sl.getIntervalFromNewStartPos(sig, newStartPos)
	sl.tree.Update(sig, newLow, newHigh)

	// Regenerate the filters
	sl.genFilters()
}

// verifyAndUpdateStartPos checks and updates the start position of the signal in the signal layout.
func (sl *SL) verifyAndUpdateStartPos(sig Signal, newStartPos int) error {
	if err := sl.verifyNewStartPos(sig, newStartPos); err != nil {
		return err
	}

	sl.updateStartPos(sig, newStartPos)

	return nil
}

func (sl *SL) getIntervalFromNewSize(sig Signal, newSize int) (int, int) {
	startPos := sig.GetRelativeStartPos()
	return startPos, startPos + newSize - 1
}

// verifyNewSize checks if setting the signal to the new size
//  does not intersect with another one.
func (sl *SL) verifyNewSize(sig Signal, newSize int) error {
	if err := sl.verifySize(newSize); err != nil {
		return err
	}

	if err := sl.verifyStartPosPlusSize(sig.GetRelativeStartPos(), newSize); err != nil {
		return err
	}

	newLow, newHigh := sl.getIntervalFromNewSize(sig, newSize)

	if !sl.tree.CanUpdate(sig, newLow, newHigh) {
		return &SignalSizeError{
			Size: newSize,
			Err:  ErrIntersects,
		}
	}

	return nil
}

// updateSize updates the size of the signal in the signal layout.
// It must be called after the verify function since it assumes
// the new size is valid.
func (sl *SL) updateSize(sig Signal, newSize int) {
	newLow, newHigh := sl.getIntervalFromNewSize(sig, newSize)
	sl.tree.Update(sig, newLow, newHigh)

	// Regenerate the filters
	sl.genFilters()
}

// verifyAndUpdateSize checks and updates the size of the signal in the signal layout.
func (sl *SL) verifyAndUpdateSize(sig Signal, newSize int) error {
	if err := sl.verifyNewSize(sig, newSize); err != nil {
		return err
	}

	sl.updateSize(sig, newSize)

	return nil
}

// Signals returns the signals in the signal layout ordered by the start position.
func (sl *SL) Signals() []Signal {
	return sl.tree.GetAllIntervals()
}

// Filters returns the signal filters of the layout.
func (sl *SL) Filters() []*SignalLayoutFilter {
	return sl.filters
}

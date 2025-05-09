package acmelib

import (
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// SignalLayoutFilter represents a filter of a [SignalLayout].
// A signle [Signal] can have multiple filters in a [SignalLayout].
type SignalLayoutFilter struct {
	signal     Signal
	byteIdx    int
	mask       uint8
	length     int
	leftOffset int
}

func (slf *SignalLayoutFilter) stringify(s *stringer.Stringer) {
	s.Write("entity_id: %s; name: %s; byte_index: %d; mask: %08b; length: %d; left_offset: %d\n",
		slf.signal.EntityID(), slf.signal.Name(), slf.byteIdx, slf.mask, slf.length, slf.leftOffset)
}

func (slf *SignalLayoutFilter) String() string {
	s := stringer.New()
	s.Write("signal_layout_filter:\n")
	slf.stringify(s)
	return s.String()
}

// Signal returns the [Signal] of the [SignalLayoutFilter].
func (slf *SignalLayoutFilter) Signal() Signal {
	return slf.signal
}

// ByteIndex returns the byte index in which the mask is located.
func (slf *SignalLayoutFilter) ByteIndex() int {
	return slf.byteIdx
}

// Mask returns the mask used for filtering.
func (slf *SignalLayoutFilter) Mask() uint8 {
	return slf.mask
}

// Length returns the length of the mask.
func (slf *SignalLayoutFilter) Length() int {
	return slf.length
}

// LeftOffset returns the amount of bits that the mask is shifted to the left.
func (slf *SignalLayoutFilter) LeftOffset() int {
	return slf.leftOffset
}

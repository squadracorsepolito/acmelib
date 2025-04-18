package acmelib

import (
	"fmt"
	"strings"
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

func (slf *SignalLayoutFilter) stringify(b *strings.Builder, tabs int) {
	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%sentity_id: %s; name: %s; byte_index: %d; mask: %08b; length: %d; left_offset: %d\n",
		tabStr, slf.signal.EntityID(), slf.signal.Name(), slf.byteIdx, slf.mask, slf.length, slf.leftOffset))
}

func (slf *SignalLayoutFilter) String() string {
	b := new(strings.Builder)
	slf.stringify(b, 0)
	return b.String()
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

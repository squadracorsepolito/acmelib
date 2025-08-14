package acmelib

import (
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

var _ Signal = (*MuxorSignal)(nil)

// MuxorSignal represents a multiplexor signal.
// It cannot be directly created since it is created when
// a [MultiplexedLayer] is added to a [SignalLayout].
type MuxorSignal struct {
	*signal

	layoutCount int
}

func newMuxorSignalFromBase(base *signal, layoutCount int) (*MuxorSignal, error) {
	if layoutCount < 0 {
		return nil, newArgError("layoutCount", ErrIsNegative)
	}

	if layoutCount == 0 {
		return nil, newArgError("layoutCount", ErrIsZero)
	}

	ms := &MuxorSignal{
		signal: base,

		layoutCount: layoutCount,
	}

	ms.signal.setSize(getSizeFromCount(layoutCount))

	return ms, nil
}

// NewMuxorSignal creates a new muxor signal with the given name and layout count.
func NewMuxorSignal(name string, layoutCount int) (*MuxorSignal, error) {
	return newMuxorSignalFromBase(newSignal(name, SignalKindMuxor), layoutCount)
}

func (ms *MuxorSignal) stringify(s *stringer.Stringer) {
	ms.signal.stringify(s)
	s.Write("layout_count: %d\n", ms.layoutCount)
}

func (ms *MuxorSignal) String() string {
	s := stringer.New()
	s.Write("muxor_signal:\n")
	ms.stringify(s)
	return s.String()
}

// ToMuxor returns the signal as a muxor signal.
func (ms *MuxorSignal) ToMuxor() (*MuxorSignal, error) {
	return ms, nil
}

// UpdateLayoutCount updates the layout count of the muxor signal.
//
// It returns:
//   - [LayoutIDError] if the new layout count is invalid.
//   - [SizeError] if the size of the muxor signal becomes invalid with
//     the new layout count.
func (ms *MuxorSignal) UpdateLayoutCount(newLayoutCount int) error {
	if newLayoutCount < 0 {
		return ms.errorf(newLayoutIDError(newLayoutCount, ErrIsNegative))
	}

	if newLayoutCount == ms.layoutCount {
		return nil
	}

	if err := ms.verifyAndUpdateSize(ms, getSizeFromCount(newLayoutCount)); err != nil {
		return ms.errorf(err)
	}

	if newLayoutCount < ms.layoutCount {
		// Check if the multiplexed layer the muxor signal is in can
		// redude the number of signal layouts
		for lID, muxLayout := range ms.parentMuxLayer.iterLayouts() {
			if lID < newLayoutCount {
				continue
			}

			if muxLayout.SignalCount() != 0 {
				return ms.errorf(newLayoutIDError(lID, ErrNotClear))
			}
		}

		// Delete the layouts in excess
		ms.parentMuxLayer.truncateLayouts(newLayoutCount)

	} else {
		// Add the layouts to reach the new layout count
		ms.parentMuxLayer.appendLayouts(newLayoutCount - ms.layoutCount)
	}

	ms.layoutCount = newLayoutCount

	return nil
}

// UpdateStartPos updates the start position of the signal.
//
// It returns a [StartPosError] if the new start position is invalid.
func (ms *MuxorSignal) UpdateStartPos(newStartPos int) error {
	return ms.signal.updateStartPos(ms, newStartPos)
}

// UpdateEncodedValue updates the current physical value of the signal.
// The value represents the layout id.
//
// It returns an [ArgError] if the given value is negative or out of bounds.
func (ms *MuxorSignal) UpdateEncodedValue(value int) error {
	if value < 0 {
		return ms.errorf(newArgError("value", ErrIsNegative))
	}

	if value >= ms.layoutCount {
		return ms.errorf(newArgError("value", ErrOutOfBounds))
	}

	ms.encodedValue = uint64(value)

	return nil
}

// ToSignal returns the signal itself.
func (ms *MuxorSignal) ToSignal() (Signal, error) {
	return ms, nil
}

// AssignAttribute assigns the given attribute/value pair to the signal.
func (ms *MuxorSignal) AssignAttribute(attribute Attribute, value any) error {
	return ms.signal.assignAttribute(ms, attribute, value)
}

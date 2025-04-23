package acmelib

import "strings"

var _ Signal = (*MuxorSignal)(nil)

type MuxorSignal struct {
	*signal

	layoutCount int
}

func newMuxorSignal(name string, layoutCount int) *MuxorSignal {
	ms := &MuxorSignal{
		signal: newSignal(name, SignalKindMuxor),

		layoutCount: layoutCount,
	}

	ms.signal.setSize(calcSizeFromValue(layoutCount))

	return ms
}

func (ms *MuxorSignal) String() string {
	b := new(strings.Builder)
	ms.stringify(b, 0)
	return b.String()
}

func (ms *MuxorSignal) ToMuxor() (*MuxorSignal, error) {
	return ms, nil
}

func (ms *MuxorSignal) AssignAttribute(attribute Attribute, value any) error {
	if err := ms.addAttributeAssignment(attribute, ms, value); err != nil {
		return ms.errorf(err)
	}
	return nil
}

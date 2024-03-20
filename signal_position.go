package acmelib

type SignalPosition struct {
	From int
	To   int
}

func NewSignalPosition(from, to int) *SignalPosition {
	return &SignalPosition{
		From: from,
		To:   to,
	}
}

func (sp *SignalPosition) Size() int {
	return sp.To - sp.From
}

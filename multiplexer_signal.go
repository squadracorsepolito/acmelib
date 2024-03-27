package acmelib

import (
	"errors"
	"fmt"
	"strings"
)

type MultiplexerSignal struct {
	*signal

	muxSignals         map[EntityID]Signal
	muxSignalNames     map[string]EntityID
	muxSignalSelValues map[EntityID]int

	signalPayloads map[int]*signalPayload

	totalSize  int
	selectSize int
}

func NewMultiplexerSignal(name, desc string, totalSize, selectSize int) (*MultiplexerSignal, error) {
	ms := &MultiplexerSignal{
		signal: newSignal(name, desc, SignalKindMultiplexer),

		muxSignals:         make(map[EntityID]Signal),
		muxSignalNames:     make(map[string]EntityID),
		muxSignalSelValues: make(map[EntityID]int),

		signalPayloads: make(map[int]*signalPayload),

		totalSize:  totalSize,
		selectSize: selectSize,
	}

	return ms, nil
}

func (ms *MultiplexerSignal) String() string {
	var builder strings.Builder

	builder.WriteString(ms.signal.String())
	builder.WriteString(fmt.Sprintf("size: %d\n", ms.GetSize()))

	builder.WriteString("\n+++END SIGNAL+++\n\n")

	return builder.String()
}

func (ms *MultiplexerSignal) GetSize() int {
	return ms.totalSize
}

func (ms *MultiplexerSignal) ToStandard() (*StandardSignal, error) {
	return nil, ms.errorf(fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindStandard, SignalKindMultiplexer))
}

func (ms *MultiplexerSignal) ToEnum() (*EnumSignal, error) {
	return nil, ms.errorf(fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindEnum, SignalKindMultiplexer))
}

func (ms *MultiplexerSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return ms, nil
}

func (ms *MultiplexerSignal) verifySelectValue(selVal int) error {
	if selVal < 0 {
		return errors.New("select value cannot be negative")
	}

	if calcSizeFromValue(selVal) > ms.selectSize {
		return fmt.Errorf(`select value "%d" size exceeds the max select value size ("%d")`, selVal, ms.selectSize)
	}

	return nil
}

func (ms *MultiplexerSignal) addSignalPayload(selVal int) *signalPayload {
	payload := newSignalPayload(ms.totalSize - ms.selectSize)
	ms.signalPayloads[selVal] = payload
	return payload
}

func (ms *MultiplexerSignal) getSignalPayload(selVal int) *signalPayload {
	if payload, ok := ms.signalPayloads[selVal]; ok {
		return payload
	}
	return nil
}

func (ms *MultiplexerSignal) getMuxSignal(sigID EntityID) Signal {
	if muxSig, ok := ms.muxSignals[sigID]; ok {
		return muxSig
	}
	return nil
}

func (ms *MultiplexerSignal) addMuxSignal(selValue int, sig Signal) {
	id := sig.EntityID()

	ms.muxSignals[id] = sig
	ms.muxSignalNames[sig.Name()] = id
	ms.muxSignalSelValues[id] = selValue

	sig.setParent(ms)
}

func (ms *MultiplexerSignal) AppendMuxSignal(selectValue int, signal Signal) error {
	if err := ms.verifySignalName(signal.Name()); err != nil {
		return ms.errorf(err)
	}

	if err := ms.verifySelectValue(selectValue); err != nil {
		return ms.errorf(err)
	}

	payload := ms.getSignalPayload(selectValue)
	if payload == nil {
		payload = ms.addSignalPayload(selectValue)
	}

	if err := payload.append(signal); err != nil {
		return ms.errorf(err)
	}

	ms.addMuxSignal(selectValue, signal)

	return nil
}

func (ms *MultiplexerSignal) InsertMuxSignal(selectValue int, signal Signal, startBit int) error {
	if err := ms.verifySignalName(signal.Name()); err != nil {
		return ms.errorf(err)
	}

	if err := ms.verifySelectValue(selectValue); err != nil {
		return ms.errorf(err)
	}

	payload := ms.getSignalPayload(selectValue)
	if payload == nil {
		payload = ms.addSignalPayload(selectValue)
	}

	if err := payload.insert(signal, startBit); err != nil {
		return ms.errorf(err)
	}

	ms.addMuxSignal(selectValue, signal)

	return nil
}

func (ms *MultiplexerSignal) GetSelectedMuxSignals(selectValue int) []Signal {
	payload := ms.getSignalPayload(selectValue)

	if payload != nil {
		return payload.signals
	}

	return []Signal{}
}

func (ms *MultiplexerSignal) getSignalParentKind() signalParentKind {
	return signalParentKindMultiplexerSignal
}

func (ms *MultiplexerSignal) modifySignalName(sigID EntityID, newName string) {
}

func (ms *MultiplexerSignal) toParentMessage() (*Message, error) {
	return nil, fmt.Errorf(`cannot convert to "%s" signal parent is of kind "%s"`,
		signalParentKindMessage, signalParentKindMultiplexerSignal)
}

func (ms *MultiplexerSignal) toParentMultiplexerSignal() (*MultiplexerSignal, error) {
	return ms, nil
}

func (ms *MultiplexerSignal) verifySignalName(name string) error {
	if _, ok := ms.muxSignalNames[name]; ok {
		return fmt.Errorf(`signal name "%s" is duplicated`, name)
	}

	if ms.hasParent() {
		return ms.parent.verifySignalName(name)
	}

	return nil
}

func (ms *MultiplexerSignal) modifySignalSize(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	sig := ms.getMuxSignal(sigID)
	if sig == nil {
		return fmt.Errorf(`multiplexed signal with id "%s" not found`, sigID)
	}

	selVal, ok := ms.muxSignalSelValues[sig.EntityID()]
	if !ok {
		return fmt.Errorf(`select value for signal "%s" not found`, sig.Name())
	}

	payload := ms.getSignalPayload(selVal)

	if amount > 0 {
		return payload.modifyStartBitsOnGrow(sig, amount)
	}

	payload.modifyStartBitsOnShrink(sig, -amount)

	return nil
}

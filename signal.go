package acmelib

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type SignalKind string

const (
	SignalKindStandard    SignalKind = "standard"
	SignalKindEnum        SignalKind = "enum"
	SignalKindMultiplexer SignalKind = "multiplexer"
)

type signalParentKind string

const (
	signalParentKindMessage           signalParentKind = "message"
	signalParentKindMultiplexerSignal signalParentKind = "multiplexer_signal"
)

type SignalParent interface {
	errorf(err error) error

	GetSignalParentKind() signalParentKind

	verifySignalName(name string) error
	modifySignalName(sigID EntityID, newName string) error

	verifySignalSizeAmount(sigID EntityID, amount int) error
	modifySignalSize(sigID EntityID, amount int) error

	ToParentMessage() (*Message, error)
	ToParentMultiplexerSignal() (*MultiplexerSignal, error)
}

type Signal interface {
	EntityID() EntityID
	Name() string
	Desc() string
	CreateTime() time.Time

	String() string

	Kind() SignalKind

	Parent() SignalParent
	setParent(parent SignalParent)

	GetStartBit() int
	getRelStartBit() int
	setRelStartBit(startBit int)

	GetSize() int

	ToStandard() (*StandardSignal, error)
	ToEnum() (*EnumSignal, error)
	ToMultiplexer() (*MultiplexerSignal, error)
}

type signal struct {
	*entity

	parent SignalParent

	kind        SignalKind
	relStartBit int
}

func newSignal(name, desc string, kind SignalKind) *signal {
	return &signal{
		entity: newEntity(name, desc),

		parent: nil,

		kind:        kind,
		relStartBit: 0,
	}
}

func (s *signal) hasParent() bool {
	return s.parent != nil
}

func (s *signal) modifySize(amount int) error {
	if s.hasParent() {
		return s.parent.modifySignalSize(s.EntityID(), amount)
	}

	return nil
}

func (s *signal) errorf(err error) error {
	sigErr := fmt.Errorf(`signal "%s": %w`, s.name, err)
	if s.hasParent() {
		return s.parent.errorf(sigErr)
	}
	return sigErr
}

func (s *signal) Kind() SignalKind {
	return s.kind
}

func (s *signal) Parent() SignalParent {
	return s.parent
}

func (s *signal) setParent(parent SignalParent) {
	s.parent = parent
}

func (s *signal) getRelStartBit() int {
	return s.relStartBit
}

func (s *signal) setRelStartBit(startBit int) {
	s.relStartBit = startBit
}

func (s *signal) String() string {
	var builder strings.Builder

	builder.WriteString("\n+++START SIGNAL+++\n\n")
	builder.WriteString(s.toString())
	builder.WriteString(fmt.Sprintf("kind: %s\n", s.kind))
	builder.WriteString(fmt.Sprintf("start_bit: %d; ", s.relStartBit))

	return builder.String()
}

func (s *signal) GetStartBit() int {
	if s.hasParent() {
		if s.parent.GetSignalParentKind() == signalParentKindMultiplexerSignal {
			muxParent, err := s.parent.ToParentMultiplexerSignal()
			if err != nil {
				panic(err)
			}

			return muxParent.GetStartBit() + muxParent.SelectSize() + s.relStartBit
		}
	}

	return s.relStartBit
}

func (s *signal) UpdateName(newName string) error {
	if s.name == newName {
		return nil
	}

	if s.hasParent() {
		if err := s.parent.verifySignalName(newName); err != nil {
			return s.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}

		if err := s.parent.modifySignalName(s.entityID, newName); err != nil {
			return s.errorf(fmt.Errorf(`cannot update name to "%s" : %w`, newName, err))
		}
	}

	s.name = newName

	return nil
}

// -----------------------
// +++ STANDARD SIGNAL +++
// -----------------------

type StandardSignal struct {
	*signal

	typ    *SignalType
	min    float64
	max    float64
	offset float64
	scale  float64
	unit   *SignalUnit
}

func NewStandardSignal(name, desc string, typ *SignalType, min, max, offset, scale float64, unit *SignalUnit) (*StandardSignal, error) {
	if typ == nil {
		return nil, errors.New("signal type cannot be nil")
	}

	return &StandardSignal{
		signal: newSignal(name, desc, SignalKindStandard),

		typ:    typ,
		min:    min,
		max:    max,
		offset: offset,
		scale:  scale,
		unit:   unit,
	}, nil
}

// ---------------------------------------
// +++ Signal interface implementation +++
// ---------------------------------------

func (ss *StandardSignal) GetSize() int {
	return ss.typ.size
}

func (ss *StandardSignal) ToStandard() (*StandardSignal, error) {
	return ss, nil
}

func (ss *StandardSignal) ToEnum() (*EnumSignal, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindEnum, SignalKindStandard)
}

func (ss *StandardSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindMultiplexer, SignalKindStandard)
}

// ----------------------
// +++ public methods +++
// ----------------------

func (ss *StandardSignal) String() string {
	var builder strings.Builder

	builder.WriteString(ss.signal.String())
	builder.WriteString(fmt.Sprintf("size: %d\n", ss.GetSize()))
	builder.WriteString(fmt.Sprintf("min: %f; max: %f; offset: %f; scale: %f\n", ss.min, ss.max, ss.offset, ss.scale))

	builder.WriteString("\n+++END SIGNAL+++\n\n")

	return builder.String()
}

func (ss *StandardSignal) Type() *SignalType {
	return ss.typ
}

func (ss *StandardSignal) SetType(typ *SignalType, min, max, offset, scale float64) error {
	if typ == nil {
		return errors.New("signal type cannot be nil")
	}

	if err := ss.modifySize(typ.size - ss.typ.size); err != nil {
		return ss.errorf(err)
	}

	ss.typ = typ
	ss.min = min
	ss.max = max
	ss.offset = offset
	ss.scale = scale

	return nil
}

func (ss *StandardSignal) Min() float64 {
	return ss.min
}

func (ss *StandardSignal) Max() float64 {
	return ss.max
}

func (ss *StandardSignal) Offset() float64 {
	return ss.offset
}

func (ss *StandardSignal) Scale() float64 {
	return ss.scale
}

func (ss *StandardSignal) Unit() *SignalUnit {
	return ss.unit
}

func (ss *StandardSignal) SetUnit(unit *SignalUnit) {
	ss.unit = unit
}

// -------------------
// +++ ENUM SIGNAL +++
// -------------------

type EnumSignal struct {
	*signal

	enum *SignalEnum
}

func NewEnumSignal(name, desc string, enum *SignalEnum) (*EnumSignal, error) {
	if enum == nil {
		return nil, errors.New("signal enum cannot be nil")
	}

	sig := &EnumSignal{
		signal: newSignal(name, desc, SignalKindEnum),

		enum: enum,
	}

	enum.addSignalRef(sig)

	return sig, nil
}

// ---------------------------------------
// +++ Signal interface implementation +++
// ---------------------------------------

func (es *EnumSignal) GetSize() int {
	return es.enum.GetSize()
}

func (es *EnumSignal) ToStandard() (*StandardSignal, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindStandard, SignalKindEnum)
}

func (es *EnumSignal) ToEnum() (*EnumSignal, error) {
	return es, nil
}

func (es *EnumSignal) ToMultiplexer() (*MultiplexerSignal, error) {
	return nil, fmt.Errorf(`cannot covert to "%s", the signal is of kind "%s"`, SignalKindMultiplexer, SignalKindEnum)
}

// ----------------------
// +++ public methods +++
// ----------------------

func (es *EnumSignal) String() string {
	var builder strings.Builder

	builder.WriteString(es.signal.String())
	builder.WriteString(fmt.Sprintf("size: %d\n", es.GetSize()))

	builder.WriteString("\n+++END SIGNAL+++\n\n")

	return builder.String()
}

func (es *EnumSignal) Enum() *SignalEnum {
	return es.enum
}

func (es *EnumSignal) SetEnum(enum *SignalEnum) error {
	if enum == nil {
		return errors.New("signal enum cannot be nil")
	}

	if err := es.modifySize(enum.GetSize() - es.GetSize()); err != nil {
		return es.errorf(err)
	}

	es.enum = enum

	es.enum.removeSignalRef(es.EntityID())
	enum.addSignalRef(es)

	return nil
}

// --------------------------
// +++ MULTIPLEXER SIGNAL +++
// --------------------------

type MultiplexerSignal struct {
	*signal

	muxSignals         map[EntityID]Signal
	muxSignalNames     map[string]EntityID
	muxSignalSelValues map[EntityID]int

	signalPayloads map[int]*signalPayload

	selValRanges map[int]int

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

		selValRanges: make(map[int]int),

		totalSize:  totalSize,
		selectSize: selectSize,
	}

	return ms, nil
}

func (ms *MultiplexerSignal) addSignalPayload(selVal int) *signalPayload {
	payload := newSignalPayload(ms.totalSize - ms.selectSize)
	ms.signalPayloads[selVal] = payload
	return payload
}

func (ms *MultiplexerSignal) getSignalPayload(selVal int) (*signalPayload, int) {
	tmpSelVal := selVal

	if len(ms.selValRanges) > 0 {
		if rangeSelVal, ok := ms.selValRanges[selVal]; ok {
			tmpSelVal = rangeSelVal
		}
	}

	if payload, ok := ms.signalPayloads[tmpSelVal]; ok {
		return payload, tmpSelVal
	}

	return nil, tmpSelVal
}

func (ms *MultiplexerSignal) addMuxSignalName(sigID EntityID, name string) {
	ms.muxSignalNames[name] = sigID

	if ms.hasParent() {
		parent := ms.Parent()
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case signalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case signalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				msgParent.addSignalName(sigID, name)
				return
			}
		}
	}
}

func (ms *MultiplexerSignal) removeMuxSignalName(name string) {
	delete(ms.muxSignalNames, name)

	if ms.hasParent() {
		parent := ms.Parent()
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case signalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case signalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				msgParent.removeSignalName(name)
				return
			}
		}
	}
}

func (ms *MultiplexerSignal) addMuxSignal(selValue int, sig Signal) {
	id := sig.EntityID()

	ms.muxSignals[id] = sig
	ms.addMuxSignalName(id, sig.Name())
	ms.muxSignalSelValues[id] = selValue

	if ms.hasParent() {
		parent := ms.Parent()
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case signalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case signalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				msgParent.addSignal(sig)
				return
			}
		}
	}

}

func (ms *MultiplexerSignal) removeMuxSignal(sigID EntityID) {
	delete(ms.muxSignalSelValues, sigID)
	delete(ms.muxSignals, sigID)

	if ms.hasParent() {
		parent := ms.Parent()
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case signalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case signalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				msgParent.removeSignal(sigID)
				return
			}
		}
	}
}

func (ms *MultiplexerSignal) getMuxSignalByID(sigID EntityID) (Signal, error) {
	if muxSig, ok := ms.muxSignals[sigID]; ok {
		return muxSig, nil
	}
	return nil, fmt.Errorf("multiplexed signal not found")
}

func (ms *MultiplexerSignal) getMuxSignalSelValue(sigID EntityID) (int, error) {
	if selVal, ok := ms.muxSignalSelValues[sigID]; ok {
		return selVal, nil
	}
	return -1, fmt.Errorf(`select value for multiplexed signal with id "%s" not found`, sigID)
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

// ---------------------------------------------
// +++ signalParent interface implementation +++
// ---------------------------------------------

func (ms *MultiplexerSignal) GetSignalParentKind() signalParentKind {
	return signalParentKindMultiplexerSignal
}

func (ms *MultiplexerSignal) modifySignalName(sigID EntityID, newName string) error {
	if ms.hasParent() {
		parent := ms.Parent()

	loop:
		for parent != nil {
			switch parent.GetSignalParentKind() {
			case signalParentKindMultiplexerSignal:
				muxParent, err := parent.ToParentMultiplexerSignal()
				if err != nil {
					panic(err)
				}
				parent = muxParent.Parent()

			case signalParentKindMessage:
				msgParent, err := parent.ToParentMessage()
				if err != nil {
					panic(err)
				}

				if err := msgParent.modifySignalName(sigID, newName); err != nil {
					return err
				}
				break loop
			}
		}
	}

	sig, err := ms.getMuxSignalByID(sigID)
	if err != nil {
		return err
	}

	oldName := sig.Name()

	ms.removeMuxSignalName(oldName)
	ms.addMuxSignalName(sigID, newName)

	return nil
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

func (ms *MultiplexerSignal) verifySignalSizeAmount(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	sig, err := ms.getMuxSignalByID(sigID)
	if err != nil {
		return err
	}

	selVal, err := ms.getMuxSignalSelValue(sigID)
	if err != nil {
		return err
	}

	payload, _ := ms.getSignalPayload(selVal)

	if amount > 0 {
		return payload.verifyBeforeGrow(sig, amount)
	}

	return payload.verifyBeforeShrink(sig, -amount)
}

func (ms *MultiplexerSignal) modifySignalSize(sigID EntityID, amount int) error {
	if amount == 0 {
		return nil
	}

	sig, err := ms.getMuxSignalByID(sigID)
	if err != nil {
		return err
	}

	selVal, err := ms.getMuxSignalSelValue(sigID)
	if err != nil {
		return err
	}

	payload, _ := ms.getSignalPayload(selVal)

	if amount > 0 {
		return payload.modifyStartBitsOnGrow(sig, amount)
	}

	return payload.modifyStartBitsOnShrink(sig, -amount)
}

func (ms *MultiplexerSignal) ToParentMessage() (*Message, error) {
	return nil, fmt.Errorf(`cannot convert to "%s" signal parent is of kind "%s"`,
		signalParentKindMessage, signalParentKindMultiplexerSignal)
}

func (ms *MultiplexerSignal) ToParentMultiplexerSignal() (*MultiplexerSignal, error) {
	return ms, nil
}

// ---------------------------------------
// +++ Signal interface implementation +++
// ---------------------------------------

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

// ----------------------
// +++ public methods +++
// ----------------------

func (ms *MultiplexerSignal) String() string {
	var builder strings.Builder

	builder.WriteString(ms.signal.String())
	builder.WriteString(fmt.Sprintf("size: %d\n", ms.GetSize()))

	builder.WriteString("\n+++END SIGNAL+++\n\n")

	return builder.String()
}

func (ms *MultiplexerSignal) GetSelectedMuxSignals(selectValue int) []Signal {
	payload, _ := ms.getSignalPayload(selectValue)

	if payload != nil {
		return payload.signals
	}

	return []Signal{}
}

func (ms *MultiplexerSignal) MuxSignals() map[int][]Signal {
	res := make(map[int][]Signal)

	for selVal, payload := range ms.signalPayloads {
		res[selVal] = payload.signals
	}

	return res
}

func (ms *MultiplexerSignal) AppendMuxSignal(selectValue int, signal Signal) error {
	if err := ms.verifySignalName(signal.Name()); err != nil {
		return ms.errorf(err)
	}

	if err := ms.verifySelectValue(selectValue); err != nil {
		return ms.errorf(err)
	}

	payload, realSelVal := ms.getSignalPayload(selectValue)
	if payload == nil {
		payload = ms.addSignalPayload(realSelVal)
	}

	if err := payload.append(signal); err != nil {
		return ms.errorf(err)
	}

	ms.addMuxSignal(realSelVal, signal)

	signal.setParent(ms)

	return nil
}

func (ms *MultiplexerSignal) InsertMuxSignal(selectValue int, signal Signal, startBit int) error {
	if err := ms.verifySignalName(signal.Name()); err != nil {
		return ms.errorf(err)
	}

	if err := ms.verifySelectValue(selectValue); err != nil {
		return ms.errorf(err)
	}

	payload, realSelVal := ms.getSignalPayload(selectValue)
	if payload == nil {
		payload = ms.addSignalPayload(realSelVal)
	}

	if err := payload.insert(signal, startBit); err != nil {
		return ms.errorf(err)
	}

	ms.addMuxSignal(realSelVal, signal)

	signal.setParent(ms)

	return nil
}

func (ms *MultiplexerSignal) ShiftMuxSignalLeft(muxSignalEntityID EntityID, amount int) int {
	selVal, err := ms.getMuxSignalSelValue(muxSignalEntityID)
	if err != nil {
		return 0
	}

	sig, err := ms.getMuxSignalByID(muxSignalEntityID)
	if err != nil {
		return 0
	}

	payload, _ := ms.getSignalPayload(selVal)
	if payload == nil {
		return 0
	}

	return payload.shiftLeft(sig, amount)
}

func (ms *MultiplexerSignal) ShiftMuxSignalRight(muxSignalEntityID EntityID, amount int) int {
	selVal, err := ms.getMuxSignalSelValue(muxSignalEntityID)
	if err != nil {
		return 0
	}

	sig, err := ms.getMuxSignalByID(muxSignalEntityID)
	if err != nil {
		return 0
	}

	payload, _ := ms.getSignalPayload(selVal)
	if payload == nil {
		return 0
	}

	return payload.shiftRight(sig, amount)
}

func (ms *MultiplexerSignal) RemoveMuxSignal(muxSignalEntityID EntityID) error {
	selVal, err := ms.getMuxSignalSelValue(muxSignalEntityID)
	if err != nil {
		return ms.errorf(fmt.Errorf(`cannot remove mux signal with id "%s" : %w`, muxSignalEntityID, err))
	}

	sig, err := ms.getMuxSignalByID(muxSignalEntityID)
	if err != nil {
		return ms.errorf(fmt.Errorf(`cannot remove mux signal with id "%s" : %w`, muxSignalEntityID, err))
	}

	payload, _ := ms.getSignalPayload(selVal)
	if payload == nil {
		return nil
	}

	ms.removeMuxSignal(muxSignalEntityID)
	ms.removeMuxSignalName(sig.Name())

	payload.remove(muxSignalEntityID)

	sig.setParent(nil)

	return nil
}

func (ms *MultiplexerSignal) RemoveAllMuxSignals() {
	for muxSigID, tmpMuxSig := range ms.muxSignals {
		tmpMuxSig.setParent(nil)
		ms.removeMuxSignalName(tmpMuxSig.Name())
		ms.removeMuxSignal(muxSigID)
	}

	for _, payload := range ms.signalPayloads {
		payload.removeAll()
	}
}

func (ms *MultiplexerSignal) SelectSize() int {
	return ms.selectSize
}

func (ms *MultiplexerSignal) AddSelectValueRange(from, to int) error {
	if from > to {
		return ms.errorf(fmt.Errorf(`cannot set select value range because from "%d" is greater then to "%d"`, from, to))
	}

	if err := ms.verifySelectValue(from); err != nil {
		return ms.errorf(fmt.Errorf(`cannot set select value range : %w`, err))
	}

	if err := ms.verifySelectValue(to); err != nil {
		return ms.errorf(fmt.Errorf(`cannot set select value range : %w`, err))
	}

	foundOne := false
	foundSelVal := from
	for i := from; i <= to; i++ {
		if _, ok := ms.signalPayloads[i]; ok {
			if foundOne {
				return ms.errorf(fmt.Errorf(`cannot set select value range because there are more than 1 payloads between "%d" an "%d"`, from, to))
			}

			foundSelVal = i
			foundOne = true
		}

		if _, ok := ms.selValRanges[i]; ok {
			return ms.errorf(fmt.Errorf(`cannot set select value range because value "%d" is already used in another range`, i))
		}
	}

	for i := from; i <= to; i++ {
		ms.selValRanges[i] = foundSelVal
	}

	return nil
}

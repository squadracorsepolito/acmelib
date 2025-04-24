package acmelib

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// SignalLayout represents a layout of signals.
// It can be generated from a [Message] or a [MultiplexerSignal] (TODO!).
type SignalLayout struct {
	size    int
	signals []Signal
	filters []*SignalLayoutFilter
}

func newSignalLayout(size int) *SignalLayout {
	return &SignalLayout{
		size:    size,
		signals: []Signal{},
		filters: []*SignalLayoutFilter{},
	}
}

// generateFilters generates the filters of the layout.
// It must be called every time the layout is changed.
func (sl *SignalLayout) generateFilters() {
	filters := []*SignalLayoutFilter{}

	for _, sig := range sl.signals {
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

			filters = append(filters, &SignalLayoutFilter{
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

				if endianness == EndiannessBigEndian {
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

			if endianness == EndiannessBigEndian {
				// For the same reason as above, shift the mask left
				leftOffset = 8 - remainingBits
				mask <<= leftOffset
			}

		appendFilter:
			filters = append(filters, &SignalLayoutFilter{
				signal:     sig,
				byteIdx:    i,
				mask:       uint8(mask),
				length:     length,
				leftOffset: leftOffset,
			})

			remainingBits -= length
		}
	}

	sl.filters = filters
}

func (sl *SignalLayout) verifyBeforeAppend(sig Signal) error {
	sigSize := sig.GetSize()
	sigCount := len(sl.signals)
	if sigCount == 0 {
		if sigSize > sl.size {
			return &SignalSizeError{
				Size: sigSize,
				Err:  ErrOutOfBounds,
			}
		}

		return nil
	}

	lastSig := sl.signals[sigCount-1]
	trailingSpace := sl.size - (lastSig.GetRelativeStartPos() + lastSig.GetSize())

	if sigSize > trailingSpace {
		return &SignalSizeError{
			Size: sigSize,
			Err:  ErrNoSpaceLeft,
		}
	}

	return nil
}

func (sl *SignalLayout) append(sig Signal) error {
	if err := sl.verifyBeforeAppend(sig); err != nil {
		return &AppendSignalError{
			EntityID: sig.EntityID(),
			Name:     sig.Name(),
			Err:      err,
		}
	}

	if len(sl.signals) == 0 {
		sig.setRelativeStartPos(0)
	} else {
		lastSig := sl.signals[len(sl.signals)-1]
		sig.setRelativeStartPos(lastSig.GetRelativeStartPos() + lastSig.GetSize())
	}

	sl.signals = append(sl.signals, sig)

	sl.generateFilters()

	return nil
}

func (sl *SignalLayout) verifyBeforeInsert(sig Signal, startBit int) error {
	if startBit < 0 {
		return &StartPosError{
			StartPos: startBit,
			Err:      ErrIsNegative,
		}
	}

	sigSize := sig.GetSize()
	endBit := startBit + sigSize

	if sigSize > sl.size {
		return &SignalSizeError{
			Size: sigSize,
			Err:  ErrOutOfBounds,
		}
	}

	if endBit > sl.size {
		return &SignalSizeError{
			Size: sigSize,
			Err:  ErrNoSpaceLeft,
		}
	}

	for _, tmpSig := range sl.signals {
		tmpStartBit := tmpSig.GetRelativeStartPos()
		tmpEndBit := tmpStartBit + tmpSig.GetSize()

		if endBit <= tmpStartBit {
			break
		}

		if startBit >= tmpEndBit {
			continue
		}

		if startBit >= tmpStartBit || endBit > tmpStartBit {
			return &StartPosError{
				StartPos: startBit,
				Err:      ErrIntersects,
			}
		}
	}

	return nil
}

func (sl *SignalLayout) insert(sig Signal, startBit int) {
	if len(sl.signals) == 0 {
		sig.setRelativeStartPos(startBit)
		sl.signals = append(sl.signals, sig)
		sl.generateFilters()
		return
	}

	inserted := false
	for idx, tmpSig := range sl.signals {
		tmpStartBit := tmpSig.GetRelativeStartPos()

		if tmpStartBit > startBit {
			inserted = true
			sl.signals = slices.Insert(sl.signals, idx, sig)
			break
		}
	}

	if !inserted {
		sl.signals = append(sl.signals, sig)
	}

	sig.setRelativeStartPos(startBit)

	sl.generateFilters()
}

func (sl *SignalLayout) verifyAndInsert(sig Signal, startBit int) error {
	if err := sl.verifyBeforeInsert(sig, startBit); err != nil {
		return &InsertSignalError{
			EntityID: sig.EntityID(),
			Name:     sig.Name(),
			StartBit: startBit,
			Err:      err,
		}
	}

	sl.insert(sig, startBit)

	return nil
}

func (sl *SignalLayout) remove(sigID EntityID) {
	sl.signals = slices.DeleteFunc(sl.signals, func(s Signal) bool { return s.EntityID() == sigID })
	sl.generateFilters()
}

func (sl *SignalLayout) removeAll() {
	sl.signals = []Signal{}
	sl.filters = []*SignalLayoutFilter{}
}

func (sl *SignalLayout) compact() {
	lastStartBit := 0
	for _, sig := range sl.signals {
		tmpStartBit := sig.GetRelativeStartPos()

		if tmpStartBit == lastStartBit {
			lastStartBit += sig.GetSize()
			continue
		}

		if lastStartBit < tmpStartBit {
			sig.setRelativeStartPos(lastStartBit)
			lastStartBit += sig.GetSize()
		}
	}

	sl.generateFilters()
}

func (sl *SignalLayout) verifyBeforeShrink(sig Signal, amount int) error {
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

func (sl *SignalLayout) modifyStartBitsOnShrink(sig Signal, amount int) error {
	if amount == 0 {
		return nil
	}

	if err := sl.verifyBeforeShrink(sig, amount); err != nil {
		return &SignalSizeError{
			Size: sig.GetSize() - amount,
			Err:  err,
		}
	}

	found := false
	for _, tmpSig := range sl.signals {
		if found {
			tmpSig.setRelativeStartPos(tmpSig.GetRelativeStartPos() - amount)
			continue
		}

		if sig.EntityID() == tmpSig.EntityID() {
			found = true
		}
	}

	sl.generateFilters()

	return nil
}

func (sl *SignalLayout) verifyBeforeGrow(sig Signal, amount int) error {
	if amount < 0 {
		return ErrIsNegative
	}

	availableSpace := 0
	prevEndBit := 0
	found := false

	for _, tmpSig := range sl.signals {
		tmpStartBit := tmpSig.GetRelativeStartPos()

		if found {
			availableSpace += tmpStartBit - prevEndBit
		} else if tmpSig.EntityID() == sig.EntityID() {
			found = true
		}

		prevEndBit = tmpStartBit + tmpSig.GetSize()
	}

	availableSpace += sl.size - prevEndBit

	if amount > availableSpace {
		return ErrNoSpaceLeft
	}

	return nil
}

func (sl *SignalLayout) modifyStartBitsOnGrow(sig Signal, amount int) error {
	if amount == 0 {
		return nil
	}

	if err := sl.verifyBeforeGrow(sig, amount); err != nil {
		return &SignalSizeError{
			Size: sig.GetSize() + amount,
			Err:  err,
		}
	}

	prevEndBit := 0
	spaces := []int{}
	nextSigIdx := 0
	found := false

	for idx, tmpSig := range sl.signals {
		tmpStartBit := tmpSig.GetRelativeStartPos()

		if found {
			space := tmpStartBit - prevEndBit
			spaces = append(spaces, space)

		} else if sig.EntityID() == tmpSig.EntityID() {
			if idx == len(sl.signals)-1 {
				return nil
			}

			found = true
			nextSigIdx = idx + 1
		}

		prevEndBit = tmpStartBit + tmpSig.GetSize()
	}

	spaces = append(spaces, sl.size-prevEndBit)

	spaceIdx := 0
	acc := amount
	for i := nextSigIdx; i < len(sl.signals); i++ {
		tmpSpace := spaces[spaceIdx]

		if tmpSpace >= acc {
			break
		}

		acc -= tmpSpace
		tmpSig := sl.signals[i]
		tmpSig.setRelativeStartPos(tmpSig.GetRelativeStartPos() + acc)
		spaceIdx++
	}

	sl.generateFilters()

	return nil
}

func (sl *SignalLayout) verifyBeforeResize(newSize int) error {
	if newSize > sl.size {
		return nil
	}

	if len(sl.signals) == 0 {
		return nil
	}

	lastSig := sl.signals[len(sl.signals)-1]
	if lastSig.GetRelativeStartPos()+lastSig.GetSize() > newSize {
		return ErrTooSmall
	}

	return nil
}

func (sl *SignalLayout) resize(newSize int) error {
	if err := sl.verifyBeforeResize(newSize); err != nil {
		return &MessageSizeError{
			Size: newSize,
			Err:  err,
		}
	}

	sl.size = newSize

	sl.generateFilters()

	return nil
}

func (sl *SignalLayout) shiftLeft(sigID EntityID, amount int) int {
	if amount <= 0 {
		return 0
	}

	perfShift := amount
	var prevSig Signal

	for idx, tmpSig := range sl.signals {
		if idx > 0 {
			prevSig = sl.signals[idx-1]
		}

		if sigID == tmpSig.EntityID() {
			tmpStartBit := tmpSig.GetRelativeStartPos()
			targetStartBit := tmpStartBit - amount

			if targetStartBit < 0 {
				targetStartBit = 0
			}

			if prevSig != nil {
				prevEndBit := prevSig.GetRelativeStartPos() + prevSig.GetSize()

				if targetStartBit < prevEndBit {
					targetStartBit = prevEndBit
				}
			}

			tmpSig.setRelativeStartPos(targetStartBit)
			perfShift = tmpStartBit - targetStartBit

			break
		}
	}

	sl.generateFilters()

	return perfShift
}

func (sl *SignalLayout) shiftRight(sigID EntityID, amount int) int {
	if amount <= 0 {
		return 0
	}

	perfShift := amount
	var nextSig Signal

	for idx, tmpSig := range sl.signals {
		if idx == len(sl.signals)-1 {
			nextSig = nil
		} else {
			nextSig = sl.signals[idx+1]
		}

		if sigID == tmpSig.EntityID() {
			tmpStartBit := tmpSig.GetRelativeStartPos()
			targetStartBit := tmpStartBit + amount
			targetEndBit := targetStartBit + tmpSig.GetSize()

			if targetEndBit > sl.size {
				targetStartBit = sl.size - tmpSig.GetSize()
			}

			if nextSig != nil {
				nextStartBit := nextSig.GetRelativeStartPos()

				if targetEndBit > nextStartBit {
					targetStartBit = nextStartBit - tmpSig.GetSize()
				}
			}

			tmpSig.setRelativeStartPos(targetStartBit)
			perfShift = targetStartBit - tmpStartBit

			break
		}
	}

	sl.generateFilters()

	return perfShift
}

func (sl *SignalLayout) stringify(b *strings.Builder, tabs int) {
	tabStr := getTabString(tabs)

	b.WriteString(fmt.Sprintf("%ssize: %d\n", tabStr, sl.size))
	b.WriteString(fmt.Sprintf("%ssignal_count: %d\n", tabStr, len(sl.signals)))

	// for _, f := range sl.filters {
	// 	//f.stringify(b, tabs)
	// }
}

func (sl *SignalLayout) String() string {
	b := new(strings.Builder)
	sl.stringify(b, 0)
	return b.String()
}

// Filters returns the signal filters of the [SignalLayout].
func (sl *SignalLayout) Filters() []*SignalLayoutFilter {
	return sl.filters
}

// Decode decodes the data according to the [SignalLayout].
// It retruns a list of signal decodings.
func (sl *SignalLayout) Decode(data []byte) []*SignalDecoding {
	signalCount := len(sl.signals)

	if signalCount == 0 {
		return nil
	}

	decodings := make([]*SignalDecoding, 0, signalCount)

	// Only used for shifting left little endian data
	consumedBits := 0

	prevEntID := EntityID("")
	var currSig Signal
	var rawValue uint64

	// Filters are sorted by entity id, so only adiacent filters belong to the same signal
	for _, filter := range sl.filters {
		entID := filter.signal.EntityID()

		// New signal to filter
		if entID != prevEntID {
			if currSig != nil {
				decodings = append(decodings, sl.decodeSignal(currSig, rawValue))
			}

			prevEntID = entID
			consumedBits = 0
			currSig = filter.signal
			rawValue = 0
		}

		tmpData := uint64((data[filter.byteIdx] & filter.mask) >> filter.leftOffset)

		// Little endian
		if filter.signal.Endianness() == EndiannessLittleEndian {
			tmpData <<= consumedBits
			rawValue |= tmpData
			consumedBits += filter.length
			continue
		}

		// Big endian
		rawValue <<= uint64(filter.length)
		rawValue |= tmpData
	}

	if currSig != nil {
		decodings = append(decodings, sl.decodeSignal(currSig, rawValue))
	}

	return decodings
}

func (sl *SignalLayout) decodeSignal(sig Signal, rawValue uint64) *SignalDecoding {
	switch sig.Kind() {
	case SignalKindStandard:
		stdSig, err := sig.ToStandard()
		if err != nil {
			panic(err)
		}

		return sl.decodeStandardSignal(stdSig, rawValue)

	case SignalKindEnum:
		enumSig, err := sig.ToEnum()
		if err != nil {
			panic(err)
		}

		return sl.decodeEnumSignal(enumSig, rawValue)
	}

	return nil
}

func (sl *SignalLayout) decodeStandardSignal(stdSig *StandardSignal, rawValue uint64) *SignalDecoding {
	var value any
	var valueType SignalValueType
	var unit string

	sigUnit := stdSig.unit
	if sigUnit != nil {
		unit = sigUnit.symbol
	}

	sigType := stdSig.typ
	switch sigType.kind {
	case SignalTypeKindFlag:
		valueType = SignalValueTypeFlag
		value = rawValue != 0

	case SignalTypeKindInteger:
		if sigType.signed {
			valueType = SignalValueTypeInt

			if rawValue&(1<<sigType.size-1) != 0 {
				// extend sign of raw value
				rawValue |= (1<<64 - 1) << sigType.size
			}

			value = int64(rawValue)*int64(sigType.scale) - int64(sigType.offset)

		} else {
			valueType = SignalValueTypeUint
			value = rawValue*uint64(sigType.scale) - uint64(sigType.offset)
		}

	case SignalTypeKindDecimal, SignalTypeKindCustom:
		valueType = SignalValueTypeFloat
		value = float64(rawValue)*sigType.scale + sigType.offset
	}

	return &SignalDecoding{
		Signal:    stdSig,
		RawValue:  rawValue,
		ValueType: valueType,
		Value:     value,
		Unit:      unit,
	}
}

func (sl *SignalLayout) decodeEnumSignal(enumSig *EnumSignal, rawValue uint64) *SignalDecoding {
	res := &SignalDecoding{
		Signal:    enumSig,
		RawValue:  rawValue,
		ValueType: SignalValueTypeEnum,
		Value:     "",
	}

	sigEnum := enumSig.enum

	for _, enumVal := range sigEnum.values.entries() {
		if enumVal.index == int(rawValue) {
			res.Value = enumVal.name
			break
		}
	}

	return res
}

// SignalValueType defines the value type of a [Signal] when decoded.
type SignalValueType string

const (
	// SignalValueTypeFlag defines a flag signal value type.
	SignalValueTypeFlag SignalValueType = "flag"
	// SignalValueTypeInt defines an integer signal value type.
	SignalValueTypeInt SignalValueType = "int"
	// SignalValueTypeUint defines an unsigned integer signal value type.
	SignalValueTypeUint SignalValueType = "uint"
	// SignalValueTypeFloat defines a float signal value type.
	SignalValueTypeFloat SignalValueType = "float"
	// SignalValueTypeEnum defines an enum signal value type.
	SignalValueTypeEnum SignalValueType = "enum"
)

// SignalDecoding represents a decoded of a signal.
type SignalDecoding struct {
	Signal    Signal
	RawValue  uint64
	ValueType SignalValueType
	Value     any
	Unit      string
}

// ValueAsFlag returns the decoded value as a flag.
// Returns false if the value type is not a flag.
func (sd *SignalDecoding) ValueAsFlag() bool {
	if sd.ValueType != SignalValueTypeFlag {
		return false
	}
	return sd.Value.(bool)
}

// ValueAsInt returns the decoded value as an integer.
// Returns 0 if the value type is not an integer.
func (sd *SignalDecoding) ValueAsInt() int64 {
	if sd.ValueType != SignalValueTypeInt {
		return 0
	}
	return sd.Value.(int64)
}

// ValueAsUint returns the decoded value as an unsigned integer.
// Returns 0 if the value type is not an unsigned integer.
func (sd *SignalDecoding) ValueAsUint() uint64 {
	if sd.ValueType != SignalValueTypeUint {
		return 0
	}
	return sd.Value.(uint64)
}

// ValueAsFloat returns the decoded value as a float.
// Returns 0 if the value type is not a float.
func (sd *SignalDecoding) ValueAsFloat() float64 {
	if sd.ValueType != SignalValueTypeFloat {
		return 0
	}
	return sd.Value.(float64)
}

// ValueAsEnum returns the decoded value as an enum.
// Returns an empty string if the value type is not an enum.
func (sd *SignalDecoding) ValueAsEnum() string {
	if sd.ValueType != SignalValueTypeEnum {
		return ""
	}
	return sd.Value.(string)
}

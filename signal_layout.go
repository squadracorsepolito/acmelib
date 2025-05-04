package acmelib

import (
	"cmp"
	"slices"

	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// SignalLayout represents the payload of a [Message]/[MultiplexedLayer]
// that carries a set of [Signal].
type SignalLayout struct {
	sizeByte int
	ibst     *collection.IBST[Signal]
	filters  []*SignalLayoutFilter

	muxLayers *collection.Map[EntityID, *MultiplexedLayer]

	parentMsg      *Message
	parentMuxLayer *MultiplexedLayer
}

func newSL(sizeByte int) *SignalLayout {
	return &SignalLayout{
		sizeByte: sizeByte,
		ibst:     collection.NewIBST[Signal](),
		filters:  []*SignalLayoutFilter{},

		muxLayers: collection.NewMap[EntityID, *MultiplexedLayer](),

		parentMsg:      nil,
		parentMuxLayer: nil,
	}
}

func (sl *SignalLayout) stringify(s *stringer.Stringer) {
	s.Write("size_byte: %d\n", sl.sizeByte)

	s.Write("interval_bst:\n")
	s.Indent()
	sl.ibst.Stringify(s)
	s.Unindent()

	if sl.ibst.Size() > 0 {
		s.Write("signals:\n")
		s.Indent()
		for sig := range sl.ibst.InOrder() {
			sig.stringify(s)
		}
		s.Unindent()
	}

	if len(sl.filters) > 0 {
		s.Write("filters:\n")
		s.Indent()
		for _, f := range sl.filters {
			f.stringify(s)
		}
		s.Unindent()
	}

	if sl.muxLayers.Size() > 0 {
		s.Write("multiplexed_layers:\n")
		s.Indent()
		for _, muxLayer := range sl.MultiplexedLayers() {
			muxLayer.stringify(s)
		}
		s.Unindent()
	}
}

func (sl *SignalLayout) String() string {
	s := stringer.New()
	s.Write("signal_layout:\n")
	sl.stringify(s)
	return s.String()
}

///////////////
// --------- //
// UTILITIES //
// --------- //
///////////////

// setParentMsg sets the parent message of the signal layout.
// It has to be called by the signal layout of a message.
func (sl *SignalLayout) setParentMsg(msg *Message) {
	sl.parentMsg = msg
}

func (sl *SignalLayout) fromMessage() bool {
	return sl.parentMsg != nil && sl.parentMuxLayer == nil
}

// setParentMuxLayer sets the parent mux layer of the signal layout.
// It has to be called by the signal layouts of a multiplexed layers.
func (sl *SignalLayout) setParentMuxLayer(muxLayer *MultiplexedLayer) {
	sl.parentMuxLayer = muxLayer
}

func (sl *SignalLayout) fromMultiplexedLayer() bool {
	return sl.parentMsg == nil && sl.parentMuxLayer != nil
}

// verifyIntersection checks if the signal intersects with any signal in the layout tree
// excluding the current signal layout.
func (sl *SignalLayout) verifyIntersection(sig Signal) error {
	// Traverse the tree of signal layouts downwards starting from the current signal layout,
	// and check the intersection in each child layout
	layoutStack := collection.NewStack[*SignalLayout]()

	// Push to the stack the multiplexed layers directly attached to the signal layout
	for ml := range sl.muxLayers.Values() {
		for _, tmpLayout := range ml.iterLayouts() {
			layoutStack.Push(tmpLayout)
		}
	}

	for !layoutStack.IsEmpty() {
		layout := layoutStack.Pop()

		// Check if the signal intersects in the current signal layout
		if intSig, ok := layout.ibst.Intersects(sig); ok {
			return newIntersectionError(intSig.Name())
		}

		// Push to the stack the multiplexed layers directly attached to the current signal layout
		for ml := range layout.muxLayers.Values() {
			for _, tmpLayout := range ml.iterLayouts() {
				layoutStack.Push(tmpLayout)
			}
		}
	}

	// If the layout is a direct child of a message, you can stop
	if sl.fromMessage() || !sl.fromMultiplexedLayer() {
		return nil
	}

	// The layout is child of a multiplexed layer, so you have to
	// traverse the tree backwards to check the intersection in parallel branches
	skipBranch := sl.parentMuxLayer.getID()
	parentLayout := sl.parentMuxLayer.attachedLayout
	if parentLayout == nil {
		return nil
	}

	for !parentLayout.fromMessage() && parentLayout.fromMultiplexedLayer() {
		if intSig, ok := parentLayout.ibst.Intersects(sig); ok {
			return newIntersectionError(intSig.Name())
		}

		skipBranch = parentLayout.parentMuxLayer.getID()
		parentLayout = parentLayout.parentMuxLayer.attachedLayout
		if parentLayout == nil {
			return nil
		}
	}

	// The tree was traversed backwards to the root, now check the intersection
	// in the parallel branches
	layoutStack.Push(parentLayout)
	for !layoutStack.IsEmpty() {
		layout := layoutStack.Pop()

		if intSig, ok := layout.ibst.Intersects(sig); ok {
			return newIntersectionError(intSig.Name())
		}

		for muxLayer := range layout.muxLayers.Values() {
			if muxLayer.getID() == skipBranch {
				continue
			}

			for _, tmpLayout := range muxLayer.iterLayouts() {
				layoutStack.Push(tmpLayout)
			}
		}
	}

	return nil
}

////////////////////
// -------------- //
// START POSITION //
// -------------- //
////////////////////

func (sl *SignalLayout) getIntervalFromNewStartPos(sig Signal, newStartPos int) (int, int) {
	return newStartPos, newStartPos + sig.Size() - 1
}

// verifyStartPos checks if the start position is valid.
func (sl *SignalLayout) verifyStartPos(startPos int) error {
	if startPos < 0 {
		return newStartPosError(startPos, ErrIsNegative)
	}

	if startPos > sl.sizeByte*8 {
		return newStartPosError(startPos, ErrOutOfBounds)
	}

	return nil
}

// verifyStartPosPlusSize checks if the start position plus the size is valid.
// It doesn't check if the size or the start position are valid, it only checks if the sum is valid.
func (sl *SignalLayout) verifyStartPosPlusSize(startPos, size int) error {
	if startPos+size > sl.sizeByte*8 {
		return newStartPosError(startPos, ErrOutOfBounds)
	}

	return nil
}

// verifyNewStartPos checks if setting the signal to the new start position
// does not intersect with another one.
func (sl *SignalLayout) verifyNewStartPos(sig Signal, newStartPos int) error {
	if err := sl.verifyStartPos(newStartPos); err != nil {
		return err
	}

	if err := sl.verifyStartPosPlusSize(newStartPos, sig.Size()); err != nil {
		return err
	}

	newLow, newHigh := sl.getIntervalFromNewStartPos(sig, newStartPos)
	if !sl.ibst.CanUpdate(sig, newLow, newHigh) {
		return newStartPosError(newStartPos, ErrIntersects)
	}

	// Set the start position to test the intersection and then reset it to the previous one
	oldStartPos := sig.StartPos()
	sig.setStartPos(newStartPos)
	defer sig.setStartPos(oldStartPos)

	return sl.verifyIntersection(sig)
}

// updateStartPos updates the start position of the signal in the signal layout.
// It must be called after the verify function since it assumes the
// new start position is valid.
func (sl *SignalLayout) updateStartPos(sig Signal, newStartPos int) {
	newLow, newHigh := sl.getIntervalFromNewStartPos(sig, newStartPos)
	sl.ibst.Update(sig, newLow, newHigh)

	// Regenerate the filters
	sl.genFilters()
}

// verifyAndUpdateStartPos checks and updates the start position of the signal in the signal layout.
func (sl *SignalLayout) verifyAndUpdateStartPos(sig Signal, newStartPos int) error {
	if err := sl.verifyNewStartPos(sig, newStartPos); err != nil {
		return err
	}

	sl.updateStartPos(sig, newStartPos)

	return nil
}

//////////
// ---- //
// SIZE //
// ---- //
//////////

func (sl *SignalLayout) getIntervalFromNewSize(sig Signal, newSize int) (int, int) {
	startPos := sig.StartPos()
	return startPos, startPos + newSize - 1
}

// verifySize checks if the size is valid.
func (sl *SignalLayout) verifySize(size int) error {
	if size < 0 {
		return newSizeError(size, ErrIsNegative)
	}

	if size > sl.sizeByte*8 {
		return newSizeError(size, ErrOutOfBounds)
	}

	return nil
}

// verifyNewSize checks if setting the signal to the new size
// does not intersect with another one.
func (sl *SignalLayout) verifyNewSize(sig Signal, newSize int) error {
	if err := sl.verifySize(newSize); err != nil {
		return err
	}

	if err := sl.verifyStartPosPlusSize(sig.StartPos(), newSize); err != nil {
		return err
	}

	newLow, newHigh := sl.getIntervalFromNewSize(sig, newSize)
	if !sl.ibst.CanUpdate(sig, newLow, newHigh) {
		return newSizeError(newSize, ErrIntersects)
	}

	oldSize := sig.Size()
	sig.setSize(newSize)
	defer sig.setSize(oldSize)

	return sl.verifyIntersection(sig)
}

// updateSize updates the size of the signal in the signal layout.
// It must be called after the verify function since it assumes
// the new size is valid.
func (sl *SignalLayout) updateSize(sig Signal, newSize int) {
	newLow, newHigh := sl.getIntervalFromNewSize(sig, newSize)
	sl.ibst.Update(sig, newLow, newHigh)

	// Regenerate the filters
	sl.genFilters()
}

// verifyAndUpdateSize checks and updates the size of the signal in the signal layout.
func (sl *SignalLayout) verifyAndUpdateSize(sig Signal, newSize int) error {
	if err := sl.verifyNewSize(sig, newSize); err != nil {
		return err
	}

	sl.updateSize(sig, newSize)

	return nil
}

////////////
// ------ //
// INSERT //
// ------ //
////////////

func (sl *SignalLayout) verifyInsert(sig Signal, startPos int) error {
	if err := sl.verifyStartPos(startPos); err != nil {
		return err
	}

	if err := sl.verifyStartPosPlusSize(startPos, sig.Size()); err != nil {
		return err
	}

	oldStartPos := sig.StartPos()
	sig.setStartPos(startPos)
	defer sig.setStartPos(oldStartPos)

	// Check if the signal intersects in the current signal layout
	if intSig, ok := sl.ibst.Intersects(sig); ok {
		return newIntersectionError(intSig.Name())
	}

	return sl.verifyIntersection(sig)
}

// insert inserts the signal into the signal layout.
// It must be called after the verify function since it assumes the signal is valid.
func (sl *SignalLayout) insert(sig Signal, startPos int) {
	sig.setStartPos(startPos)
	sl.ibst.Insert(sig)

	sig.setLayout(sl)

	// Regenerate the filters
	sl.genFilters()
}

// verifyAndInsert checks and inserts the signal into the signal layout.
func (sl *SignalLayout) verifyAndInsert(sig Signal, startPos int) error {
	if err := sl.verifyInsert(sig, startPos); err != nil {
		return err
	}

	sl.insert(sig, startPos)

	return nil
}

////////////
// ------ //
// DELETE //
// ------ //
////////////

// delete removes the signal from the signal layout.
func (sl *SignalLayout) delete(sig Signal) {
	sl.ibst.Delete(sig)

	sig.setLayout(nil)

	// Regenerate the filters
	sl.genFilters()
}

// clear removes all signals from the signal layout.
func (sl *SignalLayout) clear() {
	sl.ibst.Clear()

	// Reset the filters
	sl.filters = []*SignalLayoutFilter{}
}

////////////
// ------ //
// RESIZE //
// ------ //
////////////

// verifySizeByte checks if the size byte is valid.
func (sl *SignalLayout) verifySizeByte(sizeByte int) error {
	if sizeByte < 0 {
		return newSizeError(sizeByte, ErrIsNegative)
	}

	if sizeByte > 8 {
		return newSizeError(sizeByte, ErrTooBig)
	}

	return nil
}

// verifyResize checks if the layout can be resized to the new size.
func (sl *SignalLayout) verifyResize(newSizeByte int) error {
	if err := sl.verifySizeByte(newSizeByte); err != nil {
		return err
	}

	// Return early if the new size is the same or there are no signals
	if newSizeByte == sl.sizeByte || sl.ibst.Size() == 0 {
		return nil
	}

	// Check if the new size is too small by looking the last signal
	endPos := 0
	for sig := range sl.ibst.ReverseOrder() {
		endPos = sig.GetHigh()
		break
	}

	if endPos >= newSizeByte*8 {
		return newSizeError(newSizeByte, ErrTooSmall)
	}

	// Recursively check the multiplexed layers
	for ml := range sl.muxLayers.Values() {
		for _, muxLayout := range ml.iterLayouts() {
			if err := muxLayout.verifyResize(newSizeByte); err != nil {
				return err
			}
		}
	}

	return nil
}

// resize resizes the signal layout.
// It must be called after the verify function since it assumes the new size is valid.
func (sl *SignalLayout) resize(newSizeByte int) {
	sl.sizeByte = newSizeByte

	// Resize recursively each signal layout of the multiplexed layers
	for ml := range sl.muxLayers.Values() {
		for _, muxLayout := range ml.iterLayouts() {
			muxLayout.resize(newSizeByte)
		}

		ml.sizeByte = newSizeByte
	}
}

// verifyAndResize checks and resizes the signal layout to the new size.
func (sl *SignalLayout) verifyAndResize(newSizeByte int) error {
	if err := sl.verifyResize(newSizeByte); err != nil {
		return err
	}

	sl.resize(newSizeByte)

	return nil
}

/////////////
// ------- //
// COMPACT //
// ------- //
/////////////

// Compact compacts the signal layout.
// It will only compact the signal layout if there are no multiplexed layers attached.
func (sl *SignalLayout) Compact() {
	if sl.ibst.Size() == 0 || sl.muxLayers.Size() != 0 {
		return
	}

	type signalToUpdate struct {
		sig             Signal
		newLow, newHigh int
	}

	signalsToUpdate := []signalToUpdate{}

	// Get the signals that need to be updated
	newStartPos := 0
	for sig := range sl.ibst.InOrder() {
		tmpSize := sig.Size()
		signalsToUpdate = append(signalsToUpdate, signalToUpdate{sig, newStartPos, newStartPos + tmpSize})
		newStartPos += tmpSize
	}

	// Update the start position of signals
	for _, sigToUpd := range signalsToUpdate {
		sl.ibst.Update(sigToUpd.sig, sigToUpd.newLow, sigToUpd.newHigh)
		sigToUpd.sig.setStartPos(sigToUpd.newLow)
	}
}

///////////////////////
// ----------------- //
// MULTIPLEXED LAYER //
// ----------------- //
///////////////////////

// MultiplexedLayers returns the multiplexed layers in the signal layout
// ordered by the start position of the muxor signal.
func (sl *SignalLayout) MultiplexedLayers() []*MultiplexedLayer {
	layers := slices.Collect(sl.muxLayers.Values())
	slices.SortFunc(layers, func(a, b *MultiplexedLayer) int {
		return cmp.Compare(a.muxor.StartPos(), b.muxor.StartPos())
	})
	return layers
}

// AddMultiplexedLayer creates and adds a new multiplexed layer to the signal layout.
// A newly created multiplexed layer is returned. It will contain a muxor signal
// with the given name, start postion, and layout count.
//
// It returns:
//   - [ArgError] if an argument is invalid.
//   - [NameError] if the given muxor name is invalid.
//   - [StartPosError] if the given muxor start position is invalid.
func (sl *SignalLayout) AddMultiplexedLayer(muxorName string, muxorStartPos, layoutCount int) (*MultiplexedLayer, error) {
	// Check if the layout count is valid
	if layoutCount < 0 {
		return nil, newArgError("layoutCount", ErrIsNegative)
	} else if layoutCount == 0 {
		return nil, newArgError("layoutCount", ErrIsZero)
	}

	// Check if the muxor name is valid
	if sl.fromMessage() {
		if err := sl.parentMsg.verifySignalName(muxorName); err != nil {
			return nil, err
		}
	} else if sl.fromMultiplexedLayer() {
		if err := sl.parentMuxLayer.verifySignalName(muxorName); err != nil {
			return nil, err
		}
	}

	// Create the muxor, but before inserting it, check if the
	// start position is valid
	muxor := newMuxorSignal(muxorName, layoutCount)
	if err := sl.verifyInsert(muxor, muxorStartPos); err != nil {
		return nil, muxor.errorf(err)
	}
	sl.insert(muxor, muxorStartPos)

	// Create the multiplexed layer and add it
	ml := newMultiplexedLayer(muxor, layoutCount, sl.sizeByte)
	ml.setAttachedLayout(sl)
	sl.muxLayers.Set(ml.getID(), ml)

	// Generate filters
	sl.genFilters()

	return ml, nil
}

// DeleteMultiplexedLayer removes the multiplexed layer
// with the given entity ID from the signal layout.
//
// It returns [ErrNotFound] if the multiplexed layer is not found.
func (sl *SignalLayout) DeleteMultiplexedLayer(entityID EntityID) error {
	ml, ok := sl.muxLayers.Get(entityID)
	if !ok {
		return ErrNotFound
	}

	sl.delete(ml.muxor)
	sl.muxLayers.Delete(entityID)
	ml.setAttachedLayout(nil)

	return nil
}

/////////////////////
// --------------- //
// EXPORTED VALUES //
// --------------- //
/////////////////////

// SignalCount returns the number of signals in the signal layout.
func (sl *SignalLayout) SignalCount() int {
	return sl.ibst.Size()
}

// Signals returns the signals in the signal layout ordered by the start position.
func (sl *SignalLayout) Signals() []Signal {
	return sl.ibst.GetInOrder()
}

// Filters returns the signal filters of the layout.
func (sl *SignalLayout) Filters() []*SignalLayoutFilter {
	return sl.filters
}

/////////////
// ------- //
// FILTERS //
// ------- //
/////////////

// genFilters generates the signal layout filters.
// It must be called every time the layout is changed.
func (sl *SignalLayout) genFilters() {
	sl.filters = []*SignalLayoutFilter{}

	for _, sig := range sl.ibst.GetInOrder() {
		sigSize := sig.Size()
		startPos := sig.StartPos()
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

	for muxLayer := range sl.muxLayers.Values() {
		for _, sl := range muxLayer.iterLayouts() {
			sl.genFilters()
		}
	}
}

////////////
// ------ //
// DECODE //
// ------ //
////////////

// decodeSignal decodes the signal with the given raw value.
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

	case SignalKindMuxor:
		muxorSig, err := sig.ToMuxor()
		if err != nil {
			panic(err)
		}

		return sl.decodeMuxorSignal(muxorSig, rawValue)
	}

	return nil
}

// decodeStandardSignal interprets the raw value of a standard signal.
func (sl *SignalLayout) decodeStandardSignal(stdSig *StandardSignal, rawValue uint64) *SignalDecoding {
	var value any
	var valueType SignalValueType
	var unit string

	sigUnit := stdSig.unit
	if sigUnit != nil {
		unit = sigUnit.symbol
	}

	sigType := stdSig.typ

	if sigType.signed {
		if rawValue&(1<<sigType.size-1) != 0 {
			// extend sign of raw value
			rawValue |= (1<<64 - 1) << sigType.size
		}
	}

	switch sigType.kind {
	case SignalTypeKindFlag:
		valueType = SignalValueTypeFlag
		value = rawValue != 0

	case SignalTypeKindInteger:
		if sigType.signed {
			valueType = SignalValueTypeInt
			value = int64(rawValue)*int64(sigType.scale) - int64(sigType.offset)
		} else {
			valueType = SignalValueTypeUint
			value = rawValue*uint64(sigType.scale) - uint64(sigType.offset)
		}

	case SignalTypeKindDecimal:
		valueType = SignalValueTypeFloat

		if sigType.signed {
			value = float64(int64(rawValue))*sigType.scale + sigType.offset
		} else {
			value = float64(rawValue)*sigType.scale + sigType.offset
		}
	}

	return newSignalDecoding(stdSig, rawValue, valueType, value, unit)
}

// decodeEnumSignal interprets the raw value of an enum signal.
func (sl *SignalLayout) decodeEnumSignal(enumSig *EnumSignal, rawValue uint64) *SignalDecoding {
	dec := newSignalDecoding(enumSig, rawValue, SignalValueTypeEnum, "", "")

	sigEnum := enumSig.enum

	for _, enumVal := range sigEnum.values {
		if enumVal.index == int(rawValue) {
			dec.Value = enumVal.name
			break
		}
	}

	if dec.Value == "" {
		dec.Value = "unknown"
	}

	return dec
}

// decodeMuxorSignal interprets the raw value of a muxor signal.
func (sl *SignalLayout) decodeMuxorSignal(muxorSig *MuxorSignal, rawValue uint64) *SignalDecoding {
	return newSignalDecoding(muxorSig, rawValue, SignalValueTypeUint, rawValue, "")
}

// decodeCurrentSignal decodes the signal with the given raw value.
// It also calls recursively the Decode method for decoding signals
// contained into a multiplexed layer.
func (sl *SignalLayout) decodeCurrentSignal(decodings *[]*SignalDecoding, data []byte, sig Signal, rawValue uint64) {
	*decodings = append(*decodings, sl.decodeSignal(sig, rawValue))

	if sig.Kind() != SignalKindMuxor {
		return
	}

	layoutID := int(rawValue)
	muxorSig, err := sig.ToMuxor()
	if err != nil {
		panic(err)
	}

	muxLayout := muxorSig.parentMuxLayer.GetLayout(layoutID)
	if muxLayout == nil {
		return
	}

	*decodings = append(*decodings, muxLayout.Decode(data)...)
}

// Decode decodes the given data using the information in the signal layout.
// It returns a slice of [SignalDecoding] structs.
// When the layout contains multiplexed layers, it calls recursively the Decode method.
// In that case, the decoded muxor signal are placed in the slice before
// the multiplexed layer signals.
func (sl *SignalLayout) Decode(data []byte) []*SignalDecoding {
	signalCount := sl.ibst.Size()

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
				sl.decodeCurrentSignal(&decodings, data, currSig, rawValue)
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
		sl.decodeCurrentSignal(&decodings, data, currSig, rawValue)
	}

	return decodings
}

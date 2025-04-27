package acmelib

import (
	"cmp"
	"slices"

	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/ibst"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

type SL struct {
	sizeByte int
	tree     *ibst.Tree[Signal]
	filters  []*SignalLayoutFilter

	muxLayers *collection.Map[EntityID, *MultiplexedLayer]

	parentMsg      *Message
	parentMuxLayer *MultiplexedLayer
}

func newSL(sizeByte int) *SL {
	return &SL{
		sizeByte: sizeByte,
		tree:     ibst.NewTree[Signal](),
		filters:  []*SignalLayoutFilter{},

		muxLayers: collection.NewMap[EntityID, *MultiplexedLayer](),

		parentMsg:      nil,
		parentMuxLayer: nil,
	}
}

// genFilters generates the signal layout filters.
// It must be called every time the layout is changed.
func (sl *SL) genFilters() {
	sl.filters = []*SignalLayoutFilter{}

	for _, sig := range sl.tree.GetInOrder() {
		sigSize := sig.GetSize()
		startPos := sig.GetStartPos()
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

func (sl *SL) stringify(s *stringer.Stringer) {
	s.Write("size_byte: %d\n", sl.sizeByte)

	s.Write("interval_bst:\n")
	s.Indent()
	sl.tree.Stringify(s)
	s.Unindent()

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

func (sl *SL) String() string {
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
func (sl *SL) setParentMsg(msg *Message) {
	sl.parentMsg = msg
}

func (sl *SL) fromMessage() bool {
	return sl.parentMsg != nil && sl.parentMuxLayer == nil
}

// setParentMuxLayer sets the parent mux layer of the signal layout.
// It has to be called by the signal layouts of a multiplexed layers.
func (sl *SL) setParentMuxLayer(muxLayer *MultiplexedLayer) {
	sl.parentMuxLayer = muxLayer
}

func (sl *SL) fromMultiplexedLayer() bool {
	return sl.parentMsg == nil && sl.parentMuxLayer != nil
}

// verifyIntersection checks if the signal intersects with any signal in the layout tree
// excluding the current signal layout.
func (sl *SL) verifyIntersection(sig Signal) error {
	// Traverse the tree of signal layouts downwards starting from the current signal layout,
	// and check the intersection in each child layout
	layoutStack := collection.NewStack[*SL]()

	// Push to the stack the multiplexed layers directly attached to the signal layout
	for ml := range sl.muxLayers.Values() {
		for _, tmpLayout := range ml.iterLayouts() {
			layoutStack.Push(tmpLayout)
		}
	}

	for !layoutStack.IsEmpty() {
		layout := layoutStack.Pop()

		// Check if the signal intersects in the current signal layout
		if intSig, ok := layout.tree.Intersects(sig); ok {
			return newIntersectError(intSig.Name())
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
		if intSig, ok := parentLayout.tree.Intersects(sig); ok {
			return newIntersectError(intSig.Name())
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

		if intSig, ok := layout.tree.Intersects(sig); ok {
			return newIntersectError(intSig.Name())
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

func (sl *SL) getIntervalFromNewStartPos(sig Signal, newStartPos int) (int, int) {
	return newStartPos, newStartPos + sig.GetSize() - 1
}

// verifyStartPos checks if the start position is valid.
func (sl *SL) verifyStartPos(startPos int) error {
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
func (sl *SL) verifyStartPosPlusSize(startPos, size int) error {
	if startPos+size > sl.sizeByte*8 {
		return newStartPosError(startPos, ErrOutOfBounds)
	}

	return nil
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
		return newStartPosError(newStartPos, ErrIntersects)
	}

	// Set the start position to test the intersection and then reset it to the previous one
	oldStartPos := sig.GetStartPos()
	sig.setStartPos(newStartPos)
	defer sig.setStartPos(oldStartPos)

	return sl.verifyIntersection(sig)
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

//////////
// ---- //
// SIZE //
// ---- //
//////////

func (sl *SL) getIntervalFromNewSize(sig Signal, newSize int) (int, int) {
	startPos := sig.GetStartPos()
	return startPos, startPos + newSize - 1
}

// verifySize checks if the size is valid.
func (sl *SL) verifySize(size int) error {
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
func (sl *SL) verifyNewSize(sig Signal, newSize int) error {
	if err := sl.verifySize(newSize); err != nil {
		return err
	}

	if err := sl.verifyStartPosPlusSize(sig.GetStartPos(), newSize); err != nil {
		return err
	}

	newLow, newHigh := sl.getIntervalFromNewSize(sig, newSize)
	if !sl.tree.CanUpdate(sig, newLow, newHigh) {
		return newSizeError(newSize, ErrIntersects)
	}

	oldSize := sig.GetSize()
	sig.setSize(newSize)
	defer sig.setSize(oldSize)

	return sl.verifyIntersection(sig)
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

////////////
// ------ //
// INSERT //
// ------ //
////////////

func (sl *SL) verifyInsert(sig Signal, startPos int) error {
	if err := sl.verifyStartPos(startPos); err != nil {
		return err
	}

	if err := sl.verifyStartPosPlusSize(startPos, sig.GetSize()); err != nil {
		return err
	}

	oldStartPos := sig.GetStartPos()
	sig.setStartPos(startPos)
	defer sig.setStartPos(oldStartPos)

	// Check if the signal intersects in the current signal layout
	if intSig, ok := sl.tree.Intersects(sig); ok {
		return newIntersectError(intSig.Name())
	}

	return sl.verifyIntersection(sig)
}

// insert inserts the signal into the signal layout.
// It must be called after the verify function since it assumes the signal is valid.
func (sl *SL) insert(sig Signal, startPos int) {
	sig.setStartPos(startPos)
	sl.tree.Insert(sig)

	sig.setLayout(sl)

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

////////////
// ------ //
// DELETE //
// ------ //
////////////

// delete removes the signal from the signal layout.
func (sl *SL) delete(sig Signal) {
	sl.tree.Delete(sig)

	sig.setLayout(nil)

	// Regenerate the filters
	sl.genFilters()
}

// clear removes all signals from the signal layout.
func (sl *SL) clear() {
	sl.tree.Clear()

	// Reset the filters
	sl.filters = []*SignalLayoutFilter{}
}

////////////
// ------ //
// RESIZE //
// ------ //
////////////

// verifySizeByte checks if the size byte is valid.
func (sl *SL) verifySizeByte(sizeByte int) error {
	if sizeByte < 0 {
		return newSizeError(sizeByte, ErrIsNegative)
	}

	if sizeByte > 8 {
		return newSizeError(sizeByte, ErrTooBig)
	}

	return nil
}

// verifyResize checks if the layout can be resized to the new size.
func (sl *SL) verifyResize(newSizeByte int) error {
	if err := sl.verifySizeByte(newSizeByte); err != nil {
		return err
	}

	// Return early if the new size is the same or there are no signals
	if newSizeByte == sl.sizeByte || sl.tree.Size() == 0 {
		return nil
	}

	// Check if the new size is too small by looking the last signal
	endPos := 0
	for sig := range sl.tree.ReverseOrder() {
		endPos = sig.GetHigh()
		break
	}

	if endPos >= newSizeByte {
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
func (sl *SL) resize(newSizeByte int) {
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
func (sl *SL) verifyAndResize(newSizeByte int) error {
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

// compact compacts the signal layout.
// It will only compact the signal layout if there are no multiplexed layers.
func (sl *SL) compact() {
	if sl.tree.Size() == 0 || sl.muxLayers.Size() != 0 {
		return
	}

	// Compact the signal layout
	newStartPos := 0
	for sig := range sl.tree.InOrder() {
		tmpSize := sig.GetSize()
		sl.tree.Update(sig, newStartPos, newStartPos+tmpSize)
		newStartPos += tmpSize
	}
}

// SignalCount returns the number of signals in the signal layout.
func (sl *SL) SignalCount() int {
	return sl.tree.Size()
}

// Signals returns the signals in the signal layout ordered by the start position.
func (sl *SL) Signals() []Signal {
	return sl.tree.GetInOrder()
}

// Filters returns the signal filters of the layout.
func (sl *SL) Filters() []*SignalLayoutFilter {
	return sl.filters
}

// MultiplexedLayers returns the multiplexed layers in the signal layout
// ordered by the start position of the muxor signal.
func (sl *SL) MultiplexedLayers() []*MultiplexedLayer {
	layers := slices.Collect(sl.muxLayers.Values())
	slices.SortFunc(layers, func(a, b *MultiplexedLayer) int {
		return cmp.Compare(a.muxor.GetStartPos(), b.muxor.GetStartPos())
	})
	return layers
}

// AddMultiplexedLayer creates and adds a new multiplexed layer to the signal layout.
// A newly created multiplexed layer is returned. It will contain a muxor signal
// with the given name, start postion, and layout count.
//
// It returns:
//   - [ArgumentError] if an argument is invalid.
//   - [NameError] if the given muxor name is invalid.
//   - [StartPosError] if the given muxor start position is invalid.
func (sl *SL) AddMultiplexedLayer(muxorName string, muxorStartPos, layoutCount int) (*MultiplexedLayer, error) {
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
func (sl *SL) DeleteMultiplexedLayer(entityID EntityID) error {
	ml, ok := sl.muxLayers.Get(entityID)
	if !ok {
		return ErrNotFound
	}

	sl.delete(ml.muxor)
	sl.muxLayers.Delete(entityID)
	ml.setAttachedLayout(nil)

	return nil
}

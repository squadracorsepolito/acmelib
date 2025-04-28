package acmelib

import (
	"iter"
	"slices"

	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// MultiplexedLayer represents a layer on top of a [SL] that has
// a [MuxorSignal] and N inner [SL] where N is the layout count.
type MultiplexedLayer struct {
	sizeByte int

	signals         *collection.Map[EntityID, Signal]
	signalNames     *collection.Map[string, EntityID]
	singalLayoutIDs *collection.Map[EntityID, []int]

	muxor          *MuxorSignal
	attachedLayout *SL

	layouts []*SL
}

func newMultiplexedLayer(muxor *MuxorSignal, layoutCount, sizeByte int) *MultiplexedLayer {
	ml := &MultiplexedLayer{
		sizeByte: sizeByte,

		signals:         collection.NewMap[EntityID, Signal](),
		signalNames:     collection.NewMap[string, EntityID](),
		singalLayoutIDs: collection.NewMap[EntityID, []int](),

		muxor: muxor,

		layouts: make([]*SL, 0, layoutCount),
	}

	ml.appendLayouts(layoutCount)

	muxor.setparentMuxLayer(ml)
	ml.signalNames.Set(muxor.Name(), muxor.EntityID())

	return ml
}

func (ml *MultiplexedLayer) getID() EntityID {
	return ml.muxor.entityID
}

func (ml *MultiplexedLayer) iterLayouts() iter.Seq2[int, *SL] {
	return func(yield func(int, *SL) bool) {
		for lID, sl := range ml.layouts {
			if !yield(lID, sl) {
				break
			}
		}
	}
}

func (ml *MultiplexedLayer) addSignal(sig Signal, layoutIDs []int) {
	ml.signals.Set(sig.EntityID(), sig)
	ml.signalNames.Set(sig.Name(), sig.EntityID())
	ml.singalLayoutIDs.Set(sig.EntityID(), layoutIDs)
	sig.setparentMuxLayer(ml)
}

func (ml *MultiplexedLayer) removeSignal(sig Signal) {
	ml.signals.Delete(sig.EntityID())
	ml.signalNames.Delete(sig.Name())
	ml.singalLayoutIDs.Delete(sig.EntityID())
	sig.setparentMuxLayer(nil)
}

func (ml *MultiplexedLayer) appendLayouts(layoutCount int) {
	for range layoutCount {
		sl := newSL(ml.sizeByte)
		sl.setParentMuxLayer(ml)
		ml.layouts = append(ml.layouts, sl)
	}
}

func (ml *MultiplexedLayer) truncateLayouts(fromLayoutID int) {
	for i := fromLayoutID; i < len(ml.layouts); i++ {
		ml.layouts[i].setParentMuxLayer(nil)
	}

	ml.layouts = slices.Delete(ml.layouts, fromLayoutID, len(ml.layouts))
}

// verifyLayoutID checks if the layout ID is valid.
//
// It returns:
//   - [ErrIsNegative] if the layout ID is negative
//   - [ErrOutOfBounds] if the layout ID is out of bounds
func (ml *MultiplexedLayer) verifyLayoutID(layoutID int) error {
	if layoutID < 0 {
		return newLayoutIDError(layoutID, ErrIsNegative)
	}

	if layoutID >= ml.GetLayoutCount() {
		return newLayoutIDError(layoutID, ErrOutOfBounds)
	}

	return nil
}

// verifySignalName checks if the signal name is already used in the multiplexed layer.
// It traverses all the multiplexed layers with BFS algorithm
// and checks if the name is already used.
func (ml *MultiplexedLayer) verifySignalName(name string) error {
	muxLayerQueue := collection.NewQueue[*MultiplexedLayer]()
	muxLayerQueue.Push(ml)

	visitedMuxLayers := collection.NewSet[EntityID]()

	for muxLayerQueue.Size() > 0 {
		muxLayer := muxLayerQueue.Pop()
		visitedMuxLayers.Add(muxLayer.getID())

		// Check if the name is present in the current multiplexed layer
		if muxLayer.signalNames.Has(name) {
			return newNameError(name, ErrIsDuplicated)
		}

		// Push multiplexed layers attached to the current multiplexed layer
		for _, innerLayout := range muxLayer.iterLayouts() {
			for innerMuxLayer := range innerLayout.muxLayers.Values() {
				if !visitedMuxLayers.Has(innerMuxLayer.muxor.entityID) {
					muxLayerQueue.Push(innerMuxLayer)
				}
			}
		}

		attachedLayout := muxLayer.attachedLayout
		if muxLayer.attachedLayout == nil {
			continue
		}

		// Push multiplexed layers directly attached to the attached layout (siblings)
		for siblingMuxLayer := range attachedLayout.muxLayers.Values() {
			if !visitedMuxLayers.Has(siblingMuxLayer.muxor.entityID) {
				muxLayerQueue.Push(siblingMuxLayer)
			}
		}

		if attachedLayout.parentMsg != nil {
			// The current multiplexed layer is directly attached to the parent message,
			// so check if the name is present in the parent message
			if attachedLayout.parentMsg.signalNames.Has(name) {
				return newNameError(name, ErrIsDuplicated)
			}
		}

		parentMuxLayer := attachedLayout.parentMuxLayer
		if parentMuxLayer != nil {
			if !visitedMuxLayers.Has(parentMuxLayer.muxor.entityID) {
				// Push the parent layer of the layout attached to
				// the current multiplexed layer
				muxLayerQueue.Push(parentMuxLayer)
			}
		}
	}

	return nil
}

func (ml *MultiplexedLayer) stringify(s *stringer.Stringer) {
	s.Write("layout_count: %d\n", ml.GetLayoutCount())

	s.Write("layouts:\n")
	s.Indent()
	for lID, sl := range ml.iterLayouts() {
		if sl.SignalCount() == 0 {
			continue
		}

		s.Write("layout_id: %d\n", lID)
		sl.stringify(s)
	}
	s.Unindent()
}

func (ml *MultiplexedLayer) String() string {
	s := stringer.New()
	s.Write("multiplexed_layer:\n")
	ml.stringify(s)
	return s.String()
}

// InsertSignal inserts a signal at the given start position in the given layout IDs.
//
// It returns:
//   - [ArgError] if the given signal is nil.
//   - [NameError] if the signal name is invalid.
//   - [StartPosError] if the given start position is invalid.
//   - [LayoutIDError] if the given layout ID is invalid.
func (ml *MultiplexedLayer) InsertSignal(signal Signal, startPos int, layoutIDs ...int) error {
	if signal == nil {
		return newArgError("signal", ErrIsNil)
	}

	if err := ml.verifySignalName(signal.Name()); err != nil {
		return signal.errorf(err)
	}

	// Check if the signal has to be inserted into all layouts
	if len(layoutIDs) == 0 {
		for lID := range ml.GetLayoutCount() {
			layoutIDs = append(layoutIDs, lID)
		}

		// Check if the start position is valid
		for _, sl := range ml.iterLayouts() {
			if err := sl.verifyInsert(signal, startPos); err != nil {
				return signal.errorf(err)
			}
		}

		// Insert the signal into all layouts
		for _, sl := range ml.iterLayouts() {
			sl.insert(signal, startPos)
		}

	} else {
		// Remove duplicate layout IDs
		layoutIDs = slices.Compact(layoutIDs)

		// Check if the signal is already present in other layouts
		prevLayoutIDs := []int{}
		if ml.singalLayoutIDs.Has(signal.EntityID()) {
			tmp, ok := ml.singalLayoutIDs.Get(signal.EntityID())
			if !ok {
				return nil
			}
			prevLayoutIDs = tmp
		}

		for _, lID := range layoutIDs {
			// Check if the layout ID is valid
			if err := ml.verifyLayoutID(lID); err != nil {
				return signal.errorf(err)
			}

			// Check if the current layout ID is already present
			if slices.Contains(prevLayoutIDs, lID) {
				return signal.errorf(newLayoutIDError(lID, ErrIsDuplicated))
			}

			// Check if the start position is valid
			if err := ml.layouts[lID].verifyInsert(signal, startPos); err != nil {
				return signal.errorf(err)
			}
		}

		// Insert the signal into the given layouts
		for _, lID := range layoutIDs {
			ml.layouts[lID].insert(signal, startPos)
		}
	}

	ml.addSignal(signal, layoutIDs)

	return nil
}

// DeleteSignal deletes the signal with the given entity ID from each layout.
//
// It returns [ErrNotFound] if the signal is not found.
func (ml *MultiplexedLayer) DeleteSignal(signalEntityID EntityID) error {
	sig, ok := ml.signals.Get(signalEntityID)
	if !ok {
		return ErrNotFound
	}

	// Get all the layout IDs of the signal
	layoutIDs, ok := ml.singalLayoutIDs.Get(signalEntityID)
	if !ok {
		return nil
	}

	// Remove the signal from all layouts
	for _, lID := range layoutIDs {
		ml.layouts[lID].delete(sig)
	}

	ml.removeSignal(sig)

	return nil
}

// ClearLayout deletes all signals from the layout with the given ID.
//
// It returns a [LayoutIDError] if the given layout ID is invalid.
func (ml *MultiplexedLayer) ClearLayout(layoutID int) error {
	if err := ml.verifyLayoutID(layoutID); err != nil {
		return err
	}

	// Get the layout
	layout := ml.layouts[layoutID]

	// Remove signals that are not present in other layouts
	for _, sig := range layout.Signals() {
		layoutIDs, ok := ml.singalLayoutIDs.Get(sig.EntityID())
		if !ok {
			return nil
		}

		// Check if the signal is present in other layouts
		if len(layoutIDs) > 1 {
			continue
		}

		// Remove the signal
		ml.removeSignal(sig)
	}

	// Clear the layout
	layout.clear()

	return nil
}

// Clear deletes all the signals from the layer.
func (ml *MultiplexedLayer) Clear() {
	// Remove all the signals
	for sig := range ml.signals.Values() {
		ml.removeSignal(sig)
	}

	// Clear all the layouts
	for _, sl := range ml.iterLayouts() {
		sl.clear()
	}
}

// Muxor returns the [MuxorSignal] of the layer.
func (ml *MultiplexedLayer) Muxor() *MuxorSignal {
	return ml.muxor
}

// Layouts returns the slice of all [SL] of the layer.
func (ml *MultiplexedLayer) Layouts() []*SL {
	return ml.layouts
}

// GetLayout returns the [SL] with the given layoud ID.
// It returns nil if the layout ID is invalid.
func (ml *MultiplexedLayer) GetLayout(layoutID int) *SL {
	if err := ml.verifyLayoutID(layoutID); err != nil {
		return nil
	}
	return ml.layouts[layoutID]
}

// GetSignalByName returns the signal with the given name.
//
// It returns [ErrNotFound] if the signal ith the given name is not found.
func (ml *MultiplexedLayer) GetSignalByName(name string) (Signal, error) {
	entID, ok := ml.signalNames.Get(name)
	if !ok {
		return nil, ErrNotFound
	}

	sig, ok := ml.signals.Get(entID)
	if !ok {
		return nil, ErrNotFound
	}

	return sig, nil
}

func (ml *MultiplexedLayer) setAttachedLayout(layout *SL) {
	ml.attachedLayout = layout
}

// AttachedLayout returns the [SL] attached to the current layer (parent layout).
func (ml *MultiplexedLayer) AttachedLayout() *SL {
	return ml.attachedLayout
}

// GetLayoutCount returns the number of layouts the layer can have.
func (ml *MultiplexedLayer) GetLayoutCount() int {
	return ml.muxor.layoutCount
}

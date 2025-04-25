package acmelib

import (
	"iter"
	"slices"

	"github.com/squadracorsepolito/acmelib/internal/collection"
	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

func verifyArgNotNil(arg any, name string) error {
	if arg == nil {
		return &ArgumentError{Name: name, Err: ErrIsNil}
	}

	return nil
}

type MultiplexedLayer struct {
	sizeByte int

	signals         *collection.Map[EntityID, Signal]
	signalNames     *collection.Map[string, EntityID]
	singalLayoutIDs *collection.Map[EntityID, []int]

	muxor          *MuxorSignal
	attachedLayout *SL

	layoutCount int
	layouts     []*SL
}

func NewMultiplexedLayer(sizeByte, layoutCount int, muxorName string) *MultiplexedLayer {
	ml := &MultiplexedLayer{
		sizeByte: sizeByte,

		signals:         collection.NewMap[EntityID, Signal](),
		signalNames:     collection.NewMap[string, EntityID](),
		singalLayoutIDs: collection.NewMap[EntityID, []int](),

		layoutCount: layoutCount,
		layouts:     make([]*SL, layoutCount),
	}

	// Generate signal layouts
	for i := range layoutCount {
		sl := newSL(sizeByte)
		sl.setParentMuxLayer(ml)
		ml.layouts[i] = sl
	}

	// Generate the muxor signal
	muxor := newMuxorSignal(muxorName, layoutCount)
	ml.muxor = muxor
	muxor.setMultiplexedLayer(ml)

	layoutIDs := make([]int, 0, layoutCount)
	for lID := range layoutCount {
		layoutIDs = append(layoutIDs, lID)
	}
	ml.addSignal(muxor, layoutIDs)

	return ml
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
	sig.setMultiplexedLayer(ml)
}

func (ml *MultiplexedLayer) removeSignal(sig Signal) {
	if sig.Kind() == SignalKindMuxor && sig.EntityID() == ml.muxor.EntityID() {
		return
	}

	ml.signals.Delete(sig.EntityID())
	ml.signalNames.Delete(sig.Name())
	ml.singalLayoutIDs.Delete(sig.EntityID())
	sig.setMultiplexedLayer(nil)
}

// verifyLayoutID checks if the layout ID is valid.
//
// It returns:
//   - [ErrIsNegative] if the layout ID is negative
//   - [ErrOutOfBounds] if the layout ID is out of bounds
func (ml *MultiplexedLayer) verifyLayoutID(layoutID int) error {
	if layoutID < 0 {
		return ErrIsNegative
	}

	if layoutID >= ml.layoutCount {
		return ErrOutOfBounds
	}

	return nil
}

func (ml *MultiplexedLayer) stringify(s *stringer.Stringer) {
	s.Write("layout_count: %d\n", ml.layoutCount)

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

func (ml *MultiplexedLayer) InsertSignal(signal Signal, startPos int, layoutIDs ...int) error {
	if err := verifyArgNotNil(signal, "signal"); err != nil {
		return err
	}

	if err := ml.verifySignalName(signal.Name()); err != nil {
		return err
	}

	// Check if it intersects with any signal of the attached layout
	if ml.attachedLayout != nil {
		if err := ml.attachedLayout.verifyInsert(signal, startPos, ml.muxor.entityID); err != nil {
			return err
		}
	}

	// Check if the signal has to be inserted into all layouts
	if len(layoutIDs) == 0 {
		for lID := range ml.layoutCount {
			layoutIDs = append(layoutIDs, lID)
		}

		// Check if the start position is valid
		for _, sl := range ml.iterLayouts() {
			if err := sl.verifyInsert(signal, startPos); err != nil {
				return err
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
				return newLayoutIDError(lID, err)
			}

			// Check if the current layout ID is already present
			if slices.Contains(prevLayoutIDs, lID) {
				return &GroupIDError{GroupID: lID, Err: ErrIsDuplicated}
			}

			// Check if the start position is valid
			if err := ml.layouts[lID].verifyInsert(signal, startPos); err != nil {
				return err
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

func (ml *MultiplexedLayer) DeleteSignal(signal Signal) error {
	sigEntID := signal.EntityID()

	if !ml.signals.Has(sigEntID) {
		return ErrNotFound
	}

	// Get all the layout IDs of the signal
	layoutIDs, ok := ml.singalLayoutIDs.Get(sigEntID)
	if !ok {
		return nil
	}

	// Remove the signal from all layouts
	for _, lID := range layoutIDs {
		ml.layouts[lID].delete(signal)
	}

	ml.removeSignal(signal)

	return nil
}

func (ml *MultiplexedLayer) ClearLayout(layoutID int) error {
	if err := ml.verifyLayoutID(layoutID); err != nil {
		return newLayoutIDError(layoutID, err)
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

func (ml *MultiplexedLayer) Muxor() *MuxorSignal {
	return ml.muxor
}

func (ml *MultiplexedLayer) Layouts() []*SL {
	return ml.layouts
}

func (ml *MultiplexedLayer) GetLayout(layoutID int) *SL {
	if err := ml.verifyLayoutID(layoutID); err != nil {
		return nil
	}
	return ml.layouts[layoutID]
}

func (ml *MultiplexedLayer) GetSignals(layoutID int) []Signal {
	if err := ml.verifyLayoutID(layoutID); err != nil {
		return nil
	}
	return ml.layouts[layoutID].Signals()
}

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

// verifySignalName checks if the signal name is already used in the multiplexed layer.
// It traverses all the multiplexed layers with BFS algorithm
// and checks if the name is already used.
func (ml *MultiplexedLayer) verifySignalName(name string) error {
	muxLayerQueue := collection.NewQueue[*MultiplexedLayer]()
	muxLayerQueue.Push(ml)

	visitedMuxLayers := collection.NewSet[EntityID]()

	for muxLayerQueue.Size() > 0 {
		muxLayer := muxLayerQueue.Pop()
		visitedMuxLayers.Add(muxLayer.muxor.entityID)

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

package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_signalPayload_append(t *testing.T) {
	assert := assert.New(t)

	payload := newSignalPayload(16)

	size4Type, err := NewIntegerSignalType("4_bits", "", 4, false)
	assert.NoError(err)
	size2Type, err := NewIntegerSignalType("2_bits", "", 2, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig3, err := NewStandardSignal("sig_3", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig4, err := NewStandardSignal("sig_4", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)

	// should get this payload after appending sig0, sig1, sig2, and sig3
	// 0 0 0 0 1 1 1 1 2 2 2 2 3 3 3 3
	assert.NoError(payload.append(sig0))
	assert.NoError(payload.append(sig1))
	assert.NoError(payload.append(sig2))
	assert.NoError(payload.append(sig3))

	expectedStartBits := []int{0, 4, 8, 12}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	payload.removeAll()

	// should get this payload after appending sig0, sig1, sig2, and sig4
	// 0 0 0 0 1 1 1 1 2 2 2 2 4 4 - -
	assert.NoError(payload.append(sig0))
	assert.NoError(payload.append(sig1))
	assert.NoError(payload.append(sig2))
	assert.NoError(payload.append(sig4))

	expectedStartBits = []int{0, 4, 8, 12}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	// should produce an error because sig3 will exceed the max payload size
	assert.Error(payload.append(sig3))
}

func Test_signalPayload_insert(t *testing.T) {
	assert := assert.New(t)

	payload := newSignalPayload(16)

	size2Type, err := NewIntegerSignalType("2_bits", "", 2, false)
	assert.NoError(err)
	size4Type, err := NewIntegerSignalType("4_bits", "", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig3, err := NewStandardSignal("sig_3", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig4, err := NewStandardSignal("sig_4", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)

	// should get this payload after inserting sig0, sig1, sig2, and sig3
	// 3 3 3 3 1 1 1 1 2 2 2 2 0 0 0 0
	assert.NoError(payload.insert(sig0, 12))
	assert.NoError(payload.insert(sig1, 4))
	assert.NoError(payload.insert(sig2, 8))
	assert.NoError(payload.insert(sig3, 0))

	expectedStartBits := []int{0, 4, 8, 12}
	expectedNames := []string{"sig_3", "sig_1", "sig_2", "sig_0"}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
		assert.Equal(expectedNames[idx], sig.Name())
	}

	// should return an error because there should be no space left
	assert.Error(payload.insert(sig4, 0))

	payload.removeAll()

	// should get this payload after inserting sig0, sig1, sig2, and sig4
	// 0 0 0 0 1 1 1 1 2 2 2 2 - - 4 4
	assert.NoError(payload.insert(sig0, 0))
	assert.NoError(payload.insert(sig1, 4))
	assert.NoError(payload.insert(sig2, 8))
	assert.NoError(payload.insert(sig4, 14))

	expectedStartBits = []int{0, 4, 8, 14}
	expectedNames = []string{"sig_0", "sig_1", "sig_2", "sig_4"}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
		assert.Equal(expectedNames[idx], sig.Name())
	}

	// should return error because there should be not enough space left
	assert.Error(payload.insert(sig3, 12))

	// should return error because there another signal is already starting at bit 8
	assert.Error(payload.insert(sig3, 8))

	payload.removeAll()

	// should return an error because sig0 of size 4 starting at 14 will exceed the payload size of 16
	assert.Error(payload.insert(sig0, 14))

	// should return an error because start bit is negative
	assert.Error(payload.insert(sig0, -10))
}

func Test_signalPayload_remove(t *testing.T) {
	assert := assert.New(t)

	payload := newSignalPayload(16)

	size4Type, err := NewIntegerSignalType("4_bits", "", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig3, err := NewStandardSignal("sig_3", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)

	// starting from this payload
	// 0 0 0 0 1 1 1 1 2 2 2 2 3 3 3 3
	assert.NoError(payload.insert(sig0, 0))
	assert.NoError(payload.insert(sig1, 4))
	assert.NoError(payload.insert(sig2, 8))
	assert.NoError(payload.insert(sig3, 12))

	// should get this after removing sig2
	// 0 0 0 0 1 1 1 1 - - - - 3 3 3 3
	payload.remove(sig2.EntityID())

	expectedStartBits := []int{0, 4, 12}
	expectedNames := []string{"sig_0", "sig_1", "sig_3"}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
		assert.Equal(expectedNames[idx], sig.Name())
	}

	// should remove all signals
	payload.remove(sig0.EntityID())
	payload.remove(sig1.EntityID())
	payload.remove(sig3.EntityID())

	assert.Equal(0, len(payload.signals))
}

func Test_signalPayload_compact(t *testing.T) {
	assert := assert.New(t)

	payload := newSignalPayload(16)

	size2Type, err := NewIntegerSignalType("2_bits", "", 2, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)
	sig3, err := NewStandardSignal("sig_3", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)

	// starting from this payload
	// - - 0 0 - - 1 1 - - 2 2 - - 3 3
	assert.NoError(payload.insert(sig0, 2))
	assert.NoError(payload.insert(sig1, 6))
	assert.NoError(payload.insert(sig2, 10))
	assert.NoError(payload.insert(sig3, 14))

	// should get this after compacting
	// 0 0 1 1 2 2 3 3 - - - - - - - -
	payload.compact()

	expectedStartBits := []int{0, 2, 4, 6}
	expectedNames := []string{"sig_0", "sig_1", "sig_2", "sig_3"}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
		assert.Equal(expectedNames[idx], sig.Name())
	}

	payload.removeAll()

	// starting from this payload
	// 0 0 - - - - 1 1 - - - 2 2 - 3 3
	assert.NoError(payload.insert(sig0, 0))
	assert.NoError(payload.insert(sig1, 6))
	assert.NoError(payload.insert(sig2, 11))
	assert.NoError(payload.insert(sig3, 14))

	// should get this after compacting
	// 0 0 1 1 2 2 3 3 - - - - - - - -
	payload.compact()

	expectedStartBits = []int{0, 2, 4, 6}
	expectedNames = []string{"sig_0", "sig_1", "sig_2", "sig_3"}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
		assert.Equal(expectedNames[idx], sig.Name())
	}
}

func Test_signalPayload_modifyStartBitsOnShrink(t *testing.T) {
	assert := assert.New(t)

	payload := newSignalPayload(16)

	size2Type, err := NewIntegerSignalType("2_bits", "", 2, false)
	assert.NoError(err)
	size4Type, err := NewIntegerSignalType("4_bits", "", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)
	sig3, err := NewStandardSignal("sig_3", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)

	// starting from this payload
	// 0 0 - - 1 1 1 1 - 2 2 - 3 3 - -
	assert.NoError(payload.insert(sig0, 0))
	assert.NoError(payload.insert(sig1, 4))
	assert.NoError(payload.insert(sig2, 9))
	assert.NoError(payload.insert(sig3, 12))

	// should get this one after shrinking sig1 by 1
	// 0 0 - - 1 1 1 - 2 2 - 3 3 - - -
	assert.NoError(payload.modifyStartBitsOnShrink(sig1, 1))

	expectedStartBits := []int{0, 4, 8, 11}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	// should get this one after shrinking sig0 by 1
	// 0 - - 1 1 1 - 2 2 - 3 3 - - - -
	assert.NoError(payload.modifyStartBitsOnShrink(sig0, 1))

	expectedStartBits = []int{0, 3, 7, 10}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	payload.removeAll()

	// starting from this payload
	// 0 0 1 1 1 1 2 2 3 3 - - - - - -
	assert.NoError(payload.insert(sig0, 0))
	assert.NoError(payload.insert(sig1, 2))
	assert.NoError(payload.insert(sig2, 6))
	assert.NoError(payload.insert(sig3, 8))

	// should get this one after shrinking sig1 by 2
	// 0 0 1 1 2 2 3 3 - - - - - - - -
	assert.NoError(payload.modifyStartBitsOnShrink(sig1, 2))

	expectedStartBits = []int{0, 2, 4, 6}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	payload.removeAll()

	// starting from this payload
	// - - - - - - - - - - 0 0 - - - -
	assert.NoError(payload.insert(sig0, 10))

	// should get this one after shrinking sig0 by 1
	// - - - - - - - - - - 0 - - - - -
	assert.NoError(payload.modifyStartBitsOnShrink(sig0, 1))
	assert.Equal(10, payload.signals[0].GetStartBit())

	// should do nothing
	assert.NoError(payload.modifyStartBitsOnShrink(sig0, 0))
	assert.Equal(10, payload.signals[0].GetStartBit())

	// should return an error because amount is negative
	assert.Error(payload.modifyStartBitsOnShrink(sig0, -1))

	payload.removeAll()

	// starting from this payload
	// 0 0 - - - - - - - - - - - - - -
	assert.NoError(payload.append(sig0))

	// should return an error because amount is greater then the signal size
	assert.Error(payload.modifyStartBitsOnShrink(sig0, 3))

	// should return an error because amount is greater equal to the signal size
	assert.Error(payload.modifyStartBitsOnShrink(sig0, 2))
}

func Test_signalPayload_modifyStartBitsOnGrow(t *testing.T) {
	assert := assert.New(t)

	payload := newSignalPayload(16)

	size2Type, err := NewIntegerSignalType("2_bits", "", 2, false)
	assert.NoError(err)
	size4Type, err := NewIntegerSignalType("4_bits", "", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig2, err := NewStandardSignal("sig_2", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)
	sig3, err := NewStandardSignal("sig_3", "", size2Type, 0, 3, 0, 1, nil)
	assert.NoError(err)

	// starting from this payload
	// 0 0 - - 1 1 1 1 - 2 2 - 3 3 - -
	assert.NoError(payload.insert(sig0, 0))
	assert.NoError(payload.insert(sig1, 4))
	assert.NoError(payload.insert(sig2, 9))
	assert.NoError(payload.insert(sig3, 12))

	// should get this one after growing sig0 by 1
	// 0 0 0 - 1 1 1 1 1 1 1 2 2 3 3 -
	assert.NoError(payload.modifyStartBitsOnGrow(sig0, 1))

	expectedStartBits := []int{0, 4, 9, 12}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	payload.removeAll()

	// starting from this payload
	// 0 0 - - 1 1 1 1 - 2 2 - 3 3 - -
	assert.NoError(payload.insert(sig0, 0))
	assert.NoError(payload.insert(sig1, 4))
	assert.NoError(payload.insert(sig2, 9))
	assert.NoError(payload.insert(sig3, 12))

	// should get this one after growing sig1 by 3
	// 0 0 - - 1 1 1 1 1 1 1 2 2 3 3 -
	assert.NoError(payload.modifyStartBitsOnGrow(sig1, 3))

	expectedStartBits = []int{0, 4, 11, 13}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	payload.removeAll()

	// starting from this payload
	// 0 0 1 1 1 1 2 2 3 3 - - - - - -
	assert.NoError(payload.insert(sig0, 0))
	assert.NoError(payload.insert(sig1, 2))
	assert.NoError(payload.insert(sig2, 6))
	assert.NoError(payload.insert(sig3, 8))

	// should get this one after growing sig0 by 6
	// 0 0 0 0 0 0 0 0 1 1 1 1 2 2 3 3
	assert.NoError(payload.modifyStartBitsOnGrow(sig0, 6))

	expectedStartBits = []int{0, 8, 12, 14}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	payload.removeAll()

	// starting from this payload
	// 0 0 1 1 1 1 2 2 3 3 - - - - - -
	assert.NoError(payload.insert(sig0, 0))
	assert.NoError(payload.insert(sig1, 2))
	assert.NoError(payload.insert(sig2, 6))
	assert.NoError(payload.insert(sig3, 8))

	// should get an error after trying to grow sig0 by 8
	assert.Error(payload.modifyStartBitsOnGrow(sig0, 8))

	// should get an error after trying to grow sig1 by 8
	assert.Error(payload.modifyStartBitsOnGrow(sig1, 8))

	// should get an error after trying to grow sig2 by 8
	assert.Error(payload.modifyStartBitsOnGrow(sig2, 8))

	payload.removeAll()

	// starting from this payload
	// 0 0 - - - - - - - - - - - - - -
	assert.NoError(payload.insert(sig0, 0))

	// should get this one after growing sig0 by 6
	// 0 0 0 0 0 0 - - - - - - - - - -
	assert.NoError(payload.modifyStartBitsOnGrow(sig0, 4))

	// should get the same as before because grow by 0 has no effect
	assert.NoError(payload.modifyStartBitsOnGrow(sig0, 0))

	// should get an error because cannot grow by a negative amount
	assert.Error(payload.modifyStartBitsOnGrow(sig0, -1))
}

func Test_signalPayload_shiftLeft(t *testing.T) {
	assert := assert.New(t)

	payload := newSignalPayload(16)

	size4Type, err := NewIntegerSignalType("4_bits", "", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)

	// starting from this payload
	// - - 0 0 0 0 1 1 1 1 - - - - - -
	assert.NoError(payload.insert(sig0, 2))
	assert.NoError(payload.insert(sig1, 6))

	// should get this one after shifting sig0 by 2
	// 0 0 0 0 - - 1 1 1 1 - - - - - -
	// and should get 2 as result
	assert.Equal(2, payload.shiftLeft(sig0, 2))

	expectedStartBits := []int{0, 6}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	// should get this one after shifting sig1 by 4
	// 0 0 0 0 1 1 1 1 - - - - - - - -
	// and should get 2 as result
	assert.Equal(2, payload.shiftLeft(sig1, 4))

	expectedStartBits = []int{0, 4}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	payload.removeAll()

	// starting from this payload
	// - - 0 0 0 0 1 1 1 1 - - - - - -
	assert.NoError(payload.insert(sig0, 2))
	assert.NoError(payload.insert(sig1, 6))

	// should get this one after shifting sig0 by 4
	// 0 0 0 0 - - 1 1 1 1 - - - - - -
	// and should get 2 as result
	assert.Equal(2, payload.shiftLeft(sig0, 4))

	expectedStartBits = []int{0, 6}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	// should get 0 as result
	assert.Equal(0, payload.shiftLeft(sig1, 0))

	payload.removeAll()

	// starting from this payload
	// - - - - - - - - - - - - 0 0 0 0
	assert.NoError(payload.insert(sig0, 12))

	// should get this one after shifting sig0 by 4
	// - - - - - - - - 0 0 0 0 - - - -
	// and should get 4 as result
	assert.Equal(4, payload.shiftLeft(sig0, 4))

	assert.Equal(8, payload.signals[0].GetStartBit())
}

func Test_signalPayload_shiftRight(t *testing.T) {
	assert := assert.New(t)

	payload := newSignalPayload(16)

	size4Type, err := NewIntegerSignalType("4_bits", "", 4, false)
	assert.NoError(err)

	sig0, err := NewStandardSignal("sig_0", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)
	sig1, err := NewStandardSignal("sig_1", "", size4Type, 0, 15, 0, 1, nil)
	assert.NoError(err)

	// starting from this payload
	// - - - - - - 0 0 0 0 1 1 1 1 - -
	assert.NoError(payload.insert(sig0, 6))
	assert.NoError(payload.insert(sig1, 10))

	// should get this one after shifting sig1 by 2
	// - - - - - - 0 0 0 0 - - 1 1 1 1
	// and should get 2 as result
	assert.Equal(2, payload.shiftRight(sig1, 2))

	expectedStartBits := []int{6, 12}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	// should get this one after shifting sig0 by 4
	// - - - - - - - - 0 0 0 0 1 1 1 1
	// and should get 2 as result
	assert.Equal(2, payload.shiftRight(sig0, 4))

	expectedStartBits = []int{8, 12}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	payload.removeAll()

	// starting from this payload
	// - - - - - - 0 0 0 0 1 1 1 1 - -
	assert.NoError(payload.insert(sig0, 6))
	assert.NoError(payload.insert(sig1, 10))

	// should get this one after shifting sig1 by 4
	// - - - - - - 0 0 0 0 - - 1 1 1 1
	// and should get 2 as result
	assert.Equal(2, payload.shiftRight(sig1, 4))

	expectedStartBits = []int{6, 12}
	for idx, sig := range payload.signals {
		assert.Equal(expectedStartBits[idx], sig.GetStartBit())
	}

	// should get 0 as result
	assert.Equal(0, payload.shiftRight(sig1, 0))

	payload.removeAll()

	// starting from this payload
	// 0 0 0 0 - - - - - - - - - - - -
	assert.NoError(payload.insert(sig0, 0))

	// should get this one after shifting sig0 by 4
	// - - - - 0 0 0 0 - - - - - - - -
	// and should get 4 as result
	assert.Equal(4, payload.shiftRight(sig0, 4))

	assert.Equal(4, payload.signals[0].GetStartBit())
}

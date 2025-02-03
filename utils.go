package acmelib

import "golang.org/x/exp/slices"

// MessageLoad is a struct that represents the load caused by a [Message] in a [Bus].
type MessageLoad struct {
	// Message is the examined message.
	Message *Message
	// BitsPerSec is the number of bits occupied by the message per second.
	BitsPerSec float64
	// Percentage is the load in percentage relative to the total number of bits sent in the bus.
	Percentage float64
}

// CalculateBusLoad returns the estimed load of the given [Bus] in the worst case scenario.
// It returns the load percentage and a slice of [MessageLoad] structs sorted from the message
// that causes the most load to the message that causes the least load.
// The default cycle time is used when a message within the bus does not have one set.
// If the bus does not have the baudrate set, it returns 0.
//
// It returns an [ArgumentError] if the given default cycle time is invalid.
func CalculateBusLoad(bus *Bus, defCycleTime int) (float64, []*MessageLoad, error) {
	msgLoads := []*MessageLoad{}

	if defCycleTime < 0 {
		return 0, msgLoads, &ArgumentError{
			Name: "defCycleTime",
			Err:  ErrIsNegative,
		}
	}

	if defCycleTime == 0 {
		return 0, msgLoads, &ArgumentError{
			Name: "defCycleTime",
			Err:  ErrIsZero,
		}
	}

	if bus.baudrate == 0 {
		return 0, msgLoads, nil
	}

	var headerBits int
	var trailerBits int
	var headerStuffingBits int
	switch bus.typ {
	case BusTypeCAN2A:
		// start of frame + id + rtr + ide + r0 + dlc
		headerBits = 19
		// crc + delim crc + slot ack + delim ack + eof
		trailerBits = 25
		// from bit stuffing section of wikipedia (https://en.wikipedia.org/wiki/CAN_bus#Bit_stuffing)
		headerStuffingBits = 34
	}

	totConsumedBitsPerSec := float64(0)
	for _, tmpInt := range bus.nodeInts.getValues() {
		for _, tmpMsg := range tmpInt.sentMessages.getValues() {
			stuffingBits := (headerStuffingBits + tmpMsg.sizeByte*8 - 1) / 4
			msgBits := tmpMsg.sizeByte*8 + headerBits + trailerBits + stuffingBits

			cycleTime := tmpMsg.cycleTime
			if cycleTime == 0 {
				cycleTime = defCycleTime
			}

			msgBitsPerSec := float64(msgBits) / float64(cycleTime) * 1000
			totConsumedBitsPerSec += msgBitsPerSec

			msgLoads = append(msgLoads, &MessageLoad{
				Message:    tmpMsg,
				BitsPerSec: msgBitsPerSec,
			})
		}
	}

	for _, tmpMsgLoad := range msgLoads {
		tmpMsgLoad.Percentage = tmpMsgLoad.BitsPerSec / totConsumedBitsPerSec * 100
	}

	slices.SortFunc(msgLoads, func(a, b *MessageLoad) int {
		diff := b.BitsPerSec - a.BitsPerSec
		return int(diff)
	})

	return totConsumedBitsPerSec / float64(bus.baudrate) * 100, msgLoads, nil
}

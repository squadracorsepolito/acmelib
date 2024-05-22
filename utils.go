package acmelib

// CalculateBusLoad returns the estimed load of the given [Bus] in the worst case scenario.
// The default cycle time is used when a message within the bus does not have one set.
// If the bus does not have the baudrate set, it returns 0.
//
// It returns an [ArgumentError] if the given default cycle time is invalid.
func CalculateBusLoad(bus *Bus, defCycleTime int) (float64, error) {
	if defCycleTime < 0 {
		return 0, &ArgumentError{
			Name: "defCycleTime",
			Err:  ErrIsNegative,
		}
	}

	if defCycleTime == 0 {
		return 0, &ArgumentError{
			Name: "defCycleTime",
			Err:  ErrIsZero,
		}
	}

	if bus.baudrate == 0 {
		return 0, nil
	}

	var headerBits int
	var trailerBits int
	var stuffingBits int
	switch bus.typ {
	case BusTypeCAN2A:
		// start of frame + id + rtr + ide + r0 + dlc
		headerBits = 19
		// crc + delim crc + slot ack + delim ack + eof
		trailerBits = 25
		// worst case scenario
		stuffingBits = 19
	}

	consumedBitsPerSec := float64(0)
	for _, tmpInt := range bus.nodeInts.getValues() {
		for _, tmpMsg := range tmpInt.messages.getValues() {
			msgBits := tmpMsg.sizeByte*8 + headerBits + trailerBits + stuffingBits

			var cycleTime int
			if tmpMsg.cycleTime == 0 {
				cycleTime = defCycleTime
			} else {
				cycleTime = tmpMsg.cycleTime
			}

			consumedBitsPerSec += float64(msgBits) / float64(cycleTime) * 1000
		}
	}

	return consumedBitsPerSec / float64(bus.baudrate) * 100, nil
}

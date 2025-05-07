package acmelib

import (
	"math"
	"strings"
)

const maxSize = 64

func getSizeFromValue(val int) int {
	if val == 0 {
		return 1
	}

	for i := range maxSize {
		if val < 1<<i {
			return i
		}
	}

	return maxSize
}

func getSizeFromCount(count int) int {
	if count == 0 {
		return 0
	}

	for i := range maxSize {
		if count <= 1<<i {
			return i
		}
	}

	return maxSize
}

func getValueFromSize(size int) int {
	if size <= 0 {
		return 1
	}
	return 1 << size
}

func getMinMaxFromSize(size int, signed bool) (float64, float64) {
	if size <= 0 {
		return 0, 0
	}

	if signed {
		min := -(1 << (size - 1))
		max := (1 << (size - 1)) - 1
		return float64(min), float64(max)
	}

	max := (1 << size) - 1
	return 0, float64(max)
}

func isDecimal(val float64) bool {
	return math.Mod(val, 1.0) != 0
}

func clearSpaces(str string) string {
	return strings.ReplaceAll(strings.TrimSpace(str), " ", "_")
}

// StartPosFromBigEndian converts the big endian start bit to a little endian start bit.
// Since the library uses little endian for storing and validating signals, this conversion function
// may be useful when you want to insert a signal into a [Message] or into a [MultiplexedLayer],
// and you only have the big endian start bit.
func StartPosFromBigEndian(bigEndianStartBit int) int {
	return bigEndianStartBit + 7 - 2*(bigEndianStartBit%8)
}

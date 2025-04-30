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

func isDecimal(val float64) bool {
	return math.Mod(val, 1.0) != 0
}

func clearSpaces(str string) string {
	return strings.ReplaceAll(strings.TrimSpace(str), " ", "_")
}

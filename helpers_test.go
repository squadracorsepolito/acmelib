package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calcSizeFromValue(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, getSizeFromValue(1))
	assert.Equal(2, getSizeFromValue(3))
	assert.Equal(2, getSizeFromValue(2))
	assert.Equal(6, getSizeFromValue(32))
	assert.Equal(7, getSizeFromValue(127))
	assert.Equal(9, getSizeFromValue(256))
}

func Test_calcValueFromSize(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, calcValueFromSize(0))
	assert.Equal(256, calcValueFromSize(8))
	assert.Equal(32, calcValueFromSize(5))
	assert.Equal(8, calcValueFromSize(3))
}

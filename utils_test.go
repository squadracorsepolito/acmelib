package acmelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calcSizeFromValue(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, calcSizeFromValue(1))
	assert.Equal(2, calcSizeFromValue(3))
	assert.Equal(2, calcSizeFromValue(2))
	assert.Equal(6, calcSizeFromValue(32))
	assert.Equal(7, calcSizeFromValue(127))
	assert.Equal(9, calcSizeFromValue(256))
}

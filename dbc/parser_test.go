package dbc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser_parseExtendedMuxRange(t *testing.T) {
	assert := assert.New(t)

	file := "10-10"

	p := NewParser("test_filename", []byte(file))

	extMuxRange, err := p.parseExtendedMuxRange()
	assert.NoError(err)

	t.Log(extMuxRange)
}

package dbc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_scanner_scanNumber(t *testing.T) {
	assert := assert.New(t)

	file := `1 - + -1 +1
		1+ 1- 1-1 1-1-
		0x1 0x1- 0x`

	expectedTokens := []tokenKind{tokenNumber, tokenPunct, tokenPunct, tokenNumber, tokenNumber,
		tokenNumber, tokenPunct, tokenNumber, tokenPunct, tokenNumberRange, tokenNumberRange, tokenPunct,
		tokenNumber, tokenNumber, tokenPunct, tokenError}

	s := newScanner(bytes.NewBuffer([]byte(file)), file)
	token := s.scan()
	idx := 0
	for !token.isEOF() {
		if !token.isSpace() {
			assert.Equal(expectedTokens[idx], token.kind)
			idx++
		}

		token = s.scan()
	}
}

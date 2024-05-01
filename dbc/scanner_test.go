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
		0x1 0x1- 0x
		32@0- @0-`

	expectedTokens := []tokenKind{tokenNumber, tokenPunct, tokenPunct, tokenNumber, tokenNumber,
		tokenNumber, tokenPunct, tokenNumber, tokenPunct, tokenNumberRange, tokenNumberRange, tokenPunct,
		tokenNumber, tokenNumber, tokenPunct, tokenError,
		tokenNumber, tokenPunct, tokenNumber, tokenPunct, tokenPunct, tokenNumber, tokenPunct,
	}

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

func Test_scanner_scanSignal(t *testing.T) {
	assert := assert.New(t)

	file := `SG_ IVTMain_Result_W : 23|32@0- (1,0) [-2147483648|2147483647] "W" HVB`

	expectedTokens := []tokenKind{tokenKeyword, tokenIdent, tokenPunct, tokenNumber, tokenPunct,
		tokenNumber, tokenPunct, tokenNumber, tokenPunct, tokenPunct, tokenNumber, tokenPunct, tokenNumber, tokenPunct,
		tokenPunct, tokenNumber, tokenPunct, tokenNumber, tokenPunct, tokenString, tokenIdent,
	}

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

package dbc

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_scanner_scan(t *testing.T) {
	assert := assert.New(t)

	expectedTokenValues := []string{
		"ident", "ident0", "ident_0",
		"1", "-1", "1.0", "-1.0", "0x1", "0X1",
		"0-0", "0-1", "10-10",
		"M", "m0", "m0M",
		"\"string 0\"", "\"string 🚀\"", "\"string °\"",
	}

	expectedTokenKinds := []tokenKind{
		tokenIdent, tokenIdent, tokenIdent,
		tokenNumber, tokenNumber, tokenNumber, tokenNumber, tokenNumber, tokenNumber,
		tokenNumberRange, tokenNumberRange, tokenNumberRange,
		tokenMuxIndicator, tokenMuxIndicator, tokenMuxIndicator,
		tokenString, tokenString, tokenString,
	}

	testStr := ""
	for _, tmpStr := range expectedTokenValues {
		testStr += "\n" + tmpStr
	}

	s := newScanner(bytes.NewBufferString(testStr))

	token := s.scan()
	idx := 0
	for !token.isEOF() {
		if !token.isSpace() {
			assert.Equal(expectedTokenKinds[idx], token.kind)
			assert.Equal(strings.ReplaceAll(expectedTokenValues[idx], "\"", ""), token.value)
			idx++
		}

		token = s.scan()
	}
}

func Test_scanner_scanNumber(t *testing.T) {
	assert := assert.New(t)

	file := `1 - + -1 +1
		1+ 1- 1-1 1-1-
		0x1 0x1-
		32@0- @0-
	`

	expectedTokens := []tokenKind{tokenNumber, tokenPunct, tokenPunct, tokenNumber, tokenNumber,
		tokenNumber, tokenPunct, tokenNumber, tokenPunct, tokenNumberRange, tokenNumberRange, tokenPunct,
		tokenNumber, tokenNumber, tokenPunct,
		tokenNumber, tokenPunct, tokenNumber, tokenPunct, tokenPunct, tokenNumber, tokenPunct,
	}

	s := newScanner(bytes.NewBufferString(file))
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

	s := newScanner(bytes.NewBufferString(file))
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

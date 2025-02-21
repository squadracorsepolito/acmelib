package dbc

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxErrorValueLength = 20

const eof = rune(0)

func isEOF(ch rune) bool {
	return ch == eof
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isNumber(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isHexNumber(ch rune) bool {
	return isNumber(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isAlphaNumeric(ch rune) bool {
	return isLetter(ch) || isNumber(ch) || ch == '_' || ch == '-'
}

type scanner struct {
	r *bufio.Reader

	value           string
	peekBytesOffset int

	lastReadCh rune

	beginToken bool

	currLine  int
	startLine int

	currCol  int
	startCol int
}

func newScanner(r io.Reader) *scanner {
	bufR := bufio.NewReader(r)

	return &scanner{
		r: bufR,

		value:           "",
		peekBytesOffset: 1,

		beginToken: true,

		currLine:  1,
		startLine: 1,

		currCol:  0,
		startCol: 0,
	}
}

func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}

	s.value += string(ch)
	s.peekBytesOffset = 0

	s.lastReadCh = ch

	s.currCol++
	if ch == '\t' {
		s.currCol += 4
	}

	if ch == '\n' {
		s.currLine++
		s.currCol = 0
	}

	if s.beginToken {
		s.startLine = s.currLine
		s.startCol = s.currCol
		s.beginToken = false
	}

	return ch
}

func (s *scanner) peek() rune {
	for peekBytes := 1; peekBytes < 4; peekBytes++ {
		b, err := s.r.Peek(peekBytes + s.peekBytesOffset)
		if err == nil {
			ch, chBytes := utf8.DecodeRune(b[s.peekBytesOffset:])
			if ch == utf8.RuneError {
				continue
			}

			s.peekBytesOffset += chBytes

			return ch
		}
	}

	return eof
}

func (s *scanner) emitToken(kind tokenKind) *token {
	val := s.value
	if kind == tokenString {
		val = s.value[1 : len(s.value)-1]
	}

	t := &token{
		kind:      kind,
		kindName:  tokenNames[kind],
		value:     val,
		startLine: s.startLine,
		startCol:  s.startCol,
		endLine:   s.currLine + 1,
		endCol:    s.currCol + 1,
	}

	s.value = ""

	s.beginToken = true

	return t
}

func (s *scanner) emitErrorToken(msg string) *token {
	val := ""
	if len(s.value) > maxErrorValueLength {
		val = fmt.Sprintf("%s : %s", msg, s.value[:maxErrorValueLength])
	} else {
		val = fmt.Sprintf("%s : %s", msg, s.value)
	}

	t := &token{
		kind:      tokenError,
		kindName:  tokenNames[tokenError],
		value:     val,
		startLine: s.startLine,
		startCol:  s.startCol,
		endLine:   s.currLine + 1,
		endCol:    s.currCol + 1,
	}

	s.value = ""

	s.beginToken = true

	return t
}

func (s *scanner) scan() *token {
	switch ch := s.read(); {
	case isEOF(ch):
		return s.emitToken(tokenEOF)

	case isSpace(ch):
		return s.scanSpace()

	case isLetter(ch):
		return s.scanText()

	case isNumber(ch) || ch == '-' || ch == '+':
		return s.scanNumber()

	case ch == '"':
		return s.scanString()

	case isPunctKeyword(ch):
		return s.emitToken(tokenPunct)
	}

	return s.emitErrorToken("unrecognized symbol")
}

func (s *scanner) scanText() *token {
	firstCh := s.lastReadCh

	buf := new(strings.Builder)
	buf.WriteRune(firstCh)

	isMuxSwitch := false
	foundSwitchNum := false
	if firstCh == 'm' {
		isMuxSwitch = true
	}

loop:
	for {
		switch ch := s.peek(); {
		case isEOF(ch):
			break loop

		case isAlphaNumeric(ch):
			if isMuxSwitch {
				if isNumber(ch) {
					foundSwitchNum = true
				} else if !foundSwitchNum || ch != 'M' {
					isMuxSwitch = false
				}
			}
			buf.WriteRune(ch)
			s.read()

		default:
			break loop
		}
	}

	if (isMuxSwitch && buf.Len() > 1) || buf.Len() == 1 && firstCh == 'M' {
		return s.emitToken(tokenMuxIndicator)
	}

	if _, ok := keywords[buf.String()]; ok {
		return s.emitToken(tokenKeyword)
	}

	return s.emitToken(tokenIdent)
}

func (s *scanner) scanSpace() *token {
	ch := s.peek()
	for isSpace(ch) {
		s.read()
		ch = s.peek()
	}
	return s.emitToken(tokenSpace)
}

func (s *scanner) scanNumber() *token {
	firstCh := s.lastReadCh
	prevCh := firstCh
	hasMore := false
	isRange := false

loop:
	for {
		switch ch := s.peek(); {
		case isEOF(ch):
			break loop

		case firstCh == '0' && (ch == 'x' || ch == 'X'):
			return s.scanHexNumber()

		case !isNumber(ch) && ch != '.':
			if ch == 'e' || ch == 'E' {
				if prevCh != '-' && prevCh != '+' && prevCh != '.' {
					return s.scanExpNumber()
				}
			}

			if ch == '-' && isNumber(firstCh) && !isRange {
				nextCh := s.peek()
				if isNumber(nextCh) {
					s.read()
					s.read()
					isRange = true
					continue loop
				}
			}

			break loop

		default:
			if ch == '.' && (prevCh == '-' || prevCh == '+') {
				break loop
			}

			hasMore = true
			s.read()
			prevCh = ch
		}
	}

	if !hasMore {
		if firstCh == '-' || firstCh == '+' {
			return s.emitToken(tokenPunct)
		}
	}

	if isRange {
		return s.emitToken(tokenNumberRange)
	}

	return s.emitToken(tokenNumber)
}

func (s *scanner) scanHexNumber() *token {
	ch := s.peek()
	if !isHexNumber(ch) {
		return s.emitErrorToken("invalid hex number")
	}

	s.read()
	s.read()

	for i := 0; i < 8; i++ {
		ch = s.peek()

		if !isHexNumber(ch) {
			break
		}

		s.read()
	}

	return s.emitToken(tokenNumber)
}

func (s *scanner) scanExpNumber() *token {
	ch := s.peek()
	if ch != '-' && ch != '+' && !isNumber(ch) {
		s.read()
		return s.emitErrorToken("invalid exponential number")
	}

	if ch == '-' || ch == '+' {
		foundExp := false
		for {
			ch = s.peek()

			if !isNumber(ch) {
				if foundExp {
					break
				}

				s.read()
				s.read()
				return s.emitErrorToken("invalid exponential number")
			}

			if !foundExp {
				foundExp = true
				s.read()
			}

			s.read()
		}

		return s.emitToken(tokenNumber)
	}

	s.read()
	for {
		ch = s.peek()
		if !isNumber(ch) {
			break
		}
		s.read()
	}

	return s.emitToken(tokenNumber)
}

func (s *scanner) scanString() *token {
	for {
		switch ch := s.read(); {
		case isEOF(ch):
			return s.emitErrorToken(`unclosed string, missing closing "`)

		case ch == '"':
			return s.emitToken(tokenString)
		}
	}
}

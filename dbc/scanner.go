package dbc

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const eof = rune(0)

func isEOF(ch rune) bool {
	return ch == eof
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
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
	r     *bufio.Reader
	input string

	pos   int
	value string
}

func newScanner(r io.Reader, str string) *scanner {
	return &scanner{
		r:     bufio.NewReader(r),
		input: str,

		pos:   0,
		value: "",
	}
}

func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}

	s.pos++

	s.value += string(ch)

	return ch
}

func (s *scanner) unread() {
	_ = s.r.UnreadRune()

	s.pos--

	s.value = s.value[:len(s.value)-1]
}

func (s *scanner) getPosition() (int, int) {
	col := 1
	line := 1

	for _, ch := range s.input[:s.pos-len(s.value)] {
		if ch == '\n' {
			col = 1
			line++
			continue
		}
		col++
	}

	return col, line
}

func (s *scanner) emitToken(kind tokenKind) *token {
	val := s.value
	if kind == tokenString {
		val = s.value[1 : len(s.value)-1]
	}

	col, line := s.getPosition()

	t := &token{
		kind:     kind,
		kindName: tokenNames[kind],
		value:    val,
		col:      col,
		line:     line,
	}

	s.value = ""

	return t
}

const maxErrorValueLength = 20

func (s *scanner) emitErrorToken(msg string) *token {
	val := ""
	if len(s.value) > maxErrorValueLength {
		val = fmt.Sprintf("%s : %s", msg, s.value[:maxErrorValueLength])
	} else {
		val = fmt.Sprintf("%s : %s", msg, s.value)
	}

	col, line := s.getPosition()

	t := &token{
		kind:     tokenError,
		kindName: tokenNames[tokenError],
		value:    val,
		col:      col,
		line:     line,
	}

	s.value = ""

	return t
}

func (s *scanner) scan() *token {
	switch ch := s.read(); {
	case isEOF(ch):
		return s.emitToken(tokenEOF)

	case isSpace(ch):
		return s.scanSpace()

	case isLetter(ch):
		s.unread()
		return s.scanText()

	case isNumber(ch) || ch == '-' || ch == '+':
		s.unread()
		return s.scanNumber()

	case ch == '"':
		return s.scanString()

	case isSyntaxKeyword(ch):
		return s.emitToken(tokenSyntax)
	}

	return s.emitErrorToken("unrecognized symbol")
}

func (s *scanner) scanText() *token {
	firstCh := s.read()

	buf := new(strings.Builder)
	buf.WriteRune(firstCh)

	isMuxSwitch := false
	foundSwitchNum := false
	if firstCh == 'm' {
		isMuxSwitch = true
	}

loop:
	for {
		switch ch := s.read(); {
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

		default:
			s.unread()
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
	ch := s.read()
	for isSpace(ch) {
		ch = s.read()
	}
	s.unread()
	return s.emitToken(tokenSpace)
}

func (s *scanner) scanNumber() *token {
	firstCh := s.read()
	hasMore := false
	isRange := false

loop:
	for {
		switch ch := s.read(); {
		case isEOF(ch):
			break loop

		case firstCh == '0' && (ch == 'x' || ch == 'X'):
			return s.scanHexNumber()

		case !isNumber(ch) && ch != '.':
			if ch == '-' && isNumber(firstCh) && !isRange {
				isRange = true
				continue
			}
			s.unread()
			break loop

		default:
			hasMore = true
		}
	}

	if !hasMore {
		if firstCh == '-' || firstCh == '+' {
			return s.emitToken(tokenSyntax)
		}
	}

	if isRange {
		return s.emitToken(tokenNumberRange)
	}

	return s.emitToken(tokenNumber)
}

func (s *scanner) scanHexNumber() *token {
	if !isHexNumber(s.read()) {
		return s.emitErrorToken("invalid hex number")
	}

	for i := 0; i < 8; i++ {
		ch := s.read()

		if !isHexNumber(ch) {
			s.unread()
			break
		}
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

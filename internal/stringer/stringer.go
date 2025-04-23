package stringer

import (
	"fmt"
	"strings"
)

type Stringer struct {
	b    *strings.Builder
	tabs int
}

func New() *Stringer {
	return &Stringer{
		b:    new(strings.Builder),
		tabs: 0,
	}
}

func (s *Stringer) String() string {
	return s.b.String()
}

func (s *Stringer) Write(format string, args ...any) {
	tmpArgs := make([]any, 0, len(args)+1)
	tmpArgs = append(tmpArgs, strings.Repeat("\t", s.tabs))
	for _, arg := range args {
		tmpArgs = append(tmpArgs, arg)
	}

	s.b.WriteString(fmt.Sprintf("%s"+format, tmpArgs...))
}

func (s *Stringer) NewLine() {
	s.b.WriteRune('\n')
}

func (s *Stringer) Indent() {
	s.tabs++
}

func (s *Stringer) Unindent() {
	s.tabs--
}

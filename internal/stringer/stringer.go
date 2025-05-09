// Package stringer provides a simple string builder with support for indentation.
package stringer

import (
	"fmt"
	"strings"
)

// Stringer is a simple string builder with support for indentation.
type Stringer struct {
	b    *strings.Builder
	tabs int
}

// New creates a new [Stringer].
func New() *Stringer {
	return &Stringer{
		b:    new(strings.Builder),
		tabs: 0,
	}
}

func (s *Stringer) String() string {
	return s.b.String()
}

// Write writes a formatted string to the stringer.
func (s *Stringer) Write(format string, args ...any) {
	tmpArgs := make([]any, 0, len(args)+1)
	tmpArgs = append(tmpArgs, strings.Repeat("\t", s.tabs))
	for _, arg := range args {
		tmpArgs = append(tmpArgs, arg)
	}

	s.b.WriteString(fmt.Sprintf("%s"+format, tmpArgs...))
}

// NewLine writes a newline to the stringer.
func (s *Stringer) NewLine() {
	s.b.WriteRune('\n')
}

// Indent increases the indentation level.
func (s *Stringer) Indent() {
	s.tabs++
}

// Unindent decreases the indentation level.
func (s *Stringer) Unindent() {
	s.tabs--
}

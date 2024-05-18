package syntax

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Position struct {
	file *string // filename
	Line int32   // 1-based line number; 0 if line unknown
	Col  int32   // 1-based column (rune) number; 0 if column unknown
}

func (p Position) add(s string) Position {
	if n := strings.Count(s, "\n"); n > 0 {
		p.Line += int32(n)
		s = s[strings.LastIndex(s, "\n")+1:]
		p.Col = 1
	}
	p.Col += int32(utf8.RuneCountInString(s))
	return p
}

func (p Position) Filename() string {
	if p.file != nil {
		return *p.file
	}
	return "<invalid>"
}

func (p Position) String() string {
	file := p.Filename()
	if p.Line > 0 {
		if p.Col > 0 {
			return fmt.Sprintf("%s:%d:%d", file, p.Line, p.Col)
		}
		return fmt.Sprintf("%s:%d", file, p.Line)
	}
	return file
}

// MakePosition returns position with the specified components.
func MakePosition(file *string, line, col int32) Position {
	return Position{file, line, col}
}

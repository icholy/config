package token

import (
	"bufio"
	"fmt"
	"strings"
)

// Pos is the position inside the file
type Pos struct {
	Line, Column, Offset int
}

// String returns the line and column as a string
func (p Pos) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Span is a range of input used to display errors
type Span struct {
	Start Pos
	Lines []string
}

func (s Span) String() string {
	var b strings.Builder
	for i := 0; i < s.Start.Column-1; i++ {
		b.WriteRune(' ')
	}
	for i, line := range s.Lines {
		fmt.Fprintf(&b, "%05d: %s\n", s.Start.Line+i, line)
	}
	return b.String()
}

// Snip returns a snippet between start and end
func Snip(input string, start, end Pos) Span {
	var lines []string
	input = string([]rune(input)[start.Offset:end.Offset])
	sc := bufio.NewScanner(strings.NewReader(input))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return Span{
		Start: start,
		Lines: lines,
	}
}

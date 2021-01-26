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
	for i, line := range s.Lines {
		var padding string
		if i == 0 {
			padding = strings.Repeat(" ", s.Start.Column-1)
		}
		fmt.Fprintf(&b, "%05d: %s%s\n", s.Start.Line+i, padding, line)
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

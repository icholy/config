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

// Snippet is a range of input used to display errors
type Snippet struct {
	Start Pos
	Lines []string
}

func (s Snippet) String() string {
	var b strings.Builder
	for i := 0; i < s.Start.Column-1; i++ {
		b.WriteRune(' ')
	}
	for i, line := range s.Lines {
		fmt.Fprintf(&b, "%02d: %s\n", s.Start.Line+i, line)
	}
	return b.String()
}

// Snip returns a snippet between start and end
func Snip(input string, start, end Pos) Snippet {
	var lines []string
	input = string([]rune(input)[start.Offset:end.Offset])
	sc := bufio.NewScanner(strings.NewReader(input))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return Snippet{
		Start: start,
		Lines: lines,
	}
}

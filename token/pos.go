package token

import "fmt"

// Pos is the position inside the file
type Pos struct {
	Line, Column, Offset int
}

// String returns the line and column as a string
func (p Pos) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

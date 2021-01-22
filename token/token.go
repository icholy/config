package token

import (
	"bufio"
	"io"
)

// Type is the token type
type Type int

const (
	INVALID Type = iota
	EOF
	IDENT
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	STRING
	COMMA
	ASSIGN
)

// Pos is the position inside the file
type Pos struct {
	Line, Column int
}

// Token contains a token's type and text
type Token struct {
	Pos  Pos
	Type Type
	Text string
}

// Lexer tokenizes a reader
type Lexer struct {
	r *bufio.Reader
	p Pos
}

// NewLexer constructs a Lexer instance
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		r: bufio.NewReader(r),
		p: Pos{Line: 1},
	}
}

// Next returns the next token
func (l *Lexer) Next() Token {
	return Token{
		Pos:  l.p,
		Type: EOF,
		Text: "",
	}
}

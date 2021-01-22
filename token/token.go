package token

import (
	"fmt"
	"strings"
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
	NUMBER
	COMMA
	ASSIGN
	BOOL
)

// String returns a string representation of the type
func (t Type) String() string {
	switch t {
	case INVALID:
		return "INVALID"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case COMMA:
		return "COMMA"
	case ASSIGN:
		return "ASSIGN"
	case BOOL:
		return "BOOL"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}

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
	data  []rune
	index int
	pos   Pos
}

// NewLexer constructs a Lexer instance
func NewLexer(input string) *Lexer {
	return &Lexer{
		// TODO: make sure this is correct
		data: []rune(input),
		pos:  Pos{Line: 1},
	}
}

// eof is a sentinel value used in place of a rune when we're
// at the end of the input
const eof = 0x00

// read a rune and advance to the next one
func (l *Lexer) read() rune {
	if l.index >= len(l.data) {
		return eof
	}
	l.index++
	l.pos.Column++
	return l.data[l.index-1]
}

// peek reveals the next rune without advancing
func (l *Lexer) peek() rune {
	if l.index >= len(l.data) {
		return eof
	}
	return l.data[l.index]
}

// Next returns the next token
func (l *Lexer) Next() Token {
	r := l.peek()
	pos := l.pos
	switch {
	case r == eof:
		return Token{
			Pos:  pos,
			Type: EOF,
			Text: "",
		}
	case isDigit(r) || r == '-':
		return Token{
			Type: NUMBER,
			Pos:  pos,
			Text: l.number(),
		}
	default:
		return Token{
			Pos:  pos,
			Type: INVALID,
			Text: string([]rune{l.read()}),
		}
	}
}

func (l *Lexer) number() string {
	var text strings.Builder
	if l.peek() == '-' || l.peek() == '+' {
		text.WriteRune(l.read())
	}
	for isDigit(l.peek()) {
		text.WriteRune(l.read())
	}
	return text.String()
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

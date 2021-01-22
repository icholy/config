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

// Next returns the next token
func (l *Lexer) Next() Token {
	l.whitespace()
	ch := l.peek()
	pos := l.pos
	switch {
	case ch == eof:
		return Token{
			Pos:  pos,
			Type: EOF,
			Text: "",
		}
	case isDigit(ch) || ch == '-':
		return Token{
			Pos:  pos,
			Type: NUMBER,
			Text: l.number(),
		}
	case ch == '"':
		return Token{
			Pos:  pos,
			Type: STRING,
			Text: l.str(),
		}
	case isAlpha(ch):
		return Token{
			Pos:  pos,
			Type: IDENT,
			Text: l.ident(),
		}
	case ch == '=':
		return l.chartok(ASSIGN)
	case ch == '{':
		return l.chartok(LBRACE)
	case ch == '}':
		return l.chartok(RBRACE)
	case ch == '[':
		return l.chartok(LBRACKET)
	case ch == ']':
		return l.chartok(RBRACKET)
	case ch == ',':
		return l.chartok(COMMA)
	default:
		return l.chartok(INVALID)
	}
}

// eof is a sentinel value used in place of a rune when we're
// at the end of the input
const eof = 0x00

// eof returns true when we're at the end of file
func (l *Lexer) eof() bool {
	return l.index >= len(l.data)
}

// read a rune and advance to the next one
func (l *Lexer) read() rune {
	if l.eof() {
		return eof
	}
	l.index++
	l.pos.Column++
	return l.data[l.index-1]
}

// peek reveals the next rune without advancing
func (l *Lexer) peek() rune {
	if l.eof() {
		return eof
	}
	return l.data[l.index]
}

// whitespace skips all whitespace
func (l *Lexer) whitespace() {
	for isWhite(l.peek()) {
		l.read()
	}
}

// chartok is a helper which returns a single character token
func (l *Lexer) chartok(typ Type) Token {
	pos := l.pos
	return Token{
		Pos:  pos,
		Type: typ,
		Text: string([]rune{l.read()}),
	}
}

// number reads a number literal
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

// str reads a string literal
func (l *Lexer) str() string {
	l.read()
	var escaped bool
	var text strings.Builder
	for !l.eof() {
		ch := l.peek()
		if escaped {
			switch ch {
			case 't':
				text.WriteByte('\t')
			case 'r':
				text.WriteByte('\r')
			case 'n':
				text.WriteByte('\n')
			default:
				text.WriteRune(ch)
			}
			escaped = false
		} else {
			if ch == '"' {
				break
			}
			if ch == '\\' {
				escaped = true
			} else {
				text.WriteRune(ch)
			}
		}
		l.read()
	}
	l.read()
	return text.String()
}

// ident reads an identifier
func (l *Lexer) ident() string {
	var text strings.Builder
	ch := l.peek()
	for isAlpha(ch) || isDigit(ch) || ch == '_' {
		text.WriteRune(l.read())
		ch = l.peek()
	}
	return text.String()
}

// isDigit returns true if ch is a digit
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isAlpha returns true if ch is a letter
func isAlpha(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

// isWhite returns true if ch is whitespace
func isWhite(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

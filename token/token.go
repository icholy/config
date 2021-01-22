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
	COMMENT
	NEWLINE
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
	case COMMENT:
		return "COMMENT"
	case NEWLINE:
		return "NEWLINE"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}

// Pos is the position inside the file
type Pos struct {
	Line, Column, Offset int
}

// Token contains a token's type and text
type Token struct {
	Start Pos
	Type  Type
	Text  string
}

// Lexer tokenizes a reader
type Lexer struct {
	data    []rune
	index   int
	current Pos
}

// NewLexer constructs a Lexer instance
func NewLexer(input string) *Lexer {
	return &Lexer{
		// TODO: make sure this is correct
		data:    []rune(input),
		current: Pos{Line: 1, Column: 1},
	}
}

// Next returns the next token
func (l *Lexer) Next() Token {
	if pos, ok := l.whitespace(); ok {
		return Token{
			Start: pos,
			Type:  NEWLINE,
		}
	}
	ch := l.peek()
	pos := l.current
	switch {
	case ch == eof:
		return Token{
			Start: pos,
			Type:  EOF,
		}
	case isDigit(ch) || ch == '-':
		return Token{
			Start: pos,
			Type:  NUMBER,
			Text:  l.number(),
		}
	case ch == '"':
		text, ok := l.str()
		if !ok {
			return l.invalid(pos, text)
		}
		return Token{
			Start: pos,
			Type:  STRING,
			Text:  text,
		}
	case isAlpha(ch):
		text := l.ident()
		return Token{
			Start: pos,
			Type:  IDENT,
			Text:  text,
		}
	case ch == '/':
		text, ok := l.comment()
		if !ok {
			return l.invalid(pos, text)
		}
		return Token{
			Start: pos,
			Type:  COMMENT,
			Text:  text,
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

// newline returns true if the next character is a newline
func (l *Lexer) newline() bool {
	return isNewline(l.peek())
}

// read a rune and advance to the next one
func (l *Lexer) read() rune {
	if l.eof() {
		return eof
	}
	ch := l.data[l.index]
	if isNewline(ch) {
		// handle CRLF
		if !(l.index > 0 && ch == '\n' && l.data[l.index-1] == '\r') {
			l.current.Line++
			l.current.Column = 1
		}
	} else {
		l.current.Column++
	}
	l.current.Offset++
	l.index++
	return ch
}

// peek reveals the next rune without advancing
func (l *Lexer) peek() rune {
	if l.eof() {
		return eof
	}
	return l.data[l.index]
}

// expect checks if the next rune is equal to ch.
// if it matches, true is returned and the tokenizer advnaces to the next rune.
func (l *Lexer) expect(ch rune) bool {
	if l.peek() != ch {
		return false
	}
	l.read()
	return true
}

// whitespace skips all whitespace and returns the position of the first newline.
// If the whitespace does not contain a newline, the second return value is false.
func (l *Lexer) whitespace() (Pos, bool) {
	var newline bool
	var pos Pos
	for isWhite(l.peek()) {
		if !newline && l.newline() {
			newline = true
			pos = l.current
		}
		l.read()
	}
	return pos, newline
}

// chartok is a helper which returns a single character token
func (l *Lexer) chartok(typ Type) Token {
	pos := l.current
	return Token{
		Start: pos,
		Type:  typ,
		Text:  string([]rune{l.read()}),
	}
}

// invalid is a helper which returns an invalid token
func (l *Lexer) invalid(pos Pos, text string) Token {
	return Token{
		Start: pos,
		Type:  INVALID,
		Text:  text,
	}
}

// number reads a number literal
func (l *Lexer) number() string {
	var text strings.Builder
	if l.peek() == '-' || l.peek() == '+' {
		text.WriteRune(l.read())
	}
	for isDigit(l.peek()) || l.peek() == '.' {
		text.WriteRune(l.read())
	}
	return text.String()
}

// str reads a string literal
func (l *Lexer) str() (string, bool) {
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
				l.read()
				return text.String(), true
			}
			if ch == '\\' {
				escaped = true
			} else {
				text.WriteRune(ch)
			}
		}
		l.read()
	}
	return text.String(), false
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

// comment reads a comment. The second bool parameter indicates if it's valid.
func (l *Lexer) comment() (string, bool) {
	var text strings.Builder
	if !l.expect('/') {
		return text.String(), false
	}
	text.WriteRune('/')
	if !l.expect('/') {
		return text.String(), false
	}
	text.WriteRune('/')
	for !l.eof() && !l.newline() {
		text.WriteRune(l.read())
	}
	return text.String(), true
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
	return ch == ' ' || ch == '\t' || isNewline(ch)
}

// isNewline returns true if ch is a newline
func isNewline(ch rune) bool {
	return ch == '\n' || ch == '\r'
}

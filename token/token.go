package token

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

const eof = 0x00

func (l *Lexer) read() rune {
	if l.index >= len(l.data) {
		return eof
	}
	l.index++
	return l.data[l.index-1]
}

func (l *Lexer) peek() rune {
	if l.index >= len(l.data) {
		return eof
	}
	return l.data[l.index]
}

// Next returns the next token
func (l *Lexer) Next() Token {

	switch l.peek() {
	case eof:
		return Token{
			Pos:  l.pos,
			Type: EOF,
			Text: "",
		}
	default:
		return Token{
			Pos:  l.pos,
			Type: INVALID,
			Text: string([]rune{l.read()}),
		}
	}
}

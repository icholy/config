package ast

import "github.com/icholy/config/token"

// Parser for the configuration language
type Parser struct {
	lex        *token.Lexer
	curr, peek token.Token
}

// NewParser constructs a new parser
func NewParser(lex *token.Lexer) *Parser {
	curr := lex.Next()
	peek := lex.Next()
	return &Parser{
		lex:  lex,
		curr: curr,
		peek: peek,
	}
}

func (p *Parser) next() error {
	p.curr = p.peek
	p.peek = p.lex.Next()
	return nil
}

func (p *Parser) parse() (*Block, error) {
	b := &Block{
		Start:   token.Pos{Line: 1, Column: 1, Offset: 0},
		Entries: []*Entry{},
	}
	e, err := p.entry()
	if err != nil {
		return nil, err
	}
	b.Entries = append(b.Entries, e)
	return b, nil
}

func (p *Parser) entry() (*Entry, error) {
	return &Entry{}, nil
}

// Parse the input
func Parse(input string) (*Block, error) {
	l := token.NewLexer(input)
	p := NewParser(l)
	return p.parse()
}

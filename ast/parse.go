package ast

import (
	"fmt"
	"strconv"

	"github.com/icholy/config/token"
)

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

func (p *Parser) expect(t token.Type) error {
	if p.curr.Type != t {
		return fmt.Errorf("unexpected token: %v", p.curr)
	}
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

func (p *Parser) number() (*Number, error) {
	num := &Number{
		Start: p.curr.Start,
	}
	if err := p.expect(token.NUMBER); err != nil {
		return nil, err
	}
	var err error
	num.Value, err = strconv.ParseFloat(p.curr.Text, 64)
	if err != nil {
		return nil, err
	}
	return num, p.next()
}

func (p *Parser) ident() (*Ident, error) {
	id := &Ident{
		Start: p.curr.Start,
	}
	if err := p.expect(token.IDENT); err != nil {
		return nil, err
	}
	id.Value = p.curr.Text
	return id, p.next()
}

func (p *Parser) entry() (*Entry, error) {
	e := &Entry{
		Start: p.curr.Start,
	}
	var err error
	// read name
	e.Name, err = p.ident()
	if err != nil {
		return nil, err
	}
	// read =
	if err := p.expect(token.ASSIGN); err != nil {
		return nil, err
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	switch p.curr.Type {
	case token.NUMBER:
		e.Value, err = p.number()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.curr)
	}
	// read value
	return e, nil
}

// Parse the input
func Parse(input string) (*Block, error) {
	l := token.NewLexer(input)
	p := NewParser(l)
	return p.parse()
}

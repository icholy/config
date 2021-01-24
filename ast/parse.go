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

// assert panics if the current token type doesn't match t.
// this helper should only be used in places where it should not ever panic.
func (p *Parser) assert(t token.Type) {
	if err := p.expect(t); err != nil {
		panic(err)
	}
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
	p.assert(token.NUMBER)
	v, err := strconv.ParseFloat(p.curr.Text, 64)
	if err != nil {
		return nil, err
	}
	n := &Number{
		Start: p.curr.Start,
		Value: v,
	}
	return n, p.next()
}

func (p *Parser) string() (*String, error) {
	p.assert(token.STRING)
	s := &String{
		Start: p.curr.Start,
		Value: p.curr.Text,
	}
	return s, p.next()
}

func (p *Parser) bool() (*Bool, error) {
	p.assert(token.IDENT)
	b := &Bool{
		Start: p.curr.Start,
	}
	switch p.curr.Text {
	case "false":
	case "true":
		b.Value = true
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.curr)
	}
	return b, nil
}

func (p *Parser) array() (*Array, error) {
	p.assert(token.LBRACKET)
	a := &Array{
		Start: p.curr.Start,
	}
	// read left bracket
	if err := p.next(); err != nil {
		return nil, err
	}
	for {
		// found right bracket, we're done
		if p.curr.Type == token.RBRACKET {
			break
		}
		// read a value
		v, err := p.value()
		if err != nil {
			return nil, err
		}
		a.Values = append(a.Values, v)
		// if there's no comma, we're done
		if p.curr.Type != token.COMMA {
			break
		}
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	if err := p.expect(token.RBRACKET); err != nil {
		return nil, err
	}
	return a, p.next()
}

func (p *Parser) ident() (*Ident, error) {
	p.assert(token.IDENT)
	id := &Ident{
		Start: p.curr.Start,
	}
	id.Value = p.curr.Text
	return id, p.next()
}

func (p *Parser) value() (Value, error) {
	switch p.curr.Type {
	case token.NUMBER:
		return p.number()
	case token.STRING:
		return p.string()
	case token.IDENT:
		return p.bool()
	case token.LBRACKET:
		return p.array()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.curr)
	}
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
	// read assign
	if err := p.expect(token.ASSIGN); err != nil {
		return nil, err
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	// read value
	e.Value, err = p.value()
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Parse the input
func Parse(input string) (*Block, error) {
	l := token.NewLexer(input)
	p := NewParser(l)
	return p.parse()
}

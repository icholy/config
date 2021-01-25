package ast

import (
	"fmt"
	"strconv"

	"github.com/icholy/config/token"
)

// Parser for the configuration language
type Parser struct {
	lex *token.Lexer
	tok token.Token
}

// NewParser constructs a new parser
func NewParser(lex *token.Lexer) *Parser {
	return &Parser{
		lex: lex,
		tok: lex.Next(),
	}
}

func (p *Parser) next() error {
	p.tok = p.lex.Next()
	return nil
}

func (p *Parser) expect(t token.Type) error {
	if p.tok.Type != t {
		return fmt.Errorf("unexpected token: %v", p.tok)
	}
	return nil
}

// newline skips newline tokens
func (p *Parser) newlines() error {
	for p.tok.Type == token.NEWLINE {
		if err := p.next(); err != nil {
			return err
		}
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
	ee, err := p.entries()
	if err != nil {
		return nil, err
	}
	b.Entries = ee
	if err := p.expect(token.EOF); err != nil {
		return nil, err
	}
	return b, nil
}

func (p *Parser) number() (*Number, error) {
	p.assert(token.NUMBER)
	v, err := strconv.ParseFloat(p.tok.Text, 64)
	if err != nil {
		return nil, err
	}
	n := &Number{
		Start: p.tok.Start,
		Value: v,
	}
	return n, p.next()
}

func (p *Parser) string() (*String, error) {
	p.assert(token.STRING)
	s := &String{
		Start: p.tok.Start,
		Value: p.tok.Text,
	}
	return s, p.next()
}

func (p *Parser) bool() (*Bool, error) {
	p.assert(token.IDENT)
	b := &Bool{
		Start: p.tok.Start,
	}
	switch p.tok.Text {
	case "false":
	case "true":
		b.Value = true
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.tok)
	}
	return b, p.next()
}

func (p *Parser) array() (*Array, error) {
	p.assert(token.LBRACKET)
	a := &Array{
		Start: p.tok.Start,
	}
	// read left bracket
	if err := p.next(); err != nil {
		return nil, err
	}
	for {
		// skip newlines
		if err := p.newlines(); err != nil {
			return nil, err
		}
		// found right bracket, we're done
		if p.tok.Type == token.RBRACKET {
			break
		}
		// read a value
		v, err := p.value()
		if err != nil {
			return nil, err
		}
		a.Values = append(a.Values, v)
		// skip newlines
		if err := p.newlines(); err != nil {
			return nil, err
		}
		// if there's no comma, we're done
		if p.tok.Type != token.COMMA {
			break
		}
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	// skip newlines
	if err := p.newlines(); err != nil {
		return nil, err
	}
	if err := p.expect(token.RBRACKET); err != nil {
		return nil, err
	}
	return a, p.next()
}

func (p *Parser) block() (*Block, error) {
	p.assert(token.LBRACE)
	b := &Block{
		Start: p.tok.Start,
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	ee, err := p.entries()
	if err != nil {
		return nil, err
	}
	b.Entries = ee
	if err := p.expect(token.RBRACE); err != nil {
		return nil, err
	}
	return b, p.next()
}

func (p *Parser) ident() (*Ident, error) {
	p.assert(token.IDENT)
	id := &Ident{
		Start: p.tok.Start,
		Value: p.tok.Text,
	}
	return id, p.next()
}

func (p *Parser) value() (Value, error) {
	switch p.tok.Type {
	case token.NUMBER:
		return p.number()
	case token.STRING:
		return p.string()
	case token.IDENT:
		return p.bool()
	case token.LBRACKET:
		return p.array()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.tok)
	}
}

func (p *Parser) entries() ([]*Entry, error) {
	var ee []*Entry
	for p.tok.Type == token.IDENT {
		e, err := p.entry()
		if err != nil {
			return nil, err
		}
		ee = append(ee, e)
	}
	return ee, nil
}

func (p *Parser) entry() (*Entry, error) {
	e := &Entry{
		Start: p.tok.Start,
	}
	var err error
	// read name
	e.Name, err = p.ident()
	if err != nil {
		return nil, err
	}
	switch p.tok.Type {
	case token.ASSIGN:
		// skip assign operator
		if err := p.next(); err != nil {
			return nil, err
		}
		// read value
		e.Value, err = p.value()
		if err != nil {
			return nil, err
		}
	case token.LBRACE:
		// read block
		e.Value, err = p.block()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.tok)
	}
	return e, nil
}

// Parse the input
func Parse(input string) (*Block, error) {
	l := token.NewLexer(input)
	p := NewParser(l)
	return p.parse()
}

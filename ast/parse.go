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

// next reads the next token from the lexer
func (p *Parser) next() {
	p.tok = p.lex.Next()
}

// expect returns an error if the current token's type doesn't match t
func (p *Parser) expect(t token.Type) error {
	if p.tok.Type != t {
		return fmt.Errorf("unexpected token: %v", p.tok)
	}
	return nil
}

// newline skips newline tokens
func (p *Parser) newlines() {
	for p.tok.Type == token.NEWLINE {
		p.next()
	}
}

// assert panics if the current token type doesn't match t.
// this helper should only be used in places where it should not ever panic.
func (p *Parser) assert(t token.Type) {
	if err := p.expect(t); err != nil {
		panic(err)
	}
}

// parse is the entry point. It parses implicit top-level block.
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

// number parses a Number
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
	p.next()
	return n, nil
}

// string parses a String
func (p *Parser) string() (*String, error) {
	p.assert(token.STRING)
	s := &String{
		Start: p.tok.Start,
		Value: p.tok.Text,
	}
	p.next()
	return s, nil
}

// bool parses a Bool
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
	p.next()
	return b, nil
}

// list parses a List
func (p *Parser) list() (*List, error) {
	p.assert(token.LBRACKET)
	l := &List{
		Start: p.tok.Start,
	}
	// read left bracket
	p.next()
	for {
		// skip newlines
		p.newlines()
		// found right bracket, we're done
		if p.tok.Type == token.RBRACKET {
			break
		}
		// read a value
		v, err := p.value()
		if err != nil {
			return nil, err
		}
		l.Values = append(l.Values, v)
		// skip newlines
		p.newlines()
		// if there's no comma, we're done
		if p.tok.Type != token.COMMA {
			break
		}
		p.next()
	}
	// skip newlines
	p.newlines()
	if err := p.expect(token.RBRACKET); err != nil {
		return nil, err
	}
	p.next()
	return l, nil
}

// block parses a Block
func (p *Parser) block() (*Block, error) {
	p.assert(token.LBRACE)
	b := &Block{
		Start: p.tok.Start,
	}
	p.next()
	ee, err := p.entries()
	if err != nil {
		return nil, err
	}
	b.Entries = ee
	p.newlines()
	if err := p.expect(token.RBRACE); err != nil {
		return nil, err
	}
	p.next()
	return b, nil
}

// ident parses an Ident
func (p *Parser) ident() (*Ident, error) {
	p.assert(token.IDENT)
	id := &Ident{
		Start: p.tok.Start,
		Value: p.tok.Text,
	}
	p.next()
	return id, nil
}

// value parses a value
func (p *Parser) value() (Value, error) {
	switch p.tok.Type {
	case token.NUMBER:
		return p.number()
	case token.STRING:
		return p.string()
	case token.IDENT:
		return p.bool()
	case token.LBRACKET:
		return p.list()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.tok)
	}
}

// entries parses a sequence of Entry nodes
func (p *Parser) entries() ([]*Entry, error) {
	var ee []*Entry
	for {
		p.newlines()
		if p.tok.Type != token.IDENT {
			break
		}
		e, err := p.entry()
		if err != nil {
			return nil, err
		}
		ee = append(ee, e)
	}
	return ee, nil
}

// entry parses an Entry
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
		p.next()
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

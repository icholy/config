package ast

import "github.com/icholy/config/token"

// Value ...
type Value interface{ value() }

// Block is a collection of entries
type Block struct {
	Start   token.Pos
	Entries []*Entry
}

func (Block) value() {}

// Ident ...
type Ident struct {
	Start token.Pos
	Value string
}

// Number ...
type Number struct {
	Start token.Pos
	Value float64
}

func (Number) value() {}

// Bool ...
type Bool struct {
	Start token.Pos
	Value bool
}

func (Bool) value() {}

// String ...
type String struct {
	Start token.Pos
	Value string
}

func (String) value() {}

type Array struct {
	Start  token.Pos
	Values []Value
}

func (Array) value() {}

// Entry is a key/value pair
type Entry struct {
	Start token.Pos
	Name  *Ident
	Value Value
}

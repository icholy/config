package ast

import (
	"encoding/json"

	"github.com/icholy/config/token"
)

// Value ...
type Value interface {
	value()
}

// Block is a collection of entries
type Block struct {
	Start   token.Pos
	Entries []*Entry
}

func (Block) value() {}

// MarshalJSON implements json.Marshaler
func (b *Block) MarshalJSON() ([]byte, error) {
	m := map[string][]Value{}
	for _, e := range b.Entries {
		m[e.Name.Value] = append(m[e.Name.Value], e.Value)
	}
	return json.Marshal(m)
}

// Ident ...
type Ident struct {
	Start token.Pos
	Value string
}

// MarshalJSON implements json.Marshaler
func (i *Ident) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Value)
}

// Number ...
type Number struct {
	Start token.Pos
	Value float64
}

func (Number) value() {}

// MarshalJSON implements json.Marshaler
func (n *Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

// Bool ...
type Bool struct {
	Start token.Pos
	Value bool
}

func (Bool) value() {}

// MarshalJSON implements json.Marshaler
func (b *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Value)
}

// String ...
type String struct {
	Start token.Pos
	Value string
}

func (String) value() {}

// MarshalJSON implements json.Marshaler
func (s *String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}

type List struct {
	Start  token.Pos
	Values []Value
}

func (List) value() {}

// MarshalJSON implements json.Marshaler
func (l *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.Values)
}

// Entry is a key/value pair
type Entry struct {
	Start token.Pos
	Name  *Ident
	Value Value
}

// MarshalJSON implements json.Marshaler
func (e *Entry) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name  *Ident
		Value Value
	}{
		Name:  e.Name,
		Value: e.Value,
	})
}

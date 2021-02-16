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
	m := map[string]interface{}{}
	for _, e := range b.Entries {
		name := e.Name.Value
		if b0, ok := e.Value.(*Block); ok {
			if blocks, ok := m[name].([]Value); ok {
				m[name] = append(blocks, b0)
			} else {
				m[name] = []Value{b0}
			}
		} else {
			m[name] = e.Value
		}
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

// List ...
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

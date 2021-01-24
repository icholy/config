package ast

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"

	"github.com/icholy/config/token"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect *Block
	}{
		{
			name:  "IntEntry",
			input: "foo=123",
			expect: &Block{
				Entries: []*Entry{
					{
						Name:  &Ident{Value: "foo"},
						Value: &Number{Value: 123},
					},
				},
			},
		},
		{
			name:  "StringEntry",
			input: `bar = "test"`,
			expect: &Block{
				Entries: []*Entry{
					{
						Name:  &Ident{Value: "bar"},
						Value: &String{Value: "test"},
					},
				},
			},
		},
		{
			name:  "BoolEntry",
			input: `baz = true`,
			expect: &Block{
				Entries: []*Entry{
					{
						Name:  &Ident{Value: "baz"},
						Value: &Bool{Value: true},
					},
				},
			},
		},
		{
			name:  "ArrayEntry",
			input: `poo = [1, false, "hello"]]`,
			expect: &Block{
				Entries: []*Entry{
					{
						Name: &Ident{Value: "poo"},
						Value: &Array{
							Values: []Value{
								&Number{Value: 1},
								&Bool{Value: false},
								&String{Value: "hello"},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := Parse(tt.input)
			assert.NilError(t, err)
			assert.DeepEqual(t, tt.expect, block, cmpopts.IgnoreTypes(token.Pos{}))
		})
	}
}

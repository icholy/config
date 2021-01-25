package ast

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"

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
			name:  "TrueEntry",
			input: "baz = true",
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
			name:  "FalseEntry",
			input: "baz = false",
			expect: &Block{
				Entries: []*Entry{
					{
						Name:  &Ident{Value: "baz"},
						Value: &Bool{Value: false},
					},
				},
			},
		},
		{
			name:  "ListEntry",
			input: `poo = [1, false, "hello"]`,
			expect: &Block{
				Entries: []*Entry{
					{
						Name: &Ident{Value: "poo"},
						Value: &List{
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
		{
			name:  "EmptyListEntry",
			input: "poo = []",
			expect: &Block{
				Entries: []*Entry{
					{
						Name:  &Ident{Value: "poo"},
						Value: &List{},
					},
				},
			},
		},
		{
			name:  "ListTrailingCommaEntry",
			input: "poo = [1,]",
			expect: &Block{
				Entries: []*Entry{
					{
						Name: &Ident{Value: "poo"},
						Value: &List{
							Values: []Value{
								&Number{Value: 1},
							},
						},
					},
				},
			},
		},
		{
			name:  "ListWithNewline",
			input: "poo = [\n1,\n]",
			expect: &Block{
				Entries: []*Entry{
					{
						Name: &Ident{Value: "poo"},
						Value: &List{
							Values: []Value{
								&Number{Value: 1},
							},
						},
					},
				},
			},
		},
		{
			name:  "EmptyBlock",
			input: "block {}",
			expect: &Block{
				Entries: []*Entry{
					{
						Name:  &Ident{Value: "block"},
						Value: &Block{},
					},
				},
			},
		},
		{
			name:  "Block",
			input: "block { foo = 1 }",
			expect: &Block{
				Entries: []*Entry{
					{
						Name: &Ident{Value: "block"},
						Value: &Block{
							Entries: []*Entry{
								{
									Name:  &Ident{Value: "foo"},
									Value: &Number{Value: 1},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "BlockWithNewlines",
			input: "block {\nfoo = 1\nbar = 2\n}",
			expect: &Block{
				Entries: []*Entry{
					{
						Name: &Ident{Value: "block"},
						Value: &Block{
							Entries: []*Entry{
								{
									Name:  &Ident{Value: "foo"},
									Value: &Number{Value: 1},
								},
								{
									Name:  &Ident{Value: "bar"},
									Value: &Number{Value: 2},
								},
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

func TestParseGolden(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"basic"},
		{"nested"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := golden.Get(t, filepath.Join(tt.name, "input.conf"))
			block, err := Parse(string(data))
			assert.NilError(t, err)
			actual, err := json.MarshalIndent(block, "", "  ")
			assert.NilError(t, err)
			golden.Assert(t, string(actual), filepath.Join(tt.name, "output.json"))
		})
	}
}

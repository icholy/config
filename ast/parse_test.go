package ast

import (
	"testing"

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
				Start: token.Pos{1, 1, 0},
				Entries: []*Entry{
					{
						Start: token.Pos{1, 1, 0},
						Name: Ident{
							Start: token.Pos{1, 1, 0},
							Value: "foo",
						},
						Value: &Number{
							Start: token.Pos{1, 5, 4},
							Value: 123,
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
			assert.DeepEqual(t, tt.expect, block)
		})
	}
}

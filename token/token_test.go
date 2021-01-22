package token

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []Token
	}{
		{
			name:  "eof",
			input: "",
			expect: []Token{
				{
					Pos:  Pos{1, 0},
					Type: EOF,
				},
			},
		},
		{
			name:  "int",
			input: "42",
			expect: []Token{
				{
					Pos:  Pos{1, 0},
					Type: NUMBER,
					Text: "42",
				},
				{
					Pos:  Pos{1, 2},
					Type: EOF,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual []Token
			lex := NewLexer(tt.input)
			for {
				tok := lex.Next()
				actual = append(actual, tok)
				if tok.Type == EOF || tok.Type == INVALID {
					break
				}
			}
			assert.DeepEqual(t, tt.expect, actual)
		})
	}
}

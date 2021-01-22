package token

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		skip   bool
		expect []Token
	}{
		{
			name:  "EOF",
			input: "",
			expect: []Token{
				{Pos{1, 0}, EOF, ""},
			},
		},
		{
			name:  "Int",
			input: "42",
			expect: []Token{
				{Pos{1, 0}, NUMBER, "42"},
				{Pos{1, 2}, EOF, ""},
			},
		},
		{
			name:  "NegativeInt",
			input: "-42",
			expect: []Token{
				{Pos{1, 0}, NUMBER, "-42"},
				{Pos{1, 3}, EOF, ""},
			},
		},
		{
			name:  "String",
			input: `"hello world"`,
			expect: []Token{
				{Pos{1, 0}, STRING, "hello world"},
				{Pos{1, 13}, EOF, ""},
			},
		},
		{
			name:  "Assign",
			input: "=",
			expect: []Token{
				{Pos{1, 0}, ASSIGN, "="},
				{Pos{1, 1}, EOF, ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}
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

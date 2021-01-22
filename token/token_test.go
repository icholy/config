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
				{Pos{1, 0}, EOF, ""},
			},
		},
		{
			name:  "int",
			input: "42",
			expect: []Token{
				{Pos{1, 0}, NUMBER, "42"},
				{Pos{1, 2}, EOF, ""},
			},
		},
		{
			name:  "int",
			input: "-42",
			expect: []Token{
				{Pos{1, 0}, NUMBER, "-42"},
				{Pos{1, 3}, EOF, ""},
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

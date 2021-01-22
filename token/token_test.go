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
				{Pos{1, 1}, EOF, ""},
			},
		},
		{
			name:  "Int",
			input: "42",
			expect: []Token{
				{Pos{1, 1}, NUMBER, "42"},
				{Pos{1, 3}, EOF, ""},
			},
		},
		{
			name:  "NegativeInt",
			input: "-42",
			expect: []Token{
				{Pos{1, 1}, NUMBER, "-42"},
				{Pos{1, 4}, EOF, ""},
			},
		},
		{
			name:  "String",
			input: `"hello world"`,
			expect: []Token{
				{Pos{1, 1}, STRING, "hello world"},
				{Pos{1, 14}, EOF, ""},
			},
		},
		{
			name:  "BadString",
			input: `"whoops`,
			expect: []Token{
				{Pos{1, 1}, INVALID, "whoops"},
			},
		},
		{
			name:  "Assign",
			input: "=",
			expect: []Token{
				{Pos{1, 1}, ASSIGN, "="},
				{Pos{1, 2}, EOF, ""},
			},
		},
		{
			name:  "Ident",
			input: "key",
			expect: []Token{
				{Pos{1, 1}, IDENT, "key"},
				{Pos{1, 4}, EOF, ""},
			},
		},
		{
			name:  "LineComment",
			input: "// this is a comment",
			expect: []Token{
				{Pos{1, 1}, COMMENT, "// this is a comment"},
				{Pos{1, 21}, EOF, ""},
			},
		},
		{
			name:  "Block",
			input: "block { }",
			expect: []Token{
				{Pos{1, 1}, IDENT, "block"},
				{Pos{1, 7}, LBRACE, "{"},
				{Pos{1, 9}, RBRACE, "}"},
				{Pos{1, 10}, EOF, ""},
			},
		},
		{
			name:  "Newline",
			input: "foo = true\nbar",
			expect: []Token{
				{Pos{1, 1}, IDENT, "foo"},
				{Pos{1, 5}, ASSIGN, "="},
				{Pos{1, 7}, IDENT, "true"},
				{Pos{1, 11}, NEWLINE, ""},
				{Pos{2, 1}, IDENT, "bar"},
				{Pos{2, 4}, EOF, ""},
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

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
			name:  "BadString",
			input: `"whoops`,
			expect: []Token{
				{Pos{1, 0}, INVALID, "whoops"},
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
		{
			name:  "Ident",
			input: "key",
			expect: []Token{
				{Pos{1, 0}, IDENT, "key"},
				{Pos{1, 3}, EOF, ""},
			},
		},
		{
			name:  "LineComment",
			input: "// this is a comment",
			expect: []Token{
				{Pos{1, 0}, COMMENT, "// this is a comment"},
				{Pos{1, 20}, EOF, ""},
			},
		},
		{
			name:  "Block",
			input: "block { }",
			expect: []Token{
				{Pos{1, 0}, IDENT, "block"},
				{Pos{1, 6}, LBRACE, "{"},
				{Pos{1, 8}, RBRACE, "}"},
				{Pos{1, 9}, EOF, ""},
			},
		},
		{
			name:  "Newline",
			input: "foo = true\nbar=123",
			expect: []Token{
				{Pos{1, 0}, IDENT, "foo"},
				{Pos{1, 4}, ASSIGN, "="},
				{Pos{1, 6}, IDENT, "true"},
				{Pos{1, 10}, NEWLINE, ""},
				{Pos{1, 11}, IDENT, "bar"},
				{Pos{1, 14}, ASSIGN, "="},
				{Pos{1, 15}, NUMBER, "123"},
				{Pos{1, 18}, EOF, ""},
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

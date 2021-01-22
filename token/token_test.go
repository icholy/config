package token

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []Token
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual []Token
			lex := NewLexer(strings.NewReader(tt.input))
			for {
				tok := lex.Next()
				actual = append(actual, tok)
				if tok.Type == EOF {
					break
				}
			}
			assert.DeepEqual(t, tt.expect, actual)
		})
	}
}

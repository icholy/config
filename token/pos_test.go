package token

import (
	"path/filepath"
	"testing"

	"gotest.tools/v3/golden"
)

func TestSnip(t *testing.T) {
	tests := []struct {
		start, end Pos
		output     string
	}{
		{
			start:  Pos{1, 1, 0},
			end:    Pos{1, 1, 0},
			output: "empty.output",
		},
		{
			start:  Pos{1, 1, 0},
			end:    Pos{11, 15, 119},
			output: "full.output",
		},
		{
			start:  Pos{3, 10, 19},
			end:    Pos{3, 19, 24},
			output: "partline.output",
		},
	}
	input := golden.Get(t, filepath.FromSlash("snippet/source.conf"))
	for _, tt := range tests {
		t.Run(tt.output, func(t *testing.T) {
			s := Snip(string(input), tt.start, tt.end)
			golden.Assert(t, s.String(), filepath.Join("snippet", tt.output))
		})
	}
}

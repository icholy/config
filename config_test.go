package config

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		input string
		dst   func() interface{}
		want  func() interface{}
	}{
		{
			input: "a=123\nb=42\nc=\"hello\"",
			dst: func() interface{} {
				m := map[string]interface{}{}
				return &m
			},
			want: func() interface{} {
				m := map[string]interface{}{
					"a": float64(123),
					"b": float64(42),
					"c": "hello",
				}
				return &m
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, want := tt.dst(), tt.want()
			err := Unmarshal([]byte(tt.input), got)
			assert.NilError(t, err)
			assert.DeepEqual(t, got, want)
		})
	}
}

package config

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestUnmarshal(t *testing.T) {

	type Foo struct {
		A float64
		B string
		C bool
	}
	type Bar struct {
		Foo *Foo
	}

	tests := []struct {
		name  string
		input string
		dst   func() interface{}
		want  func() interface{}
	}{
		{
			name:  "FlatBlockToMap",
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
		{
			name:  "FlatBlockToStruct",
			input: "A=123\nB=\"hello\"",
			dst: func() interface{} {
				return &Foo{}
			},
			want: func() interface{} {
				return &Foo{A: 123, B: "hello"}
			},
		},
		{
			name:  "NilMap",
			input: "a=true",
			dst: func() interface{} {
				var m map[string]interface{}
				return &m
			},
			want: func() interface{} {
				m := map[string]interface{}{
					"a": true,
				}
				return &m
			},
		},
		{
			name:  "NilStruct",
			input: "Foo { A=123\nC=false }",
			dst: func() interface{} {
				return &Bar{}
			},
			want: func() interface{} {
				return &Bar{
					Foo: &Foo{A: 123, C: false},
				}
			},
		},
		{
			name:  "ExistingKey",
			input: "exists { A=123 }",
			dst: func() interface{} {
				m := map[string]interface{}{
					"existing": &Foo{B: "test"},
				}
				return &m
			},
			want: func() interface{} {
				m := map[string]interface{}{
					"existing": &Foo{A: 123, B: "test"},
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

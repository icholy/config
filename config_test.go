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

	type Baz struct {
		Int64   int64
		Uint32  int32
		Int     int
		Float32 float32
	}

	type List struct {
		Items []interface{}
	}

	type MultiBlock struct {
		Foo []*Foo
	}

	tests := []struct {
		name  string
		input string
		dst   func() interface{}
		want  func() interface{}
		skip  bool
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
					"exists": &Foo{B: "test"},
				}
				return &m
			},
			want: func() interface{} {
				m := map[string]interface{}{
					"exists": &Foo{A: 123, B: "test"},
				}
				return &m
			},
		},
		{
			name:  "NestedMap",
			input: "foo { bar { baz = 123 } }",
			dst: func() interface{} {
				var v interface{}
				return &v
			},
			want: func() interface{} {
				var v interface{} = map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": map[string]interface{}{
							"baz": float64(123),
						},
					},
				}
				return &v
			},
			skip: true,
		},
		{
			name:  "ConvertableTypes",
			input: "Int64=3\nUint32=234\nInt=44\nFloat32=10.10",
			dst: func() interface{} {
				return &Baz{}
			},
			want: func() interface{} {
				return &Baz{
					Int64:   3,
					Uint32:  234,
					Int:     44,
					Float32: 10.10,
				}
			},
		},
		{
			name:  "List",
			input: "Items=[1,2,3,4]",
			dst: func() interface{} {
				return &List{}
			},
			want: func() interface{} {
				return &List{
					Items: []interface{}{
						float64(1),
						float64(2),
						float64(3),
						float64(4),
					},
				}
			},
		},
		{
			name:  "MultiBlockKey",
			input: "Foo { A = 123 } Foo { A = 321 }",
			dst: func() interface{} {
				return &MultiBlock{}
			},
			want: func() interface{} {
				return &MultiBlock{
					Foo: []*Foo{
						{A: 123},
						{A: 321},
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.SkipNow()
			}
			got, want := tt.dst(), tt.want()
			err := Unmarshal([]byte(tt.input), got)
			assert.NilError(t, err)
			assert.DeepEqual(t, got, want)
		})
	}
}

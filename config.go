package config

import (
	"fmt"
	"reflect"

	"github.com/icholy/config/ast"
)

// Unmarshal ...
func Unmarshal(data []byte, v interface{}) error {
	block, err := ast.Parse(string(data))
	if err != nil {
		return err
	}
	return decodeValue(block, reflect.ValueOf(v))
}

func decodeBlock(b *ast.Block, dst reflect.Value) error {
	switch dst.Kind() {
	case reflect.Map:
		return decodeBlockToMap(b, dst)
	case reflect.Struct:
		return decodeBlockToStruct(b, dst)
	default:
		return fmt.Errorf("cannot decode block to: %v", dst)
	}
}

func decodeBlockToMap(b *ast.Block, dst reflect.Value) error {
	if dst.IsNil() {
		dst.Set(reflect.MakeMap(dst.Type()))
	}
	for _, e := range b.Entries {
		elem := reflect.New(dst.Type().Elem()).Elem()
		if err := decodeValue(e.Value, elem); err != nil {
			return err
		}
		dst.SetMapIndex(reflect.ValueOf(e.Name.Value), elem)
	}
	return nil
}

func decodeBlockToStruct(b *ast.Block, dst reflect.Value) error {
	typ := dst.Type()
	for _, e := range b.Entries {
		field, ok := typ.FieldByName(e.Name.Value)
		if !ok {
			return fmt.Errorf("no matching field: %q", e.Name.Value)
		}
		if field.Anonymous {
			return fmt.Errorf("anonymous fields are not supported: %q", e.Name.Value)
		}
		if err := decodeValue(e.Value, dst.Field(field.Index[0])); err != nil {
			return err
		}
	}
	return nil
}

func decodeNumber(n *ast.Number, dst reflect.Value) error {
	dst.Set(reflect.ValueOf(n.Value))
	return nil
}

func decodeString(s *ast.String, dst reflect.Value) error {
	dst.Set(reflect.ValueOf(s.Value))
	return nil
}

func decodeBool(b *ast.Bool, dst reflect.Value) error {
	dst.Set(reflect.ValueOf(b.Value))
	return nil
}

func decodeValue(v ast.Value, dst reflect.Value) error {
	for dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = reflect.Indirect(dst)
	}
	switch v := v.(type) {
	case *ast.Block:
		return decodeBlock(v, dst)
	case *ast.Number:
		return decodeNumber(v, dst)
	case *ast.String:
		return decodeString(v, dst)
	case *ast.Bool:
		return decodeBool(v, dst)
	default:
		return fmt.Errorf("not implemented: %T", v)
	}
}

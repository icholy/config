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
	return decodeBlockToMap(b, dst)
}

func decodeBlockToMap(b *ast.Block, dst reflect.Value) error {
	for _, e := range b.Entries {
		elem := reflect.New(dst.Type().Elem()).Elem()
		if err := decodeValue(e.Value, elem); err != nil {
			return err
		}
		dst.SetMapIndex(reflect.ValueOf(e.Name.Value), elem)
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

func decodeValue(v ast.Value, dst reflect.Value) error {
	for dst.Kind() == reflect.Ptr {
		dst = reflect.Indirect(dst)
	}
	switch v := v.(type) {
	case *ast.Block:
		return decodeBlock(v, dst)
	case *ast.Number:
		return decodeNumber(v, dst)
	case *ast.String:
		return decodeString(v, dst)
	default:
		return fmt.Errorf("not implemented: %T", v)
	}
}

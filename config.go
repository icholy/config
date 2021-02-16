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
	case reflect.Interface:
		if dst.IsNil() {
			m := map[string]interface{}{}
			dst.Set(reflect.ValueOf(m))
		}
		return decodeValue(b, dst.Elem())
	case reflect.Map:
		return decodeBlockToMap(b, dst)
	case reflect.Struct:
		return decodeBlockToStruct(b, dst)
	case reflect.Slice:
		return decodeBlockToSlice(b, dst)
	default:
		return fmt.Errorf("cannot decode block to: %v", dst.Type())
	}
}

func decodeBlockToSlice(b *ast.Block, dst reflect.Value) error {
	elem := reflect.New(dst.Type().Elem()).Elem()
	if err := decodeValue(b, elem); err != nil {
		return err
	}
	dst.Set(reflect.Append(dst, elem))
	return nil
}

func decodeBlockToMap(b *ast.Block, dst reflect.Value) error {
	if dst.IsNil() {
		dst.Set(reflect.MakeMap(dst.Type()))
	}
	for _, e := range b.Entries {
		key := reflect.ValueOf(e.Name.Value)
		elem := dst.MapIndex(key)
		if !elem.IsValid() {
			elem = reflect.New(dst.Type().Elem()).Elem()
		}
		if err := decodeValue(e.Value, elem); err != nil {
			return err
		}
		dst.SetMapIndex(key, elem)
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

func decodeList(l *ast.List, dst reflect.Value) error {
	switch dst.Kind() {
	case reflect.Interface:
		if dst.IsNil() {
			s := []interface{}{}
			dst.Set(reflect.ValueOf(s))
		}
		return decodeValue(l, dst.Elem())
	case reflect.Slice:
		for _, v := range l.Values {
			elem := reflect.New(dst.Type().Elem()).Elem()
			if err := decodeValue(v, elem); err != nil {
				return err
			}
			dst.Set(reflect.Append(dst, elem))
		}
		return nil
	default:
		return fmt.Errorf("cannot decode block to: %v", dst.Type())
	}
}

func decodePrimitive(primitive interface{}, dst reflect.Value) error {
	v := reflect.ValueOf(primitive)
	if v.Type().ConvertibleTo(dst.Type()) {
		v = v.Convert(dst.Type())
	}
	if !v.Type().AssignableTo(dst.Type()) {
		return fmt.Errorf("cannot assign %v to %v", v.Type(), dst.Type())
	}
	dst.Set(v)
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
	case *ast.List:
		return decodeList(v, dst)
	case *ast.Number:
		return decodePrimitive(v.Value, dst)
	case *ast.String:
		return decodePrimitive(v.Value, dst)
	case *ast.Bool:
		return decodePrimitive(v.Value, dst)
	default:
		return fmt.Errorf("not implemented: %T", v)
	}
}

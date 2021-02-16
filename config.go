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

func byName(ee []*ast.Entry) map[string][]*ast.Entry {
	groups := map[string][]*ast.Entry{}
	for _, e := range ee {
		groups[e.Name.Value] = append(groups[e.Name.Value], e)
	}
	return groups
}

func decodeEntry(e *ast.Entry, dst reflect.Value, typ reflect.Type, multi bool) (reflect.Value, error) {
	if !dst.IsValid() {
		dst = reflect.New(typ).Elem()
	}
	if err := decodeValue(e.Value, dst); err != nil {
		return reflect.Value{}, err
	}
	return dst, nil
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
	for name, entries := range byName(b.Entries) {
		for _, e := range entries {
			key := reflect.ValueOf(name)
			val, err := decodeEntry(e, dst.MapIndex(key), dst.Type().Elem(), len(entries) > 1)
			if err != nil {
				return err
			}
			dst.SetMapIndex(key, val)
		}
	}
	return nil
}

func decodeBlockToStruct(b *ast.Block, dst reflect.Value) error {
	for name, entries := range byName(b.Entries) {
		for _, e := range entries {
			field, ok := dst.Type().FieldByName(name)
			if !ok {
				return fmt.Errorf("no matching field: %q", name)
			}
			if field.Anonymous {
				return fmt.Errorf("anonymous fields are not supported: %q", name)
			}
			idx := field.Index[0]
			val, err := decodeEntry(e, dst.Field(idx), field.Type, len(entries) > 1)
			if err != nil {
				return err
			}
			dst.Field(idx).Set(val)
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

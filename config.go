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
	return decodeValue(block, reflect.ValueOf(v), false)
}

func byName(ee []*ast.Entry) map[string][]*ast.Entry {
	groups := map[string][]*ast.Entry{}
	for _, e := range ee {
		groups[e.Name.Value] = append(groups[e.Name.Value], e)
	}
	return groups
}

func decodeBlock(b *ast.Block, dst reflect.Value, multi bool) error {
	dst = realise(dst)
	switch dst.Kind() {
	case reflect.Interface:
		if dst.IsNil() {
			s := []map[string]interface{}{}
			dst.Set(reflect.ValueOf(s))
		}
		return decodeBlock(b, dst.Elem(), multi)
	case reflect.Map:
		if dst.IsNil() {
			dst.Set(reflect.MakeMap(dst.Type()))
		}
		for name, entries := range byName(b.Entries) {
			for _, e := range entries {
				key := reflect.ValueOf(name)
				val := dst.MapIndex(key)
				ptr := reflect.New(dst.Type().Elem())
				if val.IsValid() {
					ptr.Elem().Set(val)
				}
				if err := decodeValue(e.Value, ptr.Elem(), len(entries) > 1); err != nil {
					return err
				}
				dst.SetMapIndex(key, ptr.Elem())
			}
		}
		return nil
	case reflect.Struct:
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
				if err := decodeValue(e.Value, dst.Field(idx), len(entries) > 1); err != nil {
					return err
				}
			}
		}
		return nil
	case reflect.Slice:
		elem := reflect.New(dst.Type().Elem()).Elem()
		if err := decodeValue(b, elem, multi); err != nil {
			return err
		}
		dst.Set(reflect.Append(dst, elem))
		return nil
	default:
		return fmt.Errorf("cannot decode block to: %v", dst.Type())
	}
}

func decodeList(l *ast.List, dst reflect.Value, multi bool) error {
	dst = realise(dst)
	switch dst.Kind() {
	case reflect.Interface:
		if dst.IsNil() {
			s := []interface{}{}
			dst.Set(reflect.ValueOf(s))
		}
		return decodeValue(l, dst.Elem(), multi)
	case reflect.Slice:
		for _, v := range l.Values {
			elem := reflect.New(dst.Type().Elem()).Elem()
			if err := decodeValue(v, elem, multi); err != nil {
				return err
			}
			dst.Set(reflect.Append(dst, elem))
		}
		return nil
	default:
		return fmt.Errorf("cannot decode block to: %v", dst.Type())
	}
}

func decodePrimitive(primitive interface{}, dst reflect.Value, multi bool) error {
	dst = realise(dst)
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

func decodeValue(v ast.Value, dst reflect.Value, multi bool) error {
	switch v := v.(type) {
	case *ast.Block:
		return decodeBlock(v, dst, multi)
	case *ast.List:
		return decodeList(v, dst, multi)
	case *ast.Number:
		return decodePrimitive(v.Value, dst, multi)
	case *ast.String:
		return decodePrimitive(v.Value, dst, multi)
	case *ast.Bool:
		return decodePrimitive(v.Value, dst, multi)
	default:
		return fmt.Errorf("not implemented: %T", v)
	}
}

func realise(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = reflect.Indirect(v)
	}
	return v
}

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
	return decodeBlock(block, reflect.ValueOf(v), false)
}

func byName(ee []*ast.Entry) map[string][]*ast.Entry {
	groups := map[string][]*ast.Entry{}
	for _, e := range ee {
		groups[e.Name.Value] = append(groups[e.Name.Value], e)
	}
	return groups
}

func decodeBlock(b *ast.Block, dst reflect.Value, multi bool) error {
	dst, update := realise(dst, func() reflect.Value {
		if multi {
			return reflect.ValueOf([]map[string]interface{}{})
		}
		return reflect.ValueOf(map[string]interface{}{})
	})
	switch dst.Kind() {
	case reflect.Map:
		if dst.IsNil() {
			dst.Set(reflect.MakeMap(dst.Type()))
		}
		for name, entries := range byName(b.Entries) {
			key := reflect.ValueOf(name)
			val := dst.MapIndex(key)
			// make an addressable copy of val
			var tmp reflect.Value
			if val.IsValid() {
				tmp = reflect.New(val.Type()).Elem()
				tmp.Set(val)
			} else {
				tmp = reflect.New(dst.Type().Elem()).Elem()
			}
			for _, e := range entries {
				if err := decodeValue(e.Value, tmp, len(entries) > 1); err != nil {
					return err
				}
			}
			dst.SetMapIndex(key, tmp)
		}
		return nil
	case reflect.Struct:
		for name, entries := range byName(b.Entries) {
			field, ok := dst.Type().FieldByName(name)
			if !ok {
				return fmt.Errorf("no matching field: %q", name)
			}
			if field.Anonymous {
				return fmt.Errorf("anonymous fields are not supported: %q", name)
			}
			idx := field.Index[0]
			for _, e := range entries {
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
		update(reflect.Append(dst, elem))
		return nil
	default:
		return fmt.Errorf("cannot decode block to: %v", dst.Type())
	}
}

func decodeList(l *ast.List, dst reflect.Value, multi bool) error {
	dst, update := realise(dst, func() reflect.Value {
		s := []interface{}{}
		return reflect.ValueOf(s)
	})
	switch dst.Kind() {
	case reflect.Slice:
		for _, v := range l.Values {
			elem := reflect.New(dst.Type().Elem()).Elem()
			if err := decodeValue(v, elem, multi); err != nil {
				return err
			}
			update(reflect.Append(dst, elem))
		}
		return nil
	default:
		return fmt.Errorf("cannot decode block to: %v", dst.Type())
	}
}

func decodePrimitive(primitive interface{}, dst reflect.Value, multi bool) error {
	dst, update := realise(dst, nil)
	v := reflect.ValueOf(primitive)
	if v.Type().ConvertibleTo(dst.Type()) {
		v = v.Convert(dst.Type())
	}
	if !v.Type().AssignableTo(dst.Type()) {
		return fmt.Errorf("cannot assign %v to %v", v.Type(), dst.Type())
	}
	update(v)
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

func realise(v reflect.Value, zero func() reflect.Value) (reflect.Value, func(reflect.Value)) {
	var settable reflect.Value
LOOP:
	for {
		switch v.Kind() {
		case reflect.Ptr:
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = reflect.Indirect(v)
		case reflect.Interface:
			settable = v
			if v.IsNil() {
				if zero != nil {
					z := zero()
					v.Set(z)
					v = z
				}
				break LOOP
			}
			v = v.Elem()
		default:
			break LOOP
		}
	}
	if !settable.IsValid() {
		settable = v
	}
	return v, func(v reflect.Value) { settable.Set(v) }
}

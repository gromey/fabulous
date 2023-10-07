package hydrogen

import (
	"errors"
	"reflect"
)

func valueFromPtr(rv reflect.Value) (reflect.Value, error) {
	if rv.Kind() == reflect.Pointer {
		return valueFromPtr(rv.Elem())
	}
	if rv.Kind() != reflect.Struct {
		return rv, errors.New("the input value is not a struct")
	}
	return rv, nil
}

func ignore(name string, fields ...string) bool {
	for _, f := range fields {
		if f == name {
			return false
		}
	}
	return true
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return v.IsZero()
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}
	return !v.IsValid()
}

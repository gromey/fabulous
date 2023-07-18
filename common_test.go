package hydrogen

import (
	"reflect"
	"testing"
)

func equal(t *testing.T, exp, got interface{}) {
	if !reflect.DeepEqual(exp, got) {
		t.Fatalf("Not equal:\nexp: %v\ngot: %v", exp, got)
	}
}

func Test_isEmptyValue(t *testing.T) {
	var foo any
	var bar any
	bar = 77
	var tests = []struct {
		value  any
		expect bool
	}{
		{
			value:  true,
			expect: false,
		},
		{
			value:  false,
			expect: true,
		},
		{
			value:  1,
			expect: false,
		},
		{
			value:  0,
			expect: true,
		},
		{
			value:  1.1,
			expect: false,
		},
		{
			value:  0.0,
			expect: true,
		},
		{
			value:  "a",
			expect: false,
		},
		{
			value:  "",
			expect: true,
		},
		{
			value:  &struct{}{},
			expect: false,
		},
		{
			value:  (*struct{})(nil),
			expect: true,
		},
		{
			value:  []int{1},
			expect: false,
		},
		{
			value:  []int{},
			expect: true,
		},
		{
			value:  foo,
			expect: true,
		},
		{
			value:  bar,
			expect: false,
		},
	}
	for _, tt := range tests {
		b := isEmptyValue(reflect.ValueOf(tt.value))
		equal(t, tt.expect, b)
	}
}

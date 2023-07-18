package hydrogen

import (
	"errors"
	"testing"
)

var tag = New("tag")

type simple struct {
	A string
	B int
	C float64
	D int
	e string
}

type nested struct {
	A string
	B int
	C float64 `tag:"name"`
	simple
	S *simple
}

type skip struct {
	A string `tag:"-"`
	B int    `tag:"-"`
	C float64
	*simple
	S simple
}

type pointer struct {
	A *int
}

var (
	ns = 12345

	s1 = simple{A: "test", B: 0, C: 3.14, D: 28, e: "not exported"}
	s2 = nested{A: "foo", B: 0, C: 3.14, simple: s1, S: &s1}
	s3 = skip{A: "foo", B: 0, C: 3.14, simple: &s1, S: s1}
	s4 = pointer{}
)

func TestEngine_Fields(t *testing.T) {
	type ts struct {
		name      string
		input     any
		omitempty bool
		fields    []string
		expect    Fields
		err       error
	}

	tests := []ts{
		{
			name:  "Not a pointer",
			input: s1,
			err:   errors.New("the input value is not a pointer"),
		},
		{
			name:  "Not a struct",
			input: &ns,
			err:   errors.New("the input value is not a struct"),
		},
		{
			name:  "All fields",
			input: &s1,
			expect: Fields{
				{Name: "A", Value: "test", Addr: &s1.A},
				{Name: "B", Value: 0, Addr: &s1.B},
				{Name: "C", Value: 3.14, Addr: &s1.C},
				{Name: "D", Value: 28, Addr: &s1.D},
			},
		},
		{
			name:      "All fields with omitempty",
			input:     &s1,
			omitempty: true,
			expect: Fields{
				{Name: "A", Value: "test", Addr: &s1.A},
				{Name: "C", Value: 3.14, Addr: &s1.C},
				{Name: "D", Value: 28, Addr: &s1.D},
			},
		},
		{
			name:   "Selected fields",
			input:  &s1,
			fields: []string{"A", "C"},
			expect: Fields{
				{Name: "A", Value: "test", Addr: &s1.A},
				{Name: "C", Value: 3.14, Addr: &s1.C},
			},
		},
		{
			name:      "Selected fields with omitempty",
			input:     &s1,
			omitempty: true,
			fields:    []string{"A", "B"},
			expect: Fields{
				{Name: "A", Value: "test", Addr: &s1.A},
			},
		},
		{
			name:  "All fields with nested",
			input: &s2,
			expect: Fields{
				{Name: "A", Value: "foo", Addr: &s2.A},
				{Name: "B", Value: 0, Addr: &s2.B},
				{Name: "name", Value: 3.14, Addr: &s2.C},
				{Name: "A", Value: "test", Addr: &s2.simple.A},
				{Name: "B", Value: 0, Addr: &s2.simple.B},
				{Name: "C", Value: 3.14, Addr: &s2.simple.C},
				{Name: "D", Value: 28, Addr: &s2.simple.D},
				{Name: "A", Value: "test", Addr: &s2.S.A},
				{Name: "B", Value: 0, Addr: &s2.S.B},
				{Name: "C", Value: 3.14, Addr: &s2.S.C},
				{Name: "D", Value: 28, Addr: &s2.S.D},
			},
		},
		{
			name:      "All fields with nested with omitempty",
			input:     &s2,
			omitempty: true,
			expect: Fields{
				{Name: "A", Value: "foo", Addr: &s2.A},
				{Name: "name", Value: 3.14, Addr: &s2.C},
				{Name: "A", Value: "test", Addr: &s2.simple.A},
				{Name: "C", Value: 3.14, Addr: &s2.simple.C},
				{Name: "D", Value: 28, Addr: &s2.simple.D},
				{Name: "A", Value: "test", Addr: &s2.S.A},
				{Name: "C", Value: 3.14, Addr: &s2.S.C},
				{Name: "D", Value: 28, Addr: &s2.S.D},
			},
		},
		{
			name:   "Selected fields with nested",
			input:  &s2,
			fields: []string{"A", "B"},
			expect: Fields{
				{Name: "A", Value: "foo", Addr: &s2.A},
				{Name: "B", Value: 0, Addr: &s2.B},
				{Name: "A", Value: "test", Addr: &s2.simple.A},
				{Name: "B", Value: 0, Addr: &s2.simple.B},
				{Name: "A", Value: "test", Addr: &s2.S.A},
				{Name: "B", Value: 0, Addr: &s2.S.B},
			},
		},
		{
			name:      "Selected fields with nested with omitempty",
			input:     &s2,
			omitempty: true,
			fields:    []string{"A", "B"},
			expect: Fields{
				{Name: "A", Value: "foo", Addr: &s2.A},
				{Name: "A", Value: "test", Addr: &s2.simple.A},
				{Name: "A", Value: "test", Addr: &s2.S.A},
			},
		},
		{
			name:  "All fields with skip",
			input: &s3,
			expect: Fields{
				{Name: "C", Value: 3.14, Addr: &s3.C},
				{Name: "A", Value: "test", Addr: &s3.simple.A},
				{Name: "B", Value: 0, Addr: &s3.simple.B},
				{Name: "C", Value: 3.14, Addr: &s3.simple.C},
				{Name: "D", Value: 28, Addr: &s3.simple.D},
				{Name: "A", Value: "test", Addr: &s3.S.A},
				{Name: "B", Value: 0, Addr: &s3.S.B},
				{Name: "C", Value: 3.14, Addr: &s3.S.C},
				{Name: "D", Value: 28, Addr: &s3.S.D},
			},
		},
		{
			name:      "All fields with skip with omitempty",
			input:     &s3,
			omitempty: true,
			expect: Fields{
				{Name: "C", Value: 3.14, Addr: &s3.C},
				{Name: "A", Value: "test", Addr: &s3.simple.A},
				{Name: "C", Value: 3.14, Addr: &s3.simple.C},
				{Name: "D", Value: 28, Addr: &s3.simple.D},
				{Name: "A", Value: "test", Addr: &s3.S.A},
				{Name: "C", Value: 3.14, Addr: &s3.S.C},
				{Name: "D", Value: 28, Addr: &s3.S.D},
			},
		},
		{
			name:   "Selected fields with skip",
			input:  &s3,
			fields: []string{"A", "C"},
			expect: Fields{
				{Name: "C", Value: 3.14, Addr: &s3.C},
				{Name: "A", Value: "test", Addr: &s3.simple.A},
				{Name: "C", Value: 3.14, Addr: &s3.simple.C},
				{Name: "A", Value: "test", Addr: &s3.S.A},
				{Name: "C", Value: 3.14, Addr: &s3.S.C},
			},
		},
		{
			name:      "Selected fields with skip with omitempty",
			input:     &s3,
			omitempty: true,
			fields:    []string{"A", "B"},
			expect: Fields{
				{Name: "A", Value: "test", Addr: &s3.simple.A},
				{Name: "A", Value: "test", Addr: &s3.S.A},
			},
		},
		{
			name:  "Nil pointer field",
			input: &s4,
			expect: Fields{
				{Name: "A", Value: (*int)(nil), Addr: &s4.A},
			},
		},
		{
			name:      "Nil pointer field",
			input:     &s4,
			omitempty: true,
			expect:    Fields{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := tag.Fields(tt.input, tt.omitempty, tt.fields...)
			if tt.err != nil {
				equal(t, tt.err.Error(), err.Error())
				return
			} else {
				equal(t, nil, err)
			}
			equal(t, len(tt.expect), len(fields))
			for i, f := range fields {
				equal(t, tt.expect[i].Name, f.Name)
				equal(t, tt.expect[i].Value, f.Value)
				equal(t, tt.expect[i].Addr, f.Addr)
			}
		})
	}
}

func TestFields_Names(t *testing.T) {
	type ts struct {
		name   string
		input  any
		expect []string
	}

	tests := []ts{
		{
			name:   "Get names",
			input:  &s2,
			expect: []string{"A", "B", "name", "A", "B", "C", "D", "A", "B", "C", "D"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := tag.Fields(tt.input, false)
			equal(t, nil, err)

			names := fields.Names()
			equal(t, tt.expect, names)
		})
	}
}

func TestFields_Values(t *testing.T) {
	type ts struct {
		name   string
		input  any
		expect []any
	}

	tests := []ts{
		{
			name:   "Get values",
			input:  &s2,
			expect: []any{"foo", 0, 3.14, "test", 0, 3.14, 28, "test", 0, 3.14, 28},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := tag.Fields(tt.input, false)
			equal(t, nil, err)

			values := fields.Values()
			equal(t, tt.expect, values)
		})
	}
}

func TestFields_Pointers(t *testing.T) {
	type ts struct {
		name   string
		input  any
		expect []any
	}

	tests := []ts{
		{
			name:   "Get pointers",
			input:  &s2,
			expect: []any{&s2.A, &s2.B, &s2.C, &s2.simple.A, &s2.simple.B, &s2.simple.C, &s2.simple.D, &s2.S.A, &s2.S.B, &s2.S.C, &s2.S.D},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := tag.Fields(tt.input, false)
			equal(t, nil, err)

			pointers := fields.Pointers()
			equal(t, tt.expect, pointers)
		})
	}
}

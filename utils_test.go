package zcalendar

import (
	"reflect"
	"testing"
)

type ParserTestCase struct {
	name string
	in   string
	out  interface{}
	err  bool
}

var (
	TypeError  = reflect.TypeOf((*error)(nil)).Elem()
	TypeString = reflect.TypeOf(string(""))
)

func testParser(t *testing.T, parserFunc interface{}, cases []ParserTestCase) {
	funcT := reflect.TypeOf(parserFunc)

	if funcT.Kind() != reflect.Func {
		t.Fatalf("parserFunc is not a func, %s given", funcT.Kind())
	}

	if funcT.NumIn() != 1 {
		t.Fatalf("parserFunc should have 1 input parameter, have %d", funcT.NumIn())
	}

	if funcT.In(0) != TypeString {
		t.Fatalf("parserFunc input parameter should be a string, %s given", funcT.In(0))
	}

	if funcT.NumOut() != 2 {
		t.Fatalf("parserFunc should have 2 output parameters, have %d", funcT.NumOut())
	}

	if !funcT.Out(1).Implements(TypeError) {
		t.Fatalf("parserFunc second output parameter should be an error, %s given", funcT.Out(1))
	}

	funcV := reflect.ValueOf(parserFunc)
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			values := funcV.Call([]reflect.Value{reflect.ValueOf(c.in)})

			out := values[0].Interface()
			var err error
			if !values[1].IsNil() {
				err = values[1].Interface().(error)
			}

			if c.err != (err != nil) {
				t.Errorf("unexpected error: got %v", err)
			}

			if c.err {
				return
			}

			if !reflect.DeepEqual(c.out, out) {
				t.Errorf("unexpected value: wanted %v, got %v", c.out, out)
			}
		})
	}
}

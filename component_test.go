package zcalendar

import (
	"reflect"
	"testing"
)

func TestParseValue(t *testing.T) {
	testParser(t, parseValue, []ParserTestCase{
		{name: "valid value", in: "1", out: component{From: 1}},
		{name: "valid repeat", in: "1/2", out: component{From: 1, Repeat: 2}},
		{name: "invalid value 1", in: "a", err: true},
		{name: "invalid value 2", in: "-1", err: true},
		{name: "invalid value 4", in: "", err: true},
		{name: "invalid repeat 1", in: "1/-2", err: true},
		{name: "invalid repeat 2", in: "1/a", err: true},
		{name: "invalid repeat 3", in: "1/2/a", err: true},
		{name: "invalid repeat 4", in: "1/", err: true},
	})
}

func TestParseRange(t *testing.T) {
	testParser(t, parseRange, []ParserTestCase{
		{name: "valid value", in: "1..2", out: component{From: 1, To: 2}},
		{name: "valid repeat", in: "1..2/3", out: component{From: 1, To: 2, Repeat: 3}},
		{name: "invalid range 1", in: "1..", err: true},
		{name: "invalid range 2", in: "..2", err: true},
		{name: "invalid range 3", in: "1..2..3", err: true},
		{name: "invalid range 4", in: "..", err: true},
		{name: "invalid range 5", in: "", err: true},
		{name: "invalid value 1", in: "a..2", err: true},
		{name: "invalid value 2", in: "-1..2", err: true},
		{name: "invalid repeat 1", in: "1..2/-2", err: true},
		{name: "invalid repeat 2", in: "1..2/a", err: true},
		{name: "invalid repeat 3", in: "1..2/3/a", err: true},
		{name: "invalid repeat 4", in: "1..2/", err: true},
		{name: "invalid bounds", in: "2..1", err: true},
	})
}

func TestParseComponents(t *testing.T) {
	testParser(t, parseComponents, []ParserTestCase{
		{name: "valid components", in: "1,3..4", out: components{{From: 1}, {From: 3, To: 4}}},
		{name: "invalid component 1", in: "1azefdwsef,3..4", err: true},
		{name: "invalid component 2", in: "1,3..4/3/4", err: true},
		{name: "empty component 1", in: "1,3..4,", err: true},
		{name: "empty component 2", in: "1,,3..4", err: true},
	})
}

func TestComponents_Values(t *testing.T) {
	type Case struct {
		name  string
		comps components
		out   []int
	}

	// We assume that the maximum value is set to 10 for simplicity's sake.
	for _, c := range []Case{
		{name: "single component", comps: components{{From: 1}}, out: []int{1}},
		{name: "multiple components", comps: components{{From: 1}, {From: 2}}, out: []int{1, 2}},
		{name: "no duplicates", comps: components{{From: 1, To: 4}, {From: 2, To: 5}}, out: []int{1, 2, 3, 4, 5}},
		{name: "no component", comps: components{}, out: []int{}},
	} {
		t.Run(c.name, func(t *testing.T) {
			out := c.comps.Values(10)
			if !reflect.DeepEqual(c.out, out) {
				t.Errorf("unexpected output: wanted %v, got %v", c.out, out)
			}
		})
	}
}

func TestComponents_Next(t *testing.T) {
	type Case struct {
		name  string
		comps components
		out   int
		ok    bool
	}

	// We assume that the maximum value is set to 10 for simplicity's sake.
	for _, c := range []Case{
		{name: "single value", comps: components{{From: 1}}, out: 1, ok: true},
		{name: "next value", comps: components{{From: 1, To: 9}}, out: 7, ok: true},
		{name: "no value", comps: components{}, out: 0, ok: false},
	} {
		t.Run(c.name, func(t *testing.T) {
			out, _, ok := c.comps.Next(7, 10)

			if ok != c.ok {
				t.Errorf("unexpected result: wanted %v, got %v", c.ok, ok)
			}

			if out != c.out {
				t.Errorf("unexpected output: wanted %v, got %v", c.out, out)
			}
		})
	}
}

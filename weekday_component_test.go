package zcalendar

import "testing"

func TestParseWeekdayValue(t *testing.T) {
	testParser(t, parseWeekdayValue, []ParserTestCase{
		{name: "valid short weekday 1", in: "Mon", out: weekdayComponent{From: 1}},
		{name: "valid short weekday 2", in: "Tue", out: weekdayComponent{From: 2}},
		{name: "valid short weekday 3", in: "Wed", out: weekdayComponent{From: 3}},
		{name: "valid short weekday 4", in: "Thu", out: weekdayComponent{From: 4}},
		{name: "valid short weekday 5", in: "Fri", out: weekdayComponent{From: 5}},
		{name: "valid short weekday 6", in: "Sat", out: weekdayComponent{From: 6}},
		{name: "valid short weekday 7", in: "Sun", out: weekdayComponent{From: 7}},
		{name: "valid weekday 1", in: "Monday", out: weekdayComponent{From: 1}},
		{name: "valid weekday 2", in: "Tuesday", out: weekdayComponent{From: 2}},
		{name: "valid weekday 3", in: "Wednesday", out: weekdayComponent{From: 3}},
		{name: "valid weekday 4", in: "Thursday", out: weekdayComponent{From: 4}},
		{name: "valid weekday 5", in: "Friday", out: weekdayComponent{From: 5}},
		{name: "valid weekday 6", in: "Saturday", out: weekdayComponent{From: 6}},
		{name: "valid weekday 7", in: "Sunday", out: weekdayComponent{From: 7}},
		{name: "valid lowercase short weekday 1", in: "mon", out: weekdayComponent{From: 1}},
		{name: "valid lowercase short weekday 2", in: "tue", out: weekdayComponent{From: 2}},
		{name: "valid lowercase short weekday 3", in: "wed", out: weekdayComponent{From: 3}},
		{name: "valid lowercase short weekday 4", in: "thu", out: weekdayComponent{From: 4}},
		{name: "valid lowercase short weekday 5", in: "fri", out: weekdayComponent{From: 5}},
		{name: "valid lowercase short weekday 6", in: "sat", out: weekdayComponent{From: 6}},
		{name: "valid lowercase short weekday 7", in: "sun", out: weekdayComponent{From: 7}},
		{name: "valid lowercase weekday 1", in: "monday", out: weekdayComponent{From: 1}},
		{name: "valid lowercase weekday 2", in: "tuesday", out: weekdayComponent{From: 2}},
		{name: "valid lowercase weekday 3", in: "wednesday", out: weekdayComponent{From: 3}},
		{name: "valid lowercase weekday 4", in: "thursday", out: weekdayComponent{From: 4}},
		{name: "valid lowercase weekday 5", in: "friday", out: weekdayComponent{From: 5}},
		{name: "valid lowercase weekday 6", in: "saturday", out: weekdayComponent{From: 6}},
		{name: "valid lowercase weekday 7", in: "sunday", out: weekdayComponent{From: 7}},
		{name: "invalid weekday 1", in: "Lundi", err: true},
		{name: "invalid weekday 2", in: "", err: true},
	})
}

func TestParseWeekdayRange(t *testing.T) {
	testParser(t, parseWeekdayRange, []ParserTestCase{
		{name: "valid range 1", in: "Mon..Tue", out: weekdayComponent{From: 1, To: 2}},
		{name: "valid range 2", in: "Monday..Tuesday", out: weekdayComponent{From: 1, To: 2}},
		{name: "valid range 3", in: "Monday..Fri", out: weekdayComponent{From: 1, To: 5}},
		{name: "invalid range 1", in: "Mon..Abe", err: true},
		{name: "invalid range 2", in: "Cjfh..Friday", err: true},
		{name: "invalid bounds", in: "Wed..Mon", err: true},
	})
}

func TestParseweekdayComponents(t *testing.T) {
	testParser(t, parseWeekdayComponents, []ParserTestCase{
		{name: "valid components", in: "Mon,Wed..Thu", out: weekdayComponents{{From: 1}, {From: 3, To: 4}}},
		{name: "invalid component 1", in: "Lundi,Wed..Thu", err: true},
		{name: "empty component 1", in: "Mon,Wed..Thu,", err: true},
		{name: "empty component 2", in: "Mon,,Wed..Thu", err: true},
	})
}

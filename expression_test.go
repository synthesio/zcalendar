package zcalendar

import (
	"reflect"
	"testing"
	"time"
)

// Those are the values used in the tests.
var EuropeParis, Zulu *time.Location

func init() {
	EuropeParis, _ = time.LoadLocation("Europe/Paris")
	Zulu, _ = time.LoadLocation("Zulu")
}

func TestParse(t *testing.T) {
	testParser(t, Parse, []ParserTestCase{
		{name: "valid expression", in: "Mon 2006-01-02 15:04:05 Europe/Paris", out: Expression{
			weekdays: []weekdayComponent{{From: 1}},
			years:    []component{{From: 2006}},
			months:   []component{{From: 1}},
			days:     []component{{From: 2}},
			hours:    []component{{From: 15}},
			minutes:  []component{{From: 4}},
			seconds:  []component{{From: 5}},
			timezone: EuropeParis,
		}},
		{name: "optional year", in: "Mon 01-02 15:04:05 Europe/Paris", out: Expression{
			weekdays: []weekdayComponent{{From: 1}},
			years:    defaultYears,
			months:   []component{{From: 1}},
			days:     []component{{From: 2}},
			hours:    []component{{From: 15}},
			minutes:  []component{{From: 4}},
			seconds:  []component{{From: 5}},
			timezone: EuropeParis,
		}},
		{name: "optional seconds", in: "Mon 2006-01-02 15:04 Europe/Paris", out: Expression{
			weekdays: []weekdayComponent{{From: 1}},
			years:    []component{{From: 2006}},
			months:   []component{{From: 1}},
			days:     []component{{From: 2}},
			hours:    []component{{From: 15}},
			minutes:  []component{{From: 4}},
			seconds:  defaultSeconds,
			timezone: EuropeParis,
		}},
		{name: "wildcards", in: "Mon *-*-* *:*:* Europe/Paris", out: Expression{
			weekdays: []weekdayComponent{{From: 1}},
			years:    allYears,
			months:   allMonths,
			days:     allDays,
			hours:    allHours,
			minutes:  allMinutes,
			seconds:  allSeconds,
			timezone: EuropeParis,
		}},
		{name: "optional timezone", in: "Mon 2006-01-02 15:04:05", out: Expression{
			weekdays: []weekdayComponent{{From: 1}},
			years:    []component{{From: 2006}},
			months:   []component{{From: 1}},
			days:     []component{{From: 2}},
			hours:    []component{{From: 15}},
			minutes:  []component{{From: 4}},
			seconds:  []component{{From: 5}},
			timezone: defaulttimezone,
		}},
		{name: "optional weekdays", in: "2006-01-02 15:04:05 Europe/Paris", out: Expression{
			weekdays: defaultWeekdays,
			years:    []component{{From: 2006}},
			months:   []component{{From: 1}},
			days:     []component{{From: 2}},
			hours:    []component{{From: 15}},
			minutes:  []component{{From: 4}},
			seconds:  []component{{From: 5}},
			timezone: EuropeParis,
		}},
		{name: "optional date", in: "Mon 15:04:05 Europe/Paris", out: Expression{
			weekdays: []weekdayComponent{{From: 1}},
			years:    defaultYears,
			months:   defaultMonths,
			days:     defaultDays,
			hours:    []component{{From: 15}},
			minutes:  []component{{From: 4}},
			seconds:  []component{{From: 5}},
			timezone: EuropeParis,
		}},
		{name: "optional time", in: "Mon 2006-01-02 Europe/Paris", out: Expression{
			weekdays: []weekdayComponent{{From: 1}},
			years:    []component{{From: 2006}},
			months:   []component{{From: 1}},
			days:     []component{{From: 2}},
			hours:    defaultHours,
			minutes:  defaultMinutes,
			seconds:  defaultSeconds,
			timezone: EuropeParis,
		}},
		{name: "UTC timezone", in: "Mon 2006-01-02 15:04:05 UTC", out: Expression{
			weekdays: []weekdayComponent{{From: 1}},
			years:    []component{{From: 2006}},
			months:   []component{{From: 1}},
			days:     []component{{From: 2}},
			hours:    []component{{From: 15}},
			minutes:  []component{{From: 4}},
			seconds:  []component{{From: 5}},
			timezone: time.UTC,
		}},
		{name: "Zulu timezone", in: "Mon 2006-01-02 15:04:05 Zulu", out: Expression{
			weekdays: []weekdayComponent{{From: 1}},
			years:    []component{{From: 2006}},
			months:   []component{{From: 1}},
			days:     []component{{From: 2}},
			hours:    []component{{From: 15}},
			minutes:  []component{{From: 4}},
			seconds:  []component{{From: 5}},
			timezone: Zulu,
		}},
		{name: "date only", in: "2006-01-02", out: Expression{
			weekdays: defaultWeekdays,
			years:    []component{{From: 2006}},
			months:   []component{{From: 1}},
			days:     []component{{From: 2}},
			hours:    defaultHours,
			minutes:  defaultMinutes,
			seconds:  defaultSeconds,
			timezone: defaulttimezone,
		}},
		{name: "time only", in: "15:04:05", out: Expression{
			weekdays: defaultWeekdays,
			years:    defaultYears,
			months:   defaultMonths,
			days:     defaultDays,
			hours:    []component{{From: 15}},
			minutes:  []component{{From: 4}},
			seconds:  []component{{From: 5}},
			timezone: defaulttimezone,
		}},
		{name: "weekdays only", in: "Mon", out: Expression{
			weekdays: []weekdayComponent{{From: 1}},
			years:    defaultYears,
			months:   defaultMonths,
			days:     defaultDays,
			hours:    defaultHours,
			minutes:  defaultMinutes,
			seconds:  defaultSeconds,
			timezone: defaulttimezone,
		}},
		{name: "empty expression", in: "", err: true},
		{name: "not an expression", in: "les sanglots longs des violons de l'automne", err: true},
		{name: "timezone only", in: "Europe/Paris", err: true},
		{name: "invalid timezone", in: "Mon 2006-01-02 15:04:05 hello", err: true},
		{name: "too many chunks", in: "Mon 2006-01-02 15:04:05 UTC hello", err: true},
		{name: "chunk after timezone", in: "Mon 15:04:05 UTC hello", err: true},
	})
}

type MarshalTestCase struct {
	name string
	in   Expression
	out  string
}

func TestExpression_MarshalText(t *testing.T) {
	var cases = []MarshalTestCase{
		{
			name: "valid expression",
			in: Expression{
				weekdays: []weekdayComponent{{From: 1}},
				years:    []component{{From: 2006}},
				months:   []component{{From: 1}},
				days:     []component{{From: 2}},
				hours:    []component{{From: 15}},
				minutes:  []component{{From: 4}},
				seconds:  []component{{From: 5}},
				timezone: EuropeParis,
			},
			out: "Mon 2006-01-02 15:04:05 Europe/Paris",
		},
		{
			name: "local timezone",
			in: Expression{
				weekdays: []weekdayComponent{{From: 1}},
				years:    []component{{From: 2006}},
				months:   []component{{From: 1}},
				days:     []component{{From: 2}},
				hours:    []component{{From: 15}},
				minutes:  []component{{From: 4}},
				seconds:  []component{{From: 5}},
				timezone: time.Local,
			},
			out: "Mon 2006-01-02 15:04:05",
		},
		{
			name: "full ranges",
			in: Expression{
				weekdays: allWeekdays,
				years:    allYears,
				months:   allMonths,
				days:     allDays,
				hours:    allHours,
				minutes:  allMinutes,
				seconds:  allSeconds,
				timezone: time.UTC,
			},
			out: "*-*-* *:*:* UTC",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			out, err := c.in.MarshalText()

			if err != nil {
				t.Errorf("unexpected error: got %v", err)
				return
			}

			if !reflect.DeepEqual([]byte(c.out), out) {
				t.Errorf("unexpected output: wanted %s, got %s", c.out, string(out))
			}
		})
	}
}

func TestExpression_Next(t *testing.T) {
	var current = time.Date(2006, 01, 02, 15, 04, 05, 0, time.UTC)

	type Case struct {
		name  string
		exp   string
		next  string
		found bool
	}

	for _, c := range []Case{
		{name: "next year", exp: "*-01-01 00:00:00 UTC", next: "2007-01-01T00:00:00Z", found: true},
		{name: "next month", exp: "*-*-01 00:00:00 UTC", next: "2006-02-01T00:00:00Z", found: true},
		{name: "next day", exp: "*-*-* 00:00:00 UTC", next: "2006-01-03T00:00:00Z", found: true},
		{name: "no next date", exp: "2005-*-* 00:00:00 UTC", next: "2006-01-03T00:00:00Z", found: false},
		{name: "next monday", exp: "Mon 00:00:00 UTC", next: "2006-01-09T00:00:00Z", found: true},
		{name: "next sunday", exp: "Sun 00:00:00 UTC", next: "2006-01-08T00:00:00Z", found: true},
	} {
		t.Run(c.name, func(t *testing.T) {
			exp, err := Parse(c.exp)
			if err != nil {
				t.Fatalf("unexpected error parsing expression: %s", err)
			}

			next, err := time.Parse(time.RFC3339, c.next)
			if err != nil {
				t.Fatalf("unexpected error parsing next time: %s", err)
			}

			out, ok := exp.Next(current)
			if ok != c.found {
				t.Fatalf("unexpected found output: wanted %v, got %v", c.found, ok)
			}

			if !ok {
				return
			}

			if !reflect.DeepEqual(next, out) {
				t.Fatalf("unexpected time output: wanted %v, got %v", next, out)
			}
		})
	}
}

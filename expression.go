package zcalendar

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// An Expression is the Go representation of a Calendar Event as per Systemd's
// specification (with some exceptions, see the Parse method).
type Expression struct {
	// First part of the expression is day of week.
	weekdays weekdayComponents

	// Second part of the expression is the date part.
	years  components
	months components
	days   components

	// Third part of the expression is the time part.
	hours   components
	minutes components
	seconds components

	// Fourth part of the expression is the timezone.
	timezone *time.Location
}

// Non unit-related boundaries.
var (
	MaxYears = 2199
	MinYears = 1970
)

// The full-range values for easy comparison.
var (
	allWeekdays = weekdayComponents{{From: 1, To: 7}}
	allYears    = components{{From: 1970, To: 2199}}
	allMonths   = components{{From: 1, To: 12}}
	allDays     = components{{From: 1, To: 31}}
	allHours    = components{{From: 0, To: 23}}
	allMinutes  = components{{From: 0, To: 59}}
	allSeconds  = components{{From: 0, To: 59}}
)

// The default values for easy manipulation.
var (
	defaultWeekdays = weekdayComponents{{From: 1, To: 7}}
	defaultYears    = components{{From: 1970, To: 2199}}
	defaultMonths   = components{{From: 1, To: 12}}
	defaultDays     = components{{From: 1, To: 31}}
	defaultHours    = components{{From: 0}}
	defaultMinutes  = components{{From: 0}}
	defaultSeconds  = components{{From: 0}}
	defaulttimezone = time.Local
)

// Parse a raw string into an expression. Follows Systemd's Calendar Events
// specification with some exceptions:
// - Any timezone can be specified, not only UTC and local
// - Sub-second aren't handled
// - The end-of-month token isn't handled
//
// Original implementation can be found here: https://github.com/systemd/systemd/blob/master/src/basic/calendarspec.c#L879
func Parse(raw string) (exp Expression, err error) {
	// By default, set all fields to the largest range available.
	exp = Expression{
		weekdays: defaultWeekdays,
		years:    defaultYears,
		months:   defaultMonths,
		days:     defaultDays,
		hours:    defaultHours,
		minutes:  defaultMinutes,
		seconds:  defaultSeconds,
		timezone: defaulttimezone,
	}

	chunks := strings.Fields(raw)

	// If there is no chunk to handle, the expression is composed of
	// whitespaces and thus invalid.
	if len(chunks) == 0 {
		return exp, errors.New("empty expression")
	}

	// If there is more than 4 chunks, the expression has whitespaces at
	// the wrong places, or is simply not an expression.
	if len(chunks) > 4 {
		return exp, errors.New("too many components")
	}

	// TODO Handle shortcuts.

	// If the first chunk has a neither a dash or a comma, then it can't be
	// a date or time, and a timezone can't be the first item, so it has to
	// be weekdays.
	if !strings.ContainsAny(chunks[0], "-:") {
		exp.weekdays, err = parseWeekdayComponents(chunks[0])
		if err != nil {
			return exp, fmt.Errorf(`parsing weekdays: %w`, err)
		}

		// If the first chunk is successfully parsed, shift it out of
		// the stack so the rest of the steps always work on the first
		// chunk of the stack.
		chunks = chunks[1:]
	}

	// If the first chunk contains a dash, it must be a date.
	if len(chunks) != 0 && strings.Contains(chunks[0], "-") {
		parts := strings.Split(chunks[0], "-")

		// A date is composed a most of 3 parts: years, months, days.
		// There is no need to check for the one part case, as it
		// wouldn't have any dash in it, and thus wouldn't enter this
		// case.
		if len(parts) > 3 {
			return exp, errors.New("invalid parts component")
		}

		// The year is optional, so add it if missing.
		if len(parts) == 2 {
			parts = append([]string{"*"}, parts...)
		}

		if parts[0] != "*" {
			exp.years, err = parseComponents(parts[0])
			if err != nil {
				return exp, fmt.Errorf(`parsing years: %w`, err)
			}
		}

		if parts[1] != "*" {
			exp.months, err = parseComponents(parts[1])
			if err != nil {
				return exp, fmt.Errorf(`parsing months: %w`, err)
			}
		}

		if parts[2] != "*" {
			exp.days, err = parseComponents(parts[2])
			if err != nil {
				return exp, fmt.Errorf(`parsing days: %w`, err)
			}
		}

		chunks = chunks[1:]
	}

	// If the first chunk contains a comma, it myst be a time.
	if len(chunks) != 0 && strings.Contains(chunks[0], ":") {
		parts := strings.Split(chunks[0], ":")

		// A time is composed at most of 3 parts: hours, minutes,
		// seconds. There is no need to check for the one part case, as
		// for the date chunk.
		if len(parts) > 3 {
			return exp, errors.New("invalid time component")
		}

		// Seconds are optional, so add them if missing.
		if len(parts) == 2 {
			parts = append(parts, "00")
		}

		if parts[0] == "*" {
			exp.hours = allHours
		} else {
			exp.hours, err = parseComponents(parts[0])
			if err != nil {
				return exp, fmt.Errorf(`parsing hours: %w`, err)
			}
		}

		if parts[1] == "*" {
			exp.minutes = allMinutes
		} else {
			exp.minutes, err = parseComponents(parts[1])
			if err != nil {
				return exp, fmt.Errorf(`parsing minutes: %w`, err)
			}
		}

		if parts[2] == "*" {
			exp.seconds = allSeconds
		} else {
			exp.seconds, err = parseComponents(parts[2])
			if err != nil {
				return exp, fmt.Errorf(`parsing seconds: %w`, err)
			}
		}

		chunks = chunks[1:]
	}

	// If there is still a chunk in the stack at this point it must be a
	// timezone.
	if len(chunks) != 0 {
		exp.timezone, err = time.LoadLocation(chunks[0])
		if err != nil {
			return exp, fmt.Errorf("invalid chunk %s", chunks[0])
		}

		chunks = chunks[1:]
	}

	// At this point, remaining items indicate unparsable chunks.
	if len(chunks) != 0 {
		return exp, fmt.Errorf("invalid chunk %s", chunks[0])
	}

	return exp, nil
}

// MustParse is like Parse but will panic in case of error.
func MustParse(raw string) (e Expression) {
	e, err := Parse(raw)
	if err != nil {
		panic(err)
	}

	return e
}

// UnmarshalText implements the encoding.TextUnmarshaler interface so an
// expression can unmarshalled from a JSON object.
func (e *Expression) UnmarshalText(raw []byte) (err error) {
	*e, err = Parse(string(raw))
	return err
}

// MarshalText implement the encoding.TextMarshaler interface.
func (e Expression) MarshalText() (text []byte, err error) {
	var buf bytes.Buffer

	// If there is actually a weekdays specification, write all parts.
	if !reflect.DeepEqual(e.weekdays, allWeekdays) {
		buf.WriteString(e.weekdays.String())
		buf.WriteString(" ")
	}

	if reflect.DeepEqual(e.years, allYears) {
		buf.WriteString("*")
	} else {
		buf.WriteString(e.years.String())
	}
	buf.WriteString("-")

	if reflect.DeepEqual(e.months, allMonths) {
		buf.WriteString("*")
	} else {
		buf.WriteString(e.months.String())
	}
	buf.WriteString("-")

	if reflect.DeepEqual(e.days, allDays) {
		buf.WriteString("*")
	} else {
		buf.WriteString(e.days.String())
	}
	buf.WriteString(" ")

	if reflect.DeepEqual(e.hours, allHours) {
		buf.WriteString("*")
	} else {
		buf.WriteString(e.hours.String())
	}
	buf.WriteString(":")

	if reflect.DeepEqual(e.minutes, allMinutes) {
		buf.WriteString("*")
	} else {
		buf.WriteString(e.minutes.String())
	}
	buf.WriteString(":")

	if reflect.DeepEqual(e.seconds, allSeconds) {
		buf.WriteString("*")
	} else {
		buf.WriteString(e.seconds.String())
	}

	if e.timezone != time.Local {
		buf.WriteString(" ")
		buf.WriteString(e.timezone.String())
	}

	return buf.Bytes(), nil
}

// String implement the fmt.Stringer interface.
func (e Expression) String() string {
	bytes, _ := e.MarshalText()
	return string(bytes)
}

// Scan implements the sql.Scanner interface, which allow to use an Expression
// as a database field and scan it.
func (e *Expression) Scan(src interface{}) (err error) {
	var raw []byte

	switch src := src.(type) {
	case []byte:
		raw = src
	case string:
		raw = []byte(src)
	default:
		return fmt.Errorf("unable to scan %T into %T", raw, e)
	}

	return e.UnmarshalText(raw)
}

// Value implements the driver.Value interface, which allows a sql database
// driver to insert it into a row.
func (e Expression) Value() (val driver.Value, err error) {
	return e.MarshalText()
}

// Next returns the next point in time that will satisfy the schedule that is
// strictly after d.
//
// Original implementation can be found here:
// https://github.com/systemd/systemd/blob/master/src/basic/calendarspec.c#L1199
func (e Expression) Next(d time.Time) (n time.Time, ok bool) {
	d = d.In(e.timezone)

	var (
		year   = d.Year()
		month  = int(d.Month())
		day    = d.Day()
		hour   = d.Hour()
		minute = d.Minute()
		second = d.Second() + 1

		diff int
	)

	// The loop works as follow: each unit is initialized with the value
	// from d. To prevent d being returned in the case it is a valid value,
	// the seconds are incremented so the returned value is necessarily
	// different from d.
	//
	// For each unit from the bigest to the smallest, get the next value
	// allowed by the expression. From this point, there is 3 possibilities:
	//
	// - If this value is equal, skip to the next unit,
	// - If the next value is bigger than the current one, reset the lower
	//   units to the first value,
	// - If the next value is  smaller than the current value, we increment
	//   the unit before, reset the lower units to the first value, and start
	//   over;
	//
	// When we reach the end of the loop, we can safely break out and
	// return the actual values as the next date.
	for {
		year, diff, ok = e.years.Next(year, MaxYears)
		if !ok {
			return
		}

		if diff < 0 {
			ok = false
			return
		}

		if diff > 0 {
			month = 1
			day = 1
			hour = 0
			minute = 0
			second = 0
		}

		month, diff, ok = e.months.Next(month, 12)
		if !ok {
			return
		}

		if diff < 0 {
			year++
			day = 1
			hour = 0
			minute = 0
			second = 0
			continue
		}

		if diff > 0 {
			day = 1
			hour = 0
			minute = 0
			second = 0
		}

		daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()

		day, diff, ok = e.days.Next(day, daysInMonth)
		if !ok {
			return
		}

		if diff < 0 {
			month++
			hour = 0
			minute = 0
			second = 0
			continue
		}

		if diff > 0 {
			hour = 0
			minute = 0
			second = 0
		}

		weekday := int(time.Date(year, time.Month(month), day, 0, 0, 0, 0, e.timezone).Weekday())
		// Go's weekdays range is Sunday=0..Saturday=6, while our weekdays are Monday=1..Sunday=7
		if weekday == 0 {
			weekday = 7
		}
		if !e.weekdays.Contains(weekday) {
			day++
			hour = 0
			minute = 0
			second = 0
			continue
		}

		hour, diff, ok = e.hours.Next(hour, 23)
		if !ok {
			return
		}

		if diff < 0 {
			day++
			minute = 0
			second = 0
			continue
		}

		if diff > 0 {
			minute = 0
			second = 0
		}

		minute, diff, ok = e.minutes.Next(minute, 59)
		if !ok {
			return
		}

		if diff < 0 {
			hour++
			second = 0
			continue
		}

		if diff > 0 {
			second = 0
		}

		second, diff, ok = e.seconds.Next(second, 59)
		if !ok {
			return
		}

		if diff < 0 {
			minute++
			continue
		}

		break
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, e.timezone), true
}

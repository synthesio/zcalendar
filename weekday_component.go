package zcalendar

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// A weekdayComponent is a single range of weekdays.
type weekdayComponent struct {
	From int
	To   int
}

// weekdays list the valid values for the weekdays in the calendar spec.
var weekdaysValues = map[string]int{
	"monday":    1,
	"mon":       1,
	"tuesday":   2,
	"tue":       2,
	"wednesday": 3,
	"wed":       3,
	"thursday":  4,
	"thu":       4,
	"friday":    5,
	"fri":       5,
	"saturday":  6,
	"sat":       6,
	"sunday":    7,
	"sun":       7,
}

// weekdays list the valid values for the weekdays in the calendar spec.
var weekdaysStrings = map[int]string{
	1: "Mon",
	2: "Tue",
	3: "Wed",
	4: "Thu",
	5: "Fri",
	6: "Sat",
	7: "Sun",
}

// parseweekdayValue create a component from the string representation of a
// weekday.
func parseWeekdayValue(raw string) (c weekdayComponent, err error) {
	v, ok := weekdaysValues[strings.ToLower(raw)]
	if !ok {
		return c, errors.New("invalid weekday")
	}
	c.From = v

	return c, nil
}

// parseweekdayRange create a component from the string representation of range
// of weekdays.
func parseWeekdayRange(raw string) (c weekdayComponent, err error) {
	bounds := strings.Split(raw, "..")
	if len(bounds) != 2 {
		return c, errors.New("invalid range")
	}

	v, ok := weekdaysValues[strings.ToLower(bounds[0])]
	if !ok {
		return c, errors.New("invalid weekday")
	}
	c.From = v

	v, ok = weekdaysValues[strings.ToLower(bounds[1])]
	if !ok {
		return c, errors.New("invalid weekday")
	}
	c.To = v

	if c.From >= c.To {
		return c, errors.New("invalid bounds")
	}

	return c, nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (c weekdayComponent) MarshalText() (text []byte, err error) {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s", weekdaysStrings[c.From])

	if c.To != 0 {
		fmt.Fprintf(&buf, "..%s", weekdaysStrings[c.To])
	}

	return buf.Bytes(), nil
}

type weekdayComponents []weekdayComponent

// parseweekdayComponents create a slice of components from a string representing a
// comma-separated list of weekdays values and ranges.
func parseWeekdayComponents(raw string) (cs weekdayComponents, err error) {
	for index, chunk := range strings.Split(raw, ",") {
		if strings.Contains(chunk, "..") {
			c, err := parseWeekdayRange(chunk)
			if err != nil {
				return cs, fmt.Errorf(`parsing range %d: %w`, index, err)
			}
			cs = append(cs, c)
			continue
		}

		c, err := parseWeekdayValue(chunk)
		if err != nil {
			return cs, fmt.Errorf(`parsing value %d: %w`, index, err)
		}
		cs = append(cs, c)
	}

	return cs, err
}

// MarshalText implements the encoding.MarshalText interface for a Component
// slice.
func (cs weekdayComponents) MarshalText() (text []byte, err error) {
	var parts [][]byte
	for _, c := range cs {
		b, _ := c.MarshalText()
		parts = append(parts, b)
	}
	return bytes.Join(parts, []byte(",")), nil
}

func (cs weekdayComponents) String() string {
	b, _ := cs.MarshalText()
	return string(b)
}

// Values return the list of actual values from the various sub-components.
func (cs weekdayComponents) Values() (values []int) {
	var seen = make(map[int]struct{})

	for _, c := range cs {
		if c.To == 0 {
			seen[c.From] = struct{}{}
		} else {
			for v := c.From; v <= c.To && v <= 7; v++ {
				seen[v] = struct{}{}
			}
		}
	}

	values = make([]int, 0, len(seen))
	for k := range seen {
		values = append(values, k)
	}
	sort.Ints(values)

	return
}

func (cs weekdayComponents) Contains(day int) (ok bool) {
	for _, v := range cs.Values() {
		if v == day {
			return true
		}
	}
	return false
}

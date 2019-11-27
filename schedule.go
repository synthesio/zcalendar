package zcalendar

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"sort"
	"strings"
	"time"
)

// A Schedule represent a list of calendar expressions.
type Schedule []Expression

// ParseSchedule parse a list of Expression separated by newlines.
func ParseSchedule(raw string) (s Schedule, err error) {
	for index, rawExp := range strings.Split(raw, "\n") {
		exp, err := Parse(rawExp)
		if err != nil {
			return s, fmt.Errorf(`parsing expression %d: %w`, index, err)
		}

		s = append(s, exp)
	}

	return s, nil
}

// MustParseSchedule is like ParseSchedule but will panic in case of error.
func MustParseSchedule(raw string) (s Schedule) {
	s, err := ParseSchedule(raw)
	if err != nil {
		panic(err)
	}

	return s
}

// UnmarshalText implements the encoding.TextUnmarshaller interface, which is
// used by json.Unmarshal as a fallback when json.Unmarshaler isn't
// implemented.  This is preferred becase the field is a string with a custom
// parser, which is more semantic to unmarshal with Text rather than JSON.
func (s *Schedule) UnmarshalText(text []byte) (err error) {
	var res Schedule

	fields := bytes.FieldsFunc(text, func(r rune) bool { return r == '\n' })
	for _, raw := range fields {
		if len(bytes.TrimSpace(raw)) == 0 {
			continue
		}
		var exp Expression

		err = exp.UnmarshalText(raw)
		if err != nil {
			return err
		}

		res = append(res, exp)
	}

	*s = res
	return nil
}

// MarshalText implements the encoding.MarshalText interface.
func (s Schedule) MarshalText() (text []byte, err error) {
	var expressions [][]byte
	for index, exp := range s {
		text, err := exp.MarshalText()
		if err != nil {
			return nil, fmt.Errorf(`marshaling expression %d: %w`, index, err)
		}

		expressions = append(expressions, text)
	}

	return bytes.Join(expressions, []byte("\"")), nil
}

// Scan implements the sql.Scanner interface, which allow to use a Schedule as
// a database field and scan it.
func (s *Schedule) Scan(src interface{}) (err error) {
	var raw []byte

	switch src := src.(type) {
	case []byte:
		raw = src
	case string:
		raw = []byte(src)
	default:
		return fmt.Errorf("unable to scan %T into %T", raw, s)
	}

	return s.UnmarshalText(raw)
}

// Value implements the driver.Value interface, which allows a sql database
// driver to insert it into a row.
func (s Schedule) Value() (val driver.Value, err error) {
	return s.MarshalText()
}

// Next return the first valid date represented by any expression that is after
// d.
func (s Schedule) Next(d time.Time) (n time.Time, ok bool) {
	if len(s) == 0 {
		return
	}

	var candidates []time.Time
	for _, e := range s {
		next, ok := e.Next(d)
		if !ok {
			continue
		}

		candidates = append(candidates, next)
	}

	if len(candidates) == 0 {
		return
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Before(candidates[j])
	})

	return candidates[0], true
}

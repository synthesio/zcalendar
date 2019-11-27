package zcalendar

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// A component is a single unit of the event expression. It represent a
// potentially repeating value or range in an unspecified time unit.
type component struct {
	From   int
	To     int
	Repeat int
}

// parseValue create a component from a string representing a simple value with
// an optional repetition.
func parseValue(raw string) (c component, err error) {
	var repeat = ""

	index := strings.Index(raw, "/")
	if index != -1 {
		raw, repeat = raw[:index], raw[index+1:]
	}

	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return c, fmt.Errorf(`invalid value: %w`, err)
	}
	if v < 0 {
		return c, errors.New("invalid negative value")
	}

	c.From = int(v)

	v, err = strconv.ParseInt(repeat, 10, 64)
	if index != -1 && err != nil {
		return c, fmt.Errorf(`invalid repeat: %w`, err)
	}
	if v < 0 {
		return c, errors.New("invalid negative repeat")
	}

	c.Repeat = int(v)

	return c, nil
}

// parseRange create a component from a string representing a range value with
// an optional repetition.
func parseRange(raw string) (c component, err error) {
	var repeat = ""

	index := strings.Index(raw, "/")
	if index != -1 {
		raw, repeat = raw[:index], raw[index+1:]
	}

	bounds := strings.Split(raw, "..")
	if len(bounds) != 2 {
		return c, errors.New("invalid range")
	}

	v, err := strconv.ParseInt(bounds[0], 10, 64)
	if err != nil {
		return c, fmt.Errorf(`invalid value: %w`, err)
	}
	if v < 0 {
		return c, errors.New("invalid negative lower bound")
	}
	c.From = int(v)

	v, err = strconv.ParseInt(bounds[1], 10, 64)
	if err != nil {
		return c, fmt.Errorf(`invalid value: %w`, err)
	}
	if v < 0 {
		return c, errors.New("invalid negative upper bound")
	}
	c.To = int(v)

	if c.From >= c.To {
		return c, errors.New("invalid bounds")
	}

	v, err = strconv.ParseInt(repeat, 10, 64)
	if index != -1 && err != nil {
		return c, fmt.Errorf(`invalid repeat: %w`, err)
	}
	if v < 0 {
		return c, errors.New("invalid negative repeat")
	}
	c.Repeat = int(v)

	return c, nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (c component) MarshalText() (text []byte, err error) {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%02d", c.From)

	if c.To != 0 {
		fmt.Fprintf(&buf, "..%02d", c.To)
	}

	if c.Repeat != 0 {
		fmt.Fprintf(&buf, "/%d", c.Repeat)
	}

	return buf.Bytes(), nil
}

type components []component

// parseComponents create a slice of components from a string representing a
// comma-separated list of values and ranges.
func parseComponents(raw string) (cs components, err error) {
	for index, chunk := range strings.Split(raw, ",") {
		if strings.Contains(chunk, "..") {
			c, err := parseRange(chunk)
			if err != nil {
				return cs, fmt.Errorf(`parsing range %d: %w`, index, err)
			}
			cs = append(cs, c)
			continue
		}

		c, err := parseValue(chunk)
		if err != nil {
			return cs, fmt.Errorf(`parsing value %d: %w`, index, err)
		}
		cs = append(cs, c)
	}

	return cs, err
}

// MarshalText implements the encoding.MarshalText interface for a component
// slice.
func (cs components) MarshalText() (text []byte, err error) {
	var parts [][]byte
	for _, c := range cs {
		b, _ := c.MarshalText()
		parts = append(parts, b)
	}
	return bytes.Join(parts, []byte(",")), nil
}

// String implements the fmt.Stringer interface.
func (cs components) String() string {
	b, _ := cs.MarshalText()
	return string(b)
}

// Values return the list of actual values from the various sub-components.
func (cs components) Values(max int) (values []int) {
	var seen = make(map[int]struct{})

	for _, c := range cs {
		for {
			if c.To == 0 {
				seen[c.From] = struct{}{}
			} else {
				for v := c.From; v <= c.To && v <= max; v++ {
					seen[v] = struct{}{}
				}
			}

			if c.Repeat == 0 {
				break
			}

			c.From += c.Repeat
			if c.To != 0 {
				c.To += c.Repeat
			}

			if c.From > max {
				break
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

// Next returns the next valid value for the components, based on the current
// value. The next value can be equal to the current value if it is valid. The
// returned value can be smaller than the current value as the values are
// considered modulo the maximum value.
func (cs components) Next(current, max int) (next int, diff int, ok bool) {
	values := cs.Values(max)
	if len(values) == 0 {
		return
	}

	// Get the first value that is greater or equal to the current value.
	var i int
	for i = 0; i < len(values) && values[i] < current; i++ {
	}

	var val = values[i%len(values)]

	return val, val - current, true
}

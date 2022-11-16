# zCalendar [![Build Status](https://travis-ci.org/synthesio/zcalendar.svg?branch=master)](https://travis-ci.org/synthesio/zcalendar) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/synthesio/zcalendar) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/synthesio/zcalendar/master/LICENSE.md)

A CalendarSpec (almost) compliant calendar event library.

## Why ?

Because cron-spec is notoriously hard to read, not intuitive, and not standardized.

## Calendar Event format

**Note** Shamelessly taken from [System's
documentation](https://www.freedesktop.org/software/systemd/man/systemd.time.html#Calendar%20Events),
with some changes to reflect what is or isn't impmemented as of the time of
writing.

Calendar events may be used to refer to one or more points in time in a single
expression.

`Thu,Fri 2012-*-1,5 11:12:13`

The above refers to 11:12:13 of the first or fifth day of any month of the year
2012, but only if that day is a Thursday or Friday.

The weekday specification is optional. If specified, it should consist of one
or more English language weekday names, either in the abbreviated (Wed) or
non-abbreviated (Wednesday) form (case does not matter), separated by commas.
Specifying two weekdays separated by ".." refers to a range of continuous
weekdays. "," and ".." may be combined freely.

Monday is considered the first day of the week and a range of weekdays must be
contained in the same week. For example, "every weekday except Thursday" cannot
be written `Fri..Wed`; it may be written `Mon..Wed,Fri..Sun` instead.

In the date and time specifications, any component may be specified as "*" in
which case any value will match. Alternatively, each component can be specified
as a list of values separated by commas. Values may be suffixed with "/" and a
repetition value, which indicates that the value itself and the value plus all
multiples of the repetition value are matched. Two values separated by ".." may
be used to indicate a range of values; ranges may also be followed with "/" and
a repetition value.

Either time or date specification may be omitted, in which case *-*-* and
00:00:00 is implied, respectively. If the year component is not specified, "*-"
is assumed. If the second component is not specified, ":00" is assumed.

A timezone specification may be added at the end of the expression. It can be
either `UTC` or a timezone as defined by the IANA (see
[here](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) for a
somewhat complete list).

Examples for valid timestamps and their normalized form:

```
  Sat,Thu,Mon..Wed,Sat..Sun → Mon..Thu,Sat,Sun *-*-* 00:00:00
      Mon,Sun 12-*-* 2,1:23 → Mon,Sun 2012-*-* 01,02:23:00
                    Wed *-1 → Wed *-*-01 00:00:00
           Wed..Wed,Wed *-1 → Wed *-*-01 00:00:00
                 Wed, 17:48 → Wed *-*-* 17:48:00
Wed..Sat,Tue 12-10-15 1:2:3 → Tue..Sat 2012-10-15 01:02:03
                *-*-7 0:0:0 → *-*-07 00:00:00
                      10-15 → *-10-15 00:00:00
        monday *-12-* 17:00 → Mon *-12-* 17:00:00
  Mon,Fri *-*-3,1,2 *:30:45 → Mon,Fri *-*-01,02,03 *:30:45
       12,14,13,12:20,10,30 → *-*-* 12,13,14:10,20,30:00
            12..14:10,20,30 → *-*-* 12..14:10,20,30:00
  mon,fri *-1/2-1,3 *:30:45 → Mon,Fri *-01/2-01,03 *:30:45
             03-05 08:05:40 → *-03-05 08:05:40
                   08:05:40 → *-*-* 08:05:40
                      05:40 → *-*-* 05:40:00
     Sat,Sun 12-05 08:05:40 → Sat,Sun *-12-05 08:05:40
           Sat,Sun 08:05:40 → Sat,Sun *-*-* 08:05:40
           2003-03-05 05:40 → 2003-03-05 05:40:00
             2003-02..04-05 → 2003-02..04-05 00:00:00
       2003-03-05 05:40 UTC → 2003-03-05 05:40:00 UTC
                 2003-03-05 → 2003-03-05 00:00:00
                      03-05 → *-03-05 00:00:00
                      *:2/3 → *-*-* *:02/3:00
```

## Usage

The main structure is `Expression`. It implements both `encoding.TextMarhaler`,
`encoding.TextUnmarshaler`, and `fmt.Stringer`, so it can be used in JSON
objects and stored in various ways.

There is also a `func Parse(raw string) (exp Expression, err error)` method to
parse a textual representation to an expression.

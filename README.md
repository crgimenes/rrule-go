# rrule-go

Go library for working with recurrence rules for calendar dates.

The rrule module offers a complete implementation of the recurrence rules documented in the [iCalendar
RFC](http://www.ietf.org/rfc/rfc2445.txt). It is a partial port of the rrule module from the excellent [python-dateutil](http://labix.org/python-dateutil/) library.

## Demo

### rrule.RRule

```go
package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/rrule-go"
)

func printTimeSlice(ts []time.Time) {
	for _, t := range ts {
		fmt.Println(t)
	}
}

func main() {
	// Daily, for 10 occurrences.
	r, err := rrule.NewRRule(rrule.ROption{
		Freq:    rrule.DAILY,
		Count:   10,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(r.String())
	printTimeSlice(r.Between(
		time.Date(1997, 9, 6, 0, 0, 0, 0, time.UTC),
		time.Date(1997, 9, 8, 0, 0, 0, 0, time.UTC),
		true,
	))

	// Every four years, the first Tuesday after a Monday in November,
	// 3 occurrences (U.S. Presidential Election day).
	r, err = rrule.NewRRule(rrule.ROption{
		Freq:       rrule.YEARLY,
		Interval:   4,
		Count:      3,
		Bymonth:    []int{11},
		Byweekday:  []rrule.Weekday{rrule.TU},
		Bymonthday: []int{2, 3, 4, 5, 6, 7, 8},
		Dtstart:    time.Date(1996, 11, 5, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(r.String())
	printTimeSlice(r.All())
}
```

### rrule.Set

```go
package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/rrule-go"
)

func printTimeSlice(ts []time.Time) {
	for _, t := range ts {
		fmt.Println(t)
	}
}

func main() {
	// Daily, for 7 occurrences.
	set := rrule.Set{}
	r, err := rrule.NewRRule(rrule.ROption{
		Freq:    rrule.DAILY,
		Count:   7,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		panic(err)
	}
	if err := set.RRule(r); err != nil {
		panic(err)
	}

	fmt.Println(set.String())
	printTimeSlice(set.All())

	// Weekly, for 4 weeks, plus one time on day 7, and not on day 16.
	set = rrule.Set{}
	r, err = rrule.NewRRule(rrule.ROption{
		Freq:    rrule.WEEKLY,
		Count:   4,
		Dtstart: time.Date(1997, 9, 2, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		panic(err)
	}
	if err := set.RRule(r); err != nil {
		panic(err)
	}
	set.RDate(time.Date(1997, 9, 7, 9, 0, 0, 0, time.UTC))
	set.ExDate(time.Date(1997, 9, 16, 9, 0, 0, 0, time.UTC))

	fmt.Println(set.String())
	printTimeSlice(set.All())
}
```

### rrule.StrToRRule

```go
package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/rrule-go"
)

func printTimeSlice(ts []time.Time) {
	for _, t := range ts {
		fmt.Println(t)
	}
}

func main() {
	r, err := rrule.StrToRRule("FREQ=DAILY;DTSTART=20060101T150405Z;COUNT=5")
	if err != nil {
		panic(err)
	}

	fmt.Println(r.OrigOptions.RRuleString())
	fmt.Println(r.OrigOptions.String())
	fmt.Println(r.String())
	printTimeSlice(r.All())
}
```

### rrule.StrToRRuleSet

```go
package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/rrule-go"
)

func printTimeSlice(ts []time.Time) {
	for _, t := range ts {
		fmt.Println(t)
	}
}

func main() {
	s, err := rrule.StrToRRuleSet("DTSTART:20060101T150405Z\nRRULE:FREQ=DAILY;COUNT=5\nEXDATE:20060102T150405Z")
	if err != nil {
		panic(err)
	}

	fmt.Println(s.String())
	printTimeSlice(s.All())
}
```

For more examples see [python-dateutil](http://labix.org/python-dateutil/) documentation.

## License

This project is licensed under the [MIT](LICENSE) license.

## Notice

This repository is a heavily modified hard fork of the original rrule-go library available at <https://github.com/teambition/rrule-go>.

When possible, consider using the original version of the library.

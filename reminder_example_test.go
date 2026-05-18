// 2017-2022, crgimenes. All rights reserved.

package rrule_test

import (
	"fmt"
	"time"

	"github.com/crgimenes/rrule-go"
)

func mustLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

func Example_reminderNextOccurrence() {
	loc := mustLocation("America/Sao_Paulo")
	r, err := rrule.StrToRRule("DTSTART;TZID=America/Sao_Paulo:20260518T090000\nRRULE:FREQ=DAILY;INTERVAL=2")
	if err != nil {
		panic(err)
	}

	now := time.Date(2026, 5, 19, 12, 0, 0, 0, loc)
	next := r.After(now, false)
	fmt.Println(next.Format(time.RFC3339))

	// Output:
	// 2026-05-20T09:00:00-03:00
}

func Example_reminderDueWindow() {
	loc := mustLocation("America/Sao_Paulo")
	r, err := rrule.StrToRRule("DTSTART;TZID=America/Sao_Paulo:20260518T090000\nRRULE:FREQ=DAILY;COUNT=4")
	if err != nil {
		panic(err)
	}

	lastChecked := time.Date(2026, 5, 18, 9, 0, 1, 0, loc)
	now := time.Date(2026, 5, 20, 9, 0, 0, 0, loc)
	for _, due := range r.Between(lastChecked, now, true) {
		fmt.Println(due.Format(time.RFC3339))
	}

	// Output:
	// 2026-05-19T09:00:00-03:00
	// 2026-05-20T09:00:00-03:00
}

func Example_reminderSkipAndSnooze() {
	loc := mustLocation("America/Sao_Paulo")
	set, err := rrule.StrToRRuleSet("DTSTART;TZID=America/Sao_Paulo:20260518T090000\nRRULE:FREQ=DAILY;COUNT=3")
	if err != nil {
		panic(err)
	}

	original := time.Date(2026, 5, 19, 9, 0, 0, 0, loc)
	snoozed := time.Date(2026, 5, 19, 10, 30, 0, 0, loc)
	set.ExDate(original)
	set.RDate(snoozed)

	windowStart := time.Date(2026, 5, 18, 0, 0, 0, 0, loc)
	windowEnd := time.Date(2026, 5, 20, 23, 59, 59, 0, loc)
	for _, due := range set.Between(windowStart, windowEnd, true) {
		fmt.Println(due.Format(time.RFC3339))
	}

	// Output:
	// 2026-05-18T09:00:00-03:00
	// 2026-05-19T10:30:00-03:00
	// 2026-05-20T09:00:00-03:00
}

func Example_reminderRoundTrip() {
	loc := mustLocation("America/Sao_Paulo")
	set, err := rrule.StrToRRuleSet("DTSTART;TZID=America/Sao_Paulo:20260518T090000\nRRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR;COUNT=3")
	if err != nil {
		panic(err)
	}
	set.ExDate(time.Date(2026, 5, 20, 9, 0, 0, 0, loc))
	set.RDate(time.Date(2026, 5, 21, 9, 0, 0, 0, loc))

	stored := set.String()
	loaded, err := rrule.StrToRRuleSet(stored)
	if err != nil {
		panic(err)
	}
	fmt.Println(loaded.String())

	// Output:
	// DTSTART;TZID=America/Sao_Paulo:20260518T090000
	// RRULE:FREQ=WEEKLY;COUNT=3;BYDAY=MO,WE,FR
	// RDATE;TZID=America/Sao_Paulo:20260521T090000
	// EXDATE;TZID=America/Sao_Paulo:20260520T090000
}

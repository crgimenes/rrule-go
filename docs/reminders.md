# Reminder Systems

Use `RRule` for the repeating schedule and `Set` when a reminder also needs skipped or extra occurrences. Store explicit `DTSTART` values with a real timezone; reminder behavior is calendar-based, not just duration-based.

## Core Model

- Store one RRULE string per repeating reminder.
- Store `DTSTART` with `TZID` for user-facing local reminders.
- Use `After(now, false)` to find the next future notification.
- Use `Between(lastChecked, now, true)` to recover due reminders after downtime.
- Use `EXDATE` for dismissed or skipped occurrences.
- Use `RDATE` for snoozed or manually added occurrences.
- Avoid `All()` for reminders without `COUNT` or `UNTIL`; open-ended reminders can generate hundreds of years of occurrences.

## Next Occurrence

```go
package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/rrule-go"
)

func main() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	r, err := rrule.StrToRRule("DTSTART;TZID=America/Sao_Paulo:20260518T090000\nRRULE:FREQ=DAILY;INTERVAL=2")
	if err != nil {
		panic(err)
	}

	now := time.Date(2026, 5, 19, 12, 0, 0, 0, loc)
	next := r.After(now, false)
	fmt.Println(next.Format(time.RFC3339))
}
```

Expected output:

```text
2026-05-20T09:00:00-03:00
```

## Due Window

Use a window query when a worker wakes up after downtime. Store the last successful check timestamp and ask for all occurrences that became due since then.

```go
package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/rrule-go"
)

func main() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	r, err := rrule.StrToRRule("DTSTART;TZID=America/Sao_Paulo:20260518T090000\nRRULE:FREQ=DAILY;COUNT=4")
	if err != nil {
		panic(err)
	}

	lastChecked := time.Date(2026, 5, 18, 9, 0, 1, 0, loc)
	now := time.Date(2026, 5, 20, 9, 0, 0, 0, loc)
	for _, due := range r.Between(lastChecked, now, true) {
		fmt.Println(due.Format(time.RFC3339))
	}
}
```

Expected output:

```text
2026-05-19T09:00:00-03:00
2026-05-20T09:00:00-03:00
```

## Skip And Snooze

Use `EXDATE` to suppress the original occurrence and `RDATE` to add the snoozed replacement.

```go
package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/rrule-go"
)

func main() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

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
}
```

Expected output:

```text
2026-05-18T09:00:00-03:00
2026-05-19T10:30:00-03:00
2026-05-20T09:00:00-03:00
```

## Persistence

Persist `set.String()` when a reminder has exclusions or extra dates. Persist `r.String()` or `option.String()` when it only has one rule.

```go
package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/rrule-go"
)

func main() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

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
}
```

Expected output:

```text
DTSTART;TZID=America/Sao_Paulo:20260518T090000
RRULE:FREQ=WEEKLY;COUNT=3;BYDAY=MO,WE,FR
RDATE;TZID=America/Sao_Paulo:20260521T090000
EXDATE;TZID=America/Sao_Paulo:20260520T090000
```

## Operational Notes

- Store timestamps truncated to seconds; this package intentionally evaluates occurrences at second precision.
- Store the user's timezone name, not only a UTC instant, when a reminder should remain at the same wall-clock time.
- Validate user-provided RRULE strings with `StrToRRule` or `StrToRRuleSet` before persisting them.
- Prefer `COUNT` or `UNTIL` for finite reminders.
- Keep generated notification delivery state outside the RRULE. RRULE describes the schedule; delivery attempts, acknowledgements, and notification ids are application state.

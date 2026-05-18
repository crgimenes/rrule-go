// Package rrule implements iCalendar recurrence rules for calendar and
// reminder systems.
//
// Rules are evaluated at second precision. Provide an explicit DTSTART with a
// real time.Location for user-facing reminders so recurrence generation follows
// the user's calendar timezone, including daylight saving transitions.
//
// Use RRule for a single recurrence rule. Use Set when a reminder needs one-off
// inclusions or exclusions, such as snoozing an occurrence with RDATE or
// dismissing a generated occurrence with EXDATE. For polling reminder workers,
// prefer After or Between over All for open-ended rules.
package rrule

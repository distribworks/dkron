package extcron

import (
	"time"
)

// SimpleSchedule represents a simple non recurring duration.
type SimpleSchedule struct {
	Date time.Time
}

// At just stores the given time for this schedule.
func At(date time.Time) SimpleSchedule {
	return SimpleSchedule{
		Date: date,
	}
}

// Next conforms to the Schedule interface but this kind of jobs
// doesn't need to be run more than once, so it doesn't return a new date but the existing one.
func (schedule SimpleSchedule) Next(t time.Time) time.Time {
	// If the date set is after the reference time return it.
	// If it's before, return a time in the past (01-01-0001)
	// so it never runs.
	if schedule.Date.After(t) {
		return schedule.Date
	}
	return time.Time{}
}

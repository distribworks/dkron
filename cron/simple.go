package cron

import (
	"time"
)

// SimpleDelaySchedule represents a simple non recurring duration.
type SimpleSchedule struct {
	Date time.Time
}

// Just store the given time for this schedule.
func At(date time.Time) SimpleSchedule {
	return SimpleSchedule{
		Date: date,
	}
}

// Next conforms to the Schedule interface but this kind of jobs
// doesn't need to be run more than once, so it doesn't return a new date but the existing one.
func (schedule SimpleSchedule) Next(t time.Time) time.Time {
	// If the date set is after the reference time return it
	// if it's before, return a virtually infinite sleep date
	// so do nothing.
	if schedule.Date.After(t) {
		return schedule.Date
	}
	return t.AddDate(10, 0, 0)
}

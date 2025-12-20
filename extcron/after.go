package extcron

import (
	"time"
)

// AfterSchedule represents a schedule that runs once at a specific time,
// with a grace period during which it can still run immediately if missed.
type AfterSchedule struct {
	Date        time.Time
	GracePeriod time.Duration
}

// After creates an AfterSchedule with the given date and grace period.
func After(date time.Time, gracePeriod time.Duration) AfterSchedule {
	return AfterSchedule{
		Date:        date,
		GracePeriod: gracePeriod,
	}
}

// Next conforms to the Schedule interface.
// It returns:
// - The scheduled date if current time is before the scheduled date
// - Current time (immediate execution) if current time is within grace period after the scheduled date
// - Zero time (never runs) if current time is beyond the grace period
func (schedule AfterSchedule) Next(t time.Time) time.Time {
	// If the date is after the reference time, return it
	if schedule.Date.After(t) {
		return schedule.Date
	}

	// If we're within the grace period (including the exact end moment), run immediately
	gracePeriodEnd := schedule.Date.Add(schedule.GracePeriod)
	if !t.After(gracePeriodEnd) {
		return t
	}

	// Beyond grace period, never run
	return time.Time{}
}

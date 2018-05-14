package cron

import (
	"time"
)

// Schedule wrapper with time zone awarness
type timezoneAwareSchedule struct {
	underlyingSchedule Schedule
	targetedTimezone  *time.Location
}

func (schedule *timezoneAwareSchedule) Next(t time.Time) time.Time {
	return schedule.underlyingSchedule.Next(t.In(schedule.targetedTimezone))
}

// Wrap a schedule inside a timezoneAwareSchedule.  timezone string
// must be a supported timezone representation by the OS otherwise
// will return an error.
func wrapSchedulerInTimezone(schedule Schedule, timezone string) (*timezoneAwareSchedule, error) {
	if location, err := time.LoadLocation(timezone); err != nil {
		return nil, err
	} else {
		return &timezoneAwareSchedule{schedule, location}, nil
	}
}

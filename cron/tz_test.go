package cron

import (
	"testing"
	"time"
)

func TestTimezoneAwareSchedule(t *testing.T) {
	timezone := "Europe/London"
	schedule, _ := Parse("0 13 15 * * *")
	now, _ := time.Parse(time.RFC3339, "2017-09-09T22:08:41+04:00")
	expt, _ := time.Parse(time.RFC3339, "2017-09-10T15:13:00+01:00")
	if tzSchedule, err := wrapSchedulerInTimezone(schedule, timezone); err != nil {
		t.Error(err)
	} else if act := tzSchedule.Next(now); !act.Equal(expt) {
		t.Error(act.Format(time.RFC3339))
	}
}

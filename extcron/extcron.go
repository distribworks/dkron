package extcron

import (
	"strings"

	"github.com/distribworks/dkron/v2/cron"
)

// ExtCron is an extended cron scheduler, which wraps, adapts and decorates
// the robfig/cron v3 library
type ExtCron struct {
	cron *cron.Cron
}

// New creates a new ExtCron instance
func New() ExtCron {
	return ExtCron{
		cron: cron.New(cron.WithParser(NewParser())),
	}
}

// Start starts the scheduler
func (e ExtCron) Start() {
	e.cron.Start()
}

// Stop stops the scheduler
func (e ExtCron) Stop() {
	e.cron.Stop()
}

// Schedule adds a Job to the Cron to be run on the given schedule.
func (e ExtCron) Schedule(schedule cron.Schedule, cmd cron.Job) {
	_ = e.cron.Schedule(schedule, cmd)
}

// Entries returns a snapshot of the cron entries.
func (e ExtCron) Entries() []*cron.Entry {
	entries := e.cron.Entries()
	entryPtrs := make([]*cron.Entry, len(entries))
	for i, entry := range entries {
		entryPtrs[i] = &entry
	}
	return entryPtrs
}

// AddFunc adds a func to the Cron to be run on the given schedule.
func (e ExtCron) AddFunc(spec string, cmd func()) error {
	_, err := e.cron.AddFunc(spec, cmd)
	return err
}

// AddJob adds a Job to the Cron to be run on the given schedule.
func (e ExtCron) AddJob(spec string, cmd cron.Job) error {
	_, err := e.cron.AddJob(spec, cmd)
	return err
}

// AddTimezoneSensitiveJob adds a Job to the Cron to be run on the given
// schedule in the given timezone. Timezone is ignored for descriptor schedules
// and schedules that have a timezone set in the spec.
func (e ExtCron) AddTimezoneSensitiveJob(spec, timezone string, cmd cron.Job) error {
	if strings.HasPrefix(spec, "@") || strings.HasPrefix(spec, "TZ=") || strings.HasPrefix(spec, "CRON_TZ=") {
		return e.AddJob(spec, cmd)
	}

	return e.AddJob("CRON_TZ="+timezone+" "+spec, cmd)
}

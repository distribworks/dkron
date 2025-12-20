package extcron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAfterScheduleEndToEnd tests the complete @after schedule flow
// from parsing to execution time calculation
func TestAfterScheduleEndToEnd(t *testing.T) {
	tests := []struct {
		name            string
		schedule        string
		currentTime     time.Time
		expectNextAt    time.Time
		expectImmediate bool
		expectNever     bool
		description     string
	}{
		{
			name:         "future job should run at scheduled time",
			schedule:     "@after 2025-01-01T00:00:00Z <PT2H",
			currentTime:  time.Date(2024, 12, 31, 23, 0, 0, 0, time.UTC),
			expectNextAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			description:  "Job created before scheduled time should run at that time",
		},
		{
			name:            "past job within grace period should run immediately",
			schedule:        "@after 2025-01-01T00:00:00Z <PT2H",
			currentTime:     time.Date(2025, 1, 1, 1, 0, 0, 0, time.UTC),
			expectImmediate: true,
			description:     "Job created 1 hour after scheduled time (within 2h grace) should run immediately",
		},
		{
			name:         "past job at end of grace period should run immediately",
			schedule:     "@after 2025-01-01T00:00:00Z <PT2H",
			currentTime:  time.Date(2025, 1, 1, 1, 59, 59, 0, time.UTC),
			expectImmediate: true,
			description:  "Job created just before end of grace period should run immediately",
		},
		{
			name:        "past job beyond grace period should never run",
			schedule:    "@after 2025-01-01T00:00:00Z <PT2H",
			currentTime: time.Date(2025, 1, 1, 2, 0, 1, 0, time.UTC),
			expectNever: true,
			description: "Job created 1 second after grace period should never run",
		},
		{
			name:         "exactly at scheduled time should run immediately",
			schedule:     "@after 2025-01-01T00:00:00Z <PT2H",
			currentTime:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expectImmediate: true,
			description:  "Job created exactly at scheduled time should run immediately",
		},
		{
			name:            "small grace period example",
			schedule:        "@after 2025-01-01T12:00:00Z <PT5M",
			currentTime:     time.Date(2025, 1, 1, 12, 3, 0, 0, time.UTC),
			expectImmediate: true,
			description:     "Job with 5 minute grace period created 3 minutes late should run immediately",
		},
		{
			name:        "small grace period expired",
			schedule:    "@after 2025-01-01T12:00:00Z <PT5M",
			currentTime: time.Date(2025, 1, 1, 12, 6, 0, 0, time.UTC),
			expectNever: true,
			description: "Job with 5 minute grace period created 6 minutes late should never run",
		},
		{
			name:            "large grace period example",
			schedule:        "@after 2025-01-01T00:00:00Z <P30D",
			currentTime:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			expectImmediate: true,
			description:     "Job with 30 day grace period created 15 days late should run immediately",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the schedule
			schedule, err := Parse(tt.schedule)
			require.NoError(t, err, "Failed to parse schedule: %s", tt.schedule)

			// Calculate the next execution time
			nextRun := schedule.Next(tt.currentTime)

			if tt.expectImmediate {
				// Should run immediately (next run == current time)
				assert.Equal(t, tt.currentTime, nextRun, tt.description)
			} else if tt.expectNever {
				// Should never run (zero time)
				assert.True(t, nextRun.IsZero(), "Expected zero time (never run), got: %v. %s", nextRun, tt.description)
			} else {
				// Should run at the expected time
				assert.Equal(t, tt.expectNextAt, nextRun, tt.description)
			}
		})
	}
}

// TestAfterScheduleUseCase demonstrates the primary use case from the issue:
// Jobs that are created slightly after their scheduled time due to network latency
func TestAfterScheduleUseCase(t *testing.T) {
	// Scenario: User wants to create a job that runs at a specific time,
	// but due to network latency, the job is created a few seconds after that time.
	
	// Job should run at 2025-01-01 12:00:00
	scheduledTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	
	// But job is created 30 seconds late due to network latency
	creationTime := scheduledTime.Add(30 * time.Second)
	
	// With a 5-minute grace period, the job should still run immediately
	schedule := "@after " + scheduledTime.Format(time.RFC3339) + " <PT5M"
	
	sched, err := Parse(schedule)
	require.NoError(t, err)
	
	nextRun := sched.Next(creationTime)
	
	// The job should run immediately (at creation time)
	assert.Equal(t, creationTime, nextRun,
		"Job created 30s late with 5min grace period should run immediately")
	
	// If the job is created much later (10 minutes), it should not run
	lateCreationTime := scheduledTime.Add(10 * time.Minute)
	nextRun = sched.Next(lateCreationTime)
	assert.True(t, nextRun.IsZero(),
		"Job created 10 minutes late with 5min grace period should never run")
}

// TestAfterScheduleVsAtSchedule compares @after with @at behavior
func TestAfterScheduleVsAtSchedule(t *testing.T) {
	scheduledTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	pastTime := scheduledTime.Add(-30 * time.Minute) // 30 minutes before scheduled time
	futureTime := scheduledTime.Add(30 * time.Minute) // 30 minutes after scheduled time
	
	// Parse @at schedule
	atSchedule, err := Parse("@at " + scheduledTime.Format(time.RFC3339))
	require.NoError(t, err)
	
	// Parse @after schedule with 1 hour grace period
	afterSchedule, err := Parse("@after " + scheduledTime.Format(time.RFC3339) + " <PT1H")
	require.NoError(t, err)
	
	// Before scheduled time: both should return scheduled time
	atNext := atSchedule.Next(pastTime)
	afterNext := afterSchedule.Next(pastTime)
	assert.Equal(t, scheduledTime, atNext, "@at should return scheduled time when before")
	assert.Equal(t, scheduledTime, afterNext, "@after should return scheduled time when before")
	
	// After scheduled time (within grace period for @after):
	// @at returns zero time (never runs)
	// @after returns current time (runs immediately)
	atNext = atSchedule.Next(futureTime)
	afterNext = afterSchedule.Next(futureTime)
	assert.True(t, atNext.IsZero(), "@at should return zero time (never run) when after")
	assert.Equal(t, futureTime, afterNext, "@after should return current time (run immediately) when within grace")
}

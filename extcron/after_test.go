package extcron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAfterScheduleNext(t *testing.T) {
	tests := []struct {
		name        string
		scheduleAt  string
		gracePeriod time.Duration
		currentTime string
		expected    string
	}{
		{
			name:        "before scheduled time - should return scheduled time",
			scheduleAt:  "2020-01-01T00:00:00Z",
			gracePeriod: 2 * time.Hour,
			currentTime: "2019-12-31T23:00:00Z",
			expected:    "2020-01-01T00:00:00Z",
		},
		{
			name:        "within grace period - should run immediately",
			scheduleAt:  "2020-01-01T00:00:00Z",
			gracePeriod: 2 * time.Hour,
			currentTime: "2020-01-01T01:00:00Z",
			expected:    "2020-01-01T01:00:00Z", // Returns current time (immediate)
		},
		{
			name:        "at end of grace period - should run immediately",
			scheduleAt:  "2020-01-01T00:00:00Z",
			gracePeriod: 2 * time.Hour,
			currentTime: "2020-01-01T01:59:59Z",
			expected:    "2020-01-01T01:59:59Z", // Returns current time (immediate)
		},
		{
			name:        "exactly at end of grace period - should run immediately",
			scheduleAt:  "2020-01-01T00:00:00Z",
			gracePeriod: 2 * time.Hour,
			currentTime: "2020-01-01T02:00:00Z",
			expected:    "2020-01-01T02:00:00Z", // Returns current time (immediate)
		},
		{
			name:        "after grace period - should never run",
			scheduleAt:  "2020-01-01T00:00:00Z",
			gracePeriod: 2 * time.Hour,
			currentTime: "2020-01-01T02:00:01Z",
			expected:    "0001-01-01T00:00:00Z", // Zero time
		},
		{
			name:        "exactly at scheduled time - should run immediately",
			scheduleAt:  "2020-01-01T00:00:00Z",
			gracePeriod: 2 * time.Hour,
			currentTime: "2020-01-01T00:00:00Z",
			expected:    "2020-01-01T00:00:00Z", // Returns current time (immediate)
		},
		{
			name:        "small grace period",
			scheduleAt:  "2020-01-01T12:00:00Z",
			gracePeriod: 5 * time.Minute,
			currentTime: "2020-01-01T12:03:00Z",
			expected:    "2020-01-01T12:03:00Z", // Returns current time (immediate)
		},
		{
			name:        "past small grace period",
			scheduleAt:  "2020-01-01T12:00:00Z",
			gracePeriod: 5 * time.Minute,
			currentTime: "2020-01-01T12:06:00Z",
			expected:    "0001-01-01T00:00:00Z", // Zero time
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduleAt, _ := time.Parse(time.RFC3339, tt.scheduleAt)
			currentTime, _ := time.Parse(time.RFC3339, tt.currentTime)
			expected, _ := time.Parse(time.RFC3339, tt.expected)

			schedule := After(scheduleAt, tt.gracePeriod)
			actual := schedule.Next(currentTime)

			assert.Equal(t, expected, actual, "Expected %v, got %v", expected, actual)
		})
	}
}

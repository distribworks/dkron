package extcron

import (
	"testing"
	"time"
)

func TestSimpleNext(t *testing.T) {
	tests := []struct {
		time     string
		date     string
		expected string
	}{
		// Simple cases
		{"2012-07-09T14:45:00Z", "2012-07-09T15:00:00Z", "2012-07-09T15:00:00Z"},
		{"2012-07-09T14:45:00Z", "2012-07-05T13:00:00Z", "0001-01-01T00:00:00Z"},
	}

	for _, c := range tests {
		now, _ := time.Parse(time.RFC3339, c.time)
		date, _ := time.Parse(time.RFC3339, c.date)
		actual := At(date).Next(now)
		expected, _ := time.Parse(time.RFC3339, c.expected)

		if actual != expected {
			t.Errorf("%s, \"%s\": (expected) %v != %v (actual)", c.time, c.date, expected, actual)
		}
	}
}

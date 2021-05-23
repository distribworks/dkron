package extcron

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomSpecSchedule(t *testing.T) {
	entries := []struct {
		expr     string
		expected cron.Schedule
	}{
		{"@at 2018-01-02T15:04:00Z", SimpleSchedule{time.Date(2018, time.January, 2, 15, 4, 0, 0, time.UTC)}},
		{"@at 2019-02-04T09:20:00+06:00", SimpleSchedule{time.Date(2019, time.February, 4, 9, 20, 0, 0, time.FixedZone("", 21600))}},
		{"@manually", SimpleSchedule{time.Time{}}},
		{"@minutely", &cron.SpecSchedule{
			Second:   0b1,
			Minute:   (1 << 63) + 0b111111111111111111111111111111111111111111111111111111111111,
			Hour:     (1 << 63) + 0b111111111111111111111111,
			Dom:      (1 << 63) + 0b11111111111111111111111111111110,
			Month:    (1 << 63) + 0b1111111111110,
			Dow:      (1 << 63) + 0b1111111,
			Location: time.Local}},
	}

	for _, c := range entries {
		actual, err := Parse(c.expr)
		require.NoError(t, err)
		assert.Equal(t, c.expected, actual, "%s => (expected) %v != %v (actual)", c.expr, c.expected, actual)
	}
}

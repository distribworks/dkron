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

func TestAfterScheduleParsing(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		wantErr  bool
		expected AfterSchedule
	}{
		{
			name:    "valid @after with 2 hour grace period",
			expr:    "@after 2020-01-01T00:00:00Z <PT2H",
			wantErr: false,
			expected: AfterSchedule{
				Date:        time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				GracePeriod: 2 * time.Hour,
			},
		},
		{
			name:    "valid @after with 1 day grace period",
			expr:    "@after 2020-01-01T12:00:00Z <P1D",
			wantErr: false,
			expected: AfterSchedule{
				Date:        time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC),
				GracePeriod: 24 * time.Hour,
			},
		},
		{
			name:    "valid @after with complex duration",
			expr:    "@after 2020-01-01T00:00:00Z <P1DT2H30M",
			wantErr: false,
			expected: AfterSchedule{
				Date:        time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				GracePeriod: 26*time.Hour + 30*time.Minute,
			},
		},
		{
			name:    "missing grace period prefix",
			expr:    "@after 2020-01-01T00:00:00Z PT2H",
			wantErr: true,
		},
		{
			name:    "invalid datetime format",
			expr:    "@after 2020-01-01 <PT2H",
			wantErr: true,
		},
		{
			name:    "invalid duration format",
			expr:    "@after 2020-01-01T00:00:00Z <invalid",
			wantErr: true,
		},
		{
			name:    "missing duration",
			expr:    "@after 2020-01-01T00:00:00Z",
			wantErr: true,
		},
		{
			name:    "too many parts",
			expr:    "@after 2020-01-01T00:00:00Z <PT2H extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Parse(tt.expr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				afterSchedule, ok := actual.(AfterSchedule)
				require.True(t, ok, "Expected AfterSchedule type")
				assert.Equal(t, tt.expected.Date, afterSchedule.Date)
				assert.Equal(t, tt.expected.GracePeriod, afterSchedule.GracePeriod)
			}
		})
	}
}

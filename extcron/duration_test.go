package extcron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseISO8601Duration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Duration
		wantErr  bool
	}{
		{
			name:     "2 hours",
			input:    "PT2H",
			expected: 2 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "2 hours without T",
			input:    "P2H",
			expected: 0,
			wantErr:  true, // Without T, hours are not recognized
		},
		{
			name:     "1 day",
			input:    "P1D",
			expected: 24 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "1 day and 2 hours",
			input:    "P1DT2H",
			expected: 26 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "30 minutes",
			input:    "PT30M",
			expected: 30 * time.Minute,
			wantErr:  false,
		},
		{
			name:     "1 hour 30 minutes",
			input:    "PT1H30M",
			expected: 90 * time.Minute,
			wantErr:  false,
		},
		{
			name:     "45 seconds",
			input:    "PT45S",
			expected: 45 * time.Second,
			wantErr:  false,
		},
		{
			name:     "1 month (30 days approximation)",
			input:    "P1M",
			expected: 30 * 24 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "1 year (365 days approximation)",
			input:    "P1Y",
			expected: 365 * 24 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "complex duration",
			input:    "P1Y2M3DT4H5M6S",
			expected: (365*24 + 60*24 + 3*24 + 4) * time.Hour + 5*time.Minute + 6*time.Second,
			wantErr:  false,
		},
		{
			name:     "invalid format",
			input:    "invalid",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "zero duration",
			input:    "P",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseISO8601Duration(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

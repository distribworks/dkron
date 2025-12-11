package ntime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullableTime_SetAndGet(t *testing.T) {
	var nt NullableTime
	assert.False(t, nt.HasValue())

	now := time.Now()
	nt.Set(now)
	assert.True(t, nt.HasValue())
	assert.Equal(t, now, nt.Get())
}

func TestNullableTime_Unset(t *testing.T) {
	var nt NullableTime
	nt.Set(time.Now())
	assert.True(t, nt.HasValue())

	nt.Unset()
	assert.False(t, nt.HasValue())
}

func TestNullableTime_GetPanicsWhenNoValue(t *testing.T) {
	var nt NullableTime
	assert.Panics(t, func() {
		nt.Get()
	})
}

func TestNullableTime_After(t *testing.T) {
	var nt1, nt2 NullableTime

	// nil after nil -> false
	assert.False(t, nt1.After(nt2))

	// nil after value -> false
	nt2.Set(time.Now())
	assert.False(t, nt1.After(nt2))

	// value after nil -> true
	nt1.Set(time.Now())
	nt2.Unset()
	assert.True(t, nt1.After(nt2))

	// value after earlier value -> true
	earlier := time.Now().Add(-time.Hour)
	later := time.Now()
	nt1.Set(later)
	nt2.Set(earlier)
	assert.True(t, nt1.After(nt2))

	// value after later value -> false
	nt1.Set(earlier)
	nt2.Set(later)
	assert.False(t, nt1.After(nt2))
}

func TestNullableTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() NullableTime
		expected string
	}{
		{
			name: "null when no value",
			setup: func() NullableTime {
				return NullableTime{}
			},
			expected: "null",
		},
		{
			name: "RFC3339 when has value",
			setup: func() NullableTime {
				var nt NullableTime
				nt.Set(time.Date(2025, 12, 19, 9, 9, 0, 0, time.UTC))
				return nt
			},
			expected: `"2025-12-19T09:09:00Z"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nt := tt.setup()
			data, err := json.Marshal(&nt)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestNullableTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectValue bool
		expectTime  time.Time
		expectError bool
	}{
		{
			name:        "null value",
			input:       "null",
			expectValue: false,
		},
		{
			name:        "RFC3339 format with Z timezone",
			input:       `"2025-12-19T09:09:00Z"`,
			expectValue: true,
			expectTime:  time.Date(2025, 12, 19, 9, 9, 0, 0, time.UTC),
		},
		{
			name:        "RFC3339 format with offset",
			input:       `"2025-12-19T09:09:00+01:00"`,
			expectValue: true,
			expectTime:  time.Date(2025, 12, 19, 9, 9, 0, 0, time.FixedZone("", 3600)),
		},
		{
			name:        "datetime-local format from HTML input",
			input:       `"2025-12-19T09:09"`,
			expectValue: true,
			expectTime:  time.Date(2025, 12, 19, 9, 9, 0, 0, time.UTC),
		},
		{
			name:        "invalid format",
			input:       `"not-a-date"`,
			expectError: true,
		},
		{
			name:        "invalid JSON",
			input:       `{invalid}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nt NullableTime
			err := json.Unmarshal([]byte(tt.input), &nt)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectValue, nt.HasValue())

			if tt.expectValue {
				assert.True(t, tt.expectTime.Equal(nt.Get()),
					"expected %v, got %v", tt.expectTime, nt.Get())
			}
		})
	}
}

func TestNullableTime_UnmarshalJSON_InStruct(t *testing.T) {
	type Job struct {
		Name     string       `json:"name"`
		StartsAt NullableTime `json:"starts_at"`
	}

	tests := []struct {
		name        string
		input       string
		expectValue bool
		expectTime  time.Time
	}{
		{
			name:        "datetime-local format in struct",
			input:       `{"name":"test","starts_at":"2025-12-19T09:09"}`,
			expectValue: true,
			expectTime:  time.Date(2025, 12, 19, 9, 9, 0, 0, time.UTC),
		},
		{
			name:        "RFC3339 format in struct",
			input:       `{"name":"test","starts_at":"2025-12-19T09:09:00Z"}`,
			expectValue: true,
			expectTime:  time.Date(2025, 12, 19, 9, 9, 0, 0, time.UTC),
		},
		{
			name:        "null in struct",
			input:       `{"name":"test","starts_at":null}`,
			expectValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var job Job
			err := json.Unmarshal([]byte(tt.input), &job)
			require.NoError(t, err)
			assert.Equal(t, "test", job.Name)
			assert.Equal(t, tt.expectValue, job.StartsAt.HasValue())

			if tt.expectValue {
				assert.True(t, tt.expectTime.Equal(job.StartsAt.Get()),
					"expected %v, got %v", tt.expectTime, job.StartsAt.Get())
			}
		})
	}
}

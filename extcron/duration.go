package extcron

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var iso8601DurationRegex = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?)?$`)

// ParseISO8601Duration parses an ISO8601 duration string (e.g., "P2H", "PT2H", "P1DT2H30M")
// into a time.Duration. Note that for month and year durations, we use approximations:
// 1 year = 365 days, 1 month = 30 days
func ParseISO8601Duration(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("empty duration string")
	}

	matches := iso8601DurationRegex.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid ISO8601 duration format: %s", s)
	}

	var duration time.Duration

	// Years (approximate: 365 days)
	if matches[1] != "" {
		years, _ := strconv.Atoi(matches[1])
		duration += time.Duration(years) * 365 * 24 * time.Hour
	}

	// Months (approximate: 30 days)
	if matches[2] != "" {
		months, _ := strconv.Atoi(matches[2])
		duration += time.Duration(months) * 30 * 24 * time.Hour
	}

	// Days
	if matches[3] != "" {
		days, _ := strconv.Atoi(matches[3])
		duration += time.Duration(days) * 24 * time.Hour
	}

	// Hours
	if matches[4] != "" {
		hours, _ := strconv.Atoi(matches[4])
		duration += time.Duration(hours) * time.Hour
	}

	// Minutes
	if matches[5] != "" {
		minutes, _ := strconv.Atoi(matches[5])
		duration += time.Duration(minutes) * time.Minute
	}

	// Seconds (can be fractional)
	if matches[6] != "" {
		seconds, _ := strconv.ParseFloat(matches[6], 64)
		duration += time.Duration(seconds * float64(time.Second))
	}

	if duration == 0 {
		return 0, errors.New("duration must be greater than zero")
	}

	return duration, nil
}

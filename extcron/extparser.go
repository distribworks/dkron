package extcron

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

// ExtParser is a parser extending robfig/cron v3 standard parser with
// several additional descriptors
type ExtParser struct {
	parser cron.Parser
}

// NewParser creates an ExtParser instance
func NewParser() cron.ScheduleParser {
	return ExtParser{cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)}
}

// Parse parses a cron schedule specification. It accepts the cron spec with
// mandatory seconds parameter, descriptors and the custom descriptors
// "@at <date>", "@manually" and "@minutely".
func (p ExtParser) Parse(spec string) (cron.Schedule, error) {
	switch spec {
	case "@manually":
		return At(time.Time{}), nil

	case "@minutely":
		spec = "0 * * * * *"
	}

	const at = "@at "
	if strings.HasPrefix(spec, at) {
		date, err := time.Parse(time.RFC3339, spec[len(at):])
		if err != nil {
			return nil, fmt.Errorf("failed to parse date %s: %s", spec, err)
		}
		return At(date), nil
	}

	// It's not a dkron specific spec: Let the regular cron schedule parser have it
	return p.parser.Parse(spec)
}

var standaloneParser = NewParser()

// Parse parses a cron schedule. This is a convenience function to not have
// to instantiate a parser with NewParser for every call.
func Parse(spec string) (cron.Schedule, error) {
	return standaloneParser.Parse(spec)
}

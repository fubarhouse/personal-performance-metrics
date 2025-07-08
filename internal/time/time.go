package time

import (
	"errors"
	"strconv"
	"time"
)

// ParseTimeString will parse input time as either unix timestamp format or one of a set of formats.
func ParseTimeString(s string) (time.Time, bool, error) {
	if s == "" {
		return time.Now(), true, nil
	}

	// Try parsing as Unix timestamp (digits only)
	if unix, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(unix, 0), true, nil
	}

	// Try parsing as DD/MM/YYYY date
	layouts := []string{
		"02/01/2006",          // DD/MM/YYYY
		"02/01/2006 15:04:05", // DD/MM/YYYY HH:MM:SS
		time.RFC3339,          // ISO-8601 / RFC-3339 (universal, includes timezone)
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, true, nil
		}
	}

	return time.Now(), false, errors.New("failed to parse time")
}

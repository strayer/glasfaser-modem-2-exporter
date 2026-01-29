package status

import (
	"math"
	"strconv"
	"strings"
)

// ParseFloat parses a float value from the ONT status.
// Returns NaN for empty strings, "--", or "N/A".
func ParseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "--" || s == "N/A" {
		return math.NaN()
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return math.NaN()
	}

	return f
}

// ParseUint parses an unsigned integer from the ONT status.
// Returns 0 for empty, invalid, or negative values.
func ParseUint(s string) uint64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}

	return u
}

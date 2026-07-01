package main

import (
	"fmt"
	"strconv"
	"time"
	"unicode"
)

const invalidDurationMessage = "expected format like 10s, 5m, 1h30m, or 1h30m10s"

// ParseDuration parses countdown durations containing h, m and s units.
// Units must appear at most once and in descending order: h, then m, then s.
func ParseDuration(input string) (time.Duration, error) {
	if input == "" {
		return 0, fmt.Errorf("empty duration: %s", invalidDurationMessage)
	}

	unitRank := map[rune]int{'h': 3, 'm': 2, 's': 1}
	lastRank := 4
	seen := map[rune]bool{}
	var total time.Duration

	for i := 0; i < len(input); {
		start := i
		for i < len(input) && unicode.IsDigit(rune(input[i])) {
			i++
		}

		if start == i {
			return 0, fmt.Errorf("invalid duration %q: %s", input, invalidDurationMessage)
		}
		if i >= len(input) {
			return 0, fmt.Errorf("invalid duration %q: missing unit", input)
		}

		unit := rune(input[i])
		rank, ok := unitRank[unit]
		if !ok {
			return 0, fmt.Errorf("invalid duration %q: unsupported unit %q", input, unit)
		}
		if seen[unit] {
			return 0, fmt.Errorf("invalid duration %q: unit %q repeated", input, unit)
		}
		if rank >= lastRank {
			return 0, fmt.Errorf("invalid duration %q: units must be ordered as h, m, s", input)
		}

		value, err := strconv.Atoi(input[start:i])
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q: %w", input, err)
		}
		if value <= 0 {
			return 0, fmt.Errorf("invalid duration %q: values must be positive", input)
		}

		switch unit {
		case 'h':
			total += time.Duration(value) * time.Hour
		case 'm':
			total += time.Duration(value) * time.Minute
		case 's':
			total += time.Duration(value) * time.Second
		}

		seen[unit] = true
		lastRank = rank
		i++
	}

	return total, nil
}

package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseUntil interprets a wall-clock time ("HH:MM" or "HH:MM:SS", 24h) as a
// target later today and returns the duration from now until then. It errors
// if the time has already passed today (we don't assume tomorrow).
func parseUntil(s string, now time.Time) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 && len(parts) != 3 {
		return 0, fmt.Errorf("invalid time %q: expected HH:MM or HH:MM:SS (24h)", s)
	}

	hour, err := parseTimeField(parts[0], 23)
	if err != nil {
		return 0, fmt.Errorf("invalid time %q: %w", s, err)
	}
	minute, err := parseTimeField(parts[1], 59)
	if err != nil {
		return 0, fmt.Errorf("invalid time %q: %w", s, err)
	}
	second := 0
	if len(parts) == 3 {
		second, err = parseTimeField(parts[2], 59)
		if err != nil {
			return 0, fmt.Errorf("invalid time %q: %w", s, err)
		}
	}

	target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())
	d := target.Sub(now)
	if d <= 0 {
		return 0, fmt.Errorf("time %q is not in the future today", s)
	}
	return d, nil
}

// parseTimeField parses a zero-padded or bare integer field within [0, max].
func parseTimeField(s string, max int) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%q is not a number", s)
	}
	if n < 0 || n > max {
		return 0, fmt.Errorf("%d out of range (0-%d)", n, max)
	}
	return n, nil
}

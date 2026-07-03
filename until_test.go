package main

import (
	"testing"
	"time"
)

func TestParseUntil(t *testing.T) {
	now := time.Date(2026, 7, 2, 14, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		in   string
		want time.Duration
	}{
		{"HH:MM later today", "15:00", time.Hour},
		{"HH:MM:SS later today", "14:30:30", 30*time.Minute + 30*time.Second},
		{"bare hour field", "14:01", time.Minute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseUntil(tt.in, now)
			if err != nil {
				t.Fatalf("parseUntil(%q) error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("parseUntil(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseUntilInvalid(t *testing.T) {
	now := time.Date(2026, 7, 2, 14, 0, 0, 0, time.UTC)

	tests := []string{
		"14:00",       // exactly now: not in the future
		"13:00",       // already passed today
		"25:00",       // hour out of range
		"12:60",       // minute out of range
		"abc",         // not a time
		"12",          // missing minute
		"12:30:40:50", // too many fields
	}
	for _, in := range tests {
		t.Run(in, func(t *testing.T) {
			if _, err := parseUntil(in, now); err == nil {
				t.Fatalf("parseUntil(%q) expected error, got nil", in)
			}
		})
	}
}

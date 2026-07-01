package main

import (
	"testing"
	"time"
)

func TestFormatRemaining(t *testing.T) {
	tests := []struct {
		name string
		in   time.Duration
		want string
	}{
		{"seconds", 10 * time.Second, "00:10"},
		{"minutes", 5 * time.Minute, "05:00"},
		{"hour", time.Hour, "01:00:00"},
		{"hour minute second", time.Hour + 30*time.Minute + 10*time.Second, "01:30:10"},
		{"negative", -5 * time.Second, "00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatRemaining(tt.in); got != tt.want {
				t.Fatalf("FormatRemaining(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

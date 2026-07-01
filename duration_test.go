package main

import (
	"testing"
	"time"
)

func TestParseDurationValid(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
	}{
		{"10s", 10 * time.Second},
		{"5m", 5 * time.Minute},
		{"1h", time.Hour},
		{"1h30m", time.Hour + 30*time.Minute},
		{"1h30m10s", time.Hour + 30*time.Minute + 10*time.Second},
		{"30m10s", 30*time.Minute + 10*time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if err != nil {
				t.Fatalf("ParseDuration(%q) returned error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("ParseDuration(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseDurationInvalid(t *testing.T) {
	tests := []string{
		"",
		"abc",
		"10x",
		"1m2h",
		"1h1h",
		"1.5h",
		"-10s",
		"0s",
		"1h0m",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			if _, err := ParseDuration(input); err == nil {
				t.Fatalf("ParseDuration(%q) expected error, got nil", input)
			}
		})
	}
}

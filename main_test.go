package main

import (
	"testing"
	"time"
)

func TestParseArgsDurationAndTitle(t *testing.T) {
	// Flag can appear after the positional duration.
	got, err := parseArgs([]string{"25m", "--title", "Deep work"})
	if err != nil {
		t.Fatalf("parseArgs error: %v", err)
	}
	if got.duration != 25*time.Minute {
		t.Errorf("duration = %v, want 25m", got.duration)
	}
	if got.title != "Deep work" {
		t.Errorf("title = %q, want %q", got.title, "Deep work")
	}
}

func TestParseArgsDefaultTitle(t *testing.T) {
	got, err := parseArgs([]string{"10s"})
	if err != nil {
		t.Fatalf("parseArgs error: %v", err)
	}
	if got.title != defaultTitle {
		t.Errorf("title = %q, want default %q", got.title, defaultTitle)
	}
}

func TestParseArgsErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"no args", nil},
		{"unknown flag", []string{"--nope"}},
		{"two durations", []string{"10s", "20s"}},
		{"duration and preset", []string{"10s", "--preset", "pomodoro"}},
		{"preset and until", []string{"--preset", "pomodoro", "--until", "15:00"}},
		{"title without value", []string{"10s", "--title"}},
		{"until without value", []string{"--until"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseArgs(tt.args); err == nil {
				t.Fatalf("parseArgs(%v) expected error, got nil", tt.args)
			}
		})
	}
}

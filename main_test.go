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

func TestParseArgsExec(t *testing.T) {
	got, err := parseArgs([]string{"10s", "--exec", "say hi"})
	if err != nil {
		t.Fatalf("parseArgs error: %v", err)
	}
	if got.execCmd != "say hi" {
		t.Errorf("execCmd = %q, want %q", got.execCmd, "say hi")
	}

	if _, err := parseArgs([]string{"10s", "--exec"}); err == nil {
		t.Fatal("--exec without a value expected error")
	}
}

func TestParseArgsCountUp(t *testing.T) {
	// --up with no duration: stopwatch, no goal.
	got, err := parseArgs([]string{"--up"})
	if err != nil {
		t.Fatalf("parseArgs(--up) error: %v", err)
	}
	if !got.countUp {
		t.Error("expected countUp = true")
	}
	if got.duration != 0 {
		t.Errorf("expected no goal, got duration %v", got.duration)
	}

	// --up with a duration: that's the goal.
	got, err = parseArgs([]string{"--up", "30m"})
	if err != nil {
		t.Fatalf("parseArgs(--up 30m) error: %v", err)
	}
	if !got.countUp || got.duration != 30*time.Minute {
		t.Errorf("expected count-up goal 30m, got countUp=%v duration=%v", got.countUp, got.duration)
	}
}

func TestParseArgsPomodoroPresetStartsCycle(t *testing.T) {
	got, err := parseArgs([]string{"--preset", "pomodoro"})
	if err != nil {
		t.Fatalf("parseArgs error: %v", err)
	}
	if !got.pomodoro {
		t.Error("expected pomodoro cycle")
	}
	if got.duration != 0 {
		t.Errorf("duration = %v, want cycle-managed duration", got.duration)
	}
	if got.title != "POMODORO" {
		t.Errorf("title = %q, want POMODORO", got.title)
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
		{"up with preset", []string{"--up", "--preset", "pomodoro"}},
		{"up with until", []string{"--up", "--until", "15:00"}},
		{"up with bad duration", []string{"--up", "nope"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseArgs(tt.args); err == nil {
				t.Fatalf("parseArgs(%v) expected error, got nil", tt.args)
			}
		})
	}
}

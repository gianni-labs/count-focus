package main

import "testing"

func TestRenderProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		progress float64
		want     string
	}{
		{"empty", 10, 0, "[░░░░░░░░░░] 0%"},
		{"full", 10, 1, "[██████████] 100%"},
		{"half", 10, 0.5, "[█████░░░░░] 50%"},
		{"clamps below zero", 10, -0.5, "[░░░░░░░░░░] 0%"},
		{"clamps above one", 10, 1.5, "[██████████] 100%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderProgressBar(tt.width, tt.progress); got != tt.want {
				t.Fatalf("renderProgressBar(%d, %v) = %q, want %q", tt.width, tt.progress, got, tt.want)
			}
		})
	}
}

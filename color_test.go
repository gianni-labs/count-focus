package main

import (
	"io"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// fakeEnviron lets us simulate NO_COLOR without touching the real process
// environment, which would leak across tests.
type fakeEnviron map[string]string

func (f fakeEnviron) Environ() []string {
	out := make([]string, 0, len(f))
	for k, v := range f {
		out = append(out, k+"="+v)
	}
	return out
}

func (f fakeEnviron) Getenv(key string) string { return f[key] }

// TestStylesRespectNoColor locks in that our lipgloss styles honor the
// NO_COLOR convention (https://no-color.org/) via termenv's built-in
// detection, so styled output degrades to plain text when set.
func TestStylesRespectNoColor(t *testing.T) {
	r := lipgloss.NewRenderer(io.Discard, termenv.WithEnvironment(fakeEnviron{"NO_COLOR": "1"}))
	style := r.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))

	if got := style.Render("hi"); got != "hi" {
		t.Fatalf("expected NO_COLOR to strip all styling, got %q", got)
	}
}

package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	countdownTickInterval = 200 * time.Millisecond
	confettiTickInterval  = 180 * time.Millisecond
)

type tickMsg time.Time
type confettiTickMsg time.Time

type model struct {
	total         time.Duration
	remaining     time.Duration
	startedAt     time.Time
	done          bool
	confettiFrame int
}

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Align(lipgloss.Center)

	timeStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")).
		Align(lipgloss.Center).
		MarginTop(1).
		MarginBottom(1)

	doneStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("42")).
		Align(lipgloss.Center).
		MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Align(lipgloss.Center)

	confettiStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Align(lipgloss.Center).
		MarginBottom(1)
)

var confettiFrames = []string{
	"✦   .   *   +   .   ✦\n  .   +   •   *   .\n*   .   ✦   .   +   *",
	"  *   ✦   .   •   +\n+   .   *   .   ✦   .\n  .   •   +   *   .",
	".   +   •   ✦   *   .\n  ✦   .   +   .   *\n+   *   .   •   .   ✦",
	"+   .   ✦   *   .   •\n  *   •   .   +   ✦\n.   ✦   +   .   *   .",
}

// RunCountdown starts the terminal UI for the given duration.
func RunCountdown(duration time.Duration) error {
	_, err := tea.NewProgram(newModel(duration), tea.WithAltScreen()).Run()
	return err
}

func newModel(duration time.Duration) model {
	return model{
		total:     duration,
		remaining: duration,
		startedAt: time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return countdownTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tickMsg:
		if m.done {
			return m, confettiTick()
		}

		remaining := m.total - time.Since(m.startedAt)
		if remaining <= 0 {
			m.remaining = 0
			m.done = true
			return m, confettiTick()
		}

		m.remaining = remaining
		return m, countdownTick()

	case confettiTickMsg:
		if m.done {
			m.confettiFrame = (m.confettiFrame + 1) % len(confettiFrames)
			return m, confettiTick()
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.done {
		return centerBlock(strings.Join([]string{
			doneStyle.Render("Done!"),
			confettiStyle.Render(confettiFrames[m.confettiFrame]),
			helpStyle.Render("Press q, Esc or Ctrl+C to quit"),
		}, "\n\n"))
	}

	return centerBlock(strings.Join([]string{
		titleStyle.Render("COUNTDOWN"),
		timeStyle.Render(FormatRemaining(m.remaining)),
		helpStyle.Render("Press q, Esc or Ctrl+C to quit"),
	}, "\n"))
}

func countdownTick() tea.Cmd {
	return tea.Tick(countdownTickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func confettiTick() tea.Cmd {
	return tea.Tick(confettiTickInterval, func(t time.Time) tea.Msg {
		return confettiTickMsg(t)
	})
}

func centerBlock(content string) string {
	return "\n\n" + lipgloss.PlaceHorizontal(40, lipgloss.Center, content) + "\n"
}

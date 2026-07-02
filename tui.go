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
	width         int
	height        int
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

	largeTimeStyle = lipgloss.NewStyle().
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

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

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
		return m.placeContent(strings.Join([]string{
			doneStyle.Render("Done!"),
			confettiStyle.Render(confettiFrames[m.confettiFrame]),
			helpStyle.Render("Press q, Esc or Ctrl+C to quit"),
		}, "\n\n"))
	}

	remaining := FormatRemaining(m.remaining)
	return m.placeContent(strings.Join([]string{
		titleStyle.Render("COUNT FOCUS"),
		m.renderTime(remaining),
		helpStyle.Render("Press q, Esc or Ctrl+C to quit"),
	}, "\n"))
}

func (m model) renderTime(remaining string) string {
	if m.canUseLargeTime(remaining) {
		return largeTimeStyle.Render(renderLargeTime(remaining))
	}

	return timeStyle.Render(remaining)
}

func (m model) canUseLargeTime(remaining string) bool {
	if m.width <= 0 || m.height <= 0 {
		return false
	}

	return m.width >= largeTimeWidth(remaining)+4 && m.height >= largeTimeGlyphHeight+7
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

func (m model) placeContent(content string) string {
	width := m.width
	height := m.height

	if width <= 0 {
		width = 40
	}
	if height <= 0 {
		return "\n\n" + lipgloss.PlaceHorizontal(width, lipgloss.Center, content) + "\n"
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

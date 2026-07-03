package main

import (
	"fmt"
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
	title         string
	total         time.Duration
	remaining     time.Duration
	elapsed       time.Duration
	runningSince  time.Time
	paused        bool
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

	progressBarStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("212")).
				Align(lipgloss.Center).
				MarginBottom(1)

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

// RunCountdown starts the terminal UI for the given duration and title.
func RunCountdown(duration time.Duration, title string) error {
	_, err := tea.NewProgram(newModel(duration, title), tea.WithAltScreen()).Run()
	return err
}

func newModel(duration time.Duration, title string) model {
	return model{
		title:        title,
		total:        duration,
		remaining:    duration,
		runningSince: time.Now(),
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
		case " ":
			return m.togglePause()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if m.done || m.paused {
			if m.done {
				return m, confettiTick()
			}
			return m, nil
		}

		remaining := m.total - m.elapsed - time.Since(m.runningSince)
		if remaining <= 0 {
			m.remaining = 0
			m.done = true
			return m, tea.Batch(confettiTick(), bell())
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

// togglePause pauses/resumes the countdown, preserving the exact remaining
// time across the transition instead of relying on a fixed startedAt.
func (m model) togglePause() (tea.Model, tea.Cmd) {
	if m.done {
		return m, nil
	}

	if m.paused {
		m.paused = false
		m.runningSince = time.Now()
		return m, countdownTick()
	}

	m.elapsed += time.Since(m.runningSince)
	m.paused = true
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
	lines := []string{
		titleStyle.Render(m.title),
		m.renderTime(remaining),
	}

	if m.canShowProgressBar() {
		lines = append(lines, progressBarStyle.Render(renderProgressBar(m.progressBarWidth(), m.progress())))
	}

	lines = append(lines, helpStyle.Render(m.helpText()))

	return m.placeContent(strings.Join(lines, "\n"))
}

func (m model) helpText() string {
	if m.paused {
		return "Paused — press Space to resume, q/Esc/Ctrl+C to quit"
	}
	return "Press Space to pause, q/Esc/Ctrl+C to quit"
}

// progress returns how much of the countdown has elapsed, in [0, 1].
func (m model) progress() float64 {
	if m.total <= 0 {
		return 1
	}
	return float64(m.total-m.remaining) / float64(m.total)
}

func (m model) canShowProgressBar() bool {
	return m.width >= progressBarMinWidth+16
}

func (m model) progressBarWidth() int {
	width := m.width - 16
	if width > progressBarMaxWidth {
		width = progressBarMaxWidth
	}
	return width
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

// bell rings the terminal bell once the countdown finishes, so the user
// notices even if they're not looking at the screen. The BEL control
// character has no effect on cursor position or the alt-screen buffer.
func bell() tea.Cmd {
	return func() tea.Msg {
		fmt.Print("\a")
		return nil
	}
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

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
	title   string
	countUp bool
	// display holds the time shown on screen: time remaining in countdown
	// mode, or time elapsed in count-up mode.
	display time.Duration
	// total is the countdown length; target is the optional count-up goal
	// (0 means no goal — count up forever).
	total         time.Duration
	target        time.Duration
	elapsed       time.Duration
	runningSince  time.Time
	paused        bool
	done          bool
	reachedGoal   bool // count-up: the goal was crossed and the bell already rang
	confettiFrame int
	width         int
	height        int
}

// elapsedNow returns the total elapsed running time, accounting for pauses.
func (m model) elapsedNow() time.Duration {
	if m.paused {
		return m.elapsed
	}
	return m.elapsed + time.Since(m.runningSince)
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

// RunCountdown starts the terminal UI for the given options.
func RunCountdown(opts options) error {
	_, err := tea.NewProgram(newModel(opts), tea.WithAltScreen()).Run()
	return err
}

func newModel(opts options) model {
	m := model{
		title:        opts.title,
		countUp:      opts.countUp,
		runningSince: time.Now(),
	}
	if opts.countUp {
		m.target = opts.duration // 0 = count up with no goal
		m.display = 0
	} else {
		m.total = opts.duration
		m.display = opts.duration
	}
	return m
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

		if m.countUp {
			m.display = m.elapsedNow()
			if m.target > 0 && !m.reachedGoal && m.display >= m.target {
				m.reachedGoal = true
				return m, tea.Batch(countdownTick(), bell())
			}
			return m, countdownTick()
		}

		remaining := m.total - m.elapsedNow()
		if remaining <= 0 {
			m.display = 0
			m.done = true
			return m, tea.Batch(confettiTick(), bell())
		}

		m.display = remaining
		return m, countdownTick()

	case confettiTickMsg:
		if m.done {
			m.confettiFrame = (m.confettiFrame + 1) % len(confettiFrames)
			return m, confettiTick()
		}
	}

	return m, nil
}

// togglePause pauses/resumes the timer, preserving the exact elapsed time
// across the transition instead of relying on a fixed start timestamp. Works
// the same in countdown and count-up modes.
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

	shown := FormatRemaining(m.display)
	lines := []string{
		titleStyle.Render(m.title),
		m.renderTime(shown),
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

// progress returns how much of the timer has elapsed, in [0, 1]. In count-up
// mode it measures progress toward the goal.
func (m model) progress() float64 {
	if m.countUp {
		if m.target <= 0 {
			return 0
		}
		p := float64(m.display) / float64(m.target)
		if p > 1 {
			return 1
		}
		return p
	}
	if m.total <= 0 {
		return 1
	}
	return float64(m.total-m.display) / float64(m.total)
}

func (m model) canShowProgressBar() bool {
	// Count-up with no goal has no total to measure against.
	if m.countUp && m.target <= 0 {
		return false
	}
	return m.width >= progressBarMinWidth+16
}

func (m model) progressBarWidth() int {
	width := m.width - 16
	if width > progressBarMaxWidth {
		width = progressBarMaxWidth
	}
	return width
}

func (m model) renderTime(shown string) string {
	// Once a count-up goal is reached, tint the time green as a signal.
	large := m.canUseLargeTime(shown)
	style := timeStyle
	if large {
		style = largeTimeStyle
	}
	if m.reachedGoal {
		style = style.Foreground(lipgloss.Color("42"))
	}

	if large {
		return style.Render(renderLargeTime(shown))
	}
	return style.Render(shown)
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

package main

import (
	"fmt"
	"os"
	"time"
)

const helpText = `Usage:
  count-focus <duration>
  count-focus --preset <name>

Examples:
  count-focus 10s
  count-focus 5m
  count-focus 1h30m
  count-focus 1h30m10s
  count-focus --preset pomodoro

Presets:
  Built-in: pomodoro (25m), short-break (5m), long-break (15m)
  Custom:   ~/.config/count-focus/presets.conf ("name = duration" per line)

Flags:
  --preset, -p     Start a named preset
  --help, -h       Show this help
  --version, -v    Show version

Keys:
  Space            Pause/Resume
  q, Esc, Ctrl+C   Quit
`

// version is set at build time via -ldflags "-X main.version=...".
// Homebrew installs set it to the release tag; local builds keep "dev".
var version = "dev"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr)
		fmt.Fprint(os.Stderr, helpText)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Print(helpText)
		return nil
	}

	if len(args) == 1 && (args[0] == "--version" || args[0] == "-v") {
		fmt.Println("count-focus " + version)
		return nil
	}

	duration, err := parseArgs(args)
	if err != nil {
		return err
	}

	return RunCountdown(duration)
}

// parseArgs resolves the CLI arguments to a countdown duration, accepting
// either a bare duration or a --preset/-p <name> pair.
func parseArgs(args []string) (time.Duration, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("missing duration")
	}

	if args[0] == "--preset" || args[0] == "-p" {
		if len(args) == 1 {
			return 0, fmt.Errorf("missing preset name")
		}
		if len(args) > 2 {
			return 0, fmt.Errorf("too many arguments")
		}
		return resolvePreset(args[1])
	}

	if len(args) > 1 {
		return 0, fmt.Errorf("too many arguments")
	}

	duration, err := ParseDuration(args[0])
	if err != nil {
		return 0, fmt.Errorf("invalid duration: %s\n%s", args[0], invalidDurationMessage)
	}
	return duration, nil
}

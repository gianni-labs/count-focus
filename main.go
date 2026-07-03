package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const defaultTitle = "COUNT FOCUS"

const helpText = `Usage:
  count-focus <duration>
  count-focus --preset <name>
  count-focus --until <HH:MM>

Examples:
  count-focus 10s
  count-focus 1h30m
  count-focus 25m --title "Deep work"
  count-focus --preset pomodoro
  count-focus --until 15:00

Presets:
  Built-in: pomodoro (25m), short-break (5m), long-break (15m)
  Custom:   ~/.config/count-focus/presets.conf ("name = duration" per line)

Flags:
  --title, -t      Set the on-screen title
  --preset, -p     Start a named preset
  --until, -u      Count down until a wall-clock time today (HH:MM or HH:MM:SS)
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

	opts, err := parseArgs(args)
	if err != nil {
		return err
	}

	return RunCountdown(opts.duration, opts.title)
}

// options holds the resolved CLI configuration for a run.
type options struct {
	duration time.Duration
	title    string
}

// parseArgs resolves the CLI arguments. The countdown length comes from exactly
// one of: a bare <duration>, --preset/-p <name>, or --until/-u <HH:MM>. An
// optional --title/-t overrides the on-screen title (a preset defaults its
// title to the preset name).
func parseArgs(args []string) (options, error) {
	var (
		opts        options
		durationArg string
		presetArg   string
		untilArg    string
		titleSet    bool
	)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--title", "-t":
			i++
			if i >= len(args) {
				return options{}, fmt.Errorf("missing value for %s", arg)
			}
			opts.title = args[i]
			titleSet = true
		case "--preset", "-p":
			i++
			if i >= len(args) {
				return options{}, fmt.Errorf("missing preset name")
			}
			presetArg = args[i]
		case "--until", "-u":
			i++
			if i >= len(args) {
				return options{}, fmt.Errorf("missing value for %s", arg)
			}
			untilArg = args[i]
		default:
			if strings.HasPrefix(arg, "-") {
				return options{}, fmt.Errorf("unknown flag %q", arg)
			}
			if durationArg != "" {
				return options{}, fmt.Errorf("too many arguments")
			}
			durationArg = arg
		}
	}

	sources := 0
	for _, s := range []string{durationArg, presetArg, untilArg} {
		if s != "" {
			sources++
		}
	}
	if sources == 0 {
		return options{}, fmt.Errorf("missing duration")
	}
	if sources > 1 {
		return options{}, fmt.Errorf("use only one of: <duration>, --preset, or --until")
	}

	switch {
	case durationArg != "":
		d, err := ParseDuration(durationArg)
		if err != nil {
			return options{}, fmt.Errorf("invalid duration: %s\n%s", durationArg, invalidDurationMessage)
		}
		opts.duration = d
	case presetArg != "":
		d, err := resolvePreset(presetArg)
		if err != nil {
			return options{}, err
		}
		opts.duration = d
		if !titleSet {
			opts.title = strings.ToUpper(presetArg)
		}
	case untilArg != "":
		d, err := parseUntil(untilArg, time.Now())
		if err != nil {
			return options{}, err
		}
		opts.duration = d
	}

	if opts.title == "" {
		opts.title = defaultTitle
	}
	return opts, nil
}

package main

import (
	"fmt"
	"os"
)

const helpText = `Usage:
  countdown <duration>

Examples:
  countdown 10s
  countdown 5m
  countdown 1h30m
  countdown 1h30m10s

Keys:
  q, Esc, Ctrl+C   Quit
`

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

	if len(args) == 0 {
		return fmt.Errorf("missing duration")
	}
	if len(args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	duration, err := ParseDuration(args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %s\n%s", args[0], invalidDurationMessage)
	}

	return RunCountdown(duration)
}

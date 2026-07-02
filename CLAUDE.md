# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build ./...          # build
go run . <duration>     # run without installing, e.g. go run . 10s
go test ./...           # run all tests
go test -run TestName ./...   # run a single test
go vet ./...             # static checks
```

There is no separate lint config — `go vet` and `gofmt` are the baseline. Format with `gofmt -l .` before committing.

## Architecture

Single-package Go CLI/TUI (`package main`, module `github.com/gianni-labs/count-focus`) built on Bubble Tea (`charmbracelet/bubbletea`) and Lip Gloss for styling. All files live at repo root, split by responsibility:

- `main.go` — arg parsing entrypoint (`run`). Delegates to `ParseDuration` then `RunCountdown`.
- `duration.go` — `ParseDuration` parses strings like `1h30m10s` into `time.Duration`. Units must appear at most once, strictly in descending order (h, m, s); anything else is a validation error listed in `invalidDurationMessage`.
- `tui.go` — the Bubble Tea `model` and `Update`/`View` loop. Countdown ticks every 200ms (`countdownTickInterval`); once `remaining <= 0`, `model.done` flips and a separate confetti animation ticks every 180ms cycling through `confettiFrames`. Quit keys: `q`, `esc`, `ctrl+c`.
- `format.go` — `FormatRemaining` renders a duration as `MM:SS` or `HH:MM:SS` (rounds up to the next whole second).
- `bigtime.go` — ASCII-art "big digit" glyphs (`largeTimeGlyphs`) used to render the time large when the terminal is big enough; `tui.go`'s `canUseLargeTime` decides based on `m.width`/`m.height` whether to use `largeTimeStyle` (big glyphs) vs the normal `timeStyle`.
- `duration_test.go`, `format_test.go` — table-driven tests for the two pure logic files above; the TUI itself has no tests (visual/interactive).

Data flow: `main.run` → `ParseDuration` (string → `time.Duration`) → `RunCountdown` → `tea.NewProgram` with alt-screen → ticks drive `remaining` down → `FormatRemaining` for display, optionally upscaled via `renderLargeTime`.

## Distribution

Released via a Homebrew tap (`Formula/count-focus.rb`), built from a GitHub release tarball with `go build` — no other packaging or CI pipeline in this repo.

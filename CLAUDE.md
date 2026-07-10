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

- `main.go` — arg parsing entrypoint (`run`). Handles `--help`/`-h` and `--version`/`-v`; `parseArgs` walks the args into an `options{countUp, duration, title}`. Countdown mode needs exactly one duration source: a bare `<duration>`, `--preset`/`-p <name>`, or `--until`/`-u <HH:MM>`. `--up` switches to count-up mode where the bare duration (if any) is an optional goal and `--preset`/`--until` are rejected. `--title`/`-t` and `--exec`/`-e <shell command>` are optional and can appear in any position. `run` then calls `RunCountdown(opts)`. `version` is a package var overridable via `-ldflags -X main.version=...` (see Distribution); local builds default to `"dev"`.
- `preset.go` — named presets. `builtinPresets` (pomodoro/short-break/long-break) merged with an optional user config at `~/.config/count-focus/presets.conf` (honors `XDG_CONFIG_HOME`), a `name = duration` per-line format parsed via `ParseDuration`. `loadPresets` merges them (config overrides/extends builtins); `resolvePreset` looks one up. A preset with no explicit `--title` defaults its title to the uppercased preset name.
- `until.go` — `parseUntil(s, now)` turns a 24h wall-clock time (`HH:MM`/`HH:MM:SS`) into a duration from `now` to that time today, erroring if it's already past. `now` is injected for testability.
- `duration.go` — `ParseDuration` parses strings like `1h30m10s` into `time.Duration`. Units must appear at most once, strictly in descending order (h, m, s); anything else is a validation error listed in `invalidDurationMessage`.
- `tui.go` — the Bubble Tea `model` and `Update`/`View` loop, shared by both modes. Ticks every 200ms (`countdownTickInterval`). Pause/resume (`Space`) works by tracking `elapsed` (accumulated running time) + `runningSince` (start of the current running segment) instead of a fixed start timestamp — see `model.togglePause`/`elapsedNow`. The `display` field is the time shown: remaining in countdown mode, elapsed in count-up mode (`model.countUp`). Countdown: once `display` hits 0, `model.done` flips, the bell rings once (`bell()`), and confetti animates. Count-up: runs forever; if a `target` goal is set, crossing it rings the bell once (`reachedGoal`) and tints the time green, without stopping. Both endpoints also fire `runExec(m.execCmd)` (the `--exec` command via `sh -c`, fire-and-forget). Quit keys: `q`, `esc`, `ctrl+c`.
- `format.go` — `FormatRemaining` renders a duration as `MM:SS` or `HH:MM:SS` (rounds up to the next whole second).
- `bigtime.go` — ASCII-art "big digit" glyphs (`largeTimeGlyphs`) used to render the time large when the terminal is big enough; `tui.go`'s `canUseLargeTime` decides based on `m.width`/`m.height` whether to use `largeTimeStyle` (big glyphs) vs the normal `timeStyle`.
- `progressbar.go` — `renderProgressBar` is a pure function producing `[████░░░░] 60%`; `tui.go`'s `canShowProgressBar`/`progressBarWidth` decide, based on terminal width, whether/how wide to render it (same width-gating pattern as `canUseLargeTime`).
- `duration_test.go`, `format_test.go`, `progressbar_test.go`, `until_test.go`, `preset_test.go`, `main_test.go`, `exec_test.go` — table-driven tests for the pure logic and arg parsing above (`exec_test.go` also runs a real command to verify a side effect).
- `color_test.go` — regression test confirming Lip Gloss/termenv honor `NO_COLOR` (https://no-color.org/) out of the box; uses a fake `termenv.Environ` so it doesn't touch real process env vars.
- The TUI's `Update`/`View` wiring itself has no tests (visual/interactive).

Data flow: `main.run` → `parseArgs` (args → `options`) → `RunCountdown(opts)` → `tea.NewProgram` with alt-screen → ticks update `display` (down in countdown, up in count-up; frozen while paused) → `FormatRemaining` renders it, optionally upscaled via `renderLargeTime`, plus `renderProgressBar` when there's room.

## Distribution

Released via a Homebrew tap (`Formula/count-focus.rb`), built from a GitHub release tarball with `go build -ldflags "-s -w -X main.version=vX.Y.Z"`. CI (`.github/workflows/ci.yml`) runs `gofmt`/`go vet`/`go test`/`go build` on push to `main`/`dev` and on PRs. Full release checklist (tagging, computing the tarball sha256, updating the Formula) is in `docs/RELEASING.md` (local-only, gitignored like the rest of `docs/`).

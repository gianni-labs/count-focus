# Count Focus

A focused, terminal-native timer for deep work.

[Español](README.es.md)

![Count Focus demo](assets/count-focus-demo.gif)

## Installation

With Homebrew:

```bash
brew tap gianni-labs/tap
brew trust --formula gianni-labs/tap/count-focus
brew install count-focus
```

This installs the `count-focus` command.

## Usage

```bash
count-focus <duration>
```

Examples:

```bash
count-focus 10s
count-focus 5m
count-focus 1h
count-focus 1h30m
count-focus 1h30m10s
```

### Title

Set a title to keep the timer tied to the task at hand:

```bash
count-focus 25m --title "Write report"
count-focus 1h -t "Deep work"
```

### Count down to a specific time

Instead of a duration, target a wall-clock time in 24-hour format. Count Focus will count down until that time today:

```bash
count-focus --until 15:00      # until 3:00 PM
count-focus -u 15:30:30        # including seconds
```

It returns an error if that time has already passed today.

### Stopwatch (count up)

Use `--up` to count up from zero:

```bash
count-focus --up          # runs until you quit
count-focus --up 30m      # optional goal; alerts and turns green when reached
```

### Run a command when time ends

Use `--exec` to run a shell command when a countdown ends, or when a goal is reached in `--up` mode:

```bash
count-focus 25m --exec "open -a Slack"
count-focus 10m --exec "say 'time is up'"    # macOS voice
```

The command starts without blocking the final screen. It can also trigger native notifications, for example on macOS:

```bash
count-focus 25m --exec 'osascript -e "display notification \"Done\" with title \"count-focus\""'
```

### Keys

While the timer is running:

- `Space` — pause / resume
- `q`, `Esc`, `Ctrl+C` — quit

## Presets

Use a named preset instead of a duration:

```bash
count-focus --preset pomodoro     # full cycle: 4 × 25m, 5m breaks, then a 15m break
count-focus -p short-break        # 5m
count-focus -p long-break         # 15m
```

`pomodoro` starts the complete standard cycle: four 25-minute focus rounds, 5-minute short breaks between rounds, and a final 15-minute long break. The screen shows the current round and phase, and the terminal bell rings at every transition. When provided, `--exec` runs after the complete cycle ends.

### Custom presets

Create this file to define your own presets or change the built-in breaks:

```
~/.config/count-focus/presets.conf
```

Use one `name = duration` entry per line:

```conf
short-break = 10m
deep-work = 90m
review = 45m
```

Custom presets override or extend the built-in `short-break` and `long-break` presets. `pomodoro` is reserved for the complete standard cycle. See [`examples/presets.conf`](examples/presets.conf) for a full example.

## Version

```bash
count-focus --version
```

## License

MIT

## Update

```bash
brew update
brew upgrade count-focus
```

## Uninstall

```bash
brew uninstall count-focus
```

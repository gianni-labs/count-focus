package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// builtinPresets are always available, even without a config file. A user
// config file can override any of these or add new ones.
var builtinPresets = map[string]time.Duration{
	"pomodoro":    25 * time.Minute,
	"short-break": 5 * time.Minute,
	"long-break":  15 * time.Minute,
}

// resolvePreset returns the duration for a named preset, consulting the user
// config file (which overrides/extends the builtins) before the builtins.
func resolvePreset(name string) (time.Duration, error) {
	presets, err := loadPresets()
	if err != nil {
		return 0, err
	}

	d, ok := presets[name]
	if !ok {
		return 0, fmt.Errorf("unknown preset %q\navailable presets: %s", name, presetNames(presets))
	}
	return d, nil
}

// presetConfigPath returns the path to the user's presets file, honoring
// XDG_CONFIG_HOME and falling back to ~/.config.
func presetConfigPath() (string, error) {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "count-focus", "presets.conf"), nil
}

// loadPresets returns the builtin presets merged with any defined in the user
// config file. A missing config file is not an error; a malformed one is.
func loadPresets() (map[string]time.Duration, error) {
	presets := make(map[string]time.Duration, len(builtinPresets))
	for name, d := range builtinPresets {
		presets[name] = d
	}

	path, err := presetConfigPath()
	if err != nil {
		// Can't locate a home/config dir: fall back to builtins only.
		return presets, nil
	}

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return presets, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		name, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("%s:%d: expected 'name = duration', got %q", path, lineNo, line)
		}
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if name == "" {
			return nil, fmt.Errorf("%s:%d: empty preset name", path, lineNo)
		}

		d, err := ParseDuration(value)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: invalid duration %q for preset %q", path, lineNo, value, name)
		}
		presets[name] = d
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return presets, nil
}

// presetNames returns the preset names sorted alphabetically, for help/errors.
func presetNames(presets map[string]time.Duration) string {
	names := make([]string, 0, len(presets))
	for name := range presets {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

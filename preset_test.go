package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// writePresetConfig points XDG_CONFIG_HOME at a temp dir and, if contents is
// non-empty, writes a presets.conf there. Returns nothing; cleanup is handled
// by t.TempDir / t.Setenv.
func writePresetConfig(t *testing.T, contents string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	if contents == "" {
		return
	}
	confDir := filepath.Join(dir, "count-focus")
	if err := os.MkdirAll(confDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(confDir, "presets.conf"), []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestResolvePresetBuiltin(t *testing.T) {
	writePresetConfig(t, "") // no config file, builtins only

	tests := map[string]time.Duration{
		"short-break": 5 * time.Minute,
		"long-break":  15 * time.Minute,
	}
	for name, want := range tests {
		got, err := resolvePreset(name)
		if err != nil {
			t.Fatalf("resolvePreset(%q) error: %v", name, err)
		}
		if got != want {
			t.Fatalf("resolvePreset(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestResolvePresetUnknown(t *testing.T) {
	writePresetConfig(t, "")
	if _, err := resolvePreset("nope"); err == nil {
		t.Fatal("resolvePreset(\"nope\") expected error, got nil")
	}
}

func TestPresetConfigOverridesAndExtends(t *testing.T) {
	writePresetConfig(t, "# my presets\nshort-break = 10m\ndeep-work = 90m\n")

	presets, err := loadPresets()
	if err != nil {
		t.Fatalf("loadPresets error: %v", err)
	}

	if presets["short-break"] != 10*time.Minute {
		t.Errorf("config should override builtin short-break: got %v", presets["short-break"])
	}
	if presets["deep-work"] != 90*time.Minute {
		t.Errorf("config should add deep-work: got %v", presets["deep-work"])
	}
	if presets["long-break"] != 15*time.Minute {
		t.Errorf("untouched builtin should remain: got %v", presets["long-break"])
	}
}

func TestPresetConfigRejectsReservedPomodoroName(t *testing.T) {
	writePresetConfig(t, "pomodoro = 30m\n")
	if _, err := loadPresets(); err == nil {
		t.Fatal("loadPresets with reserved pomodoro expected error")
	}
}

func TestPresetConfigInvalid(t *testing.T) {
	writePresetConfig(t, "broken = notaduration\n")
	if _, err := loadPresets(); err == nil {
		t.Fatal("loadPresets with invalid duration expected error, got nil")
	}
}

func TestParseArgsPreset(t *testing.T) {
	writePresetConfig(t, "")

	got, err := parseArgs([]string{"--preset", "pomodoro"})
	if err != nil {
		t.Fatalf("parseArgs preset error: %v", err)
	}
	if !got.pomodoro {
		t.Fatal("parseArgs --preset pomodoro should start the Pomodoro cycle")
	}
	if got.title != "POMODORO" {
		t.Fatalf("preset should default title to its name: got %q", got.title)
	}

	if _, err := parseArgs([]string{"-p"}); err == nil {
		t.Fatal("parseArgs -p without name expected error")
	}
}

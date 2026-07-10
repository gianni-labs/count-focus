package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunExecEmpty(t *testing.T) {
	if runExec("") != nil {
		t.Fatal("runExec(\"\") should return a nil command (no-op)")
	}
}

func TestRunExecRunsCommand(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "done")
	cmd := runExec("touch " + marker)
	if cmd == nil {
		t.Fatal("runExec returned nil for a non-empty command")
	}

	cmd() // fire the command (fire-and-forget; the child runs asynchronously)

	// Poll for the side effect, since the command runs in the background.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(marker); err == nil {
			return // success
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("command did not create %s within the timeout", marker)
}

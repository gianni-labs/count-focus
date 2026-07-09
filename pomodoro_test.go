package main

import (
	"testing"
	"time"
)

func TestPomodoroTransitions(t *testing.T) {
	m := newModel(options{pomodoro: true, title: "POMODORO"})
	if m.display != pomodoroFocusDuration || m.pomodoroRound != 1 || m.pomodoroPhase != pomodoroFocus {
		t.Fatalf("unexpected initial Pomodoro state: %+v", m)
	}

	// Finish focus in round 1: begin the short break.
	m.runningSince = time.Now().Add(-pomodoroFocusDuration)
	next, _ := m.Update(tickMsg(time.Now()))
	m = next.(model)
	if m.pomodoroPhase != pomodoroBreak || m.display != pomodoroShortBreak || m.done {
		t.Fatalf("after first focus = phase %v, display %v, done %v", m.pomodoroPhase, m.display, m.done)
	}

	// Finish the short break: begin focus in round 2.
	m.runningSince = time.Now().Add(-pomodoroShortBreak)
	next, _ = m.Update(tickMsg(time.Now()))
	m = next.(model)
	if m.pomodoroPhase != pomodoroFocus || m.pomodoroRound != 2 || m.display != pomodoroFocusDuration {
		t.Fatalf("after first break = round %d, phase %v, display %v", m.pomodoroRound, m.pomodoroPhase, m.display)
	}
}

func TestPomodoroLongBreakCompletesCycle(t *testing.T) {
	m := newModel(options{pomodoro: true, title: "POMODORO"})
	m.pomodoroRound = pomodoroRounds
	m.runningSince = time.Now().Add(-pomodoroFocusDuration)

	next, _ := m.Update(tickMsg(time.Now()))
	m = next.(model)
	if m.pomodoroPhase != pomodoroBreak || m.display != pomodoroLongBreak || m.done {
		t.Fatalf("after last focus = phase %v, display %v, done %v", m.pomodoroPhase, m.display, m.done)
	}

	m.runningSince = time.Now().Add(-pomodoroLongBreak)
	next, _ = m.Update(tickMsg(time.Now()))
	m = next.(model)
	if !m.done || m.display != 0 {
		t.Fatalf("after long break = done %v, display %v", m.done, m.display)
	}
}

func TestPomodoroStatus(t *testing.T) {
	m := model{pomodoro: true, pomodoroRound: 4, pomodoroPhase: pomodoroBreak}
	if got, want := m.pomodoroStatus(), "Round 4/4 — LONG BREAK"; got != want {
		t.Errorf("pomodoroStatus() = %q, want %q", got, want)
	}
}

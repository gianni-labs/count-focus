package main

import (
	"fmt"
	"strings"
)

const (
	progressBarMinWidth = 10
	progressBarMaxWidth = 40
)

// renderProgressBar renders a filled/empty bar plus percentage, e.g. "[████████░░░░] 60%".
func renderProgressBar(width int, progress float64) string {
	if width < 1 {
		width = 1
	}
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filled := int(float64(width)*progress + 0.5)
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	pct := int(progress*100 + 0.5)

	return fmt.Sprintf("[%s] %d%%", bar, pct)
}

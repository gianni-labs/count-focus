package main

import (
	"fmt"
	"time"
)

// FormatRemaining formats a remaining duration as MM:SS or HH:MM:SS.
func FormatRemaining(d time.Duration) string {
	if d < 0 {
		d = 0
	}

	totalSeconds := int((d + time.Second - 1).Truncate(time.Second).Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

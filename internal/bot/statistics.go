package bot

import (
	"fmt"
	"runtime"
	"time"
)

// The time the bot was started.
var runningSince time.Time

// Return a human readable string containing the time the application has been running.
func runningDurationString() string {
	uptime := time.Since(runningSince)

	return formatDuration(uptime)
}

// Format the duration into a human readable string.
func formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	minutes := int(d.Minutes())
	hours := int(d.Hours())
	days := int(d.Hours()) / 24
	months := days / 30
	years := months / 12

	// Choose the appropriate format based on the duration
	switch {
	case years > 0:
		s := fmt.Sprintf("%d year", years)
		if years > 1 {
			s += "s"
		}
		return s
	case months > 0:
		s := fmt.Sprintf("%d month", months)
		if months > 1 {
			s += "s"
		}
		return s
	case days > 0:
		s := fmt.Sprintf("%d day", days)
		if days > 1 {
			s += "s"
		}
		return s
	case hours > 0:
		s := fmt.Sprintf("%d hour", hours)
		if hours > 1 {
			s += "s"
		}
		return s
	case minutes > 0:
		s := fmt.Sprintf("%d minute", minutes)
		if minutes > 1 {
			s += "s"
		}
		return s
	default:
		s := fmt.Sprintf("%d second", seconds)
		if seconds > 1 {
			s += "s"
		}
		return s
	}
}

// Returns a string of the current memory usage of the application.
func memoryUsageString() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if m.Alloc >= 1024*1024*1024 {
		return fmt.Sprintf("%.2f GB\n", bToGb(m.Alloc))
	}
	return fmt.Sprintf("%.2f MB\n", bToMb(m.Alloc))
}

// Converts bytes to MB.
func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// Converts bytes to GB.
func bToGb(b uint64) float64 {
	return float64(b) / 1024 / 1024 / 1024
}

package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type unit struct {
	Symbol   string
	Duration int64
}

func HumanizeTimeAgo(timestamp time.Time) string {
	if timestamp.IsZero() {
		return ""
	}

	timeDelta := time.Now().Unix() - timestamp.Unix()

	units := []unit{
		{"y", 60 * 60 * 24 * 365},
		{"w", 60 * 60 * 24 * 7},
		{"d", 60 * 60 * 24},
		{"h", 60 * 60},
		{"m", 60},
		{"s", 1},
	}

	for _, unit := range units {
		if timeDelta >= unit.Duration {
			return fmt.Sprintf("%3d%s ago", timeDelta/unit.Duration, unit.Symbol)
		}
	}

	return "just now"
}

func HumanizeDuration(durationMs int64) string {
	if durationMs == 0 {
		return ""
	}

	units := []unit{
		{"d", 24 * 60 * 60 * 1000},
		{"h", 60 * 60 * 1000},
		{"m", 60 * 1000},
		{"s", 1000},
		{"ms", 1},
	}

	for _, unit := range units {
		if durationMs >= unit.Duration {
			return fmt.Sprintf("%3d%s", durationMs/unit.Duration, unit.Symbol)
		}
	}

	return fmt.Sprintf("%dms", durationMs)
}

func HumanizePath(path string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if strings.HasPrefix(path, homeDir) {
		return strings.Replace(path, homeDir, "~", 1)
	}

	return path
}

func CenterString(text string, width int, format string) string {
	padding := (width + len(text)) / 2
	return fmt.Sprintf(format, width, fmt.Sprintf("%*s", padding, text))
}

package utils

import (
	"fmt"
	"strings"
	"time"
)

func parseTime(timeStr string, date time.Time) time.Time {
	timeStr = strings.ToLower(strings.TrimSpace(timeStr))
	timeStr = strings.ReplaceAll(timeStr, " ", "")

	var hour, minute int
	isPM := strings.HasSuffix(timeStr, "p") || strings.HasSuffix(timeStr, "pm")
	timeStr = strings.TrimSuffix(strings.TrimSuffix(timeStr, "p"), "m")
	timeStr = strings.TrimSuffix(strings.TrimSuffix(timeStr, "a"), "m")

	if strings.Contains(timeStr, ":") {
		fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	} else {
		fmt.Sscanf(timeStr, "%d", &hour)
	}

	if isPM && hour != 12 {
		hour += 12
	} else if !isPM && hour == 12 {
		hour = 0
	}

	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())
}


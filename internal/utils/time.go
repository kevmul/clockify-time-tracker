package utils

import (
	"fmt"
	"strings"
	"time"
)

// parseTimeRange splits a time range string like "9a - 5p" into start and end times
func ParseTimeRange(timeRange string, date time.Time) (time.Time, time.Time, error) {
	// Split on the dash separator
	parts := strings.Split(timeRange, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("expected format: '9a - 5p'")
	}

	// Parse each part into a time
	startTime := ParseTime(strings.TrimSpace(parts[0]), date)
	endTime := ParseTime(strings.TrimSpace(parts[1]), date)

	return startTime, endTime, nil
}

// parseTime converts a time string like "9a" or "3:30p" to a full time.Time
// It handles various formats: 9a, 9:30a, 9, 9:30
func ParseTime(timeStr string, date time.Time) time.Time {
	// Normalize the string: lowercase, remove spaces
	timeStr = strings.ToLower(strings.TrimSpace(timeStr))
	timeStr = strings.ReplaceAll(timeStr, " ", "")

	var hour, minute int

	// Check if PM (afternoon/evening)
	isPM := strings.HasSuffix(timeStr, "p") || strings.HasSuffix(timeStr, "pm")

	// Remove the am/pm suffix
	timeStr = strings.TrimSuffix(strings.TrimSuffix(timeStr, "p"), "m")
	timeStr = strings.TrimSuffix(strings.TrimSuffix(timeStr, "a"), "m")

	// Parse hour and optional minutes
	if strings.Contains(timeStr, ":") {
		fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	} else {
		fmt.Sscanf(timeStr, "%d", &hour)
	}

	// Convert to 24-hour format
	if isPM && hour != 12 {
		hour += 12 // 1pm = 13, 2pm = 14, etc.
	} else if !isPM && hour == 12 {
		hour = 0 // 12am = midnight = 0
	}

	// Combine the date with our parsed time
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())
}

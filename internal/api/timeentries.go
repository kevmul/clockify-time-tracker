// Functions for creating and fetching time entries in Clockify
package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CreateTimeEntry creates a new time entry in Clockify
// Takes all the necessary parameters and returns an error if creation fails
func (c *Client) CreateTimeEntry(workspaceID, projectID, description, timeRange string, date time.Time) error {
	// Parse the time range string (e.g., "9a - 5p") into actual times
	startTime, endTime, err := parseTimeRange(timeRange, date)
	if err != nil {
		return fmt.Errorf("invalid time range: %w", err)
	}

	// Build the request payload
	entry := TimeEntryRequest{
		Start:       startTime.Format(time.RFC3339), // Convert to RFC3339 format
		End:         endTime.Format(time.RFC3339),
		ProjectID:   projectID,
		Description: description,
	}

	// Build endpoint and make POST request
	endpoint := fmt.Sprintf("/workspaces/%s/time-entries", workspaceID)
	_, err = c.post(endpoint, entry)
	if err != nil {
		return fmt.Errorf("failed to create time entry: %w", err)
	}

	return nil
}

// GetTasks fetches previous time entries and extracts unique task descriptions
// This gives us autocomplete suggestions for the user
func (c *Client) GetTasks(workspaceID, userID string) ([]string, error) {
	// Build endpoint for user's time entries
	endpoint := fmt.Sprintf("/workspaces/%s/user/%s/time-entries", workspaceID, userID)
	
	// Make GET request
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	// Parse response into slice of time entries
	var entries []TimeEntryResponse
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse time entries: %w", err)
	}

	// Extract unique task descriptions using a map as a set
	taskMap := make(map[string]bool)
	var tasks []string
	
	for _, entry := range entries {
		// Only add non-empty descriptions that we haven't seen before
		if entry.Description != "" && !taskMap[entry.Description] {
			taskMap[entry.Description] = true
			tasks = append(tasks, entry.Description)
		}
	}

	return tasks, nil
}

// parseTimeRange splits a time range string like "9a - 5p" into start and end times
func parseTimeRange(timeRange string, date time.Time) (time.Time, time.Time, error) {
	// Split on the dash separator
	parts := strings.Split(timeRange, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("expected format: '9a - 5p'")
	}

	// Parse each part into a time
	startTime := parseTime(strings.TrimSpace(parts[0]), date)
	endTime := parseTime(strings.TrimSpace(parts[1]), date)

	return startTime, endTime, nil
}

// parseTime converts a time string like "9a" or "3:30p" to a full time.Time
// It handles various formats: 9a, 9:30a, 9, 9:30
func parseTime(timeStr string, date time.Time) time.Time {
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

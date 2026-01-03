// Functions for fetching projects from Clockify
package api

import (
	"encoding/json"
	"fmt"
)

// GetProjects fetches all projects for a given workspace
// Returns a slice of Project structs or an error
func (c *Client) GetProjects(workspaceID string) ([]Project, error) {
	// Build the endpoint URL with the workspace ID
	endpoint := fmt.Sprintf("/workspaces/%s/projects", workspaceID)
	
	// Make the GET request
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response into a slice of Project structs
	var projects []Project
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}

	return projects, nil
}


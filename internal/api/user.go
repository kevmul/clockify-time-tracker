// Functions for fetching user information from Clockify
package api

import (
	"encoding/json"
	"fmt"
)

// GetUserInfo fetches the current user's information from Clockify
// This includes their user ID and default workspace ID
// Returns UserInfo or an error if the request fails
func (c *Client) GetUserInfo() (*UserInfo, error) {
	// Make a GET request to /user endpoint
	body, err := c.get("/user")
	if err != nil {
		return nil, err
	}

	// Parse the JSON response into our UserInfo struct
	var user UserInfo
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &user, nil
}


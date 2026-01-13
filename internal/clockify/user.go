// Functions for fetching user information from Clockify
package clockify

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// GetUserInfo fetches the current user's information from Clockify
// This includes their user ID and default workspace ID
// Returns UserInfo or an error if the request fails
func (c *Client) GetUserInfo() (*UserInfo, error) {
	// Make a GET request to /user endpoint
	body, err := c.Get("/user")
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

// fetchUserInfo returns a command that fetches user information
// When complete, it sends a userInfoMsg back to Update()
func fetchUserInfo(apiKey string) tea.Cmd {
	return func() tea.Msg {
		// Create API client and fetch user info
		client := NewClient(apiKey)
		userInfo, err := client.GetUserInfo()

		// If error, return error message
		if err != nil {
			return ErrMsg(err)
		}

		// Success - return user info message with workspace and user IDs
		return UserInfoMsg{
			WorkspaceID: userInfo.DefaultWorkspace,
			UserID:      userInfo.ID,
		}
	}
}

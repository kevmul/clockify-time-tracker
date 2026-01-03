// internal/ui/commands.go
// Wraps API calls into Bubble Tea commands
// Commands are functions that return messages - they bridge the API and UI layers
package ui

import (
	"time"

	"clockify-time-tracker/internal/api"

	tea "github.com/charmbracelet/bubbletea"
)

// fetchUserInfo returns a command that fetches user information
// When complete, it sends a userInfoMsg back to Update()
func fetchUserInfo(apiKey string) tea.Cmd {
	return func() tea.Msg {
		// Create API client and fetch user info
		client := api.NewClient(apiKey)
		userInfo, err := client.GetUserInfo()
		
		// If error, return error message
		if err != nil {
			return errMsg(err)
		}

		// Success - return user info message with workspace and user IDs
		return userInfoMsg{
			workspaceID: userInfo.DefaultWorkspace,
			userID:      userInfo.ID,
		}
	}
}

// fetchProjects returns a command that fetches all projects
// When complete, it sends a projectsMsg back to Update()
func fetchProjects(apiKey, workspaceID string) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient(apiKey)
		projects, err := client.GetProjects(workspaceID)
		
		if err != nil {
			return errMsg(err)
		}

		// Wrap the projects slice in projectsMsg type
		// This is crucial - it converts []api.Project to projectsMsg
		return projectsMsg(projects)
	}
}

// fetchTasks returns a command that fetches recent task descriptions
// When complete, it sends a tasksMsg back to Update()
func fetchTasks(apiKey, workspaceID, userID string) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient(apiKey)
		tasks, err := client.GetTasks(workspaceID, userID)
		
		if err != nil {
			return errMsg(err)
		}

		// Wrap the tasks slice in tasksMsg type
		return tasksMsg(tasks)
	}
}

// createTimeEntry returns a command that creates a time entry
// When complete, it sends either submitSuccessMsg or errMsg
func createTimeEntry(apiKey, workspaceID, projectID, description, timeRange string, date time.Time) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient(apiKey)
		err := client.CreateTimeEntry(workspaceID, projectID, description, timeRange, date)
		
		if err != nil {
			return errMsg(err)
		}

		// Success - return success message
		return submitSuccessMsg{}
	}
}

// internal/ui/commands.go
// Wraps API calls into Bubble Tea commands
// Commands are functions that return messages - they bridge the API and UI layers
package timeentry

import (
	"clockify-time-tracker/internal/api"
	"clockify-time-tracker/internal/debug"
	"clockify-time-tracker/internal/messages"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// fetchUserInfo returns a command that fetches user information
// When complete, it sends a userInfoMsg back to Update()
func fetchUserInfo(apiKey string) tea.Cmd {
	return func() tea.Msg {
		debug.Log("fetchUserInfo called with apiKey: %s...", apiKey[:10])
		
		// Create API client and fetch user info
		client := api.NewClient(apiKey)
		userInfo, err := client.GetUserInfo()

		// If error, return error message
		if err != nil {
			debug.Log("fetchUserInfo error: %v", err)
			return messages.ErrMsg(err)
		}

		// Debug: Check if we got valid workspace ID
		if userInfo.DefaultWorkspace == "" {
			debug.Log("user has no default workspace")
			return messages.ErrMsg(fmt.Errorf("user has no default workspace"))
		}

		debug.Log("fetchUserInfo success - WorkspaceID: %s, UserID: %s", 
			userInfo.DefaultWorkspace, userInfo.ID)

		// Success - return user info message with workspace and user IDs
		return messages.UserInfoMsg{
			WorkspaceID: userInfo.DefaultWorkspace,
			UserID:      userInfo.ID,
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
			return messages.ErrMsg(err)
		}

		// Debug: Add some logging to see what we got
		// TODO: Remove this debug code later
		if len(projects) == 0 {
			return messages.ErrMsg(fmt.Errorf("no projects found for workspace %s", workspaceID))
		}

		// Wrap the projects slice in projectsMsg type
		// This is crucial - it converts []api.Project to projectsMsg
		return messages.ProjectsMsg(projects)
	}
}

// fetchTasks returns a command that fetches recent task descriptions
// When complete, it sends a tasksMsg back to Update()
func fetchTasks(apiKey, workspaceID, userID string) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient(apiKey)
		tasks, err := client.GetTasks(workspaceID, userID)

		if err != nil {
			return messages.ErrMsg(err)
		}

		// Wrap the tasks slice in tasksMsg type
		return messages.TasksMsg(tasks)
	}
}

// createTimeEntry returns a command that creates a time entry
// When complete, it sends either submitSuccessMsg or errMsg
func createTimeEntry(apiKey, workspaceID, projectID, description, timeRange string, date time.Time) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient(apiKey)
		err := client.CreateTimeEntry(workspaceID, projectID, description, timeRange, date)

		if err != nil {
			return messages.ErrMsg(err)
		}

		// Success - return success message
		return messages.SubmitSuccessMsg{}
	}
}

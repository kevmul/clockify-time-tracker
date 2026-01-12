// internal/ui/commands.go
// Wraps API calls into Bubble Tea commands
// Commands are functions that return messages - they bridge the API and UI layers
package timeform

import (
	"clockify-time-tracker/internal/clockify"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// fetchUserInfo returns a command that fetches user information
// When complete, it sends a userInfoMsg back to Update()
func fetchUserInfo(apiKey string) tea.Cmd {
	return func() tea.Msg {

		// Create API client and fetch user info
		client := clockify.NewClient(apiKey)
		userInfo, err := client.GetUserInfo()

		// If error, return error message
		if err != nil {
			return clockify.ErrMsg(err)
		}

		if userInfo.DefaultWorkspace == "" {
			return clockify.ErrMsg(fmt.Errorf("user has no default workspace"))
		}

		// Success - return user info message with workspace and user IDs
		return clockify.UserInfoMsg{
			WorkspaceID: userInfo.DefaultWorkspace,
			UserID:      userInfo.ID,
		}
	}
}

// fetchProjects returns a command that fetches all projects
// When complete, it sends a projectsMsg back to Update()
func fetchProjects(apiKey, workspaceID string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(3 * time.Second)
		client := clockify.NewClient(apiKey)
		projects, err := client.GetProjects(workspaceID)

		if err != nil {
			return clockify.ErrMsg(err)
		}

		// TODO: Remove this debug code later
		if len(projects) == 0 {
			return clockify.ErrMsg(fmt.Errorf("no projects found for workspace %s", workspaceID))
		}

		// Wrap the projects slice in projectsMsg type
		// This is crucial - it converts []api.Project to projectsMsg
		return clockify.ProjectsMsg(projects)
	}
}

// fetchTasks returns a command that fetches recent task descriptions
// When complete, it sends a tasksMsg back to Update()
func fetchTasks(apiKey, workspaceID, userID string) tea.Cmd {
	return func() tea.Msg {
		client := clockify.NewClient(apiKey)
		tasks, err := client.GetTasks(workspaceID, userID)

		if err != nil {
			return clockify.ErrMsg(err)
		}

		// Wrap the tasks slice in tasksMsg type
		return clockify.TasksMsg(tasks)
	}
}

// createTimeEntry returns a command that creates a time entry
// When complete, it sends either submitSuccessMsg or errMsg
func createTimeEntry(apiKey, workspaceID, projectID, description, timeRange string, date time.Time) tea.Cmd {
	return func() tea.Msg {
		client := clockify.NewClient(apiKey)
		err := client.CreateTimeEntry(workspaceID, projectID, description, timeRange, date)

		if err != nil {
			return clockify.ErrMsg(err)
		}

		// Success - return success message
		return clockify.SubmitSuccessMsg{}
	}
}

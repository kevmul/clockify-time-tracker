package timelist

import (
	"clockify-time-tracker/internal/clockify"
	"clockify-time-tracker/internal/config"
	"clockify-time-tracker/internal/ui/styles"
	"clockify-time-tracker/internal/ui/views/timeform"
	debug "clockify-time-tracker/internal/utils"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	apiKey        string
	cursor        int
	entries       []clockify.Entry
	projects      []clockify.Project
	userID        string
	workspaceID   string
	loading       bool
	spinner       spinner.Model
	timeEntryForm timeform.Model
}

func New(cfg *config.Config) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.ColorPrimary)

	return Model{
		cursor:        0,
		apiKey:        cfg.APIKey,
		loading:       false,
		spinner:       s,
		timeEntryForm: timeform.New(cfg),
	}
}

func (m Model) Init() tea.Cmd {
	m.Reset()

	return tea.Batch(
		m.spinner.Tick,
		fetchUserInfo(m.apiKey),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	// Handle key inputs
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, func() tea.Msg {
				return clockify.QuittingAppMsg{}
			}
		// Up arrow or 'k' (vim style) - move cursor up in lists
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// Down arrow or 'j' (vim style) - move cursor down in lists
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}

		case "n":
			return m, func() tea.Msg {
				return clockify.CreateOrEditEntryMsg{
					Type: clockify.Create,
				}
			}
		}

	// Handle messages
	case clockify.UserInfoMsg:
		m.workspaceID = msg.WorkspaceID
		m.userID = msg.UserID
		cmds = append(cmds, func() tea.Msg {
			return clockify.SetLoadingMsg{}
		})

	case clockify.SetLoadingMsg:
		m.loading = true
		cmds = append(
			cmds,
			m.getEntries,
			func() tea.Msg { return m.fetchProjects },
		)

	case clockify.EntriesMsg:
		// time.Sleep(2 * time.Second)
		m.loading = false
		m.entries = msg.Entries
		m.projects = msg.Projects
	}

	return m, tea.Batch(cmds...)
}

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

func (m Model) fetchEntries() ([]clockify.Entry, error) {
	client := clockify.NewClient(m.apiKey)
	endpoint := fmt.Sprintf("/workspaces/%s/user/%s/time-entries", m.workspaceID, m.userID)

	body, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	// Parse the response into slice of time entries
	var entries []clockify.Entry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse time entries: %w", err)
	}

	return entries, nil
}

func (m Model) getEntries() tea.Msg {
	entries, err := m.fetchEntries()

	if err != nil {
		return clockify.ErrMsg(err)
	}

	projects, perr := m.fetchProjects()
	if perr != nil {
		return clockify.ErrMsg(perr)
	}

	return clockify.EntriesMsg{
		Entries:  entries,
		Projects: projects,
	}
}

func (m Model) Reset() Model {
	debug.Log("Resetting TimeList Model")
	m.entries = nil
	return m
}

// fetchProjects returns a command that fetches all projects
// When complete, it sends a projectsMsg back to Update()
func (m Model) fetchProjects() ([]clockify.Project, error) {
	client := clockify.NewClient(m.apiKey)
	projects, err := client.GetProjects(m.workspaceID)

	if err != nil {
		return nil, err
	}

	// Wrap the projects slice in projectsMsg type
	// This is crucial - it converts []api.Project to projectsMsg
	return projects, nil
}

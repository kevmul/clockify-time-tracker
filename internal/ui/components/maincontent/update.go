// Handles all state updates in response to messages
// This is the heart of the Bubble Tea architecture
package maincontent

import (
	"clockify-time-tracker/internal/api"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Message types that can be sent to Update()
// These are custom types that wrap the actual data
type projectsMsg []api.Project // List of projects from API
type tasksMsg []string         // List of task descriptions
type userInfoMsg struct {      // User info from API
	workspaceID string
	userID      string
}
type errMsg error              // Error that occurred
type submitSuccessMsg struct{} // Empty struct signals success

// Update is called whenever a message is received
// It's the only place where we modify the Model
// Returns the updated Model and any commands to run
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle text input FIRST before checking message types
	// This ensures text inputs get all key events
	if m.step == stepTimeInput || m.step == stepTaskInput || (m.step == stepProjectSelect && m.projectSearch.Focused()) {
		var cmd tea.Cmd
		if m.step == stepTimeInput {
			m.timeRange, cmd = m.timeRange.Update(msg)
		} else if m.step == stepTaskInput {
			m.taskName, cmd = m.taskName.Update(msg)
		} else if m.step == stepProjectSelect && m.projectSearch.Focused() {
			m.projectSearch, cmd = m.projectSearch.Update(msg)
			// Reset cursor when search changes
			m.cursor = 0
		}

		// Still check for special keys like Enter and quit keys
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			// case "ctrl+c", "q":
			// return m, tea.Quit
			case "enter":
				return m.handleEnter()
			case "esc":
				if m.step == stepProjectSelect && m.projectSearch.Focused() {
					m.projectSearch.Blur()
					m.projectSearch.SetValue("")
					m.cursor = 0
				}
			}
		}

		return m, cmd
	}

	// Check what type of message we received
	switch msg := msg.(type) {

	// Keyboard input from the user
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	// User info was fetched successfully
	case userInfoMsg:
		m.workspaceID = msg.workspaceID
		m.userID = msg.userID
		// Now fetch projects and tasks in parallel using tea.Batch
		return m, tea.Batch(
			fetchProjects(m.apiKey, m.workspaceID),
			fetchTasks(m.apiKey, m.workspaceID, m.userID),
		)

	// Projects were fetched successfully
	case projectsMsg:
		m.projects = msg
		return m, nil

	// Tasks were fetched successfully
	case tasksMsg:
		m.tasks = msg
		return m, nil

	// An error occurred
	case errMsg:
		m.err = msg
		return m, tea.Quit // Quit the program on error

	// Time entry was created successfully
	case submitSuccessMsg:
		m.success = true
		m.step = stepComplete
		return m, tea.Quit // Quit after success

	// Window was resized (we don't handle this yet)
	case tea.WindowSizeMsg:
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes all keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	// Quit keys - always available
	case "ctrl+c":
		return m, tea.Quit

	// Quit keys - always available unless in search
	case "q":
		if m.step == stepProjectSelect && m.projectSearch.Focused() {
			return m, nil
		}
		return m, tea.Quit

	case "t":
		if m.step == stepDateSelect {
			m.date = time.Now() // Default to today
		}

	// Up arrow or 'k' (vim style) - move cursor up in lists
	case "up", "k":
		if m.step == stepProjectSelect && !m.projectSearch.Focused() && m.cursor > 0 {
			filteredProjects := m.filterProjects()
			if m.cursor < len(filteredProjects) {
				m.cursor--
			}
		}

	// Down arrow or 'j' (vim style) - move cursor down in lists
	case "down", "j":
		if m.step == stepProjectSelect && !m.projectSearch.Focused() {
			filteredProjects := m.filterProjects()
			if m.cursor < len(filteredProjects)-1 {
				m.cursor++
			}
		}

	// Forward slash - focus search
	case "/":
		if m.step == stepProjectSelect {
			m.projectSearch.Focus()
			return m, textinput.Blink
		}

	// Left arrow or 'h' (vim style) - previous day
	case "left", "h":
		if m.step == stepDateSelect {
			m.date = m.date.AddDate(0, 0, -1)
		}

	// Right arrow or 'l' (vim style) - next day
	case "right", "l":
		if m.step == stepDateSelect {
			m.date = m.date.AddDate(0, 0, 1)
		}

	// Enter key - confirm current step and move to next
	case "enter":
		return m.handleEnter()
	}

	return m, nil
}

// handleEnter processes the Enter key - advances to next step
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {

	// Date selected - move to project selection
	case stepDateSelect:
		m.step = stepProjectSelect
		m.cursor = 0 // Reset cursor to top of project list
		return m, nil

	// Project selected - move to time input
	case stepProjectSelect:
		if m.projectSearch.Focused() {
			m.projectSearch.Blur()
			return m, nil
		}
		filteredProjects := m.filterProjects()
		if len(filteredProjects) > 0 && m.cursor < len(filteredProjects) {
			m.selectedProj = filteredProjects[m.cursor] // Save selected project
			m.step = stepTimeInput
			m.timeRange.Focus()       // Focus the time input field
			return m, textinput.Blink // Start cursor blinking in text input
		}

	// Time entered - move to task input
	case stepTimeInput:
		if m.timeRange.Value() != "" { // Only proceed if they entered something
			m.step = stepTaskInput
			m.timeRange.Blur() // Unfocus the time input
			m.taskName.Focus() // Focus the task input field
			return m, textinput.Blink
		}

	// Task entered - move to confirmation
	case stepTaskInput:
		if m.taskName.Value() != "" { // Only proceed if they entered something
			m.step = stepConfirm
		}

	// Confirmed - submit the entry
	case stepConfirm:
		return m, m.submitTimeEntry()
	}

	return m, nil
}

// handleTextInput updates text input fields when user types
func (m Model) handleTextInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Update the appropriate text input based on current step
	if m.step == stepTimeInput {
		m.timeRange, cmd = m.timeRange.Update(msg)
		return m, cmd
	}

	if m.step == stepTaskInput {
		m.taskName, cmd = m.taskName.Update(msg)
		return m, cmd
	}

	return m, nil
}

// submitTimeEntry creates a command to submit the time entry
func (m Model) submitTimeEntry() tea.Cmd {
	return createTimeEntry(
		m.apiKey,
		m.workspaceID,
		m.selectedProj.ID,
		m.taskName.Value(),
		m.timeRange.Value(),
		m.date,
	)
}

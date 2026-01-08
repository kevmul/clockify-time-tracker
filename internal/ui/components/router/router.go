package router

import (
	timeentry "clockify-time-tracker/internal/ui/components/timeEntry"
	"clockify-time-tracker/internal/ui/messages"
	"clockify-time-tracker/internal/ui/styles"
	"clockify-time-tracker/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	currentView string

	// Each view has its own component
	timeEntry timeentry.Model
	// dashboard dashboard.Model
	// projects projects.Model
}

func (m Model) View() string {
	var content string

	switch m.currentView {
	case "Time Entry":
		content = m.timeEntry.View()
	default:
		content = "View not found: " + m.currentView
	}

	return styles.MainContentStyle.Render(content)
}

func New(cfg *utils.Config) Model {
	return Model{
		currentView: "Dashboard",
		timeEntry:   timeentry.New(cfg),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle view changes
	switch msg := msg.(type) {
	case messages.NavigationMsg:
		m.currentView = msg.Item
		return m, nil
	}

	// Route messages to the active component
	switch m.currentView {
	case "Time Entry":
		timeEntryModel, cmd := m.timeEntry.Update(msg)
		m.timeEntry = timeEntryModel.(timeentry.Model)

		return m, cmd
	}

	return m, cmd
}

func (m Model) SetView(view messages.NavigationMsg) Model {
	m.currentView = view.Item
	return m
}

package app

import (
	"clockify-time-tracker/internal/config"
	"clockify-time-tracker/internal/messages"
	"clockify-time-tracker/internal/ui/styles"
	timeentry "clockify-time-tracker/internal/ui/views/timeform"

	tea "github.com/charmbracelet/bubbletea"
)

type RouterModel struct {
	currentView string

	// Each view has its own component
	timeEntry timeentry.Model
	// dashboard dashboard.Model
	// projects projects.Model
}

func (m RouterModel) View() string {
	var content string

	switch m.currentView {
	case "Time Entry":
		content = m.timeEntry.View()
	default:
		content = "View not found: " + m.currentView
	}

	return styles.MainContentStyle.Render(content)
}

func NewRouter(cfg *config.Config) RouterModel {
	return RouterModel{
		currentView: "Dashboard",
		timeEntry:   timeentry.New(cfg),
	}
}

func (m RouterModel) Init() tea.Cmd {
	// Don't initialize any view by default - only when user navigates
	return nil
}

func (m RouterModel) Update(msg tea.Msg) (RouterModel, tea.Cmd) {
	var cmd tea.Cmd

	// Handle view changes
	switch msg := msg.(type) {
	case messages.NavigationMsg:
		oldView := m.currentView
		m.currentView = msg.Item

		if oldView != m.currentView {
			cmd := m.initCurrentView()
			return m, cmd
		}
		return m, nil
	}

	// Route UI messages to the active component
	switch m.currentView {
	case "Time Entry":
		m.timeEntry, cmd = m.timeEntry.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m RouterModel) initCurrentView() tea.Cmd {

	switch m.currentView {
	case "Time Entry":
		return m.timeEntry.Init()
	default:
		return nil
	}
}

package router

import (
	"clockify-time-tracker/internal/debug"
	"clockify-time-tracker/internal/messages"
	timeentry "clockify-time-tracker/internal/ui/components/timeEntry"
	"clockify-time-tracker/internal/ui/styles"
	"clockify-time-tracker/internal/utils"
	"fmt"

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

func (m Model) Init() tea.Cmd {
	// Don't initialize any view by default - only when user navigates
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle view changes
	switch msg := msg.(type) {
	case messages.NavigationMsg:
		debug.Log("Router received NavigationMsg: %s", msg.Item)
		oldView := m.currentView
		m.currentView = msg.Item

		fmt.Printf("Old: %s | New: %s", oldView, m.currentView)

		if oldView != m.currentView {
			debug.Log("Router calling initCurrentView() for: %s", m.currentView)
			cmd := m.initCurrentView()
			debug.Log("Router initCurrentView() returned command: %v", cmd != nil)
			return m, cmd
		}
		return m, nil
	}

	// Route UI messages to the active component
	switch m.currentView {
	case "Time Entry":
		timeEntryModel, cmd := m.timeEntry.Update(msg)
		m.timeEntry = timeEntryModel.(timeentry.Model)

		return m, cmd
	}

	return m, cmd
}

func (m Model) initCurrentView() tea.Cmd {
	debug.Log("initCurrentView called for: %s", m.currentView)
	switch m.currentView {
	case "Time Entry":
		debug.Log("Calling timeEntry.Init()")
		return m.timeEntry.Init()
	default:
		debug.Log("No init needed for view: %s", m.currentView)
		return nil
	}
}

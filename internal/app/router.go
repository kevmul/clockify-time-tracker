package app

import (
	"clockify-time-tracker/internal/clockify"
	"clockify-time-tracker/internal/config"
	"clockify-time-tracker/internal/ui/views/dashboard"
	timeentry "clockify-time-tracker/internal/ui/views/timeform"
	"clockify-time-tracker/internal/ui/views/timelist"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type RouterModel struct {
	currentView clockify.ViewState

	// Each view has its own component
	dashboard dashboard.Model
	timeList  timelist.Model
	timeEntry timeentry.Model
	// projects projects.Model

}

func (m RouterModel) View() string {
	var content string

	switch m.currentView {
	case clockify.ViewDashboard:
		content = m.dashboard.View()
	case clockify.ViewTimeEntry:
		content = m.timeEntry.View()
	case clockify.ViewTimeList:
		content = m.timeList.View()
	default:
		content = fmt.Sprintf("View not found: %s", m.currentView)
	}

	// width, height, _ := term.GetSize(os.Stdout.Fd())

	return content
	// return mainContentStyle.Render(content)
}

func NewRouter(cfg *config.Config) RouterModel {
	return RouterModel{
		currentView: clockify.ViewDashboard,
		timeEntry:   timeentry.New(cfg),
		timeList:    timelist.New(cfg),
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
	case clockify.NavigationMsg:
		oldView := m.currentView
		m.currentView = msg.View

		if oldView != m.currentView {
			cmd := m.initCurrentView()
			return m, cmd
		}
		return m, nil
	}

	// Route UI messages to the active component
	switch m.currentView {
	case clockify.ViewDashboard:
		return m, nil
	case clockify.ViewTimeEntry:
		m.timeEntry, cmd = m.timeEntry.Update(msg)
		return m, cmd
	case clockify.ViewTimeList:
		m.timeList, cmd = m.timeList.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m RouterModel) initCurrentView() tea.Cmd {

	switch m.currentView {
	case clockify.ViewDashboard:
		return nil
		// return m.dashboard.Init()
	case clockify.ViewTimeEntry:
		return m.timeEntry.Init()
	case clockify.ViewTimeList:
		return m.timeList.Init()
	default:
		return nil
	}
}

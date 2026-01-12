package app

import (
	"clockify-time-tracker/internal/clockify"
	"clockify-time-tracker/internal/config"
	"clockify-time-tracker/internal/ui/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	sidebar     SidebarModel
	router      RouterModel
	focusedPane string // "sidebar" or "content"
	quitting    bool
}

func New(config *config.Config) Model {
	items := []string{"Dashboard", "Time Entry", "Reports"}

	return Model{
		sidebar:     NewSidebar(items, 10),
		router:      NewRouter(config),
		focusedPane: "sidebar", // Start with the sidebar focused
		quitting:    false,
	}
}

func (m Model) Init() tea.Cmd {
	return m.router.Init()
	// return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// case "ctrl+c", "q":
		// 	m.quitting = true
		// 	return m, tea.Quit
		case "tab": // We can change this later
			if m.focusedPane == "sidebar" {
				m.focusedPane = "content"
			} else {
				m.focusedPane = "sidebar"
			}
			return m, nil
		case "esc":
			// Always return to the sidebar on Escape
			m.focusedPane = "sidebar"
			return m, nil
		}
	case clockify.NavigationMsg:
		// Handle the navigation by sending to router
		var cmd tea.Cmd
		m.router, cmd = m.router.Update(msg)
		cmds = append(cmds, cmd)

		if msg.View == clockify.ViewTimeList {
			m.focusedPane = "content"
		}
		return m, tea.Batch(cmds...)
	case clockify.QuittingAppMsg:
		m.quitting = true
		return m, tea.Quit
	}

	// Route UI messages based on focus
	if m.focusedPane == "sidebar" {
		m.sidebar, cmd = m.sidebar.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.router, cmd = m.router.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	// Clear the screen when quitting the app
	if m.quitting {
		// return "\n"
	}

	sidebarView := m.sidebar.View()
	if m.focusedPane == "sidebar" {
		sidebarView = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Render(sidebarView)
	} else {
		sidebarView = lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder()).
			Render(sidebarView)
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebarView,
		m.router.View(),
	)

	title := styles.TitleStyle.Render("⏱️  Clockify Time Tracker")
	app := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
	)

	help := fmt.Sprintf("\nFocused %s | [Tab] to switch focus | [Esc] to return to Sidebar | [q] to quit", m.focusedPane)

	return app + help
}

package ui

import (
	"clockify-time-tracker/internal/ui/components/router"
	"clockify-time-tracker/internal/ui/components/sidebar"
	"clockify-time-tracker/internal/ui/messages"
	"clockify-time-tracker/internal/utils"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	sidebar     sidebar.Model
	router      router.Model
	focusedPane string // "sidebar" or "content"
}

func New(config *utils.Config) Model {
	items := []string{"Dashboard", "Time Entry", "Reports"}

	return Model{
		sidebar:     sidebar.New(items, 10),
		router:      router.New(config),
		focusedPane: "sidebar", // Start with the sidebar focused
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
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
	case messages.NavigationMsg:
		// Handle the navigation
		m.router = m.router.SetView(msg)
		// Could trigger other actions here like loading data

		if msg.Item == "Time Entry" {
			m.focusedPane = "content"
		}
		return m, nil
	}

	// Route messages based on focus
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
		m.sidebar.View(),
		m.router.View(),
	)

	app := content

	help := fmt.Sprintf("\nFocused %s | [Tab] to swithc focus | [Esc] to return to Sidebar | [q] to quit", m.focusedPane)

	return app + help
}

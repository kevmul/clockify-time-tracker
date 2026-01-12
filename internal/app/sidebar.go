package app

import (
	"clockify-time-tracker/internal/clockify"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	sidebarStyle = lipgloss.NewStyle().
			Padding(1, 1).
			Width(20)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Bold(true)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

type SidebarModel struct {
	items       []string
	cursor      int
	viewportTop int
	height      int
}

func (m SidebarModel) Init() tea.Cmd {
	return nil
}

func NewSidebar(items []string, height int) SidebarModel {
	return SidebarModel{
		items:  items,
		height: height,
	}
}

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				// Scroll viewport if needed
				if m.cursor < m.viewportTop {
					m.viewportTop = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
				// Scroll viewport if needed
				if m.cursor >= m.viewportTop+m.height {
					m.viewportTop = m.cursor - m.height + 1
				}
			}
		case "enter":
			// Send a message when something is selected
			return m, func() tea.Msg {
				var view clockify.ViewState
				switch m.SelectedItem() {
				case "Dashboard":
					view = clockify.ViewDashboard
				case "Time Entry":
					view = clockify.ViewTimeList
				case "Reports":
					view = clockify.ViewSettings
				default:
					view = clockify.ViewDashboard
				}
				return clockify.NavigationMsg{
					View: view,
				}
			}
		}
	}
	return m, nil
}

type ViewChangedMsg struct {
	View string
}

func (m SidebarModel) SelectedItem() string {
	if m.cursor < len(m.items) {
		return m.items[m.cursor]
	}
	return ""
}

func (m SidebarModel) SelectedIndex() int {
	return m.cursor
}

func (m SidebarModel) View() string {
	var content string
	visibleEnd := m.viewportTop + m.height
	if visibleEnd > len(m.items) {
		visibleEnd = len(m.items)
	}

	for i := m.viewportTop; i < visibleEnd; i++ {
		cursor := " "
		style := normalItemStyle
		if i == m.cursor {
			cursor = "»"
			style = selectedItemStyle
		}

		content += fmt.Sprintf("%s %s\n", cursor, style.Render(m.items[i]))
	}

	return sidebarStyle.
		Height(m.height + 2).
		Render(content)
}

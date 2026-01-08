package sidebar

import (
	"fmt"

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

func (m Model) View() string {
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

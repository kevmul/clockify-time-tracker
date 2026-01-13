package dialog

import (
	"clockify-time-tracker/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

func (m Modal) View(termWidth, termHeight int) string {
	if !m.Visible {
		return ""
	}

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Padding(1, 2).
		Width(m.Width).
		Height(m.Height)

	contentStyle := lipgloss.NewStyle().
		Width(m.Width - 4).
		Height(m.Height - 4)

	title := styles.TitleStyle.Render(m.Title)
	content := contentStyle.Render(m.Content)
	modal := modalStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, content))

	// Center the modal on the screen
	return modal
}

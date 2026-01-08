package ui

import (
	"clockify-time-tracker/internal/ui/components/maincontent"
	"clockify-time-tracker/internal/ui/components/sidebar"
	"clockify-time-tracker/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	sidebar sidebar.Model
	content maincontent.Model
}

func New(config *utils.Config) Model {
	items := []string{"Dashboard", "Projects", "Tasks", "Settings"}

	return Model{
		sidebar: sidebar.New(items, 10),
		content: maincontent.New(config),
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
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	m.sidebar, cmd = m.sidebar.Update(msg)
	cmds = append(cmds, cmd)

	contentModel, cmd := m.content.Update(m.sidebar.SelectedItem())
	m.content = contentModel.(maincontent.Model)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	appWrapperStyle := lipgloss.NewStyle()
	// Border(lipgloss.RoundedBorder()).
	// BorderForeground(styles.ColorBorder)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.sidebar.View(),
		m.content.View(),
	)

	app := appWrapperStyle.Render(content)

	return app
}

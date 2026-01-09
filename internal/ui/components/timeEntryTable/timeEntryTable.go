package timeentrytable

import (
	"clockify-time-tracker/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	apiKey string
}

func New(config *utils.Config) Model {
	return Model{
		apiKey: config.APIKey,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		fetchUserInfo(m.apiKey),
	)
}

package sidebar

import (
	"clockify-time-tracker/internal/messages"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	items       []string
	cursor      int
	viewportTop int
	height      int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New(items []string, height int) Model {
	return Model{
		items:  items,
		height: height,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
				return messages.NavigationMsg{
					Item:  m.SelectedItem(),
					Index: m.SelectedIndex(),
				}
			}
		}
	}
	return m, nil
}

type ViewChangedMsg struct {
	View string
}

func (m Model) SelectedItem() string {
	if m.cursor < len(m.items) {
		return m.items[m.cursor]
	}
	return ""
}

func (m Model) SelectedIndex() int {
	return m.cursor
}

package app

import (
	"clockify-time-tracker/internal/ui/styles"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	Help key.Binding
	Quit key.Binding
	Up   key.Binding
	Down key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "show help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
	}
}

func GenerateHelpContent(keys KeyMap) string {
	var helpText strings.Builder

	helpItems := []struct {
		keys string
		desc string
	}{
		{keys.Help.Help().Key, keys.Help.Help().Desc},
		{keys.Up.Help().Key, keys.Up.Help().Desc},
		{keys.Down.Help().Key, keys.Down.Help().Desc},
		{keys.Quit.Help().Key, keys.Quit.Help().Desc},
		{"esc", "close modal"},
	}

	keyStyle := lipgloss.NewStyle().
		Foreground(styles.ColorHeaderText).
		Bold(true).
		Width(12)

	descStyle := lipgloss.NewStyle().
		Foreground(styles.ColorMuted)

	for _, item := range helpItems {
		helpText.WriteString(keyStyle.Render(item.keys))
		helpText.WriteString(descStyle.Render(item.desc))
		helpText.WriteString("\n")
	}

	return helpText.String()
}

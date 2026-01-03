// Defines all visual styles for the terminal UI using lipgloss
package ui

import "github.com/charmbracelet/lipgloss"

// Styles for different UI elements
// These are package-level variables so they can be used throughout the ui package

var (
	// titleStyle is used for the main app title at the top
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")). // Pink/magenta color
			MarginBottom(1)

	// selectedStyle highlights the currently selected item
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")). // Purple color
			Bold(true)

	// errorStyle is used for error messages
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // Red color
			Bold(true)

	// successStyle is used for success messages
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")). // Green color
			Bold(true)
)

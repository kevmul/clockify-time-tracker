// Defines all visual styles for the terminal UI using lipgloss
package styles

import "github.com/charmbracelet/lipgloss"

// Styles for different UI elements
// These are package-level variables so they can be used throughout the ui package

var (
	ColorPrimary   = lipgloss.Color("170")
	ColorSecondary = lipgloss.Color("63")
	ColorMuted     = lipgloss.Color("240")
	ColorBorder    = lipgloss.Color("62")

	ColorError   = lipgloss.Color("196")
	ColorSuccess = lipgloss.Color("42")

	ColorHeaderText = lipgloss.Color("205")
)

var (
	// titleStyle is used for the main app title at the top
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorHeaderText). // Pink/magenta color
			MarginBottom(1)

	// selectedStyle highlights the currently selected item
	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary). // Purple color
			Bold(true)

	// errorStyle is used for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError). // Red color
			Bold(true)

	// successStyle is used for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess). // Green color
			Bold(true)
)

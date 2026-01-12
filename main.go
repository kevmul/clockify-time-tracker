// Entry point for the Clockify Time Tracker CLI application
package main

import (
	"clockify-time-tracker/internal/app"
	"clockify-time-tracker/internal/config"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load configuration from .env file
	// This will read CLOCKIFY_API_KEY and return an error if not found
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create a new Bubble Tea program with our UI model
	// The ui.New() function initializes the model with our config
	p := tea.NewProgram(app.New(config), tea.WithAltScreen())

	// Run the program - this starts the interactive TUI
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

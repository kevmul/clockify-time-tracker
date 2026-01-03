// Entry point for the Clockify Time Tracker CLI application
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	
	"clockify-time-tracker/internal/ui"
	"clockify-time-tracker/internal/utils"
)

func main() {
	// Load configuration from .env file
	// This will read CLOCKIFY_API_KEY and return an error if not found
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create a new Bubble Tea program with our UI model
	// The ui.New() function initializes the model with our config
	p := tea.NewProgram(ui.New(config))
	
	// Run the program - this starts the interactive TUI
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

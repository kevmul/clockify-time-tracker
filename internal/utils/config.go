// Handles loading and validating configuration from environment variables
package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
// We only need the API key at startup - workspace and user IDs are fetched later
type Config struct {
	APIKey string
}

// LoadConfig loads environment variables from .env file and validates them
// Returns a Config struct or an error if required variables are missing
func LoadConfig() (*Config, error) {
	// Load .env file - ignore error if file doesn't exist (e.g., in production)
	// The underscore _ means we're intentionally ignoring the return value
	homedir, _ := os.UserHomeDir()
	_ = godotenv.Load("./.env")                                   // try local .env after. Will not override.
	_ = godotenv.Load(homedir + "/.config/clockify-tracker/.env") // Load config env first

	// Get the API key from environment
	apiKey := os.Getenv("CLOCKIFY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("CLOCKIFY_API_KEY not set in .env file")
	}

	// Return the config struct
	return &Config{
		APIKey: apiKey,
	}, nil
}

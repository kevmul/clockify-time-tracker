package timeentry

import (
	"clockify-time-tracker/internal/api"
	"clockify-time-tracker/internal/ui/styles"
	"clockify-time-tracker/internal/utils"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// model represents the entire state of our application
// This is the "single source of truth" for what's currently happening in the UI
type Model struct {
	// Current step in the workflow (which screen we're on)
	step int

	// Data from API
	projects []api.Project // List of available projects
	tasks    []string      // Recent task descriptions for suggestions

	// Navigation state
	cursor   int // Current position in lists (for arrow key navigation)
	selected int // Index of selected item (not currently used but kept for future)

	// User inputs
	date          time.Time       // Selected date for time entry
	timeRange     textinput.Model // Text input for time range (e.g., "9a - 5p")
	taskName      textinput.Model // Text input for task description
	projectSearch textinput.Model // Text input for project search
	selectedProj  api.Project     // The project user selected

	// API credentials and IDs
	apiKey      string // Clockify API key
	workspaceID string // User's workspace ID (fetched from API)
	userID      string // User's ID (fetched from API)

	// Status flags
	err        error // Any error that occurred
	submitting bool  // Whether we're currently submitting (not used yet)
	success    bool  // Whether submission was successful

	// Loading state
	spinner spinner.Model
}

// New creates and initializes a new model with the provided configuration
// This is called from main.go when the app starts
func New(config *utils.Config) Model {
	// Create and configure the time range text input
	ti := textinput.New()
	ti.Placeholder = "9a - 5p" // Show example format
	ti.CharLimit = 20          // Reasonable limit for time strings
	ti.Width = 30

	// Create and configure the task name text input
	taskInput := textinput.New()
	taskInput.Placeholder = "Enter task description"
	taskInput.CharLimit = 200 // Clockify's description limit
	taskInput.Width = 50

	// Create and configure the project search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search projects..."
	searchInput.Width = 50

	// Loading state
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.ColorPrimary)

	// Return a new Model with initial state
	return Model{
		step:          stepDateSelect, // Start at date selection
		date:          time.Now(),     // Default to today
		timeRange:     ti,
		taskName:      taskInput,
		projectSearch: searchInput,
		cursor:        0,             // Start at first item in lists
		apiKey:        config.APIKey, // Store API key from config
		spinner:       s,
	}
}

// Init is called once when the program starts
// It returns a command that will fetch the user's info from Clockify
// This is part of the Bubble Tea architecture - Init returns initial commands to run
func (m Model) Init() tea.Cmd {
	// Fetch user info (workspace ID and user ID) as our first action
	return tea.Batch(
		m.spinner.Tick,
		fetchUserInfo(m.apiKey),
	)
}

# Clockify Time Tracker CLI

A beautiful, interactive command-line tool for logging time entries to Clockify.

## Features

- 📅 Interactive date selection (defaults to today)
- 📂 Project selection with arrow key navigation
- ⏰ Simple time range input (e.g., "9a - 5p")
- 📝 Task description with suggestions from your previous entries
- ✨ Clean, colorful terminal UI using Bubble Tea

## Project Structure

```
clockify-time-tracker/
├── main.go                           # Entry point
├── go.mod                            # Go module definition
├── .env                              # Your API key (not in git)
├── .env.example                      # Template for .env
├── .gitignore
├── README.md
│
└── internal/
    ├── api/                          # API layer - talks to Clockify
    │   ├── client.go                 # HTTP client for API requests
    │   ├── types.go                  # Data structures (Project, TimeEntry, etc.)
    │   ├── user.go                   # User-related API calls
    │   ├── projects.go               # Project-related API calls
    │   └── timeentries.go            # Time entry API calls
    │
    ├── ui/                           # UI layer - Bubble Tea components
    │   ├── model.go                  # Application state
    │   ├── update.go                 # State updates (handles messages)
    │   ├── view.go                   # Rendering (displays UI)
    │   ├── commands.go               # Wraps API calls as Bubble Tea commands
    │   ├── steps.go                  # Step/screen constants
    │   └── styles.go                 # Visual styles (colors, formatting)
    │
    └── utils/                        # Utilities
        └── config.go                 # Configuration loading
```

## How It Works

### Architecture Overview

This project follows the **Bubble Tea** (Elm Architecture) pattern with clear separation of concerns:

1. **API Layer** (`internal/api/`):
   - Pure Go functions that make HTTP requests to Clockify
   - Returns data or errors - no UI logic
   - Reusable and testable

2. **UI Layer** (`internal/ui/`):
   - **Model**: Holds all application state
   - **Update**: Receives messages, updates state, returns commands
   - **View**: Renders the current state as a string
   - **Commands**: Bridge between API and UI - async operations that send messages

3. **Flow**:
   ```
   User Input → Update() → New State → View() → Display
                    ↓
                Commands (API calls)
                    ↓
                Messages → Update() → ...
   ```

### Key Concepts

**Messages**: Data sent to `Update()` to trigger state changes

- `userInfoMsg` - User info fetched from API
- `projectsMsg` - Projects fetched from API
- `tasksMsg` - Recent tasks fetched from API
- `errMsg` - An error occurred
- `submitSuccessMsg` - Time entry created successfully

**Commands**: Functions that return messages asynchronously

- Wrap API calls
- Run in background
- Send results back as messages

**Steps**: Different screens in the UI flow

- Date selection → Project selection → Time input → Task input → Confirm → Complete

## Prerequisites

- Go 1.21 or higher
- A Clockify account and API key

## Installation

1. Clone this repository:

```bash
git clone <your-repo-url>
cd clockify-time-tracker
```

2. Install dependencies:

```bash
go mod download
```

3. Create your `.env` file:

```bash
cp .env.example .env
```

4. Edit `.env` and add your Clockify API key:

```
CLOCKIFY_API_KEY=your_api_key_here
```

## Getting Your Clockify API Key

1. Log in to [Clockify](https://clockify.me)
2. Go to Settings → Profile
3. Scroll down to "API" section
4. Generate or copy your API key

## Usage

Run the tool:

```bash
go run main.go
```

Or build and run:

```bash
go build -o clockify-tracker
./clockify-tracker
```

### Navigation

- **Date Selection**: Use `←`/`→` arrow keys to change dates, `Enter` to confirm
- **Project Selection**: Use `↑`/`↓` arrow keys to navigate, `Enter` to select
- **Time Range**: Type in format like `9a - 5p` or `9:30a - 3:45p`
- **Task Description**: Type your task description
- **Quit**: Press `q` or `Ctrl+C` at any time

### Time Format Examples

- `9a - 5p` → 9:00 AM to 5:00 PM
- `9:30a - 3:45p` → 9:30 AM to 3:45 PM
- `10a - 2p` → 10:00 AM to 2:00 PM

## Building for Distribution

Build for your platform:

```bash
# macOS
GOOS=darwin GOARCH=amd64 go build -o clockify-tracker-mac

# Linux
GOOS=linux GOARCH=amd64 go build -o clockify-tracker-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o clockify-tracker.exe
```

## Code Walkthrough

### Adding a New Feature

Let's say you want to add a "billable" toggle:

1. **Add to model** (`ui/model.go`):

   ```go
   type model struct {
       // ... existing fields
       billable bool
   }
   ```

2. **Add a step** (`ui/steps.go`):

   ```go
   const (
       // ... existing steps
       stepBillableSelect
   )
   ```

3. **Handle in Update** (`ui/update.go`):

   ```go
   case stepBillableSelect:
       // Handle billable selection
   ```

4. **Render in View** (`ui/view.go`):

   ```go
   case stepBillableSelect:
       return m.renderBillableSelect()
   ```

5. **Update API call** (`api/timeentries.go`):
   ```go
   type TimeEntryRequest struct {
       // ... existing fields
       Billable bool `json:"billable"`
   }
   ```

## Troubleshooting

**"Error loading .env file"**

- Make sure you've created a `.env` file in the project root
- Copy `.env.example` to `.env` and fill in your API key

**"CLOCKIFY_API_KEY not set"**

- Check that your `.env` file contains `CLOCKIFY_API_KEY=your_key`
- Make sure there are no extra spaces around the `=`

**"Loading projects..." never finishes**

- Check your internet connection
- Verify your API key is correct
- Make sure you have projects in your Clockify workspace

**"Failed to create entry"**

- Verify your API key is correct
- Check that the time format is valid (e.g., `9a - 5p`)
- Ensure you have access to the selected project

## Contributing

Feel free to open issues or submit pull requests!

## License

MIT License - feel free to use this for personal or commercial projects.

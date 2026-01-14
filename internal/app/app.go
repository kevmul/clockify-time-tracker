package app

import (
	"clockify-time-tracker/internal/clockify"
	"clockify-time-tracker/internal/config"
	"clockify-time-tracker/internal/ui/components/dialog"
	"clockify-time-tracker/internal/ui/styles"
	"clockify-time-tracker/internal/ui/views/timeform"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type Model struct {
	width  int
	height int

	sidebar SidebarModel
	router  RouterModel

	// Help Modal
	keys      KeyMap
	helpModal dialog.Modal

	// Time entry Modal
	timeEntryModal dialog.Modal
	timeEntryForm  timeform.Model

	focusedPane string // "sidebar" or "content"
	quitting    bool
}

func New(config *config.Config) Model {
	items := []string{"Dashboard", "Time Entries", "Reports"}

	keys := DefaultKeyMap()

	return Model{
		width:       80,
		height:      24,
		sidebar:     NewSidebar(items, 10),
		router:      NewRouter(config),
		focusedPane: "sidebar", // Start with the sidebar focused
		quitting:    false,

		keys:      keys,
		helpModal: dialog.NewModal("Keyboard Shortcuts"),

		timeEntryModal: dialog.NewModal("Create a New Time Entry"),
		timeEntryForm:  timeform.New(config),
	}
}

func (m Model) Init() tea.Cmd {
	return m.router.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:

		if m.helpModal.Visible {
			switch msg.String() {
			case "esc", "?":
				m.helpModal.Hide()
				return m, nil
			}
			return m, nil
		}

		if m.timeEntryModal.Visible {
			switch msg.String() {
			case "esc":
				m.timeEntryModal.Hide()
				return m, nil
			default:
				// Pass all keys to the form
				m.timeEntryForm, cmd = m.timeEntryForm.Update(msg)

				m.timeEntryForm.Init()

				// Update the modal content with the form
				m.timeEntryModal.Content = m.timeEntryForm.View()

				return m, cmd
			}
		}

		switch {
		case key.Matches(msg, m.keys.Help):
			m.helpModal.Show(GenerateHelpContent(m.keys))
			return m, nil
		}
		switch msg.String() {
		// We will remove this for Msg in the future.
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "tab": // We can change this later
			if m.focusedPane == "sidebar" {
				m.focusedPane = "content"
			} else {
				m.focusedPane = "sidebar"
			}
			return m, nil
		case "esc":
			// Always return to the sidebar on Escape
			m.focusedPane = "sidebar"
			return m, nil
		}
	case clockify.NavigationMsg:
		// Handle the navigation by sending to router
		var cmd tea.Cmd
		m.router, cmd = m.router.Update(msg)
		cmds = append(cmds, cmd)

		if msg.View == clockify.ViewTimeList {
			m.focusedPane = "content"
		}
		return m, tea.Batch(cmds...)

	case clockify.CreateOrEditEntryMsg:
		if msg.Type == clockify.Edit {
			m.timeEntryModal = dialog.NewModal("Editing time stamp")
		}
		m.timeEntryModal.Show(m.timeEntryForm.View())
		m.timeEntryForm, cmd = m.timeEntryForm.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case clockify.QuittingAppMsg:
		m.quitting = true
		return m, tea.Quit
	}

	// Route UI messages based on focus
	if m.focusedPane == "sidebar" {
		m.sidebar, cmd = m.sidebar.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.router, cmd = m.router.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	// Clear the screen when quitting the app
	if m.quitting {
		return "\n"
	}

	sidebarView := m.sidebar.View()
	mainContentView := m.router.View()
	if m.focusedPane == "sidebar" {
		sidebarView = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Render(sidebarView)

		mainContentView = lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder()).
			Padding(1, 2).
			Render(mainContentView)
	} else {
		sidebarView = lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder()).
			Render(sidebarView)

		width, height, _ := term.GetSize(os.Stdout.Fd())

		mainContentView = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorder).
			Padding(1, 2).
			Width(width - 24). // 20 is the sidebar width
			Height(height - 100).
			Render(mainContentView)
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebarView,
		mainContentView,
	)

	title := styles.TitleStyle.Render("⏱️  Clockify Time Tracker")
	app := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
	)
	help := fmt.Sprintf("\nFocused %s | [Tab] to switch focus | [Esc] to return to Sidebar | [q] to quit", m.focusedPane)

	if m.helpModal.Visible {
		return m.renderWithModal(app+help, m.helpModal)
	}

	if m.timeEntryModal.Visible {
		return m.renderWithModal(app+help, m.timeEntryModal)
	}

	return app + help
}

func (m Model) renderWithModal(baseContent string, modal dialog.Modal) string {
	// Split base content into lines
	baseLines := strings.Split(baseContent, "\n")

	// Ensure we have enough lines for the terminal height
	for len(baseLines) < m.height {
		baseLines = append(baseLines, "")
	}

	// Render the modal
	modalContent := modal.View(m.width, m.height)
	modalLines := strings.Split(modalContent, "\n")

	// Calculate starting position to center the modal
	modalHeight := len(modalLines)
	startRow := (m.height - modalHeight) / 2
	if startRow < 0 {
		startRow = 0
	}

	// Find the actual width of the modal
	modalWidth := 0
	for _, line := range modalLines {
		lineLen := lipgloss.Width(line)
		if lineLen > modalWidth {
			modalWidth = lineLen
		}
	}

	startCol := (m.width - modalWidth) / 2
	if startCol < 0 {
		startCol = 0
	}

	// Helper to truncate string at visual width (ANSI-aware)
	truncateAt := func(s string, width int) string {
		if width <= 0 {
			return ""
		}
		var result strings.Builder
		currentWidth := 0
		inEscape := false

		for _, r := range s {
			if r == '\x1b' {
				inEscape = true
			}

			if inEscape {
				result.WriteRune(r)
				if r == 'm' {
					inEscape = false
				}
				continue
			}

			if currentWidth >= width {
				break
			}

			result.WriteRune(r)
			currentWidth++
		}
		return result.String()
	}

	// Helper to skip first N visual characters (ANSI-aware)
	skipChars := func(s string, n int) string {
		if n <= 0 {
			return s
		}

		skipped := 0
		inEscape := false
		var result strings.Builder
		started := false

		for _, r := range s {
			if r == '\x1b' {
				inEscape = true
			}

			if started || inEscape {
				result.WriteRune(r)
			}

			if inEscape {
				if r == 'm' {
					inEscape = false
				}
				continue
			}

			if !started {
				skipped++
				if skipped > n {
					started = true
					result.WriteRune(r)
				}
			}
		}
		return result.String()
	}

	// Overlay modal lines onto base lines
	for i, modalLine := range modalLines {
		row := startRow + i
		if row >= 0 && row < len(baseLines) {
			baseLine := baseLines[row]
			baseWidth := lipgloss.Width(baseLine)

			// Extract left part (before modal)
			leftPart := truncateAt(baseLine, startCol)

			// Extract right part (after modal)
			endCol := startCol + lipgloss.Width(modalLine)
			var rightPart string
			if endCol < baseWidth {
				rightPart = skipChars(baseLine, endCol)
			}

			// Pad if needed
			leftWidth := lipgloss.Width(leftPart)
			if leftWidth < startCol {
				leftPart += strings.Repeat(" ", startCol-leftWidth)
			}

			baseLines[row] = leftPart + modalLine + rightPart
		}
	}

	return strings.Join(baseLines, "\n")
}

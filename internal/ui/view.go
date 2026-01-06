// Renders the UI based on current model state
// This is called every time the state changes
package ui

import (
	"clockify-time-tracker/internal/api"
	"fmt"
	"strings"
)

// View returns a string representation of the UI
// Bubble Tea calls this function to render the screen
func (m model) View() string {
	// Handle error state - show error and quit
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("\n❌ Error: %v\n", m.err))
	}

	// Handle success state - show success message
	if m.step == stepComplete && m.success {
		return successStyle.Render("\n✅ Time entry created successfully!\n\n")
	}

	// Start building the UI string
	// We use a string builder for efficiency
	s := titleStyle.Render("⏱️  Clockify Time Tracker") + "\n\n"

	// Render different content based on current step
	switch m.step {
	case stepDateSelect:
		s += m.renderDateSelect()
	case stepProjectSelect:
		s += m.renderProjectSelect()
	case stepTimeInput:
		s += m.renderTimeInput()
	case stepTaskInput:
		s += m.renderTaskInput()
	case stepConfirm:
		s += m.renderConfirm()
	}

	return s
}

// renderDateSelect shows the date selection screen
func (m model) renderDateSelect() string {
	s := "Select date (use ←/→ to change, Enter to confirm):\n\n"
	s += fmt.Sprintf("  📅 %s\n\n", m.date.Format("Monday, January 2, 2006"))
	s += "  [Enter] Confirm  [←/→] Change day  [q] Quit\n"
	return s
}

// renderProjectSelect shows the project selection list
func (m model) renderProjectSelect() string {
	// If no projects loaded yet, show loading message
	if len(m.projects) == 0 {
		return "Loading projects...\n"
	}

	var sb strings.Builder

	sb.WriteString("Select a project:\n\n")

	// Show search input
	sb.WriteString("🔍 " + m.projectSearch.View() + "\n\n")

	// Filter projects based on search
	filteredProjects := m.filterProjects()

	if len(filteredProjects) == 0 {
		sb.WriteString("  No projects match your search.\n\n")
		sb.WriteString("  [/] Search  [Esc] Clear search  [q] Quit\n")
		return sb.String()
	}

	// Calculate visible range for scrolling
	const visibleItems = 10 // Show 10 items at a time
	start := 0
	end := len(filteredProjects)

	// If we have more projects than can fit, show a window around cursor
	if len(filteredProjects) > visibleItems {
		// Center the cursor in the window
		start = m.cursor - visibleItems/2
		end = start + visibleItems

		// Adjust if we're near the beginning
		if start < 0 {
			start = 0
			end = visibleItems
		}

		// Adjust if we're near the end
		if end > len(filteredProjects) {
			end = len(filteredProjects)
			start = end - visibleItems
			if start < 0 {
				start = 0
			}
		}

		// Show indicator if there are items above
		if start > 0 {
			sb.WriteString(fmt.Sprintf("  ↑ %d more above...\n", start))
		}
	}

	// Show visible projects
	for i := start; i < end; i++ {
		proj := filteredProjects[i]

		// Format project name with client if available
		displayName := proj.Name
		if proj.ClientName != "" {
			displayName = fmt.Sprintf("%s (%s)", proj.Name, proj.ClientName)
		}

		if m.cursor == i {
			// This is the selected item -= highlight it
			sb.WriteString(selectedStyle.Render(fmt.Sprintf("❯ %s", displayName)) + "\n")
		} else {
			// Regular rendering for unselected items
			sb.WriteString(fmt.Sprintf("  %s\n", displayName))
		}
	}

	// Show indicator if there are items below
	if len(filteredProjects) > visibleItems && end < len(filteredProjects) {
		sb.WriteString(fmt.Sprintf("  ↓ %d more below...\n", len(filteredProjects)-end))
	}

	sb.WriteString("\n  [↑/↓] Navigate  [Enter] Select  [/] Search  [Esc] Clear  [q] Quit\n")
	return sb.String()
}

// filterProjects returns projects that match the current search query
func (m model) filterProjects() []api.Project {
	query := strings.ToLower(strings.TrimSpace(m.projectSearch.Value()))
	if query == "" {
		return m.projects
	}

	var filtered []api.Project
	for _, proj := range m.projects {
		if strings.Contains(strings.ToLower(proj.Name), query) {
			filtered = append(filtered, proj)
		}
	}
	return filtered
}

// renderTimeInput shows the time range input field
func (m model) renderTimeInput() string {
	s := fmt.Sprintf("Project: %s\n", selectedStyle.Render(m.selectedProj.Name))
	s += fmt.Sprintf("Date: %s\n\n", m.date.Format("Jan 2, 2006"))
	s += "Enter time range (e.g., 9a - 5p):\n\n"
	s += m.timeRange.View() // Render the text input
	s += "\n\n  [Enter] Continue  [q] Quit\n"
	return s
}

// renderTaskInput shows the task description input field
func (m model) renderTaskInput() string {
	s := fmt.Sprintf("Project: %s\n", selectedStyle.Render(m.selectedProj.Name))
	s += fmt.Sprintf("Date: %s\n", m.date.Format("Jan 2, 2006"))
	s += fmt.Sprintf("Time: %s\n\n", m.timeRange.Value())
	s += "Enter task description:\n\n"
	s += m.taskName.View() // Render the text input

	// Show recent tasks as suggestions if available
	if len(m.tasks) > 0 {
		// Show up to 3 recent tasks
		recentCount := min(3, len(m.tasks))
		recentTasks := m.tasks[:recentCount]
		s += "\n\n  Recent tasks: " + strings.Join(recentTasks, ", ")
	}

	s += "\n\n  [Enter] Continue  [q] Quit\n"
	return s
}

// renderConfirm shows the confirmation screen with all entered details
func (m model) renderConfirm() string {
	s := "Confirm time entry:\n\n"
	s += fmt.Sprintf("  Project: %s\n", selectedStyle.Render(m.selectedProj.Name))
	s += fmt.Sprintf("  Date: %s\n", m.date.Format("Jan 2, 2006"))
	s += fmt.Sprintf("  Time: %s\n", m.timeRange.Value())
	s += fmt.Sprintf("  Task: %s\n\n", m.taskName.Value())
	s += "  [Enter] Submit  [q] Cancel\n"
	return s
}

// min returns the smaller of two integers
// Helper function used for limiting the number of recent tasks shown
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

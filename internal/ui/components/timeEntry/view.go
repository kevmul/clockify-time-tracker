// Renders the styles.based on current Model state
// This is called every time the state changes
package timeentry

import (
	"clockify-time-tracker/internal/api"
	"clockify-time-tracker/internal/ui/styles"
	"fmt"
	"strings"
)

// View returns a string representation of the UI
// Bubble Tea calls this function to render the screen
func (m Model) View() string {
	// Handle error state - show error and quit
	if m.err != nil {
		return styles.ErrorStyle.Render(fmt.Sprintf("\n❌ Error: %v\n", m.err))
	}

	// Handle success state - show success message
	if m.step == stepComplete && m.success {
		return styles.SuccessStyle.Render("\n✅ Time entry created successfully!\n\n")
	}

	// Start building the UI string
	// We use a string builder for efficiency
	s := styles.TitleStyle.Render("⏱️  Clockify Time Tracker") + "\n\n"

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
	// return mainContentStyle.Render(s)
}

// renderDateSelect shows the date selection screen
func (m Model) renderDateSelect() string {
	s := "Select date (use ←/→ to change, Enter to confirm):\n\n"
	s += fmt.Sprintf("  📅 %s\n\n", m.date.Format("Monday, January 2, 2006"))
	s += fmt.Sprintf("  %s %s %s %s\n", EnterPrompt, VerticalArrowPrompt, TodayPrompt, QuitPrompt)
	return s
}

// renderProjectSelect shows the project selection list
func (m Model) renderProjectSelect() string {
	// If no projects loaded yet, show loading message
	if len(m.projects) == 0 {
		str := fmt.Sprintf("%s Loading...", m.spinner.View())
		// return "Loading projects...\n"
		return str
	}

	var sb strings.Builder

	sb.WriteString("Select a project:\n\n")

	// Show search input
	sb.WriteString("🔍 " + m.projectSearch.View() + "\n\n")

	// Filter projects based on search
	filteredProjects := m.filterProjects()

	if len(filteredProjects) == 0 {
		sb.WriteString("  No projects match your search.\n\n")
		sb.WriteString(fmt.Sprintf("  %s %s %s", SearchPrompt, ClearPrompt, QuitPrompt))
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
			sb.WriteString(styles.SelectedStyle.Render(fmt.Sprintf("❯ %s", displayName)) + "\n")
		} else {
			// Regular rendering for unselected items
			sb.WriteString(fmt.Sprintf("  %s\n", displayName))
		}
	}

	// Show indicator if there are items below
	if len(filteredProjects) > visibleItems && end < len(filteredProjects) {
		sb.WriteString(fmt.Sprintf("  ↓ %d more below...\n", len(filteredProjects)-end))
	}

	sb.WriteString(fmt.Sprintf("\n  %s %s %s %s %s", ArrowPrompt, EnterPrompt, SearchPrompt, ClearPrompt, QuitPrompt))
	return sb.String()
}

// filterProjects returns projects that match the current search query
func (m Model) filterProjects() []api.Project {
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
func (m Model) renderTimeInput() string {
	s := fmt.Sprintf("Project: %s\n", styles.SelectedStyle.Render(m.selectedProj.Name))
	s += fmt.Sprintf("Date: %s\n\n", m.date.Format("Jan 2, 2006"))
	s += "Enter time range (e.g., 9a - 5p):\n\n"
	s += m.timeRange.View() // Render the text input
	s += fmt.Sprintf("\n\n  %s %s", EnterPrompt, QuitPrompt)
	return s
}

// renderTaskInput shows the task description input field
func (m Model) renderTaskInput() string {
	s := fmt.Sprintf("Project: %s\n", styles.SelectedStyle.Render(m.selectedProj.Name))
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

	s += fmt.Sprintf("\n\n  %s %s", EnterPrompt, QuitPrompt)
	return s
}

// renderConfirm shows the confirmation screen with all entered details
func (m Model) renderConfirm() string {
	s := "Confirm time entry:\n\n"
	s += fmt.Sprintf("  Project: %s\n", styles.SelectedStyle.Render(m.selectedProj.Name))
	s += fmt.Sprintf("  Date: %s\n", m.date.Format("Jan 2, 2006"))
	s += fmt.Sprintf("  Time: %s\n", m.timeRange.Value())
	s += fmt.Sprintf("  Task: %s\n\n", m.taskName.Value())
	s += fmt.Sprintf("  %s %s", EnterPrompt, QuitPrompt)
	return s
}

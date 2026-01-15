package timelist

import (
	"clockify-time-tracker/internal/ui/styles"
	"fmt"
	"time"
)

func (m Model) View() string {
	s := ""

	if m.loading {
		s += fmt.Sprintf("%s Loading...", m.spinner.View())
	}

	const visibleItems = 5
	start := 0
	end := len(m.entries)

	if len(m.entries) > visibleItems {
		start = m.cursor - visibleItems/2
		end = start + visibleItems

		if start < 0 {
			start = 0
			end = visibleItems
		}

		if end > len(m.entries) {
			end = len(m.entries)
			start = end - visibleItems
			if start < 0 {
				start = 0
			}
		}

		if start > 0 {
			s += fmt.Sprintf(" ↑ %d more above...\n", start)
		} else {
			s += "\n" // Make space for "More Above" indicator
		}
	}

	for i := start; i < end; i++ {
		entry := m.entries[i]
		st, sterr := time.Parse(time.RFC3339, entry.TimeInterval.Start)
		et, eterr := time.Parse(time.RFC3339, entry.TimeInterval.End)
		startTime := entry.TimeInterval.Start
		endTime := entry.TimeInterval.End
		if sterr == nil {
			startTime = st.In(time.Local).Format("3:04PM")
		}
		if eterr == nil {
			endTime = et.In(time.Local).Format("3:04PM")
		}

		description := entry.Description
		if description == "" {
			description = "??"
		}

		projectName := "No Project"
		for _, proj := range m.projects {
			if proj.ID == entry.ProjectId {
				projectName = fmt.Sprintf("%s (%s)", proj.Name, proj.ClientName)
				break
			}
		}

		displayTime := fmt.Sprintf("%s\n  %s (%s - %s)", projectName, description, startTime, endTime)

		if m.cursor == i {
			s += styles.SelectedListStyle.Width(60).Render(fmt.Sprintf("» %s", displayTime)) + "\n"
		} else {
			s += styles.NormalListStyle.Width(60).Render(fmt.Sprintf("  %s", displayTime)) + "\n"
		}
	}

	if len(m.entries) > visibleItems && end < len(m.entries) {
		s += fmt.Sprintf("  ↓ %d more below...\n", len(m.entries)-end)
	} else {
		s += "\n"
	}

	return s
}

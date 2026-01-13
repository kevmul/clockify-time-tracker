package timelist

import (
	"clockify-time-tracker/internal/ui/styles"
	"fmt"
	"time"
)

func (m Model) View() string {
	s := ""

	s += fmt.Sprintf("Cursor Position: %d\n\n", m.cursor)

	if m.loading {
		s += fmt.Sprintf("%s Loading...", m.spinner.View())
	}

	const visibleItems = 10
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
		}
	}

	for i := start; i < end; i++ {
		entry := m.entries[i]
		st, sterr := time.Parse(time.RFC3339, entry.TimeInterval.Start)
		et, eterr := time.Parse(time.RFC3339, entry.TimeInterval.End)
		startTime := entry.TimeInterval.Start
		endTime := entry.TimeInterval.End
		if sterr == nil {
			startTime = st.Format("3:04 PM")
		}
		if eterr == nil {
			endTime = et.Format("3:04 PM")
		}

		displayName := fmt.Sprintf("%s (%s - %s)", entry.Description, startTime, endTime)

		if m.cursor == i {
			s += styles.SelectedStyle.Render(fmt.Sprintf("» %s", displayName)) + "\n"
		} else {
			s += styles.SelectedStyle.Render(fmt.Sprintf("  %s", displayName)) + "\n"
		}
	}

	if len(m.entries) > visibleItems && end < len(m.entries) {
		s += fmt.Sprintf("  ↓ %d more below...\n", len(m.entries)-end)
	}

	return s
}

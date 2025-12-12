package tui

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	var s strings.Builder

	s.WriteString("noctisyn\n\n")

	if len(m.list.Items()) == 0 {
		if m.searching {
			s.WriteString("searching...\n")
		} else {
			s.WriteString(m.input.View())
			s.WriteString("\n\n")
			if m.err != nil {
				s.WriteString(fmt.Sprintf("error: %v\n", m.err))
			} else {
				s.WriteString("type a query and press enter to search.\n")
			}
		}
	} else {
		s.WriteString(m.list.View())
		s.WriteString("\n\npress enter to download • esc to search again • ctrl+c to quit\n")
	}

	if len(m.jobs) > 0 {
		s.WriteString("\ndownloads:\n")
		for _, job := range m.jobs {
			var status string
			switch job.status {
			case statusQueued:
				status = "queued"
			case statusFetching:
				status = "fetching"
			case statusWriting:
				status = "writing"
			case statusDone:
				status = fmt.Sprintf("done -> %s", job.outPath)
			case statusFailed:
				status = fmt.Sprintf("failed: %v", job.err)
			}

			title := job.title
			if len(title) > 40 {
				title = title[:37] + "..."
			}
			s.WriteString(fmt.Sprintf("  [%s] %s\n", status, title))
		}
	}

	return s.String()
}

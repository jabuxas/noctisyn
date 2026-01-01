package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if !m.ready {
		return "\ninitializing..."
	}

	downloadsPanelWidth := max(int(float64(m.width)*0.25), 30)
	mainPanelWidth := m.width - downloadsPanelWidth + 2

	mainPanel := m.renderMainPanel(mainPanelWidth)

	downloadsPanel := m.renderDownloadsPanel(downloadsPanelWidth, m.height-2)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		mainPanel,
		downloadsPanel,
	)

	return content
}

func (m model) renderMainPanel(width int) string {
	var s strings.Builder

	title := titleStyle.Width(width - 4).Render("noctisyn")
	s.WriteString(title)
	s.WriteString("\n\n")

	if len(m.list.Items()) == 0 {
		if m.searching {
			searchingText := lipgloss.NewStyle().
				Foreground(secondaryColor).
				Render("Searching...")
			s.WriteString(searchingText)
		} else {
			s.WriteString(m.input.View())
			s.WriteString("\n\n")

			if m.err != nil {
				errorText := lipgloss.NewStyle().
					Foreground(errorColor).
					Render(fmt.Sprintf("error: %v", m.err))
				s.WriteString(errorText)
			} else {
				helpText := helpStyle.Render("type a query and press enter to search")
				s.WriteString(helpText)
			}
		}
	} else {
		s.WriteString(m.list.View())
		s.WriteString("\n\n")

		help := helpStyle.Render("enter: download • esc: new search • ctrl+c: quit")
		s.WriteString(help)
	}

	return mainPanelStyle.
		Width(width - 4).
		Height(m.height - 2).
		Render(s.String())
}

func (m model) renderDownloadsPanel(width int, height int) string {
	var s strings.Builder

	header := downloadsHeaderStyle.
		Width(width - 4).
		Render("downloads")
	s.WriteString(header)
	s.WriteString("\n")

	if len(m.jobs) == 0 {
		emptyText := helpStyle.Render("no active downloads")
		s.WriteString(emptyText)
	} else {
		maxJobsToShow := max((height-8)/3, 1)

		startIdx := 0
		if len(m.jobs) > maxJobsToShow {
			startIdx = len(m.jobs) - maxJobsToShow
		}

		for i := startIdx; i < len(m.jobs); i++ {
			job := m.jobs[i]
			s.WriteString(m.renderJob(job, width-6))
			if i < len(m.jobs)-1 {
				s.WriteString("\n")
				divider := dividerStyle.Render(strings.Repeat("─", width-6))
				s.WriteString(divider)
				s.WriteString("\n")
			}
		}

		if startIdx > 0 {
			s.WriteString("\n")
			moreText := helpStyle.Render(fmt.Sprintf("... and %d more", startIdx))
			s.WriteString(moreText)
		}
	}

	return downloadsPanelStyle.
		Width(width - 2).
		Height(height).
		Render(s.String())
}

func (m model) renderJob(job downloadJob, width int) string {
	var s strings.Builder

	title := job.title
	if len(title) > width-10 {
		title = title[:width-13] + "..."
	}

	var statusText string
	switch job.status {
	case statusQueued:
		statusText = statusQueuedStyle.Render("⋯  queued")
	case statusFetching:
		if job.totalChapters > 0 {
			percentage := int(float64(job.currentChapter) / float64(job.totalChapters) * 100)
			progressLine := fmt.Sprintf("↓ fetching: %d/%d (%d%%)", job.currentChapter, job.totalChapters, percentage)
			statusText = statusFetchingStyle.Render(progressLine)
		} else {
			statusText = statusFetchingStyle.Render("↓ fetching")
		}
	case statusWriting:
		statusText = statusFetchingStyle.Render("✎ writing")
	case statusDone:
		statusText = statusDoneStyle.Render("✓ done")
	case statusFailed:
		statusText = statusFailedStyle.Render("✗ failed")
	}

	s.WriteString(lipgloss.NewStyle().Bold(true).Render(title))
	s.WriteString("\n")
	s.WriteString(statusText)

	if job.status == statusFetching && job.estimatedTimeMs > 0 {
		s.WriteString("\n")
		estimatedSec := job.estimatedTimeMs / 1000
		timeText := helpStyle.Render(fmt.Sprintf("→ estimated: %ds remaining", estimatedSec))
		s.WriteString(timeText)
	}

	if job.status == statusDone && job.outPath != "" {
		s.WriteString("\n")
		pathText := helpStyle.Render(fmt.Sprintf("→ %s", job.outPath))
		s.WriteString(pathText)
	}

	if job.status == statusFailed && job.err != nil {
		s.WriteString("\n")
		errText := statusFailedStyle.Render(fmt.Sprintf("→ %v", job.err))
		s.WriteString(errText)
	}

	return s.String()
}

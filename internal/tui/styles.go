package tui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor   = lipgloss.Color("255")
	secondaryColor = lipgloss.Color("250")
	successColor   = lipgloss.Color("248")
	errorColor     = lipgloss.Color("240")
	warningColor   = lipgloss.Color("244")
	mutedColor     = lipgloss.Color("237")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1)

	mainPanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(1, 2)

	downloadsPanelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor).
				Padding(1, 1)

	downloadsHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(secondaryColor).
				BorderStyle(lipgloss.Border{Bottom: "─"}).
				BorderBottom(true).
				BorderForeground(mutedColor).
				MarginBottom(1).
				Width(100)

	statusQueuedStyle = lipgloss.NewStyle().
				Foreground(warningColor)

	statusFetchingStyle = lipgloss.NewStyle().
				Foreground(secondaryColor)

	statusDoneStyle = lipgloss.NewStyle().
			Foreground(successColor)

	statusFailedStyle = lipgloss.NewStyle().
				Foreground(errorColor)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	dividerStyle = lipgloss.NewStyle().
			Foreground(mutedColor)
)

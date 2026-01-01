package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type jobStatus string

const (
	statusQueued   jobStatus = "queued"
	statusFetching jobStatus = "fetching"
	statusWriting  jobStatus = "writing"
	statusDone     jobStatus = "done"
	statusFailed   jobStatus = "failed"
)

type downloadJob struct {
	id              int
	title           string
	url             string
	status          jobStatus
	err             error
	outPath         string
	currentChapter  int
	totalChapters   int
	estimatedTimeMs int64
	sub             chan tea.Msg
}

type model struct {
	input     textinput.Model
	list      list.Model
	jobs      []downloadJob
	nextJobID int
	searching bool
	err       error
	width     int
	height    int
	ready     bool
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "search novel"
	ti.Focus()
	ti.Prompt = "> "
	ti.CharLimit = 64
	ti.Width = 50

	delegate := list.NewDefaultDelegate()

	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("255")).
		BorderForeground(lipgloss.Color("250"))

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("248")).
		BorderForeground(lipgloss.Color("250"))

	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("250"))

	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(lipgloss.Color("240"))

	l := list.New([]list.Item{}, delegate, 0, 0)

	l.Styles.Title = l.Styles.Title.
		Foreground(lipgloss.Color("255")).
		Bold(true)

	l.Title = "search results"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return model{
		input:     ti,
		list:      l,
		jobs:      []downloadJob{},
		nextJobID: 1,
	}
}

package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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
	id      int
	title   string
	url     string
	status  jobStatus
	err     error
	outPath string
}

type model struct {
	input     textinput.Model
	list      list.Model
	jobs      []downloadJob
	nextJobID int
	searching bool
	err       error
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "search novel"
	ti.Focus()
	ti.Prompt = "> "
	ti.CharLimit = 64

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
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

package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jabuxas/noctisyn/internal/scraper"
)

type model struct {
	input   textinput.Model
	loading bool
	result  string
	err     error
}

// message types
type searchMsg struct {
	url string
	err error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "search novel"
	ti.Focus()
	ti.Prompt = "> "
	ti.CharLimit = 64

	return model{
		input: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			q := m.input.Value()
			if q == "" {
				return m, nil
			}
			m.loading = true
			m.result = ""
			m.err = nil
			return m, doSearch(q)
		}
	case searchMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.result = msg.url
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := "noctisyn\n\n"
	s += m.input.View() + "\n\n"

	if m.loading {
		s += "searching...\n"
	} else if m.err != nil {
		s += fmt.Sprintf("error: %v\n", m.err)
	} else if m.result != "" {
		s += fmt.Sprintf("match URL: %s\n", m.result)
		s += "press enter again to search another, or ctrl+c to quit.\n"
	} else {
		s += "type a query and press enter.\n"
	}

	return s
}

func doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		url, err := scraper.Search(query)
		return searchMsg{url: url, err: err}
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

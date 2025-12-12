package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jabuxas/noctisyn/internal/scraper"
)

type model struct {
	input    textinput.Model
	loading  bool
	results  []*scraper.Book
	selected int
	err      error
}

// message types
type searchMsg struct {
	books []*scraper.Book
	err   error
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
			if len(m.results) > 0 {
				return m, nil
			}
			q := m.input.Value()
			if q == "" {
				return m, nil
			}
			m.loading = true
			m.err = nil
			m.results = nil
			m.selected = 0
			return m, doSearch(q)
		case tea.KeyUp:
			if len(m.results) > 0 && m.selected > 0 {
				m.selected--
			}
		case tea.KeyDown:
			if len(m.results) > 0 && m.selected < len(m.results)-1 {
				m.selected++
			}
		}
	case searchMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.results = msg.books
		m.selected = 0
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
	} else if len(m.results) > 0 {
		s += "results (↑/↓ to move, enter to download):\n\n"
		for i, u := range m.results {
			cursor := " "
			if i == m.selected {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, u.Title)
		}
	} else {
		s += "type a query and press enter.\n"
	}

	return s
}

func doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		url, err := scraper.Search(query)
		return searchMsg{books: url, err: err}
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

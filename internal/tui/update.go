package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		if len(m.list.Items()) > 0 {
			switch msg.Type {
			case tea.KeyEnter:
				item := m.list.SelectedItem().(bookItem)

				sub := make(chan tea.Msg)
				job := downloadJob{
					id:     m.nextJobID,
					title:  item.book.Title,
					url:    item.book.SourceURL,
					status: statusQueued,
					sub:    sub,
				}
				m.jobs = append(m.jobs, job)
				m.nextJobID++

				return m, doDownload(job.id, item.book, sub)

			case tea.KeyEsc:
				m.list.SetItems([]list.Item{})
				m.input.Focus()
				m.input.SetValue("")
				return m, nil
			}

			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

		if msg.Type == tea.KeyEnter {
			q := m.input.Value()
			if q == "" {
				return m, nil
			}
			m.searching = true
			m.err = nil
			return m, doSearch(q)
		}

	case searchMsg:
		m.searching = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		items := make([]list.Item, len(msg.books))
		for i, b := range msg.books {
			items[i] = bookItem{book: b}
		}
		m.list.SetItems(items)
		m.input.Blur()
		return m, nil

	case downloadProgressMsg:
		for i := range m.jobs {
			if m.jobs[i].id == msg.jobID {
				m.jobs[i].currentChapter = msg.currentChapter
				m.jobs[i].totalChapters = msg.totalChapters
				m.jobs[i].estimatedTimeMs = msg.estimatedTimeMs
				if m.jobs[i].status == statusQueued {
					m.jobs[i].status = statusFetching
				}
				return m, waitForMsg(m.jobs[i].sub)
			}
		}
		return m, nil

	case downloadDoneMsg:
		for i := range m.jobs {
			if m.jobs[i].id == msg.jobID {
				if msg.err != nil {
					m.jobs[i].status = statusFailed
					m.jobs[i].err = msg.err
				} else {
					m.jobs[i].status = statusDone
					m.jobs[i].outPath = msg.outPath
				}
				break
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		downloadsPanelWidth := max(int(float64(m.width)*0.25), 30)

		mainPanelWidth := m.width - downloadsPanelWidth - 4

		listHeight := max(m.height-12, 5)

		m.list.SetSize(mainPanelWidth-4, listHeight)
		return m, nil
	}

	if len(m.list.Items()) == 0 {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

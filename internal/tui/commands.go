package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jabuxas/noctisyn/internal/epubgen"
	"github.com/jabuxas/noctisyn/internal/scraper"
)

func doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		books, err := scraper.Search(query)
		return searchMsg{books: books, err: err}
	}
}

func doDownload(jobID int, book *scraper.Book) tea.Cmd {
	return func() tea.Msg {
		novel, err := scraper.Fetch(book.SourceURL)
		if err != nil {
			return downloadDoneMsg{jobID: jobID, err: err}
		}

		filename := safeFilename(novel.Title) + ".epub"
		err = epubgen.WriteEPUB(novel, filename)
		return downloadDoneMsg{jobID: jobID, err: err, outPath: filename}
	}
}

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

func doDownload(jobID int, book *scraper.Book, sub chan tea.Msg) tea.Cmd {
	go func() {
		defer close(sub)

		novel, err := scraper.FetchWithProgress(book.SourceURL, func(current, total int, estimatedTimeMs int64) {
			sub <- downloadProgressMsg{
				jobID:           jobID,
				currentChapter:  current,
				totalChapters:   total,
				estimatedTimeMs: estimatedTimeMs,
			}
		})

		if err != nil {
			sub <- downloadDoneMsg{jobID: jobID, err: err}
			return
		}

		filename := safeFilename(novel.Title) + ".epub"
		err = epubgen.WriteEPUB(novel, filename)
		sub <- downloadDoneMsg{jobID: jobID, err: err, outPath: filename}
	}()

	return waitForMsg(sub)
}

func waitForMsg(sub chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

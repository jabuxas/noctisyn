package tui

import "github.com/jabuxas/noctisyn/internal/scraper"

type searchMsg struct {
	books []*scraper.Book
	err   error
}

type downloadStartedMsg struct {
	jobID int
}

type downloadProgressMsg struct {
	jobID           int
	currentChapter  int
	totalChapters   int
	estimatedTimeMs int64
}

type downloadDoneMsg struct {
	jobID   int
	err     error
	outPath string
}

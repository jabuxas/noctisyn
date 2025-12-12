package tui

import "github.com/jabuxas/noctisyn/internal/scraper"

type searchMsg struct {
	books []*scraper.Book
	err   error
}

type downloadStartedMsg struct {
	jobID int
}

type downloadDoneMsg struct {
	jobID   int
	err     error
	outPath string
}

package tui

import "github.com/jabuxas/noctisyn/internal/scraper"

type bookItem struct {
	book *scraper.Book
}

func (i bookItem) Title() string       { return i.book.Title }
func (i bookItem) Description() string { return i.book.SourceURL }
func (i bookItem) FilterValue() string { return i.book.Title }

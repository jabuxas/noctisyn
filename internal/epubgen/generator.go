package epubgen

import (
	"fmt"

	epub "github.com/go-shiori/go-epub"
	"github.com/jabuxas/noctisyn/internal/scraper"
)

func WriteEPUB(book *scraper.Book, outPath string) error {
	e, err := epub.NewEpub(book.Title)
	if err != nil {
		fmt.Println(err.Error())
	}
	e.SetAuthor(book.Author)
	e.SetIdentifier(book.SourceURL)

	// cssPath, _ := e.AddCSS("style.css", "style.css")

	for _, ch := range book.Chapters {
		body := fmt.Sprintf("<h1>%s</h1>\n<p>%s</p>", ch.Title, ch.Text)
		internal := fmt.Sprintf("ch-%04d.xhtml", ch.Index)
		if _, err := e.AddSection(body, ch.Title, internal, ""); err != nil {
			return err
		}
	}

	return e.Write(outPath)
}

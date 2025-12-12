package epubgen

import (
	"fmt"

	epub "github.com/go-shiori/go-epub"
	"github.com/jabuxas/noctisyn/internal/scraper"
)

func WriteEPUB(book *scraper.Book, outPath string) error {
	e, err := epub.NewEpub(book.Title)
	if err != nil {
		return err
	}
	e.SetAuthor(book.Author)
	e.SetDescription(book.Description)
	e.SetLang("en_US")
	e.SetIdentifier(book.SourceURL)
	coverImgPath, err := e.AddImage(book.CoverURL, "cover.jpg")
	if err != nil {
		return err
	}

	if err := e.SetCover(coverImgPath, ""); err != nil {
		return err
	}

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

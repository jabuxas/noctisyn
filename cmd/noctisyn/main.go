package main

import (
	"fmt"

	"github.com/jabuxas/noctisyn/internal/epubgen"
	"github.com/jabuxas/noctisyn/internal/scraper"
)

func main() {
	url, err := scraper.Search("warlock magus")
	if err != nil {
		fmt.Println(err.Error())
	}
	book, err := scraper.Fetch(url)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = epubgen.WriteEPUB(book, "dev/test.epub")
	if err != nil {
		fmt.Println(err.Error())
	}
}

package main

import (
	"fmt"

	"github.com/jabuxas/noctisyn/internal/novel"
)

func main() {
	url, err := novel.GetNovelURL("warlock magus")
	if err != nil {
		fmt.Println(err.Error())
	}
	book, err := novel.FetchBook(url)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(book.Chapters[38].Text)

	fmt.Printf("Scraped book: %s by %s\n(%d chapters)\nDesc: %s\n", book.Title, book.Author, len(book.Chapters), book.Description)
}

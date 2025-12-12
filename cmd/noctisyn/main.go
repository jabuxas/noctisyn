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

	// for _, ch := range book.Chapters {
	// 	fmt.Printf("ch index: %d\nch title: %s\nch html: %s\nch url: %s\n", ch.Index, ch.Title, ch.HTML, ch.URL)
	// }

	fmt.Printf("Scraped book: %s by %s\n(%d chapters)\nDesc: %s\n", book.Title, book.Author, len(book.Chapters), book.Description)
}

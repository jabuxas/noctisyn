package main

import (
	"github.com/jabuxas/noctisyn/internal/novel"
)

func main() {
	url, _ := novel.GetNovelURL("warlock magus")
	novel.FetchBook(url)
}

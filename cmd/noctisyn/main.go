package main

import (
	"fmt"

	"github.com/jabuxas/noctisyn/internal/novel"
)

func main() {
	url, _ := novel.GetNovelURL("warlock magus")
	fmt.Println(url)
}

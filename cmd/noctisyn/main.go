package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		chapterHref := strings.Contains(link, "chapter")
		if !chapterHref {
			h.Request.Visit(link)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	c.Visit("https://novelfull.net")
}

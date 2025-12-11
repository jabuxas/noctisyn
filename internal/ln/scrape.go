package ln

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

var collector *colly.Collector

func init() {
	collector = colly.NewCollector(
		colly.AllowedDomains("novelfull.net", "novgo.net"),
		colly.IgnoreRobotsTxt(),
	)

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})
}

func Scrape(query string) {
	c := collector

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		title := strings.TrimSpace(h.Text)

		href := h.Attr("href")

		if wordsInOrder(title, query) || wordsInOrder(href, query) {
			fmt.Println("matched: ", title, "->", h.Request.AbsoluteURL(href))
		}
	})

	search := fmt.Sprintf("https://novelfull.net/search?keyword=%s", url.QueryEscape(query))

	err := c.Visit(search)
	if err != nil {
		fmt.Println(err.Error())
	}
}

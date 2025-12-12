package scraper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func wordsInOrder(s, query string) bool {
	sLower := strings.ToLower(s)
	words := strings.Fields(strings.ToLower(query))

	pos := 0
	for _, w := range words {
		idx := strings.Index(sLower[pos:], w)
		if idx == -1 {
			return false
		}
		pos += idx + len(w)
	}
	return true
}

func extractChapterText(h *colly.HTMLElement) string {
	var parts []string

	h.DOM.Find("p").Each(func(_ int, p *goquery.Selection) {
		t := strings.TrimSpace(p.Text())
		if t != "" {
			parts = append(parts, t)
		}
	})

	return strings.Join(parts, "\n\n")
}

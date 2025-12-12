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
	sel := h.DOM.Clone()

	sel.Find("script, style, iframe, noscript").Remove()
	sel.Find(".ads, .ads-holder, .ads-middle, .text-center, #frame").Remove()

	var parts []string
	sel.Find("p").Each(func(_ int, p *goquery.Selection) {
		html, err := goquery.OuterHtml(p)
		if err == nil {
			html = strings.TrimSpace(html)
			if html != "<p></p>" && html != "" {
				parts = append(parts, html)
			}
		}
	})

	return strings.Join(parts, "\n\n")
}

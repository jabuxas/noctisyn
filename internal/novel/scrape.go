package novel

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Chapter struct {
	Index int
	Title string
	URL   string
	HTML  string
}

type Book struct {
	Title       string
	Author      string
	SourceURL   string
	Description string
	Chapters    []Chapter
}

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

func GetNovelURL(query string) (string, error) {
	c := collector.Clone()

	var matchURL string

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		title := strings.TrimSpace(h.Text)
		href := h.Attr("href")
		absURL := h.Request.AbsoluteURL(href)

		if matchURL != "" {
			return
		}
		// when scraping i get 3 urls, and only the 2nd one matches
		if wordsInOrder(title, query) && wordsInOrder(href, query) && !strings.Contains(h.Text, "Search") {
			matchURL = absURL
		}
	})

	search := fmt.Sprintf("https://novelfull.net/search?keyword=%s", url.QueryEscape(query))
	if err := c.Visit(search); err != nil {
		return "", err
	}

	if matchURL == "" {
		return "", fmt.Errorf("no matching novel found for '%s'", query)
	}
	return matchURL, nil
}

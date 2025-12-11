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

func FetchBook(novelURL string) (*Book, error) {
	c := collector.Clone()

	book := &Book{SourceURL: novelURL}

	c.OnHTML("body", func(h *colly.HTMLElement) { // runs on novel page
		book.Title = strings.TrimSpace(h.ChildText("h3.title"))
		book.Author = strings.TrimSpace(h.ChildText(".info a"))
		book.Description = strings.TrimSpace(h.ChildText(".desc-text"))

		var chIndex = 1
		h.ForEach("ul.list-chapter li a, div.chapters a[href*='chapter']", func(_ int, el *colly.HTMLElement) {
			chURL := el.Request.AbsoluteURL(el.Attr("href"))
			title := strings.TrimSpace(el.Text)
			book.Chapters = append(book.Chapters, Chapter{
				Index: chIndex,
				Title: title,
				URL:   chURL,
			})
			chIndex++
		})
	})

	if err := c.Visit(novelURL); err != nil {
		return nil, err
	}

	fmt.Printf("Scraped book:\n%s\nby\n%s\n(%d chapters)\n", book.Title, book.Author, len(book.Chapters))
	return book, nil
}

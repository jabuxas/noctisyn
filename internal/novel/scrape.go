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
		colly.Async(true),
	)
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

	c.Wait()

	if matchURL == "" {
		return "", fmt.Errorf("no matching novel found for '%s'", query)
	}
	return matchURL, nil
}

func FetchBook(novelURL string) (*Book, error) {
	infoCollector := collector.Clone()
	chapterCollector := collector.Clone()

	book := &Book{SourceURL: novelURL}

	infoCollector.OnHTML("body", func(h *colly.HTMLElement) {
		book.Title = strings.TrimSpace(h.ChildText(".books h3.title"))
		book.Author = strings.TrimSpace(h.ChildText("div.info > div:first-child > a:nth-child(2)"))
		book.Description = strings.TrimSpace(h.ChildText("div.desc-text"))

		chIndex := 1
		h.ForEach("ul.list-chapter li a", func(_ int, el *colly.HTMLElement) {
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

	chapterCollector.OnHTML("#chapter-content", func(h *colly.HTMLElement) {
		html, err := h.DOM.Html()
		if err != nil {
			fmt.Printf("failed to get html for chapter: %v\n", err)
		}
		u := h.Request.URL.String()
		for i := range book.Chapters {
			if book.Chapters[i].URL == u {
				book.Chapters[i].HTML = html
				break
			}
		}
	})

	if err := infoCollector.Visit(novelURL); err != nil {
		return nil, err
	}

	infoCollector.Wait()

	for _, ch := range book.Chapters {
		if err := chapterCollector.Visit(ch.URL); err != nil {
			fmt.Printf("queue failed %s: %v\n", ch.Title, err)
		}
	}

	chapterCollector.Wait()

	return book, nil
}

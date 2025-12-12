package scraper

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
	Text  string
}

type Book struct {
	Title       string
	Author      string
	SourceURL   string
	Description string
	CoverURL    string
	Chapters    []Chapter
}

var collector *colly.Collector

func init() {
	collector = colly.NewCollector(
		colly.AllowedDomains("novelfull.net", "novgo.net"),
		colly.IgnoreRobotsTxt(),
		colly.Async(true),
		colly.AllowURLRevisit(),
	)
}

func Search(query string) ([]*Book, error) {
	c := collector.Clone()

	var matched []*Book

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		title := strings.TrimSpace(h.Text)
		href := h.Attr("href")
		absURL := h.Request.AbsoluteURL(href)

		// when scraping i get 3 urls, and only the 2nd one matches
		if wordsInOrder(title, query) && wordsInOrder(href, query) && !strings.Contains(h.Text, "Search") {
			matched = append(matched, &Book{Title: title, SourceURL: absURL})
		}
	})

	search := fmt.Sprintf("https://novelfull.net/search?keyword=%s", url.QueryEscape(query))
	if err := c.Visit(search); err != nil {
		return matched, err
	}

	c.Wait()

	return matched, nil
}

func Fetch(novelURL string) (*Book, error) {
	infoCollector := collector.Clone()
	chapterCollector := collector.Clone()

	book := &Book{SourceURL: novelURL}

	infoCollector.OnHTML("body", func(h *colly.HTMLElement) {
		book.Title = strings.TrimSpace(h.ChildText(".books h3.title"))
		book.Author = strings.TrimSpace(h.ChildText("div.info > div:first-child > a:nth-child(2)"))
		book.Description = strings.TrimSpace(h.ChildText("div.desc-text"))

		chIndex := len(book.Chapters)
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

	infoCollector.OnHTML("img", func(h *colly.HTMLElement) {
		src := h.Attr("src")
		if src == "" {
			return
		}
		coverURL := h.Request.AbsoluteURL(src)
		book.CoverURL = coverURL
	})

	infoCollector.OnHTML("li.next > a[href]", func(h *colly.HTMLElement) {
		err := infoCollector.Visit(h.Request.AbsoluteURL(h.Attr("href")))
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	chapterCollector.OnHTML("#chapter-content", func(h *colly.HTMLElement) {
		text := extractChapterText(h)
		u := h.Request.URL.String()
		for i := range book.Chapters {
			if book.Chapters[i].URL == u {
				book.Chapters[i].Text = text
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
			return nil, err
		}
	}

	chapterCollector.Wait()

	return book, nil
}

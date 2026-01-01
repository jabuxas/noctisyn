package scraper

import (
	"fmt"
	"math/rand/v2"
	"net/url"
	"regexp"
	"strings"
	"time"

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

type ProgressCallback func(current, total int, estimatedTimeMs int64)

var (
	collector *colly.Collector
	mirrors   []string = []string{
		"https://novgo.net/",
		"https://novelfull.net/",
	}
)

func init() {
	collector = colly.NewCollector(
		colly.AllowedDomains("novelfull.net", "novgo.net"),
		colly.IgnoreRobotsTxt(),
		colly.Async(true),
		colly.AllowURLRevisit(),
	)
}

func randomSelectMirror() string {
	i := rand.IntN(2)
	return mirrors[i]
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

	search := fmt.Sprintf("%ssearch?keyword=%s", randomSelectMirror(), url.QueryEscape(query))
	if err := c.Visit(search); err != nil {
		return matched, err
	}

	c.Wait()

	return matched, nil
}

func FetchWithProgress(novelURL string, progressCb ProgressCallback) (*Book, error) {
	infoCollector := collector.Clone()
	chapterCollector := collector.Clone()

	chapterCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1000})

	re := regexp.MustCompile(`^https.*.net\/`)

	book := &Book{SourceURL: novelURL}
	
	var startTime time.Time
	fetchedCount := 0

	infoCollector.OnHTML("li.next > a[href]", func(h *colly.HTMLElement) {
		err := infoCollector.Visit(h.Request.AbsoluteURL(h.Attr("href")))
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	infoCollector.OnHTML("body", func(h *colly.HTMLElement) {
		if h.Request.Depth == 1 {
			book.Title = strings.TrimSpace(h.ChildText(".books h3.title"))
			book.Author = strings.TrimSpace(h.ChildText("div.info > div:first-child > a:nth-child(2)"))
			book.Description = strings.TrimSpace(h.ChildText("div.desc-text"))

			h.ForEach("img[src]", func(_ int, imgEl *colly.HTMLElement) {
				src := imgEl.Request.AbsoluteURL(imgEl.Attr("src"))
				if src != "" {
					book.CoverURL = src
				}
			})
		}

		chIndex := len(book.Chapters)
		h.ForEach("ul.list-chapter li a", func(_ int, el *colly.HTMLElement) {
			chURL := el.Request.AbsoluteURL(el.Attr("href"))
			chURL = re.ReplaceAllString(chURL, randomSelectMirror())
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
		text := extractChapterText(h)
		u := h.Request.URL.String()
		for i := range book.Chapters {
			if book.Chapters[i].URL == u {
				book.Chapters[i].Text = text
				fetchedCount++
				
				if progressCb != nil {
					elapsed := time.Since(startTime).Milliseconds()
					avgTimePerChapter := elapsed / int64(fetchedCount)
					remaining := len(book.Chapters) - fetchedCount
					estimatedTimeMs := avgTimePerChapter * int64(remaining)
					progressCb(fetchedCount, len(book.Chapters), estimatedTimeMs)
				}
				break
			}
		}
	})

	if err := infoCollector.Visit(novelURL); err != nil {
		return nil, err
	}

	infoCollector.Wait()
	
	startTime = time.Now()

	for _, ch := range book.Chapters {
		if err := chapterCollector.Visit(ch.URL); err != nil {
			return nil, err
		}
	}

	chapterCollector.Wait()

	return book, nil
}

func Fetch(novelURL string) (*Book, error) {
	return FetchWithProgress(novelURL, nil)
}

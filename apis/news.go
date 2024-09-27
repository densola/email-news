/*
news.go contains news-related operations.
*/
package apis

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type News struct {
	HNThreads        HNThreads
	TLDRTechArticles []Article
	MBArticles       []Article // On web, not on emails.
}

type Article struct {
	Title    string
	Link     string
	Overview string
}

type HNThreads struct {
	Items []HNThread `xml:"channel>item"`
}

type HNThread struct {
	Title    string `xml:"title"`
	Link     string `xml:"link"`
	Comments string `xml:"comments"`
}

func GetContent(date string) (News, error) {
	var n News

	err := n.getHN()
	if err != nil {
		return n, fmt.Errorf("getting hacker news content: %w", err)
	}

	err = n.getTLDR(date)
	if err != nil {
		return n, fmt.Errorf("getting tldr content: %w", err)
	}

	err = n.getMB()
	if err != nil {
		return n, fmt.Errorf("getting morning brew content: %w", err)
	}

	return n, nil
}

func (n *News) getHN() error {
	r, err := http.Get("https://hnrss.org/newest?points=100&comments=25&description=0")
	if err != nil {
		return fmt.Errorf("getting rss xml: %w", err)
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	var threads HNThreads
	err = xml.Unmarshal(body, &threads)
	if err != nil {
		return fmt.Errorf("parsing xml body: %w", err)
	}
	n.HNThreads = threads

	return nil
}

func (n *News) getTLDR(date string) error {
	date = date[0:4] + "-" + date[4:6] + "-" + date[6:8]

	url := "https://tldr.tech/tech/" + date

	c := colly.NewCollector()

	c.OnHTML("article.mt-3:not(.text-center)", func(elem *colly.HTMLElement) {
		t := strings.TrimSpace(elem.ChildText("h3"))
		l := strings.TrimSpace(elem.ChildAttr("a", "href"))
		o := strings.TrimSpace(elem.ChildText("div.newsletter-html"))

		if t == "" || l == "" || o == "" {
			return
		}

		l = strings.ReplaceAll(l, "utm_source=tldrnewsletter", "")

		// Remove the "?" at the end of the URL for aesthetic purposes
		if l[len(l)-1:] == "?" {
			l = l[:len(l)-1]
		}

		if !strings.Contains(t, "(Sponsor)") {
			article := Article{
				Title:    t,
				Link:     l,
				Overview: o,
			}

			n.TLDRTechArticles = append(n.TLDRTechArticles, article)
		}

	})

	c.Visit(url)

	return nil
}

func (n *News) getMB() error {
	// Exclude sunday news
	if time.Now().Weekday() == time.Weekday(0) {
		return nil
	}

	url := "https://www.morningbrew.com/daily/issues/latest"
	c := colly.NewCollector()

	c.OnHTML("td.card", func(elem *colly.HTMLElement) {
		t := strings.TrimSpace(elem.ChildText("td.tag-title h1 a"))
		o := elem.ChildText("td.story-content")

		if len(t) != 0 && len(o) == 0 {
			o = elem.ChildText("td.card-content")
			article := Article{
				Title: t,
				// No links needed
				Overview: o,
			}
			n.MBArticles = append(n.MBArticles, article)
		} else if len(t) != 0 {
			article := Article{
				Title:    t,
				Overview: o,
			}
			n.MBArticles = append(n.MBArticles, article)
		}
	})

	c.Visit(url)

	// Remove the last advertisement card
	if len(n.MBArticles) != 0 {
		n.MBArticles = n.MBArticles[:len(n.MBArticles)-1]
	}

	return nil
}

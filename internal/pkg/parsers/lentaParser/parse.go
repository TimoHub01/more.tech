package lentaParser

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	store2 "hack/internal/pkg/store"
	"time"
)

type store interface {
	CreateNew(ctx context.Context, topic, date, link, text string) (store2.New, error)
}

type Parser struct {
	store store
}

func NewParser(s store) *Parser {
	return &Parser{store: s}
}

func (p *Parser) postNews(ctx context.Context, link, topic, date string) {
	link = "https://lenta.ru/" + link

	text := textOfNew(link)
	_, err := p.store.CreateNew(ctx, topic, date, link, text)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *Parser) Parse(ctx context.Context) {
	currentTime := time.Date(2021, 10, 7, 10, 10, 10, 10, time.Local)
	for currentTime.Before(time.Now()) {
		c := colly.NewCollector()
		c.OnHTML("ul[class='archive-page__container']", func(e *colly.HTMLElement) {
			e.ForEach("a.card-full-news", func(_ int, el *colly.HTMLElement) {
				p.postNews(ctx, el.Attr("href"), el.ChildText("h3"), el.ChildText("time"))
			})
		})
		err := c.Visit("https://lenta.ru/news/" + currentTime.Format("2006/01/02"))
		if err != nil {
			fmt.Println(err)
		}
		currentTime = currentTime.AddDate(0, 0, 1)
	}
}

func textOfNew(href string) string {
	c := colly.NewCollector()
	text := ""
	c.OnHTML(".topic-body__content", func(e *colly.HTMLElement) {
		e.ForEach(".topic-body__content-text", func(_ int, el *colly.HTMLElement) {
			text += " " + el.Text + el.ChildText("a")
		})
	})
	err := c.Visit(href)
	if err != nil {
		fmt.Println(err)
	}
	return text
}

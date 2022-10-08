package newslab

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
	link = "https://newslab.ru" + link

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
		c.OnHTML(".n-default__content-wrapper", func(e *colly.HTMLElement) {
			e.ForEach("li", func(_ int, el *colly.HTMLElement) {
				go p.postNews(ctx, el.ChildAttr("a.n-list__item__link", "href"), el.ChildText("a.n-list__item__link"), currentTime.Format("2006-1-2"))
			})
		})
		err := c.Visit("https://newslab.ru/news/archive/" + currentTime.Format("2006/1/2"))
		if err != nil {
			fmt.Println(err)
		}
		currentTime = currentTime.AddDate(0, 0, 1)
	}
}

func textOfNew(href string) string {
	c := colly.NewCollector()
	text := ""
	c.OnHTML(".di3-body__text_news", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			text += " " + el.Text + el.ChildText("a")
		})
	})
	err := c.Visit(href)
	if err != nil {
		fmt.Println(err)
	}
	return text
}

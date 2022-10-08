package kommersantParser

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
	link = "https://www.kommersant.ru" + link

	text := textOfNew(link)
	_, err := p.store.CreateNew(ctx, topic, date, link, text)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *Parser) Parse(ctx context.Context) {
	fmt.Println("1")
	currentTime := time.Date(2021, 10, 7, 10, 10, 10, 10, time.Local)
	for currentTime.Before(time.Now()) {
		c := colly.NewCollector()
		c.OnHTML("div.grid-col", func(e *colly.HTMLElement) {
			e.ForEach(".js-article", func(_ int, el *colly.HTMLElement) {
				go p.postNews(ctx, el.ChildAttr("a", "href"), el.ChildText("a"), currentTime.Format("2006-01-02"))
			})
		})
		err := c.Visit("https://www.kommersant.ru/archive/news/day/" + currentTime.Format("2006-01-02"))
		if err != nil {
			fmt.Println(err)
		}
		currentTime = currentTime.AddDate(0, 0, 1)
	}
}

func textOfNew(href string) string {
	c := colly.NewCollector()
	text := ""
	c.OnHTML(".doc__body", func(e *colly.HTMLElement) {
		e.ForEach(".doc__text", func(_ int, el *colly.HTMLElement) {
			text += " " + el.Text + el.ChildText("a")
		})
	})
	err := c.Visit(href)
	if err != nil {
		fmt.Println(err)
	}
	return text
}

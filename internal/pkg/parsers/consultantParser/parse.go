package consultantParser

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	store2 "hack/internal/pkg/store"
	"strconv"
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

func (p *Parser) postNews(ctx context.Context, date, topic, link string) {
	link = "https://www.consultant.ru" + link
	text := textFromNew(link)
	_, err := p.store.CreateNew(ctx, topic, date, link, text)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *Parser) Parse(ctx context.Context, days int) {
	for page := 1; page <= days; page++ {
		c := colly.NewCollector()
		c.OnHTML(".listing-news__list", func(e *colly.HTMLElement) {
			e.ForEach(".listing-news__item", func(_ int, el *colly.HTMLElement) {
				p.postNews(ctx, el.ChildText("div[class='listing-news__item-date']"), el.ChildText("span"), el.ChildAttr("a", "href"))
			})
		})
		err := c.Visit("https://www.consultant.ru/legalnews/?page=" + strconv.Itoa(page))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func textFromNew(href string) string {
	c := colly.NewCollector()
	text := ""
	c.OnHTML("div[class='news-page__content']", func(e *colly.HTMLElement) {
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

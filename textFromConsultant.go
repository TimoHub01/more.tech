package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func textFromConsultant(href string) string {
	c := colly.NewCollector()
	text := ""
	c.OnHTML("div[class='news-page__content']", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			text += " " + e.Text + e.ChildText("a")
		})
	})
	err := c.Visit("https://www.consultant.ru/" + href)
	if err != nil {
		fmt.Println(err)
	}
	return text
}

package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func textFromLenta(href string) string {
	c := colly.NewCollector()
	text := ""
	c.OnHTML("div[class='topic-body _news']", func(e *colly.HTMLElement) {
		e.ForEach("div", func(_ int, el *colly.HTMLElement) {
			text += " " + e.Text + e.ChildText("a")
		})
	})
	err := c.Visit("https://lenta.ru/" + href)
	if err != nil {
		fmt.Println(err)
	}
	return text
}

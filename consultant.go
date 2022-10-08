package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
)

func consultant() {
	fName := "consultant.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Could not create file, err: %q", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for page := 1; page <= 940; page++ {
		c := colly.NewCollector()
		c.OnHTML("div[class='listing-news__item']", func(el *colly.HTMLElement) {
			writer.Write([]string{
				el.ChildText("div[class='listing-news__item-date']"),
				el.ChildText("span"),
				el.ChildAttr("a", "href"),
				textFromConsultant(el.ChildAttr("a", "href")),
			})
		})
		err = c.Visit("https://www.consultant.ru/legalnews/?page=" + strconv.Itoa(page))
		if err != nil {
			fmt.Println(err)
		}
	}
}

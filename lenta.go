package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"time"
)

func lenta() {
	fName := "lenta.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Could not create file, err: %q", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	currentTime := time.Date(2022, 10, 5, 10, 10, 10, 10, time.Local)
	for currentTime.Before(time.Now()) {
		c := colly.NewCollector()
		c.OnHTML("ul[class='archive-page__container']", func(e *colly.HTMLElement) {
			e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
				writer.Write([]string{
					el.Attr("href"),
					el.ChildText("h3"),
					el.ChildText("time"),
					textFromLenta(el.Attr("href")),
				})
			})
		})
		err = c.Visit("https://lenta.ru/news/" + currentTime.Format("2006/01/02"))
		if err != nil {
			fmt.Println(err)
		}
		currentTime = currentTime.AddDate(0, 0, 1)
	}
}

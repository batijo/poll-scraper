package scraper

import (
	"log"
	"os"
	"strings"

	"github.com/batijo/poll-scraper/models"
	"github.com/gocolly/colly/v2"
)

func Scrape() []models.Data {
	var data []models.Data
	c := colly.NewCollector()
	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		tds := e.ChildTexts(".pdg")
		if !(len(tds) < 2) {
			data = append(data, models.Data{
				Name:  tds[0],
				Value: tds[1],
			})
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.Visit(os.Getenv("PS_LINK"))
	return data
}

func ScrapeWithEquals() []models.Data {
	var data []models.Data
	c := colly.NewCollector()
	c.OnHTML("p", func(e *colly.HTMLElement) {
		tds := strings.Split(e.Text, "=")
		if !(len(tds) < 2) {
			data = append(data, models.Data{
				Name:  tds[0],
				Value: tds[1],
			})
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.Visit(os.Getenv("PS_LINK"))
	return data
}

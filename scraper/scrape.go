package scraper

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"

	"github.com/batijo/poll-scraper/models"
)

func ScrapeAll(links []string) []models.Data {
	var (
		nData []models.Data
		data  []models.Data
	)
	for _, l := range links {
		if os.Getenv("PS_WITH_EQ") == "true" {
			nData = ScrapeWithEquals(l)
		} else {
			nData = ScrapeWithoutEquals(l)
		}
		data = append(data, nData...)
	}
	return data
}

func ScrapeWithoutEquals(link string) []models.Data {
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
		slog.Error(fmt.Sprint("Request URL:", r.Request.URL, "failed"), "err", err)
	})
	c.Visit(link)
	return data
}

func ScrapeWithEquals(link string) []models.Data {
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
		slog.Error(fmt.Sprint("Request URL:", r.Request.URL, "failed"), "err", err)
	})
	c.Visit(link)
	return data
}

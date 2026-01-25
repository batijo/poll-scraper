package scraper

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gocolly/colly/v2"

	"github.com/batijo/poll-scraper/models"
)

const minParts = 2

func ScrapeAll(links []string, withEq bool) []models.Data {
	var (
		nData []models.Data
		data  []models.Data
	)
	for _, l := range links {
		if withEq {
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
		if len(tds) >= minParts {
			data = append(data, models.Data{
				Name:  tds[0],
				Value: tds[1],
			})
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		slog.Error(fmt.Sprint("Request URL:", r.Request.URL, "failed"), "err", err)
	})
	if err := c.Visit(link); err != nil {
		slog.Error("failed to visit link", "link", link, "err", err)
	}
	return data
}

func ScrapeWithEquals(link string) []models.Data {
	var data []models.Data
	c := colly.NewCollector()
	c.OnHTML("p", func(e *colly.HTMLElement) {
		tds := strings.Split(e.Text, "=")
		if len(tds) >= minParts {
			data = append(data, models.Data{
				Name:  tds[0],
				Value: tds[1],
			})
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		slog.Error(fmt.Sprint("Request URL:", r.Request.URL, "failed"), "err", err)
	})
	if err := c.Visit(link); err != nil {
		slog.Error("failed to visit link", "link", link, "err", err)
	}
	return data
}

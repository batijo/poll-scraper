package scraper

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gocolly/colly/v2"

	"github.com/batijo/poll-scraper/models"
)

const minParts = 2

func ScrapeURL(link string, withEq bool) []models.Data {
	if withEq {
		return ScrapeWithEquals(link)
	}
	return ScrapeWithoutEquals(link)
}

func ScrapeAll(links []string, withEq bool) []models.Data {
	slog.Debug("scraping all URLs", "count", len(links), "with_eq", withEq)
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
		slog.Debug("scraped URL", "url", l, "lines", len(nData))
		data = append(data, nData...)
	}
	slog.Debug("scraping complete", "total_lines", len(data))
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

package scraper

import (
	"fmt"
	"os"

	"github.com/gocolly/colly/v2"
)

func Scrape() {
	c := colly.NewCollector()
	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		fmt.Println(e.ChildText(".pdg"))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.Visit(os.Getenv("POLL_DATA_LINK"))
}

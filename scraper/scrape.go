package scraper

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"

	"github.com/batijo/poll-scraper/models"
)

const minParts = 2

var errLoginRequired = errors.New("page requires login, add a username and password to this link")

// ScrapeURL fetches link without a login session.
func ScrapeURL(link string, withEq bool) ([]models.Data, error) {
	data, loginURL, err := visit(newCollector(nil), link, withEq)
	if err != nil {
		return nil, err
	}
	if loginURL != "" {
		return nil, errLoginRequired
	}
	return data, nil
}

func newCollector(jar http.CookieJar) *colly.Collector {
	c := colly.NewCollector()
	if jar != nil {
		c.SetCookieJar(jar)
	}
	c.SetRequestTimeout(requestTimeout)
	return c
}

// visit returns the parsed lines and, when the server redirected to a login
// form, that form's URL.
func visit(c *colly.Collector, link string, withEq bool) ([]models.Data, string, error) {
	var data []models.Data
	var loginURL string

	if withEq {
		c.OnHTML("p", func(e *colly.HTMLElement) {
			if d, ok := parseEquals(e.Text); ok {
				data = append(data, d)
			}
		})
	} else {
		// Plain "tr": the live page has no <tbody>, so "tbody tr" only matched
		// because the HTML parser happens to insert one.
		c.OnHTML("tr", func(e *colly.HTMLElement) {
			tds := e.ChildTexts(".pdg")
			if len(tds) >= minParts {
				data = append(data, models.Data{Name: tds[0], Value: tds[1]})
			}
		})
	}
	c.OnResponse(func(r *colly.Response) {
		if isLoginPage(r.Request.URL) {
			loginURL = r.Request.URL.String()
		}
	})

	if err := c.Visit(link); err != nil {
		return nil, "", err
	}
	return data, loginURL, nil
}

// parseEquals splits on the first "=" only, so a value may itself contain "=".
func parseEquals(text string) (models.Data, bool) {
	name, value, ok := strings.Cut(text, "=")
	if !ok {
		return models.Data{}, false
	}
	return models.Data{Name: strings.TrimSpace(name), Value: strings.TrimSpace(value)}, true
}

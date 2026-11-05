package scraper

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"

	"github.com/batijo/poll-scraper/config"
	"github.com/batijo/poll-scraper/models"
)

const (
	requestTimeout = 30 * time.Second
	maxLoginBody   = 1 << 20
)

type Session struct {
	jar   http.CookieJar
	creds config.Link
	mu    sync.Mutex
}

// NewSession returns nil unless both credentials are provided, so callers can
// keep using the original scraper for links that do not need a login.
func NewSession(link config.Link) *Session {
	if link.Name == "" || link.Password == "" {
		return nil
	}
	jar, _ := cookiejar.New(nil)
	return &Session{jar: jar, creds: link}
}

func NewSessions(links []config.Link) []*Session {
	sessions := make([]*Session, len(links))
	for i, link := range links {
		sessions[i] = NewSession(link)
	}
	return sessions
}

func (s *Session) Scrape(link string, withEq bool) []models.Data {
	if s == nil {
		return ScrapeURL(link, withEq)
	}

	data, loginURL := s.scrape(link, withEq)
	if loginURL == "" {
		return data
	}

	slog.Info("login required, signing in", "user", s.creds.Name)
	if err := s.login(loginURL); err != nil {
		slog.Error("login failed", "user", s.creds.Name, "err", err)
		return nil
	}

	data, loginURL = s.scrape(link, withEq)
	if loginURL != "" {
		slog.Error("still on login page after signing in", "user", s.creds.Name)
		return nil
	}
	return data
}

func (s *Session) scrape(link string, withEq bool) ([]models.Data, string) {
	var data []models.Data
	var loginURL string

	c := colly.NewCollector()
	c.SetCookieJar(s.jar)
	c.SetRequestTimeout(requestTimeout)

	if withEq {
		c.OnHTML("p", func(e *colly.HTMLElement) {
			tds := strings.Split(e.Text, "=")
			if len(tds) >= minParts {
				data = append(data, models.Data{Name: tds[0], Value: tds[1]})
			}
		})
	} else {
		c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
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
	c.OnError(func(r *colly.Response, err error) {
		slog.Error("request failed", "url", r.Request.URL, "err", err)
	})

	if err := c.Visit(link); err != nil {
		slog.Error("failed to visit link", "link", link, "err", err)
	}
	return data, loginURL
}

func (s *Session) login(loginURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	client := &http.Client{Jar: s.jar, Timeout: requestTimeout}
	resp, err := client.Get(loginURL)
	if err != nil {
		return err
	}
	body, err := readBody(resp)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	values := url.Values{
		"Login1$UserName": {s.creds.Name},
		"Login1$Password": {s.creds.Password},
		"Login1$Button1":  {"Prisijungti"},
	}
	doc.Find("input[type=hidden]").Each(func(_ int, input *goquery.Selection) {
		name, _ := input.Attr("name")
		value, _ := input.Attr("value")
		if name != "" {
			values.Set(name, value)
		}
	})

	resp, err = client.PostForm(loginURL, values)
	if err != nil {
		return err
	}
	if _, err := readBody(resp); err != nil {
		return err
	}
	if isLoginPage(resp.Request.URL) {
		return errors.New("credentials rejected")
	}
	return nil
}

func readBody(resp *http.Response) ([]byte, error) {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Debug("failed to close response body", "err", err)
		}
	}()
	return io.ReadAll(io.LimitReader(resp.Body, maxLoginBody))
}

func isLoginPage(u *url.URL) bool {
	return u != nil && strings.Contains(strings.ToLower(u.Path), "login")
}

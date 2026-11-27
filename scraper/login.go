package scraper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

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

func (s *Session) Scrape(link string, withEq bool) ([]models.Data, error) {
	if s == nil {
		return ScrapeURL(link, withEq)
	}

	data, loginURL, err := visit(newCollector(s.jar), link, withEq)
	if err != nil || loginURL == "" {
		return data, err
	}

	slog.Info("login required, signing in", "user", s.creds.Name)
	if err := s.login(loginURL); err != nil {
		return nil, fmt.Errorf("login as %s: %w", s.creds.Name, err)
	}

	data, loginURL, err = visit(newCollector(s.jar), link, withEq)
	if err != nil {
		return nil, err
	}
	if loginURL != "" {
		return nil, fmt.Errorf("still on login page after signing in as %s", s.creds.Name)
	}
	return data, nil
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

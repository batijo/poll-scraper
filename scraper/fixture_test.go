package scraper

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	"github.com/batijo/poll-scraper/config"
	"github.com/batijo/poll-scraper/models"
)

// Fixtures in testdata are real responses saved from the live voting services,
// so these tests break if parsing drifts from what the servers actually send.

var wantEqualsFixture = []models.Data{
	{Name: "1", Value: "0"},
	{Name: "2", Value: "0"},
	{Name: "3", Value: "3"},
	{Name: "4", Value: "720"},
	{Name: "5", Value: "94800"},
	{Name: "6", Value: "0"},
	{Name: "7", Value: "0"},
	{Name: "8", Value: "0"},
	{Name: "9", Value: "95523"},
}

var wantTableFixture = []models.Data{
	{Name: "1343_wosub", Value: "94800"},
	{Name: "mb5", Value: "720"},
	{Name: "mb3", Value: "3"},
	{Name: "mb10", Value: "0"},
	{Name: "mb20", Value: "0"},
	{Name: "Teo", Value: "0"},
	{Name: "Aukok.lt", Value: "0"},
	{Name: "MaistoBankas", Value: "0"},
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}
	return body
}

func serveFixture(t *testing.T, name string) *httptest.Server {
	t.Helper()
	body := readFixture(t, name)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(body)
	}))
	t.Cleanup(ts.Close)
	return ts
}

func TestScrapeURL_RealEqualsPage(t *testing.T) {
	ts := serveFixture(t, "equals.html")

	got := ScrapeURL(ts.URL, true)

	if !reflect.DeepEqual(got, wantEqualsFixture) {
		t.Errorf("got %v\nwant %v", got, wantEqualsFixture)
	}
}

func TestScrapeURL_RealTablePage(t *testing.T) {
	ts := serveFixture(t, "table.html")

	got := ScrapeURL(ts.URL, false)

	if !reflect.DeepEqual(got, wantTableFixture) {
		t.Errorf("got %v\nwant %v", got, wantTableFixture)
	}
}

func hiddenValue(t *testing.T, page []byte, name string) string {
	t.Helper()
	m := regexp.MustCompile(`name="` + regexp.QuoteMeta(name) + `"[^>]*value="([^"]*)"`).FindSubmatch(page)
	if m == nil {
		t.Fatalf("hidden field %s not found in login fixture", name)
	}
	return string(m[1])
}

// Mirrors the live savitarna.mobi.lt flow: the data page redirects to
// Login.aspx, which serves the real ASP.NET form, and a valid post sets a cookie.
func TestSession_RealLoginPage(t *testing.T) {
	loginPage := readFixture(t, "login.html")
	tablePage := readFixture(t, "table.html")
	wantValidation := hiddenValue(t, loginPage, "__EVENTVALIDATION")
	wantVSID := hiddenValue(t, loginPage, "__VS_ID")

	mux := http.NewServeMux()
	mux.HandleFunc("/voting/VotesForScreen.aspx", func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(".ASPXAUTH"); err != nil || c.Value != "ok" {
			http.Redirect(w, r, "/Login.aspx?ReturnUrl="+url.QueryEscape(r.URL.RequestURI()), http.StatusFound)
			return
		}
		_, _ = w.Write(tablePage)
	})
	mux.HandleFunc("/Login.aspx", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse login form: %v", err)
			}
			f := r.PostForm
			if f.Get("__EVENTVALIDATION") == wantValidation && f.Get("__VS_ID") == wantVSID &&
				f.Get("Login1$UserName") == "user" && f.Get("Login1$Password") == "pass" &&
				f.Get("Login1$Button1") == "Prisijungti" {
				http.SetCookie(w, &http.Cookie{Name: ".ASPXAUTH", Value: "ok", Path: "/"})
				http.Redirect(w, r, r.URL.Query().Get("ReturnUrl"), http.StatusFound)
				return
			}
		}
		_, _ = w.Write(loginPage)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	session := NewSession(config.Link{Name: "user", Password: "pass"})
	got := session.Scrape(ts.URL+"/voting/VotesForScreen.aspx?id=5AB0479B", false)

	if !reflect.DeepEqual(got, wantTableFixture) {
		t.Errorf("got %v\nwant %v", got, wantTableFixture)
	}
}

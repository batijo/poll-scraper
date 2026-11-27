package scraper

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/batijo/poll-scraper/models"
)

func serveHTML(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(ts.Close)
	return ts
}

func TestScrapeURL_EqualsTrimsAndSplitsOnFirstEquals(t *testing.T) {
	ts := serveHTML(t, http.StatusOK, `<p> 5 = 100 </p><p>formula=a=b</p><p>no separator</p>`)

	got, err := ScrapeURL(ts.URL, true)

	if err != nil {
		t.Fatalf("ScrapeURL() error = %v", err)
	}
	want := []models.Data{{Name: "5", Value: "100"}, {Name: "formula", Value: "a=b"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestScrapeURL_TableRowsNeedTwoCells(t *testing.T) {
	ts := serveHTML(t, http.StatusOK, `<table>
		<tr><td class="pdg">Item1</td><td class="pdg">100</td></tr>
		<tr><td>&nbsp;</td><td>&nbsp;</td></tr>
	</table>`)

	got, err := ScrapeURL(ts.URL, false)

	if err != nil {
		t.Fatalf("ScrapeURL() error = %v", err)
	}
	want := []models.Data{{Name: "Item1", Value: "100"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestScrapeURL_ServerErrorIsAnError(t *testing.T) {
	ts := serveHTML(t, http.StatusInternalServerError, `<p>1=0</p>`)

	got, err := ScrapeURL(ts.URL, true)

	if err == nil {
		t.Fatal("ScrapeURL() error = nil, want error for HTTP 500")
	}
	if len(got) != 0 {
		t.Errorf("got %d items, want 0", len(got))
	}
}

func TestScrapeURL_UnreachableHostIsAnError(t *testing.T) {
	ts := serveHTML(t, http.StatusOK, "")
	ts.Close()

	if _, err := ScrapeURL(ts.URL, true); err == nil {
		t.Fatal("ScrapeURL() error = nil, want connection error")
	}
}

func TestScrapeURL_LoginRedirectWithoutCredentials(t *testing.T) {
	ts, _ := newAuthServer(t, testUser, testPassword)

	_, err := ScrapeURL(ts.URL+"/data", false)

	if !errors.Is(err, errLoginRequired) {
		t.Errorf("ScrapeURL() error = %v, want errLoginRequired", err)
	}
}

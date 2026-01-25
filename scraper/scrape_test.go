package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScrapeWithoutEquals(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`
			<html><body><table><tbody>
				<tr><td class="pdg">Item1</td><td class="pdg">100</td></tr>
				<tr><td class="pdg">Item2</td><td class="pdg">200</td></tr>
			</tbody></table></body></html>
		`))
	}))
	defer ts.Close()

	data := ScrapeWithoutEquals(ts.URL)

	if len(data) != 2 {
		t.Fatalf("got %d items, want 2", len(data))
	}
	if data[0].Name != "Item1" || data[0].Value != "100" {
		t.Errorf("data[0] = {%q, %q}, want {%q, %q}", data[0].Name, data[0].Value, "Item1", "100")
	}
}

func TestScrapeWithEquals(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><body><p>Name1=Value1</p><p>Name2=Value2</p></body></html>`))
	}))
	defer ts.Close()

	data := ScrapeWithEquals(ts.URL)

	if len(data) != 2 {
		t.Fatalf("got %d items, want 2", len(data))
	}
	if data[0].Name != "Name1" || data[0].Value != "Value1" {
		t.Errorf("data[0] = {%q, %q}, want {%q, %q}", data[0].Name, data[0].Value, "Name1", "Value1")
	}
}

func TestScrapeAll_WithEquals(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><body><p>Key=Val</p></body></html>`))
	}))
	defer ts.Close()

	data := ScrapeAll([]string{ts.URL}, true)

	if len(data) != 1 {
		t.Fatalf("got %d items, want 1", len(data))
	}
	if data[0].Name != "Key" || data[0].Value != "Val" {
		t.Errorf("data[0] = {%q, %q}, want {%q, %q}", data[0].Name, data[0].Value, "Key", "Val")
	}
}

func TestScrapeAll_WithoutEquals(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`
			<html><body><table><tbody>
				<tr><td class="pdg">Item1</td><td class="pdg">100</td></tr>
			</tbody></table></body></html>
		`))
	}))
	defer ts.Close()

	data := ScrapeAll([]string{ts.URL}, false)

	if len(data) != 1 {
		t.Fatalf("got %d items, want 1", len(data))
	}
	if data[0].Name != "Item1" || data[0].Value != "100" {
		t.Errorf("data[0] = {%q, %q}, want {%q, %q}", data[0].Name, data[0].Value, "Item1", "100")
	}
}

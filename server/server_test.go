package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithMiddleware_SetsHeaders(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := withMiddleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if got := rec.Header().Get("Server"); got != "poll-scraper" {
		t.Errorf("Server header = %q, want %q", got, "poll-scraper")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "*")
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got != "GET" {
		t.Errorf("Access-Control-Allow-Methods = %q, want %q", got, "GET")
	}
}

func TestWithMiddleware_OptionsRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for OPTIONS request")
	})

	wrapped := withMiddleware(handler)
	req := httptest.NewRequest(http.MethodOptions, "/", http.NoBody)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestWithMiddleware_CustomOrigins(t *testing.T) {
	t.Setenv("PS_DOMAINS", "https://example.com")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := withMiddleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "https://example.com")
	}
}

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	WriteJSON(rec, http.StatusOK, data)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
	want := `{"key":"value"}` + "\n"
	if got := rec.Body.String(); got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, http.StatusNotFound, "not found")

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	want := `{"message":"not found"}` + "\n"
	if got := rec.Body.String(); got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestNew_ReturnsServer(t *testing.T) {
	srv := New()

	if srv == nil {
		t.Fatal("New() returned nil")
	}
	if srv.Server == nil {
		t.Error("Server.Server is nil")
	}
	if srv.mux == nil {
		t.Error("Server.mux is nil")
	}
	if srv.Handler == nil {
		t.Error("Server.Handler is nil")
	}
}

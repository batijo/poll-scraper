package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/batijo/poll-scraper/models"
)

func TestData_ReturnsJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()

	Data(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	var data []models.Data
	if err := json.NewDecoder(rec.Body).Decode(&data); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
}

func TestData_WithFilters(t *testing.T) {
	t.Setenv("PS_LINKS", "")
	t.Setenv("PS_FILTER_LINES", "")
	t.Setenv("PS_ADD_LINES", "")
	t.Setenv("PS_ADD_SUM", "")

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()

	Data(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var data []models.Data
	if err := json.NewDecoder(rec.Body).Decode(&data); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
}

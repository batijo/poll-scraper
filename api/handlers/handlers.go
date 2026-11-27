package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/batijo/poll-scraper/config"
	"github.com/batijo/poll-scraper/models"
	"github.com/batijo/poll-scraper/scraper"
)

func Data(cfg *config.Config) http.HandlerFunc {
	sessions := scraper.NewSessions(cfg.Links)
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("HTTP request received", "method", r.Method, "remote", r.RemoteAddr)
		var data []models.Data
		for i, link := range cfg.Links {
			urlData, err := sessions[i].Scrape(link.URL, cfg.WithEq)
			if err != nil {
				// Partial data would shift every later line, so fail the whole response.
				slog.Error("HTTP scrape failed", "url", link.URL, "err", err)
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
			data = append(data, urlData...)
		}
		data = cfg.Process(data)
		slog.Debug("HTTP response", "lines", len(data))
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			slog.Error("failed to encode response", "err", err)
		}
	}
}

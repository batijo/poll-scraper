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
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("HTTP request received", "method", r.Method, "remote", r.RemoteAddr)
		data := scraper.ScrapeAll(cfg.Links, cfg.WithEq)
		lines := cfg.FilterLinesZeroIndexed()
		if len(lines) > 0 {
			data = models.FilterData(lines, data)
		}
		if len(cfg.AddLines) > 0 {
			addLines := make([]models.Data, len(cfg.AddLines))
			for i, l := range cfg.AddLines {
				addLines[i] = models.Data{Name: l.Name, Value: l.Value}
			}
			data = models.AddLines(data, addLines)
		}
		if cfg.AddSum {
			data = models.SumData(data, cfg.SumSymbols)
		}
		slog.Debug("HTTP response", "lines", len(data))
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			slog.Error("failed to encode response", "err", err)
		}
	}
}

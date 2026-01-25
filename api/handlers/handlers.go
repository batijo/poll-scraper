package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/batijo/poll-scraper/models"
	"github.com/batijo/poll-scraper/scraper"
	"github.com/batijo/poll-scraper/utils"
)

func Data(w http.ResponseWriter, r *http.Request) {
	links := strings.Split(os.Getenv("PS_LINKS"), " ")
	data := scraper.ScrapeAll(links)
	lines, err := utils.GetFilterLines("PS_FILTER_LINES")
	if err != nil {
		slog.Error("cannot parse filter lines", "err", err)
	}
	if len(lines) > 0 {
		data = models.FilterData(lines, data)
	}
	if os.Getenv("PS_ADD_LINES") != "" {
		data = models.AddLines(data)
	}
	if os.Getenv("PS_ADD_SUM") == utils.EnvTrue {
		data = models.SumData(data)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "err", err)
	}
}

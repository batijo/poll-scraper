package file

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/batijo/poll-scraper/config"
	"github.com/batijo/poll-scraper/models"
	"github.com/batijo/poll-scraper/scraper"
	"github.com/batijo/poll-scraper/utils"
)

type EventEmitter interface {
	EmitScraperData(data, rawData []models.Data)
	EmitScraperState(state string)
	EmitScraperError(message string)
	EmitURLStatus(statuses []models.URLStatus)
	RequestScraperStop()
}

func StartWriting(cfg *config.Config, emitter EventEmitter) (context.CancelFunc, error) {
	if cfg.UpdateInterval < 0 {
		slog.Error("update_interval cannot be negative")
		return nil, fmt.Errorf("invalid value")
	}
	if cfg.UpdateInterval < utils.MinIntervalWarn {
		slog.Warn("setting update_interval too low might cause high CPU usage and/or server load")
	}
	slog.Info("scraper started", "interval", cfg.UpdateInterval, "urls", len(cfg.Links))
	ctx, cancel := context.WithCancel(context.Background())
	go writer(ctx, cfg, emitter)
	return cancel, nil
}

//nolint:gocyclo,funlen // main scrape loop with inherent complexity
func writer(ctx context.Context, cfg *config.Config, emitter EventEmitter) {
	cycle := 0
	expectedLineCounts := make(map[string]int)

	for {
		select {
		case <-ctx.Done():
			slog.Info("scraper stopped")
			return
		default:
		}
		emitter.EmitScraperState("scraping")
		cycle++

		start := time.Now()
		lines := cfg.FilterLinesZeroIndexed()

		var data []models.Data
		statuses := make([]models.URLStatus, 0, len(cfg.Links))
		lineCountChanged := false
		for _, link := range cfg.Links {
			urlData := scraper.ScrapeURL(link, cfg.WithEq)
			status := models.URLStatus{
				URL:       link,
				HasData:   len(urlData) > 0,
				LineCount: len(urlData),
			}
			if len(urlData) == 0 {
				slog.Warn("no data from URL", "url", link)
			}

			if expected, ok := expectedLineCounts[link]; ok {
				if len(urlData) != expected {
					slog.Error("URL line count changed", "url", link, "expected", expected, "got", len(urlData))
					emitter.EmitScraperError(fmt.Sprintf("URL line count changed for %s: expected %d, got %d", link, expected, len(urlData)))
					status.Error = true
					lineCountChanged = true
				}
			}
			expectedLineCounts[link] = len(urlData)

			statuses = append(statuses, status)
			data = append(data, urlData...)
		}
		emitter.EmitURLStatus(statuses)

		if lineCountChanged && cfg.StopOnLineCountChange {
			emitter.EmitScraperState("error")
			emitter.RequestScraperStop()
			return
		}

		// rawData = URL-scraped data only (for frontend filter modal)
		rawData := make([]models.Data, len(data))
		copy(rawData, data)

		rawCount := len(data)
		if len(lines) > 0 {
			data = models.FilterData(lines, data)
			slog.Debug("filtered lines", "before", rawCount, "after", len(data))
		}

		// Add non-filtered custom lines after URL filtering
		if len(cfg.AddLines) > 0 {
			var addLines []models.Data
			for _, l := range cfg.AddLines {
				if !l.Filtered {
					addLines = append(addLines, models.Data{Name: l.Name, Value: l.Value})
				}
			}
			if len(addLines) > 0 {
				data = models.AddLines(data, addLines)
				slog.Debug("added custom lines", "count", len(addLines))
			}
		}
		if cfg.AddSum {
			data = models.SumData(data, cfg.SumSymbols)
		}

		hasError := false
		if cfg.WriteToCSV {
			if err := writeToCsv(data, cfg.CSVPath); err != nil {
				slog.Error("failed to write to CSV file", "err", err)
				emitter.EmitScraperError(fmt.Sprintf("failed to write to CSV file: %v", err))
				hasError = true
			} else {
				slog.Debug("wrote CSV", "path", cfg.CSVPath, "lines", len(data))
			}
		}
		if cfg.WriteToTXT {
			if err := writeToTxt(data, cfg.TXTPath, cfg.DatasetName, cfg.TXTEncoding); err != nil {
				slog.Warn("TXT write returned non-fatal error", "err", err)
			} else {
				slog.Debug("wrote TXT", "path", cfg.TXTPath, "lines", len(data))
			}
		}

		emitter.EmitScraperData(data, rawData)

		elapsed := time.Since(start)
		slog.Debug("scrape cycle complete", "cycle", cycle, "lines", len(data), "took", elapsed.Round(time.Millisecond))

		if hasError {
			emitter.EmitScraperState("error")
		} else {
			emitter.EmitScraperState("idle")
		}
		remaining := time.Duration(cfg.UpdateInterval)*time.Millisecond - elapsed
		if remaining > 0 {
			timer := time.NewTimer(remaining)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		} else {
			slog.Warn(fmt.Sprintf("scrape took %v, which is longer than the update interval of %d ms\n", elapsed, cfg.UpdateInterval))
		}
	}
}

func writeToCsv(data []models.Data, csvPath string) (err error) {
	cleanPath := filepath.Clean(csvPath)
	f, err := os.OpenFile(cleanPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, utils.FileMode)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	writer := csv.NewWriter(f)
	defer writer.Flush()
	for _, d := range data {
		row := []string{d.Name, d.Value}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func writeToTxt(data []models.Data, txtPath, datasetName, txtEncoding string) (err error) {
	cleanPath := filepath.Clean(txtPath)
	f, err := os.OpenFile(cleanPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, utils.FileMode)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	var writer *bufio.Writer
	if txtEncoding == "utf-8" {
		writer = bufio.NewWriter(f)
	} else {
		ansiEncoder := charmap.Windows1252.NewEncoder()
		encodedFile := transform.NewWriter(f, ansiEncoder)
		writer = bufio.NewWriter(encodedFile)
	}
	if datasetName != "" {
		_, err = fmt.Fprintf(writer, "[%s]\nCount=%v\n", datasetName, len(data))
	} else {
		_, err = fmt.Fprintf(writer, "Count=%v\n", len(data))
	}
	if err != nil {
		return err
	}
	for i, d := range data {
		_, err = fmt.Fprintf(writer, "Value%v=%v\n", i+1, d.Value)
		if err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}

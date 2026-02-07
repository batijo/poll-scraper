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
	EmitScraperData(data []models.Data)
	EmitScraperState(state string)
	EmitScraperError(message string)
	EmitURLStatus(statuses []models.URLStatus)
}

func StartWriting(cfg *config.Config, emitter EventEmitter) (context.CancelFunc, error) {
	if cfg.UpdateInterval < 0 {
		slog.Error("update_interval cannot be negative")
		return nil, fmt.Errorf("invalid value")
	}
	if cfg.UpdateInterval < utils.MinIntervalWarn {
		slog.Warn("setting update_interval too low might cause high CPU usage and/or server load")
	}
	ctx, cancel := context.WithCancel(context.Background())
	go writer(ctx, cfg, emitter)
	return cancel, nil
}

func writer(ctx context.Context, cfg *config.Config, emitter EventEmitter) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		emitter.EmitScraperState("scraping")

		start := time.Now()
		lines := cfg.FilterLinesZeroIndexed()

		var data []models.Data
		statuses := make([]models.URLStatus, 0, len(cfg.Links))
		for _, link := range cfg.Links {
			urlData := scraper.ScrapeURL(link, cfg.WithEq)
			statuses = append(statuses, models.URLStatus{
				URL:     link,
				HasData: len(urlData) > 0,
			})
			data = append(data, urlData...)
		}
		emitter.EmitURLStatus(statuses)
		if len(lines) > 0 {
			data = models.FilterData(lines, data)
		}
		if len(cfg.AddLines) > 0 {
			data = models.AddLines(data, cfg.AddLines)
		}
		if cfg.AddSum {
			data = models.SumData(data, cfg.SumSymbols)
		}

		hasError := false
		if cfg.WriteToCSV {
			err := writeToCsv(data, cfg.CSVPath)
			if err != nil {
				slog.Error("failed to write to CSV file", "err", err)
				emitter.EmitScraperError(fmt.Sprintf("failed to write to CSV file: %v", err))
				hasError = true
			}
		}
		if cfg.WriteToTXT {
			err := writeToTxt(data, cfg.TXTPath, cfg.DatasetName)
			if err != nil {
				slog.Error("failed to write to TXT file", "err", err)
				emitter.EmitScraperError(fmt.Sprintf("failed to write to TXT file: %v", err))
				hasError = true
			}
		}

		emitter.EmitScraperData(data)

		if hasError {
			emitter.EmitScraperState("error")
		} else {
			emitter.EmitScraperState("idle")
		}

		elapsed := time.Since(start)
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

func writeToTxt(data []models.Data, txtPath, datasetName string) (err error) {
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
	ansiEncoder := charmap.Windows1252.NewEncoder()
	encodedFile := transform.NewWriter(f, ansiEncoder)
	writer := bufio.NewWriter(encodedFile)
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

package file

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/batijo/poll-scraper/models"
	"github.com/batijo/poll-scraper/scraper"
	"github.com/batijo/poll-scraper/utils"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func StartWriting() error {
	interval, err := strconv.Atoi(os.Getenv("PS_UPDATE_INTERVAL"))
	if err != nil {
		slog.Error("PS_UPDATE_INTERVAL might have incorrect symbols")
		return err
	} else if interval < 0 {
		slog.Error("PS_UPDATE_INTERVAL cannot be negative")
		return fmt.Errorf("invalid value")
	}
	if interval < 500 {
		slog.Warn("setting PS_UPDATE_INTERVAL might cause high CPU usage and/or server load")
	}
	go writer(interval)
	return nil
}

func writer(interval int) {
	for {
		start := time.Now()
		lines, err := utils.GetFilterLines("PS_FILTER_LINES")
		if err != nil {
			slog.Error("cannot parse filter lines", "err", err)
		}
		links := strings.Split(os.Getenv("PS_LINKS"), " ")
		data := scraper.ScrapeAll(links)
		if len(lines) > 0 {
			data = models.FilterData(lines, data)
		}
		if os.Getenv("PS_ADD_LINES") != "" {
			data = models.AddLines(data)
		}
		if os.Getenv("PS_ADD_SUM") == "true" {
			data = models.SumData(data)
		}
		if os.Getenv("PS_WRITE_TO_CSV") == "true" {
			err = writeToCsv(data)
			if err != nil {
				slog.Error("failed to write to CSV file", "err", err)
			}
		}
		if os.Getenv("PS_WRITE_TO_TXT") == "true" {
			err = writeToTxt(data)
			if err != nil {
				slog.Error("failed to write to TXT file", "err", err)
			}
		}
		elapsed := time.Since(start)
		remaining := time.Duration(interval)*time.Millisecond - elapsed
		if remaining > 0 {
			time.Sleep(remaining)
		} else {
			slog.Warn(fmt.Sprintf("scrape took %v, which is longer than the update interval of %d ms\n", elapsed, interval))
		}
	}
}

func writeToCsv(data []models.Data) error {
	f, err := os.OpenFile(os.Getenv("PS_CSV_PATH"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
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

func writeToTxt(data []models.Data) error {
	filePath := os.Getenv("PS_TXT_PATH")
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	ansiEncoder := charmap.Windows1252.NewEncoder()
	encodedFile := transform.NewWriter(f, ansiEncoder)
	writer := bufio.NewWriter(encodedFile)
	if os.Getenv("PS_DATASET_NAME") != "" {
		_, err = fmt.Fprintf(writer, "[%s]\nCount=%v\n", os.Getenv("PS_DATASET_NAME"), len(data))
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

package csv

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/batijo/poll-scraper/models"
	"github.com/batijo/poll-scraper/scraper"
)

func StartWriting() error {
	interval, err := strconv.Atoi(os.Getenv("PS_CSV_UPDATE_INTERVAL"))
	if err != nil {
		return err
	}
	go writer(interval)
	return nil
}

func writer(interval int) {
	if interval <= 50 {
		log.Println("WARNING: setting PS_CSV_UPDATE_INTERVAL too low might overload CPU")
	}
	for {
		var err error
		if os.Getenv("PS_WITH_EQ") == "true" {
			err = writeToCsv(scraper.ScrapeWithEquals())
		} else {
			err = writeToCsv(scraper.Scrape())
		}
		if err != nil {
			log.Println("ERROR: failed to write to CSV file")
			log.Println(err)
			return
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

func writeToCsv(data []models.Data) error {
	file, err := os.OpenFile(os.Getenv("PS_CSV_PATH"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		file.Close()
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, d := range data {
		row := []string{d.Name, d.Value}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

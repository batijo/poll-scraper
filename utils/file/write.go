package file

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/batijo/poll-scraper/models"
	"github.com/batijo/poll-scraper/scraper"
)

func StartWriting() error {
	interval, err := strconv.Atoi(os.Getenv("PS_UPDATE_INTERVAL"))
	if err != nil {
		log.Println("WARNING: PS_UPDATE_INTERVAL might have incorrect symbols")
		return err
	}
	go writer(interval)
	return nil
}

func writer(interval int) {
	if interval <= 50 {
		log.Println("WARNING: setting PS_UPDATE_INTERVAL too low might overload CPU")
	}
	for {
		var (
			err  error
			data []models.Data
		)
		if os.Getenv("PS_WITH_EQ") == "true" {
			data = scraper.ScrapeWithEquals("PS_LINK")
		} else {
			data = scraper.Scrape("PS_LINK")
		}
		if os.Getenv("PS_WRITE_TO_CSV") == "true" {
			err = writeToCsv(data)
			if err != nil {
				log.Println("ERROR: failed to write to CSV file")
				log.Println(err)
			}
		}
		if os.Getenv("PS_WRITE_TO_TXT") == "true" {
			err = writeToTxt(data)
			if err != nil {
				log.Println("ERROR: failed to write to TXT file")
				log.Println(err)
			}
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

func writeToCsv(data []models.Data) error {
	f, err := os.OpenFile(os.Getenv("PS_CSV_PATH"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		f.Close()
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
	f, err := os.OpenFile(os.Getenv("PS_TXT_PATH"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		f.Close()
		return err
	}
	defer f.Close()
	// err = f.Truncate(0)
	// if err != nil {
	// 	return err
	// }
	// _, err = f.Seek(0, 0)
	// if err != nil {
	// 	return err
	// }
	_, err = fmt.Fprintf(f, "Count=%v\n", len(data))
	if err != nil {
		return err
	}
	for i, d := range data {
		_, err = fmt.Fprintf(f, "Value%v=%v\n", i+1, d.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
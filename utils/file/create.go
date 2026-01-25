package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/batijo/poll-scraper/config"
	"github.com/batijo/poll-scraper/utils"
)

func InitFiles(cfg *config.Config) error {
	if cfg.WriteToCSV {
		if filepath.Ext(cfg.CSVPath) != ".csv" {
			return fmt.Errorf("file must be of type CSV")
		}
		err := create(cfg.CSVPath)
		if err != nil {
			return err
		}
	}
	if cfg.WriteToTXT {
		if filepath.Ext(cfg.TXTPath) != ".txt" {
			return fmt.Errorf("file must be of type TXT")
		}
		err := create(cfg.TXTPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func create(path string) error {
	cleanPath := filepath.Clean(path)
	f, err := os.OpenFile(cleanPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, utils.FileMode)
	if err != nil {
		return err
	}
	return f.Close()
}

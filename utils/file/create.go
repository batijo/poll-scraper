package file

import (
	"fmt"
	"log/slog"
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
		if err := create(cfg.CSVPath); err != nil {
			return err
		}
		slog.Info("output file initialized", "type", "CSV", "path", cfg.CSVPath)
	}
	if cfg.WriteToTXT {
		if filepath.Ext(cfg.TXTPath) != ".txt" {
			return fmt.Errorf("file must be of type TXT")
		}
		if err := create(cfg.TXTPath); err != nil {
			return err
		}
		slog.Info("output file initialized", "type", "TXT", "path", cfg.TXTPath)
	}
	return nil
}

func create(path string) error {
	cleanPath := filepath.Clean(path)
	f, err := os.OpenFile(cleanPath, os.O_RDWR|os.O_CREATE, utils.FileMode)
	if err != nil {
		return err
	}
	return f.Close()
}

package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/batijo/poll-scraper/utils"
)

func InitFiles() error {
	if os.Getenv("PS_WRITE_TO_CSV") == utils.EnvTrue {
		path := os.Getenv("PS_CSV_PATH")
		if filepath.Ext(path) != ".csv" {
			return fmt.Errorf("file must be of type CSV")
		}
		err := create(path)
		if err != nil {
			return err
		}
	}
	if os.Getenv("PS_WRITE_TO_TXT") == utils.EnvTrue {
		path := os.Getenv("PS_TXT_PATH")
		if filepath.Ext(path) != ".txt" {
			return fmt.Errorf("file must be of type TXT")
		}
		err := create(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func create(path string) error {
	cleanPath := filepath.Clean(path)
	f, err := os.OpenFile(cleanPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	return f.Close()
}

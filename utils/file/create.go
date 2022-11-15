package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitFiles() error {
	if os.Getenv("PS_WRITE_TO_CSV") == "true" {
		path := os.Getenv("PS_CSV_PATH")
		if filepath.Ext(path) != ".csv" {
			return fmt.Errorf("file must be of type CSV")
		}
		err := create(path)
		if err != nil {
			return err
		}
	}
	if os.Getenv("PS_WRITE_TO_TXT") == "true" {
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
	f, err := os.Create(path)
	if err != nil {
		f.Close()
		return err
	}
	defer f.Close()
	return nil
}

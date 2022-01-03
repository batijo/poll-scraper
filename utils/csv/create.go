package csv

import (
	"fmt"
	"os"
	"path/filepath"
)

func Create(path string) error {
	var err error
	if filepath.Ext(path) != ".csv" {
		return fmt.Errorf("file must be of type CSV")
	}
	file, err := os.Create(path)
	if err != nil {
		file.Close()
		return err
	}
	defer file.Close()
	return nil
}

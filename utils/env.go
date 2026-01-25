package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetFilterLines(env string) ([]int, error) {
	stringLines := strings.Split(os.Getenv(env), ",")
	if len(stringLines) == 0 || stringLines[0] == "" {
		return nil, nil
	}
	var lines []int
	for _, sl := range stringLines {
		l, err := strconv.Atoi(sl)
		if err != nil {
			return nil, err
		}
		lines = append(lines, l-1)
	}
	return lines, nil
}

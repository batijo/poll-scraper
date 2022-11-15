package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetFilterLines(env string) ([]int, error) {
	stringLines := strings.Split(os.Getenv(env), ",")
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
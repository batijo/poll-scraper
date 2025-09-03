package models

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func FilterData(lines []int, unfilteredData []Data) []Data {
	var data []Data
	for _, l := range lines {
		if l < len(unfilteredData) && l > -1 {
			data = append(data, unfilteredData[l])
		}
	}
	return data
}

func SumData(data []Data) []Data {
	var sum int
	for _, d := range data {
		v, err := strconv.Atoi(d.Value)
		if err != nil {
			slog.Warn(fmt.Sprintf("cannot convert value of [%s] to integer\n", d.Value), "err", err)
			continue
		}
		sum += v
	}
	data = append(data, Data{Name: "sum", Value: strconv.Itoa(sum)})
	if os.Getenv("PS_SUM_SYMBOLS") != "" {
		data = append(data, Data{Name: "sum_symbol", Value: strconv.Itoa(sum) + os.Getenv("PS_SUM_SYMBOLS")})
	}
	return data
}

func AddLines(data []Data) []Data {
	envValue := os.Getenv("PS_ADD_LINES")
	if envValue == "" {
		return data
	}
	if envValue[0] != '[' || envValue[len(envValue)-1] != ']' {
		slog.Warn("PS_ADD_LINES should be an array")
		return data
	}
	envValue = strings.Trim(envValue, "[]")
	values := strings.Split(envValue, ",")
	if values[0] == "" {
		slog.Warn("PS_ADD_LINES should contain at least one value")
		return data
	}
	for _, v := range values {
		slog.Info(fmt.Sprintf("adding line [%s] as number %v\n", v, len(data)+1))
		data = append(data, Data{Name: strconv.Itoa(len(data) + 1), Value: v})
	}
	return data
}

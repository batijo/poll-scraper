package models

import (
	"fmt"
	"log/slog"
	"strconv"
)

type Data struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type URLStatus struct {
	URL     string `json:"url"`
	HasData bool   `json:"hasData"`
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

func SumData(data []Data, sumSymbols string) []Data {
	var sum int
	for _, d := range data {
		v, err := strconv.Atoi(d.Value)
		if err != nil {
			slog.Debug(fmt.Sprintf("cannot convert value of [%s] to integer", d.Value), "err", err)
			continue
		}
		sum += v
	}
	data = append(data, Data{Name: "sum", Value: strconv.Itoa(sum)})
	if sumSymbols != "" {
		data = append(data, Data{Name: "sum_symbol", Value: strconv.Itoa(sum) + sumSymbols})
	}
	return data
}

func AddLines(data []Data, values []string) []Data {
	if len(values) == 0 {
		return data
	}
	for _, v := range values {
		slog.Debug(fmt.Sprintf("adding line [%s] as number %v", v, len(data)+1))
		data = append(data, Data{Name: strconv.Itoa(len(data) + 1), Value: v})
	}
	return data
}

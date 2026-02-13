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
	URL       string `json:"url"`
	HasData   bool   `json:"hasData"`
	LineCount int    `json:"lineCount"`
	Error     bool   `json:"error"`
}

type PreviewResult struct {
	RawData  []Data      `json:"rawData"`
	Data     []Data      `json:"data"`
	Statuses []URLStatus `json:"statuses"`
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

func AddLines(data, lines []Data) []Data {
	if len(lines) == 0 {
		return data
	}
	for _, l := range lines {
		slog.Debug(fmt.Sprintf("adding line [%s: %s] as number %v", l.Name, l.Value, len(data)+1))
		data = append(data, l)
	}
	return data
}

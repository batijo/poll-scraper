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

// SumLines returns the sum line, plus a sum_symbol line when sumSymbols is set.
func SumLines(data []Data, sumSymbols string) []Data {
	var sum int
	for _, d := range data {
		v, err := strconv.Atoi(d.Value)
		if err != nil {
			slog.Debug(fmt.Sprintf("cannot convert value of [%s] to integer", d.Value), "err", err)
			continue
		}
		sum += v
	}
	lines := []Data{{Name: "sum", Value: strconv.Itoa(sum)}}
	if sumSymbols != "" {
		lines = append(lines, Data{Name: "sum_symbol", Value: strconv.Itoa(sum) + sumSymbols})
	}
	return lines
}

// AddSymbols inserts, right after each line whose 1-based position is a key
// of symbols, a copy of that line with the symbol appended to its value.
func AddSymbols(data []Data, symbols map[int]string) []Data {
	if len(symbols) == 0 {
		return data
	}
	result := make([]Data, 0, len(data)+len(symbols))
	for i, d := range data {
		result = append(result, d)
		if symbol, ok := symbols[i+1]; ok {
			result = append(result, Data{Name: d.Name + "_symbol", Value: d.Value + symbol})
		}
	}
	return result
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

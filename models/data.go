package models

import (
	"strconv"
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
			return data
		}
		sum += v
	}
	data = append(data, Data{Name: strconv.Itoa(len(data)), Value: strconv.Itoa(sum)})
	return data
}

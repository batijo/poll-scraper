package config

import "github.com/batijo/poll-scraper/models"

// Process turns raw scraped data into the output lines shared by the
// TXT/CSV writer, the HTTP API, and the preview.
func (c *Config) Process(raw []models.Data) []models.Data {
	data := raw
	if lines := c.FilterLinesZeroIndexed(); len(lines) > 0 {
		data = models.FilterData(lines, data)
	}

	var addLines []models.Data
	for _, l := range c.AddLines {
		if !l.Filtered {
			addLines = append(addLines, models.Data{Name: l.Name, Value: l.Value})
		}
	}
	data = models.AddLines(data, addLines)

	// The sum covers only original lines, so symbol copies are inserted after it is computed.
	var sumLines []models.Data
	if c.AddSum {
		sumLines = models.SumLines(data, c.SumSymbols)
	}

	symbols := make(map[int]string, len(c.SymbolLines))
	for _, s := range c.SymbolLines {
		symbols[s.Number] = s.Symbol
	}
	data = append(models.AddSymbols(data, symbols), sumLines...)

	if data == nil {
		return []models.Data{}
	}
	return data
}

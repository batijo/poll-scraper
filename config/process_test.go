package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/batijo/poll-scraper/models"
)

func TestProcessSymbolLines(t *testing.T) {
	cfg := defaultConfig()
	cfg.FilterLines = []int{1, 3}
	cfg.AddLines = []AddLine{
		{Name: "hidden", Value: "100", Filtered: true},
		{Name: "custom", Value: "5"},
	}
	cfg.SymbolLines = []SymbolLine{{Number: 2, Symbol: " eur"}, {Number: 3, Symbol: "%"}}
	cfg.AddSum = true
	cfg.SumSymbols = " $"
	raw := []models.Data{
		{Name: "a", Value: "10"},
		{Name: "b", Value: "20"},
		{Name: "c", Value: "30"},
	}

	got := cfg.Process(raw)

	// Symbol copies must not be summed: 10 + 30 + 5 = 45.
	want := []models.Data{
		{Name: "a", Value: "10"},
		{Name: "c", Value: "30"},
		{Name: "c_symbol", Value: "30 eur"},
		{Name: "custom", Value: "5"},
		{Name: "custom_symbol", Value: "5%"},
		{Name: "sum", Value: "45"},
		{Name: "sum_symbol", Value: "45 $"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestLoadRejectsInvalidSymbolLines(t *testing.T) {
	tests := map[string]string{
		"zero number": `[{"number": 0, "symbol": "x"}]`,
		"duplicate":   `[{"number": 2, "symbol": "x"}, {"number": 2, "symbol": "y"}]`,
	}
	for name, symbolLines := range tests {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			content := `{"port": 3000, "symbol_lines": ` + symbolLines + `}`
			if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
				t.Fatalf("failed to write config: %v", err)
			}

			_, err := Load(path)
			if err == nil || !strings.Contains(err.Error(), "symbol_lines") {
				t.Errorf("Load() error = %v, want symbol_lines error", err)
			}
		})
	}
}

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"

	"github.com/batijo/poll-scraper/utils"
)

const (
	defaultPort           = 3000
	defaultUpdateInterval = 1000
)

type AddLine struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Filtered bool   `json:"filtered"`
}

type Config struct {
	Links                 []string  `json:"links"`
	Port                  int       `json:"port"`
	IP                    string    `json:"ip"`
	Domains               []string  `json:"domains"`
	EnableServer          bool      `json:"enable_server"`
	WithEq                bool      `json:"with_eq"`
	FilterLines           []int     `json:"filter_lines"`
	AddLines              []AddLine `json:"add_lines"`
	AddSum                bool      `json:"add_sum"`
	SumSymbols            string    `json:"sum_symbols"`
	UpdateInterval        int       `json:"update_interval"`
	WriteToCSV            bool      `json:"write_to_csv"`
	CSVPath               string    `json:"csv_path"`
	WriteToTXT            bool      `json:"write_to_txt"`
	TXTPath               string    `json:"txt_path"`
	TXTEncoding           string    `json:"txt_encoding"`
	DatasetName           string    `json:"dataset_name"`
	Debug                 bool      `json:"debug"`
	StopOnLineCountChange bool      `json:"stop_on_line_count_change"`
}

func defaultConfig() *Config {
	return &Config{
		Links:          []string{},
		Port:           defaultPort,
		IP:             "localhost",
		Domains:        []string{},
		EnableServer:   true,
		FilterLines:    []int{},
		AddLines:       []AddLine{},
		UpdateInterval: defaultUpdateInterval,
	}
}

func Load(path string) (*Config, error) {
	slog.Debug("loading config", "path", path)
	cleanPath := filepath.Clean(path)
	file, err := os.Open(cleanPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Info("config file not found, creating default", "path", path)
			cfg := defaultConfig()
			if err := cfg.Save(path); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			slog.Error("failed to close config file", "err", cerr)
		}
	}()

	var cfg Config
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cfg.applyDefaults()
	cfg.sortFilters()
	slog.Debug("config loaded",
		"links", len(cfg.Links),
		"filter_lines", len(cfg.FilterLines),
		"add_lines", len(cfg.AddLines),
		"server", cfg.EnableServer,
	)
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if c.UpdateInterval < 0 {
		return fmt.Errorf("update_interval cannot be negative")
	}
	if c.WriteToCSV && c.CSVPath == "" {
		return fmt.Errorf("csv_path is required when write_to_csv is true")
	}
	if c.WriteToTXT && c.TXTPath == "" {
		return fmt.Errorf("txt_path is required when write_to_txt is true")
	}
	if c.WriteToTXT && c.DatasetName == "" {
		return fmt.Errorf("dataset_name is required when write_to_txt is true")
	}
	return nil
}

func (c *Config) Save(path string) error {
	c.applyDefaults()
	c.sortFilters()

	if err := c.validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	c.warnEmptyValues()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filepath.Clean(path), data, utils.FileMode); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) warnEmptyValues() {
	if len(c.Links) == 0 {
		slog.Warn("no URLs configured")
	}
	for i, line := range c.AddLines {
		if line.Name == "" || line.Value == "" {
			slog.Warn("custom line has empty field", "index", i, "name", line.Name, "value", line.Value)
		}
	}
	if c.AddSum && c.SumSymbols == "" {
		slog.Warn("add_sum enabled but sum_symbols is empty")
	}
}

func (c *Config) applyDefaults() {
	if c.IP == "" {
		c.IP = "localhost"
	}
	if c.UpdateInterval == 0 {
		c.UpdateInterval = defaultUpdateInterval
	}
}

func (c *Config) sortFilters() {
	sort.Ints(c.FilterLines)
}

// FilterLinesZeroIndexed returns filter lines converted to zero-based indices
func (c *Config) FilterLinesZeroIndexed() []int {
	if len(c.FilterLines) == 0 {
		return nil
	}
	result := make([]int, len(c.FilterLines))
	for i, line := range c.FilterLines {
		result[i] = line - 1
	}
	return result
}

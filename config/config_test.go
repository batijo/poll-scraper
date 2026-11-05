package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWithLinkCredentials(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{
		"links": [
			{"url": "https://example.com/data", "name": "user", "password": "pass"}
		],
		"port": 3000
	}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Links) != 1 {
		t.Fatalf("got %d links, want 1", len(cfg.Links))
	}
	link := cfg.Links[0]
	if link.URL != "https://example.com/data" || link.Name != "user" || link.Password != "pass" {
		t.Errorf("link = %+v, want url with user/pass", link)
	}
}

func TestSaveLoadLinksRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	cfg := defaultConfig()
	cfg.Links = []Link{
		{URL: "https://example.com/data", Name: "first", Password: "one"},
		{URL: "https://example.com/other"},
	}

	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(loaded.Links) != 2 {
		t.Fatalf("got %d links, want 2", len(loaded.Links))
	}
	if loaded.Links[0].Name != "first" || loaded.Links[0].Password != "one" {
		t.Errorf("link 0 = %+v, want first/one", loaded.Links[0])
	}
	if loaded.Links[1].URL != "https://example.com/other" {
		t.Errorf("link 1 = %+v, want url without credentials", loaded.Links[1])
	}
}

func TestLoadWithPartialCredentials(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{
		"links": [
			{"url": "https://example.com/data", "name": "user"},
			{"url": "https://example.com/other", "password": "pass"}
		],
		"port": 3000
	}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Links) != 2 {
		t.Fatalf("got %d links, want 2", len(cfg.Links))
	}
	if cfg.Links[0].Name != "user" || cfg.Links[1].Password != "pass" {
		t.Errorf("links = %+v, want partial credentials preserved", cfg.Links)
	}
}

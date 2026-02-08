package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/batijo/poll-scraper/config"
	"github.com/batijo/poll-scraper/models"
	"github.com/batijo/poll-scraper/scraper"
	"github.com/batijo/poll-scraper/server"
	"github.com/batijo/poll-scraper/utils"
	"github.com/batijo/poll-scraper/utils/file"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx            context.Context
	cfg            *config.Config
	srv            *server.Server
	stopWriter     context.CancelFunc
	scraperRunning bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	cfg, err := config.Load("config.json")
	if err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("failed to load config: %v", err))
		return
	}
	a.cfg = cfg

	if err := a.initLogger(cfg.Debug); err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("failed to init logger: %v", err))
		return
	}

	slog.Info("config loaded",
		"urls", len(cfg.Links),
		"server_enabled", cfg.EnableServer,
		"interval", cfg.UpdateInterval,
		"debug", cfg.Debug,
	)

	if err := file.InitFiles(cfg); err != nil {
		slog.Error("failed to init files", "err", err)
	}

	a.scraperRunning = false
	slog.Info("application started", "scraper_state", "stopped")
}

func (a *App) Shutdown(ctx context.Context) {
	slog.Info("application shutting down")
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, from Go!", name)
}

func (a *App) GetConfig() *config.Config {
	return a.cfg
}

func (a *App) UpdateConfig(cfg config.Config) error {
	slog.Info("config update requested")

	if err := cfg.Save("config.json"); err != nil {
		slog.Error("failed to save config", "err", err)
		return err
	}
	slog.Debug("config saved to disk")

	oldCfg := a.cfg
	a.cfg = &cfg

	// Log what changed
	a.logConfigChanges(oldCfg, &cfg)

	// Restart scraper only if it was running
	wasRunning := a.scraperRunning
	if wasRunning {
		a.StopScraper()
	}

	// Restart server if enabled, stop if disabled
	a.stopServer()
	if cfg.EnableServer {
		a.startServer()
	}

	// Reinit output files if paths or toggles changed
	if oldCfg.WriteToCSV != cfg.WriteToCSV || oldCfg.CSVPath != cfg.CSVPath ||
		oldCfg.WriteToTXT != cfg.WriteToTXT || oldCfg.TXTPath != cfg.TXTPath {
		slog.Debug("output file config changed, reinitializing")
		if err := file.InitFiles(a.cfg); err != nil {
			slog.Error("failed to reinit files", "err", err)
		}
	}

	// Reinit logger if debug mode changed
	if oldCfg.Debug != cfg.Debug {
		if err := a.initLogger(cfg.Debug); err != nil {
			slog.Error("failed to reinit logger", "err", err)
		}
	}

	// Restart scraper if it was running before config update
	if wasRunning {
		if err := a.StartScraper(); err != nil {
			slog.Error("failed to restart scraper after config update", "err", err)
			return err
		}
	}

	slog.Info("config updated successfully")
	return nil
}

func (a *App) StartScraper() error {
	if a.scraperRunning {
		return fmt.Errorf("scraper is already running")
	}
	slog.Info("starting scraper")

	if a.cfg.EnableServer {
		a.startServer()
	}

	stopWriter, err := file.StartWriting(a.cfg, a)
	if err != nil {
		slog.Error("failed to start scraper", "err", err)
		a.stopServer()
		return err
	}
	a.stopWriter = stopWriter
	a.scraperRunning = true
	return nil
}

func (a *App) StopScraper() {
	if !a.scraperRunning {
		return
	}
	slog.Info("stopping scraper")
	if a.stopWriter != nil {
		a.stopWriter()
		a.stopWriter = nil
	}
	a.stopServer()
	a.scraperRunning = false
	a.EmitScraperState("stopped")
}

func (a *App) IsScraperRunning() bool {
	return a.scraperRunning
}

func (a *App) RequestScraperStop() {
	go a.StopScraper()
}

func (a *App) PreviewScrape() models.PreviewResult {
	slog.Info("preview scrape requested")

	lines := a.cfg.FilterLinesZeroIndexed()
	var rawData []models.Data
	statuses := make([]models.URLStatus, 0, len(a.cfg.Links))

	for _, link := range a.cfg.Links {
		urlData := scraper.ScrapeURL(link, a.cfg.WithEq)
		statuses = append(statuses, models.URLStatus{
			URL:       link,
			HasData:   len(urlData) > 0,
			LineCount: len(urlData),
		})
		if len(urlData) == 0 {
			slog.Warn("no data from URL", "url", link)
		}
		rawData = append(rawData, urlData...)
	}

	// rawData = URL-scraped data only (for frontend filter modal)
	processedData := make([]models.Data, len(rawData))
	copy(processedData, rawData)

	if len(lines) > 0 {
		processedData = models.FilterData(lines, processedData)
	}

	// Add non-filtered custom lines after URL filtering
	if len(a.cfg.AddLines) > 0 {
		var addLines []models.Data
		for _, l := range a.cfg.AddLines {
			if !l.Filtered {
				addLines = append(addLines, models.Data{Name: l.Name, Value: l.Value})
			}
		}
		if len(addLines) > 0 {
			processedData = models.AddLines(processedData, addLines)
		}
	}
	if a.cfg.AddSum {
		processedData = models.SumData(processedData, a.cfg.SumSymbols)
	}

	if rawData == nil {
		rawData = []models.Data{}
	}
	if processedData == nil {
		processedData = []models.Data{}
	}

	slog.Info("preview scrape complete", "raw_lines", len(rawData), "processed_lines", len(processedData))

	return models.PreviewResult{
		RawData:  rawData,
		Data:     processedData,
		Statuses: statuses,
	}
}

func (a *App) logConfigChanges(oldCfg, newCfg *config.Config) {
	if oldCfg.UpdateInterval != newCfg.UpdateInterval {
		slog.Info("config changed", "field", "update_interval", "old", oldCfg.UpdateInterval, "new", newCfg.UpdateInterval)
	}
	if oldCfg.EnableServer != newCfg.EnableServer {
		slog.Info("config changed", "field", "enable_server", "old", oldCfg.EnableServer, "new", newCfg.EnableServer)
	}
	if oldCfg.IP != newCfg.IP {
		slog.Info("config changed", "field", "ip", "old", oldCfg.IP, "new", newCfg.IP)
	}
	if oldCfg.Port != newCfg.Port {
		slog.Info("config changed", "field", "port", "old", oldCfg.Port, "new", newCfg.Port)
	}
	if oldCfg.WithEq != newCfg.WithEq {
		slog.Info("config changed", "field", "with_eq", "old", oldCfg.WithEq, "new", newCfg.WithEq)
	}
	if oldCfg.WriteToCSV != newCfg.WriteToCSV {
		slog.Info("config changed", "field", "write_to_csv", "old", oldCfg.WriteToCSV, "new", newCfg.WriteToCSV)
	}
	if oldCfg.CSVPath != newCfg.CSVPath {
		slog.Info("config changed", "field", "csv_path", "old", oldCfg.CSVPath, "new", newCfg.CSVPath)
	}
	if oldCfg.WriteToTXT != newCfg.WriteToTXT {
		slog.Info("config changed", "field", "write_to_txt", "old", oldCfg.WriteToTXT, "new", newCfg.WriteToTXT)
	}
	if oldCfg.TXTPath != newCfg.TXTPath {
		slog.Info("config changed", "field", "txt_path", "old", oldCfg.TXTPath, "new", newCfg.TXTPath)
	}
	if oldCfg.DatasetName != newCfg.DatasetName {
		slog.Info("config changed", "field", "dataset_name", "old", oldCfg.DatasetName, "new", newCfg.DatasetName)
	}
	if oldCfg.Debug != newCfg.Debug {
		slog.Info("config changed", "field", "debug", "old", oldCfg.Debug, "new", newCfg.Debug)
	}
	if oldCfg.AddSum != newCfg.AddSum {
		slog.Info("config changed", "field", "add_sum", "old", oldCfg.AddSum, "new", newCfg.AddSum)
	}
	if oldCfg.SumSymbols != newCfg.SumSymbols {
		slog.Info("config changed", "field", "sum_symbols", "old", oldCfg.SumSymbols, "new", newCfg.SumSymbols)
	}
	if !reflect.DeepEqual(oldCfg.Links, newCfg.Links) {
		slog.Info("config changed", "field", "links", "old_count", len(oldCfg.Links), "new_count", len(newCfg.Links))
	}
	if !reflect.DeepEqual(oldCfg.Domains, newCfg.Domains) {
		slog.Info("config changed", "field", "domains", "old_count", len(oldCfg.Domains), "new_count", len(newCfg.Domains))
	}
	if !reflect.DeepEqual(oldCfg.FilterLines, newCfg.FilterLines) {
		slog.Info("config changed", "field", "filter_lines", "old_count", len(oldCfg.FilterLines), "new_count", len(newCfg.FilterLines))
	}
	if !reflect.DeepEqual(oldCfg.AddLines, newCfg.AddLines) {
		slog.Info("config changed", "field", "add_lines", "old_count", len(oldCfg.AddLines), "new_count", len(newCfg.AddLines))
	}
	if oldCfg.StopOnLineCountChange != newCfg.StopOnLineCountChange {
		slog.Info("config changed", "field", "stop_on_line_count_change", "old", oldCfg.StopOnLineCountChange, "new", newCfg.StopOnLineCountChange)
	}
}

func (a *App) EmitScraperData(data []models.Data, rawData []models.Data) {
	if data == nil {
		data = []models.Data{}
	}
	if rawData == nil {
		rawData = []models.Data{}
	}
	payload := map[string]interface{}{
		"data":      data,
		"rawData":   rawData,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	runtime.EventsEmit(a.ctx, "polled:data", payload)
}

func (a *App) EmitScraperState(state string) {
	runtime.EventsEmit(a.ctx, "polled:state", state)
}

func (a *App) PreviewURL(url string) []models.Data {
	slog.Debug("previewing URL", "url", url)
	data := scraper.ScrapeURL(url, a.cfg.WithEq)
	slog.Debug("preview complete", "url", url, "lines", len(data))
	return data
}

func (a *App) EmitURLStatus(statuses []models.URLStatus) {
	runtime.EventsEmit(a.ctx, "polled:url-status", statuses)
}

func (a *App) EmitLog(entry utils.LogEntry) {
	runtime.EventsEmit(a.ctx, "polled:log", entry)
}

func (a *App) EmitScraperError(message string) {
	payload := map[string]interface{}{
		"message":   message,
		"timestamp": time.Now().Unix(),
	}
	runtime.EventsEmit(a.ctx, "polled:error", payload)
}

func (a *App) startServer() {
	srv := server.New(a.cfg)
	srv.Addr = fmt.Sprintf("%s:%d", a.cfg.IP, a.cfg.Port)
	a.srv = srv
	go func() {
		slog.Info("server starting", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server stopped", "err", err)
		}
	}()
}

func (a *App) stopServer() {
	if a.srv != nil {
		slog.Debug("stopping HTTP server")
		if err := a.srv.Close(); err != nil {
			slog.Error("failed to stop server", "err", err)
		}
		a.srv = nil
	}
}

func (a *App) initLogger(debug bool) error {
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, utils.FileMode)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	var inner slog.Handler
	if debug {
		inner = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		inner = slog.NewTextHandler(logFile, &slog.HandlerOptions{
			Level:     slog.LevelError,
			AddSource: true,
		})
	}

	// Wrap with frontend handler to forward logs to the UI
	handler := utils.NewFrontendHandler(inner, a, debug)
	slog.SetDefault(slog.New(handler))

	return nil
}

package file

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/batijo/poll-scraper/config"
	"github.com/batijo/poll-scraper/models"
)

type recordingEmitter struct {
	mu            sync.Mutex
	errors        []string
	stopRequested bool
	data          chan []models.Data
	failures      chan string
}

func newRecordingEmitter() *recordingEmitter {
	return &recordingEmitter{data: make(chan []models.Data, 16), failures: make(chan string, 16)}
}

// Sends never block, so a slow test cannot stall the writer goroutine.
func send[T any](ch chan<- T, v T) {
	select {
	case ch <- v:
	default:
	}
}

func (e *recordingEmitter) EmitScraperData(data, _ []models.Data) { send(e.data, data) }
func (e *recordingEmitter) EmitScraperState(string)               {}
func (e *recordingEmitter) EmitURLStatus([]models.URLStatus)      {}

func (e *recordingEmitter) EmitScraperError(message string) {
	e.mu.Lock()
	e.errors = append(e.errors, message)
	e.mu.Unlock()
	send(e.failures, message)
}

func (e *recordingEmitter) RequestScraperStop() {
	e.mu.Lock()
	e.stopRequested = true
	e.mu.Unlock()
}

func receive[T any](t *testing.T, ch <-chan T, what string) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
		var zero T
		return zero
	}
}

func TestWriter_FailedScrapeKeepsLastOutputAndRecovers(t *testing.T) {
	var failing atomic.Bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if failing.Load() {
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(`<p>1=10</p><p>2=20</p>`))
	}))
	t.Cleanup(ts.Close)

	txtPath := filepath.Join(t.TempDir(), "data.txt")
	cfg := &config.Config{
		Links:                 []config.Link{{URL: ts.URL}},
		WithEq:                true,
		WriteToTXT:            true,
		TXTPath:               txtPath,
		TXTEncoding:           "utf-8",
		DatasetName:           "data",
		UpdateInterval:        10,
		StopOnLineCountChange: true,
	}
	emitter := newRecordingEmitter()
	cancel, err := StartWriting(cfg, emitter)
	if err != nil {
		t.Fatalf("StartWriting() error = %v", err)
	}
	t.Cleanup(cancel)

	receive(t, emitter.data, "first successful cycle")
	good, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	failing.Store(true)
	receive(t, emitter.failures, "first failed cycle")
	receive(t, emitter.failures, "second failed cycle")
	afterFailure, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	if string(afterFailure) != string(good) {
		t.Errorf("output changed during failure:\n%s\nwant:\n%s", afterFailure, good)
	}

	failing.Store(false)
	for len(emitter.data) > 0 {
		<-emitter.data
	}
	recovered := receive(t, emitter.data, "recovered cycle")
	if len(recovered) != 2 {
		t.Errorf("recovered cycle got %d lines, want 2", len(recovered))
	}

	emitter.mu.Lock()
	defer emitter.mu.Unlock()
	if emitter.stopRequested {
		t.Error("scraper stopped on a failed request, want it to keep retrying")
	}
	for _, msg := range emitter.errors {
		if strings.Contains(msg, "line count changed") {
			t.Errorf("failed request reported as a line count change: %q", msg)
		}
	}
}

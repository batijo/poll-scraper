package utils

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// LogEmitter sends log entries to the frontend.
type LogEmitter interface {
	EmitLog(entry LogEntry)
}

type LogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

// FrontendHandler is a slog.Handler that forwards records to the frontend
// while delegating to an inner handler for normal logging.
type FrontendHandler struct {
	inner    slog.Handler
	emitter  LogEmitter
	minLevel slog.Level
}

func NewFrontendHandler(inner slog.Handler, emitter LogEmitter, debug bool) *FrontendHandler {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	return &FrontendHandler{inner: inner, emitter: emitter, minLevel: level}
}

func (h *FrontendHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.minLevel
}

//nolint:gocritic // slog.Handler interface requires value receiver for slog.Record
func (h *FrontendHandler) Handle(ctx context.Context, r slog.Record) error {
	var attrs []string
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value))
		return true
	})

	msg := r.Message
	if len(attrs) > 0 {
		msg += " " + strings.Join(attrs, " ")
	}

	// Always emit to frontend regardless of inner handler level
	h.emitter.EmitLog(LogEntry{
		Level:   r.Level.String(),
		Message: msg,
		Time:    r.Time.Format("15:04:05"),
	})

	// Inner handler applies its own level filtering
	if h.inner.Enabled(ctx, r.Level) {
		return h.inner.Handle(ctx, r)
	}
	return nil
}

func (h *FrontendHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &FrontendHandler{inner: h.inner.WithAttrs(attrs), emitter: h.emitter, minLevel: h.minLevel}
}

func (h *FrontendHandler) WithGroup(name string) slog.Handler {
	return &FrontendHandler{inner: h.inner.WithGroup(name), emitter: h.emitter, minLevel: h.minLevel}
}

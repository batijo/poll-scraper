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
	inner   slog.Handler
	emitter LogEmitter
}

func NewFrontendHandler(inner slog.Handler, emitter LogEmitter) *FrontendHandler {
	return &FrontendHandler{inner: inner, emitter: emitter}
}

func (h *FrontendHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *FrontendHandler) Handle(ctx context.Context, r slog.Record) error {
	// Build attrs string
	var attrs []string
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value))
		return true
	})

	msg := r.Message
	if len(attrs) > 0 {
		msg += " " + strings.Join(attrs, " ")
	}

	h.emitter.EmitLog(LogEntry{
		Level:   r.Level.String(),
		Message: msg,
		Time:    r.Time.Format("15:04:05"),
	})

	return h.inner.Handle(ctx, r)
}

func (h *FrontendHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &FrontendHandler{inner: h.inner.WithAttrs(attrs), emitter: h.emitter}
}

func (h *FrontendHandler) WithGroup(name string) slog.Handler {
	return &FrontendHandler{inner: h.inner.WithGroup(name), emitter: h.emitter}
}

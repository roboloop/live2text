package logger

import (
	"context"
	"log/slog"
	"sync"
)

type CaptureHandler struct {
	logs  []LogEntry
	attrs []slog.Attr
	mu    sync.Mutex
}

type LogEntry struct {
	Level slog.Level
	Msg   string
	Attrs []slog.Attr
}

func (ch *CaptureHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (ch *CaptureHandler) Handle(_ context.Context, r slog.Record) error {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	var attrs []slog.Attr
	attrs = append(attrs, ch.attrs...)
	r.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})

	entry := LogEntry{
		Level: r.Level,
		Msg:   r.Message,
		Attrs: attrs,
	}

	ch.logs = append(ch.logs, entry)

	return nil
}

func (ch *CaptureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.attrs = append(ch.attrs, attrs...)
	return ch
}

func (ch *CaptureHandler) WithGroup(string) slog.Handler {
	return ch
}

func (ch *CaptureHandler) GetLog(msg string) (LogEntry, bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for _, log := range ch.logs {
		if log.Msg == msg {
			return log, true
		}
	}
	return LogEntry{}, false
}

func (ch *CaptureHandler) Logs() []LogEntry {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	logs := make([]LogEntry, len(ch.logs))
	copy(logs, ch.logs)

	return logs
}

func (e LogEntry) GetAttr(key string) (any, bool) {
	for _, attr := range e.Attrs {
		if attr.Key == key {
			return attr.Value.Any(), true
		}
	}
	return nil, false
}

func NewCaptureLogger() (*slog.Logger, *CaptureHandler) {
	captureHandler := &CaptureHandler{}
	return slog.New(captureHandler), captureHandler
}

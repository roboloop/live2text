package logger

import (
	"context"
	"log/slog"
)

type CaptureHandler struct {
	Logs  []LogEntry
	attrs []slog.Attr
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
	var attrs []slog.Attr
	// attrs := make()
	attrs = append(attrs, ch.attrs...)
	r.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})
	ch.Logs = append(ch.Logs, LogEntry{
		Level: r.Level,
		Msg:   r.Message,
		Attrs: attrs,
	})

	return nil
}

func (ch *CaptureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	ch.attrs = append(ch.attrs, attrs...)
	return ch
}

func (ch *CaptureHandler) WithGroup(string) slog.Handler {
	return ch
}

func (ch *CaptureHandler) GetLog(msg string) (LogEntry, bool) {
	for _, log := range ch.Logs {
		if log.Msg == msg {
			return log, true
		}
	}
	return LogEntry{}, false
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

package utils

import (
	"context"
	"log/slog"
)

type TestHandler struct {
	Logs []LogEntry
}

type LogEntry struct {
	Level slog.Level
	Msg   string
	Attrs []slog.Attr
}

func (t *TestHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (t *TestHandler) Handle(_ context.Context, r slog.Record) error {
	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	t.Logs = append(t.Logs, LogEntry{
		Level: r.Level,
		Msg:   r.Message,
		Attrs: attrs,
	})

	return nil
}

func (t *TestHandler) WithAttrs([]slog.Attr) slog.Handler {
	return t
}

func (t *TestHandler) WithGroup(string) slog.Handler {
	return t
}

func NewTestLogger() (*slog.Logger, *TestHandler) {
	testHandler := &TestHandler{}
	return slog.New(testHandler), testHandler
}

package logger

import (
	"context"
	"log/slog"
)

type nilHandler struct{}

func (n *nilHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (n *nilHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (n *nilHandler) WithAttrs([]slog.Attr) slog.Handler {
	return n
}

func (n *nilHandler) WithGroup(string) slog.Handler {
	return n
}

var NilLogger = slog.New(&nilHandler{})

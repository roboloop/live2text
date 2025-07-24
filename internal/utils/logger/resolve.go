package logger

import (
	"context"
	"errors"
	"log/slog"
)

func ResolveLevel(err error) slog.Level {
	if err == nil {
		return slog.LevelInfo
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return slog.LevelInfo
	}

	return slog.LevelError
}

package logger_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestResolveLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    error
		expect slog.Level
	}{
		{
			name:   "nil",
			err:    nil,
			expect: slog.LevelInfo,
		},
		{
			name:   "context canceled error",
			err:    context.Canceled,
			expect: slog.LevelInfo,
		},
		{
			name:   "context deadline exceeded error",
			err:    context.DeadlineExceeded,
			expect: slog.LevelInfo,
		},
		{
			name:   "other error",
			err:    errors.New("something happened"),
			expect: slog.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			level := logger.ResolveLevel(tt.err)
			require.Equal(t, tt.expect, level)
		})
	}
}

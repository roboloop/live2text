package logger_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestNilLogger(t *testing.T) {
	t.Parallel()

	nilLogger := logger.NilLogger

	nilLogger.InfoContext(t.Context(), "msg")

	assert.True(t, nilLogger.Enabled(t.Context(), slog.LevelInfo))
	assert.Equal(t, nilLogger, nilLogger.WithGroup("foo"))
	assert.Equal(t, nilLogger, nilLogger.With("key", "value"))
}

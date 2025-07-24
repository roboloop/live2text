package logger_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"live2text/internal/utils/logger"
)

func TestCaptureLogger(t *testing.T) {
	t.Parallel()

	captureLogger, captureHandler := logger.NewCaptureLogger()

	assert.True(t, captureLogger.Enabled(t.Context(), slog.LevelInfo))
	assert.Equal(t, captureLogger, captureLogger.WithGroup("foo"))
	assert.Equal(t, captureLogger, captureLogger.With("key", "value"))

	captureLogger.InfoContext(t.Context(), "this is a message")
	require.Len(t, captureHandler.Logs, 1)
	require.Equal(t, "this is a message", captureHandler.Logs[0].Msg)
	require.Equal(t, slog.LevelInfo, captureHandler.Logs[0].Level)

	attrs := captureHandler.Logs[0].Attrs
	require.Len(t, attrs, 1)
	require.Equal(t, "key", attrs[0].Key)
	require.Equal(t, "value", attrs[0].Value.String())

	log, ok := captureHandler.GetLog("this is a message")
	require.True(t, ok)
	require.Equal(t, captureHandler.Logs[0], log)
	value, ok := log.GetAttr("key")
	require.True(t, ok)
	require.Equal(t, "value", value)

	log, ok = captureHandler.GetLog("non existing message")
	require.False(t, ok)
	require.Empty(t, log)
	entry, ok := log.GetAttr("key")
	require.False(t, ok)
	require.Nil(t, entry)
}

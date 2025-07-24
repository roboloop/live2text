package middleware_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"live2text/internal/api/middleware"
	"live2text/internal/utils/logger"
)

func TestLoggerMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("logs attributes on error", func(t *testing.T) {
		t.Parallel()

		l, h := logger.NewCaptureLogger()
		m := middleware.LoggerMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusGatewayTimeout)
			_, _ = w.Write([]byte("error message"))
		}), l)

		req := httptest.NewRequest(http.MethodPost, "/foo", nil)
		req.Header.Add("User-Agent", "foo agent")
		req.RemoteAddr = "127.0.0.1:1234"
		w := httptest.NewRecorder()

		m.ServeHTTP(w, req)

		require.Len(t, h.Logs, 1)
		entry := h.Logs[0]
		require.Equal(t, slog.LevelError, entry.Level)
		require.Equal(t, "HTTP request", entry.Msg)

		assertAttr(t, entry.Attrs, "method", http.MethodPost)
		assertAttr(t, entry.Attrs, "path", "/foo")
		assertAttr(t, entry.Attrs, "status", http.StatusGatewayTimeout)
		assertAttr(t, entry.Attrs, "IP", "127.0.0.1")
		assertAttr(t, entry.Attrs, "user_agent", "foo agent")
		assertAttr(t, entry.Attrs, "error", "error message")
	})
}

func assertAttr(t *testing.T, attrs []slog.Attr, expectedKey string, expectedValue any) {
	t.Helper()

	i := slices.IndexFunc(attrs, func(attr slog.Attr) bool {
		return attr.Key == expectedKey
	})
	if i == -1 {
		require.FailNow(t, "Not contains key:", expectedKey)
	}

	require.Truef(
		t,
		attrs[i].Equal(slog.Any(expectedKey, expectedValue)),
		"Expected on key '%s' on value '%v', got '%v'",
		expectedKey,
		expectedValue,
		attrs[i].Value,
	)
}

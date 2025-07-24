package middleware_test

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"live2text/internal/api/middleware"
	"live2text/internal/utils/logger"
)

type errorCloser struct{}

func (e *errorCloser) Read([]byte) (int, error) {
	return 0, nil
}

func (e *errorCloser) Close() error {
	return errors.New("dummy error")
}

func TestBodyCloser(t *testing.T) {
	t.Parallel()

	t.Run("logs error on body close", func(t *testing.T) {
		t.Parallel()

		l, h := logger.NewCaptureLogger()
		m := middleware.BodyCloserMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// nothing
		}), l)

		req := httptest.NewRequest(http.MethodGet, "/", &errorCloser{})
		w := httptest.NewRecorder()

		m.ServeHTTP(w, req)

		require.Len(t, h.Logs, 1)
		require.True(t, slices.ContainsFunc(h.Logs[0].Attrs, func(attr slog.Attr) bool {
			return attr.Key == "error" && attr.Value.String() == "dummy error"
		}))
	})
}

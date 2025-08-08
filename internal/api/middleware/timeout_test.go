package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/api/middleware"
)

func TestTimeoutMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("bypasses timeout when path is excluded", func(t *testing.T) {
		t.Parallel()

		called := false
		h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			called = true
		})
		m := middleware.TimeoutMiddleware(h, 0*time.Minute, []string{"/foo"}, false)
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, req)
		require.True(t, called)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("bypasses timeout when debug mode is enabled", func(t *testing.T) {
		t.Parallel()

		called := false
		h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			called = true
		})
		m := middleware.TimeoutMiddleware(h, 1*time.Nanosecond, nil, true)
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, req)
		require.True(t, called)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("returns 503 when handler exceeds timeout", func(t *testing.T) {
		t.Parallel()

		h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			time.Sleep(100 * time.Millisecond)
		})
		m := middleware.TimeoutMiddleware(h, 10*time.Millisecond, nil, false)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	})

	t.Run("responds with 200 when handler completes in time", func(t *testing.T) {
		t.Parallel()

		h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			// nothing
		})
		m := middleware.TimeoutMiddleware(h, 10*time.Millisecond, nil, false)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

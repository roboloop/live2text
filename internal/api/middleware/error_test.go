package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/api/middleware"
)

func TestErrorMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("capture panic", func(t *testing.T) {
		t.Parallel()

		m := middleware.ErrorMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("foo")
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		require.Contains(t, w.Body.String(), "Internal Server Error: foo")
	})
}

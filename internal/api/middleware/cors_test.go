package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"live2text/internal/api/middleware"
)

func TestCORSMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("sets CORS headers", func(t *testing.T) {
		t.Parallel()

		m := middleware.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, r)

		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("responds with 204 and skips handler on OPTIONS request", func(t *testing.T) {
		t.Parallel()

		dummyHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {})
		assert.HTTPStatusCode(
			t,
			middleware.CORSMiddleware(dummyHandler),
			http.MethodOptions,
			"/",
			nil,
			http.StatusNoContent,
		)
	})
}

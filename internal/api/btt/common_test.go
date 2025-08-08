package btt_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	bttapi "github.com/roboloop/live2text/internal/api/btt"
	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func setupServer(
	t *testing.T,
	setupMocks func(*minimock.Controller, *services.ServicesMock),
	l *slog.Logger,
) *bttapi.Server {
	t.Helper()

	mc := minimock.NewController(t)
	s := services.NewServicesMock(mc)

	if setupMocks != nil {
		setupMocks(mc, s)
	}
	if l == nil {
		l = logger.NilLogger
	}

	return bttapi.NewServer(l, s)
}

func performRequest(t *testing.T, handler http.HandlerFunc, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	return w
}

func assertResponse(
	t *testing.T,
	w *httptest.ResponseRecorder,
	expectedCode int,
	expectedHeaders map[string]string,
	expectedBody string,
) {
	t.Helper()

	require.Equal(t, expectedCode, w.Code)
	for key, expectedValue := range expectedHeaders {
		actualValue := w.Header().Get(key)
		require.Equalf(t, expectedValue, actualValue, "header %q: expected %q, got %q", key, expectedValue, actualValue)
	}
	if expectedBody != "" {
		require.Contains(t, w.Body.String(), expectedBody)
	}
}

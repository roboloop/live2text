package core_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/api/core"
	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	t.Run("returns ok", func(t *testing.T) {
		t.Parallel()

		mc := minimock.NewController(t)
		s := services.NewServicesMock(mc)
		server := core.NewServer(logger.NilLogger, s)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		server.Health(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), `"ok"`)
	})
}

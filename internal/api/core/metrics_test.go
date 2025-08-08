package core_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/api/core"
	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/services/metrics"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestMetrics(t *testing.T) {
	t.Parallel()

	t.Run("returns metrics", func(t *testing.T) {
		t.Parallel()

		mc := minimock.NewController(t)
		m := metrics.NewMetricsMock(mc)
		m.WritePrometheusMock.Times(1).Return()
		s := services.NewServicesMock(mc)
		s.MetricsMock.Return(m)
		server := core.NewServer(logger.NilLogger, s)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		server.Metrics(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

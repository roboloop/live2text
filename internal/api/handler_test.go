package api_test

import (
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"live2text/internal/api"
	"live2text/internal/services"
	"live2text/internal/utils/logger"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		mc := minimock.NewController(t)
		s := services.NewServicesMock(mc)

		handler := api.NewHandler(logger.NilLogger, s)

		_, ok := handler.(http.HandlerFunc)
		require.True(t, ok)
	})
}

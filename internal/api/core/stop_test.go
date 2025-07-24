package core_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"live2text/internal/api/core"
	"live2text/internal/services"
	"live2text/internal/services/recognition"
	"live2text/internal/utils/logger"
)

func TestStop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		body     string
		mockStop func() error

		expectedCode int
		expectedBody string
	}{
		{
			name:         "bad request",
			body:         "invalid json",
			expectedCode: http.StatusBadRequest,
			expectedBody: "cannot decode request",
		},
		{
			name: "no device busy",
			body: `{"id":"foo"}`,
			mockStop: func() error {
				return recognition.ErrNoDeviceBusy
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "no device busy",
		},
		{
			name: "device search error",
			body: `{"id":"foo"}`,
			mockStop: func() error {
				return errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			body: `{"id":"foo"}`,
			mockStop: func() error {
				return nil
			},
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			s := services.NewServicesMock(mc)
			if tt.mockStop != nil {
				r := recognition.NewRecognitionMock(mc)
				r.StopMock.Return(tt.mockStop())
				s.RecognitionMock.Return(r)
			}

			server := core.NewServer(logger.NilLogger, s)
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.Stop(w, req)

			require.Equal(t, tt.expectedCode, w.Code)
			require.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"live2text/internal/services"
	"live2text/internal/services/btt"
)

func TestIsRunning(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockIsRunning func() (bool, error)

		expectedCode int
		expectedBody string
	}{
		{
			name: "getting status error",
			mockIsRunning: func() (bool, error) {
				return false, errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			mockIsRunning: func() (bool, error) {
				return true, nil
			},
			expectedCode: http.StatusOK,
			expectedBody: `true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := setupServer(t, func(mc *minimock.Controller, s *services.ServicesMock) {
				b := btt.NewBttMock(mc)
				b.IsRunningMock.Return(tt.mockIsRunning())
				s.BttMock.Return(b)
			}, nil)

			w := performRequest(t, server.IsRunning, "")

			assertResponse(t, w, tt.expectedCode, nil, tt.expectedBody)
		})
	}
}

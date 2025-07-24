package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"live2text/internal/services"
	"live2text/internal/services/btt"
)

func TestToggleListening(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		mockToggleListening func() error

		expectedCode int
		expectedBody string
	}{
		{
			name: "toggle listening failed",
			mockToggleListening: func() error {
				return errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			mockToggleListening: func() error {
				return nil
			},
			expectedCode: http.StatusOK,
			expectedBody: `"ok"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := setupServer(t, func(mc *minimock.Controller, s *services.ServicesMock) {
				b := btt.NewBttMock(mc)
				b.ToggleListeningMock.Return(tt.mockToggleListening())
				s.BttMock.Return(b)
			}, nil)

			w := performRequest(t, server.ToggleListening, "")

			assertResponse(t, w, tt.expectedCode, nil, tt.expectedBody)
		})
	}
}

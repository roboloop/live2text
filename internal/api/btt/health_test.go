package btt_test

import (
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/services/btt"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		mockHealth func() bool

		expectedCode    int
		expectedHeaders map[string]string
		expectedBody    string
	}{
		{
			name: "btt is not running",
			mockHealth: func() bool {
				return false
			},
			expectedCode:    http.StatusInternalServerError,
			expectedHeaders: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			expectedBody:    "btt is not running",
		},
		{
			name: "btt is running",
			mockHealth: func() bool {
				return true
			},
			expectedCode:    http.StatusOK,
			expectedHeaders: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			expectedBody:    "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := setupServer(t, func(mc *minimock.Controller, s *services.ServicesMock) {
				b := btt.NewBttMock(mc)
				b.HealthMock.Return(tt.mockHealth())
				s.BttMock.Return(b)
			}, nil)

			w := performRequest(t, server.Health, "")

			assertResponse(t, w, tt.expectedCode, tt.expectedHeaders, tt.expectedBody)
		})
	}
}

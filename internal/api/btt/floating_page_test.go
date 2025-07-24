package btt_test

import (
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"live2text/internal/services"
	"live2text/internal/services/btt"
)

func TestFloatingPage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mockPage func() string

		expectedCode    int
		expectedHeaders map[string]string
		expectedBody    string
	}{
		{
			name: "ok",
			mockPage: func() string {
				return "sample page"
			},
			expectedCode: http.StatusOK,
			expectedHeaders: map[string]string{
				"Content-Type":  "text/html; charset=utf-8",
				"Cache-Control": "no-cache",
			},
			expectedBody: "sample page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := setupServer(t, func(mc *minimock.Controller, s *services.ServicesMock) {
				b := btt.NewBttMock(mc)
				b.FloatingPageMock.Return(tt.mockPage())
				s.BttMock.Return(b)
			}, nil)

			w := performRequest(t, server.FloatingPage, "")

			assertResponse(t, w, tt.expectedCode, tt.expectedHeaders, tt.expectedBody)
		})
	}
}

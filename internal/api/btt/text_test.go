package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/services/btt"
)

func TestText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mockText func() (string, error)

		expectedCode int
		expectedBody string
	}{
		{
			name: "toggle listening error",
			mockText: func() (string, error) {
				return "", errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			mockText: func() (string, error) {
				return "sample text", nil
			},
			expectedCode: http.StatusOK,
			expectedBody: "sample text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := setupServer(t, func(mc *minimock.Controller, s *services.ServicesMock) {
				b := btt.NewBttMock(mc)
				b.TextMock.Return(tt.mockText())
				s.BttMock.Return(b)
			}, nil)

			w := performRequest(t, server.Text, "")

			assertResponse(t, w, tt.expectedCode, nil, tt.expectedBody)
		})
	}
}

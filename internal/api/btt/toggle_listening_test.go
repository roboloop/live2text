package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/services/btt"
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
			name: "toggle listening error",
			mockToggleListening: func() error {
				return errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "device is not selected",
			mockToggleListening: func() error {
				return btt.ErrDeviceNotSelected
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "device is not selected",
		},
		{
			name: "device is unavailable",
			mockToggleListening: func() error {
				return btt.ErrDeviceIsUnavailable
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "device is unavailable",
		},
		{
			name: "language is not selected",
			mockToggleListening: func() error {
				return btt.ErrLanguageNotSelected
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "language is not selected",
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

package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"live2text/internal/services"
	"live2text/internal/services/btt"
)

func TestSelectedFloatingState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		mockSelectedFloatingState func() (btt.Floating, error)

		expectedCode    int
		expectedHeaders map[string]string
		expectedBody    string
	}{
		{
			name: "getting selected floating state failed",
			mockSelectedFloatingState: func() (btt.Floating, error) {
				return "", errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			mockSelectedFloatingState: func() (btt.Floating, error) {
				return "foo", nil
			},
			expectedCode:    http.StatusOK,
			expectedHeaders: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			expectedBody:    "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := setupServer(t, func(mc *minimock.Controller, s *services.ServicesMock) {
				b := btt.NewBttMock(mc)
				b.SelectedFloatingMock.Return(tt.mockSelectedFloatingState())
				s.BttMock.Return(b)
			}, nil)

			w := performRequest(t, server.SelectedFloatingState, "")

			assertResponse(t, w, tt.expectedCode, tt.expectedHeaders, tt.expectedBody)
		})
	}
}

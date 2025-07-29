package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"live2text/internal/services"
	"live2text/internal/services/btt"
)

func TestSelectedClipboard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		mockSelectedClipboard func() (btt.Clipboard, error)

		expectedCode    int
		expectedHeaders map[string]string
		expectedBody    string
	}{
		{
			name: "getting selected clipboard error",
			mockSelectedClipboard: func() (btt.Clipboard, error) {
				return "", errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			mockSelectedClipboard: func() (btt.Clipboard, error) {
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
				b.SelectedClipboardMock.Return(tt.mockSelectedClipboard())
				s.BttMock.Return(b)
			}, nil)

			w := performRequest(t, server.SelectedClipboard, "")

			assertResponse(t, w, tt.expectedCode, tt.expectedHeaders, tt.expectedBody)
		})
	}
}

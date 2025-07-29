package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"live2text/internal/services"
	"live2text/internal/services/btt"
)

func TestSelectClipboard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		body                string
		mockSelectClipboard func() error

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
			name: "validation error",
			body: `{"clipboard":"invalid"}`,

			expectedCode: http.StatusBadRequest,
			expectedBody: `{"clipboard":"clipboard is not valid"}`,
		},
		{
			name: "select clipboard error",
			body: `{"clipboard":"Shown"}`,
			mockSelectClipboard: func() error {
				return errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			body: `{"clipboard":"Shown"}`,
			mockSelectClipboard: func() error {
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
				if tt.mockSelectClipboard != nil {
					b := btt.NewBttMock(mc)
					b.SelectClipboardMock.Return(tt.mockSelectClipboard())
					s.BttMock.Return(b)
				}
			}, nil)

			w := performRequest(t, server.SelectClipboard, tt.body)

			assertResponse(t, w, tt.expectedCode, nil, tt.expectedBody)
		})
	}
}

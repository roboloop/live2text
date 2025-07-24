package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"live2text/internal/services"
	"live2text/internal/services/btt"
)

func TestSelectViewMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		body               string
		mockSelectViewMode func() error

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
			name: "validation failed",
			body: `{"view_mode":"invalid"}`,

			expectedCode: http.StatusBadRequest,
			expectedBody: `{"view_mode":"view_mode is not valid"}`,
		},
		{
			name: "select view mode failed",
			body: `{"view_mode":"Embed"}`,
			mockSelectViewMode: func() error {
				return errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			body: `{"view_mode":"Embed"}`,
			mockSelectViewMode: func() error {
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
				if tt.mockSelectViewMode != nil {
					b := btt.NewBttMock(mc)
					b.SelectViewModeMock.Return(tt.mockSelectViewMode())
					s.BttMock.Return(b)
				}
			}, nil)

			w := performRequest(t, server.SelectViewMode, tt.body)

			assertResponse(t, w, tt.expectedCode, nil, tt.expectedBody)
		})
	}
}

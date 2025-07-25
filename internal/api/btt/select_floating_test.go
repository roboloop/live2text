package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"live2text/internal/services"
	"live2text/internal/services/btt"
)

func TestSelectFloating(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		body               string
		mockSelectFloating func() error

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
			body: `{"floating":"invalid"}`,

			expectedCode: http.StatusBadRequest,
			expectedBody: `{"floating":"floating is not valid"}`,
		},
		{
			name: "select floating failed",
			body: `{"floating":"Shown"}`,
			mockSelectFloating: func() error {
				return errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			body: `{"floating":"Shown"}`,
			mockSelectFloating: func() error {
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
				if tt.mockSelectFloating != nil {
					b := btt.NewBttMock(mc)
					b.SelectFloatingMock.Return(tt.mockSelectFloating())
					s.BttMock.Return(b)
				}
			}, nil)

			w := performRequest(t, server.SelectFloating, tt.body)

			assertResponse(t, w, tt.expectedCode, nil, tt.expectedBody)
		})
	}
}

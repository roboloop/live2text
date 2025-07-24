package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"live2text/internal/services"
	"live2text/internal/services/btt"
)

func TestSelectLanguage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		body               string
		mockSelectLanguage func() error

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
			body: `{"language":"invalid-code"}`,

			expectedCode: http.StatusBadRequest,
			expectedBody: `{"language":"language is not valid"}`,
		},
		{
			name: "select language failed",
			body: `{"language":"en-US"}`,
			mockSelectLanguage: func() error {
				return errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			body: `{"language":"en-US"}`,
			mockSelectLanguage: func() error {
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
				if tt.mockSelectLanguage != nil {
					b := btt.NewBttMock(mc)
					b.SelectLanguageMock.Return(tt.mockSelectLanguage())
					s.BttMock.Return(b)
				}
			}, nil)

			w := performRequest(t, server.SelectLanguage, tt.body)

			assertResponse(t, w, tt.expectedCode, nil, tt.expectedBody)
		})
	}
}

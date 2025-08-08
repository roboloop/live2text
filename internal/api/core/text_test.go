package core_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/api/core"
	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/services/recognition"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		body     string
		writer   func() http.ResponseWriter
		mockText func() (string, error)

		expectedCode int
		expectedBody string
	}{
		{
			name: "bad request",
			body: "invalid json",
			writer: func() http.ResponseWriter {
				return httptest.NewRecorder()
			},

			expectedCode: http.StatusBadRequest,
			expectedBody: "cannot decode request",
		},
		{
			name: "no task found",
			body: `{"id": "foo"}`,
			writer: func() http.ResponseWriter {
				return httptest.NewRecorder()
			},
			mockText: func() (string, error) {
				return "", recognition.ErrNoTaskFound
			},

			expectedCode: http.StatusBadRequest,
			expectedBody: "no task found",
		},
		{
			name: "task search error",
			body: `{"id": "foo"}`,
			writer: func() http.ResponseWriter {
				return httptest.NewRecorder()
			},
			mockText: func() (string, error) {
				return "", errors.New("dummy error")
			},

			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "error during writing",
			body: `{"id": "foo"}`,
			writer: func() http.ResponseWriter {
				return &errorResponseWriter{}
			},
			mockText: func() (string, error) {
				return "sample text", nil
			},

			expectedBody: "",
			expectedCode: http.StatusOK,
		},
		{
			name: "ok",
			body: `{"id": "foo"}`,
			writer: func() http.ResponseWriter {
				return httptest.NewRecorder()
			},
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

			mc := minimock.NewController(t)
			s := services.NewServicesMock(mc)

			if tt.mockText != nil {
				r := recognition.NewRecognitionMock(mc)
				r.TextMock.Return(tt.mockText())
				s.RecognitionMock.Return(r)
			}

			server := core.NewServer(logger.NilLogger, s)
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.Text(w, req)

			resp := w.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

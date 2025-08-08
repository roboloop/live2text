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
	"github.com/roboloop/live2text/internal/services/audio"
	"github.com/roboloop/live2text/internal/services/recognition"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		mockList  func() ([]string, error)
		mockStart func() (string, string, error)

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
			name: "audio list error",
			body: "{}",
			mockList: func() ([]string, error) {
				return nil, errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "validation error",
			body: `{"device":"foo","language":"bar"}`,
			mockList: func() ([]string, error) {
				return []string{"baz"}, nil
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"device":"device not found","language":"language is not valid"}`,
		},
		{
			name: "device is busy",
			body: `{"device":"foo","language":"en-US"}`,
			mockList: func() ([]string, error) {
				return []string{"foo"}, nil
			},
			mockStart: func() (string, string, error) {
				return "", "", recognition.ErrDeviceIsBusy
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "device is busy",
		},
		{
			name: "device search error",
			body: `{"device":"foo","language":"en-US"}`,
			mockList: func() ([]string, error) {
				return []string{"foo"}, nil
			},
			mockStart: func() (string, string, error) {
				return "", "", errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			body: `{"device":"foo","language":"en-US"}`,
			mockList: func() ([]string, error) {
				return []string{"foo"}, nil
			},
			mockStart: func() (string, string, error) {
				return "bar", "/path/to/socket", nil
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"bar","socketPath":"/path/to/socket"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			s := services.NewServicesMock(mc)
			if tt.mockList != nil {
				a := audio.NewAudioMock(mc)
				a.ListOfNamesMock.Return(tt.mockList())
				s.AudioMock.Return(a)
			}
			if tt.mockStart != nil {
				r := recognition.NewRecognitionMock(mc)
				r.StartMock.Return(tt.mockStart())
				s.RecognitionMock.Return(r)
			}

			server := core.NewServer(logger.NilLogger, s)
			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.Start(w, req)

			require.Equal(t, tt.expectedCode, w.Code)
			require.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

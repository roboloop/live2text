package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/services/btt"
)

func TestSelectedViewMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		mockSelectedViewMode func() (btt.ViewMode, error)

		expectedCode    int
		expectedHeaders map[string]string
		expectedBody    string
	}{
		{
			name: "getting selected view mode error",
			mockSelectedViewMode: func() (btt.ViewMode, error) {
				return "", errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			mockSelectedViewMode: func() (btt.ViewMode, error) {
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
				b.SelectedViewModeMock.Return(tt.mockSelectedViewMode())
				s.BttMock.Return(b)
			}, nil)

			w := performRequest(t, server.SelectedViewMode, "")

			assertResponse(t, w, tt.expectedCode, tt.expectedHeaders, tt.expectedBody)
		})
	}
}

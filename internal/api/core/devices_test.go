package core_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/api/core"
	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/services/audio"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestDevices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mockList func(mc *minimock.Controller) audio.Audio

		expectedCode   int
		expectedHeader string
		expectedBody   string
	}{
		{
			"returns list of audio devices",
			func(mc *minimock.Controller) audio.Audio {
				a := audio.NewAudioMock(mc)
				a.ListOfNamesMock.Return([]string{"foo", "bar"}, nil)

				return a
			},
			http.StatusOK,
			"application/json",
			`{"devices":["foo","bar"]}`,
		},
		{
			"handles audio service error",
			func(mc *minimock.Controller) audio.Audio {
				a := audio.NewAudioMock(mc)
				a.ListOfNamesMock.Return([]string{}, errors.New("dummy error"))

				return a
			},
			http.StatusInternalServerError,
			"text/plain; charset=utf-8",
			`dummy error`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			a := tt.mockList(mc)
			s := services.NewServicesMock(mc)
			s.AudioMock.Return(a)

			server := core.NewServer(logger.NilLogger, s)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			server.Devices(w, req)

			resp := w.Result()
			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Equal(t, tt.expectedHeader, resp.Header.Get("Content-Type"))
			require.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

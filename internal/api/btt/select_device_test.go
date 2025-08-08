package btt_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"

	"github.com/roboloop/live2text/internal/services"
	"github.com/roboloop/live2text/internal/services/audio"
	"github.com/roboloop/live2text/internal/services/btt"
)

func TestSelectDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		body             string
		mockList         func() ([]string, error)
		mockSelectDevice func() error

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
			body: `{"device":"foo"}`,
			mockList: func() ([]string, error) {
				return []string{"baz"}, nil
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"device":"device not found"}`,
		},
		{
			name: "select device error",
			body: `{"device":"foo"}`,
			mockList: func() ([]string, error) {
				return []string{"foo"}, nil
			},
			mockSelectDevice: func() error {
				return errors.New("dummy error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "dummy error",
		},
		{
			name: "ok",
			body: `{"device":"foo"}`,
			mockList: func() ([]string, error) {
				return []string{"foo"}, nil
			},
			mockSelectDevice: func() error {
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
				if tt.mockList != nil {
					a := audio.NewAudioMock(mc)
					a.ListOfNamesMock.Return(tt.mockList())
					s.AudioMock.Return(a)
				}
				if tt.mockSelectDevice != nil {
					b := btt.NewBttMock(mc)
					b.SelectDeviceMock.Return(tt.mockSelectDevice())
					s.BttMock.Return(b)
				}
			}, nil)

			w := performRequest(t, server.SelectDevice, tt.body)

			assertResponse(t, w, tt.expectedCode, nil, tt.expectedBody)
		})
	}
}

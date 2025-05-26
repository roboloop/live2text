package btt_test

import (
	"errors"
	"live2text/internal/services/btt"
	btthttp "live2text/internal/services/btt/http"
	"live2text/internal/utils"
	"strings"
	"testing"
)

func TestSelectedDevice(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name string

		mockHTTPClient *btthttp.MockClient

		expected    string
		expectedErr string
	}{
		{
			name:           "Cannot get selected device variable",
			mockHTTPClient: &btthttp.MockClient{SendError: []error{errors.New("bad happened")}},
			expectedErr:    "cannot get selected device variable: cannot get string variable: bad happened",
		},
		{
			name:           "Get selected device",
			mockHTTPClient: &btthttp.MockClient{SendResponse: [][]byte{[]byte("foo")}},
			expected:       "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAudio, mockRecognition, _, mockExecClient, cfg := newMocks()
			b := btt.NewBtt(utils.NilLogger, mockAudio, mockRecognition, tt.mockHTTPClient, mockExecClient, cfg)
			device, err := b.SelectedDevice(ctx)
			if tt.expectedErr != "" && err != nil {
				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if device != tt.expected {
				t.Errorf("Expected device name: %v, got: %v", tt.expected, device)
				return
			}
		})
	}
}

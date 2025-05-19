package btt_test

import (
	"context"
	"errors"
	"live2text/internal/services/btt"
	btthttp "live2text/internal/services/btt/http"
	"live2text/internal/utils"
	"strings"
	"testing"
)

func TestSelectedLanguage(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		mockHttpClient *btthttp.MockClient

		expected    string
		expectedErr string
	}{
		{
			name:           "Cannot get selected language variable",
			mockHttpClient: &btthttp.MockClient{SendError: []error{errors.New("bad happened")}},
			expectedErr:    "cannot get selected language variable: cannot get string variable: bad happened",
		},
		{
			name:           "Get selected language",
			mockHttpClient: &btthttp.MockClient{SendResponse: [][]byte{[]byte("foo")}},
			expected:       "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAudio, mockRecognition, _, mockExecClient, cfg := newMocks()
			b := btt.NewBtt(utils.NilLogger, mockAudio, mockRecognition, tt.mockHttpClient, mockExecClient, cfg)
			language, err := b.SelectedLanguage(ctx)
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
			if language != tt.expected {
				t.Errorf("Expected language name: %v, got: %v", tt.expected, language)
				return
			}
		})
	}
}

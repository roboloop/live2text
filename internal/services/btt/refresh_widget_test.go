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

func TestRefreshWidget(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		mockHttpClient *btthttp.MockClient

		expectedErr string
	}{
		{
			name:           "Cannot refresh widget",
			mockHttpClient: &btthttp.MockClient{SendError: []error{errors.New("bad happened")}},
			expectedErr:    "cannot refresh widget: bad happened",
		},
		{
			name:           "Widget refreshed",
			mockHttpClient: &btthttp.MockClient{SendResponse: [][]byte{[]byte("foo")}},
			expectedErr:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAudio, mockRecognition, _, mockExecClient, cfg := newMocks()
			b := btt.NewBtt(utils.NilLogger, mockAudio, mockRecognition, tt.mockHttpClient, mockExecClient, cfg)
			err := b.RefreshWidget(ctx, "DUMMY-UUID")
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
		})
	}
}

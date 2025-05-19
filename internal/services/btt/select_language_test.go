package btt_test

import (
	"context"
	"errors"
	"live2text/internal/services/btt"
	bttexec "live2text/internal/services/btt/exec"
	btthttp "live2text/internal/services/btt/http"
	"live2text/internal/utils"
	"strings"
	"testing"
)

func TestSelectLanguage(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		mockHttpClient *btthttp.MockClient
		mockExecClient *bttexec.MockClient

		language string

		expectedErr string
	}{
		{
			name:           "Cannot save selected language in btt",
			mockHttpClient: &btthttp.MockClient{SendError: []error{errors.New("bad happened")}},
			mockExecClient: &bttexec.MockClient{},
			language:       "en-US",
			expectedErr:    "cannot set selected language: cannot set string variable: bad happened",
		},
		{
			name:           "Cannot get selected language trigger uuid",
			mockHttpClient: &btthttp.MockClient{SendResponse: [][]byte{{}}},
			mockExecClient: &bttexec.MockClient{ExecError: []error{errors.New("bad happened")}},
			language:       "en-US",
			expectedErr:    "cannot get selected language trigger: cannot exec get_triggers command: bad happened",
		},
		{
			name:           "Cannot get selected refresh widget",
			mockHttpClient: &btthttp.MockClient{SendResponse: [][]byte{[]byte{}}, SendError: []error{nil, errors.New("bad happened")}},
			mockExecClient: &bttexec.MockClient{ExecResponse: [][]byte{[]byte(`[{"BTTTriggerTypeDescription": "Selected Language", "BTTUUID": "DUMMY-UUID"}]`)}},
			language:       "en-US",
			expectedErr:    "cannot refresh selected language widget: cannot refresh widget: bad happened",
		},
		{
			name:           "Language was selected",
			mockHttpClient: &btthttp.MockClient{SendResponse: [][]byte{{}, {}}},
			mockExecClient: &bttexec.MockClient{ExecResponse: [][]byte{[]byte(`[{"BTTTriggerTypeDescription": "Selected Language", "BTTUUID": "DUMMY-UUID"}]`)}},
			language:       "en-US",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAudio, mockRecognition, _, _, cfg := newMocks()
			b := btt.NewBtt(utils.NilLogger, mockAudio, mockRecognition, tt.mockHttpClient, tt.mockExecClient, cfg)

			err := b.SelectLanguage(ctx, tt.language)
			if tt.expectedErr != "" && err != nil {
				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}

}

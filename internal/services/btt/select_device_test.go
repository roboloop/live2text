package btt_test

import (
	"errors"
	"live2text/internal/services/btt"
	bttexec "live2text/internal/services/btt/exec"
	btthttp "live2text/internal/services/btt/http"
	"live2text/internal/utils"
	"strings"
	"testing"
)

func TestSelectDevice(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name string

		mockHTTPClient *btthttp.MockClient
		mockExecClient *bttexec.MockClient

		deviceName string

		expectedErr string
	}{
		{
			name:           "Cannot save selected device in btt",
			mockHTTPClient: &btthttp.MockClient{SendError: []error{errors.New("bad happened")}},
			mockExecClient: &bttexec.MockClient{},
			deviceName:     "foo",
			expectedErr:    "cannot set selected device: cannot set string variable: bad happened",
		},
		{
			name:           "Cannot get selected device trigger uuid",
			mockHTTPClient: &btthttp.MockClient{SendResponse: [][]byte{{}}},
			mockExecClient: &bttexec.MockClient{ExecError: []error{errors.New("bad happened")}},
			deviceName:     "foo",
			expectedErr:    "cannot get selected device trigger: cannot exec get_triggers command: bad happened",
		},
		{
			name: "Cannot get selected refresh widget",
			mockHTTPClient: &btthttp.MockClient{
				SendResponse: [][]byte{[]byte{}},
				SendError:    []error{nil, errors.New("bad happened")},
			},
			mockExecClient: &bttexec.MockClient{
				ExecResponse: [][]byte{
					[]byte(`[{"BTTTriggerTypeDescription": "Selected Device", "BTTUUID": "DUMMY-UUID"}]`),
				},
			},
			deviceName:  "foo",
			expectedErr: "cannot refresh selected device widget: cannot refresh widget: bad happened",
		},
		{
			name:           "Device was selected",
			mockHTTPClient: &btthttp.MockClient{SendResponse: [][]byte{{}, {}}},
			mockExecClient: &bttexec.MockClient{
				ExecResponse: [][]byte{
					[]byte(`[{"BTTTriggerTypeDescription": "Selected Device", "BTTUUID": "DUMMY-UUID"}]`),
				},
			},
			deviceName: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAudio, mockRecognition, _, _, cfg := newMocks()
			b := btt.NewBtt(utils.NilLogger, mockAudio, mockRecognition, tt.mockHTTPClient, tt.mockExecClient, cfg)

			err := b.SelectDevice(ctx, tt.deviceName)
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

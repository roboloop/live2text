package btt_test

import (
	"context"
	"errors"
	"github.com/gordonklaus/portaudio"
	"live2text/internal/services/audio"
	"live2text/internal/services/btt"
	bttexec "live2text/internal/services/btt/exec"
	btthttp "live2text/internal/services/btt/http"
	"live2text/internal/utils"
	"strings"
	"testing"
)

func TestLoadDevices(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string

		mockAudio      *audio.MockAudio
		mockHttpClient *btthttp.MockClient
		mockExecClient *bttexec.MockClient

		expectedErr string
	}{
		{
			name:        "Cannot get list of devices",
			mockAudio:   &audio.MockAudio{ListDeviceInfo: nil, ListError: errors.New("bad happened")},
			expectedErr: "cannot get list of devices: bad happened",
		},
		{
			name:           "Cannot get trigger uuid",
			mockAudio:      &audio.MockAudio{ListDeviceInfo: []*portaudio.DeviceInfo{{Name: "foo"}}},
			mockExecClient: &bttexec.MockClient{ExecError: []error{errors.New("bad happened")}},
			expectedErr:    "cannot find device group: cannot exec get_triggers command: bad happened",
		},
		{
			name:      "Cannot get triggers",
			mockAudio: &audio.MockAudio{ListDeviceInfo: []*portaudio.DeviceInfo{{Name: "foo"}}},
			mockExecClient: &bttexec.MockClient{
				ExecResponse: [][]byte{[]byte(`[{"BTTTriggerTypeDescription": "Device", "BTTUUID": "FOO-UUID"}]`)},
				ExecError:    []error{nil, errors.New("bad happened")},
			},
			expectedErr: "cannot get triggers: cannot exec get_triggers command: bad happened",
		},
		{
			name:      "Cannot get triggers",
			mockAudio: &audio.MockAudio{ListDeviceInfo: []*portaudio.DeviceInfo{{Name: "foo"}}},
			mockExecClient: &bttexec.MockClient{
				ExecResponse: [][]byte{[]byte(`[{"BTTTriggerTypeDescription": "Device", "BTTUUID": "FOO-UUID"}]`)},
				ExecError:    []error{nil, errors.New("bad happened")},
			},
			expectedErr: "cannot get triggers: cannot exec get_triggers command: bad happened",
		},
		// TODO: more tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mockRecognition, _, _, cfg := newMocks()
			b := btt.NewBtt(utils.NilLogger, tt.mockAudio, mockRecognition, tt.mockHttpClient, tt.mockExecClient, cfg)
			err := b.LoadDevices(ctx)
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

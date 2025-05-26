package audio_test

import (
	"errors"
	"live2text/internal/services/audio"
	audiowrapper "live2text/internal/services/audio_wrapper"
	"live2text/internal/utils"
	"reflect"
	"testing"

	"github.com/gordonklaus/portaudio"
)

func TestList(t *testing.T) {
	validDevice := &portaudio.DeviceInfo{Name: "foo", MaxInputChannels: 1}
	invalidDevice := &portaudio.DeviceInfo{Name: "bar", MaxInputChannels: 0}

	tests := []struct {
		name string

		mockAudioWrapper *audiowrapper.MockAudio

		expected    []*portaudio.DeviceInfo
		expectedErr string
	}{
		{
			name:             "DefaultHostAPI fails",
			mockAudioWrapper: &audiowrapper.MockAudio{DefaultHostAPIError: errors.New("host api failed")},
			expectedErr:      "cannot list host apis: host api failed",
		},
		{
			name: "List of devices",
			mockAudioWrapper: &audiowrapper.MockAudio{
				DefaultHostAPIHostAPIInfo: &portaudio.HostApiInfo{
					Devices: []*portaudio.DeviceInfo{validDevice, invalidDevice},
				},
			},
			expected: []*portaudio.DeviceInfo{validDevice},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := audio.NewAudio(utils.NilLogger, nil, tt.mockAudioWrapper)
			devices, err := a.List()
			if tt.expectedErr != "" {
				if tt.expectedErr != err.Error() {
					t.Errorf("List() got %v, expected %v", err, tt.expectedErr)
				}
				return
			}
			if !reflect.DeepEqual(devices, tt.expected) {
				t.Errorf("List() got %v, expected %v", devices, tt.expected)
			}
		})
	}
}

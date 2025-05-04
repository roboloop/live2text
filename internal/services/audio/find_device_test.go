package audio_test

import (
	"errors"
	"github.com/gordonklaus/portaudio"
	"live2text/internal/services/audio"
	"live2text/internal/services/audio_wrapper"
	"live2text/internal/utils"
	"reflect"
	"testing"
)

func TestFindInputDevice(t *testing.T) {
	tests := []struct {
		name string

		mockAudioWrapper *audio_wrapper.MockAudio

		deviceName string

		expected    *portaudio.DeviceInfo
		expectedErr string
	}{
		{
			name:             "DefaultHostApi fails",
			mockAudioWrapper: &audio_wrapper.MockAudio{DefaultHostApiError: errors.New("host api failed")},
			deviceName:       "mic1",
			expectedErr:      "cannot list host apis: host api failed",
		},
		{
			name:             "Device not found",
			mockAudioWrapper: &audio_wrapper.MockAudio{DefaultHostApiHostApiInfo: &portaudio.HostApiInfo{Devices: []*portaudio.DeviceInfo{{Name: "bar"}}}},
			deviceName:       "foo",
			expectedErr:      "device was not found",
		},
		{
			name:             "No input channel",
			mockAudioWrapper: &audio_wrapper.MockAudio{DefaultHostApiHostApiInfo: &portaudio.HostApiInfo{Devices: []*portaudio.DeviceInfo{{Name: "foo"}}}},
			deviceName:       "foo",
			expectedErr:      "device hasn't input channels",
		},
		{
			name:             "Device found",
			mockAudioWrapper: &audio_wrapper.MockAudio{DefaultHostApiHostApiInfo: &portaudio.HostApiInfo{Devices: []*portaudio.DeviceInfo{{Name: "foo", MaxInputChannels: 1}}}},
			deviceName:       "foo",
			expected:         &portaudio.DeviceInfo{Name: "foo", MaxInputChannels: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := audio.NewAudio(utils.NilLogger, nil, tt.mockAudioWrapper)
			device, err := a.FindInputDevice(tt.deviceName)
			if tt.expectedErr != "" {
				if tt.expectedErr != err.Error() {
					t.Errorf("FindInputDevice() got %v, expected %v", err, tt.expectedErr)
				}
				return
			}
			if !reflect.DeepEqual(device, tt.expected) {
				t.Errorf("FindInputDevice() got %#v, expected %#v", device, tt.expected)
				return
			}
		})
	}
}

package audio_test

import (
	"errors"
	"live2text/internal/services/audio"
	audiowrapper "live2text/internal/services/audio_wrapper"
	"live2text/internal/utils"
	"testing"

	"github.com/gordonklaus/portaudio"
)

func TestListenDevice(t *testing.T) {
	ctx := t.Context()

	var (
		deviceName   = "foo"
		device       = &portaudio.DeviceInfo{Name: deviceName, MaxInputChannels: 1}
		hostAPIInfo  = &portaudio.HostApiInfo{Devices: []*portaudio.DeviceInfo{device}}
		listenerInfo = &audio.ListenerInfo{Device: device, Channels: 1, ChunkSizeMs: 100}
	)

	for _, tt := range []struct {
		name string

		mockAudioWrapper *audiowrapper.MockAudio

		deviceName string

		expected    *audio.ListenerInfo
		expectedErr string
	}{
		{
			name:             "No device",
			mockAudioWrapper: &audiowrapper.MockAudio{DefaultHostAPIError: errors.New("internal")},
			deviceName:       deviceName,
			expected:         listenerInfo,
			expectedErr:      "cannot find input device: cannot list host apis: internal",
		},
		{
			name:             "Happy path",
			mockAudioWrapper: &audiowrapper.MockAudio{DefaultHostAPIHostAPIInfo: hostAPIInfo, OpenStreamStream: &audiowrapper.MockStream{}},
			deviceName:       deviceName,
			expected:         listenerInfo,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			a := audio.NewAudio(utils.NilLogger, nil, tt.mockAudioWrapper)
			listener, err := a.ListenDevice(ctx, tt.deviceName)

			if tt.expectedErr != "" {
				if err == nil {
					t.Fatalf("ListenDevice() expected error %v, got nil", tt.expectedErr)
				}

				if tt.expectedErr != err.Error() {
					t.Errorf("Expected error: %v, got %v", tt.expectedErr, err.Error())
				}
				return
			}

			if listener.Channels != tt.expected.Channels ||
				listener.SampleRate != tt.expected.SampleRate ||
				listener.ChunkSizeMs != tt.expected.ChunkSizeMs {
				t.Errorf("ListenDevice() got %#v, expected %#v", listener, tt.expected)
				return
			}
		})
	}
}

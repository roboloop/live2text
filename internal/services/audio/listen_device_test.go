package audio_test

import (
	"context"
	"errors"
	"github.com/gordonklaus/portaudio"
	"live2text/internal/services/audio"
	"live2text/internal/services/audio_wrapper"
	"live2text/internal/utils"
	"testing"
)

func TestListenDevice(t *testing.T) {
	ctx := context.Background()

	var (
		deviceName   = "foo"
		device       = &portaudio.DeviceInfo{Name: deviceName, MaxInputChannels: 1}
		hostApiInfo  = &portaudio.HostApiInfo{Devices: []*portaudio.DeviceInfo{device}}
		listenerInfo = &audio.ListenerInfo{Device: device, Channels: 1, ChunkSizeMs: 100}
	)

	for _, tt := range []struct {
		name string

		mockAudioWrapper *audio_wrapper.MockAudio

		deviceName string

		expected    *audio.ListenerInfo
		expectedErr string
	}{
		{
			name:             "No stream",
			mockAudioWrapper: &audio_wrapper.MockAudio{DefaultHostApiHostApiInfo: hostApiInfo, OpenStreamError: errors.New("internal")},
			deviceName:       deviceName,
			expected:         listenerInfo,
			expectedErr:      "could not open the stream: internal",
		},
		{
			name:             "Stream not started",
			mockAudioWrapper: &audio_wrapper.MockAudio{DefaultHostApiHostApiInfo: hostApiInfo, OpenStreamStream: &audio_wrapper.MockStream{StartError: errors.New("internal")}},
			deviceName:       deviceName,
			expected:         listenerInfo,
			expectedErr:      "could not start the stream: internal",
		},
		{
			name:             "Read failed",
			mockAudioWrapper: &audio_wrapper.MockAudio{DefaultHostApiHostApiInfo: hostApiInfo, OpenStreamStream: &audio_wrapper.MockStream{ReadError: errors.New("internal")}},
			deviceName:       deviceName,
			expected:         listenerInfo,
			expectedErr:      "could not read stream: internal",
		},
		{
			name:             "Happy path",
			mockAudioWrapper: &audio_wrapper.MockAudio{DefaultHostApiHostApiInfo: hostApiInfo, OpenStreamStream: &audio_wrapper.MockStream{}},
			deviceName:       deviceName,
			expected:         listenerInfo,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			a := audio.NewAudio(utils.NilLogger, nil, tt.mockAudioWrapper)
			listener, err := a.ListenDevice(ctx, tt.deviceName)

			if tt.expectedErr != "" {
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

func TestListening(t *testing.T) {

}

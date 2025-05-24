package audio_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gordonklaus/portaudio"

	"live2text/internal/services/audio"
	"live2text/internal/services/audio_wrapper"
	"live2text/internal/utils"
)

// TestFindInputDevice tests the FindInputDevice function with various scenarios
// to ensure it correctly identifies audio input devices and handles error cases.
func TestFindInputDevice(t *testing.T) {
	// Define test cases
	tests := []struct {
		name             string        // Test case name
		mockAudioWrapper *audio_wrapper.MockAudio // Mock audio wrapper configuration
		deviceName       string        // Device name to search for
		expected         *portaudio.DeviceInfo // Expected device info result
		expectedErr      string        // Expected error message (if any)
	}{
		{
			name:             "DefaultHostApi fails",
			mockAudioWrapper: &audio_wrapper.MockAudio{DefaultHostApiError: errors.New("host api failed")},
			deviceName:       "mic1",
			expectedErr:      "cannot list host apis: host api failed",
		},
		{
			name: "Device not found",
			mockAudioWrapper: &audio_wrapper.MockAudio{
				DefaultHostApiHostApiInfo: &portaudio.HostApiInfo{
					Devices: []*portaudio.DeviceInfo{{Name: "bar"}},
				},
			},
			deviceName:  "foo",
			expectedErr: "device not found", // Match the actual error message in implementation
		},
		{
			name: "No input channels",
			mockAudioWrapper: &audio_wrapper.MockAudio{
				DefaultHostApiHostApiInfo: &portaudio.HostApiInfo{
					Devices: []*portaudio.DeviceInfo{{Name: "foo"}},
				},
			},
			deviceName:  "foo",
			expectedErr: "device hasn't input channels",
		},
		{
			name: "Device found successfully",
			mockAudioWrapper: &audio_wrapper.MockAudio{
				DefaultHostApiHostApiInfo: &portaudio.HostApiInfo{
					Devices: []*portaudio.DeviceInfo{{Name: "foo", MaxInputChannels: 1}},
				},
			},
			deviceName: "foo",
			expected:   &portaudio.DeviceInfo{Name: "foo", MaxInputChannels: 1},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize audio service with mock
			audioService := audio.NewAudio(utils.NilLogger, nil, tc.mockAudioWrapper)

			// Call the function under test
			device, err := audioService.FindInputDevice(tc.deviceName)

			// Check error cases
			if tc.expectedErr != "" {
				if err == nil {
					t.Fatalf("FindInputDevice() expected error %q, got nil", tc.expectedErr)
				}

				if err.Error() != tc.expectedErr {
					t.Errorf("FindInputDevice() error = %q, want %q", err.Error(), tc.expectedErr)
				}
			} else {
				// Check success case
				if err != nil {
					t.Fatalf("FindInputDevice() unexpected error: %v", err)
				}

				if !reflect.DeepEqual(device, tc.expected) {
					t.Errorf("FindInputDevice() device = %#v, want %#v", device, tc.expected)
				}
			}
		})
	}
}

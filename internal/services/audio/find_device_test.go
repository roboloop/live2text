package audio_test

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	audiowrapper "live2text/internal/services/audio_wrapper"
	"live2text/internal/services/metrics"
)

func TestFindInputDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		mockDevices func() ([]*audiowrapper.DeviceInfo, error)

		deviceName     string
		expectedDevice *audiowrapper.DeviceInfo
		expectedErr    string
	}{
		{
			name: "cannot get a list of devices",
			mockDevices: func() ([]*audiowrapper.DeviceInfo, error) {
				return nil, errors.New("dummy error")
			},
			expectedErr: "cannot get a list of devices",
		},
		{
			name: "device not found",
			mockDevices: func() ([]*audiowrapper.DeviceInfo, error) {
				return nil, nil
			},
			deviceName:  "mic1",
			expectedErr: "device not found",
		},
		{
			name: "device hasn't input channels",
			mockDevices: func() ([]*audiowrapper.DeviceInfo, error) {
				return []*audiowrapper.DeviceInfo{{Name: "mic1", MaxInputChannels: 0}}, nil
			},
			deviceName:  "mic1",
			expectedErr: "device hasn't input channels",
		},
		{
			name: "device found successfully",
			mockDevices: func() ([]*audiowrapper.DeviceInfo, error) {
				return []*audiowrapper.DeviceInfo{{Name: "mic1", MaxInputChannels: 1}}, nil
			},
			deviceName:     "mic1",
			expectedDevice: &audiowrapper.DeviceInfo{Name: "mic1", MaxInputChannels: 1},
			expectedErr:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := setupAudio(t, func(mc *minimock.Controller, m *metrics.MetricsMock, ea *audiowrapper.AudioMock) {
				ea.DevicesMock.Return(tt.mockDevices())
			}, nil)

			device, err := a.FindInputDevice(tt.deviceName)

			if tt.expectedErr != "" {
				require.Nil(t, device)
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedDevice, device)
		})
	}
}

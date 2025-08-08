package audio_test

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	audiowrapper "github.com/roboloop/live2text/internal/services/audio_wrapper"
	"github.com/roboloop/live2text/internal/services/metrics"
)

func TestListOfNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		mockDevices func() ([]*audiowrapper.DeviceInfo, error)

		expectedNames []string
		expectedErr   string
	}{
		{
			name: "cannot get a list of devices",
			mockDevices: func() ([]*audiowrapper.DeviceInfo, error) {
				return nil, errors.New("dummy error")
			},
			expectedNames: nil,
			expectedErr:   "cannot get a list of devices",
		},
		{
			name: "getting a list of devices",
			mockDevices: func() ([]*audiowrapper.DeviceInfo, error) {
				return []*audiowrapper.DeviceInfo{
					{Name: "mic1", MaxInputChannels: 0, DefaultSampleRate: 48000},
					{Name: "mic2", MaxInputChannels: 1, DefaultSampleRate: 24000},
					{Name: "mic3", MaxInputChannels: 2, DefaultSampleRate: 12000},
				}, nil
			},
			expectedNames: []string{"mic2", "mic3"},
			expectedErr:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := setupAudio(t, func(_ *minimock.Controller, m *metrics.MetricsMock, ea *audiowrapper.AudioMock) {
				ea.DevicesMock.Return(tt.mockDevices())
			}, nil)

			names, err := a.ListOfNames()

			if tt.expectedErr != "" {
				require.Nil(t, names)
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedNames, names)
		})
	}
}

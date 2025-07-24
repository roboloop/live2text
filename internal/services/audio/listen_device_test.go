package audio_test

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"live2text/internal/services/audio"
	audiowrapper "live2text/internal/services/audio_wrapper"
	"live2text/internal/services/metrics"
)

func TestListenDevice(t *testing.T) {
	t.Parallel()

	t.Run("cannot find the input device", func(t *testing.T) {
		t.Parallel()

		a := setupAudio(t, func(_ *minimock.Controller, m *metrics.MetricsMock, ea *audiowrapper.AudioMock) {
			ea.DevicesMock.Return(nil, errors.New("dummy error"))
		}, nil)

		listener, err := a.ListenDevice("mic1")

		require.Nil(t, listener)
		require.Error(t, err)
		require.ErrorContains(t, err, "cannot find the input device")
	})

	t.Run("getting the device listener", func(t *testing.T) {
		t.Parallel()

		a := setupAudio(t, func(_ *minimock.Controller, m *metrics.MetricsMock, ea *audiowrapper.AudioMock) {
			ea.DevicesMock.Return([]*audiowrapper.DeviceInfo{
				{Name: "mic1", MaxInputChannels: 1, DefaultSampleRate: 24000},
			}, nil)
		}, nil)

		listener, err := a.ListenDevice("mic1")

		require.NoError(t, err)
		require.Implements(t, (*audio.DeviceListener)(nil), listener)
	})
}

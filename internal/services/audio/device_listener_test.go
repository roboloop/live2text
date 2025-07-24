package audio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"live2text/internal/services/audio"
	audiowrapper "live2text/internal/services/audio_wrapper"
	"live2text/internal/services/metrics"
	"live2text/internal/utils/logger"
)

func TestDeviceListener(t *testing.T) {
	t.Parallel()

	t.Run("assert default device listener parameters", func(t *testing.T) {
		t.Parallel()

		mc := minimock.NewController(t)
		m := metrics.NewMetricsMock(mc)
		ea := audiowrapper.NewAudioMock(mc)
		d := &audiowrapper.DeviceInfo{DefaultSampleRate: 24000}

		listener := audio.NewDeviceListener(logger.NilLogger, m, ea, d)

		parameters := listener.GetParameters()
		require.Equal(t, 1, parameters.Channels)
		require.Equal(t, 100, parameters.ChunkSizeMs)
		require.Equal(t, 12000, parameters.SampleRate)

		ch := listener.GetChannel()
		require.Equal(t, 1024, cap(ch))
	})

	tests := []struct {
		name         string
		setupMocks   func(mc *minimock.Controller, m *metrics.MetricsMock, ea *audiowrapper.AudioMock)
		expectedLogs func(t *testing.T, h *logger.CaptureHandler)
		expectedErr  string
	}{
		{
			name: "cannot open the stream",
			setupMocks: func(_ *minimock.Controller, _ *metrics.MetricsMock, ea *audiowrapper.AudioMock) {
				ea.StreamDeviceMock.Return(nil, errors.New("dummy error"))
			},
			expectedErr: "cannot open the stream",
		},
		{
			name: "cannot start the stream",
			setupMocks: func(mc *minimock.Controller, _ *metrics.MetricsMock, ea *audiowrapper.AudioMock) {
				s := audiowrapper.NewStreamMock(mc)
				s.StartMock.Return(errors.New("dummy error"))
				ea.StreamDeviceMock.Return(s, nil)

				// defer
				s.CloseMock.Return(nil)
			},
			expectedErr: "cannot start the stream",
		},
		{
			name: "cannot read the stream",
			setupMocks: func(mc *minimock.Controller, _ *metrics.MetricsMock, ea *audiowrapper.AudioMock) {
				s := audiowrapper.NewStreamMock(mc)
				s.StartMock.Return(nil)
				s.ReadMock.Return(nil, errors.New("dummy error"))
				ea.StreamDeviceMock.Return(s, nil)

				// defer
				s.CloseMock.Return(nil)
				s.StopMock.Return(nil)
			},
			expectedErr: "cannot read the stream",
		},
		{
			name: "reads data from the stream successfully",
			setupMocks: func(mc *minimock.Controller, m *metrics.MetricsMock, ea *audiowrapper.AudioMock) {
				s := audiowrapper.NewStreamMock(mc)
				s.StartMock.Return(nil)
				s.ReadMock.Return([]int16{0, 1, 2}, nil)
				ea.StreamDeviceMock.Return(s, nil)

				// defer
				s.CloseMock.Return(nil)
				s.StopMock.Return(nil)

				m.AddBytesReadFromAudioMock.Expect(6)
			},
			expectedErr: "",
		},
		// {
		// 	name: "the channel is full, a segment dropped",
		// 	setupMocks: func(mc *minimock.Controller, m *metrics.MetricsMock, ea *audiowrapper.AudioMock) {
		//
		// 		s := audiowrapper.NewStreamMock(mc)
		// 		s.StartMock.Return(nil)
		// 		s.ReadMock.Times(1025).Return([]int16{0, 1, 2}, nil)
		// 		ea.StreamDeviceMock.Return(s, nil)
		//
		// 		// defer
		// 		s.CloseMock.Return(nil)
		// 		s.StopMock.Return(nil)
		//
		// 		m.AddBytesReadFromAudioMock.Return()
		// 	},
		// 	expectedLogs: func(t *testing.T, h *logger.CaptureHandler) {
		// 		require.Len(t, h.Logs, 1)
		// 		require.Equal(t, slog.LevelError, h.Logs[0].Level)
		// 		require.Equal(t, "The channel is full, the segment was dropped", h.Logs[0].Msg)
		// 	},
		// },
		{
			name: "handles context cancellation",
			setupMocks: func(mc *minimock.Controller, m *metrics.MetricsMock, ea *audiowrapper.AudioMock) {
				s := audiowrapper.NewStreamMock(mc)
				s.StartMock.Return(nil)
				s.ReadMock.Return([]int16{0, 1, 2}, nil)
				ea.StreamDeviceMock.Return(s, nil)

				// defer
				s.CloseMock.Return(nil)
				s.StopMock.Return(nil)

				m.AddBytesReadFromAudioMock.Return()
			},
			expectedLogs: nil,
			expectedErr:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			l, h := logger.NewCaptureLogger()
			listener := setupDeviceListener(t, tt.setupMocks, l, nil)

			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			err := listener.Listen(ctx)

			if tt.expectedLogs != nil {
				tt.expectedLogs(t, h)
			}

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

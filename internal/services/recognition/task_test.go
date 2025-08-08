package recognition_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/audio"
	audiowrapper "github.com/roboloop/live2text/internal/services/audio_wrapper"
	"github.com/roboloop/live2text/internal/services/recognition"
	"github.com/roboloop/live2text/internal/services/recognition/components"
	"github.com/roboloop/live2text/internal/services/recognition/text"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(mc *minimock.Controller, a *audio.AudioMock, dl *audio.DeviceListenerMock, b *components.BurnerComponentMock, r *components.RecognizerComponentMock, s *components.SocketComponentMock, o *components.OutputComponentMock)
		device     string
		language   string
		socketPath string

		expectErr    string
		expectErrMsg string
	}{
		{
			name: "cannot listen to a device",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, dl *audio.DeviceListenerMock, b *components.BurnerComponentMock, r *components.RecognizerComponentMock, s *components.SocketComponentMock, o *components.OutputComponentMock) {
				a.ListenDeviceMock.Return(nil, errors.New("something happened"))
			},
			expectErr: "cannot listen to a device",
		},
		{
			name: "something happened",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, dl *audio.DeviceListenerMock, b *components.BurnerComponentMock, r *components.RecognizerComponentMock, s *components.SocketComponentMock, o *components.OutputComponentMock) {
				a.ListenDeviceMock.Return(dl, nil)
				dl.ListenMock.Return(nil)
				dl.GetParametersMock.Return(&audiowrapper.StreamParameters{})
				dl.GetChannelMock.Return(nil)
				b.SaveAudioMock.Return(nil)
				r.RecognizeMock.Return(nil)
				o.ToConsoleMock.Return(nil)
				o.ToFileMock.Return(nil)
				o.PrintMock.Return(nil)
				s.ListenMock.Return(errors.New("something happened"))
			},
			expectErr: "something happened",
		},
		{
			name: "context cancelled",
			setupMocks: func(mc *minimock.Controller, a *audio.AudioMock, dl *audio.DeviceListenerMock, b *components.BurnerComponentMock, r *components.RecognizerComponentMock, s *components.SocketComponentMock, o *components.OutputComponentMock) {
				a.ListenDeviceMock.Expect("foo").Return(dl, nil)
				dl.ListenMock.Expect(minimock.AnyContext).Return(nil)
				dl.GetParametersMock.Return(&audiowrapper.StreamParameters{
					Channels:    1,
					ChunkSizeMs: 100,
					SampleRate:  24000,
				})
				dl.GetChannelMock.Return(make(chan []int16, 10))
				b.SaveAudioMock.ExpectParametersParam3(components.BurnerParameters{
					Channels:   1,
					SampleRate: 24000,
				}).Return(nil)
				r.RecognizeMock.ExpectParametersParam2(components.RecognizeParameters{
					Channels:   1,
					SampleRate: 24000,
					Language:   "bar",
				}).Return(nil)
				o.ToConsoleMock.Return(nil)
				o.ToFileMock.Return(nil)
				o.PrintMock.Return(nil)
				s.ListenMock.Set(func(ctx context.Context, socketPath string, formatter *text.Formatter) error {
					require.Equal(t, "/path/to/file", socketPath)

					time.Sleep(20 * time.Millisecond)
					return ctx.Err()
				})
			},
			device:     "foo",
			language:   "bar",
			socketPath: "/path/to/file",
			expectErr:  "context canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			l, h := logger.NewCaptureLogger()
			task := setupTask(t, tt.setupMocks, l, tt.device, tt.language, tt.socketPath)
			ctx, cancel := context.WithCancel(t.Context())
			go func() {
				time.Sleep(10 * time.Millisecond)
				cancel()
			}()

			err := task.Run(ctx)

			require.Error(t, err)
			require.ErrorContains(t, err, tt.expectErr)
			if tt.expectErrMsg != "" {
				requireLogError(t, h, slog.LevelError, "Recognition Task failed", "something happened")
			}
		})
	}
}

func setupTask(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, a *audio.AudioMock, dl *audio.DeviceListenerMock, b *components.BurnerComponentMock, r *components.RecognizerComponentMock, s *components.SocketComponentMock, o *components.OutputComponentMock),
	l *slog.Logger,
	device, language, socketPath string,
) *recognition.Task {
	t.Helper()

	mc := minimock.NewController(t)
	a := audio.NewAudioMock(mc)
	dl := audio.NewDeviceListenerMock(mc)
	b := components.NewBurnerComponentMock(mc)
	r := components.NewRecognizerComponentMock(mc)
	s := components.NewSocketComponentMock(mc)
	o := components.NewOutputComponentMock(mc)

	if setupMocks != nil {
		setupMocks(mc, a, dl, b, r, s, o)
	}
	if l == nil {
		l = logger.NilLogger
	}

	return recognition.NewTask(l, a, b, r, s, o, device, language, socketPath)
}

func requireLogError(t *testing.T, h *logger.CaptureHandler, level slog.Level, msg, errorMsg string) {
	t.Helper()

	logEntry, ok := h.GetLog(msg)
	require.Truef(t, ok, "cannot find the log entry")
	require.Equal(t, level, logEntry.Level)

	errAttr, ok := logEntry.GetAttr("error")
	require.Truef(t, ok, "cannot find the error attribute")
	require.Implements(t, (*error)(nil), errAttr)
	require.ErrorContains(t, errAttr.(error), errorMsg)
}

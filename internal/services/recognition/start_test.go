package recognition_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/background"
	"github.com/roboloop/live2text/internal/services/audio"
	"github.com/roboloop/live2text/internal/services/recognition"
	"github.com/roboloop/live2text/internal/services/recognition/components"
)

func TestStart(t *testing.T) {
	t.Parallel()

	t.Run("device is busy", func(t *testing.T) {
		t.Parallel()

		r, tm := setupRecognition(t, nil, nil)
		defer tm.Wait()
		_ = tm.Go("foo", background.RunnableTaskFunc(func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}))

		_, _, err := r.Start(t.Context(), "foo", "bar")
		require.Error(t, err)
		require.ErrorContains(t, err, "device is busy")
	})

	t.Run("task created", func(t *testing.T) {
		t.Parallel()

		r, tm := setupRecognition(
			t,
			func(mc *minimock.Controller, a *audio.AudioMock, tf *recognition.TaskFactoryMock) {
				tf.NewTaskMock.
					Expect("foo", "bar").
					Return(setupTask(t, func(mc *minimock.Controller, a *audio.AudioMock, dl *audio.DeviceListenerMock, b *components.BurnerComponentMock, r *components.RecognizerComponentMock, s *components.SocketComponentMock, o *components.OutputComponentMock) {
						a.ListenDeviceMock.Return(nil, errors.New("something happened"))
					}, nil, "foo", "bar", "/path/to/file"))
			},
			nil,
		)
		defer tm.Wait()

		id, socketPath, err := r.Start(t.Context(), "foo", "bar")
		require.Equal(t, "foo", id)
		require.Equal(t, "/path/to/file", socketPath)
		require.NoError(t, err)
	})
}

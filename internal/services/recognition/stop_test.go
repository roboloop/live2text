package recognition_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/background"
	"github.com/roboloop/live2text/internal/services/recognition"
)

func TestStop(t *testing.T) {
	t.Parallel()

	t.Run("no device busy", func(t *testing.T) {
		t.Parallel()

		r, tm := setupRecognition(t, nil, nil)
		defer tm.Wait()

		err := r.Stop(t.Context(), "foo")

		require.Error(t, err)
		require.ErrorIs(t, err, recognition.ErrNoDeviceBusy)
	})

	t.Run("task canceled", func(t *testing.T) {
		t.Parallel()

		r, tm := setupRecognition(t, nil, nil)
		defer tm.Wait()
		_ = tm.Go("foo", background.RunnableTaskFunc(func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}))

		err := r.Stop(t.Context(), "foo")
		require.NoError(t, err)
	})
}

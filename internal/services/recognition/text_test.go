package recognition_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/background"
	"github.com/roboloop/live2text/internal/services/recognition"
)

func TestText(t *testing.T) {
	t.Parallel()

	t.Run("no task found", func(t *testing.T) {
		t.Parallel()

		r, _ := setupRecognition(t, nil, nil)

		_, err := r.Text(t.Context(), "foo")

		require.Error(t, err)
		require.Error(t, recognition.ErrNoTaskFound)
	})

	t.Run("not a recognize Task", func(t *testing.T) {
		t.Parallel()

		r, tm := setupRecognition(t, nil, nil)
		defer tm.Wait()
		_ = tm.Go("foo", background.RunnableTaskFunc(func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}))

		_, err := r.Text(t.Context(), "foo")

		require.Error(t, err)
		require.Error(t, recognition.ErrNotRecognizeTask)
	})
}

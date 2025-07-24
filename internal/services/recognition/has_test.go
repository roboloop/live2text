package recognition_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"live2text/internal/background"
)

func TestHas(t *testing.T) {
	t.Parallel()

	t.Run("Task exists", func(t *testing.T) {
		t.Parallel()

		r, tm := setupRecognition(t, nil, nil)
		defer tm.Wait()
		_ = tm.Go("foo", background.RunnableTaskFunc(func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}))

		require.True(t, r.Has("foo"))
		require.False(t, r.Has("bar"))
	})
}

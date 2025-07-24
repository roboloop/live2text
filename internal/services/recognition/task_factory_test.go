package recognition_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"live2text/internal/services/recognition"
)

func TestNewTask(t *testing.T) {
	t.Parallel()

	t.Run("task created", func(t *testing.T) {
		t.Parallel()

		tf := setupTaskFactory(t, nil, nil)
		task := tf.NewTask("foo", "bar")
		require.IsType(t, &recognition.Task{}, task)
	})
}

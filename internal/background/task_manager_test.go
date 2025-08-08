package background_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/background"
	"github.com/roboloop/live2text/internal/utils/logger"
)

type dummyTask struct{}

func (t *dummyTask) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

type emptyTask struct{}

func (t *emptyTask) Run(_ context.Context) error {
	return nil
}

func TestTaskManager(t *testing.T) {
	t.Parallel()

	t.Run("task canceled by context", func(t *testing.T) {
		t.Parallel()

		tm, cancel := newTaskManager(t.Context())
		err := tm.Go("foo", &dummyTask{})
		require.NoError(t, err)
		cancel()
		tm.Wait()
	})

	t.Run("task canceled by call", func(t *testing.T) {
		t.Parallel()

		tm, _ := newTaskManager(t.Context())
		err := tm.Go("foo", &dummyTask{})
		require.NoError(t, err)

		ok := tm.Cancel("foo")
		require.True(t, ok)

		ok = tm.Cancel("foo")
		require.False(t, ok)

		tm.Wait()
	})

	t.Run("wait until all tasks are finished", func(t *testing.T) {
		t.Parallel()

		tm, _ := newTaskManager(t.Context())
		err := tm.Go("foo", &emptyTask{})
		require.NoError(t, err)
		tm.Wait()
	})

	t.Run("cannot run a task with the same id", func(t *testing.T) {
		t.Parallel()

		tm, cancel := newTaskManager(t.Context())
		err := tm.Go("foo", &emptyTask{})
		require.NoError(t, err)

		err = tm.Go("foo", &dummyTask{})
		require.ErrorIs(t, err, background.ErrTaskIsRunning)
		cancel()
		tm.Wait()
	})

	t.Run("total running tasks", func(t *testing.T) {
		t.Parallel()

		tm, cancel := newTaskManager(t.Context())
		require.NoError(t, tm.Go("foo", &dummyTask{}))
		require.NoError(t, tm.Go("bar", &dummyTask{}))

		total := tm.TotalRunningTasks()
		require.Equal(t, 2, total)
		cancel()
		tm.Wait()
	})

	t.Run("get task by id", func(t *testing.T) {
		t.Parallel()

		tm, cancel := newTaskManager(t.Context())
		task := &dummyTask{}
		require.NoError(t, tm.Go("foo", task))

		gotTask := tm.Get("foo")
		require.Equal(t, task, gotTask)

		noTask := tm.Get("bar")
		require.Nil(t, noTask)

		cancel()
		tm.Wait()
	})
}

func newTaskManager(ctx context.Context) (*background.TaskManager, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	tm := background.NewTaskManager(ctx, logger.NilLogger)

	return tm, cancel
}

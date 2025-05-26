package background_test

import (
	"context"
	"errors"
	"live2text/internal/background"
	"slices"
	"testing"
)

type dummyTask struct {
}

func (t *dummyTask) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

type emptyTask struct {
}

func (t *emptyTask) Run(_ context.Context) error {
	return nil
}

func TestTaskManager(t *testing.T) {
	t.Run("Task canceled by context", func(t *testing.T) {
		tm, cancel := newTaskManager(t.Context())
		assertNoError(t, tm.Go("foo", &dummyTask{}))
		cancel()
		tm.Wait()
	})

	t.Run("Task canceled by call", func(t *testing.T) {
		tm, _ := newTaskManager(t.Context())
		assertNoError(t, tm.Go("foo", &dummyTask{}))
		if ok := tm.Cancel("foo"); !ok {
			t.Fatalf("Cancel() returned mismatch value: got %v, expected %v", ok, true)
		}
		// double cancel
		if ok := tm.Cancel("foo"); ok {
			t.Fatalf("Cancel() returned mismatch value: got %v, expected %v", ok, false)
		}

		tm.Wait()
	})

	t.Run("Wait until all tasks are finished", func(t *testing.T) {
		tm, _ := newTaskManager(t.Context())
		assertNoError(t, tm.Go("foo", &emptyTask{}))
		tm.Wait()
	})

	t.Run("Cannot run a task with the same name", func(t *testing.T) {
		tm, cancel := newTaskManager(t.Context())
		assertNoError(t, tm.Go("foo", &dummyTask{}))
		err := tm.Go("foo", &dummyTask{})
		if !errors.Is(err, background.ErrTaskIsRunning) {
			t.Fatalf("Go() returned unexpected error: got %v, expected %v", err, background.ErrTaskIsRunning)
		}
		cancel()
		tm.Wait()
	})

	t.Run("Task manager status", func(t *testing.T) {
		tm, cancel := newTaskManager(t.Context())
		assertNoError(t, tm.Go("foo", &dummyTask{}))
		assertNoError(t, tm.Go("bar", &dummyTask{}))
		status := tm.Status()
		if status.TotalTasks != 2 {
			t.Fatalf("TaskManagerStatus() mismatch: got %v, expected %v", status.TotalTasks, 2)
		}
		expected := []string{"bar", "foo"}
		slices.Sort(status.TaskNames)
		if !slices.Equal(status.TaskNames, expected) {
			t.Fatalf("TaskManagerStatus() mismatch: got %v, expected %v", status.TaskNames, expected)
		}

		cancel()
		tm.Wait()
	})
}

func newTaskManager(ctx context.Context) (*background.TaskManager, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	tm := background.NewTaskManager(ctx)

	return tm, cancel
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

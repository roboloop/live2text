package background

import (
	"context"
	"errors"
	"maps"
	"sync"
)

var ErrTaskIsRunning = errors.New("task is already running")

type Status int

const (
	statusPreparing Status = iota
	statusRunning
	statusFinished
)

type RunnableTask interface {
	Run(ctx context.Context) error
}

type task struct {
	task   RunnableTask
	cancel context.CancelFunc
	status Status
	err    error
}

type TaskManager struct {
	ctx context.Context
	mu  sync.Mutex
	wg  sync.WaitGroup

	tasks map[string]*task
}

type TaskManagerStatus struct {
	TotalTasks int
	TaskNames  []string
}

func NewTaskManager(ctx context.Context) *TaskManager {
	return &TaskManager{ctx: ctx, tasks: map[string]*task{}}
}

func (tm *TaskManager) Wait() {
	tm.wg.Wait()
}

func (tm *TaskManager) Go(name string, runnableTask RunnableTask) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, ok := tm.tasks[name]; ok {
		return ErrTaskIsRunning
	}

	taskCtx, cancel := context.WithCancel(tm.ctx)
	t := &task{
		task:   runnableTask,
		cancel: cancel,
		status: statusPreparing,
		err:    nil,
	}
	tm.tasks[name] = t

	tm.wg.Add(1)
	go func() {
		defer tm.wg.Done()

		t.status = statusRunning
		t.err = runnableTask.Run(taskCtx)

		t.status = statusFinished
	}()

	return nil
}

func (tm *TaskManager) Cancel(name string) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, ok := tm.tasks[name]
	if !ok {
		return false
	}
	t.cancel()
	delete(tm.tasks, name)

	return true
}

func (tm *TaskManager) Status() TaskManagerStatus {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	total := len(tm.tasks)
	names := make([]string, 0, total)
	for name := range maps.Keys(tm.tasks) {
		names = append(names, name)
	}

	return TaskManagerStatus{total, names}
}

func (tm *TaskManager) Get(name string) RunnableTask {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, ok := tm.tasks[name]
	if !ok {
		return nil
	}

	return t.task
}

func (tm *TaskManager) Has(name string) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	_, ok := tm.tasks[name]
	return ok
}

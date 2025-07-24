package background

import (
	"context"
	"errors"
	"log/slog"
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

type RunnableTaskFunc func(ctx context.Context) error

func (r RunnableTaskFunc) Run(ctx context.Context) error {
	return r(ctx)
}

type task struct {
	task   RunnableTask
	cancel context.CancelFunc
	status Status
	err    error
}

type TaskManager struct {
	ctx context.Context
	mu  sync.RWMutex
	wg  sync.WaitGroup

	tasks  map[string]*task
	logger *slog.Logger
}

func NewTaskManager(ctx context.Context, logger *slog.Logger) *TaskManager {
	// TODO: add logging
	return &TaskManager{
		ctx:    ctx,
		tasks:  map[string]*task{},
		logger: logger,
	}
}

func (tm *TaskManager) Wait() {
	tm.wg.Wait()
}

func (tm *TaskManager) Go(id string, runnableTask RunnableTask) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, ok := tm.tasks[id]; ok {
		return ErrTaskIsRunning
	}

	taskCtx, cancel := context.WithCancel(tm.ctx)
	t := &task{
		task:   runnableTask,
		cancel: cancel,
		status: statusPreparing,
		err:    nil,
	}
	tm.tasks[id] = t

	tm.wg.Add(1)
	go func() {
		defer tm.wg.Done()

		t.status = statusRunning
		t.err = runnableTask.Run(taskCtx)

		t.status = statusFinished
	}()

	return nil
}

func (tm *TaskManager) Cancel(id string) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, ok := tm.tasks[id]
	if !ok {
		return false
	}
	t.cancel()
	delete(tm.tasks, id)

	return true
}

func (tm *TaskManager) TotalRunningTasks() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	return len(tm.tasks)
}

func (tm *TaskManager) Get(id string) RunnableTask {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, ok := tm.tasks[id]
	if !ok {
		return nil
	}

	return t.task
}

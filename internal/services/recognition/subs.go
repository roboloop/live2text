package recognition

import (
	"context"
	"errors"
)

var ErrNoTask = errors.New("no task found")

func (r *recognition) Subs(_ context.Context, id string) (string, error) {
	task := r.taskManager.Get(id)
	if task == nil {
		return "", ErrNoTask
	}

	recognizeTask, ok := task.(*RecognizeTask)
	if !ok {
		return "", errors.New("task is not recognize task")
	}

	return recognizeTask.Subs(), nil
}

package recognition

import (
	"context"
	"errors"
)

var NoTaskError = errors.New("no task found")

func (r *recognition) Subs(_ context.Context, device string) (string, error) {
	task := r.taskManager.Get(device)
	if task == nil {
		return "", NoTaskError
	}

	recognizeTask, ok := task.(*RecognizeTask)
	if !ok {
		return "", errors.New("task is not recognize task")
	}

	return recognizeTask.Subs(), nil
}

package recognition

import (
	"context"
	"errors"
)

var (
	ErrNoTaskFound      = errors.New("no task found")
	ErrNotRecognizeTask = errors.New("not a recognize Task")
)

func (r *recognition) Text(_ context.Context, id string) (string, error) {
	t := r.taskManager.Get(id)
	if t == nil {
		return "", ErrNoTaskFound
	}

	recognizeTask, ok := t.(*Task)
	if !ok {
		return "", ErrNotRecognizeTask
	}

	return recognizeTask.Text(), nil
}

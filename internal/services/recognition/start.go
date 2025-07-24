package recognition

import (
	"context"
	"errors"
	"fmt"
)

var ErrDeviceIsBusy = errors.New("device is busy")

func (r *recognition) Start(_ context.Context, device, language string) (string, string, error) {
	if r.taskManager.Get(device) != nil {
		return "", "", ErrDeviceIsBusy
	}

	id := device
	task := r.taskFactory.NewTask(device, language)
	if err := r.taskManager.Go(id, task); err != nil {
		return "", "", fmt.Errorf("cannot run the Task: %w", err)
	}

	return id, task.SocketPath(), nil
}

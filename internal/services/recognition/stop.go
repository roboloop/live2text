package recognition

import (
	"context"
	"errors"
)

var ErrNoDeviceBusy = errors.New("no device busy")

func (r *recognition) Stop(_ context.Context, id string) error {
	if r.taskManager.Get(id) == nil {
		return ErrNoDeviceBusy
	}

	r.taskManager.Cancel(id)

	return nil
}

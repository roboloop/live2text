package recognition

import (
	"context"
	"errors"
)

var ErrNoDeviceBusy = errors.New("no device busy")

func (r *recognition) Stop(_ context.Context, device string) error {
	if !r.taskManager.Has(device) {
		return ErrNoDeviceBusy
	}

	r.taskManager.Cancel(device)

	return nil
}

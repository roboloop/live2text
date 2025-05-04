package recognition

import (
	"context"
	"errors"
)

var NoDeviceBusyError = errors.New("no device busy")

func (r *recognition) Stop(ctx context.Context, device string) error {
	if !r.taskManager.Has(device) {
		return NoDeviceBusyError
	}

	r.taskManager.Cancel(device)

	return nil
}

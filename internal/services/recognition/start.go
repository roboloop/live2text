package recognition

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
)

var ErrDeviceIsBusy = errors.New("device is busy")

func (r *recognition) Start(_ context.Context, device string, language string) (string, string, error) {
	if r.taskManager.Has(device) {
		return "", "", ErrDeviceIsBusy
	}

	id := device
	path := fmt.Sprintf("%s-%d.sock", "recognizer", rand.Uint64()) //nolint:gosec
	socketPath := filepath.Join(os.TempDir(), path)
	task := NewRecognizeTask(
		r.logger,
		r.metrics,
		r.audio,
		r.burner,
		r.socketManager,
		r.speechClient,
		device,
		language,
		socketPath,
	)
	if err := r.taskManager.Go(id, task); err != nil {
		return "", "", fmt.Errorf("cannot run the task: %w", err)
	}

	return id, socketPath, nil
}

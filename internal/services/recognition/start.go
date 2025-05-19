package recognition

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

var DeviceIsBusyError = errors.New("device is busy")

func (r *recognition) Start(_ context.Context, device string, language string) (string, string, error) {
	if r.taskManager.Has(device) {
		return "", "", DeviceIsBusyError
	}

	id := device
	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d.sock", "recognizer", rand.Uint64()))
	task := NewRecognizeTask(r.logger, r.metrics, r.audio, r.burner, r.socketManager, r.speechClient, device, language, socketPath)
	if err := r.taskManager.Go(id, task); err != nil {
		return "", "", fmt.Errorf("cannot run the task: %w", err)
	}

	return id, socketPath, nil
}

package recognition

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"path/filepath"

	"github.com/roboloop/live2text/internal/services/audio"
	"github.com/roboloop/live2text/internal/services/recognition/components"
)

type taskFactory struct {
	logger *slog.Logger
	audio  audio.Audio

	burner     components.BurnerComponent
	recognizer components.RecognizerComponent
	socket     components.SocketComponent
	output     components.OutputComponent
}

func NewTaskFactory(
	logger *slog.Logger,
	audio audio.Audio,
	burner components.BurnerComponent,
	recognizer components.RecognizerComponent,
	socket components.SocketComponent,
	output components.OutputComponent,
) TaskFactory {
	return &taskFactory{
		logger:     logger,
		audio:      audio,
		burner:     burner,
		recognizer: recognizer,
		socket:     socket,
		output:     output,
	}
}

func (tf *taskFactory) NewTask(device, language string) *Task {
	path := fmt.Sprintf("%d.sock", rand.Uint64()) //nolint:gosec
	socketPath := filepath.Join(os.TempDir(), path)

	return NewTask(
		tf.logger,
		tf.audio,
		tf.burner,
		tf.recognizer,
		tf.socket,
		tf.output,

		device,
		language,
		socketPath,
	)
}

package recognition

import (
	"live2text/internal/background"
	"live2text/internal/services/audio"
	"live2text/internal/services/burner"
	"live2text/internal/services/metrics"
	speechwrapper "live2text/internal/services/speech_wrapper"
	"log/slog"
)

type recognition struct {
	logger        *slog.Logger
	metrics       metrics.Metrics
	audio         audio.Audio
	burner        burner.Burner
	speechClient  speechwrapper.Client
	taskManager   *background.TaskManager
	socketManager *background.SocketManager
}

func NewRecognition(
	logger *slog.Logger,
	metrics metrics.Metrics,
	audio audio.Audio,
	burner burner.Burner,
	speechClient speechwrapper.Client,
	taskManager *background.TaskManager,
	socketManager *background.SocketManager,
) Recognition {
	return &recognition{
		logger.With("service", "Recognition"),
		metrics,
		audio,
		burner,
		speechClient,
		taskManager,
		socketManager,
	}
}

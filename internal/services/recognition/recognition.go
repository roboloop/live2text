package recognition

import (
	"log/slog"

	"github.com/roboloop/live2text/internal/background"
)

type recognition struct {
	logger      *slog.Logger
	taskManager *background.TaskManager
	taskFactory TaskFactory
}

func NewRecognition(
	logger *slog.Logger,
	taskManager *background.TaskManager,
	taskFactory TaskFactory,
) Recognition {
	return &recognition{
		logger:      logger,
		taskManager: taskManager,
		taskFactory: taskFactory,
	}
}

package audio

import (
	"log/slog"

	audiowrapper "github.com/roboloop/live2text/internal/services/audio_wrapper"
	"github.com/roboloop/live2text/internal/services/metrics"
)

type audio struct {
	logger        *slog.Logger
	metrics       metrics.Metrics
	externalAudio audiowrapper.Audio
}

func NewAudio(logger *slog.Logger, metrics metrics.Metrics, externalAudio audiowrapper.Audio) Audio {
	return &audio{logger: logger, metrics: metrics, externalAudio: externalAudio}
}

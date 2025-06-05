package audio

import (
	"log/slog"

	audiowrapper "live2text/internal/services/audio_wrapper"
	"live2text/internal/services/metrics"
)

type audio struct {
	logger        *slog.Logger
	metrics       metrics.Metrics
	externalAudio audiowrapper.Audio
}

func NewAudio(logger *slog.Logger, metrics metrics.Metrics, externalAudio audiowrapper.Audio) Audio {
	return &audio{logger.With("service", "Audio"), metrics, externalAudio}
}

package audio

import (
	"live2text/internal/services/audio_wrapper"
	"live2text/internal/services/metrics"
	"log/slog"
)

type audio struct {
	logger        *slog.Logger
	metrics       metrics.Metrics
	externalAudio audio_wrapper.Audio
}

func NewAudio(logger *slog.Logger, metrics metrics.Metrics, externalAudio audio_wrapper.Audio) Audio {
	return &audio{logger, metrics, externalAudio}
}

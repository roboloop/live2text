package services

import (
	"live2text/internal/services/audio"
	"live2text/internal/services/audio_wrapper"
	"live2text/internal/services/btt"
	"live2text/internal/services/burner"
	"live2text/internal/services/metrics"
	"live2text/internal/services/recognition"
)

type Services interface {
	Audio() audio.Audio
	AudioWrapper() audio_wrapper.Audio
	Burner() burner.Burner
	Recognition() recognition.Recognition
	Metrics() metrics.Metrics
	Btt() btt.Btt
}

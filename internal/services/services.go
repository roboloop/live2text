package services

import (
	"live2text/internal/services/audio"
	"live2text/internal/services/audio_wrapper"
	"live2text/internal/services/burner"
	"live2text/internal/services/metrics"
	"live2text/internal/services/recognition"
)

type services struct {
	audio        audio.Audio
	audioWrapper audio_wrapper.Audio
	burner       burner.Burner
	recognition  recognition.Recognition
	metrics      metrics.Metrics
}

func NewServices(
	audio audio.Audio,
	audioWrapper audio_wrapper.Audio,
	burner burner.Burner,
	recognition recognition.Recognition,
	metrics metrics.Metrics,
) Services {
	return &services{audio, audioWrapper, burner, recognition, metrics}
}

func (s *services) Audio() audio.Audio {
	return s.audio
}

func (s *services) AudioWrapper() audio_wrapper.Audio {
	return s.audioWrapper
}

func (s *services) Burner() burner.Burner {
	return s.burner
}

func (s *services) Recognition() recognition.Recognition {
	return s.recognition
}

func (s *services) Metrics() metrics.Metrics {
	return s.metrics
}

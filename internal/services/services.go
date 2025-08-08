package services

import (
	"github.com/roboloop/live2text/internal/services/audio"
	audiowrapper "github.com/roboloop/live2text/internal/services/audio_wrapper"
	"github.com/roboloop/live2text/internal/services/btt"
	"github.com/roboloop/live2text/internal/services/burner"
	"github.com/roboloop/live2text/internal/services/metrics"
	"github.com/roboloop/live2text/internal/services/recognition"
)

type services struct {
	audio        audio.Audio
	audioWrapper audiowrapper.Audio
	burner       burner.Burner
	recognition  recognition.Recognition
	metrics      metrics.Metrics
	btt          btt.Btt
}

func NewServices(
	audio audio.Audio,
	audioWrapper audiowrapper.Audio,
	burner burner.Burner,
	recognition recognition.Recognition,
	metrics metrics.Metrics,
	btt btt.Btt,
) Services {
	return &services{audio, audioWrapper, burner, recognition, metrics, btt}
}

func (s *services) Audio() audio.Audio {
	return s.audio
}

func (s *services) AudioWrapper() audiowrapper.Audio {
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

func (s *services) Btt() btt.Btt {
	return s.btt
}

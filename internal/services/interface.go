package services

import (
	"github.com/roboloop/live2text/internal/services/audio"
	audiowrapper "github.com/roboloop/live2text/internal/services/audio_wrapper"
	"github.com/roboloop/live2text/internal/services/btt"
	"github.com/roboloop/live2text/internal/services/burner"
	"github.com/roboloop/live2text/internal/services/metrics"
	"github.com/roboloop/live2text/internal/services/recognition"
)

//go:generate minimock -g -i Services -s _mock.go -o .

type Services interface {
	Audio() audio.Audio
	AudioWrapper() audiowrapper.Audio
	Burner() burner.Burner
	Recognition() recognition.Recognition
	Metrics() metrics.Metrics
	Btt() btt.Btt
}

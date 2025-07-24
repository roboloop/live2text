package services

import (
	"live2text/internal/services/audio"
	audiowrapper "live2text/internal/services/audio_wrapper"
	"live2text/internal/services/btt"
	"live2text/internal/services/burner"
	"live2text/internal/services/metrics"
	"live2text/internal/services/recognition"
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

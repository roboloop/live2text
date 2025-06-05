package btt

import (
	"context"
	"log/slog"

	"live2text/internal/config"
	"live2text/internal/services/audio"
	"live2text/internal/services/btt/http"
	"live2text/internal/services/btt/tmpl"
	"live2text/internal/services/recognition"
)

const (
	appName                    = "Live2Text"
	appTitle                   = "App"
	settingsTitle              = "Settings"
	deviceGroupTitle           = "Device"
	languageGroupTitle         = "Language"
	floatingStateGroupTitle    = "Floating State"
	streamingTextTitle         = "Streaming Text"
	selectedLanguageTitle      = "Selected Language"
	selectedDeviceTitle        = "Selected Device"
	selectedFloatingStateTitle = "Selected Floating State"

	defaultInterval = 0.25
)

type btt struct {
	logger      *slog.Logger
	audio       audio.Audio
	recognition recognition.Recognition
	httpClient  http.Client
	renderer    *tmpl.Renderer

	appAddress string
	bttAddress string
	languages  []string

	interval float64
}

func NewBtt(
	logger *slog.Logger,
	audio audio.Audio,
	recognition recognition.Recognition,
	httpClient http.Client,
	cfg *config.Config,
) Btt {
	debug := logger.Handler().Enabled(context.Background(), slog.LevelDebug)

	return &btt{
		logger:      logger.With("service", "Btt"),
		audio:       audio,
		recognition: recognition,
		renderer:    tmpl.NewRenderer(appName, debug),
		httpClient:  httpClient,

		appAddress: cfg.AppAddress,
		bttAddress: cfg.BttAddress,
		languages:  cfg.Languages,

		interval: defaultInterval,
	}
}

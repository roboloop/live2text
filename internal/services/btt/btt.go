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

// Must be uniq names.
const (
	appName                    = "Live2Text"
	appTitle                   = "App"
	settingsTitle              = "Settings"
	cleanViewTitle             = "Clean view"
	cleanViewAppTitle          = "Clean view App"
	deviceGroupTitle           = "Device"
	languageGroupTitle         = "Language"
	viewModeTitle              = "View Mode"
	floatingStateGroupTitle    = "Floating State"
	metricsGroupTitle          = "Metrics"
	streamingTextTitle         = "Streaming Text"
	selectedLanguageTitle      = "Selected Language"
	selectedDeviceTitle        = "Selected Device"
	selectedFloatingStateTitle = "Selected Floating State"
	selectedViewModeTitle      = "Selected View Mode"

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

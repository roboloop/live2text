package btt

import (
	"context"
	"live2text/internal/config"
	"live2text/internal/services/audio"
	"live2text/internal/services/btt/exec"
	"live2text/internal/services/btt/http"
	"live2text/internal/services/btt/tmpl"
	"live2text/internal/services/recognition"
	"log/slog"
)

const bttName = "BetterTouchTool"
const appName = "live2text"

const (
	settingsTitle         = appName + " Settings"
	appTitle              = appName + " App"
	deviceGroupTitle      = "Device"
	languageGroupTitle    = "Language"
	selectedLanguageTitle = "Selected Language"
	selectedDeviceTitle   = "Selected Device"

	defaultInterval = 0.25
)

type btt struct {
	logger      *slog.Logger
	audio       audio.Audio
	recognition recognition.Recognition
	httpClient  http.Client
	execClient  exec.Client
	renderer    *tmpl.Renderer

	appAddress string
	bttAddress string
	languages  []string

	interval float64
}

func NewBtt(logger *slog.Logger, audio audio.Audio, recognition recognition.Recognition, httpClient http.Client, execClient exec.Client, cfg *config.Config) Btt {
	debug := logger.Handler().Enabled(context.Background(), slog.LevelDebug)

	return &btt{
		logger:      logger.With("service", "Btt"),
		audio:       audio,
		recognition: recognition,
		renderer:    tmpl.NewRenderer(appName, debug),
		httpClient:  httpClient,
		execClient:  execClient,

		appAddress: cfg.AppAddress,
		bttAddress: cfg.BttAddress,
		languages:  cfg.Languages,

		interval: defaultInterval,
	}
}

package api

import (
	"log/slog"
	"net/http"
	"time"

	"live2text/internal/api/btt"
	"live2text/internal/api/core"
	"live2text/internal/api/middleware"
	"live2text/internal/env"
	"live2text/internal/services"
)

type Server struct {
	core *core.Server
	btt  *btt.Server
}

func NewHandler(logger *slog.Logger, services services.Services) http.Handler {
	server := &Server{
		core: core.NewServer(logger, services),
		btt:  btt.NewServer(logger, services),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", server.core.Health)
	mux.HandleFunc("GET /api/devices", server.core.Devices)
	mux.HandleFunc("POST /api/start", server.core.Start)
	mux.HandleFunc("POST /api/stop", server.core.Stop)
	mux.HandleFunc("POST /api/text", server.core.Text)
	mux.HandleFunc("GET /metrics", server.core.Metrics)

	mux.HandleFunc("GET /api/btt/selected-device", server.btt.SelectedDevice)
	mux.HandleFunc("GET /api/btt/selected-language", server.btt.SelectedLanguage)
	mux.HandleFunc("GET /api/btt/selected-view-mode", server.btt.SelectedViewMode)
	mux.HandleFunc("GET /api/btt/selected-floating", server.btt.SelectedFloating)
	mux.HandleFunc("GET /api/btt/selected-clipboard", server.btt.SelectedClipboard)
	mux.HandleFunc("POST /api/btt/select-device", server.btt.SelectDevice)
	mux.HandleFunc("POST /api/btt/select-language", server.btt.SelectLanguage)
	mux.HandleFunc("POST /api/btt/select-view-mode", server.btt.SelectViewMode)
	mux.HandleFunc("POST /api/btt/select-clipboard", server.btt.SelectClipboard)

	mux.HandleFunc("POST /api/btt/load-devices", server.btt.LoadDevices)
	mux.HandleFunc("POST /api/btt/toggle-listening", server.btt.ToggleListening)
	mux.HandleFunc("GET /api/btt/is-running", server.btt.IsRunning)

	mux.HandleFunc("GET /btt/floating-page", server.btt.FloatingPage)
	mux.HandleFunc("GET /api/btt/text-stream", server.btt.TextStream)
	mux.HandleFunc("GET /api/btt/text", server.btt.Text)
	mux.HandleFunc("GET /api/btt/health", server.btt.Health)

	var handler http.Handler = mux
	handler = middleware.LoggerMiddleware(handler, logger)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.ErrorMiddleware(handler)
	handler = middleware.TimeoutMiddleware(handler, 15*time.Second, []string{"/api/btt/text-stream"}, env.IsDebugMode())
	handler = middleware.BodyCloserMiddleware(handler, logger)

	return handler
}

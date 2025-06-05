package api

import (
	"log/slog"
	"net/http"
	"time"

	"live2text/internal/api/btt"
	"live2text/internal/api/core"
	"live2text/internal/api/middleware"
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
	mux.HandleFunc("GET /api/subs", server.core.Subs)
	mux.HandleFunc("GET /metrics", server.core.Metrics)

	mux.HandleFunc("GET /api/btt/selected-device", server.btt.SelectedDevice)
	mux.HandleFunc("GET /api/btt/selected-language", server.btt.SelectedLanguage)
	mux.HandleFunc("GET /api/btt/selected-floating-state", server.btt.SelectedFloatingState)
	mux.HandleFunc("POST /api/btt/select-device", server.btt.SelectDevice)
	mux.HandleFunc("POST /api/btt/select-language", server.btt.SelectLanguage)
	mux.HandleFunc("POST /api/btt/select-floating-state", server.btt.SelectFloatingState)

	mux.HandleFunc("POST /api/btt/load-devices", server.btt.LoadDevices)
	mux.HandleFunc("POST /api/btt/toggle-listening", server.btt.ToggleListening)

	mux.HandleFunc("GET /btt/floating-page", server.btt.FloatingPage)
	mux.HandleFunc("GET /api/btt/text-stream", server.btt.TextStream)

	var handler http.Handler = mux
	handler = middleware.LoggerMiddleware(handler, logger)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.ErrorMiddleware(handler)
	handler = middleware.TimeoutMiddleware(handler, 15*time.Second, []string{"/api/btt/text-stream"})

	return handler
}

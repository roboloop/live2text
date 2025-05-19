package api

import (
	"live2text/internal/api/middleware"
	"live2text/internal/services"
	"log/slog"
	"net/http"
)

type Server struct {
	logger   *slog.Logger
	services services.Services
}

func NewHandler(logger *slog.Logger, services services.Services) http.Handler {
	server := &Server{
		logger,
		services,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", server.Health)
	mux.HandleFunc("GET /api/devices", server.Devices)
	mux.HandleFunc("POST /api/start", server.Start)
	mux.HandleFunc("POST /api/stop", server.Stop)
	mux.HandleFunc("GET /api/subs", server.Subs)
	mux.HandleFunc("GET /metrics", server.Metrics)

	mux.HandleFunc("GET /api/selected-device", server.SelectedDevice)
	mux.HandleFunc("GET /api/selected-language", server.SelectedLanguage)
	mux.HandleFunc("POST /api/select-device", server.SelectDevice)
	mux.HandleFunc("POST /api/select-language", server.SelectLanguage)

	mux.HandleFunc("POST /api/load-devices", server.LoadDevices)
	mux.HandleFunc("POST /api/toggle-listening", server.ToggleListening)

	var handler http.Handler = mux
	handler = middleware.ErrorMiddleware(handler)
	handler = middleware.LoggerMiddleware(handler, logger)

	return handler
}

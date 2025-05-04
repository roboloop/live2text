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
	mux.HandleFunc("/api/health", server.Health)
	mux.HandleFunc("/api/devices", server.Devices)
	mux.HandleFunc("/api/start", server.Start)
	mux.HandleFunc("/api/stop", server.Stop)
	mux.HandleFunc("/api/subs", server.Subs)
	mux.HandleFunc("/metrics", server.Metrics)

	var handler http.Handler = mux
	handler = middleware.ErrorMiddleware(handler)
	handler = middleware.LoggerMiddleware(handler, logger)

	return handler
}

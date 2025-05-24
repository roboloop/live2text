package btt

import (
	"live2text/internal/services"
	"log/slog"
)

type Server struct {
	logger   *slog.Logger
	services services.Services
}

func NewServer(logger *slog.Logger, services services.Services) *Server {
	return &Server{logger: logger, services: services}
}
